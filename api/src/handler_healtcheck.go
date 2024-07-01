package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// HealthcheckHandler is a function that handles requests to the /healthcheck endpoint.
func (s *Service) HealthcheckHandler(c echo.Context) error {
	payload := "OK"
	AddCustomAttributes(c, S.Any("payload", payload))
	return c.JSON(http.StatusOK, payload)
}
