package engine

import (
	"context"
	"database/sql"
	"math/rand/v2"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/plamen-v/tic-tac-toe-models/models"
	"github.com/plamen-v/tic-tac-toe/src/repository"
)

const (
	DefaultBoardTile         byte   = '_'
	DefaultBoard             string = "_________"
	XMark                    byte   = 'X'
	OMark                    byte   = 'O'
	MaxRoomTitleLength       int    = 30
	MaxRoomDescriptionLength int    = 150
)

var (
	PlayerPartOfOtherRoomError      error = models.NewValidationError("player is part of other room")
	TitleRequiredError              error = models.NewValidationError("title is required")
	TitleTooLongError               error = models.NewValidationErrorf("title is too long. Max length is %d", MaxRoomTitleLength)
	DescriptionTooLongError         error = models.NewValidationErrorf("description is too long required. Max length is %d", MaxRoomDescriptionLength)
	FullRoomError                   error = models.NewValidationError("room is full")
	PlayerPartOfTheRoomAsHostError  error = models.NewValidationErrorf("player is already part of the room as host")
	PlayerPartOfTheRoomAsGuestError error = models.NewValidationErrorf("player is already part of the room as guest")
	PlayerNotInRoomError            error = models.NewValidationErrorf("player is not part of the room")
	GameInProgressError             error = models.NewValidationErrorf("can't creat new game. Previous game is in progress.")
	GameCompletedError              error = models.NewValidationError("game is completed")
	PlayerNotInTurnError            error = models.NewValidationError("player not in turn")
	InvalidBoardPositionError       error = models.NewValidationError("invalid position index")
	BoardPositionOcopiedError       error = models.NewValidationError("position ocopied")
)

type GameEngineService interface {
	GetRoom(context.Context, uuid.UUID) (*models.Room, error)
	GetOpenRooms(context.Context) ([]*models.Room, error)
	CreateRoom(context.Context, uuid.UUID, string, string) (uuid.UUID, error)
	PlayerJoinRoom(context.Context, uuid.UUID, uuid.UUID) error
	PlayerLeaveRoom(context.Context, uuid.UUID, uuid.UUID) error
	CreateGame(context.Context, uuid.UUID, uuid.UUID) (uuid.UUID, error)
	GetGameState(context.Context, uuid.UUID, uuid.UUID) (*models.Game, error)
	PlayerMakeMove(context.Context, uuid.UUID, uuid.UUID, int) error
	GetRanking(context.Context) ([]*models.Player, error)
}

type gameEngineServiceImpl struct {
	db                      *sql.DB
	playerRepositoryFactory func(q repository.Querier) repository.PlayerRepository
	gameRepositoryFactory   func(q repository.Querier) repository.GameRepository
	roomRepositoryFactory   func(q repository.Querier) repository.RoomRepository
}

func NewGameEngineService(db *sql.DB,
	playerRepositoryFactory func(q repository.Querier) repository.PlayerRepository,
	gameRepositoryFactory func(q repository.Querier) repository.GameRepository,
	roomRepositoryFactory func(q repository.Querier) repository.RoomRepository) GameEngineService {
	return &gameEngineServiceImpl{
		db:                      db,
		playerRepositoryFactory: playerRepositoryFactory,
		gameRepositoryFactory:   gameRepositoryFactory,
		roomRepositoryFactory:   roomRepositoryFactory,
	}
}

func (g *gameEngineServiceImpl) GetRoom(ctx context.Context, playerID uuid.UUID) (*models.Room, error) {
	return g.roomRepositoryFactory(g.db).GetByPlayerID(ctx, playerID)
}

func (g *gameEngineServiceImpl) GetOpenRooms(ctx context.Context) ([]*models.Room, error) {
	return g.roomRepositoryFactory(g.db).GetList(ctx, models.RoomPhaseOpen)
}

func (g *gameEngineServiceImpl) CreateRoom(ctx context.Context, playerID uuid.UUID, title string, description string) (id uuid.UUID, err error) {
	room := &models.Room{
		Host: models.RoomPlayer{
			ID:             playerID,
			RequestNewGame: true,
		},
		Title:       title,
		Description: description,
		Phase:       models.RoomPhaseOpen,
	}
	roomRepository := g.roomRepositoryFactory(g.db)
	err = g.validateCreateRoom(ctx, roomRepository, room, room.Host.ID)
	if err != nil {
		return uuid.Nil, err
	}

	id, err = roomRepository.Create(ctx, room)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func (g *gameEngineServiceImpl) validateCreateRoom(ctx context.Context, roomRepository repository.RoomRepository, room *models.Room, playerID uuid.UUID) error {
	playerRoom, err := roomRepository.GetByPlayerID(ctx, playerID)
	if err != nil && !models.IsNotFoundError(err) {
		return err
	}

	if playerRoom != nil {
		return PlayerPartOfOtherRoomError
	}

	if room != nil && len(room.Title) == 0 {
		return TitleRequiredError
	}

	if room != nil && len(room.Title) > MaxRoomTitleLength {
		return TitleTooLongError
	}

	if room != nil && len(room.Description) > MaxRoomDescriptionLength {
		return DescriptionTooLongError
	}

	return nil
}

func (g *gameEngineServiceImpl) PlayerJoinRoom(ctx context.Context, roomID uuid.UUID, playerID uuid.UUID) error {
	return withTransaction(ctx, g.db, func(tx *sql.Tx) error {
		roomRepository := g.roomRepositoryFactory(tx)
		room, err := roomRepository.Get(ctx, roomID, true)
		if err != nil {
			return err
		}

		err = g.validatePlayerJoinRoom(ctx, roomRepository, room, playerID)
		if err != nil {
			return err
		}

		room.Guest = &models.RoomPlayer{
			ID:             playerID,
			RequestNewGame: true,
		}

		gameRepository := g.gameRepositoryFactory(tx)
		err = g.createGame(ctx, gameRepository, room)
		if err != nil {
			return err
		}

		room.Phase = models.RoomPhaseFull
		err = roomRepository.Update(ctx, room)
		if err != nil {
			return err
		}

		return nil
	})
}

func (g *gameEngineServiceImpl) validatePlayerJoinRoom(ctx context.Context, roomRepository repository.RoomRepository, room *models.Room, playerID uuid.UUID) error {
	if room.Phase == models.RoomPhaseFull {
		return FullRoomError
	}

	if room.Host.ID == playerID {
		return PlayerPartOfTheRoomAsHostError
	}

	if room.Guest != nil && room.Guest.ID == playerID {
		return PlayerPartOfTheRoomAsGuestError
	}

	if _, err := roomRepository.GetByPlayerID(ctx, playerID); err == nil {
		return PlayerPartOfOtherRoomError
	} else if models.IsNotFoundError(err) {
		return nil
	} else {
		return err
	}
}

func (g *gameEngineServiceImpl) PlayerLeaveRoom(ctx context.Context, roomID uuid.UUID, playerID uuid.UUID) (err error) {
	return withTransaction(ctx, g.db, func(tx *sql.Tx) error {
		roomRepository := g.roomRepositoryFactory(tx)

		room, err := roomRepository.Get(ctx, roomID, true)
		if err != nil {
			return err
		}
		err = g.validatePlayerLeaveRoom(room, playerID)
		if err != nil {
			return err
		}

		playerIsHost := playerID == room.Host.ID
		if room.GameID != nil {
			gameRepository := g.gameRepositoryFactory(tx)
			game, err := gameRepository.Get(ctx, *room.GameID)
			if err != nil {
				return err
			}

			if game.Phase == models.GamePhaseInProgress {
				playerRepository := g.playerRepositoryFactory(tx)
				host, err := playerRepository.Get(ctx, room.Host.ID)
				if err != nil {
					return err
				}

				guest, err := playerRepository.Get(ctx, room.Guest.ID)
				if err != nil {
					return err
				}

				if playerIsHost {
					g.finalizeGameWithWin(game, guest, host)
				} else {
					g.finalizeGameWithWin(game, host, guest)
				}

				err = playerRepository.UpdateStats(ctx, host)
				if err != nil {
					return err
				}

				err = playerRepository.UpdateStats(ctx, guest)
				if err != nil {
					return err
				}

				err = gameRepository.Update(ctx, game)
				if err != nil {
					return err
				}
			}
		}

		emptyRoom := false
		if playerIsHost {
			if room.Guest != nil {
				room.Host = *room.Guest
				room.Guest = nil
			} else {
				emptyRoom = true
			}
		} else {
			room.Guest = nil
		}

		if emptyRoom {
			err = roomRepository.Delete(ctx, room.ID)
			if err != nil {
				return err
			}
		} else {
			err = roomRepository.Update(ctx, room)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (g *gameEngineServiceImpl) validatePlayerLeaveRoom(room *models.Room, playerID uuid.UUID) error {
	if room.Host.ID != playerID &&
		(room.Guest == nil || room.Guest.ID != playerID) {
		return PlayerNotInRoomError
	}

	return nil
}

func (g *gameEngineServiceImpl) CreateGame(ctx context.Context, roomID uuid.UUID, playerID uuid.UUID) (uuid.UUID, error) {
	return withTransactionT(ctx, g.db, func(tx *sql.Tx) (uuid.UUID, error) {
		roomRepository := g.roomRepositoryFactory(tx)
		room, err := roomRepository.Get(ctx, roomID, true)
		if err != nil {
			return uuid.Nil, err
		}

		err = g.validateCreateGame(ctx, room, playerID)
		if err != nil {
			return uuid.Nil, err
		}

		if playerID == room.Host.ID {
			room.Host.RequestNewGame = true
		}

		if room.Guest != nil && room.Guest.ID == playerID {
			room.Guest.RequestNewGame = true
		}

		var gameID uuid.UUID
		gameRepository := g.gameRepositoryFactory(tx)
		if room.Guest != nil {
			if room.Guest.RequestNewGame && room.Host.RequestNewGame {
				err := g.createGame(ctx, gameRepository, room)
				if err != nil {
					return uuid.Nil, err
				}

				gameID = *room.GameID
			}
		}

		err = roomRepository.Update(ctx, room)
		if err != nil {
			return uuid.Nil, err
		}
		return gameID, nil
	})
}

func (g *gameEngineServiceImpl) validateCreateGame(ctx context.Context, room *models.Room, playerID uuid.UUID) error {
	if room.Host.ID != playerID &&
		(room.Guest == nil || room.Guest.ID != playerID) {
		return PlayerNotInRoomError
	}

	return nil
}

func (g *gameEngineServiceImpl) GetGameState(ctx context.Context, roomID uuid.UUID, playerID uuid.UUID) (*models.Game, error) {
	roomRepository := g.roomRepositoryFactory(g.db)
	room, err := roomRepository.Get(ctx, roomID, false)
	if err != nil {
		return nil, err
	}

	err = g.validateGetGameState(room, playerID)
	if err != nil {
		return nil, err
	}

	gameID := uuid.Nil
	if room.GameID != nil {
		gameID = *room.GameID
	}
	gameRepository := g.gameRepositoryFactory(g.db)
	game, err := gameRepository.Get(ctx, gameID)
	if err != nil {
		return nil, err
	}
	return game, nil
}

func (g *gameEngineServiceImpl) validateGetGameState(room *models.Room, playerID uuid.UUID) error {
	if room.Host.ID != playerID &&
		(room.Guest == nil || room.Guest.ID != playerID) {
		return PlayerNotInRoomError
	}

	return nil
}

func (g *gameEngineServiceImpl) PlayerMakeMove(ctx context.Context, roomID uuid.UUID, playerID uuid.UUID, position int) error {
	return withTransaction(ctx, g.db, func(tx *sql.Tx) error {
		roomRepository := g.roomRepositoryFactory(tx)
		room, err := roomRepository.Get(ctx, roomID, true)
		if err != nil {
			return err
		}

		gameRepository := g.gameRepositoryFactory(tx)
		if room.GameID != nil {
			game, err := gameRepository.Get(ctx, *room.GameID)
			if err != nil {
				return err
			}

			err = g.validatePlayerMakeMove(game, playerID, position)
			if err != nil {
				return err
			}

			mark := []byte(game.Host.Mark)[0]
			if playerID != room.Host.ID {
				mark = []byte(game.Guest.Mark)[0]
			}

			boardBytes := []byte(game.Board)
			boardBytes[position-1] = mark
			game.Board = string(boardBytes)

			win := g.inWinState(game)
			if win || !strings.Contains(game.Board, string(DefaultBoardTile)) {
				playerRepository := g.playerRepositoryFactory(tx)
				host, err := playerRepository.Get(ctx, room.Host.ID)
				if err != nil {
					return err
				}

				guest, err := playerRepository.Get(ctx, room.Guest.ID)
				if err != nil {
					return err
				}

				winner := host
				loser := guest
				if game.CurrentPlayerID == guest.ID {
					winner = guest
					loser = host
				}
				if win {
					g.finalizeGameWithWin(game, winner, loser)
				} else {
					g.finalizeGameWithDraw(game, host, guest)
				}

				err = playerRepository.UpdateStats(ctx, host)
				if err != nil {
					return err
				}

				err = playerRepository.UpdateStats(ctx, guest)
				if err != nil {
					return err
				}
			} else {
				if game.CurrentPlayerID == game.Host.ID {
					game.CurrentPlayerID = game.Guest.ID
				} else {
					game.CurrentPlayerID = game.Host.ID
				}
			}

			err = gameRepository.Update(ctx, game)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (g *gameEngineServiceImpl) validatePlayerMakeMove(game *models.Game, playerID uuid.UUID, position int) error {
	if game.Host.ID != playerID && game.Guest.ID != playerID {
		return PlayerNotInRoomError
	}

	if game.Phase == models.GamePhaseCompleted {
		return GameCompletedError
	}

	if game.CurrentPlayerID != playerID {
		return PlayerNotInTurnError
	}

	if position < 1 || position > len(DefaultBoard) {
		return InvalidBoardPositionError
	}
	if game.Board[position-1] != DefaultBoardTile {
		return BoardPositionOcopiedError
	}

	return nil
}

func (g *gameEngineServiceImpl) GetRanking(ctx context.Context) ([]*models.Player, error) {
	playerRepository := g.playerRepositoryFactory(g.db)
	players, err := playerRepository.GetRanking(ctx)
	if err != nil {
		return nil, err
	}

	return players, nil
}

func (g *gameEngineServiceImpl) createGame(ctx context.Context, gameRepository repository.GameRepository, room *models.Room) error {

	game, err := g.initializeGame(ctx, gameRepository, room)
	if err != nil {
		return err
	}
	game.ID, err = gameRepository.Create(ctx, game)
	if err != nil {
		return err
	}
	room.GameID = &game.ID
	room.Host.RequestNewGame = false
	room.Guest.RequestNewGame = false

	return nil
}

func (g *gameEngineServiceImpl) initializeGame(ctx context.Context, gameRepository repository.GameRepository, room *models.Room) (*models.Game, error) {
	marks := []byte{XMark, OMark}
	rand.Shuffle(len(marks), func(i, j int) {
		marks[i], marks[j] = marks[j], marks[i]
	})
	playerIDs := []uuid.UUID{room.Host.ID, room.Guest.ID}

	game := &models.Game{
		Host:            models.GamePlayer{ID: room.Host.ID, Mark: string(marks[0])},
		Guest:           models.GamePlayer{ID: room.Guest.ID, Mark: string(marks[1])},
		CurrentPlayerID: playerIDs[rand.IntN(2)],
		Board:           DefaultBoard,
		Phase:           models.GamePhaseInProgress,
	}

	if room.GameID != nil {
		prevGame, err := gameRepository.Get(ctx, *room.GameID)
		if err != nil {
			return nil, err
		}

		if prevGame.Phase == models.GamePhaseInProgress {
			return nil, GameInProgressError
		}

		if prevGame.Host.ID == room.Host.ID && prevGame.Guest.ID == room.Guest.ID {
			game.Host.Mark = prevGame.Host.Mark
			game.Guest.Mark = prevGame.Guest.Mark

			game.CurrentPlayerID = room.Host.ID
			if prevGame.CurrentPlayerID == room.Host.ID {
				game.CurrentPlayerID = room.Guest.ID
			}
		}
	}

	return game, nil
}

func (g *gameEngineServiceImpl) inWinState(game *models.Game) bool {

	if (game.Board[0] != DefaultBoardTile && game.Board[0] == game.Board[1] && game.Board[0] == game.Board[2]) ||
		(game.Board[3] != DefaultBoardTile && game.Board[3] == game.Board[4] && game.Board[3] == game.Board[5]) ||
		(game.Board[6] != DefaultBoardTile && game.Board[6] == game.Board[7] && game.Board[6] == game.Board[8]) ||

		(game.Board[0] != DefaultBoardTile && game.Board[0] == game.Board[3] && game.Board[0] == game.Board[6]) ||
		(game.Board[1] != DefaultBoardTile && game.Board[1] == game.Board[4] && game.Board[1] == game.Board[7]) ||
		(game.Board[2] != DefaultBoardTile && game.Board[2] == game.Board[5] && game.Board[2] == game.Board[8]) ||

		(game.Board[0] != DefaultBoardTile && game.Board[0] == game.Board[4] && game.Board[0] == game.Board[8]) ||
		(game.Board[2] != DefaultBoardTile && game.Board[2] == game.Board[4] && game.Board[2] == game.Board[6]) {
		return true
	}

	return false
}

func (g *gameEngineServiceImpl) finalizeGameWithWin(game *models.Game, winner *models.Player, loser *models.Player) {
	game.Phase = models.GamePhaseCompleted
	game.WinnerID = &winner.ID
	winner.Stats.Wins++
	loser.Stats.Losses++
}

func (g *gameEngineServiceImpl) finalizeGameWithDraw(game *models.Game, host *models.Player, guest *models.Player) {
	game.Phase = models.GamePhaseCompleted
	host.Stats.Draws++
	guest.Stats.Draws++
}

func withTransactionT[T any](ctx context.Context, db *sql.DB, fn func(tx *sql.Tx) (T, error)) (result T, err error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return result, err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	result, err = fn(tx)
	return result, err
}

func withTransaction(ctx context.Context, db *sql.DB, fn func(tx *sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	err = fn(tx)
	return err
}
