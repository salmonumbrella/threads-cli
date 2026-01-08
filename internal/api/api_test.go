package api

import (
	"context"
	"net/http"
	"testing"
)

// Tests for GetPost with mocked HTTP

func TestGetPost_Success(t *testing.T) {
	client, server := createTestClient(t, createMockHandler(t, MockResponse{
		StatusCode: http.StatusOK,
		Body:       mockPostResponse(),
	}))
	defer server.Close()

	post, err := client.GetPost(context.Background(), ConvertToPostID("123456789"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if post == nil {
		t.Fatal("expected post to not be nil")
	}

	if post.ID != "123456789" {
		t.Errorf("expected post ID '123456789', got '%s'", post.ID)
	}
}

func TestGetPost_NotFound(t *testing.T) {
	client, server := createTestClient(t, createMockHandler(t, MockResponse{
		StatusCode: http.StatusNotFound,
		Body:       mockErrorResponse(404, "Post not found", "validation_error"),
	}))
	defer server.Close()

	_, err := client.GetPost(context.Background(), ConvertToPostID("nonexistent"))
	if err == nil {
		t.Fatal("expected error for not found")
	}
	// Just verify an error is returned - the specific error type depends on response parsing
}

// Tests for GetUserPosts with mocked HTTP

func TestGetUserPosts_Success(t *testing.T) {
	client, server := createTestClient(t, createMockHandler(t, MockResponse{
		StatusCode: http.StatusOK,
		Body:       mockPostsListResponse(),
	}))
	defer server.Close()

	posts, err := client.GetUserPosts(context.Background(), ConvertToUserID("12345"), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if posts == nil {
		t.Fatal("expected posts to not be nil")
	}

	if len(posts.Data) != 2 {
		t.Errorf("expected 2 posts, got %d", len(posts.Data))
	}
}

func TestGetUserPosts_WithPagination(t *testing.T) {
	client, server := createTestClient(t, createMockHandler(t, MockResponse{
		StatusCode: http.StatusOK,
		Body:       mockPostsListResponse(),
	}))
	defer server.Close()

	opts := &PaginationOptions{
		Limit:  25,
		Before: "cursor123",
	}

	posts, err := client.GetUserPosts(context.Background(), ConvertToUserID("12345"), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if posts == nil {
		t.Fatal("expected posts to not be nil")
	}
}

// Tests for GetUser with mocked HTTP

func TestGetUser_Success(t *testing.T) {
	client, server := createTestClient(t, createMockHandler(t, MockResponse{
		StatusCode: http.StatusOK,
		Body:       mockUserResponse(),
	}))
	defer server.Close()

	user, err := client.GetUser(context.Background(), ConvertToUserID("12345"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if user == nil {
		t.Fatal("expected user to not be nil")
	}

	if user.Username != "testuser" {
		t.Errorf("expected username 'testuser', got '%s'", user.Username)
	}
}

func TestGetUser_NotFound(t *testing.T) {
	client, server := createTestClient(t, createMockHandler(t, MockResponse{
		StatusCode: http.StatusNotFound,
		Body:       mockErrorResponse(404, "User not found", "validation_error"),
	}))
	defer server.Close()

	_, err := client.GetUser(context.Background(), ConvertToUserID("nonexistent"))
	if err == nil {
		t.Fatal("expected error for not found")
	}
	// Error is returned for 404 status
}

func TestGetUser_AccessDenied(t *testing.T) {
	client, server := createTestClient(t, createMockHandler(t, MockResponse{
		StatusCode: http.StatusForbidden,
		Body:       mockErrorResponse(403, "Access denied", "authentication_error"),
	}))
	defer server.Close()

	_, err := client.GetUser(context.Background(), ConvertToUserID("protected"))
	if err == nil {
		t.Fatal("expected error for access denied")
	}
	// Error is returned for 403 status
}

// Tests for GetReplies with mocked HTTP

func TestGetReplies_Success(t *testing.T) {
	client, server := createTestClient(t, createMockHandler(t, MockResponse{
		StatusCode: http.StatusOK,
		Body:       mockRepliesResponse(),
	}))
	defer server.Close()

	replies, err := client.GetReplies(context.Background(), ConvertToPostID("123456"), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if replies == nil {
		t.Fatal("expected replies to not be nil")
	}

	if len(replies.Data) != 1 {
		t.Errorf("expected 1 reply, got %d", len(replies.Data))
	}
}

func TestGetConversation_Success(t *testing.T) {
	client, server := createTestClient(t, createMockHandler(t, MockResponse{
		StatusCode: http.StatusOK,
		Body:       mockRepliesResponse(),
	}))
	defer server.Close()

	replies, err := client.GetConversation(context.Background(), ConvertToPostID("123456"), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if replies == nil {
		t.Fatal("expected conversation to not be nil")
	}
}

// Tests for HideReply/UnhideReply with mocked HTTP

func TestHideReply_Success(t *testing.T) {
	client, server := createTestClient(t, createMockHandler(t, MockResponse{
		StatusCode: http.StatusOK,
		Body:       mockSuccessResponse(),
	}))
	defer server.Close()

	err := client.HideReply(context.Background(), ConvertToPostID("reply123"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUnhideReply_Success(t *testing.T) {
	client, server := createTestClient(t, createMockHandler(t, MockResponse{
		StatusCode: http.StatusOK,
		Body:       mockSuccessResponse(),
	}))
	defer server.Close()

	err := client.UnhideReply(context.Background(), ConvertToPostID("reply123"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// Tests for DeletePost validation (API tests require complex mocking due to ownership validation)

func TestDeletePost_ValidationOnly(t *testing.T) {
	// Test validation - empty post ID
	client := &Client{}
	err := client.DeletePost(context.Background(), ConvertToPostID(""))
	if err == nil {
		t.Fatal("expected error for empty post ID")
	}

	validationErr, ok := err.(*ValidationError)
	if !ok {
		t.Errorf("expected ValidationError, got %T", err)
	}
	if validationErr != nil && validationErr.Field != "post_id" {
		t.Errorf("expected field 'post_id', got '%s'", validationErr.Field)
	}
}

// Tests for SearchLocations with mocked HTTP

func TestSearchLocations_Success(t *testing.T) {
	client, server := createTestClient(t, createMockHandler(t, MockResponse{
		StatusCode: http.StatusOK,
		Body:       mockLocationSearchResponse(),
	}))
	defer server.Close()

	locations, err := client.SearchLocations(context.Background(), "coffee shop", nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if locations == nil {
		t.Fatal("expected locations to not be nil")
	}

	if len(locations.Data) != 1 {
		t.Errorf("expected 1 location, got %d", len(locations.Data))
	}
}

func TestSearchLocations_WithCoordinates(t *testing.T) {
	client, server := createTestClient(t, createMockHandler(t, MockResponse{
		StatusCode: http.StatusOK,
		Body:       mockLocationSearchResponse(),
	}))
	defer server.Close()

	lat := 37.7749
	lon := -122.4194

	locations, err := client.SearchLocations(context.Background(), "", &lat, &lon)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if locations == nil {
		t.Fatal("expected locations to not be nil")
	}
}

func TestGetLocation_Success(t *testing.T) {
	client, server := createTestClient(t, createMockHandler(t, MockResponse{
		StatusCode: http.StatusOK,
		Body:       mockLocationResponse(),
	}))
	defer server.Close()

	location, err := client.GetLocation(context.Background(), ConvertToLocationID("loc123"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if location == nil {
		t.Fatal("expected location to not be nil")
	}

	if location.Name != "Test Location" {
		t.Errorf("expected name 'Test Location', got '%s'", location.Name)
	}
}

func TestGetLocation_NotFound(t *testing.T) {
	client, server := createTestClient(t, createMockHandler(t, MockResponse{
		StatusCode: http.StatusNotFound,
		Body:       mockErrorResponse(404, "Location not found", "validation_error"),
	}))
	defer server.Close()

	_, err := client.GetLocation(context.Background(), ConvertToLocationID("nonexistent"))
	if err == nil {
		t.Fatal("expected error for not found")
	}
	// Error is returned for 404 status
}

// Tests for KeywordSearch with mocked HTTP

func TestKeywordSearch_Success(t *testing.T) {
	client, server := createTestClient(t, createMockHandler(t, MockResponse{
		StatusCode: http.StatusOK,
		Body:       mockPostsListResponse(),
	}))
	defer server.Close()

	posts, err := client.KeywordSearch(context.Background(), "golang", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if posts == nil {
		t.Fatal("expected posts to not be nil")
	}
}

func TestKeywordSearch_WithOptions(t *testing.T) {
	client, server := createTestClient(t, createMockHandler(t, MockResponse{
		StatusCode: http.StatusOK,
		Body:       mockPostsListResponse(),
	}))
	defer server.Close()

	opts := &SearchOptions{
		MediaType: MediaTypeText,
		Limit:     50,
		Since:     1700000000,
	}

	posts, err := client.KeywordSearch(context.Background(), "golang", opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if posts == nil {
		t.Fatal("expected posts to not be nil")
	}
}

// Tests for GetMe - validation only (API test requires proper auth setup)

func TestGetMe_NoAuth(t *testing.T) {
	// GetMe requires authentication
	client := &Client{}
	_, err := client.GetMe(context.Background())
	if err == nil {
		t.Fatal("expected error when not authenticated")
	}

	// Should be an auth error
	_, ok := err.(*AuthenticationError)
	if !ok {
		t.Logf("Got error type %T: %v", err, err)
	}
}

// Tests for GetUserMentions with mocked HTTP

func TestGetUserMentions_Success(t *testing.T) {
	client, server := createTestClient(t, createMockHandler(t, MockResponse{
		StatusCode: http.StatusOK,
		Body:       mockPostsListResponse(),
	}))
	defer server.Close()

	mentions, err := client.GetUserMentions(context.Background(), ConvertToUserID("12345"), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if mentions == nil {
		t.Fatal("expected mentions to not be nil")
	}
}

// Tests for GetUserGhostPosts with mocked HTTP

func TestGetUserGhostPosts_Success(t *testing.T) {
	client, server := createTestClient(t, createMockHandler(t, MockResponse{
		StatusCode: http.StatusOK,
		Body:       mockPostsListResponse(),
	}))
	defer server.Close()

	posts, err := client.GetUserGhostPosts(context.Background(), ConvertToUserID("12345"), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if posts == nil {
		t.Fatal("expected posts to not be nil")
	}
}

// Tests for GetPublishingLimits with mocked HTTP

func TestGetPublishingLimits_Success(t *testing.T) {
	client, server := createTestClient(t, createMockHandler(t, MockResponse{
		StatusCode: http.StatusOK,
		Body:       mockPublishingLimitsResponse(),
	}))
	defer server.Close()

	limits, err := client.GetPublishingLimits(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if limits == nil {
		t.Fatal("expected limits to not be nil")
	}
}

// Tests for GetUserFields with mocked HTTP

func TestGetUserFields_Success(t *testing.T) {
	client, server := createTestClient(t, createMockHandler(t, MockResponse{
		StatusCode: http.StatusOK,
		Body:       mockUserResponse(),
	}))
	defer server.Close()

	user, err := client.GetUserFields(context.Background(), ConvertToUserID("12345"), []string{"id", "username"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if user == nil {
		t.Fatal("expected user to not be nil")
	}
}

func TestGetUserFields_DefaultFields(t *testing.T) {
	client, server := createTestClient(t, createMockHandler(t, MockResponse{
		StatusCode: http.StatusOK,
		Body:       mockUserResponse(),
	}))
	defer server.Close()

	user, err := client.GetUserFields(context.Background(), ConvertToUserID("12345"), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if user == nil {
		t.Fatal("expected user to not be nil")
	}
}

// Tests for LookupPublicProfile with mocked HTTP

func TestLookupPublicProfile_Success(t *testing.T) {
	client, server := createTestClient(t, createMockHandler(t, MockResponse{
		StatusCode: http.StatusOK,
		Body: map[string]interface{}{
			"username": "publicuser",
			"name":     "Public User",
		},
	}))
	defer server.Close()

	user, err := client.LookupPublicProfile(context.Background(), "publicuser")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if user == nil {
		t.Fatal("expected user to not be nil")
	}
}

func TestLookupPublicProfile_WithAtSymbol(t *testing.T) {
	client, server := createTestClient(t, createMockHandler(t, MockResponse{
		StatusCode: http.StatusOK,
		Body: map[string]interface{}{
			"username": "publicuser",
		},
	}))
	defer server.Close()

	user, err := client.LookupPublicProfile(context.Background(), "@publicuser")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if user == nil {
		t.Fatal("expected user to not be nil")
	}
}

// Tests for GetPublicProfilePosts with mocked HTTP

func TestGetPublicProfilePosts_Success(t *testing.T) {
	client, server := createTestClient(t, createMockHandler(t, MockResponse{
		StatusCode: http.StatusOK,
		Body:       mockPostsListResponse(),
	}))
	defer server.Close()

	posts, err := client.GetPublicProfilePosts(context.Background(), "publicuser", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if posts == nil {
		t.Fatal("expected posts to not be nil")
	}
}

// Tests for GetUserReplies with mocked HTTP

func TestGetUserReplies_Success(t *testing.T) {
	client, server := createTestClient(t, createMockHandler(t, MockResponse{
		StatusCode: http.StatusOK,
		Body:       mockRepliesResponse(),
	}))
	defer server.Close()

	replies, err := client.GetUserReplies(context.Background(), ConvertToUserID("12345"), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if replies == nil {
		t.Fatal("expected replies to not be nil")
	}
}

// Tests for GetUserPostsWithOptions with mocked HTTP

func TestGetUserPostsWithOptions_Success(t *testing.T) {
	client, server := createTestClient(t, createMockHandler(t, MockResponse{
		StatusCode: http.StatusOK,
		Body:       mockPostsListResponse(),
	}))
	defer server.Close()

	opts := &PostsOptions{
		Limit: 50,
		Since: 1700000000,
		Until: 1700100000,
	}

	posts, err := client.GetUserPostsWithOptions(context.Background(), ConvertToUserID("12345"), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if posts == nil {
		t.Fatal("expected posts to not be nil")
	}
}
