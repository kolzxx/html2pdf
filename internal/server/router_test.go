package server

import (
	"context"
	"testing"

	"github.com/kolzxx/html2pdf/internal/logger"
	"github.com/stretchr/testify/assert"
)

func TestRouter(t *testing.T) {

	t.Run("Routers", func(t *testing.T) {
		logger := logger.NewFakeLogger()

		ctx := context.Background()

		s := NewServer(ServerOptions{
			Context: ctx,
			Logger:  logger,
		})

		s.SetupServer()
		s.SetupSwagger()
		s.SetupMiddlewares()
		s.RegisterRoutes()
		assert.NotNil(t, s)
	})

}
