package main

import (
	"context"

	"github.com/kolzxx/html2pdf/internal/logger"
	"github.com/kolzxx/html2pdf/internal/server"
)

// @title           payroll-loan.html2pdf
// @version         1.0
// @description.markdown

// @BasePath  /

// @schemes https http
func main() {
	ctx := context.Background()

	s := server.NewServer(server.ServerOptions{
		Context: ctx,
		Logger:  logger.NewEcsLogger(ctx),
	})
	s.Start()
}
