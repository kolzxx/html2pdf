package controllers_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kolzxx/html2pdf/internal/controllers"
	"github.com/kolzxx/html2pdf/internal/logger"
	"github.com/stretchr/testify/assert"
)

func TestHandleGetHealthCheck(t *testing.T) {
	t.Parallel()

	t.Run("HandleGetHealthCheck", func(t *testing.T) {
		logger := logger.NewFakeLogger()

		w := httptest.NewRecorder()
		_, r := gin.CreateTestContext(w)

		hc := controllers.NewHealthControler(logger)

		path := "/ping"
		r.GET(path, hc.HandleGetHealthCheck)

		req, _ := http.NewRequest("GET", path, nil)
		r.ServeHTTP(w, req)

		expectedResponse := `{"message":"ok"}`
		responseData, _ := io.ReadAll(w.Body)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, expectedResponse, string(responseData))
	})
}
