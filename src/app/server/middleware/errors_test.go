package middleware_test

import (
	"encoding/json"
	"errors"
	"net/http"

	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/plamen-v/tic-tac-toe-models/models"
	"github.com/plamen-v/tic-tac-toe/src/app/server/middleware"

	. "github.com/onsi/ginkgo/v2"

	. "github.com/onsi/gomega"
)

var _ = Describe("ErrorHandler", func() {
	var (
		router  *gin.Engine
		request *http.Request
	)

	errorMsg := "test error"
	notFoundErrorHandler := func(c *gin.Context) {
		_ = c.Error(models.NewNotFoundError(errorMsg))
		c.Abort()
	}
	validationErrorHandler := func(c *gin.Context) {
		_ = c.Error(models.NewValidationError(errorMsg))
		c.Abort()
	}
	authorizationErrorHandler := func(c *gin.Context) {
		_ = c.Error(models.NewAuthorizationError(errorMsg))
		c.Abort()
	}
	genericErrorHandler := func(c *gin.Context) {
		_ = c.Error(models.NewGenericError(errorMsg))
		c.Abort()
	}

	unknownErrorHandler := func(c *gin.Context) {
		_ = c.Error(errors.New(errorMsg))
		c.Abort()
	}

	okHandler := func(c *gin.Context) {
		c.String(http.StatusOK, "test")
	}

	BeforeEach(func() {

		gin.SetMode(gin.TestMode)
		router = gin.Default()
		request = httptest.NewRequest(http.MethodGet, "/test", nil)
	})

	It("should return NotFoundError error", func() {
		errorz := middleware.ErrorHandler()

		router.GET("/test", errorz, notFoundErrorHandler)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, request)
		Expect(w.Code).To(Equal(http.StatusNotFound))

		var resp models.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		Expect(err).To(BeNil())

		Expect(*&resp.Code).To(Equal(string(models.NotFoundErrorCode)))
	})

	It("should return ValidationError error", func() {
		errorz := middleware.ErrorHandler()

		router.GET("/test", errorz, validationErrorHandler)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, request)
		Expect(w.Code).To(Equal(http.StatusBadRequest))

		var resp models.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		Expect(err).To(BeNil())

		Expect(*&resp.Code).To(Equal(string(models.BadRequestErrorCode)))
	})

	It("should return AuthorizationError error", func() {
		errorz := middleware.ErrorHandler()

		router.GET("/test", errorz, authorizationErrorHandler)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, request)
		Expect(w.Code).To(Equal(http.StatusUnauthorized))

		var resp models.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		Expect(err).To(BeNil())

		Expect(*&resp.Code).To(Equal(string(models.UnauthorizedErrorCode)))
	})

	It("should return GenericError error", func() {
		errorz := middleware.ErrorHandler()

		router.GET("/test", errorz, genericErrorHandler)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, request)
		Expect(w.Code).To(Equal(http.StatusInternalServerError))

		var resp models.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		Expect(err).To(BeNil())

		Expect(*&resp.Message).To(Equal(string(models.InternalServerErrorMessage)))
	})

	It("should return status 500", func() {
		errorz := middleware.ErrorHandler()

		router.GET("/test", errorz, unknownErrorHandler)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, request)
		Expect(w.Code).To(Equal(http.StatusInternalServerError))
	})

	It("should return status 200", func() {
		errorz := middleware.ErrorHandler()

		router.GET("/test", errorz, okHandler)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, request)
		Expect(w.Code).To(Equal(http.StatusOK))
	})

})
