package middlewares

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kolzxx/html2pdf/internal/logger"
)

// @see https://gin-gonic.com/docs/examples/custom-middleware/
func CorrelationIdMiddleware(l logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		correlationHeader := "X-Correlation-ID"
		correlationID := c.Request.Header.Get(correlationHeader)
		if strings.TrimSpace(correlationID) == "" {
			correlationID = uuid.NewString()
			c.Request.Header.Add(correlationHeader, correlationID)
		}
		l.SetCorrelationId(correlationID)
	}
}
