package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	. "github.com/onsi/ginkgo/v2"
	"github.com/plamen-v/tic-tac-toe-models/models"
	"github.com/plamen-v/tic-tac-toe/src/app/server/handlers"
	"github.com/plamen-v/tic-tac-toe/src/app/server/middleware"
	"github.com/plamen-v/tic-tac-toe/src/services/auth/mocks"
	"github.com/stretchr/testify/mock"

	. "github.com/onsi/gomega"
)

var _ = Describe("LoginHandler", func() {
	var (
		mockAuthenticationService *mocks.MockAuthenticationService
		router                    *gin.Engine
	)

	BeforeEach(func() {
		mockAuthenticationService = new(mocks.MockAuthenticationService)
		gin.SetMode(gin.TestMode)
		router = gin.Default()
		router.Use(middleware.ErrorHandler())
	})

	It("should return 400 if request is invalid", func() {
		loginHandler := handlers.LoginHandler(mockAuthenticationService)

		player := &models.Player{
			ID: uuid.Must(uuid.NewV4()),
		}
		validTokenStr := "valid-token"

		mockAuthenticationService.On("Authenticate", mock.Anything, mock.Anything, mock.Anything).Return(player, nil)
		mockAuthenticationService.On("CreateToken", player, mock.Anything).Return(validTokenStr, nil)

		response := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(response)

		req, err := http.NewRequest("POST", "/login", nil)
		Expect(err).To(BeNil())

		c.Request = req
		loginHandler(c)
		Expect(response.Code).To(Equal(http.StatusBadRequest))
	})

	It("should return 200 if player is valid", func() {
		loginHandler := handlers.LoginHandler(mockAuthenticationService)

		player := &models.Player{
			ID: uuid.Must(uuid.NewV4()),
		}
		validTokenStr := "valid-token"

		mockAuthenticationService.On("Authenticate", mock.Anything, mock.Anything, mock.Anything).Return(player, nil)
		mockAuthenticationService.On("CreateToken", player, mock.Anything).Return(validTokenStr, nil)

		response := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(response)

		loginRequest := models.LoginRequest{
			Login:    "login",
			Password: "password",
		}
		requestBody, err := json.Marshal(loginRequest)
		Expect(err).To(BeNil())

		req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(requestBody))
		Expect(err).To(BeNil())

		c.Request = req
		loginHandler(c)
		Expect(response.Code).To(Equal(http.StatusOK))

		var loginResponse models.LoginResponse
		err = json.Unmarshal(response.Body.Bytes(), &loginResponse)
		Expect(err).To(BeNil())

		Expect(loginResponse.Player.ID).To(Equal(player.ID))
	})
})
