package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Healthz godoc
// @Summary      Health check
// @Description  Returns server health status
// @Tags         health
// @Produce      json
// @Success      200  {object}  map[string]string
// @Router       /healthz [get]
func Healthz(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}
