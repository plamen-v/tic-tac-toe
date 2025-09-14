package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/plamen-v/tic-tac-toe/src/services/logger"
)

func Logger(l logger.LoggerService) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)

		fields := []logger.Field{
			logger.String("method", c.Request.Method),
			logger.String("path", strconv.Quote(c.Request.URL.Path)),
			logger.Int("status", c.Writer.Status()),
			logger.Duration("duration", duration),
			logger.String("client_ip", strconv.Quote(c.ClientIP())),
		}

		if len(c.Errors) > 0 {
			fields = append(fields, logger.String("errors", c.Errors.String()))
			l.Error("request", fields...)
		} else {
			l.Info("request", fields...)
		}
	}
}
