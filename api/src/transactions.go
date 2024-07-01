package main

import (
	"time"

	"github.com/labstack/echo/v4"
)

// TODO: Try and refactor this to use a single Transaction struct for all handlers, or a TransactionInterface for each handler

// TransactionData represents the request and response information gathered and tracked across the lifetime of a transaction
type Transaction struct {
	// optional Error message
	Error *echo.HTTPError
	// RTT is the round trip time for the transaction
	RTT time.Duration
	// StartTime is the time the transaction started
	StartTime time.Time
	// Data is an interface for the data that is collected or collated as part of the transaction
	Data interface{}
	// Payload is the payload of the transaction to be returned as a JSON object
	Payload interface{}
	// TokensSent is the number of tokens sent to the AI API
	TokensSent int
	// TokensReceived is the number of tokens received from the AI API
	TokensReceived int
}

// SetRTT sets the RTT for the transaction
func (t *Transaction) SetRTT() {
	if t.StartTime.IsZero() {
		t.RTT = 0
		return
	}
	t.RTT = time.Duration(time.Since(t.StartTime).Milliseconds())
}

// SetStartTime sets the start time for the transaction
func (t *Transaction) SetStartTime(startTime time.Time) {
	t.StartTime = startTime
}

// SetErr sets the error for the transaction
func (t *Transaction) SetErr(err *echo.HTTPError) {
	t.Error = err
}

// SetResult sets the result for the transaction
func (t *Transaction) SetData(data interface{}) {
	t.Data = data
}

// SetPayload sets the payload for the transaction
func (t *Transaction) SetPayload(payload interface{}) {
	t.Payload = payload
}

// SetTokensSent sets the number of tokens sent to the AI API
func (t *Transaction) SetTokensSent(tokensSent int) {
	t.TokensSent = tokensSent
}

// SetTokensReceived sets the number of tokens received from the AI API
func (t *Transaction) SetTokensReceived(tokensReceived int) {
	t.TokensReceived = tokensReceived
}

// NewTransaction creates a new Transaction object
func NewTransaction() *Transaction {
	return &Transaction{}
}

// Log logs the transaction data to the console
func (t *Transaction) Log() {
	// log the transaction data
}

// TransactionInterface is an interface for the Transaction struct
type TransactionInterface interface {
	SetRTT()
	SetStartTime()
	SetErr(err *echo.HTTPError)
	SetResult(result interface{})
	SetPayload(payload interface{})
	SetTokensSent(tokensSent int)
	SetTokensReceived(tokensReceived int)
	Log()
}
