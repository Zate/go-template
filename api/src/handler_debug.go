package main

import (
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/labstack/echo/v4"
	"golang.org/x/exp/slices"
)

// TODO: Set this to only allow debug requests from localhost and only when running locally under comtek/air

// DebugHandler is a function that handles requests to the /debug endpoint.
func (s *Service) DebugHandler(c echo.Context) error {
	env := os.Environ()
	// sort the []string alphabetically using slices.Sort
	slices.Sort(sort.StringSlice(env))

	payload := s.GenDebugInfo(env)
	return c.JSONPretty(http.StatusOK, payload, "  ")
}

// GenDebugInfo is a function that turns a []string of key=value pairs into a map[string]string
func (s *Service) GenDebugInfo(env []string) map[string]string {
	info := make(map[string]string)
	for _, pair := range env {
		kv := strings.Split(pair, "=")
		// If the env var ends with 'TOKEN' or 'KEY', then replace the value with 'REDACTED'
		if strings.HasSuffix(kv[0], "_TOKEN") || strings.HasSuffix(kv[0], "_KEY") {
			info[kv[0]] = "REDACTED"
			continue
		}
		info[kv[0]] = kv[1]
	}
	return info
}
