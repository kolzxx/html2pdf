package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kolzxx/html2pdf/internal/logger"
)

type HealthControler struct {
	logger logger.Logger
}

func NewHealthControler(logger logger.Logger) HealthControler {
	return HealthControler{
		logger: logger,
	}
}

type GetHealthCheckResponse struct {
	Message string `json:"message"`
}

// @BasePath /
// @version		1.0
// @Summary Check API health
// @Schemes
// @Description Check API health
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} GetHealthCheckResponse
// @Router /healthcheck [get]
func (h *HealthControler) HandleGetHealthCheck(c *gin.Context) {
	h.logger.Info("Example log already set up!")

	// implement your health check logic here...
	c.JSON(http.StatusOK, GetHealthCheckResponse{
		Message: "ok",
	})
}
