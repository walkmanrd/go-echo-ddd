package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/walkmanrd/assessment/types"
)

// HealthCheckController is a struct for health check controller
type HealthCheckController struct{}

// GET /health-check
// Index is a function to get health check
func (c *HealthCheckController) Index(e echo.Context) error {
	return e.JSON(http.StatusOK, &types.HealthCheck{
		Message: "success",
		Status:  true,
	})
}
