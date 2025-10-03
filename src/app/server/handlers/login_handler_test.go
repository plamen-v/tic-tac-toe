package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
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

		request, err := http.NewRequest("POST", "/test", nil)
		Expect(err).To(BeNil())

		response := httptest.NewRecorder()
		router.POST("/test", loginHandler)
		router.ServeHTTP(response, request)

		Expect(response.Code).To(Equal(http.StatusBadRequest))
	})

	It("should return 401 if player is invalid", func() {
		loginHandler := handlers.LoginHandler(mockAuthenticationService)
		mockAuthenticationService.On("Authenticate", mock.Anything, mock.Anything, mock.Anything).Return(nil, "", models.NewAuthorizationErrorf("invalid password"))

		loginRequest := models.LoginRequest{
			Login:    "login",
			Password: "password",
		}
		requestBody, err := json.Marshal(loginRequest)
		Expect(err).To(BeNil())
		request, err := http.NewRequest("POST", "/test", bytes.NewBuffer(requestBody))
		Expect(err).To(BeNil())
		request.Header.Set("Content-Type", "application/json")

		response := httptest.NewRecorder()
		router.POST("/test", loginHandler)
		router.ServeHTTP(response, request)

		Expect(response.Code).To(Equal(http.StatusUnauthorized))
	})

	It("should return 200 if player is valid", func() {
		loginHandler := handlers.LoginHandler(mockAuthenticationService)

		loginRequest := models.LoginRequest{
			Login:    "login",
			Password: "password",
		}

		requestBody, err := json.Marshal(loginRequest)
		Expect(err).To(BeNil())
		request, err := http.NewRequest("POST", "/test", bytes.NewBuffer(requestBody))
		Expect(err).To(BeNil())
		request.Header.Set("Content-Type", "application/json")

		response := httptest.NewRecorder()
		router.POST("/test", loginHandler)
		mockAuthenticationService.On("Authenticate", mock.Anything, mock.Anything, mock.Anything).Return(&models.Player{}, "valid-token", nil)

		router.ServeHTTP(response, request)
		Expect(response.Code).To(Equal(http.StatusOK))
	})
})
