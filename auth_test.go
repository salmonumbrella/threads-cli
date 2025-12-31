package threads

import (
	"context"
	"testing"
	"time"
)

// TestExchangeCodeForToken_EmptyCode tests that ExchangeCodeForToken returns an error for empty code
func TestExchangeCodeForToken_EmptyCode(t *testing.T) {
	client := &Client{}

	err := client.ExchangeCodeForToken(context.TODO(), "")
	if err == nil {
		t.Error("expected error for empty code")
		return
	}

	validationErr, ok := err.(*ValidationError)
	if !ok {
		t.Errorf("expected ValidationError, got %T", err)
		return
	}

	if validationErr.Field != "code" {
		t.Errorf("expected field 'code', got '%s'", validationErr.Field)
	}
}

// TestGetLongLivedToken_NoToken tests that GetLongLivedToken returns an error when no token is set
func TestGetLongLivedToken_NoToken(t *testing.T) {
	client := &Client{}

	err := client.GetLongLivedToken(context.TODO())
	if err == nil {
		t.Error("expected error when no token is set")
		return
	}

	authErr, ok := err.(*AuthenticationError)
	if !ok {
		t.Errorf("expected AuthenticationError, got %T", err)
		return
	}

	if authErr.Code != 401 {
		t.Errorf("expected error code 401, got %d", authErr.Code)
	}
}

// TestRefreshToken_NoToken tests that RefreshToken returns an error when no token is set
func TestRefreshToken_NoToken(t *testing.T) {
	client := &Client{}

	err := client.RefreshToken(context.TODO())
	if err == nil {
		t.Error("expected error when no token is set")
		return
	}

	authErr, ok := err.(*AuthenticationError)
	if !ok {
		t.Errorf("expected AuthenticationError, got %T", err)
		return
	}

	if authErr.Code != 401 {
		t.Errorf("expected error code 401, got %d", authErr.Code)
	}
}

// TestGetAccessToken_Empty tests GetAccessToken returns empty string when no token set
func TestGetAccessToken_Empty(t *testing.T) {
	client := &Client{}

	token := client.GetAccessToken()
	if token != "" {
		t.Errorf("expected empty token, got '%s'", token)
	}
}

// TestDebugToken_NoToken tests that DebugToken returns an error when no token is set
func TestDebugToken_NoToken(t *testing.T) {
	client := &Client{}

	_, err := client.DebugToken(context.TODO(), "")
	if err == nil {
		t.Error("expected error when no token is set")
		return
	}

	authErr, ok := err.(*AuthenticationError)
	if !ok {
		t.Errorf("expected AuthenticationError, got %T", err)
		return
	}

	if authErr.Code != 401 {
		t.Errorf("expected error code 401, got %d", authErr.Code)
	}
}

// TestSetTokenFromDebugInfo_NilResponse tests that SetTokenFromDebugInfo returns an error for nil response
func TestSetTokenFromDebugInfo_NilResponse(t *testing.T) {
	client := &Client{}

	err := client.SetTokenFromDebugInfo("token", nil)
	if err == nil {
		t.Error("expected error for nil debug response")
	}
}

// TestSetTokenFromDebugInfo_InvalidToken tests that SetTokenFromDebugInfo returns an error for invalid token
func TestSetTokenFromDebugInfo_InvalidToken(t *testing.T) {
	client := &Client{}

	debugResp := &DebugTokenResponse{}
	debugResp.Data.IsValid = false

	err := client.SetTokenFromDebugInfo("token", debugResp)
	if err == nil {
		t.Error("expected error for invalid token")
		return
	}

	authErr, ok := err.(*AuthenticationError)
	if !ok {
		t.Errorf("expected AuthenticationError, got %T", err)
		return
	}

	if authErr.Code != 401 {
		t.Errorf("expected error code 401, got %d", authErr.Code)
	}
}

// TestTokenResponse_Structure tests the TokenResponse struct
func TestTokenResponse_Structure(t *testing.T) {
	resp := &TokenResponse{
		AccessToken: "test-token",
		TokenType:   "Bearer",
		ExpiresIn:   3600,
		UserID:      12345,
	}

	if resp.AccessToken != "test-token" {
		t.Errorf("expected AccessToken 'test-token', got '%s'", resp.AccessToken)
	}
	if resp.TokenType != "Bearer" {
		t.Errorf("expected TokenType 'Bearer', got '%s'", resp.TokenType)
	}
	if resp.ExpiresIn != 3600 {
		t.Errorf("expected ExpiresIn 3600, got %d", resp.ExpiresIn)
	}
	if resp.UserID != 12345 {
		t.Errorf("expected UserID 12345, got %d", resp.UserID)
	}
}

// TestLongLivedTokenResponse_Structure tests the LongLivedTokenResponse struct
func TestLongLivedTokenResponse_Structure(t *testing.T) {
	resp := &LongLivedTokenResponse{
		AccessToken: "long-lived-token",
		TokenType:   "Bearer",
		ExpiresIn:   5184000, // 60 days in seconds
	}

	if resp.AccessToken != "long-lived-token" {
		t.Errorf("expected AccessToken 'long-lived-token', got '%s'", resp.AccessToken)
	}
	if resp.TokenType != "Bearer" {
		t.Errorf("expected TokenType 'Bearer', got '%s'", resp.TokenType)
	}
	if resp.ExpiresIn != 5184000 {
		t.Errorf("expected ExpiresIn 5184000, got %d", resp.ExpiresIn)
	}
}

// TestDebugTokenResponse_Structure tests the DebugTokenResponse struct
func TestDebugTokenResponse_Structure(t *testing.T) {
	resp := &DebugTokenResponse{}
	resp.Data.Type = "USER"
	resp.Data.Application = "Test App"
	resp.Data.DataAccessExpiresAt = 1700000000
	resp.Data.ExpiresAt = 1700100000
	resp.Data.IsValid = true
	resp.Data.IssuedAt = 1699900000
	resp.Data.Scopes = []string{"threads_basic", "threads_content_publish"}
	resp.Data.UserID = "12345"

	if resp.Data.Type != "USER" {
		t.Errorf("expected Type 'USER', got '%s'", resp.Data.Type)
	}
	if resp.Data.Application != "Test App" {
		t.Errorf("expected Application 'Test App', got '%s'", resp.Data.Application)
	}
	if !resp.Data.IsValid {
		t.Error("expected IsValid to be true")
	}
	if len(resp.Data.Scopes) != 2 {
		t.Errorf("expected 2 scopes, got %d", len(resp.Data.Scopes))
	}
	if resp.Data.UserID != "12345" {
		t.Errorf("expected UserID '12345', got '%s'", resp.Data.UserID)
	}
}

// TestTokenInfo_Structure tests the TokenInfo struct
func TestTokenInfo_Structure(t *testing.T) {
	now := time.Now()
	expiresAt := now.Add(time.Hour)

	tokenInfo := &TokenInfo{
		AccessToken: "test-token",
		TokenType:   "Bearer",
		ExpiresAt:   expiresAt,
		UserID:      "12345",
		CreatedAt:   now,
	}

	if tokenInfo.AccessToken != "test-token" {
		t.Errorf("expected AccessToken 'test-token', got '%s'", tokenInfo.AccessToken)
	}
	if tokenInfo.TokenType != "Bearer" {
		t.Errorf("expected TokenType 'Bearer', got '%s'", tokenInfo.TokenType)
	}
	if tokenInfo.UserID != "12345" {
		t.Errorf("expected UserID '12345', got '%s'", tokenInfo.UserID)
	}
	if tokenInfo.ExpiresAt != expiresAt {
		t.Errorf("expected ExpiresAt %v, got %v", expiresAt, tokenInfo.ExpiresAt)
	}
	if tokenInfo.CreatedAt != now {
		t.Errorf("expected CreatedAt %v, got %v", now, tokenInfo.CreatedAt)
	}
}

// TestGetAuthURL tests the GetAuthURL method
func TestGetAuthURL(t *testing.T) {
	config := NewConfig()
	config.ClientID = "test-client-id"
	config.ClientSecret = "test-client-secret"
	config.RedirectURI = "https://example.com/callback"

	client, err := NewClient(config)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	tests := []struct {
		name   string
		scopes []string
	}{
		{"with default scopes", nil},
		{"with empty scopes", []string{}},
		{"with custom scopes", []string{"threads_basic"}},
		{"with multiple custom scopes", []string{"threads_basic", "threads_content_publish", "threads_manage_replies"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := client.GetAuthURL(tt.scopes)
			if url == "" {
				t.Error("expected non-empty URL")
				return
			}

			// Check that URL contains required components
			if !contains(url, "https://www.threads.net/oauth/authorize") {
				t.Error("URL should start with threads authorization endpoint")
			}
			if !contains(url, "client_id=test-client-id") {
				t.Error("URL should contain client_id")
			}
			if !contains(url, "response_type=code") {
				t.Error("URL should contain response_type=code")
			}
			if !contains(url, "state=") {
				t.Error("URL should contain state parameter")
			}
		})
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i+len(substr) <= len(s); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// TestGenerateState tests the generateState function indirectly through GetAuthURL
func TestGenerateState(t *testing.T) {
	config := NewConfig()
	config.ClientID = "test-client-id"
	config.ClientSecret = "test-client-secret"
	config.RedirectURI = "https://example.com/callback"

	client, err := NewClient(config)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Generate multiple URLs and check that states are different (random)
	url1 := client.GetAuthURL(nil)
	url2 := client.GetAuthURL(nil)

	// URLs should be different due to random state
	if url1 == url2 {
		t.Error("consecutive GetAuthURL calls should generate different states")
	}
}
