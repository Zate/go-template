package main

import (
	"log/slog"
	"os"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// getEnv is a function to get an environment variable or return a default value
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// getDebugInfo is a function to get debug information from runtime.Caller to add to the logging output
func getDebugInfo() (string, string, int) {
	pc, filename, line, _ := runtime.Caller(1)
	return runtime.FuncForPC(pc).Name(), filename, line
}

// Log is a function to log messages to the console
func (s *Service) Log(level string, msg string, args ...any) {
	if Debug {
		_, filename, line, _ := runtime.Caller(1)
		args = append(args, "filename", filename, "line", line)
	}

	switch level {
	case "debug":
		if Debug {
			s.Logger.Info(msg, args...)
		}
	case "info":
		s.Logger.Info(msg, args...)
	case "warn":
		s.Logger.Warn(msg, args...)
	case "error":
		s.Logger.Error(msg, args...)
	default:
		s.Logger.Info(msg, args...)
	}
}

// Any is a function to generate a slog.Attr from a key value pair
// it will return a slog.Attr with the key and value
func (s *Service) Any(key string, value any) slog.Attr {
	return slog.Any(key, value)
}

// LogHTTPErr is a function to generate a slog.Attr from a given *echo.HTTPError
// it will return a slog.Attr with the key "error" and the value of the error
func (s *Service) LogHTTPErr(name string, err *echo.HTTPError) slog.Attr {
	if err.Internal != nil {
		return slog.Group("error", s.Any("function", name), s.Any("err.message", err.Message), s.Any("err.code", err.Code), s.Any("err.internal", err.Internal))
	}
	return slog.Group("error", s.Any("function", name), s.Any("err.message", err.Message), s.Any("err.code", err.Code))
}

// TrackTime is a function to track the time elapsed since a given time
func TrackTime(pre time.Time) time.Duration {
	return time.Duration(time.Since(pre).Milliseconds())
}

// logCustomAttr is a function to log custom attributes to the log
func (s *Service) logCustomAttr(groupName string, data map[string]interface{}) {
	var attrs []any
	for key, value := range data {
		attrs = append(attrs, s.Any(key, value))
	}

	if os.Getenv("MY_ENVTYPE") != "staging" {
		if data["error"] != nil {
			pc, filename, line, _ := runtime.Caller(1)
			funcName := runtime.FuncForPC(pc).Name()
			attrs = append(attrs, s.Any("filename", filename), s.Any("line", line), s.Any("function", funcName))
		}
	}

	group := slog.Group(groupName, attrs...)

	AddCustomAttributes(*s.C, group)
}

// logCustomAttr is a function to log custom attributes to the log
func (s *Service) genCustomAttr(groupName string, data map[string]interface{}) slog.Attr {
	var attrs []any
	for key, value := range data {
		attrs = append(attrs, s.Any(key, value))
	}

	if os.Getenv("MY_ENVTYPE") != "staging" {
		if data["error"] != nil {
			pc, filename, line, _ := runtime.Caller(1)
			funcName := runtime.FuncForPC(pc).Name()
			attrs = append(attrs, s.Any("filename", filename), s.Any("line", line), s.Any("function", funcName))
		}
	}
	return slog.Group(groupName, attrs...)
}

func (s *Service) AddCustomAttributes(attr slog.Attr) {
	c := *s.C
	v := c.Get(customAttributesCtxKey)
	if v == nil {
		c.Set(customAttributesCtxKey, []slog.Attr{attr})
		return
	}

	switch attrs := v.(type) {
	case []slog.Attr:
		c.Set(customAttributesCtxKey, append(attrs, attr))
	}
	s.C = &c
}

// Utility function to transform the path string
func transformPath(path string) string {
	if path == "/" {
		return "root"
	}

	segments := strings.Split(path, "/")
	for i, segment := range segments {
		if segment != "" {
			if i == 0 {
				segments[i] = strings.ToLower(segment)
				continue
			}
			segments[i] = cases.Title(language.English).String(segment)
		}
	}

	return strings.Join(segments, "")
}

func vulnValidation(fl validator.FieldLevel) bool {
	vuln := fl.Field().String()
	regex := `^VULN-[0-9]+$`
	re := regexp.MustCompile(regex)
	return re.MatchString(vuln)
}
