package cmd

import "testing"

func TestPostsCmd_Structure(t *testing.T) {
	f := newTestFactory(t)
	cmd := NewPostsCmd(f)

	if cmd.Use != "posts" {
		t.Errorf("expected Use=posts, got %s", cmd.Use)
	}

	if cmd.Short == "" {
		t.Error("expected Short description to be set")
	}
}

func TestPostsCmd_Subcommands(t *testing.T) {
	f := newTestFactory(t)
	cmd := NewPostsCmd(f)

	expectedSubs := map[string]bool{
		"create":     true,
		"get":        true,
		"list":       true,
		"delete":     true,
		"carousel":   true,
		"quote":      true,
		"repost":     true,
		"unrepost":   true,
		"ghost-list": true,
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

func TestPostsCreateCmd_Flags(t *testing.T) {
	f := newTestFactory(t)
	cmd := newPostsCreateCmd(f)

	flags := []struct {
		name      string
		shorthand string
	}{
		{"text", "t"},
		{"image", ""},
		{"video", ""},
		{"alt-text", ""},
		{"reply-to", ""},
		{"poll", ""},
		{"ghost", ""},
		{"topic", ""},
		{"location", ""},
		{"reply-control", ""},
		{"gif", ""},
	}

	for _, f := range flags {
		flag := cmd.Flag(f.name)
		if flag == nil {
			t.Errorf("missing flag: %s", f.name)
			continue
		}
		if f.shorthand != "" && flag.Shorthand != f.shorthand {
			t.Errorf("flag %s expected shorthand %q, got %q", f.name, f.shorthand, flag.Shorthand)
		}
	}
}

func TestPostsGetCmd_Structure(t *testing.T) {
	f := newTestFactory(t)
	cmd := newPostsGetCmd(f)

	if cmd.Use != "get [post-id]" {
		t.Errorf("expected Use='get [post-id]', got %s", cmd.Use)
	}

	if cmd.Args == nil {
		t.Error("expected Args validator for exactly 1 arg")
	}

	if cmd.RunE == nil {
		t.Error("expected RunE to be set")
	}
}

func TestPostsListCmd_Structure(t *testing.T) {
	f := newTestFactory(t)
	cmd := newPostsListCmd(f)

	if cmd.Use != "list" {
		t.Errorf("expected Use=list, got %s", cmd.Use)
	}

	if cmd.RunE == nil {
		t.Error("expected RunE to be set")
	}
}

func TestPostsDeleteCmd_Structure(t *testing.T) {
	f := newTestFactory(t)
	cmd := newPostsDeleteCmd(f)

	if cmd.Use != "delete [post-id]" {
		t.Errorf("expected Use='delete [post-id]', got %s", cmd.Use)
	}

	if cmd.Args == nil {
		t.Error("expected Args validator for exactly 1 arg")
	}

	if cmd.RunE == nil {
		t.Error("expected RunE to be set")
	}
}

func TestPostsQuoteCmd_Structure(t *testing.T) {
	f := newTestFactory(t)
	cmd := newPostsQuoteCmd(f)

	if cmd.Use != "quote [post-id]" {
		t.Errorf("expected Use='quote [post-id]', got %s", cmd.Use)
	}

	if cmd.Args == nil {
		t.Error("expected Args validator")
	}

	if cmd.RunE == nil {
		t.Error("expected RunE to be set")
	}
}

func TestPostsQuoteCmd_Flags(t *testing.T) {
	f := newTestFactory(t)
	cmd := newPostsQuoteCmd(f)

	flags := []string{"text", "image", "video"}
	for _, flag := range flags {
		if cmd.Flag(flag) == nil {
			t.Errorf("missing flag: %s", flag)
		}
	}
}

func TestPostsQuoteCmd_HasExample(t *testing.T) {
	f := newTestFactory(t)
	cmd := newPostsQuoteCmd(f)

	if cmd.Example == "" {
		t.Error("expected Example to be set for quote command")
	}
}

func TestPostsRepostCmd_Structure(t *testing.T) {
	f := newTestFactory(t)
	cmd := newPostsRepostCmd(f)

	if cmd.Use != "repost [post-id]" {
		t.Errorf("expected Use='repost [post-id]', got %s", cmd.Use)
	}

	if cmd.Args == nil {
		t.Error("expected Args validator")
	}

	if cmd.RunE == nil {
		t.Error("expected RunE to be set")
	}
}

func TestPostsRepostCmd_HasExample(t *testing.T) {
	f := newTestFactory(t)
	cmd := newPostsRepostCmd(f)

	if cmd.Example == "" {
		t.Error("expected Example to be set for repost command")
	}
}

func TestPostsUnrepostCmd_Structure(t *testing.T) {
	f := newTestFactory(t)
	cmd := newPostsUnrepostCmd(f)

	if cmd.Use != "unrepost [repost-id]" {
		t.Errorf("expected Use='unrepost [repost-id]', got %s", cmd.Use)
	}

	if cmd.Args == nil {
		t.Error("expected Args validator")
	}

	if cmd.RunE == nil {
		t.Error("expected RunE to be set")
	}
}

func TestPostsUnrepostCmd_HasExample(t *testing.T) {
	f := newTestFactory(t)
	cmd := newPostsUnrepostCmd(f)

	if cmd.Example == "" {
		t.Error("expected Example to be set for unrepost command")
	}
}

func TestPostsUnrepostCmd_HasLongDescription(t *testing.T) {
	f := newTestFactory(t)
	cmd := newPostsUnrepostCmd(f)

	if cmd.Long == "" {
		t.Error("expected Long description to be set for unrepost command")
	}
}

func TestPostsCarouselCmd_Structure(t *testing.T) {
	f := newTestFactory(t)
	cmd := newPostsCarouselCmd(f)

	if cmd.Use != "carousel" {
		t.Errorf("expected Use=carousel, got %s", cmd.Use)
	}

	if cmd.RunE == nil {
		t.Error("expected RunE to be set")
	}
}

func TestPostsCarouselCmd_Flags(t *testing.T) {
	f := newTestFactory(t)
	cmd := newPostsCarouselCmd(f)

	flags := []string{"items", "text", "alt-text", "reply-to", "timeout"}
	for _, flag := range flags {
		if cmd.Flag(flag) == nil {
			t.Errorf("missing flag: %s", flag)
		}
	}

	itemsFlag := cmd.Flag("items")
	if itemsFlag == nil {
		t.Fatal("--items flag not found")
	}
}

func TestPostsCarouselCmd_TimeoutDefault(t *testing.T) {
	f := newTestFactory(t)
	cmd := newPostsCarouselCmd(f)

	timeoutFlag := cmd.Flag("timeout")
	if timeoutFlag == nil {
		t.Fatal("missing timeout flag")
	}

	if timeoutFlag.DefValue != "300" {
		t.Errorf("expected timeout default=300, got %s", timeoutFlag.DefValue)
	}
}

func TestPostsCarouselCmd_HasExample(t *testing.T) {
	f := newTestFactory(t)
	cmd := newPostsCarouselCmd(f)

	if cmd.Example == "" {
		t.Error("expected Example to be set for carousel command")
	}
}

func TestDetectMediaType_Image(t *testing.T) {
	tests := []struct {
		url      string
		expected string
	}{
		{"https://example.com/image.jpg", "IMAGE"},
		{"https://example.com/image.jpeg", "IMAGE"},
		{"https://example.com/image.png", "IMAGE"},
		{"https://example.com/image.gif", "IMAGE"},
		{"https://example.com/image.webp", "IMAGE"},
		{"https://example.com/image.JPG", "IMAGE"},
		{"https://example.com/image.PNG", "IMAGE"},
	}

	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			result := detectMediaType(tt.url)
			if result != tt.expected {
				t.Errorf("detectMediaType(%q) = %q, want %q", tt.url, result, tt.expected)
			}
		})
	}
}

func TestDetectMediaType_Video(t *testing.T) {
	tests := []struct {
		url      string
		expected string
	}{
		{"https://example.com/video.mp4", "VIDEO"},
		{"https://example.com/video.mov", "VIDEO"},
		{"https://example.com/video.m4v", "VIDEO"},
		{"https://example.com/video.webm", "VIDEO"},
		{"https://example.com/video.MP4", "VIDEO"},
		{"https://example.com/video.MOV", "VIDEO"},
	}

	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			result := detectMediaType(tt.url)
			if result != tt.expected {
				t.Errorf("detectMediaType(%q) = %q, want %q", tt.url, result, tt.expected)
			}
		})
	}
}

func TestDetectMediaType_WithQueryParams(t *testing.T) {
	tests := []struct {
		url      string
		expected string
	}{
		{"https://example.com/image.jpg?width=100", "IMAGE"},
		{"https://example.com/video.mp4?quality=hd", "VIDEO"},
		{"https://example.com/file.png?token=abc123", "IMAGE"},
	}

	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			result := detectMediaType(tt.url)
			if result != tt.expected {
				t.Errorf("detectMediaType(%q) = %q, want %q", tt.url, result, tt.expected)
			}
		})
	}
}

func TestDetectMediaType_DefaultToImage(t *testing.T) {
	tests := []string{
		"https://example.com/file",
		"https://example.com/file.txt",
		"https://example.com/file.pdf",
		"https://example.com/file.unknown",
	}

	for _, url := range tests {
		t.Run(url, func(t *testing.T) {
			result := detectMediaType(url)
			if result != "IMAGE" {
				t.Errorf("detectMediaType(%q) = %q, want IMAGE (default)", url, result)
			}
		})
	}
}

func TestPostsGhostListCmd_Structure(t *testing.T) {
	f := newTestFactory(t)
	cmd := newPostsGhostListCmd(f)

	if cmd.Use != "ghost-list" {
		t.Errorf("expected Use='ghost-list', got %s", cmd.Use)
	}

	if cmd.RunE == nil {
		t.Error("expected RunE to be set")
	}

	if cmd.Long == "" {
		t.Error("expected Long description to be set")
	}
}
