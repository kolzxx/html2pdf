package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/kolzxx/html2pdf/internal/logger"
	"go.uber.org/zap"
)

type Request struct {
	Method string `json:"method"`
}

type Response struct {
	StatusCode int `json:"status_code"`
}

type Http struct {
	Request  Request  `json:"request"`
	Response Response `json:"response"`
}

type Url struct {
	Path   string `json:"path"`
	Scheme string `json:"scheme"`
	Domain string `json:"domain"`
}

// @see https://gin-gonic.com/docs/examples/custom-middleware/
func ECSMiddleware(log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// before request
		c.Next()
		// after request
		scheme := "http"
		if c.Request.TLS != nil {
			scheme = "https"
		}
		log.Info("access log",
			zap.Any("http", Http{
				Request: Request{
					Method: c.Request.Method,
				},
				Response: Response{
					StatusCode: c.Writer.Status(),
				},
			}),
			zap.Any("url", Url{
				Domain: c.Request.Host,
				Path:   c.Request.URL.Path,
				Scheme: scheme,
			}),
		)
	}
}
