package threads

import (
	"errors"
	"testing"
	"time"
)

// TestBaseError_Error tests the BaseError.Error() method
func TestBaseError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *BaseError
		contains []string
	}{
		{
			name: "with details",
			err: &BaseError{
				Code:    400,
				Message: "Test message",
				Type:    "test_error",
				Details: "Test details",
			},
			contains: []string{"400", "Test message", "test_error", "Test details"},
		},
		{
			name: "without details",
			err: &BaseError{
				Code:    500,
				Message: "Server error",
				Type:    "server_error",
			},
			contains: []string{"500", "Server error", "server_error"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errStr := tt.err.Error()
			for _, substr := range tt.contains {
				if !containsSubstr(errStr, substr) {
					t.Errorf("error string should contain '%s', got '%s'", substr, errStr)
				}
			}
		})
	}
}

// TestNewAuthenticationError tests creating authentication errors
func TestNewAuthenticationError(t *testing.T) {
	err := NewAuthenticationError(401, "Unauthorized", "Token expired")

	if err.Code != 401 {
		t.Errorf("expected Code 401, got %d", err.Code)
	}
	if err.Message != "Unauthorized" {
		t.Errorf("expected Message 'Unauthorized', got '%s'", err.Message)
	}
	if err.Details != "Token expired" {
		t.Errorf("expected Details 'Token expired', got '%s'", err.Details)
	}
	if err.Type != "authentication_error" {
		t.Errorf("expected Type 'authentication_error', got '%s'", err.Type)
	}
}

// TestNewRateLimitError tests creating rate limit errors
func TestNewRateLimitError(t *testing.T) {
	retryAfter := 60 * time.Second
	err := NewRateLimitError(429, "Rate limited", "Too many requests", retryAfter)

	if err.Code != 429 {
		t.Errorf("expected Code 429, got %d", err.Code)
	}
	if err.Message != "Rate limited" {
		t.Errorf("expected Message 'Rate limited', got '%s'", err.Message)
	}
	if err.Details != "Too many requests" {
		t.Errorf("expected Details 'Too many requests', got '%s'", err.Details)
	}
	if err.Type != "rate_limit_error" {
		t.Errorf("expected Type 'rate_limit_error', got '%s'", err.Type)
	}
	if err.RetryAfter != retryAfter {
		t.Errorf("expected RetryAfter %v, got %v", retryAfter, err.RetryAfter)
	}
}

// TestNewValidationError tests creating validation errors
func TestNewValidationError(t *testing.T) {
	err := NewValidationError(400, "Invalid input", "Field must be non-empty", "username")

	if err.Code != 400 {
		t.Errorf("expected Code 400, got %d", err.Code)
	}
	if err.Message != "Invalid input" {
		t.Errorf("expected Message 'Invalid input', got '%s'", err.Message)
	}
	if err.Details != "Field must be non-empty" {
		t.Errorf("expected Details 'Field must be non-empty', got '%s'", err.Details)
	}
	if err.Type != "validation_error" {
		t.Errorf("expected Type 'validation_error', got '%s'", err.Type)
	}
	if err.Field != "username" {
		t.Errorf("expected Field 'username', got '%s'", err.Field)
	}
}

// TestNewNetworkError tests creating network errors
func TestNewNetworkError(t *testing.T) {
	tests := []struct {
		name      string
		code      int
		message   string
		details   string
		temporary bool
	}{
		{"temporary error", 0, "Connection failed", "DNS lookup failed", true},
		{"permanent error", 0, "Invalid host", "Host not found", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewNetworkError(tt.code, tt.message, tt.details, tt.temporary)

			if err.Code != tt.code {
				t.Errorf("expected Code %d, got %d", tt.code, err.Code)
			}
			if err.Message != tt.message {
				t.Errorf("expected Message '%s', got '%s'", tt.message, err.Message)
			}
			if err.Details != tt.details {
				t.Errorf("expected Details '%s', got '%s'", tt.details, err.Details)
			}
			if err.Type != "network_error" {
				t.Errorf("expected Type 'network_error', got '%s'", err.Type)
			}
			if err.Temporary != tt.temporary {
				t.Errorf("expected Temporary %v, got %v", tt.temporary, err.Temporary)
			}
		})
	}
}

// TestNewAPIError tests creating API errors
func TestNewAPIError(t *testing.T) {
	err := NewAPIError(500, "Internal error", "Server crashed", "req-123")

	if err.Code != 500 {
		t.Errorf("expected Code 500, got %d", err.Code)
	}
	if err.Message != "Internal error" {
		t.Errorf("expected Message 'Internal error', got '%s'", err.Message)
	}
	if err.Details != "Server crashed" {
		t.Errorf("expected Details 'Server crashed', got '%s'", err.Details)
	}
	if err.Type != "api_error" {
		t.Errorf("expected Type 'api_error', got '%s'", err.Type)
	}
	if err.RequestID != "req-123" {
		t.Errorf("expected RequestID 'req-123', got '%s'", err.RequestID)
	}
}

// TestIsAuthenticationError tests the IsAuthenticationError helper
func TestIsAuthenticationError(t *testing.T) {
	authErr := NewAuthenticationError(401, "Unauthorized", "")
	validationErr := NewValidationError(400, "Invalid", "", "")
	regularErr := errors.New("regular error")

	if !IsAuthenticationError(authErr) {
		t.Error("IsAuthenticationError should return true for AuthenticationError")
	}
	if IsAuthenticationError(validationErr) {
		t.Error("IsAuthenticationError should return false for ValidationError")
	}
	if IsAuthenticationError(regularErr) {
		t.Error("IsAuthenticationError should return false for regular error")
	}
	if IsAuthenticationError(nil) {
		t.Error("IsAuthenticationError should return false for nil")
	}
}

// TestIsRateLimitError tests the IsRateLimitError helper
func TestIsRateLimitError(t *testing.T) {
	rateLimitErr := NewRateLimitError(429, "Rate limited", "", 60*time.Second)
	authErr := NewAuthenticationError(401, "Unauthorized", "")
	regularErr := errors.New("regular error")

	if !IsRateLimitError(rateLimitErr) {
		t.Error("IsRateLimitError should return true for RateLimitError")
	}
	if IsRateLimitError(authErr) {
		t.Error("IsRateLimitError should return false for AuthenticationError")
	}
	if IsRateLimitError(regularErr) {
		t.Error("IsRateLimitError should return false for regular error")
	}
	if IsRateLimitError(nil) {
		t.Error("IsRateLimitError should return false for nil")
	}
}

// TestIsValidationError tests the IsValidationError helper
func TestIsValidationError(t *testing.T) {
	validationErr := NewValidationError(400, "Invalid", "", "field")
	authErr := NewAuthenticationError(401, "Unauthorized", "")
	regularErr := errors.New("regular error")

	if !IsValidationError(validationErr) {
		t.Error("IsValidationError should return true for ValidationError")
	}
	if IsValidationError(authErr) {
		t.Error("IsValidationError should return false for AuthenticationError")
	}
	if IsValidationError(regularErr) {
		t.Error("IsValidationError should return false for regular error")
	}
	if IsValidationError(nil) {
		t.Error("IsValidationError should return false for nil")
	}
}

// TestIsNetworkError tests the IsNetworkError helper
func TestIsNetworkError(t *testing.T) {
	networkErr := NewNetworkError(0, "Network failed", "", true)
	authErr := NewAuthenticationError(401, "Unauthorized", "")
	regularErr := errors.New("regular error")

	if !IsNetworkError(networkErr) {
		t.Error("IsNetworkError should return true for NetworkError")
	}
	if IsNetworkError(authErr) {
		t.Error("IsNetworkError should return false for AuthenticationError")
	}
	if IsNetworkError(regularErr) {
		t.Error("IsNetworkError should return false for regular error")
	}
	if IsNetworkError(nil) {
		t.Error("IsNetworkError should return false for nil")
	}
}

// TestIsAPIError tests the IsAPIError helper
func TestIsAPIError(t *testing.T) {
	apiErr := NewAPIError(500, "Server error", "", "")
	authErr := NewAuthenticationError(401, "Unauthorized", "")
	regularErr := errors.New("regular error")

	if !IsAPIError(apiErr) {
		t.Error("IsAPIError should return true for APIError")
	}
	if IsAPIError(authErr) {
		t.Error("IsAPIError should return false for AuthenticationError")
	}
	if IsAPIError(regularErr) {
		t.Error("IsAPIError should return false for regular error")
	}
	if IsAPIError(nil) {
		t.Error("IsAPIError should return false for nil")
	}
}

// TestErrorsAs tests that errors work with errors.As
func TestErrorsAs(t *testing.T) {
	authErr := NewAuthenticationError(401, "Unauthorized", "Token expired")

	var target *AuthenticationError
	if !errors.As(authErr, &target) {
		t.Error("errors.As should work with AuthenticationError")
	}
	if target.Message != "Unauthorized" {
		t.Errorf("expected Message 'Unauthorized', got '%s'", target.Message)
	}
}

// Helper function
func containsSubstr(s, substr string) bool {
	for i := 0; i+len(substr) <= len(s); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
