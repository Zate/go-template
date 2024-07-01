package main

import (
	"embed"
	"fmt"
	"html/template"

	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	slogformatter "github.com/samber/slog-formatter"

	"log/slog"
)

//go:embed templates/*.tmpl
var templates embed.FS

// Service is the main struct for our API service
type Service struct {
	Logger *slog.Logger
	Port   int
	Server *echo.Echo
	C      *echo.Context
	Path   string
}

type CustomValidator struct {
	validator *validator.Validate
}

// Validate is a function to validate a struct using go-playground/validator/v10
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

// Run starts the service
func (s *Service) Run() error {
	e, err := s.BindRoutes()
	if err != nil {
		return err
	}
	s.Server = e

	listenAddress := ":" + fmt.Sprint(s.Port)

	e.Logger.Fatal(e.Start(listenAddress))
	return nil
}

// BindRoutes binds the routes to the service
func (s *Service) BindRoutes() (*echo.Echo, error) {
	e := echo.New()
	e.HTTPErrorHandler = e.DefaultHTTPErrorHandler

	e.HideBanner = true
	e.HidePort = true
	e.Pre(middleware.RemoveTrailingSlash())
	// Custom validator
	v := validator.New()
	v.RegisterValidation("vuln", vulnValidation)
	e.Validator = &CustomValidator{validator: v}
	config := LoggingConfig{
		WithUserAgent: true,
		WithRequestID: true,
		Message:       "REQUEST",
	}
	e.Use(s.ContextMiddleware)
	e.Use(NewLoggingMiddlewareWithConfig(s.Logger, config))
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize:    1 << 10, // 1 KB
		LogLevel:     log.Lvl(slog.LevelError),
		LogErrorFunc: s.PanicErrorFunc(),
	}))
	e.Use(middleware.RequestID())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:8080"},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowCredentials: true,
		AllowMethods:     []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))
	e.Use(middleware.Gzip())
	e.Use(middleware.Secure())

	// TODO: Add template rendering here
	t := &TemplateRegistry{
		templates: template.Must(template.ParseFS(templates, "templates/*.tmpl")),
	}
	e.Renderer = t

	// Generic and util endpoints
	root := e.Group("/")
	root.RouteNotFound("*", s.NotFoundHandler)
	root.GET("healthcheck", s.HealthcheckHandler)
	root.GET("status", s.StatusHandler)
	root.GET("debug", s.DebugHandler)

	return e, nil
}

// NewService creates a new service
func NewService(port int) (*Service, error) {
	handlerOptions := &slog.HandlerOptions{}
	handler := slog.NewJSONHandler(os.Stdout, handlerOptions)

	logger := slog.New(
		slogformatter.NewFormatterHandler(
			slogformatter.TimezoneConverter(time.UTC),
			slogformatter.TimeFormatter(time.RFC3339, nil),
		)(
			handler,
		),
	)
	newService := &Service{
		Logger: logger,
		Port:   port,
	}
	return newService, nil
}

func (s *Service) ContextMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Set("start_time", time.Now())
		s.Path = transformPath(c.Path())
		s.C = &c
		return next(c)
	}
}
