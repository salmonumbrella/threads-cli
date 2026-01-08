package api

import (
	"context"
	"testing"
)

// TestGetReplies_InvalidPostID tests that GetReplies returns an error for empty post IDs
func TestGetReplies_InvalidPostID(t *testing.T) {
	client := &Client{}

	_, err := client.GetReplies(context.TODO(), ConvertToPostID(""), nil)
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

// TestGetConversation_InvalidPostID tests that GetConversation returns an error for empty post IDs
func TestGetConversation_InvalidPostID(t *testing.T) {
	client := &Client{}

	_, err := client.GetConversation(context.TODO(), ConvertToPostID(""), nil)
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

// TestHideReply_InvalidReplyID tests that HideReply returns an error for empty reply IDs
func TestHideReply_InvalidReplyID(t *testing.T) {
	client := &Client{}

	err := client.HideReply(context.TODO(), ConvertToPostID(""))
	if err == nil {
		t.Error("expected error for empty reply ID")
		return
	}

	validationErr, ok := err.(*ValidationError)
	if !ok {
		t.Errorf("expected ValidationError, got %T", err)
		return
	}

	if validationErr.Field != "reply_id" {
		t.Errorf("expected field 'reply_id', got '%s'", validationErr.Field)
	}
}

// TestUnhideReply_InvalidReplyID tests that UnhideReply returns an error for empty reply IDs
func TestUnhideReply_InvalidReplyID(t *testing.T) {
	client := &Client{}

	err := client.UnhideReply(context.TODO(), ConvertToPostID(""))
	if err == nil {
		t.Error("expected error for empty reply ID")
		return
	}

	validationErr, ok := err.(*ValidationError)
	if !ok {
		t.Errorf("expected ValidationError, got %T", err)
		return
	}

	if validationErr.Field != "reply_id" {
		t.Errorf("expected field 'reply_id', got '%s'", validationErr.Field)
	}
}

// TestBuildRepliesParams tests the buildRepliesParams helper function
func TestBuildRepliesParams(t *testing.T) {
	tests := []struct {
		name        string
		opts        *RepliesOptions
		maxLimit    int
		limitDesc   string
		shouldErr   bool
		errField    string
		checkFields map[string]string
	}{
		{
			name:      "nil options",
			opts:      nil,
			maxLimit:  100,
			limitDesc: "replies per request",
			shouldErr: false,
		},
		{
			name:      "empty options",
			opts:      &RepliesOptions{},
			maxLimit:  100,
			limitDesc: "replies per request",
			shouldErr: false,
		},
		{
			name:      "valid limit",
			opts:      &RepliesOptions{Limit: 50},
			maxLimit:  100,
			limitDesc: "replies per request",
			shouldErr: false,
			checkFields: map[string]string{
				"limit": "50",
			},
		},
		{
			name:      "limit at max",
			opts:      &RepliesOptions{Limit: 100},
			maxLimit:  100,
			limitDesc: "replies per request",
			shouldErr: false,
			checkFields: map[string]string{
				"limit": "100",
			},
		},
		{
			name:      "limit exceeds max",
			opts:      &RepliesOptions{Limit: 101},
			maxLimit:  100,
			limitDesc: "replies per request",
			shouldErr: true,
			errField:  "limit",
		},
		{
			name:      "with before cursor",
			opts:      &RepliesOptions{Before: "cursor123"},
			maxLimit:  100,
			limitDesc: "replies per request",
			shouldErr: false,
			checkFields: map[string]string{
				"before": "cursor123",
			},
		},
		{
			name:      "with after cursor",
			opts:      &RepliesOptions{After: "cursor456"},
			maxLimit:  100,
			limitDesc: "replies per request",
			shouldErr: false,
			checkFields: map[string]string{
				"after": "cursor456",
			},
		},
		{
			name: "with reverse flag true",
			opts: func() *RepliesOptions {
				reverse := true
				return &RepliesOptions{Reverse: &reverse}
			}(),
			maxLimit:  100,
			limitDesc: "replies per request",
			shouldErr: false,
			checkFields: map[string]string{
				"reverse": "true",
			},
		},
		{
			name: "with reverse flag false",
			opts: func() *RepliesOptions {
				reverse := false
				return &RepliesOptions{Reverse: &reverse}
			}(),
			maxLimit:  100,
			limitDesc: "replies per request",
			shouldErr: false,
			checkFields: map[string]string{
				"reverse": "false",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params, err := buildRepliesParams(tt.opts, tt.maxLimit, tt.limitDesc)

			if tt.shouldErr {
				if err == nil {
					t.Error("expected error but got none")
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
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			// Check that fields parameter always exists
			if params.Get("fields") == "" {
				t.Error("expected fields parameter to be set")
			}

			// Check specific field values
			for key, expected := range tt.checkFields {
				actual := params.Get(key)
				if actual != expected {
					t.Errorf("expected %s='%s', got '%s'", key, expected, actual)
				}
			}
		})
	}
}

// TestReplyFieldsConstant tests that ReplyFields is defined
func TestReplyFieldsConstant(t *testing.T) {
	if ReplyFields == "" {
		t.Error("ReplyFields should not be empty")
	}
}

// TestRepliesOptions_Structure tests the RepliesOptions structure
func TestRepliesOptions_Structure(t *testing.T) {
	reverse := true
	opts := &RepliesOptions{
		Limit:   50,
		Before:  "before_cursor",
		After:   "after_cursor",
		Reverse: &reverse,
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
	if opts.Reverse == nil || !*opts.Reverse {
		t.Error("expected Reverse to be true")
	}
}
