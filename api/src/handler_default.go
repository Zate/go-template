package main

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

// DefaultInfo represents the response structure for the default page.
type DefaultInfo struct {
	Title          string `json:"title"`
	Environment    string `json:"environment"`
	InstanceType   string `json:"instance_type"`
	Service        string `json:"service"`
	ServiceVersion string `json:"service_version"`
}

// DefaultError represents a default JSON error response
type DefaultError struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

// NotFoundHandler is a function that handles requests that are not supported or handled by any other handler
func (s *Service) NotFoundHandler(c echo.Context) error {
	payload := DefaultError{
		Error: "Not Supported",
		Code:  http.StatusBadRequest,
	}
	return c.JSONPretty(http.StatusOK, payload, "  ")

}

// DefaultInfo is a function that returns the default info for the default page.
func (s *Service) DefaultInfo(c echo.Context) DefaultInfo {
	var payload DefaultInfo
	payload.Title = "Welcome to the default page"
	payload.Environment = s.GetEnv("MY_ENVTYPE", "local")
	payload.InstanceType = s.GetEnv("MY_INSTANCETYPE", "local")
	payload.Service = s.GetEnv("SERVICE_NAME", "default")
	payload.ServiceVersion = s.GetEnv("SERVICE_VERSION", "v1")
	return payload
}

// GetEnv is a function that takes in an environment variable to check and the default to return if it's not found.
// returns string
func (s *Service) GetEnv(envVar string, defaultVal string) string {
	env := os.Getenv(envVar)
	if env == "" {
		return defaultVal
	}
	return env
}
