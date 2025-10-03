package middleware_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/plamen-v/tic-tac-toe-models/models"
	"github.com/plamen-v/tic-tac-toe/src/app/server/middleware"
	"github.com/plamen-v/tic-tac-toe/src/services/logger"
	"github.com/plamen-v/tic-tac-toe/src/services/logger/mocks"

	"github.com/stretchr/testify/mock"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Logger", func() {
	var (
		mockLoggerService *mocks.MockLoggerService
		router            *gin.Engine
		request           *http.Request
	)

	errorMsg := "test error"
	errorHandler := func(c *gin.Context) {
		_ = c.Error(models.NewValidationError(errorMsg))
		c.Status(http.StatusBadRequest)
		c.Abort()
	}

	infoHandler := func(c *gin.Context) {
		c.String(http.StatusOK, "test")
	}

	BeforeEach(func() {
		mockLoggerService = new(mocks.MockLoggerService)
		gin.SetMode(gin.TestMode)
		router = gin.Default()
		router.Use(middleware.ErrorHandler())
		request = httptest.NewRequest(http.MethodGet, "/test", nil)
	})

	It("should log an Error when context has errors", func() {
		logz := middleware.Logger(mockLoggerService)

		mockLoggerService.On("Error", "request", mock.Anything).Run(func(args mock.Arguments) {
			fields := args.Get(1).([]logger.Field)
			Expect(len(fields)).To(Equal(6))
		}).Return()

		router.GET("/test", logz, errorHandler)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, request)
		Expect(w.Code).To(Equal(http.StatusBadRequest))
		mockLoggerService.AssertCalled(GinkgoT(), "Error", "request", mock.Anything)
	})

	It("should log an Info when context has no errors", func() {
		logz := middleware.Logger(mockLoggerService)

		mockLoggerService.On("Info", "request", mock.Anything).Run(func(args mock.Arguments) {
			fields := args.Get(1).([]logger.Field)
			Expect(len(fields)).To(Equal(5))
		}).Return()

		router.GET("/test", logz, infoHandler)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, request)
		Expect(w.Code).To(Equal(http.StatusOK))
		mockLoggerService.AssertCalled(GinkgoT(), "Info", "request", mock.Anything)
	})
})
