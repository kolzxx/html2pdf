package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/kolzxx/html2pdf/api"
	"github.com/kolzxx/html2pdf/configs"

	"github.com/kolzxx/html2pdf/internal/logger"
	"github.com/kolzxx/html2pdf/internal/middlewares"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server interface {
	Start()
	SetupServer()
	SetupSwagger()
	SetupMiddlewares()
	RegisterRoutes()
}

type ServerOptions struct {
	Logger  logger.Logger
	Context context.Context
}

type server struct {
	router *gin.Engine
	ServerOptions
}

func NewServer(sopt ServerOptions) Server {
	return server{
		router:        gin.New(),
		ServerOptions: sopt,
	}
}

func (s server) Start() {
	s.SetupServer()
	s.SetupSwagger()
	s.SetupMiddlewares()
	s.RegisterRoutes()

	// s.engine.Run(port)

	s.startWithGracefulShutdown()
}

func (s server) SetupServer() {
	// @see https://gin-gonic.com/docs/examples/define-format-for-the-log-of-routes/
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, _ string, _ int) {
		s.Logger.Info(fmt.Sprintf("Mapped [%v %v] route", httpMethod, absolutePath))
	}

	// @see https://pkg.go.dev/github.com/gin-gonic/gin#Engine.SetTrustedProxies
	s.router.SetTrustedProxies(nil)
}

func (s server) SetupMiddlewares() {
	s.router.Use(middlewares.CorrelationIdMiddleware(s.Logger))
	s.router.Use(middlewares.ECSMiddleware(s.Logger))
	s.router.Use(gin.Recovery())
}

func (s server) SetupSwagger() {
	if configs.GetConfig().Swagger.Enabled {
		s.router.StaticFile("/swagger.json", "./api/swagger.json")
		s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}
}

// @see https://gin-gonic.com/docs/examples/graceful-restart-or-stop/
func (s server) startWithGracefulShutdown() {
	port := fmt.Sprintf(":%s", configs.GetConfig().Server.Port)

	srv := &http.Server{
		Addr:    port,
		Handler: s.router,
	}

	go func() {
		s.Logger.Info(fmt.Sprintf("Starting server on port %s...", port))
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.Logger.Error(fmt.Sprintf("Start server error! %s", err.Error()))
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	sigs := make(chan os.Signal, 1)

	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigs

	s.Logger.Info(fmt.Sprintf("Server Shutdown [%s]...", sig))

	ctx, cancel := context.WithTimeout(s.Context, 5*time.Second)

	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		s.Logger.Error("Server Shutdown:", err)
		os.Exit(1)
	}

	<-ctx.Done()

	s.Logger.Info("Server Shutdown: Done!")
}
