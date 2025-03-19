package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kolzxx/html2pdf/internal/logger"
	"github.com/kolzxx/html2pdf/internal/middlewares"
	"github.com/stretchr/testify/assert"
)

func TestCorrelationIdMiddleware(t *testing.T) {
	t.Parallel()

	t.Run("Checks if the correlation ID has been set on the logger", func(t *testing.T) {
		r := gin.Default()
		logger := logger.NewFakeLogger()

		r.Use(middlewares.CorrelationIdMiddleware(logger))

		r.GET("/ping", func(c *gin.Context) {
			c.String(200, "pong")
		})

		w := httptest.NewRecorder()

		req, _ := http.NewRequest("GET", "/ping", nil)
		r.ServeHTTP(w, req)

		assert.NotNil(t, logger.GetCorrelationId())
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
