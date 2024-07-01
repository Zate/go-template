package main

import (
	"errors"
	"net/http"
)

var (
	ErrDocumentNotFound     = errors.New("DocumentNotFound")
	ErrInternalServiceError = errors.New("InternalServiceError")
	ErrMethodNotAllowed     = errors.New("MethodNotAllowed")
	ErrInvalidContentType   = errors.New("InvalidContentType")
	ErrInvalidContentLength = errors.New("InvalidContentLength")
	ErrBadGateway           = errors.New("BadGateway")
	ErrGatewayTimeout       = errors.New("GatewayTimeout")
	ErrUnauthorized         = errors.New("Unauthorized")
	ErrForbidden            = errors.New("Forbidden")
	ErrLengthRequired       = errors.New("LengthRequired")
	ErrPayloadTooLarge      = errors.New("PayloadTooLarge")
	ErrURITooLong           = errors.New("URITooLong")
	ErrUnsupportedMediaType = errors.New("UnsupportedMediaType")
	ErrImATeaPot            = errors.New("ImATeaPot")
	ErrTooManyRequests      = errors.New("TooManyRequests")
)

func NewErrorStatusCodeMaps() map[error]int {

	var errorStatusCodeMaps = make(map[error]int)
	errorStatusCodeMaps[ErrDocumentNotFound] = http.StatusNotFound
	errorStatusCodeMaps[ErrInternalServiceError] = http.StatusInternalServerError
	errorStatusCodeMaps[ErrMethodNotAllowed] = http.StatusMethodNotAllowed
	errorStatusCodeMaps[ErrInvalidContentType] = http.StatusBadRequest
	errorStatusCodeMaps[ErrInvalidContentLength] = http.StatusBadRequest
	errorStatusCodeMaps[ErrBadGateway] = http.StatusBadGateway
	errorStatusCodeMaps[ErrGatewayTimeout] = http.StatusGatewayTimeout
	errorStatusCodeMaps[ErrUnauthorized] = http.StatusUnauthorized
	errorStatusCodeMaps[ErrForbidden] = http.StatusForbidden
	errorStatusCodeMaps[ErrLengthRequired] = http.StatusLengthRequired
	errorStatusCodeMaps[ErrPayloadTooLarge] = http.StatusRequestEntityTooLarge
	errorStatusCodeMaps[ErrURITooLong] = http.StatusRequestURITooLong
	errorStatusCodeMaps[ErrUnsupportedMediaType] = http.StatusUnsupportedMediaType
	errorStatusCodeMaps[ErrImATeaPot] = http.StatusTeapot
	errorStatusCodeMaps[ErrTooManyRequests] = http.StatusTooManyRequests
	return errorStatusCodeMaps
}
