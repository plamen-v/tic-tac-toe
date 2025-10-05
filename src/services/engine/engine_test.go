package engine_test

import (
	"context"
	"database/sql"
	"strings"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gofrs/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/plamen-v/tic-tac-toe-models/models"
	"github.com/plamen-v/tic-tac-toe/src/repository"
	"github.com/plamen-v/tic-tac-toe/src/repository/mocks"
	"github.com/plamen-v/tic-tac-toe/src/services/engine"
	tmock "github.com/stretchr/testify/mock"
)

var _ = Describe("GameEngine", func() {
	var (
		db                   *sql.DB
		mock                 sqlmock.Sqlmock
		ctx                  context.Context
		mockRoomRepository   *mocks.MockRoomRepository
		mockGameRepository   *mocks.MockGameRepository
		mockPlayerRepository *mocks.MockPlayerRepository
		gameEngineService    engine.GameEngineService
		err                  error
	)

	BeforeEach(func() {
		ctx = context.TODO()
		db, mock, err = sqlmock.New()
		Expect(err).ToNot(HaveOccurred())
		mockRoomRepository = new(mocks.MockRoomRepository)
		mockGameRepository = new(mocks.MockGameRepository)
		mockPlayerRepository = new(mocks.MockPlayerRepository)
		gameEngineService = engine.NewGameEngineService(
			db,
			func(db repository.Querier) repository.PlayerRepository {
				return mockPlayerRepository
			},
			func(db repository.Querier) repository.GameRepository {
				return mockGameRepository
			},
			func(db repository.Querier) repository.RoomRepository {
				return mockRoomRepository
			},
		)

	})

	AfterEach(func() {
		err = mock.ExpectationsWereMet()
		Expect(err).ToNot(HaveOccurred())
		db.Close()
	})

	Context("GetRoom", func() {
		It("should returns expected room", func() {
			roomID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			expectedRoom := &models.Room{
				ID: roomID,
			}

			playerID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			mockRoomRepository.
				On("GetByPlayerID", ctx, playerID).
				Return(expectedRoom, nil)

			room, err := gameEngineService.GetRoom(ctx, playerID)

			Expect(err).ToNot(HaveOccurred())
			Expect(*room).To(Equal(*expectedRoom))
		})
	})

	Context("GetOpenRooms", func() {
		It("should returns list of rooms", func() {
			roomID1, err := uuid.NewV4()
			Expect(err).To(BeNil())
			expectedRooms := []*models.Room{
				{ID: roomID1,
					Phase: models.RoomPhaseOpen,
				},
			}

			mockRoomRepository.
				On("GetList", ctx, models.RoomPhaseOpen).
				Return(expectedRooms, nil)
			rooms, err := gameEngineService.GetOpenRooms(ctx)

			Expect(err).ToNot(HaveOccurred())
			Expect(len(rooms)).To(Equal(len(expectedRooms)))
			for i := 0; i < len(expectedRooms); i++ {
				Expect(rooms[i]).To(Equal(expectedRooms[i]))
			}
		})
	})

	Context("GetGameState", func() {
		It("should returns the game status if the player is host in a existing room that has a game", func() {
			gameID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			expectedGame := &models.Game{
				ID: gameID,
			}

			roomID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			playerID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			expectedRoom := &models.Room{
				ID: roomID,
				Host: models.RoomPlayer{
					ID: playerID,
				},
				GameID: &expectedGame.ID,
			}

			mockRoomRepository.
				On("Get", ctx, expectedRoom.ID, false).
				Return(expectedRoom, nil)

			mockGameRepository.
				On("Get", ctx, expectedGame.ID).
				Return(expectedGame, nil)

			game, err := gameEngineService.GetGameState(ctx, roomID, playerID)

			Expect(err).ToNot(HaveOccurred())
			Expect(*game).To(Equal(*expectedGame))
		})

		It("should returns the game status if the player is guest in a existing room that has a game", func() {
			gameID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			expectedGame := &models.Game{
				ID: gameID,
			}

			roomID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			playerID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			expectedRoom := &models.Room{
				ID: roomID,
				Guest: &models.RoomPlayer{
					ID: playerID,
				},
				GameID: &expectedGame.ID,
			}

			mockRoomRepository.
				On("Get", ctx, expectedRoom.ID, false).
				Return(expectedRoom, nil)

			mockGameRepository.
				On("Get", ctx, expectedGame.ID).
				Return(expectedGame, nil)

			game, err := gameEngineService.GetGameState(ctx, roomID, playerID)

			Expect(err).ToNot(HaveOccurred())
			Expect(*game).To(Equal(*expectedGame))
		})

		It("should returns an error if player is not in the room", func() {
			roomID, err := uuid.NewV4()
			Expect(err).To(BeNil())

			hostID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			guestID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			expectedRoom := &models.Room{
				ID: roomID,
				Host: models.RoomPlayer{
					ID: hostID,
				},
				Guest: &models.RoomPlayer{
					ID: guestID,
				},
			}

			mockRoomRepository.
				On("Get", ctx, expectedRoom.ID, false).
				Return(expectedRoom, nil)

			playerID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			_, err = gameEngineService.GetGameState(ctx, roomID, playerID)

			expectedErrorMessage := engine.PlayerNotInRoomErrorMessage
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(expectedErrorMessage))
			mockGameRepository.AssertExpectations(GinkgoT())
		})
	})

	Context("GetRanking", func() {
		It("should return the ranking", func() {
			playerID1, err := uuid.NewV4()
			Expect(err).To(BeNil())
			playerID2, err := uuid.NewV4()
			Expect(err).To(BeNil())
			expectedRanking := []*models.Player{
				{ID: playerID1, Stats: models.PlayerStats{Wins: 1, Losses: 1, Draws: 1}},
				{ID: playerID2, Stats: models.PlayerStats{Wins: 1, Losses: 1, Draws: 1}},
			}

			mockPlayerRepository.
				On("GetRanking", ctx).
				Return(expectedRanking, 1, 1, 1, nil)

			ranking, _, _, _, err := gameEngineService.GetRanking(ctx, 1, 1)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(ranking)).To(Equal(len(expectedRanking)))
			for i := 0; i < len(expectedRanking); i++ {
				Expect(ranking[i]).To(Equal(expectedRanking[i]))
			}
		})
	})

	Context("CreateRoom", func() {
		It("should returns id of the created room if input is valid and player is not in other room", func() {
			roomID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			playerID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			expectedRoom := &models.Room{
				ID: roomID,
				Host: models.RoomPlayer{
					ID: playerID,
				},
				Title:       "title",
				Description: "description",
				Phase:       models.RoomPhaseOpen,
			}

			mockRoomRepository.
				On("GetByPlayerID", ctx, playerID).
				Return(nil, models.NewNotFoundError("error"))

			mockRoomRepository.
				On("Create", ctx, tmock.Anything).
				Return(expectedRoom.ID, nil)

			resultRoomID, err := gameEngineService.CreateRoom(ctx, playerID, expectedRoom.Title, expectedRoom.Description)
			Expect(err).ToNot(HaveOccurred())
			Expect(resultRoomID).To(Equal(expectedRoom.ID))

		})

		It("should returns error if player is in other room", func() {
			roomID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			playerID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			expectedRoom := &models.Room{
				ID: roomID,
				Host: models.RoomPlayer{
					ID: playerID,
				},
				Title:       "title",
				Description: "description",
				Phase:       models.RoomPhaseOpen,
			}

			playerRoom := &models.Room{
				ID: roomID,
				Host: models.RoomPlayer{
					ID: playerID,
				},
			}
			mockRoomRepository.
				On("GetByPlayerID", ctx, playerID).
				Return(playerRoom, nil)

			_, err = gameEngineService.CreateRoom(ctx, playerID, expectedRoom.Title, expectedRoom.Description)

			expectedErrorMessage := engine.PlayerPartOfOtherRoomErrorMessage
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(expectedErrorMessage))
			mockGameRepository.AssertExpectations(GinkgoT())
		})

		It("should returns error if input is invalid and title is empty", func() {
			roomID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			playerID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			expectedRoom := &models.Room{
				ID: roomID,
				Host: models.RoomPlayer{
					ID: playerID,
				},
				Title: "",
				Phase: models.RoomPhaseOpen,
			}

			mockRoomRepository.
				On("GetByPlayerID", ctx, playerID).
				Return(nil, models.NewNotFoundError("error"))

			_, err = gameEngineService.CreateRoom(ctx, playerID, expectedRoom.Title, expectedRoom.Description)

			expectedErrorMessage := engine.TitleRequiredErrorMessage
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(expectedErrorMessage))
			mockGameRepository.AssertExpectations(GinkgoT())
		})

		It("should returns error if input is invalid and title is too long", func() {
			roomID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			playerID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			expectedRoom := &models.Room{
				ID: roomID,
				Host: models.RoomPlayer{
					ID: playerID,
				},
				Title: strings.Repeat(string('x'), engine.MaxRoomTitleLength+1),
				Phase: models.RoomPhaseOpen,
			}

			mockRoomRepository.
				On("GetByPlayerID", ctx, playerID).
				Return(nil, models.NewNotFoundError("error"))

			_, err = gameEngineService.CreateRoom(ctx, playerID, expectedRoom.Title, expectedRoom.Description)

			expectedErrorMessage := engine.TitleTooLongErrorMessage
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(expectedErrorMessage))
			mockGameRepository.AssertExpectations(GinkgoT())
		})

		It("should returns error if input is invalid and title is too long", func() {
			roomID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			playerID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			expectedRoom := &models.Room{
				ID: roomID,
				Host: models.RoomPlayer{
					ID: playerID,
				},
				Title:       "title",
				Description: strings.Repeat(string('x'), engine.MaxRoomDescriptionLength+1),
				Phase:       models.RoomPhaseOpen,
			}

			mockRoomRepository.
				On("GetByPlayerID", ctx, playerID).
				Return(nil, models.NewNotFoundError("error"))

			_, err = gameEngineService.CreateRoom(ctx, playerID, expectedRoom.Title, expectedRoom.Description)

			expectedErrorMessage := engine.DescriptionTooLongErrorMessage
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(expectedErrorMessage))
			mockGameRepository.AssertExpectations(GinkgoT())
		})
	})

	Context("PlayerJoinRoom", func() {
		It("should return no error if room is not full and player is not part of other room", func() {
			mock.ExpectBegin()
			mock.ExpectCommit()

			roomID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			room := &models.Room{
				ID: roomID,
			}
			mockRoomRepository.
				On("Get", ctx, roomID, true).
				Return(room, nil)

			playerID, err := uuid.NewV4()
			Expect(err).To(BeNil())

			mockRoomRepository.
				On("GetByPlayerID", ctx, playerID).
				Return(nil, models.NewNotFoundError("not found"))

			mockRoomRepository.
				On("Update", ctx, room).
				Return(nil)

			gameID, err := uuid.NewV4()
			Expect(err).To(BeNil())

			mockGameRepository.
				On("Create", ctx, tmock.Anything).Return(gameID, nil)

			err = gameEngineService.PlayerJoinRoom(ctx, roomID, playerID)

			Expect(err).ToNot(HaveOccurred())
			mockGameRepository.AssertExpectations(GinkgoT())
		})

		It("should return error if room is full", func() {
			mock.ExpectBegin()
			mock.ExpectRollback()

			roomID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			room := &models.Room{
				ID:    roomID,
				Phase: models.RoomPhaseFull,
			}
			mockRoomRepository.
				On("Get", ctx, roomID, true).
				Return(room, nil)

			playerID, err := uuid.NewV4()
			Expect(err).To(BeNil())

			err = gameEngineService.PlayerJoinRoom(ctx, roomID, playerID)

			expectedErrorMessage := engine.FullRoomErrorMessage
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(expectedErrorMessage))
			mockGameRepository.AssertExpectations(GinkgoT())
		})

		It("should return error if player is in the room as host", func() {
			mock.ExpectBegin()
			mock.ExpectRollback()

			playerID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			roomID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			room := &models.Room{
				ID:    roomID,
				Host:  models.RoomPlayer{ID: playerID},
				Phase: models.RoomPhaseOpen,
			}
			mockRoomRepository.
				On("Get", ctx, roomID, true).
				Return(room, nil)

			err = gameEngineService.PlayerJoinRoom(ctx, roomID, playerID)

			expectedErrorMessage := engine.PlayerPartOfTheRoomAsHostErrorMessage
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(expectedErrorMessage))
			mockGameRepository.AssertExpectations(GinkgoT())
		})

		It("should return error if player is in the room as guest", func() {
			mock.ExpectBegin()
			mock.ExpectRollback()

			playerID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			roomID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			room := &models.Room{
				ID:    roomID,
				Guest: &models.RoomPlayer{ID: playerID},
				Phase: models.RoomPhaseOpen,
			}
			mockRoomRepository.
				On("Get", ctx, roomID, true).
				Return(room, nil)

			err = gameEngineService.PlayerJoinRoom(ctx, roomID, playerID)

			expectedErrorMessage := engine.PlayerPartOfTheRoomAsGuestErrorMessage
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(expectedErrorMessage))
			mockGameRepository.AssertExpectations(GinkgoT())
		})

		It("should return error if player is in other room", func() {
			mock.ExpectBegin()
			mock.ExpectRollback()

			playerID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			roomID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			room := &models.Room{
				ID:    roomID,
				Phase: models.RoomPhaseOpen,
			}
			mockRoomRepository.
				On("Get", ctx, roomID, true).
				Return(room, nil)

			playerRoom := &models.Room{
				ID: roomID,
				Host: models.RoomPlayer{
					ID: playerID,
				},
			}
			mockRoomRepository.
				On("GetByPlayerID", ctx, playerID).
				Return(playerRoom, nil)

			err = gameEngineService.PlayerJoinRoom(ctx, roomID, playerID)

			expectedErrorMessage := engine.PlayerPartOfOtherRoomErrorMessage
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(expectedErrorMessage))
			mockGameRepository.AssertExpectations(GinkgoT())
		})

	})

	Context("PlayerLeaveRoom", func() {

		It("should return no error if room exist, game is in progress and player is in the room as host", func() {
			mock.ExpectBegin()
			mock.ExpectCommit()

			playerID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			host := &models.Player{
				ID: playerID,
			}

			guestID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			guest := &models.Player{
				ID: guestID,
			}

			gameID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			game := &models.Game{
				ID: gameID,
				Host: models.GamePlayer{
					ID: host.ID,
				},
				Guest: models.GamePlayer{
					ID: guest.ID,
				},
				Phase: models.GamePhaseInProgress,
			}

			roomID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			room := &models.Room{
				ID: roomID,
				Host: models.RoomPlayer{
					ID: host.ID,
				},
				Guest: &models.RoomPlayer{
					ID: guest.ID,
				},
				GameID: &game.ID,
				Phase:  models.RoomPhaseFull,
			}

			mockRoomRepository.
				On("Get", ctx, roomID, true).
				Return(room, nil)

			mockGameRepository.
				On("Get", ctx, gameID).
				Return(game, nil)

			mockPlayerRepository.
				On("Get", ctx, playerID).
				Return(host, nil)

			mockPlayerRepository.
				On("Get", ctx, guestID).
				Return(guest, nil)

			mockPlayerRepository.
				On("UpdateStats", ctx, tmock.Anything).
				Return(nil)

			mockGameRepository.
				On("Update", ctx, game).
				Return(nil)

			mockRoomRepository.
				On("Update", ctx, room).
				Return(nil)

			err = gameEngineService.PlayerLeaveRoom(ctx, roomID, playerID)

			Expect(err).ToNot(HaveOccurred())
			mockGameRepository.AssertExpectations(GinkgoT())
		})

		It("should return no error if room exist, game is in progress and player is in the room as guest", func() {
			mock.ExpectBegin()
			mock.ExpectCommit()

			hostID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			host := &models.Player{
				ID: hostID,
			}

			playerID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			guest := &models.Player{
				ID: playerID,
			}

			gameID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			game := &models.Game{
				ID: gameID,
				Host: models.GamePlayer{
					ID: host.ID,
				},
				Guest: models.GamePlayer{
					ID: guest.ID,
				},
				Phase: models.GamePhaseInProgress,
			}

			roomID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			room := &models.Room{
				ID: roomID,
				Host: models.RoomPlayer{
					ID: host.ID,
				},
				Guest: &models.RoomPlayer{
					ID: guest.ID,
				},
				GameID: &game.ID,
				Phase:  models.RoomPhaseFull,
			}

			mockRoomRepository.
				On("Get", ctx, roomID, true).
				Return(room, nil)

			mockGameRepository.
				On("Get", ctx, gameID).
				Return(game, nil)

			mockPlayerRepository.
				On("Get", ctx, host.ID).
				Return(host, nil)

			mockPlayerRepository.
				On("Get", ctx, playerID).
				Return(guest, nil)

			mockPlayerRepository.
				On("UpdateStats", ctx, tmock.Anything).
				Return(nil)

			mockGameRepository.
				On("Update", ctx, game).
				Return(nil)

			mockRoomRepository.
				On("Update", ctx, room).
				Return(nil)

			err = gameEngineService.PlayerLeaveRoom(ctx, roomID, playerID)

			Expect(err).ToNot(HaveOccurred())
			mockGameRepository.AssertExpectations(GinkgoT())
		})

		It("should return error if player is not in the room", func() {
			mock.ExpectBegin()
			mock.ExpectRollback()

			hostID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			host := &models.Player{
				ID: hostID,
			}

			guestID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			guest := &models.Player{
				ID: guestID,
			}

			roomID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			room := &models.Room{
				ID: roomID,
				Host: models.RoomPlayer{
					ID: host.ID,
				},
				Guest: &models.RoomPlayer{
					ID: guest.ID,
				},
			}

			mockRoomRepository.
				On("Get", ctx, roomID, true).
				Return(room, nil)

			playerID, err := uuid.NewV4()
			Expect(err).To(BeNil())

			err = gameEngineService.PlayerLeaveRoom(ctx, roomID, playerID)

			expectedErrorMessage := engine.PlayerNotInRoomErrorMessage
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(expectedErrorMessage))
			mockGameRepository.AssertExpectations(GinkgoT())
		})

		It("should return no error if player leave last the room as host", func() {
			mock.ExpectBegin()
			mock.ExpectCommit()

			playerID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			host := &models.Player{
				ID: playerID,
			}

			roomID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			room := &models.Room{
				ID: roomID,
				Host: models.RoomPlayer{
					ID: host.ID,
				},
				Phase: models.RoomPhaseOpen,
			}

			mockRoomRepository.
				On("Get", ctx, roomID, true).
				Return(room, nil)

			mockRoomRepository.
				On("Delete", ctx, room.ID).
				Return(nil)

			err = gameEngineService.PlayerLeaveRoom(ctx, roomID, playerID)

			Expect(err).ToNot(HaveOccurred())
			mockGameRepository.AssertExpectations(GinkgoT())
		})
	})

	Context("CreateGame", func() {
		It("should return new game id if both players have new game requests", func() {
			mock.ExpectBegin()
			mock.ExpectCommit()

			playerID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			host := &models.Player{
				ID: playerID,
			}

			guestID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			guest := &models.Player{
				ID: guestID,
			}

			prevGameID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			prevGame := &models.Game{
				ID:    prevGameID,
				Phase: models.GamePhaseCompleted,
				Host: models.GamePlayer{
					ID: host.ID,
				},
				Guest: models.GamePlayer{
					ID: guest.ID,
				},
				CurrentPlayerID: host.ID,
			}

			roomID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			room := &models.Room{
				ID: roomID,
				Host: models.RoomPlayer{
					ID: host.ID,
				},
				Guest: &models.RoomPlayer{
					ID:             guest.ID,
					RequestNewGame: true,
				},
				GameID: &prevGame.ID,
				Phase:  models.RoomPhaseFull,
			}

			mockRoomRepository.
				On("Get", ctx, roomID, true).
				Return(room, nil)

			mockGameRepository.
				On("Get", ctx, prevGameID).
				Return(prevGame, nil)

			expectedNewGameID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			mockGameRepository.
				On("Create", ctx, tmock.Anything).
				Return(expectedNewGameID, nil)

			mockRoomRepository.
				On("Update", ctx, room).
				Return(nil)

			newGameID, err := gameEngineService.CreateGame(ctx, roomID, playerID)

			Expect(err).ToNot(HaveOccurred())
			Expect(newGameID).To(Equal(expectedNewGameID))
			mockGameRepository.AssertExpectations(GinkgoT())
		})

		It("should return error if previous game is in progress", func() {
			mock.ExpectBegin()
			mock.ExpectRollback()

			playerID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			host := &models.Player{
				ID: playerID,
			}

			guestID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			guest := &models.Player{
				ID: guestID,
			}

			prevGameID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			prevGame := &models.Game{
				ID:    prevGameID,
				Phase: models.GamePhaseInProgress,
				Host: models.GamePlayer{
					ID: host.ID,
				},
				Guest: models.GamePlayer{
					ID: guest.ID,
				},
				CurrentPlayerID: host.ID,
			}

			roomID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			room := &models.Room{
				ID: roomID,
				Host: models.RoomPlayer{
					ID: host.ID,
				},
				Guest: &models.RoomPlayer{
					ID:             guest.ID,
					RequestNewGame: true,
				},
				GameID: &prevGame.ID,
				Phase:  models.RoomPhaseFull,
			}

			mockRoomRepository.
				On("Get", ctx, roomID, true).
				Return(room, nil)

			mockGameRepository.
				On("Get", ctx, prevGameID).
				Return(prevGame, nil)

			_, err = gameEngineService.CreateGame(ctx, roomID, playerID)

			expectedErrorMessage := engine.GameInProgressErrorMessage
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(expectedErrorMessage))
			mockGameRepository.AssertExpectations(GinkgoT())
		})

		It("should return error if player not in room", func() {
			mock.ExpectBegin()
			mock.ExpectRollback()

			hostID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			host := &models.Player{
				ID: hostID,
			}

			guestID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			guest := &models.Player{
				ID: guestID,
			}

			roomID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			room := &models.Room{
				ID: roomID,
				Host: models.RoomPlayer{
					ID: host.ID,
				},
				Guest: &models.RoomPlayer{
					ID:             guest.ID,
					RequestNewGame: true,
				},
			}

			mockRoomRepository.
				On("Get", ctx, roomID, true).
				Return(room, nil)

			_, err = gameEngineService.CreateGame(ctx, roomID, uuid.Nil)

			expectedErrorMessage := engine.PlayerNotInRoomErrorMessage
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(expectedErrorMessage))
			mockGameRepository.AssertExpectations(GinkgoT())
		})
	})

	Context("PlayerMakeMove", func() {
		It("should return no error if guest is in turn", func() {
			mock.ExpectBegin()
			mock.ExpectCommit()

			hostID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			host := &models.Player{
				ID: hostID,
			}

			playerID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			guest := &models.Player{
				ID: playerID,
			}

			gameID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			game := &models.Game{
				ID:    gameID,
				Phase: models.GamePhaseInProgress,
				Host: models.GamePlayer{
					ID:   host.ID,
					Mark: string(engine.XMark),
				},
				Guest: models.GamePlayer{
					ID:   guest.ID,
					Mark: string(engine.OMark),
				},
				CurrentPlayerID: guest.ID,
				Board:           "XO_______",
			}

			roomID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			room := &models.Room{
				ID: roomID,
				Host: models.RoomPlayer{
					ID: host.ID,
				},
				Guest: &models.RoomPlayer{
					ID:             guest.ID,
					RequestNewGame: true,
				},
				GameID: &game.ID,
				Phase:  models.RoomPhaseFull,
			}

			mockRoomRepository.
				On("Get", ctx, roomID, true).
				Return(room, nil)

			mockGameRepository.
				On("Get", ctx, gameID).
				Return(game, nil)

			mockGameRepository.
				On("Update", ctx, game).
				Return(nil)

			mockPlayerRepository.
				On("Get", ctx, guest.ID).
				Return(guest, nil)

			mockPlayerRepository.
				On("Get", ctx, host.ID).
				Return(host, nil)

			mockPlayerRepository.
				On("UpdateStats", ctx, tmock.Anything).
				Return(nil)

			mockRoomRepository.
				On("Update", ctx, room).
				Return(nil)

			position := 3
			err = gameEngineService.PlayerMakeMove(ctx, roomID, playerID, position)

			Expect(err).ToNot(HaveOccurred())
			mockGameRepository.AssertExpectations(GinkgoT())
		})

		It("should return no error if guest is in turn and win", func() {
			mock.ExpectBegin()
			mock.ExpectCommit()

			hostID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			host := &models.Player{
				ID: hostID,
			}

			playerID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			guest := &models.Player{
				ID: playerID,
			}

			gameID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			game := &models.Game{
				ID:    gameID,
				Phase: models.GamePhaseInProgress,
				Host: models.GamePlayer{
					ID:   host.ID,
					Mark: string(engine.XMark),
				},
				Guest: models.GamePlayer{
					ID:   guest.ID,
					Mark: string(engine.OMark),
				},
				CurrentPlayerID: guest.ID,
				Board:           "XX__O_O__",
			}

			roomID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			room := &models.Room{
				ID: roomID,
				Host: models.RoomPlayer{
					ID: host.ID,
				},
				Guest: &models.RoomPlayer{
					ID: guest.ID,
				},
				GameID: &game.ID,
				Phase:  models.RoomPhaseFull,
			}

			mockRoomRepository.
				On("Get", ctx, roomID, true).
				Return(room, nil)

			mockGameRepository.
				On("Get", ctx, gameID).
				Return(game, nil)

			mockGameRepository.
				On("Update", ctx, game).
				Return(nil)

			mockPlayerRepository.
				On("Get", ctx, guest.ID).
				Return(guest, nil)

			mockPlayerRepository.
				On("Get", ctx, host.ID).
				Return(host, nil)

			mockPlayerRepository.
				On("UpdateStats", ctx, tmock.Anything).
				Return(nil)

			mockRoomRepository.
				On("Update", ctx, room).
				Return(nil)

			position := 3
			err = gameEngineService.PlayerMakeMove(ctx, roomID, playerID, position)

			Expect(err).ToNot(HaveOccurred())
			mockGameRepository.AssertExpectations(GinkgoT())
		})

		It("should return error if player is make incorrect move", func() {
			mock.ExpectBegin()
			mock.ExpectRollback()

			hostID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			host := &models.Player{
				ID: hostID,
			}

			playerID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			guest := &models.Player{
				ID: playerID,
			}

			gameID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			game := &models.Game{
				ID:    gameID,
				Phase: models.GamePhaseInProgress,
				Host: models.GamePlayer{
					ID:   host.ID,
					Mark: string(engine.XMark),
				},
				Guest: models.GamePlayer{
					ID:   guest.ID,
					Mark: string(engine.OMark),
				},
				CurrentPlayerID: guest.ID,
				Board:           "_________",
			}

			roomID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			room := &models.Room{
				ID: roomID,
				Host: models.RoomPlayer{
					ID: host.ID,
				},
				Guest: &models.RoomPlayer{
					ID: guest.ID,
				},
				GameID: &game.ID,
				Phase:  models.RoomPhaseFull,
			}

			mockRoomRepository.
				On("Get", ctx, roomID, true).
				Return(room, nil)

			mockGameRepository.
				On("Get", ctx, gameID).
				Return(game, nil)

			position := 99
			err = gameEngineService.PlayerMakeMove(ctx, roomID, playerID, position)

			expectedErrorMessage := engine.InvalidBoardPositionErrorMessage
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(expectedErrorMessage))
			mockGameRepository.AssertExpectations(GinkgoT())
		})

		It("should return error if player is make incorrect move on occupied zone", func() {
			mock.ExpectBegin()
			mock.ExpectRollback()

			hostID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			host := &models.Player{
				ID: hostID,
			}

			playerID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			guest := &models.Player{
				ID: playerID,
			}

			gameID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			game := &models.Game{
				ID:    gameID,
				Phase: models.GamePhaseInProgress,
				Host: models.GamePlayer{
					ID:   host.ID,
					Mark: string(engine.XMark),
				},
				Guest: models.GamePlayer{
					ID:   guest.ID,
					Mark: string(engine.OMark),
				},
				CurrentPlayerID: guest.ID,
				Board:           "X________",
			}

			roomID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			room := &models.Room{
				ID: roomID,
				Host: models.RoomPlayer{
					ID: host.ID,
				},
				Guest: &models.RoomPlayer{
					ID: guest.ID,
				},
				GameID: &game.ID,
				Phase:  models.RoomPhaseFull,
			}

			mockRoomRepository.
				On("Get", ctx, roomID, true).
				Return(room, nil)

			mockGameRepository.
				On("Get", ctx, gameID).
				Return(game, nil)

			position := 1
			err = gameEngineService.PlayerMakeMove(ctx, roomID, playerID, position)

			expectedErrorMessage := engine.BoardPositionOcopiedErrorMessage
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(expectedErrorMessage))
			mockGameRepository.AssertExpectations(GinkgoT())
		})

		It("should return error if player is not in turn", func() {
			mock.ExpectBegin()
			mock.ExpectRollback()

			hostID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			host := &models.Player{
				ID: hostID,
			}

			playerID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			guest := &models.Player{
				ID: playerID,
			}

			gameID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			game := &models.Game{
				ID:    gameID,
				Phase: models.GamePhaseInProgress,
				Host: models.GamePlayer{
					ID:   host.ID,
					Mark: string(engine.XMark),
				},
				Guest: models.GamePlayer{
					ID:   guest.ID,
					Mark: string(engine.OMark),
				},
				CurrentPlayerID: host.ID,
				Board:           "X________",
			}

			roomID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			room := &models.Room{
				ID: roomID,
				Host: models.RoomPlayer{
					ID: host.ID,
				},
				Guest: &models.RoomPlayer{
					ID: guest.ID,
				},
				GameID: &game.ID,
				Phase:  models.RoomPhaseFull,
			}

			mockRoomRepository.
				On("Get", ctx, roomID, true).
				Return(room, nil)

			mockGameRepository.
				On("Get", ctx, gameID).
				Return(game, nil)

			position := 2
			err = gameEngineService.PlayerMakeMove(ctx, roomID, playerID, position)

			expectedErrorMessage := engine.PlayerNotInTurnErrorMessage
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(expectedErrorMessage))
			mockGameRepository.AssertExpectations(GinkgoT())
		})

		It("should return error if game is completed", func() {
			mock.ExpectBegin()
			mock.ExpectRollback()

			hostID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			host := &models.Player{
				ID: hostID,
			}

			playerID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			guest := &models.Player{
				ID: playerID,
			}

			gameID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			game := &models.Game{
				ID:    gameID,
				Phase: models.GamePhaseCompleted,
				Host: models.GamePlayer{
					ID:   host.ID,
					Mark: string(engine.XMark),
				},
				Guest: models.GamePlayer{
					ID:   guest.ID,
					Mark: string(engine.OMark),
				},
				CurrentPlayerID: host.ID,
				Board:           "X________",
			}

			roomID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			room := &models.Room{
				ID: roomID,
				Host: models.RoomPlayer{
					ID: host.ID,
				},
				Guest: &models.RoomPlayer{
					ID: guest.ID,
				},
				GameID: &game.ID,
				Phase:  models.RoomPhaseFull,
			}

			mockRoomRepository.
				On("Get", ctx, roomID, true).
				Return(room, nil)

			mockGameRepository.
				On("Get", ctx, gameID).
				Return(game, nil)

			position := 2
			err = gameEngineService.PlayerMakeMove(ctx, roomID, playerID, position)

			expectedErrorMessage := engine.GameCompletedErrorMessage
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(expectedErrorMessage))
			mockGameRepository.AssertExpectations(GinkgoT())
		})

		It("should return error if player not in room", func() {
			mock.ExpectBegin()
			mock.ExpectRollback()

			hostID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			host := &models.Player{
				ID: hostID,
			}

			guestID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			guest := &models.Player{
				ID: guestID,
			}

			gameID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			game := &models.Game{
				ID:    gameID,
				Phase: models.GamePhaseCompleted,
				Host: models.GamePlayer{
					ID:   host.ID,
					Mark: string(engine.XMark),
				},
				Guest: models.GamePlayer{
					ID:   guest.ID,
					Mark: string(engine.OMark),
				},
				CurrentPlayerID: host.ID,
				Board:           "X________",
			}

			roomID, err := uuid.NewV4()
			Expect(err).To(BeNil())
			room := &models.Room{
				ID: roomID,
				Host: models.RoomPlayer{
					ID: host.ID,
				},
				Guest: &models.RoomPlayer{
					ID: guest.ID,
				},
				GameID: &game.ID,
				Phase:  models.RoomPhaseFull,
			}

			mockRoomRepository.
				On("Get", ctx, roomID, true).
				Return(room, nil)

			mockGameRepository.
				On("Get", ctx, gameID).
				Return(game, nil)
			position := 2
			err = gameEngineService.PlayerMakeMove(ctx, roomID, uuid.Nil, position)

			expectedErrorMessage := engine.PlayerNotInRoomErrorMessage
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(expectedErrorMessage))
			mockGameRepository.AssertExpectations(GinkgoT())
		})
	})
})
