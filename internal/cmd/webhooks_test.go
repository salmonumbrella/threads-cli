package cmd

import (
	"testing"

	threads "github.com/salmonumbrella/threads-go"
)

func TestWebhooksCmd_Structure(t *testing.T) {
	cmd := webhooksCmd

	if cmd.Use != "webhooks" {
		t.Errorf("expected Use=webhooks, got %s", cmd.Use)
	}

	if cmd.Short == "" {
		t.Error("expected Short description to be set")
	}

	if cmd.Long == "" {
		t.Error("expected Long description to be set")
	}
}

func TestWebhooksCmd_Subcommands(t *testing.T) {
	cmd := webhooksCmd

	expectedSubs := map[string]bool{
		"subscribe": true,
		"list":      true,
		"delete":    true,
	}

	for _, sub := range cmd.Commands() {
		name := sub.Name()
		if !expectedSubs[name] {
			t.Errorf("unexpected subcommand: %s", name)
		}
		delete(expectedSubs, name)
	}

	for name := range expectedSubs {
		t.Errorf("missing subcommand: %s", name)
	}
}

func TestWebhooksCmd_SubcommandCount(t *testing.T) {
	cmd := webhooksCmd
	subcommands := cmd.Commands()

	expectedCount := 3 // subscribe, list, delete
	if len(subcommands) != expectedCount {
		t.Errorf("expected %d subcommands, got %d", expectedCount, len(subcommands))
	}
}

func TestWebhooksSubscribeCmd_Structure(t *testing.T) {
	cmd := newWebhooksSubscribeCmd()

	if cmd.Use != "subscribe" {
		t.Errorf("expected Use=subscribe, got %s", cmd.Use)
	}

	if cmd.Short == "" {
		t.Error("expected Short description to be set")
	}

	if cmd.Long == "" {
		t.Error("expected Long description to be set")
	}

	if cmd.Example == "" {
		t.Error("expected Example to be set")
	}

	if cmd.RunE == nil {
		t.Error("expected RunE to be set")
	}
}

func TestWebhooksSubscribeCmd_Flags(t *testing.T) {
	cmd := newWebhooksSubscribeCmd()

	requiredFlags := []string{"url", "event"}
	for _, flag := range requiredFlags {
		f := cmd.Flag(flag)
		if f == nil {
			t.Errorf("missing flag: %s", flag)
			continue
		}
	}

	optionalFlags := []string{"verify-token"}
	for _, flag := range optionalFlags {
		if cmd.Flag(flag) == nil {
			t.Errorf("missing optional flag: %s", flag)
		}
	}
}

func TestWebhooksListCmd_Structure(t *testing.T) {
	cmd := newWebhooksListCmd()

	if cmd.Use != "list" {
		t.Errorf("expected Use=list, got %s", cmd.Use)
	}

	if cmd.RunE == nil {
		t.Error("expected RunE to be set")
	}
}

func TestWebhooksDeleteCmd_Structure(t *testing.T) {
	cmd := newWebhooksDeleteCmd()

	if cmd.Use != "delete [subscription-id]" {
		t.Errorf("expected Use='delete [subscription-id]', got %s", cmd.Use)
	}

	if cmd.Args == nil {
		t.Error("expected Args validator for exactly 1 arg")
	}

	if cmd.RunE == nil {
		t.Error("expected RunE to be set")
	}
}

func TestWebhooksDeleteCmd_HasExample(t *testing.T) {
	cmd := newWebhooksDeleteCmd()

	if cmd.Example == "" {
		t.Error("expected Example to be set")
	}
}

func TestWebhookSubscriptionToMap(t *testing.T) {
	sub := &threads.WebhookSubscription{
		ID:          "sub123",
		Object:      "user",
		CallbackURL: "https://example.com/webhook",
		Fields: []threads.WebhookField{
			{Name: "mentions", Version: "v1"},
			{Name: "publishes", Version: "v1"},
		},
		Active:      true,
		CreatedTime: "2024-01-01T00:00:00Z",
	}

	result := webhookSubscriptionToMap(sub)

	if result["id"] != "sub123" {
		t.Errorf("expected id=sub123, got %v", result["id"])
	}
	if result["object"] != "user" {
		t.Errorf("expected object=user, got %v", result["object"])
	}
	if result["callback_url"] != "https://example.com/webhook" {
		t.Errorf("expected callback_url to be set, got %v", result["callback_url"])
	}
	if result["active"] != true {
		t.Errorf("expected active=true, got %v", result["active"])
	}
	if result["created_time"] != "2024-01-01T00:00:00Z" {
		t.Errorf("expected created_time to be set, got %v", result["created_time"])
	}

	fields, ok := result["fields"].([]string)
	if !ok {
		t.Fatalf("expected fields to be []string, got %T", result["fields"])
	}
	if len(fields) != 2 {
		t.Errorf("expected 2 fields, got %d", len(fields))
	}
	if fields[0] != "mentions" || fields[1] != "publishes" {
		t.Errorf("unexpected fields: %v", fields)
	}
}

func TestWebhookSubscriptionToMap_EmptyFields(t *testing.T) {
	sub := &threads.WebhookSubscription{
		ID:          "sub456",
		Object:      "user",
		CallbackURL: "https://example.com/empty",
		Fields:      []threads.WebhookField{},
		Active:      false,
	}

	result := webhookSubscriptionToMap(sub)

	fields, ok := result["fields"].([]string)
	if !ok {
		t.Fatalf("expected fields to be []string, got %T", result["fields"])
	}
	if len(fields) != 0 {
		t.Errorf("expected 0 fields, got %d", len(fields))
	}
}

func TestFormatWebhookFields(t *testing.T) {
	tests := []struct {
		name     string
		fields   []threads.WebhookField
		expected string
	}{
		{
			name:     "empty fields",
			fields:   []threads.WebhookField{},
			expected: "-",
		},
		{
			name:     "nil fields",
			fields:   nil,
			expected: "-",
		},
		{
			name: "single field",
			fields: []threads.WebhookField{
				{Name: "mentions"},
			},
			expected: "mentions",
		},
		{
			name: "multiple fields",
			fields: []threads.WebhookField{
				{Name: "mentions"},
				{Name: "publishes"},
				{Name: "deletes"},
			},
			expected: "mentions, publishes, deletes",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatWebhookFields(tt.fields)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestTruncateURL(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		maxLen   int
		expected string
	}{
		{
			name:     "short URL unchanged",
			url:      "https://example.com",
			maxLen:   40,
			expected: "https://example.com",
		},
		{
			name:     "URL at exact length",
			url:      "https://example.com/webhook",
			maxLen:   27,
			expected: "https://example.com/webhook",
		},
		{
			name:     "long URL truncated",
			url:      "https://example.com/very/long/webhook/path/that/exceeds/limit",
			maxLen:   30,
			expected: "https://example.com/very/lo...",
		},
		{
			name:     "very short maxLen",
			url:      "https://example.com",
			maxLen:   10,
			expected: "https:/...",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := truncateURL(tt.url, tt.maxLen)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
			if len(result) > tt.maxLen {
				t.Errorf("result length %d exceeds maxLen %d", len(result), tt.maxLen)
			}
		})
	}
}

func TestTruncateURL_EdgeCases(t *testing.T) {
	// Empty URL
	result := truncateURL("", 10)
	if result != "" {
		t.Errorf("expected empty string for empty URL, got %q", result)
	}

	// URL shorter than ellipsis
	result = truncateURL("ab", 5)
	if result != "ab" {
		t.Errorf("expected 'ab', got %q", result)
	}
}
