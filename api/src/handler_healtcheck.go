package main

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
)

// HealthcheckHandler is a function that handles requests to the /healthcheck endpoint.
func (s *Service) HealthcheckHandler(c echo.Context) error {
	payload := []string{"OK"}
	s.AddCustomAttributes(slog.Group(s.Path, s.Any("health", payload)))
	return c.JSON(http.StatusOK, payload)
}
