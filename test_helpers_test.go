package threads

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// MockResponse represents a mock HTTP response for testing
type MockResponse struct {
	StatusCode int
	Body       interface{}
	Headers    map[string]string
}

// noopTokenStorage is a token storage that does nothing (for testing)
type noopTokenStorage struct{}

func (n *noopTokenStorage) Store(tokenInfo *TokenInfo) error { return nil }
func (n *noopTokenStorage) Load() (*TokenInfo, error)        { return nil, nil }
func (n *noopTokenStorage) Delete() error                    { return nil }

// createTestClient creates a client configured to use a test server
func createTestClient(t *testing.T, handler http.HandlerFunc) (*Client, *httptest.Server) {
	t.Helper()

	server := httptest.NewServer(handler)

	config := NewConfig()
	config.ClientID = "test-client-id"
	config.ClientSecret = "test-client-secret"
	config.RedirectURI = "https://example.com/callback"

	client, err := NewClient(config)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Update the HTTP client to use the test server
	client.httpClient.baseURL = server.URL

	// Replace token storage with a no-op storage
	client.tokenStorage = &noopTokenStorage{}

	// Set a valid token to bypass authentication using the proper method
	tokenInfo := &TokenInfo{
		AccessToken: "test-access-token",
		TokenType:   "Bearer",
		ExpiresAt:   time.Now().Add(time.Hour),
		UserID:      "12345",
		CreatedAt:   time.Now(),
	}

	// Use SetTokenInfo which properly sets both tokenInfo and accessToken with mutex
	if err := client.SetTokenInfo(tokenInfo); err != nil {
		t.Fatalf("failed to set token info: %v", err)
	}

	return client, server
}

// createMockHandler creates a handler that responds with the given mock response
func createMockHandler(t *testing.T, mock MockResponse) http.HandlerFunc {
	t.Helper()

	return func(w http.ResponseWriter, r *http.Request) {
		// Set headers
		for key, value := range mock.Headers {
			w.Header().Set(key, value)
		}

		// Set status code
		w.WriteHeader(mock.StatusCode)

		// Write body
		if mock.Body != nil {
			body, err := json.Marshal(mock.Body)
			if err != nil {
				t.Errorf("failed to marshal mock body: %v", err)
				return
			}
			_, _ = w.Write(body)
		}
	}
}

// mockPostResponse creates a mock Post response
func mockPostResponse() map[string]interface{} {
	return map[string]interface{}{
		"id":         "123456789",
		"media_type": "TEXT",
		"text":       "Test post content",
		"username":   "testuser",
		"timestamp":  "2024-01-01T00:00:00Z",
	}
}

// mockPostsListResponse creates a mock posts list response
func mockPostsListResponse() map[string]interface{} {
	return map[string]interface{}{
		"data": []map[string]interface{}{
			{
				"id":         "123456789",
				"media_type": "TEXT",
				"text":       "Test post 1",
				"username":   "testuser",
			},
			{
				"id":         "987654321",
				"media_type": "IMAGE",
				"text":       "Test post 2",
				"username":   "testuser",
			},
		},
		"paging": map[string]interface{}{
			"cursors": map[string]interface{}{
				"before": "cursor_before",
				"after":  "cursor_after",
			},
		},
	}
}

// mockUserResponse creates a mock User response
func mockUserResponse() map[string]interface{} {
	return map[string]interface{}{
		"id":                          "12345",
		"username":                    "testuser",
		"threads_profile_picture_url": "https://example.com/pic.jpg",
		"threads_biography":           "Test bio",
	}
}

// mockRepliesResponse creates a mock replies list response
func mockRepliesResponse() map[string]interface{} {
	return map[string]interface{}{
		"data": []map[string]interface{}{
			{
				"id":         "reply1",
				"media_type": "TEXT",
				"text":       "Test reply 1",
				"username":   "replyuser",
			},
		},
	}
}

// mockLocationResponse creates a mock Location response
func mockLocationResponse() map[string]interface{} {
	return map[string]interface{}{
		"id":        "loc123",
		"name":      "Test Location",
		"latitude":  37.7749,
		"longitude": -122.4194,
	}
}

// mockLocationSearchResponse creates a mock location search response
func mockLocationSearchResponse() map[string]interface{} {
	return map[string]interface{}{
		"data": []map[string]interface{}{
			{
				"id":        "loc123",
				"name":      "Test Location 1",
				"latitude":  37.7749,
				"longitude": -122.4194,
			},
		},
	}
}

// mockErrorResponse creates a mock error response
func mockErrorResponse(code int, message, errorType string) map[string]interface{} {
	return map[string]interface{}{
		"error": map[string]interface{}{
			"message": message,
			"type":    errorType,
			"code":    code,
		},
	}
}

// mockPublishingLimitsResponse creates a mock publishing limits response
func mockPublishingLimitsResponse() map[string]interface{} {
	return map[string]interface{}{
		"data": []map[string]interface{}{
			{
				"config": map[string]interface{}{
					"quota_duration": 86400,
				},
				"quota_usage": 10,
			},
		},
	}
}

// mockSuccessResponse creates a simple success response
func mockSuccessResponse() map[string]interface{} {
	return map[string]interface{}{
		"success": true,
	}
}
