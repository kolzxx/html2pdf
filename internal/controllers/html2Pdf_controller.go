package controllers

import (
	"context"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/kolzxx/html2pdf/internal/dtos"
	"github.com/kolzxx/html2pdf/internal/interfaces"
	"github.com/kolzxx/html2pdf/internal/logger"
	"github.com/kolzxx/html2pdf/internal/services"
)

type Http2PdfController struct {
	html2PdfService interfaces.Html2PdfServiceInterface
	chromedpService *services.ChromedpService
	logger          logger.Logger
}

func NewHtml2PdfController(logger logger.Logger) *Http2PdfController {
	numCPUS := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPUS)

	var app Http2PdfController
	app.logger = logger
	app.chromedpService = services.NewChromedpService(context.Background(), logger)
	err := app.chromedpService.RunChromeDp()
	if err != nil {
		logger.Error("Error running chromedp - Run Chromedp", err)
	}

	app.html2PdfService = services.NewHtml2PdfService(logger, app.chromedpService)
	return &app
}

// @Summary API Convert html to pdf
// @Description Retrieve the pdf file of a html
// @Tags HTML PDF
// @Produce json
// @Version 1.0
// @Param Request body dtos.HtmlRequest true "The input HtmlRequest struct"
// @Success 200 {object} dtos.BaseResponse "success"
// @Failure 400 {object} dtos.BaseResponse "error"
// @Router /v1/html2pdf [post]
func (h *Http2PdfController) HandleHttp2Pdf(c *gin.Context) {
	h.logger.Info("Http2Pdf - Started")
	var request dtos.HtmlRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, dtos.WithError(err.Error(), 40))
		return
	}

	response, err := h.html2PdfService.HtmlToPdf(request)

	if err != nil {
		c.JSON(500, dtos.WithError(err.Error(), 40))
		return
	}

	if len(response.Content) == 0 {
		c.JSON(500, dtos.WithError("Timeout", 40))
		return
	}

	c.JSON(200, dtos.WithSuccess("html converted successfully", 200, response))
	h.logger.Info("Http2Pdf - Finished")
}
