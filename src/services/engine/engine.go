package engine

import (
	"database/sql"
	"math/rand/v2"

	"github.com/plamen-v/tic-tac-toe-models/models"
	"github.com/plamen-v/tic-tac-toe-models/models/errors"
	"github.com/plamen-v/tic-tac-toe/src/services/repository"
)

// todo! not set win lose . be good
type GameEngineService interface {
	GetOpenRooms(string, string, string, models.RoomPhase) ([]*models.Room, error)
	CreateRoom(*models.Room) (int64, error)
	GetRoomState(int64, int64) (*models.Room, error)
	HostLeave(int64, int64) error
	RegisterGuest(int64, int64) error
	GuestLeave(int64, int64) error
	CreateGame(int64, int64) (int64, error)
	GetGameState(int64, int64, int64) (*models.Game, error)
	SetMark(int64, int64, int64, int) error
}

type gameEngineService struct {
	repo repository.Repository
}

const (
	defaultBoard string = "         "
	xMark               = "X"
	oMark               = "O"
)

func NewGameEngineService(repo repository.Repository) GameEngineService {
	return &gameEngineService{
		repo: repo,
	}
}

func (g *gameEngineService) GetOpenRooms(host string, title string, description string, phase models.RoomPhase) ([]*models.Room, error) {
	return g.repo.Rooms().GetList(host, title, description, phase)
}

func (g *gameEngineService) CreateRoom(room *models.Room) (id int64, err error) {
	tx, err := g.repo.Begin()
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	host, err := g.repo.Players().Get(room.Host.ID, tx)
	if err != nil {
		return
	}

	err = g.validateCreateRoom(room, host)
	if err != nil {
		return
	}

	id, err = g.repo.Rooms().Create(room, tx)
	if err != nil {
		return
	}

	host.RoomID = &id
	err = g.repo.Players().Update(host, tx)
	if err != nil {
		return
	}

	return
}

func (g *gameEngineService) validateCreateRoom(room *models.Room, host *models.Player) error {
	if host.RoomID != nil {
		return errors.NewValidationErrorf("player has room")
	}
	if len(room.Title) == 0 {
		return errors.NewValidationErrorf("title is required")
	}

	return nil
}

func (g *gameEngineService) GetRoomState(roomID int64, playerID int64) (*models.Room, error) {
	room, err := g.repo.Rooms().Get(roomID, nil)
	if err != nil {
		return nil, err
	}
	err = g.validateGetRoomState(room, playerID)
	if err != nil {
		return nil, err
	}

	return room, nil
}

func (g *gameEngineService) validateGetRoomState(room *models.Room, playerID int64) error {
	if room.Host.ID != playerID ||
		(room.Guest != nil && room.Guest.ID != playerID) {
		return errors.NewAuthorizationError("player not in room")
	}
	return nil
}

func (g *gameEngineService) HostLeave(roomID int64, hostID int64) (err error) {
	var (
		room  *models.Room
		host  *models.Player
		guest *models.Player
		game  *models.Game
	)

	tx, err := g.repo.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	room, err = g.repo.Rooms().Get(roomID, tx)
	if err != nil {
		return err
	}
	err = g.validateHostLeave(room, hostID)
	if err != nil {
		return err
	}

	host, err = g.repo.Players().Get(room.Host.ID, tx)
	if err != nil {
		return err
	}
	host.RoomID = nil

	if room.Guest != nil {
		guest, err = g.repo.Players().Get(room.Guest.ID, tx)
		if err != nil {
			return err
		}
		guest.RoomID = nil

		if room.GameID != nil {
			game, err = g.repo.Games().Get(*room.GameID, tx)
			if err != nil {
				return err
			}
			g.closeGame(game, host, guest)

			err = g.repo.Games().Update(game, tx)
			if err != nil {
				return err
			}
		}

		err = g.repo.Players().Update(guest, tx)
		if err != nil {
			return err
		}
	}

	err = g.repo.Players().Update(host, tx)
	if err != nil {
		return err
	}

	err = g.repo.Rooms().Delete(room.ID, tx)
	if err != nil {
		return err
	}

	return nil
}

func (g *gameEngineService) validateHostLeave(room *models.Room, playerID int64) error {
	if room.Host.ID != playerID {
		return errors.NewAuthorizationError("player not host of the room")
	}
	return nil
}

func (g *gameEngineService) RegisterGuest(roomID int64, guestID int64) error {
	var (
		room  *models.Room
		host  *models.Player
		guest *models.Player
	)

	tx, err := g.repo.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	room, err = g.repo.Rooms().Get(roomID, tx)
	if err != nil {
		return err
	}
	guest, err = g.repo.Players().Get(guestID, tx)
	if err != nil {
		return err
	}
	err = g.validateRegisterGuest(room, guest)
	if err != nil {
		return err
	}

	guest.RoomID = &room.ID
	room.Guest = &models.RoomParticipant{
		ID:       guest.ID,
		Continue: true,
	}
	room.Phase = models.RoomPhaseFull

	host, err = g.repo.Players().Get(room.Host.ID, tx)
	if err != nil {
		return err
	}

	err = g.createGameHelper(room, host, guest, tx)
	if err != nil {
		return err
	}

	err = g.repo.Players().Update(host, tx)
	if err != nil {
		return err
	}

	err = g.repo.Players().Update(guest, tx)
	if err != nil {
		return err
	}

	err = g.repo.Rooms().Update(room, tx)
	if err != nil {
		return err
	}

	return nil
}

func (g *gameEngineService) validateRegisterGuest(room *models.Room, guest *models.Player) error {
	if room.Phase == models.RoomPhaseFull {
		return errors.NewValidationError("room is full")
	}

	if guest.RoomID != nil {
		return errors.NewValidationErrorf("player has room")
	}

	if room.Host.ID == guest.ID {
		return errors.NewValidationErrorf("player is host of the room")
	}

	return nil
}

func (g *gameEngineService) GuestLeave(roomID int64, guestID int64) error {
	var (
		room  *models.Room
		host  *models.Player
		guest *models.Player
		game  *models.Game
	)

	tx, err := g.repo.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	room, err = g.repo.Rooms().Get(roomID, tx)
	if err != nil {
		return err
	}
	err = g.validateGuestLeave(room, guestID)
	if err != nil {
		return err
	}

	guest, err = g.repo.Players().Get(guestID, tx)
	if err != nil {
		return err
	}
	guest.RoomID = nil
	room.Phase = models.RoomPhaseOpen

	if room.GameID != nil {
		host, err = g.repo.Players().Get(room.Host.ID, tx)
		if err != nil {
			return err
		}
		game, err = g.repo.Games().Get(*room.GameID, tx)
		if err != nil {
			return err
		}

		g.closeGame(game, host, guest)

		err = g.repo.Players().Update(host, tx)
		if err != nil {
			return err
		}

		err = g.repo.Games().Update(game, tx)
		if err != nil {
			return err
		}
	}

	err = g.repo.Players().Update(guest, tx)
	if err != nil {
		return err
	}

	err = g.repo.Rooms().Update(room, tx)
	if err != nil {
		return err
	}

	return nil
}

func (g *gameEngineService) validateGuestLeave(room *models.Room, playerID int64) error {
	if room.Guest != nil && room.Guest.ID != playerID {
		return errors.NewAuthorizationError("player not guest of the room")
	}
	return nil
}

func (g *gameEngineService) CreateGame(roomID int64, playerID int64) (int64, error) {
	var (
		room   *models.Room
		host   *models.Player
		guest  *models.Player
		gameID int64
		err    error
	)

	tx, err := g.repo.Begin()
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	room, err = g.repo.Rooms().Get(roomID, tx)
	if err != nil {
		return 0, err
	}
	err = g.validateCreateGame(room, playerID)
	if err != nil {
		return 0, err
	}

	if playerID == room.Host.ID {
		room.Host.Continue = true
	}

	if room.Guest != nil && room.Guest.ID == playerID {
		room.Guest.Continue = true
	}

	if room.Guest != nil {
		if room.Guest.Continue && room.Host.Continue {
			err := g.createGameHelper(room, host, guest, tx)
			if err != nil {
				return 0, err
			}

			err = g.repo.Players().Update(host, tx)
			if err != nil {
				return 0, err
			}

			err = g.repo.Players().Update(guest, tx)
			if err != nil {
				return 0, err
			}

			gameID = *room.GameID
		}
	}

	err = g.repo.Rooms().Update(room, tx)
	if err != nil {
		return 0, err
	}

	return gameID, nil
}

func (g *gameEngineService) validateCreateGame(room *models.Room, playerID int64) error {
	if room.Host.ID != playerID || (room.Guest != nil && room.Guest.ID != playerID) {
		return errors.NewAuthorizationError("player not in the room")
	}
	return nil
}

func (g *gameEngineService) createGameHelper(room *models.Room, host *models.Player, guest *models.Player, tx *sql.Tx) error {
	var (
		game *models.Game
		err  error
	)

	game, err = g.initializeGame(room, host, guest, tx)
	if err != nil {
		return err
	}
	game.ID, err = g.repo.Games().Create(game, tx)
	if err != nil {
		return err
	}
	room.PrevGameID = room.GameID
	room.GameID = &game.ID
	room.Host.Continue = false
	room.Guest.Continue = false
	host.GameID = &game.ID
	guest.GameID = &game.ID

	return nil
}

func (g *gameEngineService) initializeGame(room *models.Room, host *models.Player, guest *models.Player, tx *sql.Tx) (*models.Game, error) {
	marks := []string{xMark, oMark}
	rand.Shuffle(len(marks), func(i, j int) {
		marks[i], marks[j] = marks[j], marks[i]
	})
	playerIDs := []int64{room.Host.ID, room.Guest.ID}

	game := &models.Game{
		HostID:          host.ID,
		HostMark:        marks[0],
		GuestID:         guest.ID,
		GuestMark:       marks[1],
		CurrentPlayerID: playerIDs[rand.IntN(2)],
		Board:           defaultBoard,
		Phase:           models.GamePhaseStarted,
	}

	if room.PrevGameID != nil {
		prevGame, err := g.repo.Games().Get(*room.GameID, tx)
		if err != nil {
			return nil, err
		}

		if prevGame.HostID == room.Host.ID && prevGame.GuestID == room.Guest.ID {
			game.HostMark = prevGame.HostMark
			game.GuestMark = prevGame.GuestMark
			if prevGame.LoserID != nil {
				game.CurrentPlayerID = *prevGame.LoserID
			} else {
				if prevGame.CurrentPlayerID == room.Host.ID {
					game.CurrentPlayerID = room.Guest.ID
				} else {
					game.CurrentPlayerID = room.Host.ID
				}
			}
		}
	}

	return game, nil
}

func (g *gameEngineService) GetGameState(roomID int64, gameID int64, playerID int64) (*models.Game, error) {
	var (
		room *models.Room
		err  error
	)

	tx, err := g.repo.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()
	room, err = g.repo.Rooms().Get(roomID, tx)
	if err != nil {
		return nil, err
	}
	err = g.validateGetGameState(room, playerID)
	if err != nil {
		return nil, err
	}
	return g.repo.Games().Get(gameID, nil)
}

func (g *gameEngineService) validateGetGameState(room *models.Room, playerID int64) error {
	if room.Host.ID != playerID || (room.Guest != nil && room.Guest.ID != playerID) {
		return errors.NewAuthorizationError("player not in the room")
	}
	return nil
}

func (g *gameEngineService) SetMark(roomID int64, gameID int64, playerID int64, position int) error {
	var (
		room *models.Room
		game *models.Game
		//player *models.Player
		err error
	)

	tx, err := g.repo.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	room, err = g.repo.Rooms().Get(roomID, tx)
	if err != nil {
		return err
	}

	// player, err = g.repo.Players().Get(playerID, tx)
	// if err != nil {
	// 	return err
	// }

	game, err = g.repo.Games().Get(gameID, tx)
	if err != nil {
		return err
	}
	err = g.validateSetMark(room, game, playerID)
	if err != nil {
		return err
	}

	return nil
}

func (g *gameEngineService) validateSetMark(room *models.Room, game *models.Game, playerID int64) error { //todo!
	if room.Host.ID != playerID || (room.Guest != nil && room.Guest.ID != playerID) {
		return errors.NewAuthorizationError("player not in the room")
	}
	//

	return nil
}

func (g *gameEngineService) checkWinState(room *models.Room, game *models.Game, playerID int64) error { //todo!
	if room.Host.ID != playerID || (room.Guest != nil && room.Guest.ID != playerID) {
		return errors.NewAuthorizationError("player not in the room")
	}
	//

	return nil
}

func (g *gameEngineService) closeGame(game *models.Game, leaver *models.Player, other *models.Player) {
	game.Phase = models.GamePhaseDone
	game.WinnerID = &other.ID
	game.LoserID = &leaver.ID
	leaver.Losses = leaver.Losses + 1
	leaver.GameID = nil
	other.Wins = other.Wins + 1
	other.GameID = nil
}
