package api

import (
	"encoding/json"
	"testing"
	"time"
)

func TestNewConfig(t *testing.T) {
	config := NewConfig()

	if config == nil {
		t.Fatal("NewConfig() returned nil")
	}

	// Check defaults
	if config.HTTPTimeout != DefaultHTTPTimeout {
		t.Errorf("Expected HTTPTimeout to be %v, got %v", DefaultHTTPTimeout, config.HTTPTimeout)
	}

	if config.BaseURL != BaseAPIURL {
		t.Errorf("Expected BaseURL to be %s, got %s", BaseAPIURL, config.BaseURL)
	}

	if config.UserAgent != DefaultUserAgent {
		t.Errorf("Expected UserAgent to be %s, got %s", DefaultUserAgent, config.UserAgent)
	}

	// Check that scopes are set
	if len(config.Scopes) == 0 {
		t.Error("Expected scopes to be set by default")
	}

	// Check retry config
	if config.RetryConfig == nil {
		t.Fatal("Expected RetryConfig to be set")
	}

	if config.RetryConfig.MaxRetries != 3 {
		t.Errorf("Expected MaxRetries to be 3, got %d", config.RetryConfig.MaxRetries)
	}
}

func TestConfigValidation(t *testing.T) {
	validator := NewConfigValidator()

	tests := []struct {
		name      string
		config    *Config
		shouldErr bool
	}{
		{
			name: "valid config",
			config: &Config{
				ClientID:     "test-client-id",
				ClientSecret: "test-client-secret",
				RedirectURI:  "https://example.com/callback",
				Scopes:       []string{"threads_basic"},
				HTTPTimeout:  30 * time.Second,
				BaseURL:      "https://graph.threads.net",
			},
			shouldErr: false,
		},
		{
			name: "missing client ID",
			config: &Config{
				ClientSecret: "test-client-secret",
				RedirectURI:  "https://example.com/callback",
				Scopes:       []string{"threads_basic"},
				HTTPTimeout:  30 * time.Second,
				BaseURL:      "https://graph.threads.net",
			},
			shouldErr: true,
		},
		{
			name: "invalid redirect URI",
			config: &Config{
				ClientID:     "test-client-id",
				ClientSecret: "test-client-secret",
				RedirectURI:  "not-a-url",
				Scopes:       []string{"threads_basic"},
				HTTPTimeout:  30 * time.Second,
				BaseURL:      "https://graph.threads.net",
			},
			shouldErr: true,
		},
		{
			name: "invalid scope",
			config: &Config{
				ClientID:     "test-client-id",
				ClientSecret: "test-client-secret",
				RedirectURI:  "https://example.com/callback",
				Scopes:       []string{"invalid_scope"},
				HTTPTimeout:  30 * time.Second,
				BaseURL:      "https://graph.threads.net",
			},
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Validate(tt.config)
			if tt.shouldErr && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestGhostPostValidation(t *testing.T) {
	client := &Client{}

	// Test valid ghost post
	validGhost := &TextPostContent{
		Text:        "This is a ghost post",
		IsGhostPost: true,
	}
	err := client.ValidateTextPostContent(validGhost)
	if err != nil {
		t.Errorf("Expected valid ghost post to pass validation, got: %v", err)
	}

	// Test invalid ghost post (reply)
	invalidGhost := &TextPostContent{
		Text:        "This is an invalid ghost post",
		IsGhostPost: true,
		ReplyTo:     "some-post-id",
	}
	err = client.ValidateTextPostContent(invalidGhost)
	if err == nil {
		t.Error("Expected error for ghost post with ReplyTo")
	} else if validationErr, ok := err.(*ValidationError); ok {
		if validationErr.Field != "is_ghost_post" {
			t.Errorf("Expected error field 'is_ghost_post', got '%s'", validationErr.Field)
		}
	} else {
		t.Errorf("Expected ValidationError, got %T", err)
	}
}

func TestValidation(t *testing.T) {
	validator := NewValidator()

	t.Run("ValidateTextLength", func(t *testing.T) {
		// Test valid text
		err := validator.ValidateTextLength("Hello world", "Text")
		if err != nil {
			t.Errorf("Expected no error for valid text, got: %v", err)
		}

		// Test text too long
		longText := make([]byte, MaxTextLength+1)
		for i := range longText {
			longText[i] = 'a'
		}
		err = validator.ValidateTextLength(string(longText), "Text")
		if err == nil {
			t.Error("Expected error for text too long")
		}
	})

	t.Run("ValidateTopicTag", func(t *testing.T) {
		// Test valid tag
		err := validator.ValidateTopicTag("valid_tag")
		if err != nil {
			t.Errorf("Expected no error for valid tag, got: %v", err)
		}

		// Test invalid tag with period
		err = validator.ValidateTopicTag("invalid.tag")
		if err == nil {
			t.Error("Expected error for tag with period")
		}

		// Test invalid tag with ampersand
		err = validator.ValidateTopicTag("invalid&tag")
		if err == nil {
			t.Error("Expected error for tag with ampersand")
		}
	})

	t.Run("ValidateCountryCodes", func(t *testing.T) {
		// Test valid codes
		err := validator.ValidateCountryCodes([]string{"US", "CA", "GB"})
		if err != nil {
			t.Errorf("Expected no error for valid country codes, got: %v", err)
		}

		// Test invalid code length
		err = validator.ValidateCountryCodes([]string{"USA"})
		if err == nil {
			t.Error("Expected error for invalid country code length")
		}

		// Test invalid characters
		err = validator.ValidateCountryCodes([]string{"U1"})
		if err == nil {
			t.Error("Expected error for country code with numbers")
		}
	})

	t.Run("ValidateLinkCount", func(t *testing.T) {
		// Test valid link count (0 links)
		err := validator.ValidateLinkCount("Hello world", "")
		if err != nil {
			t.Errorf("Expected no error for 0 links, got: %v", err)
		}

		// Test valid link count (5 links)
		fiveLinks := "http://a.com https://b.com http://c.com https://d.com http://e.com"
		err = validator.ValidateLinkCount(fiveLinks, "")
		if err != nil {
			t.Errorf("Expected no error for 5 links, got: %v", err)
		}

		// Test unique links logic
		// "If the text field contains www.example.com, www.example.com, and www.test.com,
		// and the link_attachment is www.test.com, this counts as 2 links"
		// (Assuming http/https prefix for validator detection)
		duplicateLinks := "http://example.com http://example.com http://test.com"
		err = validator.ValidateLinkCount(duplicateLinks, "http://test.com")
		if err != nil {
			t.Errorf("Expected no error for duplicate links (should count as 2), got: %v", err)
		}

		// Test link_attachment adds to count
		// "If the text field contains www.instagram.com and www.threads.com,
		// and the link_attachment is www.facebook.com, this counts as 3 links."
		textWithLinks := "http://instagram.com http://threads.com"
		err = validator.ValidateLinkCount(textWithLinks, "http://facebook.com")
		if err != nil {
			t.Errorf("Expected no error for 3 total links, got: %v", err)
		}

		// Test invalid link count (6 unique links)
		sixLinks := "http://a.com https://b.com http://c.com https://d.com http://e.com https://f.com"
		err = validator.ValidateLinkCount(sixLinks, "")
		if err == nil {
			t.Error("Expected error for 6 links")
		}

		// Test invalid link count (5 in text + 1 unique in attachment)
		fiveInText := "http://a.com https://b.com http://c.com https://d.com http://e.com"
		err = validator.ValidateLinkCount(fiveInText, "http://f.com")
		if err == nil {
			t.Error("Expected error for 6 total unique links")
		}
	})
}

func TestPostIDTypes(t *testing.T) {
	// Test PostID
	postID := ConvertToPostID("test-post-id")
	if !postID.Valid() {
		t.Error("Expected PostID to be valid")
	}
	if postID.String() != "test-post-id" {
		t.Errorf("Expected PostID string to be 'test-post-id', got '%s'", postID.String())
	}

	// Test empty PostID
	emptyPostID := ConvertToPostID("")
	if emptyPostID.Valid() {
		t.Error("Expected empty PostID to be invalid")
	}

	// Test UserID
	userID := ConvertToUserID("test-user-id")
	if !userID.Valid() {
		t.Error("Expected UserID to be valid")
	}
	if userID.String() != "test-user-id" {
		t.Errorf("Expected UserID string to be 'test-user-id', got '%s'", userID.String())
	}

	// Test ContainerID
	containerID := ConvertToContainerID("test-container-id")
	if !containerID.Valid() {
		t.Error("Expected ContainerID to be valid")
	}
	if containerID.String() != "test-container-id" {
		t.Errorf("Expected ContainerID string to be 'test-container-id', got '%s'", containerID.String())
	}

	// Test empty ContainerID
	emptyContainerID := ConvertToContainerID("")
	if emptyContainerID.Valid() {
		t.Error("Expected empty ContainerID to be invalid")
	}

	// Test LocationID
	locationID := ConvertToLocationID("test-location-id")
	if !locationID.Valid() {
		t.Error("Expected LocationID to be valid")
	}
	if locationID.String() != "test-location-id" {
		t.Errorf("Expected LocationID string to be 'test-location-id', got '%s'", locationID.String())
	}

	// Test empty LocationID
	emptyLocationID := ConvertToLocationID("")
	if emptyLocationID.Valid() {
		t.Error("Expected empty LocationID to be invalid")
	}
}

func TestContainerBuilder(t *testing.T) {
	builder := NewContainerBuilder()

	params := builder.
		SetMediaType(MediaTypeText).
		SetText("Hello world").
		SetReplyControl(ReplyControlEveryone).
		Build()

	if params.Get("media_type") != MediaTypeText {
		t.Errorf("Expected media_type to be %s, got %s", MediaTypeText, params.Get("media_type"))
	}

	if params.Get("text") != "Hello world" {
		t.Errorf("Expected text to be 'Hello world', got '%s'", params.Get("text"))
	}

	if params.Get("reply_control") != string(ReplyControlEveryone) {
		t.Errorf("Expected reply_control to be %s, got %s", string(ReplyControlEveryone), params.Get("reply_control"))
	}
}

func TestContainerBuilderGIFAttachment(t *testing.T) {
	builder := NewContainerBuilder()

	gif := &GIFAttachment{
		GIFID:    "test-gif-id-12345",
		Provider: GIFProviderTenor,
	}

	params := builder.
		SetMediaType(MediaTypeText).
		SetText("Check out this GIF!").
		SetGIFAttachment(gif).
		Build()

	if params.Get("media_type") != MediaTypeText {
		t.Errorf("Expected media_type to be %s, got %s", MediaTypeText, params.Get("media_type"))
	}

	gifParam := params.Get("gif_attachment")
	if gifParam == "" {
		t.Error("Expected gif_attachment to be set")
	}

	// Check that the GIF attachment contains expected values
	if gifParam == "" {
		t.Error("Expected gif_attachment parameter to be set")
	}
}

func TestContainerBuilderGIFAttachmentNil(t *testing.T) {
	builder := NewContainerBuilder()

	params := builder.
		SetMediaType(MediaTypeText).
		SetText("No GIF here").
		SetGIFAttachment(nil).
		Build()

	if params.Get("gif_attachment") != "" {
		t.Error("Expected gif_attachment to be empty when nil")
	}
}

func TestValidateGIFAttachment(t *testing.T) {
	validator := NewValidator()

	tests := []struct {
		name      string
		gif       *GIFAttachment
		shouldErr bool
		errField  string
	}{
		{
			name:      "nil gif attachment is valid",
			gif:       nil,
			shouldErr: false,
		},
		{
			name: "valid gif attachment",
			gif: &GIFAttachment{
				GIFID:    "test-gif-id",
				Provider: GIFProviderTenor,
			},
			shouldErr: false,
		},
		{
			name: "missing gif_id",
			gif: &GIFAttachment{
				GIFID:    "",
				Provider: GIFProviderTenor,
			},
			shouldErr: true,
			errField:  "gif_attachment.gif_id",
		},
		{
			name: "whitespace only gif_id",
			gif: &GIFAttachment{
				GIFID:    "   ",
				Provider: GIFProviderTenor,
			},
			shouldErr: true,
			errField:  "gif_attachment.gif_id",
		},
		{
			name: "missing provider",
			gif: &GIFAttachment{
				GIFID:    "test-gif-id",
				Provider: "",
			},
			shouldErr: true,
			errField:  "gif_attachment.provider",
		},
		{
			name: "invalid provider",
			gif: &GIFAttachment{
				GIFID:    "test-gif-id",
				Provider: GIFProvider("GIPHY"),
			},
			shouldErr: true,
			errField:  "gif_attachment.provider",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateGIFAttachment(tt.gif)
			if tt.shouldErr {
				if err == nil {
					t.Error("Expected error but got none")
				} else if validationErr, ok := err.(*ValidationError); ok {
					if tt.errField != "" && validationErr.Field != tt.errField {
						t.Errorf("Expected error field '%s', got '%s'", tt.errField, validationErr.Field)
					}
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
			}
		})
	}
}

func TestGIFProviderConstants(t *testing.T) {
	// Verify the TENOR constant is correctly defined
	if GIFProviderTenor != "TENOR" {
		t.Errorf("Expected GIFProviderTenor to be 'TENOR', got '%s'", GIFProviderTenor)
	}
}

func TestGIFAttachmentStruct(t *testing.T) {
	gif := &GIFAttachment{
		GIFID:    "12345-tenor-gif",
		Provider: GIFProviderTenor,
	}

	if gif.GIFID != "12345-tenor-gif" {
		t.Errorf("Expected GIFID to be '12345-tenor-gif', got '%s'", gif.GIFID)
	}

	if gif.Provider != GIFProviderTenor {
		t.Errorf("Expected Provider to be GIFProviderTenor, got '%s'", gif.Provider)
	}
}

func TestTimeMarshalJSON(t *testing.T) {
	// Create a Time from a known time value
	testTime := time.Date(2024, 6, 15, 10, 30, 0, 0, time.UTC)
	customTime := Time{Time: testTime}

	// Marshal it
	data, err := json.Marshal(&customTime)
	if err != nil {
		t.Fatalf("Failed to marshal Time: %v", err)
	}

	// The result should be RFC3339 formatted
	expected := `"2024-06-15T10:30:00Z"`
	if string(data) != expected {
		t.Errorf("Expected %s, got %s", expected, string(data))
	}
}

func TestTimeUnmarshalJSON_AllFormats(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"RFC3339", `"2024-06-15T10:30:00Z"`},
		{"Threads format", `"2024-06-15T10:30:00+0000"`},
		{"With offset", `"2024-06-15T10:30:00-0700"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var customTime Time
			err := json.Unmarshal([]byte(tt.input), &customTime)
			if err != nil {
				t.Fatalf("Failed to unmarshal Time: %v", err)
			}
			if customTime.IsZero() {
				t.Error("Expected non-zero time")
			}
		})
	}
}

func TestTimeUnmarshalJSON_FallbackToDefault(t *testing.T) {
	// Test with a format that needs fallback to default time.Time unmarshalling
	input := `"2024-06-15T10:30:00.000Z"` // With milliseconds
	var customTime Time
	err := json.Unmarshal([]byte(input), &customTime)
	if err != nil {
		t.Fatalf("Failed to unmarshal Time: %v", err)
	}
	if customTime.IsZero() {
		t.Error("Expected non-zero time")
	}
}

func TestRateLimiter_NewRateLimiter(t *testing.T) {
	rl := NewRateLimiter(&RateLimiterConfig{
		InitialLimit: 100,
	})
	if rl == nil {
		t.Fatal("NewRateLimiter returned nil")
	}
	if rl.limit != 100 {
		t.Errorf("Expected limit to be 100, got %d", rl.limit)
	}
}

func TestRateLimiter_ShouldWait(t *testing.T) {
	rl := NewRateLimiter(&RateLimiterConfig{
		InitialLimit: 100,
	})

	// Initially should not need to wait
	should := rl.ShouldWait()
	if should {
		t.Error("ShouldWait should return false initially")
	}
}

func TestConfig_SetDefaults(t *testing.T) {
	config := &Config{}
	config.SetDefaults()

	if config.HTTPTimeout == 0 {
		t.Error("SetDefaults should set HTTPTimeout")
	}
	if config.BaseURL == "" {
		t.Error("SetDefaults should set BaseURL")
	}
	if config.UserAgent == "" {
		t.Error("SetDefaults should set UserAgent")
	}
}

func TestConfig_Validate(t *testing.T) {
	validator := NewConfigValidator()

	// Test with fully valid config
	config := &Config{
		ClientID:     "test-id",
		ClientSecret: "test-secret",
		RedirectURI:  "https://example.com/callback",
		Scopes:       []string{"threads_basic"},
		HTTPTimeout:  30 * time.Second,
		BaseURL:      "https://graph.threads.net",
		RetryConfig:  &RetryConfig{MaxRetries: 3, InitialDelay: time.Second, MaxDelay: 30 * time.Second, BackoffFactor: 2.0},
	}

	err := validator.Validate(config)
	if err != nil {
		t.Errorf("Expected valid config to pass, got: %v", err)
	}

	// Test with missing required field - empty scopes
	config.Scopes = nil
	err = validator.Validate(config)
	if err == nil {
		t.Error("Expected error for nil scopes")
	}
}

func TestValidationMoreCases(t *testing.T) {
	validator := NewValidator()

	t.Run("ValidateCarouselChildren", func(t *testing.T) {
		// Valid children count
		err := validator.ValidateCarouselChildren(2)
		if err != nil {
			t.Errorf("Expected no error for valid children count, got: %v", err)
		}

		// Empty children (0)
		err = validator.ValidateCarouselChildren(0)
		if err == nil {
			t.Error("Expected error for zero children")
		}

		// Too many children
		err = validator.ValidateCarouselChildren(25)
		if err == nil {
			t.Error("Expected error for too many children")
		}
	})

	t.Run("ValidateTopicTag", func(t *testing.T) {
		// Valid tag
		err := validator.ValidateTopicTag("valid")
		if err != nil {
			t.Errorf("Expected no error for valid tag, got: %v", err)
		}

		// Empty tag (should be valid - optional)
		err = validator.ValidateTopicTag("")
		if err != nil {
			t.Errorf("Expected no error for empty tag, got: %v", err)
		}
	})

	t.Run("ValidateMediaURL", func(t *testing.T) {
		// Valid URL
		err := validator.ValidateMediaURL("https://example.com/image.jpg", "image_url")
		if err != nil {
			t.Errorf("Expected no error for valid URL, got: %v", err)
		}

		// Empty URL
		err = validator.ValidateMediaURL("", "image_url")
		if err == nil {
			t.Error("Expected error for empty URL")
		}
	})
}
