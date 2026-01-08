package api

import (
	"testing"
)

// TestContainerBuilder_SetImageURL tests the SetImageURL method
func TestContainerBuilder_SetImageURL(t *testing.T) {
	tests := []struct {
		name     string
		imageURL string
		expected string
	}{
		{"valid URL", "https://example.com/image.jpg", "https://example.com/image.jpg"},
		{"empty URL", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewContainerBuilder()
			builder.SetImageURL(tt.imageURL)
			params := builder.Build()

			actual := params.Get("image_url")
			if actual != tt.expected {
				t.Errorf("expected image_url='%s', got '%s'", tt.expected, actual)
			}
		})
	}
}

// TestContainerBuilder_SetVideoURL tests the SetVideoURL method
func TestContainerBuilder_SetVideoURL(t *testing.T) {
	tests := []struct {
		name     string
		videoURL string
		expected string
	}{
		{"valid URL", "https://example.com/video.mp4", "https://example.com/video.mp4"},
		{"empty URL", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewContainerBuilder()
			builder.SetVideoURL(tt.videoURL)
			params := builder.Build()

			actual := params.Get("video_url")
			if actual != tt.expected {
				t.Errorf("expected video_url='%s', got '%s'", tt.expected, actual)
			}
		})
	}
}

// TestContainerBuilder_SetAltText tests the SetAltText method
func TestContainerBuilder_SetAltText(t *testing.T) {
	tests := []struct {
		name     string
		altText  string
		expected string
	}{
		{"valid alt text", "Image description", "Image description"},
		{"empty alt text", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewContainerBuilder()
			builder.SetAltText(tt.altText)
			params := builder.Build()

			actual := params.Get("alt_text")
			if actual != tt.expected {
				t.Errorf("expected alt_text='%s', got '%s'", tt.expected, actual)
			}
		})
	}
}

// TestContainerBuilder_SetReplyTo tests the SetReplyTo method
func TestContainerBuilder_SetReplyTo(t *testing.T) {
	tests := []struct {
		name      string
		replyToID string
		expected  string
	}{
		{"valid reply ID", "post123", "post123"},
		{"empty reply ID", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewContainerBuilder()
			builder.SetReplyTo(tt.replyToID)
			params := builder.Build()

			actual := params.Get("reply_to_id")
			if actual != tt.expected {
				t.Errorf("expected reply_to_id='%s', got '%s'", tt.expected, actual)
			}
		})
	}
}

// TestContainerBuilder_SetTopicTag tests the SetTopicTag method
func TestContainerBuilder_SetTopicTag(t *testing.T) {
	tests := []struct {
		name     string
		tag      string
		expected string
	}{
		{"valid tag", "golang", "golang"},
		{"empty tag", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewContainerBuilder()
			builder.SetTopicTag(tt.tag)
			params := builder.Build()

			actual := params.Get("topic_tag")
			if actual != tt.expected {
				t.Errorf("expected topic_tag='%s', got '%s'", tt.expected, actual)
			}
		})
	}
}

// TestContainerBuilder_SetLocationID tests the SetLocationID method
func TestContainerBuilder_SetLocationID(t *testing.T) {
	tests := []struct {
		name       string
		locationID string
		expected   string
	}{
		{"valid location ID", "loc123", "loc123"},
		{"empty location ID", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewContainerBuilder()
			builder.SetLocationID(tt.locationID)
			params := builder.Build()

			actual := params.Get("location_id")
			if actual != tt.expected {
				t.Errorf("expected location_id='%s', got '%s'", tt.expected, actual)
			}
		})
	}
}

// TestContainerBuilder_SetQuotePostID tests the SetQuotePostID method
func TestContainerBuilder_SetQuotePostID(t *testing.T) {
	tests := []struct {
		name        string
		quotePostID string
		expected    string
	}{
		{"valid quote post ID", "quote123", "quote123"},
		{"empty quote post ID", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewContainerBuilder()
			builder.SetQuotePostID(tt.quotePostID)
			params := builder.Build()

			actual := params.Get("quote_post_id")
			if actual != tt.expected {
				t.Errorf("expected quote_post_id='%s', got '%s'", tt.expected, actual)
			}
		})
	}
}

// TestContainerBuilder_SetLinkAttachment tests the SetLinkAttachment method
func TestContainerBuilder_SetLinkAttachment(t *testing.T) {
	tests := []struct {
		name     string
		linkURL  string
		expected string
	}{
		{"valid link", "https://example.com", "https://example.com"},
		{"empty link", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewContainerBuilder()
			builder.SetLinkAttachment(tt.linkURL)
			params := builder.Build()

			actual := params.Get("link_attachment")
			if actual != tt.expected {
				t.Errorf("expected link_attachment='%s', got '%s'", tt.expected, actual)
			}
		})
	}
}

// TestContainerBuilder_SetPollAttachment tests the SetPollAttachment method
func TestContainerBuilder_SetPollAttachment(t *testing.T) {
	t.Run("with poll", func(t *testing.T) {
		builder := NewContainerBuilder()
		poll := &PollAttachment{
			OptionA: "Option 1",
			OptionB: "Option 2",
		}
		builder.SetPollAttachment(poll)
		params := builder.Build()

		actual := params.Get("poll_attachment")
		if actual == "" {
			t.Error("expected poll_attachment to be set")
		}
	})

	t.Run("nil poll", func(t *testing.T) {
		builder := NewContainerBuilder()
		builder.SetPollAttachment(nil)
		params := builder.Build()

		actual := params.Get("poll_attachment")
		if actual != "" {
			t.Error("expected poll_attachment to be empty for nil poll")
		}
	})
}

// TestContainerBuilder_SetAllowlistedCountryCodes tests the SetAllowlistedCountryCodes method
func TestContainerBuilder_SetAllowlistedCountryCodes(t *testing.T) {
	builder := NewContainerBuilder()
	codes := []string{"US", "CA", "GB"}
	builder.SetAllowlistedCountryCodes(codes)
	params := builder.Build()

	actual := params["allowlisted_country_codes"]
	if len(actual) != 3 {
		t.Errorf("expected 3 country codes, got %d", len(actual))
	}
}

// TestContainerBuilder_AddChild tests the AddChild method
func TestContainerBuilder_AddChild(t *testing.T) {
	builder := NewContainerBuilder()
	builder.AddChild("child1")
	builder.AddChild("child2")
	params := builder.Build()

	children := params["children"]
	if len(children) != 2 {
		t.Errorf("expected 2 children, got %d", len(children))
	}
	if children[0] != "child1" {
		t.Errorf("expected first child 'child1', got '%s'", children[0])
	}
	if children[1] != "child2" {
		t.Errorf("expected second child 'child2', got '%s'", children[1])
	}
}

// TestContainerBuilder_SetChildren tests the SetChildren method
func TestContainerBuilder_SetChildren(t *testing.T) {
	builder := NewContainerBuilder()
	childIDs := []string{"child1", "child2", "child3"}
	builder.SetChildren(childIDs)
	params := builder.Build()

	children := params["children"]
	if len(children) != 3 {
		t.Errorf("expected 3 children in 'children', got %d", len(children))
	}

	// Check indexed parameters
	if params.Get("children[0]") != "child1" {
		t.Errorf("expected children[0]='child1', got '%s'", params.Get("children[0]"))
	}
	if params.Get("children[1]") != "child2" {
		t.Errorf("expected children[1]='child2', got '%s'", params.Get("children[1]"))
	}
	if params.Get("children[2]") != "child3" {
		t.Errorf("expected children[2]='child3', got '%s'", params.Get("children[2]"))
	}
}

// TestContainerBuilder_SetAutoPublishText tests the SetAutoPublishText method
func TestContainerBuilder_SetAutoPublishText(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		builder := NewContainerBuilder()
		builder.SetAutoPublishText(true)
		params := builder.Build()

		if params.Get("auto_publish_text") != "true" {
			t.Error("expected auto_publish_text to be 'true'")
		}
	})

	t.Run("false", func(t *testing.T) {
		builder := NewContainerBuilder()
		builder.SetAutoPublishText(false)
		params := builder.Build()

		if params.Get("auto_publish_text") != "" {
			t.Error("expected auto_publish_text to be empty for false")
		}
	})
}

// TestContainerBuilder_SetIsCarouselItem tests the SetIsCarouselItem method
func TestContainerBuilder_SetIsCarouselItem(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		builder := NewContainerBuilder()
		builder.SetIsCarouselItem(true)
		params := builder.Build()

		if params.Get("is_carousel_item") != "true" {
			t.Error("expected is_carousel_item to be 'true'")
		}
	})

	t.Run("false", func(t *testing.T) {
		builder := NewContainerBuilder()
		builder.SetIsCarouselItem(false)
		params := builder.Build()

		if params.Get("is_carousel_item") != "" {
			t.Error("expected is_carousel_item to be empty for false")
		}
	})
}

// TestContainerBuilder_SetTextEntities tests the SetTextEntities method
func TestContainerBuilder_SetTextEntities(t *testing.T) {
	t.Run("with entities", func(t *testing.T) {
		builder := NewContainerBuilder()
		entities := []TextEntity{
			{Offset: 0, Length: 10},
			{Offset: 15, Length: 5},
		}
		builder.SetTextEntities(entities)
		params := builder.Build()

		actual := params.Get("text_entities")
		if actual == "" {
			t.Error("expected text_entities to be set")
		}
	})

	t.Run("empty entities", func(t *testing.T) {
		builder := NewContainerBuilder()
		builder.SetTextEntities([]TextEntity{})
		params := builder.Build()

		actual := params.Get("text_entities")
		if actual != "" {
			t.Error("expected text_entities to be empty for empty slice")
		}
	})
}

// TestContainerBuilder_SetIsSpoilerMedia tests the SetIsSpoilerMedia method
func TestContainerBuilder_SetIsSpoilerMedia(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		builder := NewContainerBuilder()
		builder.SetIsSpoilerMedia(true)
		params := builder.Build()

		if params.Get("is_spoiler_media") != "true" {
			t.Error("expected is_spoiler_media to be 'true'")
		}
	})

	t.Run("false", func(t *testing.T) {
		builder := NewContainerBuilder()
		builder.SetIsSpoilerMedia(false)
		params := builder.Build()

		if params.Get("is_spoiler_media") != "" {
			t.Error("expected is_spoiler_media to be empty for false")
		}
	})
}

// TestContainerBuilder_SetTextAttachment tests the SetTextAttachment method
func TestContainerBuilder_SetTextAttachment(t *testing.T) {
	t.Run("with attachment", func(t *testing.T) {
		builder := NewContainerBuilder()
		attachment := &TextAttachment{
			Plaintext: "Test plaintext content",
		}
		builder.SetTextAttachment(attachment)
		params := builder.Build()

		actual := params.Get("text_attachment")
		if actual == "" {
			t.Error("expected text_attachment to be set")
		}
	})

	t.Run("nil attachment", func(t *testing.T) {
		builder := NewContainerBuilder()
		builder.SetTextAttachment(nil)
		params := builder.Build()

		actual := params.Get("text_attachment")
		if actual != "" {
			t.Error("expected text_attachment to be empty for nil")
		}
	})
}

// TestContainerBuilder_SetIsGhostPost tests the SetIsGhostPost method
func TestContainerBuilder_SetIsGhostPost(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		builder := NewContainerBuilder()
		builder.SetIsGhostPost(true)
		params := builder.Build()

		if params.Get("is_ghost_post") != "true" {
			t.Error("expected is_ghost_post to be 'true'")
		}
	})

	t.Run("false", func(t *testing.T) {
		builder := NewContainerBuilder()
		builder.SetIsGhostPost(false)
		params := builder.Build()

		if params.Get("is_ghost_post") != "" {
			t.Error("expected is_ghost_post to be empty for false")
		}
	})
}

// TestContainerBuilder_toString tests the toString helper method
func TestContainerBuilder_toString(t *testing.T) {
	builder := NewContainerBuilder()

	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{"string input", "hello", "hello"},
		{"int zero", 0, "0"},
		{"positive int", 123, "123"},
		{"negative int", -456, "-456"},
		{"large int", 1000000, "1000000"},
		{"other type", 3.14, ""},
		{"nil", nil, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := builder.toString(tt.input)
			if actual != tt.expected {
				t.Errorf("expected '%s', got '%s'", tt.expected, actual)
			}
		})
	}
}

// TestContainerBuilder_Chaining tests that methods can be chained
func TestContainerBuilder_Chaining(t *testing.T) {
	params := NewContainerBuilder().
		SetMediaType(MediaTypeImage).
		SetImageURL("https://example.com/img.jpg").
		SetAltText("Test image").
		SetReplyControl(ReplyControlEveryone).
		SetTopicTag("test").
		SetLocationID("loc123").
		Build()

	if params.Get("media_type") != MediaTypeImage {
		t.Errorf("expected media_type='%s', got '%s'", MediaTypeImage, params.Get("media_type"))
	}
	if params.Get("image_url") != "https://example.com/img.jpg" {
		t.Error("expected image_url to be set")
	}
	if params.Get("alt_text") != "Test image" {
		t.Error("expected alt_text to be set")
	}
	if params.Get("reply_control") != string(ReplyControlEveryone) {
		t.Error("expected reply_control to be set")
	}
	if params.Get("topic_tag") != "test" {
		t.Error("expected topic_tag to be set")
	}
	if params.Get("location_id") != "loc123" {
		t.Error("expected location_id to be set")
	}
}
