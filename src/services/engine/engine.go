package engine

import (
	"context"
	"database/sql"
	"math/rand/v2"
	"strings"

	"github.com/plamen-v/tic-tac-toe-models/models"
	"github.com/plamen-v/tic-tac-toe/src/repository"
)

const (
	defaultBoardTile byte   = '_'
	defaultBoard     string = "_________"
	xMark            byte   = 'X'
	oMark            byte   = 'O'
)

type GameEngineService interface {
	GetOpenRooms(context.Context, string, string, string, models.RoomPhase) ([]*models.Room, error)
	CreateRoom(context.Context, *models.Room) (int64, error)
	PlayerJoinRoom(context.Context, int64, int64) error
	PlayerLeaveRoom(context.Context, int64, int64) error
	CreateGame(context.Context, int64, int64) (int64, error)
	GetGameState(context.Context, int64, int64) (*models.Game, error)
	PlayerMakeMove(context.Context, int64, int64, int) error
}

type gameEngineServiceImpl struct {
	db *sql.DB
}

func NewGameEngineService(db *sql.DB) GameEngineService {
	return &gameEngineServiceImpl{db: db}
}

func (s *gameEngineServiceImpl) playerRepositoryFactory(q repository.Querier) repository.PlayerRepository {
	return repository.NewPlayerRepository(q)
}

func (s *gameEngineServiceImpl) gameRepositoryFactory(q repository.Querier) repository.GameRepository {
	return repository.NewGameRepository(q)
}

func (s *gameEngineServiceImpl) roomRepositoryFactory(q repository.Querier) repository.RoomRepository {
	return repository.NewRoomRepository(q)
}

func (g *gameEngineServiceImpl) GetOpenRooms(ctx context.Context, host string, title string, description string, phase models.RoomPhase) ([]*models.Room, error) {
	return g.roomRepositoryFactory(g.db).GetList(ctx, host, title, description, phase)
}

func (g *gameEngineServiceImpl) CreateRoom(ctx context.Context, room *models.Room) (id int64, err error) {
	roomRepository := g.roomRepositoryFactory(g.db)
	err = g.validateCreateRoom(ctx, roomRepository, room, room.Host.ID)
	if err != nil {
		return 0, err
	}

	id, err = roomRepository.Create(ctx, room)
	if err != nil {
		return 0, nil
	}
	return id, nil
}

func (g *gameEngineServiceImpl) validateCreateRoom(ctx context.Context, roomRepository repository.RoomRepository, room *models.Room, playerID int64) error {
	playerRoom, err := roomRepository.GetByPlayerID(ctx, playerID)
	if err != nil && !models.IsNotFoundError(err) {
		return err
	}

	if playerRoom != nil {
		return models.NewValidationErrorf("player is part of other room")
	}

	if room != nil && len(room.Title) == 0 {
		return models.NewValidationErrorf("title is required")
	}

	return nil
}

func (g *gameEngineServiceImpl) PlayerJoinRoom(ctx context.Context, roomID int64, playerID int64) error {
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

		room.Guest = &models.RoomParticipant{
			ID:       playerID,
			Continue: true,
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

func (g *gameEngineServiceImpl) validatePlayerJoinRoom(ctx context.Context, roomRepository repository.RoomRepository, room *models.Room, playerID int64) error {
	if room.Guest != nil && room.Guest.ID != playerID && room.Phase == models.RoomPhaseFull {
		return models.NewValidationError("room is full")
	}

	if room.Host.ID == playerID {
		return models.NewValidationErrorf("player is already part of the room as host")
	}

	if room.Guest != nil && room.Guest.ID == playerID {
		return models.NewValidationErrorf("player is already part of the room as guest")
	}

	playerRoom, err := roomRepository.GetByPlayerID(ctx, playerID)
	if err != nil {
		if models.IsNotFoundError(err) {
			return nil
		}
		return err
	}
	if playerRoom != nil {
		return models.NewValidationErrorf("player is part of other room")
	}

	return nil
}

func (g *gameEngineServiceImpl) PlayerLeaveRoom(ctx context.Context, roomID int64, playerID int64) (err error) {
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

		if playerIsHost {
			err = roomRepository.Delete(ctx, room.Host.ID)
			if err != nil {
				return err
			}
		} else {
			room.Guest = nil
			err = roomRepository.Update(ctx, room)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (g *gameEngineServiceImpl) validatePlayerLeaveRoom(room *models.Room, playerID int64) error {
	if room.Host.ID != playerID &&
		(room.Guest == nil || room.Guest.ID != playerID) {
		return models.NewValidationErrorf("player is not part of the room")
	}

	return nil
}

func (g *gameEngineServiceImpl) CreateGame(ctx context.Context, roomID int64, playerID int64) (int64, error) {
	return withTransactionT(ctx, g.db, func(tx *sql.Tx) (int64, error) {
		roomRepository := g.roomRepositoryFactory(tx)
		room, err := roomRepository.Get(ctx, roomID, true)
		if err != nil {
			return 0, err
		}

		gameRepository := g.gameRepositoryFactory(tx)
		err = g.validateCreateGame(ctx, gameRepository, room, playerID)
		if err != nil {
			return 0, err
		}

		if playerID == room.Host.ID {
			room.Host.Continue = true
		}

		if room.Guest != nil && room.Guest.ID == playerID {
			room.Guest.Continue = true
		}

		var gameID int64
		if room.Guest != nil {
			if room.Guest.Continue && room.Host.Continue {
				err := g.createGame(ctx, gameRepository, room)
				if err != nil {
					return 0, err
				}

				gameID = *room.GameID
			}
		}

		err = roomRepository.Update(ctx, room)
		if err != nil {
			return 0, err
		}
		return gameID, nil
	})
}

func (g *gameEngineServiceImpl) validateCreateGame(ctx context.Context, gameRepository repository.GameRepository, room *models.Room, playerID int64) error {
	if room.Host.ID != playerID &&
		(room.Guest == nil || room.Guest.ID != playerID) {
		return models.NewValidationErrorf("player is not part of the room")
	}

	if room.GameID != nil {
		game, err := gameRepository.Get(ctx, *room.GameID)
		if err != nil {
			return err
		}

		if game.Phase == models.GamePhaseInProgress {
			return models.NewValidationErrorf("can't creat new game. Previous game is in progress.")
		}
	}
	return nil
}

func (g *gameEngineServiceImpl) GetGameState(ctx context.Context, roomID int64, playerID int64) (*models.Game, error) {
	roomRepository := g.roomRepositoryFactory(g.db)
	room, err := roomRepository.Get(ctx, roomID, false)
	if err != nil {
		return nil, err
	}

	err = g.validateGetGameState(room, playerID)
	if err != nil {
		return nil, err
	}

	gameRepository := g.gameRepositoryFactory(g.db)
	game, err := gameRepository.Get(ctx, *room.GameID)
	if err != nil {
		return nil, err
	}
	return game, nil
}

func (g *gameEngineServiceImpl) validateGetGameState(room *models.Room, playerID int64) error {
	if room.Host.ID != playerID &&
		(room.Guest == nil || room.Guest.ID != playerID) {
		return models.NewValidationError("player is not part of the room")
	}

	return nil
}

func (g *gameEngineServiceImpl) PlayerMakeMove(ctx context.Context, roomID int64, playerID int64, position int) error {
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

			mark := []byte(game.HostMark)[0]
			if playerID != room.Host.ID {
				mark = []byte(game.GuestMark)[0]
			}

			boardBytes := []byte(game.Board)
			boardBytes[position-1] = mark
			game.Board = string(boardBytes)

			win := g.inWinState(game)
			if win || !strings.Contains(game.Board, string(defaultBoardTile)) {
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
				if game.CurrentPlayerID == game.HostID {
					game.CurrentPlayerID = game.GuestID
				} else {
					game.CurrentPlayerID = game.HostID
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

func (g *gameEngineServiceImpl) validatePlayerMakeMove(game *models.Game, playerID int64, position int) error {
	if game.HostID != playerID && game.GuestID != playerID {
		return models.NewValidationError("player is not part of the game")
	}

	if game.Phase == models.GamePhaseCompleted {
		return models.NewValidationError("game is completed")
	}

	if game.CurrentPlayerID != playerID {
		return models.NewValidationError("player not in turn")
	}

	if position < 1 || position > len(defaultBoard) {
		return models.NewValidationError("invalid position index")
	}
	if game.Board[position-1] != defaultBoardTile {
		return models.NewValidationError("position ocopied")
	}

	return nil
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
	room.Host.Continue = false
	room.Guest.Continue = false

	return nil
}

func (g *gameEngineServiceImpl) initializeGame(ctx context.Context, gameRepository repository.GameRepository, room *models.Room) (*models.Game, error) {
	marks := []byte{xMark, oMark}
	rand.Shuffle(len(marks), func(i, j int) {
		marks[i], marks[j] = marks[j], marks[i]
	})
	playerIDs := []int64{room.Host.ID, room.Guest.ID}

	game := &models.Game{
		HostID:          room.Host.ID,
		HostMark:        string(marks[0]),
		GuestID:         room.Guest.ID,
		GuestMark:       string(marks[1]),
		CurrentPlayerID: playerIDs[rand.IntN(2)],
		Board:           defaultBoard,
		Phase:           models.GamePhaseInProgress,
	}

	if room.GameID != nil {
		prevGame, err := gameRepository.Get(ctx, *room.GameID)
		if err != nil {
			return nil, err
		}

		if prevGame.HostID == room.Host.ID && prevGame.GuestID == room.Guest.ID {
			game.HostMark = prevGame.HostMark
			game.GuestMark = prevGame.GuestMark

			game.CurrentPlayerID = room.Host.ID
			if prevGame.CurrentPlayerID == room.Host.ID {
				game.CurrentPlayerID = room.Guest.ID
			}
		}
	}

	return game, nil
}

func (g *gameEngineServiceImpl) inWinState(game *models.Game) bool {

	if (game.Board[0] != defaultBoardTile && game.Board[0] == game.Board[1] && game.Board[0] == game.Board[2]) ||
		(game.Board[3] != defaultBoardTile && game.Board[3] == game.Board[4] && game.Board[3] == game.Board[5]) ||
		(game.Board[6] != defaultBoardTile && game.Board[6] == game.Board[7] && game.Board[6] == game.Board[8]) ||

		(game.Board[0] != defaultBoardTile && game.Board[0] == game.Board[3] && game.Board[0] == game.Board[6]) ||
		(game.Board[1] != defaultBoardTile && game.Board[1] == game.Board[4] && game.Board[1] == game.Board[7]) ||
		(game.Board[2] != defaultBoardTile && game.Board[2] == game.Board[5] && game.Board[2] == game.Board[8]) ||

		(game.Board[0] != defaultBoardTile && game.Board[0] == game.Board[4] && game.Board[0] == game.Board[8]) ||
		(game.Board[2] != defaultBoardTile && game.Board[2] == game.Board[4] && game.Board[2] == game.Board[6]) {
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
