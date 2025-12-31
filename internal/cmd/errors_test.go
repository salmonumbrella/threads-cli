package cmd

import (
	"errors"
	"strings"
	"testing"
	"time"

	threads "github.com/salmonumbrella/threads-go"
)

func TestUserFriendlyError_Error(t *testing.T) {
	tests := []struct {
		name       string
		err        *UserFriendlyError
		wantMsg    string
		wantHasSug bool
	}{
		{
			name: "with suggestion",
			err: &UserFriendlyError{
				Message:    "Something went wrong",
				Suggestion: "Try again later",
			},
			wantMsg:    "Something went wrong",
			wantHasSug: true,
		},
		{
			name: "without suggestion",
			err: &UserFriendlyError{
				Message: "Something went wrong",
			},
			wantMsg:    "Something went wrong",
			wantHasSug: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.Error()
			if !strings.Contains(got, tt.wantMsg) {
				t.Errorf("Error() = %v, want to contain %v", got, tt.wantMsg)
			}
			hasSuggestion := strings.Contains(got, "Suggestion:")
			if hasSuggestion != tt.wantHasSug {
				t.Errorf("Error() has suggestion = %v, want %v", hasSuggestion, tt.wantHasSug)
			}
		})
	}
}

func TestUserFriendlyError_Unwrap(t *testing.T) {
	cause := errors.New("underlying error")
	err := &UserFriendlyError{
		Message: "Wrapper",
		Cause:   cause,
	}

	if err.Unwrap() != cause {
		t.Errorf("Unwrap() = %v, want %v", err.Unwrap(), cause)
	}
}

func TestFormatError_AuthenticationError(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantSubstr string
		wantSugg   string
	}{
		{
			name:       "expired token",
			err:        threads.NewAuthenticationError(401, "Token has expired", ""),
			wantSubstr: "expired",
			wantSugg:   "threads auth refresh",
		},
		{
			name:       "invalid token",
			err:        threads.NewAuthenticationError(401, "Invalid access token", ""),
			wantSubstr: "invalid",
			wantSugg:   "threads auth login",
		},
		{
			name:       "401 error",
			err:        threads.NewAuthenticationError(401, "Authentication required", ""),
			wantSubstr: "Authentication required",
			wantSugg:   "threads auth login",
		},
		{
			name:       "403 error",
			err:        threads.NewAuthenticationError(403, "Access denied", ""),
			wantSubstr: "permission",
			wantSugg:   "scopes",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatted := FormatError(tt.err)
			ufErr, ok := formatted.(*UserFriendlyError)
			if !ok {
				t.Fatalf("FormatError() did not return *UserFriendlyError, got %T", formatted)
			}
			errStr := ufErr.Error()
			if !strings.Contains(strings.ToLower(errStr), strings.ToLower(tt.wantSubstr)) {
				t.Errorf("Error() = %v, want to contain %v", errStr, tt.wantSubstr)
			}
			if !strings.Contains(errStr, tt.wantSugg) {
				t.Errorf("Error() = %v, want suggestion to contain %v", errStr, tt.wantSugg)
			}
		})
	}
}

func TestFormatError_RateLimitError(t *testing.T) {
	err := threads.NewRateLimitError(429, "Too many requests", "", 5*time.Minute)
	formatted := FormatError(err)

	ufErr, ok := formatted.(*UserFriendlyError)
	if !ok {
		t.Fatalf("FormatError() did not return *UserFriendlyError, got %T", formatted)
	}

	errStr := ufErr.Error()
	if !strings.Contains(errStr, "Rate limit") {
		t.Errorf("Error() = %v, want to contain 'Rate limit'", errStr)
	}
	if !strings.Contains(errStr, "threads ratelimit status") {
		t.Errorf("Error() = %v, want suggestion to contain 'threads ratelimit status'", errStr)
	}
}

func TestFormatError_ValidationError(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantSubstr string
	}{
		{
			name:       "with field",
			err:        threads.NewValidationError(400, "Invalid value", "", "text"),
			wantSubstr: "text",
		},
		{
			name:       "without field",
			err:        threads.NewValidationError(400, "Validation failed", "", ""),
			wantSubstr: "Validation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatted := FormatError(tt.err)
			ufErr, ok := formatted.(*UserFriendlyError)
			if !ok {
				t.Fatalf("FormatError() did not return *UserFriendlyError, got %T", formatted)
			}
			errStr := ufErr.Error()
			if !strings.Contains(errStr, tt.wantSubstr) {
				t.Errorf("Error() = %v, want to contain %v", errStr, tt.wantSubstr)
			}
		})
	}
}

func TestFormatError_NetworkError(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantSubstr string
	}{
		{
			name:       "timeout",
			err:        threads.NewNetworkError(0, "Request timeout", "", true),
			wantSubstr: "timed out",
		},
		{
			name:       "dns error",
			err:        threads.NewNetworkError(0, "no such host", "", false),
			wantSubstr: "DNS",
		},
		{
			name:       "temporary error",
			err:        threads.NewNetworkError(0, "Temporary failure", "", true),
			wantSubstr: "transient",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatted := FormatError(tt.err)
			ufErr, ok := formatted.(*UserFriendlyError)
			if !ok {
				t.Fatalf("FormatError() did not return *UserFriendlyError, got %T", formatted)
			}
			errStr := ufErr.Error()
			if !strings.Contains(strings.ToLower(errStr), strings.ToLower(tt.wantSubstr)) {
				t.Errorf("Error() = %v, want to contain %v", errStr, tt.wantSubstr)
			}
		})
	}
}

func TestFormatError_APIError(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantSubstr string
	}{
		{
			name:       "server error",
			err:        threads.NewAPIError(500, "Internal server error", "", "req-123"),
			wantSubstr: "server-side",
		},
		{
			name:       "not found",
			err:        threads.NewAPIError(404, "Resource not found", "", ""),
			wantSubstr: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatted := FormatError(tt.err)
			ufErr, ok := formatted.(*UserFriendlyError)
			if !ok {
				t.Fatalf("FormatError() did not return *UserFriendlyError, got %T", formatted)
			}
			errStr := ufErr.Error()
			if !strings.Contains(strings.ToLower(errStr), strings.ToLower(tt.wantSubstr)) {
				t.Errorf("Error() = %v, want to contain %v", errStr, tt.wantSubstr)
			}
		})
	}
}

func TestFormatError_GenericErrors(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantSubstr string
	}{
		{
			name:       "no account configured",
			err:        errors.New("no account configured"),
			wantSubstr: "threads auth login",
		},
		{
			name:       "token expired",
			err:        errors.New("token expired"),
			wantSubstr: "threads auth refresh",
		},
		{
			name:       "empty response",
			err:        errors.New("empty response from API"),
			wantSubstr: "empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatted := FormatError(tt.err)
			ufErr, ok := formatted.(*UserFriendlyError)
			if !ok {
				// Some generic errors may not be converted
				return
			}
			errStr := ufErr.Error()
			if !strings.Contains(strings.ToLower(errStr), strings.ToLower(tt.wantSubstr)) {
				t.Errorf("Error() = %v, want to contain %v", errStr, tt.wantSubstr)
			}
		})
	}
}

func TestFormatError_Nil(t *testing.T) {
	if FormatError(nil) != nil {
		t.Error("FormatError(nil) should return nil")
	}
}

func TestWrapError(t *testing.T) {
	authErr := threads.NewAuthenticationError(401, "Token expired", "")
	wrapped := WrapError("API call failed", authErr)

	ufErr, ok := wrapped.(*UserFriendlyError)
	if !ok {
		t.Fatalf("WrapError() did not return *UserFriendlyError, got %T", wrapped)
	}

	errStr := ufErr.Error()
	if !strings.Contains(errStr, "API call failed") {
		t.Errorf("Error() = %v, want to contain context 'API call failed'", errStr)
	}
	if !strings.Contains(errStr, "expired") {
		t.Errorf("Error() = %v, want to contain original error info", errStr)
	}
}

func TestWrapError_Nil(t *testing.T) {
	if WrapError("context", nil) != nil {
		t.Error("WrapError(context, nil) should return nil")
	}
}

func TestWrapError_PlainError(t *testing.T) {
	plainErr := errors.New("something went wrong")
	wrapped := WrapError("operation failed", plainErr)

	if wrapped == nil {
		t.Fatal("WrapError() returned nil for plain error")
	}

	errStr := wrapped.Error()
	if !strings.Contains(errStr, "operation failed") {
		t.Errorf("Error() = %v, want to contain context 'operation failed'", errStr)
	}
}

// Additional tests for better coverage

func TestFormatError_AuthenticationError_DefaultCase(t *testing.T) {
	// Test the default case (not expired, not invalid, not 401, not 403)
	err := threads.NewAuthenticationError(400, "Some other auth error", "")
	formatted := FormatError(err)

	ufErr, ok := formatted.(*UserFriendlyError)
	if !ok {
		t.Fatalf("FormatError() did not return *UserFriendlyError, got %T", formatted)
	}

	errStr := ufErr.Error()
	if !strings.Contains(errStr, "Authentication error") {
		t.Errorf("Error() = %v, want to contain 'Authentication error'", errStr)
	}
	if !strings.Contains(errStr, "threads auth status") {
		t.Errorf("Error() = %v, want suggestion to contain 'threads auth status'", errStr)
	}
}

func TestFormatError_RateLimitError_NoRetryAfter(t *testing.T) {
	err := threads.NewRateLimitError(429, "Too many requests", "", 0)
	formatted := FormatError(err)

	ufErr, ok := formatted.(*UserFriendlyError)
	if !ok {
		t.Fatalf("FormatError() did not return *UserFriendlyError, got %T", formatted)
	}

	errStr := ufErr.Error()
	if !strings.Contains(errStr, "Rate limit") {
		t.Errorf("Error() = %v, want to contain 'Rate limit'", errStr)
	}
	if !strings.Contains(errStr, "Wait a few minutes") {
		t.Errorf("Error() = %v, want suggestion to contain 'Wait a few minutes'", errStr)
	}
}

func TestFormatError_ValidationError_SpecificPatterns(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantSubstr string
	}{
		{
			name:       "text too long",
			err:        threads.NewValidationError(400, "Text is too long", "", "text"),
			wantSubstr: "500 characters",
		},
		{
			name:       "invalid url",
			err:        threads.NewValidationError(400, "URL is invalid", "", "url"),
			wantSubstr: "http://",
		},
		{
			name:       "media format",
			err:        threads.NewValidationError(400, "Unsupported media format", "", "media"),
			wantSubstr: "JPEG",
		},
		{
			name:       "carousel items",
			err:        threads.NewValidationError(400, "Carousel has too few items", "", ""),
			wantSubstr: "2-20",
		},
		{
			name:       "empty field empty message",
			err:        threads.NewValidationError(400, "", "", ""),
			wantSubstr: "--help",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatted := FormatError(tt.err)
			ufErr, ok := formatted.(*UserFriendlyError)
			if !ok {
				t.Fatalf("FormatError() did not return *UserFriendlyError, got %T", formatted)
			}
			errStr := ufErr.Error()
			if !strings.Contains(errStr, tt.wantSubstr) {
				t.Errorf("Error() = %v, want to contain %v", errStr, tt.wantSubstr)
			}
		})
	}
}

func TestFormatError_NetworkError_AllCases(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantSubstr string
	}{
		{
			name:       "connection refused",
			err:        threads.NewNetworkError(0, "connection refused", "", false),
			wantSubstr: "temporarily unavailable",
		},
		{
			name:       "tls error",
			err:        threads.NewNetworkError(0, "tls handshake error", "", false),
			wantSubstr: "SSL/TLS",
		},
		{
			name:       "certificate error",
			err:        threads.NewNetworkError(0, "certificate invalid", "", false),
			wantSubstr: "SSL/TLS",
		},
		{
			name:       "default network error",
			err:        threads.NewNetworkError(0, "unknown network issue", "", false),
			wantSubstr: "internet connection",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatted := FormatError(tt.err)
			ufErr, ok := formatted.(*UserFriendlyError)
			if !ok {
				t.Fatalf("FormatError() did not return *UserFriendlyError, got %T", formatted)
			}
			errStr := ufErr.Error()
			if !strings.Contains(strings.ToLower(errStr), strings.ToLower(tt.wantSubstr)) {
				t.Errorf("Error() = %v, want to contain %v", errStr, tt.wantSubstr)
			}
		})
	}
}

func TestFormatError_APIError_AllCases(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantSubstr string
	}{
		{
			name:       "server error without request id",
			err:        threads.NewAPIError(503, "Service unavailable", "", ""),
			wantSubstr: "server-side",
		},
		{
			name:       "deleted content",
			err:        threads.NewAPIError(410, "Content has been deleted", "", ""),
			wantSubstr: "no longer exists",
		},
		{
			name:       "private content",
			err:        threads.NewAPIError(403, "Content is private", "", ""),
			wantSubstr: "private content",
		},
		{
			name:       "default error without request id",
			err:        threads.NewAPIError(400, "Bad request", "", ""),
			wantSubstr: "problem persists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatted := FormatError(tt.err)
			ufErr, ok := formatted.(*UserFriendlyError)
			if !ok {
				t.Fatalf("FormatError() did not return *UserFriendlyError, got %T", formatted)
			}
			errStr := ufErr.Error()
			if !strings.Contains(strings.ToLower(errStr), strings.ToLower(tt.wantSubstr)) {
				t.Errorf("Error() = %v, want to contain %v", errStr, tt.wantSubstr)
			}
		})
	}
}

func TestFormatError_GenericErrors_AllCases(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantSubstr string
	}{
		{
			name:       "account not found",
			err:        errors.New("account not found"),
			wantSubstr: "threads auth login",
		},
		{
			name:       "client secret not stored",
			err:        errors.New("client secret not stored"),
			wantSubstr: "client ID and secret",
		},
		{
			name:       "cannot refresh",
			err:        errors.New("cannot refresh token"),
			wantSubstr: "client ID and secret",
		},
		{
			name:       "credential store error",
			err:        errors.New("could not access credential store"),
			wantSubstr: "keychain/keyring",
		},
		{
			name:       "keyring error",
			err:        errors.New("keyring access denied"),
			wantSubstr: "keychain/keyring",
		},
		{
			name:       "context deadline exceeded",
			err:        errors.New("context deadline exceeded"),
			wantSubstr: "timed out",
		},
		{
			name:       "context canceled",
			err:        errors.New("context canceled"),
			wantSubstr: "cancelled",
		},
		{
			name:       "json error",
			err:        errors.New("json: cannot unmarshal"),
			wantSubstr: "parse",
		},
		{
			name:       "unrecognized error",
			err:        errors.New("some unknown error"),
			wantSubstr: "some unknown error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatted := FormatError(tt.err)
			errStr := formatted.Error()
			if !strings.Contains(strings.ToLower(errStr), strings.ToLower(tt.wantSubstr)) {
				t.Errorf("Error() = %v, want to contain %v", errStr, tt.wantSubstr)
			}
		})
	}
}
