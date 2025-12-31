package threads

import (
	"context"
	"testing"
)

// TestDeletePost_InvalidPostID tests that DeletePost returns an error for empty post IDs
func TestDeletePost_InvalidPostID(t *testing.T) {
	client := &Client{}

	err := client.DeletePost(context.TODO(), ConvertToPostID(""))
	if err == nil {
		t.Error("expected error for empty post ID")
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
}

// TestDeletePostWithConfirmation_InvalidPostID tests that DeletePostWithConfirmation returns an error for empty post IDs
func TestDeletePostWithConfirmation_InvalidPostID(t *testing.T) {
	client := &Client{}

	confirmCalled := false
	callback := func(post *Post) bool {
		confirmCalled = true
		return true
	}

	err := client.DeletePostWithConfirmation(context.TODO(), ConvertToPostID(""), callback)
	if err == nil {
		t.Error("expected error for empty post ID")
		return
	}

	if confirmCalled {
		t.Error("callback should not be called for empty post ID")
	}

	validationErr, ok := err.(*ValidationError)
	if !ok {
		t.Errorf("expected ValidationError, got %T", err)
		return
	}

	if validationErr.Field != "post_id" {
		t.Errorf("expected field 'post_id', got '%s'", validationErr.Field)
	}
}

// TestDeletePostWithConfirmation_NilCallback tests that DeletePostWithConfirmation returns an error for nil callback
func TestDeletePostWithConfirmation_NilCallback(t *testing.T) {
	client := &Client{}

	err := client.DeletePostWithConfirmation(context.TODO(), ConvertToPostID("valid-post-id"), nil)
	if err == nil {
		t.Error("expected error for nil callback")
		return
	}

	validationErr, ok := err.(*ValidationError)
	if !ok {
		t.Errorf("expected ValidationError, got %T", err)
		return
	}

	if validationErr.Field != "confirmation_callback" {
		t.Errorf("expected field 'confirmation_callback', got '%s'", validationErr.Field)
	}
}
