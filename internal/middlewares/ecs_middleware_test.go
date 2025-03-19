package middlewares_test

import (
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kolzxx/html2pdf/internal/logger"
	"github.com/kolzxx/html2pdf/internal/middlewares"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestECSMiddleware(t *testing.T) {
	t.Parallel()

	verifyLogHasBeenRecorded := func(t *testing.T, w *httptest.ResponseRecorder, expectedFields []zapcore.Field, logger *logger.FakeLogger) {
		t.Helper()

		assert.ElementsMatch(t, expectedFields, logger.Fields)
		assert.Equal(t, "access log", logger.Message)
		assert.Greater(t, len(logger.Fields), 0)
		assert.Equal(t, http.StatusOK, w.Code)
	}

	t.Run("Checks if the access log has been recorded", func(t *testing.T) {
		w := httptest.NewRecorder()
		_, r := gin.CreateTestContext(w)
		logger := logger.NewFakeLogger().(*logger.FakeLogger)
		path := "/ping"

		r.Use(middlewares.ECSMiddleware(logger))
		r.GET(path, func(c *gin.Context) {
			c.String(200, "pong")
		})

		req, _ := http.NewRequest("GET", path, nil)
		r.ServeHTTP(w, req)

		expectedFields := []zap.Field{
			zap.Any("http", middlewares.Http{
				Request: middlewares.Request{
					Method: "GET",
				},
				Response: middlewares.Response{
					StatusCode: http.StatusOK,
				},
			}),
			zap.Any("url", middlewares.Url{
				Domain: "",
				Path:   path,
				Scheme: "http",
			}),
		}

		verifyLogHasBeenRecorded(t, w, expectedFields, logger)
	})

	t.Run("Checks if the access log has been recorded (with TLS)", func(t *testing.T) {
		w := httptest.NewRecorder()
		_, r := gin.CreateTestContext(w)
		logger := logger.NewFakeLogger().(*logger.FakeLogger)
		path := "/ping"

		r.Use(middlewares.ECSMiddleware(logger))
		r.GET(path, func(c *gin.Context) {
			c.Request.TLS = &tls.ConnectionState{}
			c.String(http.StatusOK, "pong")
		})

		req, _ := http.NewRequest("GET", path, nil)
		r.ServeHTTP(w, req)

		expectedFields := []zap.Field{
			zap.Any("http", middlewares.Http{
				Request: middlewares.Request{
					Method: "GET",
				},
				Response: middlewares.Response{
					StatusCode: http.StatusOK,
				},
			}),
			zap.Any("url", middlewares.Url{
				Domain: "",
				Path:   path,
				Scheme: "https",
			}),
		}

		verifyLogHasBeenRecorded(t, w, expectedFields, logger)
	})
}
