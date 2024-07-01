package main

import (
	"fmt"
	"html"
	"log/slog"

	"net/http"
	"time"

	"github.com/enescakir/emoji"
	"github.com/labstack/echo/v4"
)

// APIResponse represents the response structure for API calls.
type APIResponse struct {
	api         string
	respBody    string
	err         error
	statusCode  int
	rawRespBody []byte
}

// Status represents the response structure for the /status endpoint.
type Status struct {
	Name    string
	Message string
	Code    int
	Emoji   string
	Error   *echo.HTTPError
	RTT     time.Duration
}

// StatusJSONResponse represents the JSON response structure for the /status endpoint.
type StatusJSONResponse struct {
	Service1 string `json:"service1"`
	Service2 string `json:"service2"`
	Service3 string `json:"service3"`
}

// StatusHandler is a function that handles requests to the /status endpoint.
func (s *Service) StatusHandler(c echo.Context) error {
	var payload StatusJSONResponse
	var service1LogAttr, service2LogAttr, service3LogAttr slog.Attr
	var numServices int

	if c.QueryParam("debug") != "" {
		Debug = true
	}

	numServices = 3

	// Create a channel to receive statuses
	statusChan := make(chan Status, numServices)

	// Launch each check in a separate goroutine
	go func() { statusChan <- s.checkService1() }()
	go func() { statusChan <- s.checkService2() }()
	go func() { statusChan <- s.checkService3() }()

	// Wait for all checks to complete
	for i := 0; i < numServices; i++ {
		status := <-statusChan

		switch status.Name {
		case "service1":
			service1LogAttr = s.logStatus(status)
			payload.Service1 = status.Emoji
		case "service2":
			service2LogAttr = s.logStatus(status)
			payload.Service2 = status.Emoji
		case "service3":
			service3LogAttr = s.logStatus(status)
			payload.Service3 = status.Emoji
		}
	}

	s.AddCustomAttributes(slog.Group(s.Path, service1LogAttr, service2LogAttr, service3LogAttr))
	return c.JSON(http.StatusOK, payload)
}

// checkStatus1 is a
func (s *Service) checkService1() Status {
	var status Status

	startTime := (*s.C).Get("start_time").(time.Time)

	status.Name = "service1"

	response, err := s.mockService(status.Name)
	if err != nil {
		s.logCustomAttr("MOCK_SERVICE1_REQUEST_ERROR", map[string]interface{}{"error": err})
		return s.statusErr(status.Name, startTime, echo.NewHTTPError(http.StatusInternalServerError, err.Error()))
	}

	status.Code = 200
	status.Emoji = response.Emoji

	status.Message = response.Message
	status.RTT = time.Duration(time.Since(startTime).Milliseconds())
	return status
}

// checkStatus2 is a function
func (s *Service) checkService2() Status {
	var status Status

	startTime := (*s.C).Get("start_time").(time.Time)

	status.Name = "service2"

	response, err := s.mockService(status.Name)
	if err != nil {
		s.logCustomAttr("MOCK_SERVICE2_REQUEST_ERROR", map[string]interface{}{"error": err})
		return s.statusErr(status.Name, startTime, echo.NewHTTPError(http.StatusInternalServerError, err.Error()))
	}

	status.Code = 200
	status.Emoji = response.Emoji

	status.Message = response.Message
	status.RTT = time.Duration(time.Since(startTime).Milliseconds())
	return status
}

// checkStatus3 is a function
func (s *Service) checkService3() Status {
	var status Status

	startTime := (*s.C).Get("start_time").(time.Time)

	status.Name = "service3"

	response, err := s.mockService(status.Name)
	if err != nil {
		s.logCustomAttr("MOCK_SERVICE3_REQUEST_ERROR", map[string]interface{}{"error": err})
		return s.statusErr(status.Name, startTime, echo.NewHTTPError(http.StatusInternalServerError, err.Error()))
	}

	status.Code = 200
	status.Emoji = response.Emoji

	status.Message = response.Message
	status.RTT = time.Duration(time.Since(startTime).Milliseconds())
	return status
}

// function that takes an error and returs a Status struct
func (s *Service) statusErr(name string, startTime time.Time, err *echo.HTTPError) Status {
	return Status{
		Name:    name,
		Message: err.Error(),
		Code:    err.Code,
		Emoji:   html.UnescapeString(emoji.Sprint(":red_circle:")),
		Error:   err,
		RTT:     time.Duration(time.Since(startTime).Milliseconds()),
	}
}

// TODO: Make this more generic to be able to log for any handler
// logStatus is a function to turn the status object into a log.Attr object
func (s *Service) logStatus(status Status) slog.Attr {
	logData := map[string]interface{}{
		"code":    status.Code,
		"message": status.Message,
		"emoji":   status.Emoji,
	}
	if status.Error != nil {
		logData["error"] = status.Error.Message
		logData["error_code"] = status.Error.Code
	}
	return s.genCustomAttr(status.Name, logData)
}

// mockService is a function to mock a service check
func (s *Service) mockService(name string) (Status, error) {
	var status Status
	startTime := (*s.C).Get("start_time").(time.Time)

	status.Name = name
	status.Code = http.StatusOK
	status.Emoji = html.UnescapeString(emoji.Sprint(":green_circle:"))
	status.Message = fmt.Sprintf("Mock service %s is up", name)
	status.RTT = time.Duration(time.Since(startTime).Milliseconds())
	return status, nil
}
