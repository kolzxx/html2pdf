package server

import (
	"context"
	"testing"

	"github.com/kolzxx/html2pdf/internal/logger"
	"github.com/stretchr/testify/assert"
)

func TestServer(t *testing.T) {

	t.Run("NewServer", func(t *testing.T) {
		logger := logger.NewFakeLogger()

		ctx := context.Background()

		s := NewServer(ServerOptions{
			Context: ctx,
			Logger:  logger,
		})
		assert.NotNil(t, s)
	})

}
