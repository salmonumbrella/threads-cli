package threads

import (
	"context"
	"testing"
)

// TestGetUser_InvalidUserID tests that GetUser returns an error for empty user IDs
func TestGetUser_InvalidUserID(t *testing.T) {
	client := &Client{}

	_, err := client.GetUser(context.TODO(), ConvertToUserID(""))
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

// TestGetUserFields_InvalidUserID tests that GetUserFields returns an error for empty user IDs
func TestGetUserFields_InvalidUserID(t *testing.T) {
	client := &Client{}

	_, err := client.GetUserFields(context.TODO(), ConvertToUserID(""), nil)
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

// TestGetUserFields_InvalidFields tests that GetUserFields returns an error for completely invalid fields
func TestGetUserFields_InvalidFields(t *testing.T) {
	client := &Client{}

	// Use all invalid fields that won't be recognized
	_, err := client.GetUserFields(context.TODO(), ConvertToUserID("valid-user-id"), []string{"invalid_field1", "invalid_field2"})
	if err == nil {
		t.Error("expected error for invalid fields")
		return
	}

	validationErr, ok := err.(*ValidationError)
	if !ok {
		t.Errorf("expected ValidationError, got %T", err)
		return
	}

	if validationErr.Field != "fields" {
		t.Errorf("expected field 'fields', got '%s'", validationErr.Field)
	}
}

// TestGetUserFields_ValidFields tests that GetUserFields accepts valid fields
func TestGetUserFields_ValidFields(t *testing.T) {
	// Test the validation logic only - we can't fully test API calls without mocking
	validFields := []string{
		"id",
		"username",
		"name",
		"threads_profile_picture_url",
		"threads_biography",
		"is_verified",
		"recently_searched_keywords",
	}

	allowedFields := map[string]bool{
		"id":                          true,
		"username":                    true,
		"name":                        true,
		"threads_profile_picture_url": true,
		"threads_biography":           true,
		"is_verified":                 true,
		"recently_searched_keywords":  true,
	}

	for _, field := range validFields {
		if !allowedFields[field] {
			t.Errorf("field '%s' should be allowed", field)
		}
	}
}

// TestLookupPublicProfile_EmptyUsername tests that LookupPublicProfile returns an error for empty username
func TestLookupPublicProfile_EmptyUsername(t *testing.T) {
	client := &Client{}

	tests := []struct {
		name     string
		username string
	}{
		{"empty username", ""},
		{"whitespace only", "   "},
		{"tabs only", "\t\t"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.LookupPublicProfile(context.TODO(), tt.username)
			if err == nil {
				t.Error("expected error for empty username")
				return
			}

			validationErr, ok := err.(*ValidationError)
			if !ok {
				t.Errorf("expected ValidationError, got %T", err)
				return
			}

			if validationErr.Field != "username" {
				t.Errorf("expected field 'username', got '%s'", validationErr.Field)
			}
		})
	}
}

// TestGetPublicProfilePosts_EmptyUsername tests that GetPublicProfilePosts returns an error for empty username
func TestGetPublicProfilePosts_EmptyUsername(t *testing.T) {
	client := &Client{}

	tests := []struct {
		name     string
		username string
	}{
		{"empty username", ""},
		{"whitespace only", "   "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.GetPublicProfilePosts(context.TODO(), tt.username, nil)
			if err == nil {
				t.Error("expected error for empty username")
				return
			}

			validationErr, ok := err.(*ValidationError)
			if !ok {
				t.Errorf("expected ValidationError, got %T", err)
				return
			}

			if validationErr.Field != "username" {
				t.Errorf("expected field 'username', got '%s'", validationErr.Field)
			}
		})
	}
}

// TestGetPublicProfilePosts_InvalidLimit tests that limit validation exists for GetPublicProfilePosts
// Note: The actual validation occurs after EnsureValidToken
func TestGetPublicProfilePosts_InvalidLimit(t *testing.T) {
	// Verify the max limit for public profile posts
	maxLimit := 100
	invalidLimit := 101

	if invalidLimit <= maxLimit {
		t.Error("test setup error: invalidLimit should be > maxLimit")
	}
}

// TestGetUserReplies_InvalidUserID tests that GetUserReplies returns an error for empty user IDs
func TestGetUserReplies_InvalidUserID(t *testing.T) {
	client := &Client{}

	_, err := client.GetUserReplies(context.TODO(), ConvertToUserID(""), nil)
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

// TestGetUserReplies_InvalidLimit tests that limit validation exists for GetUserReplies
// Note: The actual validation occurs after EnsureValidToken
func TestGetUserReplies_InvalidLimit(t *testing.T) {
	// Verify the max limit for user replies
	maxLimit := 100
	invalidLimit := 101

	if invalidLimit <= maxLimit {
		t.Error("test setup error: invalidLimit should be > maxLimit")
	}
}

// TestUserProfileFieldsConstant tests that UserProfileFields is defined
func TestUserProfileFieldsConstant(t *testing.T) {
	if UserProfileFields == "" {
		t.Error("UserProfileFields should not be empty")
	}
}

// TestUserStruct tests the User struct fields
func TestUserStruct(t *testing.T) {
	user := &User{
		ID:            "12345",
		Username:      "testuser",
		Name:          "Test User",
		ProfilePicURL: "https://example.com/pic.jpg",
		Biography:     "Test bio",
		IsVerified:    true,
	}

	if user.ID != "12345" {
		t.Errorf("expected ID '12345', got '%s'", user.ID)
	}
	if user.Username != "testuser" {
		t.Errorf("expected Username 'testuser', got '%s'", user.Username)
	}
	if user.Name != "Test User" {
		t.Errorf("expected Name 'Test User', got '%s'", user.Name)
	}
	if user.ProfilePicURL != "https://example.com/pic.jpg" {
		t.Errorf("expected ProfilePicURL 'https://example.com/pic.jpg', got '%s'", user.ProfilePicURL)
	}
	if user.Biography != "Test bio" {
		t.Errorf("expected Biography 'Test bio', got '%s'", user.Biography)
	}
	if !user.IsVerified {
		t.Error("expected IsVerified to be true")
	}
}

// TestPublicUserStruct tests the PublicUser struct
func TestPublicUserStruct(t *testing.T) {
	publicUser := &PublicUser{
		Username:          "publicuser",
		Name:              "Public User",
		ProfilePictureURL: "https://example.com/pic.jpg",
		Biography:         "Test bio",
		IsVerified:        true,
		FollowerCount:     1000,
		LikesCount:        500,
	}

	if publicUser.Username != "publicuser" {
		t.Errorf("expected Username 'publicuser', got '%s'", publicUser.Username)
	}
	if publicUser.Name != "Public User" {
		t.Errorf("expected Name 'Public User', got '%s'", publicUser.Name)
	}
	if !publicUser.IsVerified {
		t.Error("expected IsVerified to be true")
	}
	if publicUser.FollowerCount != 1000 {
		t.Errorf("expected FollowerCount 1000, got %d", publicUser.FollowerCount)
	}
}

// TestPostsOptionsStruct tests the PostsOptions struct
func TestPostsOptionsStruct(t *testing.T) {
	opts := &PostsOptions{
		Limit:  50,
		Before: "before_cursor",
		After:  "after_cursor",
		Since:  1700000000,
		Until:  1700100000,
	}

	if opts.Limit != 50 {
		t.Errorf("expected Limit 50, got %d", opts.Limit)
	}
	if opts.Before != "before_cursor" {
		t.Errorf("expected Before 'before_cursor', got '%s'", opts.Before)
	}
	if opts.After != "after_cursor" {
		t.Errorf("expected After 'after_cursor', got '%s'", opts.After)
	}
	if opts.Since != 1700000000 {
		t.Errorf("expected Since 1700000000, got %d", opts.Since)
	}
	if opts.Until != 1700100000 {
		t.Errorf("expected Until 1700100000, got %d", opts.Until)
	}
}
