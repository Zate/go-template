package main

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"go.opentelemetry.io/otel/trace"
)

const (
	customAttributesCtxKey = "slog-echo.custom-attributes"
)

var (
	RequestBodyMaxSize  = 64 * 1024 // 64KB
	ResponseBodyMaxSize = 64 * 1024 // 64KB

	HiddenRequestHeaders = map[string]struct{}{
		"authorization": {},
		"cookie":        {},
		"set-cookie":    {},
		"x-auth-token":  {},
		"x-csrf-token":  {},
		"x-xsrf-token":  {},
	}
	HiddenResponseHeaders = map[string]struct{}{
		"set-cookie": {},
	}
)

type LoggingConfig struct {
	DefaultLevel     slog.Level
	ClientErrorLevel slog.Level
	ServerErrorLevel slog.Level

	WithUserAgent      bool
	WithRequestID      bool
	WithRequestBody    bool
	WithRequestHeader  bool
	WithResponseBody   bool
	WithResponseHeader bool
	WithSpanID         bool
	WithTraceID        bool

	Message string

	Filters []Filter
}

// New returns a echo.MiddlewareFunc (middleware) that logs requests using slog.
//
// Requests with errors are logged using slog.Error().
// Requests without errors are logged using slog.Info().
func NewLoggingMiddleware(logger *slog.Logger) echo.MiddlewareFunc {
	return NewLoggingMiddlewareWithConfig(logger, LoggingConfig{
		DefaultLevel:     slog.LevelInfo,
		ClientErrorLevel: slog.LevelWarn,
		ServerErrorLevel: slog.LevelError,

		WithUserAgent:      false,
		WithRequestID:      true,
		WithRequestBody:    false,
		WithRequestHeader:  false,
		WithResponseBody:   false,
		WithResponseHeader: false,
		WithSpanID:         false,
		WithTraceID:        false,

		Message: "Incoming request",

		Filters: []Filter{},
	})
}

// NewWithFilters returns a echo.MiddlewareFunc (middleware) that logs requests using slog.
//
// Requests with errors are logged using slog.Error().
// Requests without errors are logged using slog.Info().
func NewLoggingMiddlewareWithFilters(logger *slog.Logger, filters ...Filter) echo.MiddlewareFunc {
	return NewLoggingMiddlewareWithConfig(logger, LoggingConfig{
		DefaultLevel:     slog.LevelInfo,
		ClientErrorLevel: slog.LevelWarn,
		ServerErrorLevel: slog.LevelError,

		WithUserAgent:      false,
		WithRequestID:      true,
		WithRequestBody:    false,
		WithRequestHeader:  false,
		WithResponseBody:   false,
		WithResponseHeader: false,
		WithSpanID:         false,
		WithTraceID:        false,

		Message: "Incoming request",

		Filters: filters,
	})
}

// NewWithConfig returns a echo.HandlerFunc (middleware) that logs requests using slog.
func NewLoggingMiddlewareWithConfig(logger *slog.Logger, config LoggingConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			req := c.Request()
			res := c.Response()
			start := time.Now()
			path := req.URL.Path
			query := req.URL.RawQuery

			params := map[string]string{}
			for i, k := range c.ParamNames() {
				params[k] = c.ParamValues()[i]
			}

			// dump request body
			br := newBodyReader(req.Body, RequestBodyMaxSize, config.WithRequestBody)
			req.Body = br

			// dump response body
			bw := newBodyWriter(res.Writer, ResponseBodyMaxSize, config.WithResponseBody)
			res.Writer = bw

			err = next(c)

			if err != nil {
				c.Error(err)
			}

			status := res.Status
			method := req.Method
			host := req.Host
			route := c.Path()
			end := time.Now()
			latency := end.Sub(start)
			userAgent := req.UserAgent()
			ip := c.RealIP()
			referer := c.Request().Referer()

			httpErr := new(echo.HTTPError)
			if err != nil && errors.As(err, &httpErr) {
				status = httpErr.Code
				if msg, ok := httpErr.Message.(string); ok {
					err = errors.New(msg)
				}
			}

			baseAttributes := []slog.Attr{}

			requestAttributes := []slog.Attr{
				slog.Time("time", start),
				slog.String("method", method),
				slog.String("host", host),
				slog.String("path", path),
				slog.String("query", query),
				slog.Any("params", params),
				slog.String("route", route),
				slog.String("ip", ip),
				slog.String("referer", referer),
			}

			responseAttributes := []slog.Attr{
				slog.Time("time", end),
				slog.Duration("latency", latency),
				slog.Int("status", status),
			}

			if config.WithRequestID {
				requestID := req.Header.Get(echo.HeaderXRequestID)
				if requestID == "" {
					requestID = res.Header().Get(echo.HeaderXRequestID)
				}
				if requestID != "" {
					baseAttributes = append(baseAttributes, slog.String("id", requestID))
				}
			}

			// otel
			if config.WithTraceID {
				traceID := trace.SpanFromContext(c.Request().Context()).SpanContext().TraceID().String()
				baseAttributes = append(baseAttributes, slog.String("trace-id", traceID))
			}
			if config.WithSpanID {
				spanID := trace.SpanFromContext(c.Request().Context()).SpanContext().SpanID().String()
				baseAttributes = append(baseAttributes, slog.String("span-id", spanID))
			}

			// request body
			requestAttributes = append(requestAttributes, slog.Int("length", br.bytes))
			if config.WithRequestBody {
				requestAttributes = append(requestAttributes, slog.String("body", br.body.String()))
			}

			// request headers
			if config.WithRequestHeader {
				for k, v := range c.Request().Header {
					if _, found := HiddenRequestHeaders[strings.ToLower(k)]; found {
						continue
					}
					requestAttributes = append(requestAttributes, slog.Group("header", slog.Any(k, v)))
				}
			}

			if config.WithUserAgent {
				requestAttributes = append(requestAttributes, slog.String("user-agent", userAgent))
			}

			xForwardedFor, ok := c.Get(echo.HeaderXForwardedFor).(string)
			if ok && len(xForwardedFor) > 0 {
				ips := lo.Map(strings.Split(xForwardedFor, ","), func(ip string, _ int) string {
					return strings.TrimSpace(ip)
				})
				requestAttributes = append(requestAttributes, slog.Any("x-forwarded-for", ips))
			}

			// response body body
			responseAttributes = append(responseAttributes, slog.Int("length", bw.bytes))
			if config.WithResponseBody {
				responseAttributes = append(responseAttributes, slog.String("body", bw.body.String()))
			}

			// response headers
			if config.WithResponseHeader {
				for k, v := range c.Response().Header() {
					if _, found := HiddenResponseHeaders[strings.ToLower(k)]; found {
						continue
					}
					responseAttributes = append(responseAttributes, slog.Group("header", slog.Any(k, v)))
				}
			}

			attributes := append(
				[]slog.Attr{
					{
						Key:   "request",
						Value: slog.GroupValue(requestAttributes...),
					},
					{
						Key:   "response",
						Value: slog.GroupValue(responseAttributes...),
					},
				},
				baseAttributes...,
			)

			// custom context values
			if v := c.Get(customAttributesCtxKey); v != nil {
				switch attrs := v.(type) {
				case []slog.Attr:
					attributes = append(attributes, attrs...)
				}
			}

			for _, filter := range config.Filters {
				if !filter(c) {
					return
				}
			}

			for _, attr := range attributes {
				if attr.Key == "msg" {
					config.Message = attr.Value.String()
					// remove this attr from the list
					attributes = lo.Filter(attributes, func(a slog.Attr, _ int) bool {
						return a.Key != "msg"
					})
				}
			}

			level := config.DefaultLevel
			msg := config.Message
			if status >= http.StatusInternalServerError {
				level = config.ServerErrorLevel
				if err != nil {
					msg = err.Error()
				} else {
					msg = http.StatusText(status)
				}
			} else if status >= http.StatusBadRequest && status < http.StatusInternalServerError {
				level = config.ClientErrorLevel
				if err != nil {
					msg = err.Error()
				} else {
					msg = http.StatusText(status)
				}
			}

			logger.LogAttrs(c.Request().Context(), level, msg, attributes...)

			return
		}
	}
}

func AddCustomAttributes(c echo.Context, attr slog.Attr) {
	v := c.Get(customAttributesCtxKey)
	if v == nil {
		c.Set(customAttributesCtxKey, []slog.Attr{attr})
		return
	}

	switch attrs := v.(type) {
	case []slog.Attr:
		c.Set(customAttributesCtxKey, append(attrs, attr))
	}
}

type Filter func(ctx echo.Context) bool

// Basic
func Accept(filter Filter) Filter { return filter }
func Ignore(filter Filter) Filter { return filter }

// Method
func AcceptMethod(methods ...string) Filter {
	return func(c echo.Context) bool {
		reqMethod := strings.ToLower(c.Request().Method)

		for _, method := range methods {
			if strings.ToLower(method) == reqMethod {
				return true
			}
		}

		return false
	}
}

func IgnoreMethod(methods ...string) Filter {
	return func(c echo.Context) bool {
		reqMethod := strings.ToLower(c.Request().Method)

		for _, method := range methods {
			if strings.ToLower(method) == reqMethod {
				return false
			}
		}

		return true
	}
}

// Status
func AcceptStatus(statuses ...int) Filter {
	return func(c echo.Context) bool {
		for _, status := range statuses {
			if status == c.Response().Status {
				return true
			}
		}

		return false
	}
}

func IgnoreStatus(statuses ...int) Filter {
	return func(c echo.Context) bool {
		for _, status := range statuses {
			if status == c.Response().Status {
				return false
			}
		}

		return true
	}
}

func AcceptStatusGreaterThan(status int) Filter {
	return func(c echo.Context) bool {
		return c.Response().Status > status
	}
}

func IgnoreStatusLessThan(status int) Filter {
	return func(c echo.Context) bool {
		return c.Response().Status < status
	}
}

func AcceptStatusGreaterThanOrEqual(status int) Filter {
	return func(c echo.Context) bool {
		return c.Response().Status >= status
	}
}

func IgnoreStatusLessThanOrEqual(status int) Filter {
	return func(c echo.Context) bool {
		return c.Response().Status <= status
	}
}

// Path
func AcceptPath(urls ...string) Filter {
	return func(c echo.Context) bool {
		for _, url := range urls {
			if c.Request().URL.Path == url {
				return true
			}
		}

		return false
	}
}

func IgnorePath(urls ...string) Filter {
	return func(c echo.Context) bool {
		for _, url := range urls {
			if c.Request().URL.Path == url {
				return false
			}
		}

		return true
	}
}

func AcceptPathContains(parts ...string) Filter {
	return func(c echo.Context) bool {
		for _, part := range parts {
			if strings.Contains(c.Request().URL.Path, part) {
				return true
			}
		}

		return false
	}
}

func IgnorePathContains(parts ...string) Filter {
	return func(c echo.Context) bool {
		for _, part := range parts {
			if strings.Contains(c.Request().URL.Path, part) {
				return false
			}
		}

		return true
	}
}

func AcceptPathPrefix(prefixs ...string) Filter {
	return func(c echo.Context) bool {
		for _, prefix := range prefixs {
			if strings.HasPrefix(c.Request().URL.Path, prefix) {
				return true
			}
		}

		return false
	}
}

func IgnorePathPrefix(prefixs ...string) Filter {
	return func(c echo.Context) bool {
		for _, prefix := range prefixs {
			if strings.HasPrefix(c.Request().URL.Path, prefix) {
				return false
			}
		}

		return true
	}
}

func AcceptPathSuffix(prefixs ...string) Filter {
	return func(c echo.Context) bool {
		for _, prefix := range prefixs {
			if strings.HasPrefix(c.Request().URL.Path, prefix) {
				return true
			}
		}

		return false
	}
}

func IgnorePathSuffix(suffixs ...string) Filter {
	return func(c echo.Context) bool {
		for _, suffix := range suffixs {
			if strings.HasSuffix(c.Request().URL.Path, suffix) {
				return false
			}
		}

		return true
	}
}

func AcceptPathMatch(regs ...regexp.Regexp) Filter {
	return func(c echo.Context) bool {
		for _, reg := range regs {
			if reg.Match([]byte(c.Request().URL.Path)) {
				return true
			}
		}

		return false
	}
}

func IgnorePathMatch(regs ...regexp.Regexp) Filter {
	return func(c echo.Context) bool {
		for _, reg := range regs {
			if reg.Match([]byte(c.Request().URL.Path)) {
				return false
			}
		}

		return true
	}
}

// Host
func AcceptHost(hosts ...string) Filter {
	return func(c echo.Context) bool {
		for _, host := range hosts {
			if c.Request().URL.Host == host {
				return true
			}
		}

		return false
	}
}

func IgnoreHost(hosts ...string) Filter {
	return func(c echo.Context) bool {
		for _, host := range hosts {
			if c.Request().URL.Host == host {
				return false
			}
		}

		return true
	}
}

func AcceptHostContains(parts ...string) Filter {
	return func(c echo.Context) bool {
		for _, part := range parts {
			if strings.Contains(c.Request().URL.Host, part) {
				return true
			}
		}

		return false
	}
}

func IgnoreHostContains(parts ...string) Filter {
	return func(c echo.Context) bool {
		for _, part := range parts {
			if strings.Contains(c.Request().URL.Host, part) {
				return false
			}
		}

		return true
	}
}

func AcceptHostPrefix(prefixs ...string) Filter {
	return func(c echo.Context) bool {
		for _, prefix := range prefixs {
			if strings.HasPrefix(c.Request().URL.Host, prefix) {
				return true
			}
		}

		return false
	}
}

func IgnoreHostPrefix(prefixs ...string) Filter {
	return func(c echo.Context) bool {
		for _, prefix := range prefixs {
			if strings.HasPrefix(c.Request().URL.Host, prefix) {
				return false
			}
		}

		return true
	}
}

func AcceptHostSuffix(prefixs ...string) Filter {
	return func(c echo.Context) bool {
		for _, prefix := range prefixs {
			if strings.HasPrefix(c.Request().URL.Host, prefix) {
				return true
			}
		}

		return false
	}
}

func IgnoreHostSuffix(suffixs ...string) Filter {
	return func(c echo.Context) bool {
		for _, suffix := range suffixs {
			if strings.HasSuffix(c.Request().URL.Host, suffix) {
				return false
			}
		}

		return true
	}
}

func AcceptHostMatch(regs ...regexp.Regexp) Filter {
	return func(c echo.Context) bool {
		for _, reg := range regs {
			if reg.Match([]byte(c.Request().URL.Host)) {
				return true
			}
		}

		return false
	}
}

func IgnoreHostMatch(regs ...regexp.Regexp) Filter {
	return func(c echo.Context) bool {
		for _, reg := range regs {
			if reg.Match([]byte(c.Request().URL.Host)) {
				return false
			}
		}

		return true
	}
}

type bodyWriter struct {
	http.ResponseWriter
	body    *bytes.Buffer
	maxSize int
	bytes   int
}

// implements gin.ResponseWriter
func (w bodyWriter) Write(b []byte) (int, error) {
	if w.body != nil {
		if w.body.Len()+len(b) > w.maxSize {
			w.body.Write(b[:w.maxSize-w.body.Len()])
		} else {
			w.body.Write(b)
		}
	}
	w.bytes += len(b)

	return w.ResponseWriter.Write(b)
}

func newBodyWriter(writer http.ResponseWriter, maxSize int, recordBody bool) *bodyWriter {
	var body *bytes.Buffer
	if recordBody {
		body = bytes.NewBufferString("")
	}

	return &bodyWriter{
		ResponseWriter: writer,
		body:           body,
		maxSize:        maxSize,
	}
}

type bodyReader struct {
	io.ReadCloser
	body    *bytes.Buffer
	maxSize int
	bytes   int
}

// implements io.Reader
func (r *bodyReader) Read(b []byte) (int, error) {
	n, err := r.ReadCloser.Read(b)
	if r.body != nil {
		if r.body.Len()+n > r.maxSize {
			r.body.Write(b[:r.maxSize-r.body.Len()])
		} else {
			r.body.Write(b)
		}
	}
	r.bytes += n
	return n, err
}

func newBodyReader(reader io.ReadCloser, maxSize int, recordBody bool) *bodyReader {
	var body *bytes.Buffer
	if recordBody {
		body = bytes.NewBufferString("")
	}

	return &bodyReader{
		ReadCloser: reader,
		body:       body,
		maxSize:    maxSize,
	}
}
