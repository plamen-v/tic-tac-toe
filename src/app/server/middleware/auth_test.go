//go:build mock
// +build mock

package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/plamen-v/tic-tac-toe-models/models"
	"github.com/plamen-v/tic-tac-toe/src/app/server/middleware"
	"github.com/plamen-v/tic-tac-toe/src/services/auth"
	"github.com/plamen-v/tic-tac-toe/src/services/auth/mocks"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
)

var _ = Describe("Authentication", func() {
	var (
		mockAuthenticationService *mocks.MockAuthenticationService
		router                    *gin.Engine
		request                   *http.Request
	)

	testHandler := func(c *gin.Context) {
		c.String(http.StatusOK, "test")
	}

	BeforeEach(func() {
		mockAuthenticationService = new(mocks.MockAuthenticationService)
		gin.SetMode(gin.TestMode)
		router = gin.Default()
		router.Use(middleware.ErrorHandler())
		request = httptest.NewRequest(http.MethodGet, "/test", nil)
	})

	It("should return 401 if Authorization header is missing", func() {
		authz := middleware.Authentication(mockAuthenticationService)

		mockAuthenticationService.On("ValidateToken", mock.Anything).Return(nil, models.NewAuthorizationError("Invalid token"))
		router.GET("/test", authz, testHandler)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, request)
		Expect(w.Code).To(Equal(http.StatusUnauthorized))
	})

	It("should  return 401 if token is successful", func() {
		authz := middleware.Authentication(mockAuthenticationService)

		invalidTokenType := &jwt.Token{
			Claims: &jwt.MapClaims{
				"exp": time.Now().Add(time.Hour).Unix(),
			},
			Valid: true,
		}

		mockAuthenticationService.On("ValidateToken", mock.Anything).Return(invalidTokenType, nil)
		router.GET("/test", authz, testHandler)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, request)
		Expect(w.Code).To(Equal(http.StatusUnauthorized))
	})

	It("should call next handler if authorization is successful", func() {
		authz := middleware.Authentication(mockAuthenticationService)

		mockToken := &jwt.Token{
			Claims: &auth.ExtendedClaims{
				PlayerID: 1,
			},
			Valid: true,
		}

		mockAuthenticationService.On("ValidateToken", mock.Anything).Return(mockToken, nil)
		router.GET("/test", authz, testHandler)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, request)
		Expect(w.Code).To(Equal(http.StatusOK))
	})
})
