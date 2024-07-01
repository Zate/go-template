package main

// // package main

// // import (
// //     "net/http"

// //     "github.com/labstack/echo/v4"
// // )

// // type RequestPayload struct {
// //     Name string `json:"name"`
// // }

// // func GreetHandler(c echo.Context) error {
// //     var payload RequestPayload
// //     if err := c.Bind(&payload); err != nil {
// //         return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid payload"})
// //     }
// //     return c.JSON(http.StatusOK, map[string]string{"message": "Hello, " + payload.Name + "!"})
// // }

// // package main

// // import (
// //     "bytes"
// //     "encoding/json"
// //     "io/ioutil"
// //     "net/http"

// //     "github.com/labstack/echo/v4"
// // )

// // type RequestPayload struct {
// //     Name string `json:"name"`
// // }

// // type ResponsePayload struct {
// //     Message string `json:"message"`
// // }

// // func ExternalService(payload RequestPayload) (ResponsePayload, error) {
// //     url := "http://external-service/endpoint"
// //     jsonData, err := json.Marshal(payload)
// //     if err != nil {
// //         return ResponsePayload{}, err
// //     }

// //     resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
// //     if err != nil {
// //         return ResponsePayload{}, err
// //     }
// //     defer resp.Body.Close()

// //     body, err := ioutil.ReadAll(resp.Body)
// //     if err != nil {
// //         return ResponsePayload{}, err
// //     }

// //     var responsePayload ResponsePayload
// //     err = json.Unmarshal(body, &responsePayload)
// //     if err != nil {
// //         return ResponsePayload{}, err
// //     }

// //     return responsePayload, nil
// // }

// // func GreetHandler(c echo.Context) error {
// //     var payload RequestPayload
// //     if err := c.Bind(&payload); err != nil {
// //         return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid payload"})
// //     }

// //     responsePayload, err := ExternalService(payload)
// //     if err != nil {
// //         return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to call external service"})
// //     }

// //     return c.JSON(http.StatusOK, responsePayload)
// // }

// // package main

// // import (
// //     "encoding/json"
// //     "net/http"
// //     "net/http/httptest"
// //     "strings"
// //     "testing"

// //     "github.com/labstack/echo/v4"
// //     "github.com/stretchr/testify/assert"
// // )

// // func TestGreetHandler(t *testing.T) {
// //     // Create a mock server for the external service
// //     mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// //         var payload RequestPayload
// //         err := json.NewDecoder(r.Body).Decode(&payload)
// //         assert.NoError(t, err)

// //         responsePayload := ResponsePayload{Message: "Hello, " + payload.Name + "!"}
// //         w.Header().Set("Content-Type", "application/json")
// //         json.NewEncoder(w).Encode(responsePayload)
// //     }))
// //     defer mockServer.Close()

// //     // Override the ExternalService function to use the mock server
// //     originalExternalService := ExternalService
// //     ExternalService = func(payload RequestPayload) (ResponsePayload, error) {
// //         url := mockServer.URL
// //         jsonData, err := json.Marshal(payload)
// //         if err != nil {
// //             return ResponsePayload{}, err
// //         }

// //         resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
// //         if err != nil {
// //             return ResponsePayload{}, err
// //         }
// //         defer resp.Body.Close()

// //         body, err := ioutil.ReadAll(resp.Body)
// //         if err != nil {
// //             return ResponsePayload{}, err
// //         }

// //         var responsePayload ResponsePayload
// //         err = json.Unmarshal(body, &responsePayload)
// //         if err != nil {
// //             return ResponsePayload{}, err
// //         }

// //         return responsePayload, nil
// //     }
// //     defer func() { ExternalService = originalExternalService }() // Restore the original function after the test

// //     // Create an instance of Echo
// //     e := echo.New()

// //     // Define the JSON payload
// //     jsonPayload := `{"name":"John"}`

// //     // Create a new POST request with the JSON payload
// //     req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(jsonPayload))
// //     req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

// //     // Create a new HTTP response recorder
// //     rec := httptest.NewRecorder()

// //     // Create a new context with the request and response
// //     c := e.NewContext(req, rec)

// //     // Invoke the handler function
// //     if assert.NoError(t, GreetHandler(c)) {
// //         // Assert the status code
// //         assert.Equal(t, http.StatusOK, rec.Code)

// //         // Assert the response body
// //         expectedResponse := `{"message":"Hello, John!"}`
// //         assert.JSONEq(t, expectedResponse, rec.Body.String())
// //     }
// // }
