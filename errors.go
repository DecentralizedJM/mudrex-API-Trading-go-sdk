package mudrex

import (
	"encoding/json"
	"fmt"
)

// MudrexError is the base error type
type MudrexError struct {
	Code    int
	Message string
	Status  int
}

func (e *MudrexError) Error() string {
	return fmt.Sprintf("Mudrex API Error (code=%d, status=%d): %s", e.Code, e.Status, e.Message)
}

// Specific error types
type AuthenticationError struct {
	*MudrexError
}

type RateLimitError struct {
	*MudrexError
}

type ValidationError struct {
	*MudrexError
}

type NotFoundError struct {
	*MudrexError
}

type ConflictError struct {
	*MudrexError
}

type ServerError struct {
	*MudrexError
}

type InsufficientBalanceError struct {
	*MudrexError
}

// RaiseForError checks HTTP status and response body and raises appropriate error
func RaiseForError(status int, body []byte) error {
	if status < 400 {
		return nil
	}
	
	var apiResp APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		// If we can't parse the response, create a generic error
		return &MudrexError{
			Code:    -1,
			Message: string(body),
			Status:  status,
		}
	}
	
	baseErr := &MudrexError{
		Code:    status,
		Message: apiResp.Message,
		Status:  status,
	}
	
	if apiResp.Error != nil {
		baseErr.Code = apiResp.Error.Code
		baseErr.Message = apiResp.Error.Message
	}
	
	// Return specific error based on status code
	switch status {
	case 401:
		return &AuthenticationError{baseErr}
	case 429:
		return &RateLimitError{baseErr}
	case 400:
		return &ValidationError{baseErr}
	case 404:
		return &NotFoundError{baseErr}
	case 409:
		return &ConflictError{baseErr}
	case 500, 502, 503, 504:
		return &ServerError{baseErr}
	default:
		if status >= 400 && status < 500 {
			// Check for insufficient balance (400 with specific message)
			if baseErr.Code == 1002 || contains(baseErr.Message, "insufficient balance") {
				return &InsufficientBalanceError{baseErr}
			}
			return &ValidationError{baseErr}
		}
		return &ServerError{baseErr}
	}
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
