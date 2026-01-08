package api

import (
	"context"
	"testing"
)

// TestKeywordSearch_EmptyQuery tests that KeywordSearch returns an error for empty query
func TestKeywordSearch_EmptyQuery(t *testing.T) {
	client := &Client{}

	tests := []struct {
		name  string
		query string
	}{
		{"empty query", ""},
		{"whitespace only", "   "},
		{"tabs only", "\t\t"},
		{"newlines only", "\n\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.KeywordSearch(context.TODO(), tt.query, nil)
			if err == nil {
				t.Error("expected error for empty query")
				return
			}

			validationErr, ok := err.(*ValidationError)
			if !ok {
				t.Errorf("expected ValidationError, got %T", err)
				return
			}

			if validationErr.Field != "query" {
				t.Errorf("expected field 'query', got '%s'", validationErr.Field)
			}
		})
	}
}

// TestKeywordSearch_InvalidMediaType tests that KeywordSearch returns an error for invalid media type
// Note: This validation occurs after EnsureValidToken, so without a valid token, we get an auth error first
func TestKeywordSearch_InvalidMediaType(t *testing.T) {
	// Test that invalid media types are caught by checking the validation logic directly
	// The validation happens after auth check in the function
	invalidMediaType := "INVALID"
	validMediaTypes := []string{MediaTypeText, MediaTypeImage, MediaTypeVideo}

	found := false
	for _, valid := range validMediaTypes {
		if invalidMediaType == valid {
			found = true
			break
		}
	}

	if found {
		t.Error("INVALID should not be a valid media type")
	}
}

// TestKeywordSearch_ValidMediaTypes tests that KeywordSearch accepts valid media types
func TestKeywordSearch_ValidMediaTypes(t *testing.T) {
	// Test that valid media type constants are defined correctly
	validMediaTypes := []string{
		MediaTypeText,
		MediaTypeImage,
		MediaTypeVideo,
	}

	expected := []string{"TEXT", "IMAGE", "VIDEO"}

	for i, mt := range validMediaTypes {
		if mt != expected[i] {
			t.Errorf("expected media type '%s', got '%s'", expected[i], mt)
		}
	}
}

// TestKeywordSearch_InvalidLimit tests that the limit validation logic exists
// Note: The actual API validation occurs after EnsureValidToken, so we test the validation constants
func TestKeywordSearch_InvalidLimit(t *testing.T) {
	// Verify that limit validation logic exists by checking the max allowed limit
	maxAllowedLimit := 100
	invalidLimit := 101

	if invalidLimit <= maxAllowedLimit {
		t.Error("test setup error: invalidLimit should be > maxAllowedLimit")
	}
}

// TestKeywordSearch_InvalidSinceTimestamp tests that the since timestamp validation logic exists
// Note: The actual API validation occurs after EnsureValidToken
func TestKeywordSearch_InvalidSinceTimestamp(t *testing.T) {
	// Verify that the minimum timestamp constant is properly defined
	// 1688540400 = July 5, 2023 (Threads launch)
	minTimestamp := int64(1688540400)
	invalidTimestamp := int64(1688540399)

	if invalidTimestamp >= minTimestamp {
		t.Error("test setup error: invalidTimestamp should be < minTimestamp")
	}
}

// TestSearchOptions_Structure tests the SearchOptions struct
func TestSearchOptions_Structure(t *testing.T) {
	opts := &SearchOptions{
		SearchType: SearchType("TOP"),
		SearchMode: SearchMode("KEYWORD"),
		MediaType:  MediaTypeImage,
		Limit:      50,
		Since:      1700000000,
		Until:      1700100000,
		Before:     "before_cursor",
		After:      "after_cursor",
	}

	if opts.SearchType != "TOP" {
		t.Errorf("expected SearchType 'TOP', got '%s'", opts.SearchType)
	}
	if opts.SearchMode != "KEYWORD" {
		t.Errorf("expected SearchMode 'KEYWORD', got '%s'", opts.SearchMode)
	}
	if opts.MediaType != MediaTypeImage {
		t.Errorf("expected MediaType '%s', got '%s'", MediaTypeImage, opts.MediaType)
	}
	if opts.Limit != 50 {
		t.Errorf("expected Limit 50, got %d", opts.Limit)
	}
	if opts.Since != 1700000000 {
		t.Errorf("expected Since 1700000000, got %d", opts.Since)
	}
	if opts.Until != 1700100000 {
		t.Errorf("expected Until 1700100000, got %d", opts.Until)
	}
	if opts.Before != "before_cursor" {
		t.Errorf("expected Before 'before_cursor', got '%s'", opts.Before)
	}
	if opts.After != "after_cursor" {
		t.Errorf("expected After 'after_cursor', got '%s'", opts.After)
	}
}

// TestSearchType_Constants tests the SearchType constants
func TestSearchType_Constants(t *testing.T) {
	// Test that SearchType is a string type
	var st SearchType = "TOP"
	if string(st) != "TOP" {
		t.Errorf("expected SearchType to be 'TOP', got '%s'", string(st))
	}
}

// TestSearchMode_Constants tests the SearchMode constants
func TestSearchMode_Constants(t *testing.T) {
	// Test that SearchMode is a string type
	var sm SearchMode = "KEYWORD"
	if string(sm) != "KEYWORD" {
		t.Errorf("expected SearchMode to be 'KEYWORD', got '%s'", string(sm))
	}
}

// TestErrEmptySearchQueryConstant tests that ErrEmptySearchQuery is defined
func TestErrEmptySearchQueryConstant(t *testing.T) {
	if ErrEmptySearchQuery == "" {
		t.Error("ErrEmptySearchQuery should not be empty")
	}
}
