package server

import (
	"github.com/kolzxx/html2pdf/internal/controllers"
)

func (s server) RegisterRoutes() {
	hc := controllers.NewHealthControler(s.Logger)
	pc := controllers.NewHtml2PdfController(s.Logger)

	s.router.GET("/healthcheck", hc.HandleGetHealthCheck)
	v1 := s.router.Group("/v1")
	{
		v1.POST("/html2pdf", pc.HandleHttp2Pdf)
	}
}
