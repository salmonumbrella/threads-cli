package api

import (
	"context"
	"testing"
)

// TestGetPost_InvalidPostID tests that GetPost returns an error for empty post IDs
func TestGetPost_InvalidPostID(t *testing.T) {
	client := &Client{}

	tests := []struct {
		name   string
		postID PostID
	}{
		{"empty post ID", ConvertToPostID("")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.GetPost(context.TODO(), tt.postID)
			if err == nil {
				t.Error("expected error for invalid post ID")
				return
			}

			validationErr, ok := err.(*ValidationError)
			if !ok {
				t.Errorf("expected ValidationError, got %T", err)
				return
			}

			if validationErr.Field != "post_id" {
				t.Errorf("expected field 'post_id', got '%s'", validationErr.Field)
			}
		})
	}
}

// TestGetUserPosts_InvalidUserID tests that GetUserPosts returns an error for empty user IDs
func TestGetUserPosts_InvalidUserID(t *testing.T) {
	client := &Client{}

	_, err := client.GetUserPosts(context.TODO(), ConvertToUserID(""), nil)
	if err == nil {
		t.Error("expected error for empty user ID")
		return
	}

	validationErr, ok := err.(*ValidationError)
	if !ok {
		t.Errorf("expected ValidationError, got %T", err)
		return
	}

	if validationErr.Field != "user_id" {
		t.Errorf("expected field 'user_id', got '%s'", validationErr.Field)
	}
}

// TestGetUserPostsWithOptions_InvalidUserID tests that GetUserPostsWithOptions returns an error for empty user IDs
func TestGetUserPostsWithOptions_InvalidUserID(t *testing.T) {
	client := &Client{}

	_, err := client.GetUserPostsWithOptions(context.TODO(), ConvertToUserID(""), nil)
	if err == nil {
		t.Error("expected error for empty user ID")
		return
	}

	validationErr, ok := err.(*ValidationError)
	if !ok {
		t.Errorf("expected ValidationError, got %T", err)
		return
	}

	if validationErr.Field != "user_id" {
		t.Errorf("expected field 'user_id', got '%s'", validationErr.Field)
	}
}

// TestGetUserMentions_InvalidUserID tests that GetUserMentions returns an error for empty user IDs
func TestGetUserMentions_InvalidUserID(t *testing.T) {
	client := &Client{}

	_, err := client.GetUserMentions(context.TODO(), ConvertToUserID(""), nil)
	if err == nil {
		t.Error("expected error for empty user ID")
		return
	}

	validationErr, ok := err.(*ValidationError)
	if !ok {
		t.Errorf("expected ValidationError, got %T", err)
		return
	}

	if validationErr.Field != "user_id" {
		t.Errorf("expected field 'user_id', got '%s'", validationErr.Field)
	}
}

// TestGetUserGhostPosts_InvalidUserID tests that GetUserGhostPosts returns an error for empty user IDs
func TestGetUserGhostPosts_InvalidUserID(t *testing.T) {
	client := &Client{}

	_, err := client.GetUserGhostPosts(context.TODO(), ConvertToUserID(""), nil)
	if err == nil {
		t.Error("expected error for empty user ID")
		return
	}

	validationErr, ok := err.(*ValidationError)
	if !ok {
		t.Errorf("expected ValidationError, got %T", err)
		return
	}

	if validationErr.Field != "user_id" {
		t.Errorf("expected field 'user_id', got '%s'", validationErr.Field)
	}
}

// TestPostsOptions_ValidPaginationOptions tests that valid pagination options are accepted
func TestPostsOptions_ValidPaginationOptions(t *testing.T) {
	validator := NewValidator()

	tests := []struct {
		name string
		opts *PaginationOptions
	}{
		{"nil options", nil},
		{"empty options", &PaginationOptions{}},
		{"valid limit", &PaginationOptions{Limit: 25}},
		{"valid before cursor", &PaginationOptions{Before: "cursor123"}},
		{"valid after cursor", &PaginationOptions{After: "cursor456"}},
		{"combined options", &PaginationOptions{Limit: 50, After: "cursor789"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidatePaginationOptions(tt.opts)
			if err != nil {
				t.Errorf("expected no error for valid options, got: %v", err)
			}
		})
	}
}

// TestPostsOptions_InvalidPaginationOptions tests that invalid pagination options are rejected
func TestPostsOptions_InvalidPaginationOptions(t *testing.T) {
	validator := NewValidator()

	tests := []struct {
		name     string
		opts     *PaginationOptions
		errField string
	}{
		{"limit too high", &PaginationOptions{Limit: 1000}, "limit"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidatePaginationOptions(tt.opts)
			if err == nil {
				t.Error("expected error for invalid options")
				return
			}

			validationErr, ok := err.(*ValidationError)
			if !ok {
				t.Errorf("expected ValidationError, got %T", err)
				return
			}

			if validationErr.Field != tt.errField {
				t.Errorf("expected field '%s', got '%s'", tt.errField, validationErr.Field)
			}
		})
	}
}

// TestPostExtendedFieldsConstant tests that PostExtendedFields is defined
func TestPostExtendedFieldsConstant(t *testing.T) {
	if PostExtendedFields == "" {
		t.Error("PostExtendedFields should not be empty")
	}
}

// TestGhostPostFieldsConstant tests that GhostPostFields is defined
func TestGhostPostFieldsConstant(t *testing.T) {
	if GhostPostFields == "" {
		t.Error("GhostPostFields should not be empty")
	}
}

// TestPublishingLimitFieldsConstant tests that PublishingLimitFields is defined
func TestPublishingLimitFieldsConstant(t *testing.T) {
	if PublishingLimitFields == "" {
		t.Error("PublishingLimitFields should not be empty")
	}
}
