package cmd

import (
	"testing"
)

func TestPostsQuoteCmd_Structure(t *testing.T) {
	cmd := newPostsQuoteCmd()

	if cmd.Use != "quote [post-id]" {
		t.Errorf("expected Use='quote [post-id]', got %s", cmd.Use)
	}

	if cmd.Args == nil {
		t.Error("expected Args validator")
	}
}

func TestPostsQuoteCmd_Flags(t *testing.T) {
	cmd := newPostsQuoteCmd()

	flags := []string{"text", "image", "video"}
	for _, flag := range flags {
		if cmd.Flag(flag) == nil {
			t.Errorf("missing flag: %s", flag)
		}
	}
}

func TestPostsRepostCmd_Structure(t *testing.T) {
	cmd := newPostsRepostCmd()

	if cmd.Use != "repost [post-id]" {
		t.Errorf("expected Use='repost [post-id]', got %s", cmd.Use)
	}

	if cmd.Args == nil {
		t.Error("expected Args validator")
	}
}

func TestPostsCarouselCmd_Flags(t *testing.T) {
	// postsCarouselCmd is a package-level var
	cmd := postsCarouselCmd

	flags := []string{"items", "text", "alt-text", "reply-to", "timeout"}
	for _, flag := range flags {
		if cmd.Flag(flag) == nil {
			t.Errorf("missing flag: %s", flag)
		}
	}

	// --items should be required
	itemsFlag := cmd.Flag("items")
	if itemsFlag == nil {
		t.Fatal("--items flag not found")
	}
}
