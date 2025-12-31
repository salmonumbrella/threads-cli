# Threads CLI v2 Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Transform the threads-cli scaffold into a production-quality CLI matching airwallex-cli patterns.

**Architecture:** Context-based dependency injection for testability. All configuration flows through context (IO, UI, format, flags). Generic helpers for list commands with pagination. Modular internal packages.

**Tech Stack:** Go 1.24, Cobra, 99designs/keyring, muesli/termenv, itchyny/gojq, golang.org/x/term

---

## Phase 1: Foundation Infrastructure

### Task 1: Add iocontext Package

**Files:**
- Create: `internal/iocontext/io.go`
- Test: `internal/iocontext/io_test.go`

**Step 1: Write the test**

```go
// internal/iocontext/io_test.go
package iocontext

import (
	"bytes"
	"context"
	"os"
	"strings"
	"testing"
)

func TestDefaultIO(t *testing.T) {
	io := DefaultIO()
	if io.Out != os.Stdout {
		t.Error("expected stdout")
	}
	if io.ErrOut != os.Stderr {
		t.Error("expected stderr")
	}
	if io.In != os.Stdin {
		t.Error("expected stdin")
	}
}

func TestWithIO(t *testing.T) {
	var buf bytes.Buffer
	io := &IO{Out: &buf, ErrOut: &buf, In: strings.NewReader("test")}
	ctx := WithIO(context.Background(), io)

	got := GetIO(ctx)
	if got != io {
		t.Error("expected injected IO")
	}
}

func TestGetIO_FallsBackToDefault(t *testing.T) {
	ctx := context.Background()
	io := GetIO(ctx)
	if io.Out != os.Stdout {
		t.Error("expected fallback to stdout")
	}
}

func TestHasIO(t *testing.T) {
	ctx := context.Background()
	if HasIO(ctx) {
		t.Error("expected no IO in empty context")
	}

	ctx = WithIO(ctx, &IO{})
	if !HasIO(ctx) {
		t.Error("expected IO after injection")
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/iocontext/... -v`
Expected: FAIL - package does not exist

**Step 3: Write implementation**

```go
// internal/iocontext/io.go
package iocontext

import (
	"context"
	"io"
	"os"
)

// IO holds input/output streams for commands
type IO struct {
	Out    io.Writer // stdout
	ErrOut io.Writer // stderr
	In     io.Reader // stdin
}

type contextKey struct{}

// DefaultIO returns IO using os.Std streams
func DefaultIO() *IO {
	return &IO{
		Out:    os.Stdout,
		ErrOut: os.Stderr,
		In:     os.Stdin,
	}
}

// WithIO injects IO into context
func WithIO(ctx context.Context, io *IO) context.Context {
	return context.WithValue(ctx, contextKey{}, io)
}

// GetIO retrieves IO from context, falling back to defaults
func GetIO(ctx context.Context) *IO {
	if io, ok := ctx.Value(contextKey{}).(*IO); ok && io != nil {
		return io
	}
	return DefaultIO()
}

// HasIO checks if IO is in context
func HasIO(ctx context.Context) bool {
	_, ok := ctx.Value(contextKey{}).(*IO)
	return ok
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./internal/iocontext/... -v`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/iocontext/io.go internal/iocontext/io_test.go
git commit -m "feat(iocontext): add IO context injection for testability"
```

---

### Task 2: Enhance outfmt Package - Column Types

**Files:**
- Modify: `internal/outfmt/outfmt.go`
- Create: `internal/outfmt/outfmt_test.go`

**Step 1: Write the test**

```go
// internal/outfmt/outfmt_test.go
package outfmt

import (
	"bytes"
	"context"
	"testing"
)

func TestColumnTypes(t *testing.T) {
	tests := []struct {
		name     string
		colType  ColumnType
		value    string
		wantDiff bool // whether output differs from input (colorized)
	}{
		{"plain stays same", ColumnPlain, "test", false},
		{"status gets color", ColumnStatus, "COMPLETED", true},
		{"amount gets color", ColumnAmount, "100.00", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Color disabled - should never change
			result := formatColumn(tt.value, tt.colType, false)
			if result != tt.value {
				t.Errorf("with color disabled, got %q want %q", result, tt.value)
			}
		})
	}
}

func TestFormatterOutput_JSON(t *testing.T) {
	var buf bytes.Buffer
	ctx := WithFormat(context.Background(), "json")
	f := FromContext(ctx, WithWriter(&buf))

	data := map[string]string{"key": "value"}
	if err := f.Output(data); err != nil {
		t.Fatal(err)
	}

	if buf.String() != `{"key":"value"}`+"\n" {
		t.Errorf("got %q", buf.String())
	}
}

func TestFormatterOutput_Text(t *testing.T) {
	var buf bytes.Buffer
	ctx := WithFormat(context.Background(), "text")
	f := FromContext(ctx, WithWriter(&buf))

	data := map[string]string{"key": "value"}
	if err := f.Output(data); err != nil {
		t.Fatal(err)
	}

	// Text mode should pretty-print
	if !bytes.Contains(buf.Bytes(), []byte("key")) {
		t.Errorf("expected key in output, got %q", buf.String())
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/outfmt/... -v`
Expected: FAIL - ColumnType not defined

**Step 3: Enhance implementation**

Add to `internal/outfmt/outfmt.go`:

```go
// ColumnType defines how a column should be formatted
type ColumnType int

const (
	ColumnPlain ColumnType = iota
	ColumnStatus
	ColumnAmount
	ColumnCurrency
	ColumnDate
	ColumnID
)

// formatColumn applies formatting based on column type
func formatColumn(value string, colType ColumnType, colorEnabled bool) string {
	if !colorEnabled {
		return value
	}

	switch colType {
	case ColumnStatus:
		return formatStatus(value)
	case ColumnAmount:
		return formatAmount(value)
	case ColumnCurrency:
		return formatCurrency(value)
	case ColumnDate:
		return formatDate(value)
	case ColumnID:
		return formatID(value)
	default:
		return value
	}
}

func formatStatus(status string) string {
	switch status {
	case "PUBLISHED", "FINISHED", "ACTIVE":
		return "\033[32m" + status + "\033[0m" // Green
	case "IN_PROGRESS", "PUBLISHING":
		return "\033[33m" + status + "\033[0m" // Yellow
	case "FAILED", "ERROR":
		return "\033[31m" + status + "\033[0m" // Red
	default:
		return status
	}
}

func formatAmount(amount string) string {
	if len(amount) > 0 && amount[0] == '-' {
		return "\033[31m" + amount + "\033[0m" // Red for negative
	}
	return "\033[32m" + amount + "\033[0m" // Green for positive
}

func formatCurrency(currency string) string {
	return "\033[36m" + currency + "\033[0m" // Cyan
}

func formatDate(date string) string {
	return "\033[90m" + date + "\033[0m" // Gray
}

func formatID(id string) string {
	return "\033[34m" + id + "\033[0m" // Blue
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./internal/outfmt/... -v`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/outfmt/outfmt.go internal/outfmt/outfmt_test.go
git commit -m "feat(outfmt): add column types for colorized output"
```

---

### Task 3: Add Formatter Table Output

**Files:**
- Modify: `internal/outfmt/outfmt.go`
- Modify: `internal/outfmt/outfmt_test.go`

**Step 1: Write the test**

Add to `internal/outfmt/outfmt_test.go`:

```go
func TestFormatter_Table(t *testing.T) {
	var buf bytes.Buffer
	ctx := WithFormat(context.Background(), "text")
	f := FromContext(ctx, WithWriter(&buf))

	headers := []string{"ID", "STATUS", "COUNT"}
	rows := [][]string{
		{"123", "ACTIVE", "10"},
		{"456", "PENDING", "20"},
	}

	if err := f.Table(headers, rows, nil); err != nil {
		t.Fatal(err)
	}

	output := buf.String()
	if !strings.Contains(output, "ID") {
		t.Error("missing header")
	}
	if !strings.Contains(output, "123") {
		t.Error("missing row data")
	}
}

func TestFormatter_TableWithColors(t *testing.T) {
	var buf bytes.Buffer
	ctx := WithFormat(context.Background(), "text")
	f := FromContext(ctx, WithWriter(&buf))

	headers := []string{"ID", "STATUS"}
	rows := [][]string{{"123", "PUBLISHED"}}
	colTypes := []ColumnType{ColumnID, ColumnStatus}

	// With color disabled, should not contain ANSI codes
	if err := f.Table(headers, rows, colTypes); err != nil {
		t.Fatal(err)
	}
	// Just verify it doesn't crash - color depends on terminal detection
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/outfmt/... -v`
Expected: FAIL - Table method not defined

**Step 3: Add Table method**

Add to `internal/outfmt/outfmt.go`:

```go
import (
	"text/tabwriter"
)

// Table outputs data in tabular format
func (f *Formatter) Table(headers []string, rows [][]string, colTypes []ColumnType) error {
	if IsJSON(f.ctx) {
		// Convert to JSON array of objects
		items := make([]map[string]string, len(rows))
		for i, row := range rows {
			item := make(map[string]string)
			for j, val := range row {
				if j < len(headers) {
					item[headers[j]] = val
				}
			}
			items[i] = item
		}
		return f.Output(items)
	}

	tw := tabwriter.NewWriter(f.out, 0, 4, 2, ' ', 0)

	// Write headers
	for i, h := range headers {
		if i > 0 {
			fmt.Fprint(tw, "\t")
		}
		fmt.Fprint(tw, h)
	}
	fmt.Fprintln(tw)

	// Write rows
	for _, row := range rows {
		for i, val := range row {
			if i > 0 {
				fmt.Fprint(tw, "\t")
			}
			// Apply column type formatting if provided
			if colTypes != nil && i < len(colTypes) {
				val = formatColumn(val, colTypes[i], f.colorEnabled())
			}
			fmt.Fprint(tw, val)
		}
		fmt.Fprintln(tw)
	}

	return tw.Flush()
}

func (f *Formatter) colorEnabled() bool {
	// Check NO_COLOR env and terminal detection
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	// For now, disable in non-TTY contexts
	return term.IsTerminal(int(os.Stdout.Fd()))
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./internal/outfmt/... -v`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/outfmt/outfmt.go internal/outfmt/outfmt_test.go
git commit -m "feat(outfmt): add Table method for tabular output"
```

---

### Task 4: Add helpers Package

**Files:**
- Create: `internal/helpers/helpers.go`
- Create: `internal/helpers/helpers_test.go`

**Step 1: Write the test**

```go
// internal/helpers/helpers_test.go
package helpers

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/salmonumbrella/threads-go/internal/iocontext"
	"github.com/salmonumbrella/threads-go/internal/outfmt"
)

func TestConfirmOrYes_WithYesFlag(t *testing.T) {
	ctx := outfmt.WithYes(context.Background(), true)
	confirmed, err := ConfirmOrYes(ctx, "Delete?")
	if err != nil {
		t.Fatal(err)
	}
	if !confirmed {
		t.Error("expected confirmation with --yes flag")
	}
}

func TestConfirmOrYes_JSONMode(t *testing.T) {
	ctx := outfmt.WithFormat(context.Background(), "json")
	confirmed, err := ConfirmOrYes(ctx, "Delete?")
	if err != nil {
		t.Fatal(err)
	}
	if !confirmed {
		t.Error("expected confirmation in JSON mode")
	}
}

func TestConfirmOrYes_UserSaysYes(t *testing.T) {
	var outBuf bytes.Buffer
	ctx := context.Background()
	ctx = outfmt.WithFormat(ctx, "text")
	ctx = outfmt.WithYes(ctx, false)
	ctx = iocontext.WithIO(ctx, &iocontext.IO{
		Out:    &outBuf,
		ErrOut: &outBuf,
		In:     strings.NewReader("yes\n"),
	})

	// Mock terminal check
	oldIsTerminal := isTerminal
	isTerminal = func() bool { return true }
	defer func() { isTerminal = oldIsTerminal }()

	confirmed, err := ConfirmOrYes(ctx, "Delete?")
	if err != nil {
		t.Fatal(err)
	}
	if !confirmed {
		t.Error("expected confirmation when user says yes")
	}
}

func TestConfirmOrYes_UserSaysNo(t *testing.T) {
	var outBuf bytes.Buffer
	ctx := context.Background()
	ctx = outfmt.WithFormat(ctx, "text")
	ctx = outfmt.WithYes(ctx, false)
	ctx = iocontext.WithIO(ctx, &iocontext.IO{
		Out:    &outBuf,
		ErrOut: &outBuf,
		In:     strings.NewReader("no\n"),
	})

	oldIsTerminal := isTerminal
	isTerminal = func() bool { return true }
	defer func() { isTerminal = oldIsTerminal }()

	confirmed, err := ConfirmOrYes(ctx, "Delete?")
	if err != nil {
		t.Fatal(err)
	}
	if confirmed {
		t.Error("expected rejection when user says no")
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/helpers/... -v`
Expected: FAIL - package does not exist

**Step 3: Write implementation**

```go
// internal/helpers/helpers.go
package helpers

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"

	"github.com/salmonumbrella/threads-go/internal/iocontext"
	"github.com/salmonumbrella/threads-go/internal/outfmt"
)

// isTerminal is a variable that can be overridden in tests
var isTerminal = func() bool {
	return term.IsTerminal(int(os.Stdin.Fd()))
}

// ConfirmOrYes prompts for confirmation unless --yes flag is set
func ConfirmOrYes(ctx context.Context, prompt string) (bool, error) {
	// Skip if --yes flag set
	if outfmt.GetYes(ctx) {
		return true, nil
	}

	// Skip in JSON mode (scripts expect non-interactive)
	if outfmt.IsJSON(ctx) {
		return true, nil
	}

	// Check if stdin is a terminal
	if !isTerminal() {
		return false, fmt.Errorf("cannot prompt for confirmation: stdin is not a terminal (use --yes to skip)")
	}

	io := iocontext.GetIO(ctx)

	// Prompt to stderr
	fmt.Fprint(io.ErrOut, prompt+" [y/N]: ")

	// Read response
	reader := bufio.NewReader(io.In)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false, fmt.Errorf("failed to read confirmation: %w", err)
	}

	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y" || response == "yes", nil
}

// MustMarkRequired marks a flag as required, panicking on error
func MustMarkRequired(cmd interface{ MarkFlagRequired(string) error }, name string) {
	if err := cmd.MarkFlagRequired(name); err != nil {
		panic(fmt.Sprintf("flag %q not defined: %v", name, err))
	}
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./internal/helpers/... -v`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/helpers/helpers.go internal/helpers/helpers_test.go
git commit -m "feat(helpers): add confirmation and utility helpers"
```

---

### Task 5: Add list_helper Generic Pagination

**Files:**
- Create: `internal/cmd/list_helper.go`
- Create: `internal/cmd/list_helper_test.go`

**Step 1: Write the test**

```go
// internal/cmd/list_helper_test.go
package cmd

import (
	"bytes"
	"context"
	"testing"

	threads "github.com/salmonumbrella/threads-go"
	"github.com/salmonumbrella/threads-go/internal/iocontext"
	"github.com/salmonumbrella/threads-go/internal/outfmt"
)

type mockPost struct {
	ID     string
	Text   string
	Status string
}

func TestNewListCommand(t *testing.T) {
	cfg := ListConfig[mockPost]{
		Use:     "list",
		Short:   "List items",
		Headers: []string{"ID", "TEXT", "STATUS"},
		RowFunc: func(p mockPost) []string {
			return []string{p.ID, p.Text, p.Status}
		},
		Fetch: func(ctx context.Context, client *threads.Client, cursor string, limit int) (ListResult[mockPost], error) {
			return ListResult[mockPost]{
				Items:   []mockPost{{ID: "1", Text: "Hello", Status: "PUBLISHED"}},
				HasMore: false,
			}, nil
		},
		EmptyMessage: "No posts found",
	}

	// Mock client getter
	getClient := func(ctx context.Context) (*threads.Client, error) {
		return nil, nil // Tests don't need real client
	}

	cmd := NewListCommand(cfg, getClient)
	if cmd.Use != "list" {
		t.Errorf("expected Use=list, got %s", cmd.Use)
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/cmd/... -run TestNewListCommand -v`
Expected: FAIL - ListConfig not defined

**Step 3: Write implementation**

```go
// internal/cmd/list_helper.go
package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	threads "github.com/salmonumbrella/threads-go"
	"github.com/salmonumbrella/threads-go/internal/outfmt"
)

// ListResult represents the result of a paginated list operation
type ListResult[T any] struct {
	Items   []T
	HasMore bool
	Cursor  string
}

// ListConfig defines how a list command behaves
type ListConfig[T any] struct {
	Use          string
	Short        string
	Long         string
	Example      string
	Headers      []string
	RowFunc      func(T) []string
	ColumnTypes  []outfmt.ColumnType
	EmptyMessage string

	// Fetch function - called with cursor and limit, returns items
	Fetch func(ctx context.Context, client *threads.Client, cursor string, limit int) (ListResult[T], error)
}

// NewListCommand creates a cobra command from ListConfig
func NewListCommand[T any](cfg ListConfig[T], getClient func(context.Context) (*threads.Client, error)) *cobra.Command {
	var limit int
	var cursor string

	cmd := &cobra.Command{
		Use:     cfg.Use,
		Short:   cfg.Short,
		Long:    cfg.Long,
		Example: cfg.Example,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Default limit
			if limit <= 0 {
				limit = 25
			}
			if limit > 100 {
				limit = 100
			}

			client, err := getClient(cmd.Context())
			if err != nil {
				return err
			}

			result, err := cfg.Fetch(cmd.Context(), client, cursor, limit)
			if err != nil {
				return err
			}

			f := outfmt.FromContext(cmd.Context())

			// Handle empty results
			if len(result.Items) == 0 {
				if outfmt.IsJSON(cmd.Context()) {
					return f.Output(map[string]interface{}{
						"items":    result.Items,
						"has_more": result.HasMore,
					})
				}
				f.Empty(cfg.EmptyMessage)
				return nil
			}

			// Build rows from items
			rows := make([][]string, len(result.Items))
			for i, item := range result.Items {
				rows[i] = cfg.RowFunc(item)
			}

			if err := f.Table(cfg.Headers, rows, cfg.ColumnTypes); err != nil {
				return err
			}

			// Show pagination hint
			if !outfmt.IsJSON(cmd.Context()) && result.HasMore {
				fmt.Fprintln(os.Stderr, "# More results available. Use --cursor to paginate.")
			}

			return nil
		},
	}

	cmd.Flags().IntVar(&limit, "limit", 25, "Maximum results (1-100)")
	cmd.Flags().StringVar(&cursor, "cursor", "", "Pagination cursor")

	return cmd
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./internal/cmd/... -run TestNewListCommand -v`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/cmd/list_helper.go internal/cmd/list_helper_test.go
git commit -m "feat(cmd): add generic list command helper with pagination"
```

---

### Task 6: Add Shell Completion Command

**Files:**
- Create: `internal/cmd/completion.go`
- Modify: `internal/cmd/root.go` (add completion command)

**Step 1: Write the implementation**

```go
// internal/cmd/completion.go
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

func newCompletionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate shell completion scripts",
		Long: `Generate shell completion scripts for threads.

To load completions:

Bash:
  $ source <(threads completion bash)
  # To load completions for each session, execute once:
  # Linux:
  $ threads completion bash > /etc/bash_completion.d/threads
  # macOS:
  $ threads completion bash > $(brew --prefix)/etc/bash_completion.d/threads

Zsh:
  # If shell completion is not already enabled in your environment,
  # you will need to enable it. You can execute the following once:
  $ echo "autoload -U compinit; compinit" >> ~/.zshrc
  # To load completions for each session, execute once:
  $ threads completion zsh > "${fpath[1]}/_threads"
  # You will need to start a new shell for this setup to take effect.

Fish:
  $ threads completion fish | source
  # To load completions for each session, execute once:
  $ threads completion fish > ~/.config/fish/completions/threads.fish

PowerShell:
  PS> threads completion powershell | Out-String | Invoke-Expression
  # To load completions for every new session, run:
  PS> threads completion powershell > threads.ps1
  # and source this file from your PowerShell profile.
`,
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		RunE: func(cmd *cobra.Command, args []string) error {
			switch args[0] {
			case "bash":
				return cmd.Root().GenBashCompletion(os.Stdout)
			case "zsh":
				return cmd.Root().GenZshCompletion(os.Stdout)
			case "fish":
				return cmd.Root().GenFishCompletion(os.Stdout, true)
			case "powershell":
				return cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
			}
			return nil
		},
	}
	return cmd
}
```

**Step 2: Add to root.go**

In `internal/cmd/root.go`, add to the init or NewRootCmd function:

```go
rootCmd.AddCommand(newCompletionCmd())
```

**Step 3: Verify it works**

Run: `go build -o threads ./cmd/threads && ./threads completion bash | head -20`
Expected: Bash completion script output

**Step 4: Commit**

```bash
git add internal/cmd/completion.go internal/cmd/root.go
git commit -m "feat(completion): add shell completion command"
```

---

## Phase 2: Missing API Commands

### Task 7: Add Carousel Post Support

**Files:**
- Modify: `internal/cmd/posts.go`
- Create: `internal/cmd/posts_test.go`

**Step 1: Write the test**

```go
// internal/cmd/posts_test.go
package cmd

import (
	"testing"
)

func TestPostsCarouselCmd_Flags(t *testing.T) {
	cmd := newPostsCarouselCmd()

	// Check required flags exist
	flags := []string{"items"}
	for _, flag := range flags {
		if cmd.Flag(flag) == nil {
			t.Errorf("missing flag: %s", flag)
		}
	}
}
```

**Step 2: Add carousel subcommand**

Add to `internal/cmd/posts.go`:

```go
func newPostsCarouselCmd() *cobra.Command {
	var items []string
	var text string
	var altTexts []string
	var replyTo string

	cmd := &cobra.Command{
		Use:   "carousel",
		Short: "Create a carousel post with multiple images/videos",
		Long: `Create a carousel post with 2-20 media items.

Each item should be a URL to an image or video. Alt text can be provided
for accessibility using --alt-text (one per item, in order).`,
		Example: `  # Create carousel with 3 images
  threads posts carousel --items url1,url2,url3

  # With caption and alt text
  threads posts carousel --items url1,url2 --text "My photos" --alt-text "First image" --alt-text "Second image"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(items) < 2 {
				return fmt.Errorf("carousel requires at least 2 items")
			}
			if len(items) > 20 {
				return fmt.Errorf("carousel supports maximum 20 items")
			}

			client, err := getClient(cmd.Context())
			if err != nil {
				return err
			}

			// Build carousel content
			content := &threads.CarouselPostContent{
				Text: text,
			}

			// Create media containers for each item
			for i, itemURL := range items {
				var altText string
				if i < len(altTexts) {
					altText = altTexts[i]
				}

				// Detect media type from URL
				mediaType := detectMediaType(itemURL)
				containerID, err := client.CreateMediaContainer(cmd.Context(), mediaType, itemURL, altText)
				if err != nil {
					return fmt.Errorf("failed to create container for item %d: %w", i+1, err)
				}

				// Wait for container to be ready
				if err := waitForContainer(cmd.Context(), client, containerID); err != nil {
					return fmt.Errorf("container %d not ready: %w", i+1, err)
				}

				content.Children = append(content.Children, string(containerID))
			}

			if replyTo != "" {
				content.ReplyToID = threads.PostID(replyTo)
			}

			post, err := client.CreateCarouselPost(cmd.Context(), content)
			if err != nil {
				return err
			}

			f := outfmt.FromContext(cmd.Context())
			return f.Output(post)
		},
	}

	cmd.Flags().StringSliceVar(&items, "items", nil, "Media URLs (comma-separated or multiple flags)")
	cmd.Flags().StringVar(&text, "text", "", "Caption text")
	cmd.Flags().StringSliceVar(&altTexts, "alt-text", nil, "Alt text for each item (in order)")
	cmd.Flags().StringVar(&replyTo, "reply-to", "", "Post ID to reply to")

	_ = cmd.MarkFlagRequired("items")

	return cmd
}

func detectMediaType(url string) string {
	lower := strings.ToLower(url)
	if strings.Contains(lower, ".mp4") || strings.Contains(lower, ".mov") {
		return "VIDEO"
	}
	return "IMAGE"
}

func waitForContainer(ctx context.Context, client *threads.Client, containerID threads.ContainerID) error {
	// Poll container status with timeout
	timeout := time.After(5 * time.Minute)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timeout:
			return fmt.Errorf("timeout waiting for container")
		case <-ticker.C:
			status, err := client.GetContainerStatus(ctx, containerID)
			if err != nil {
				return err
			}
			switch status.Status {
			case "FINISHED":
				return nil
			case "ERROR":
				return fmt.Errorf("container error: %s", status.Error)
			}
			// Still processing, continue waiting
		}
	}
}
```

**Step 3: Add to posts command**

In the `newPostsCmd()` function, add:
```go
cmd.AddCommand(newPostsCarouselCmd())
```

**Step 4: Run test and verify**

Run: `go test ./internal/cmd/... -run TestPostsCarousel -v`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/cmd/posts.go internal/cmd/posts_test.go
git commit -m "feat(posts): add carousel post creation"
```

---

### Task 8: Add Quote Post Support

**Files:**
- Modify: `internal/cmd/posts.go`

**Step 1: Add quote subcommand**

Add to `internal/cmd/posts.go`:

```go
func newPostsQuoteCmd() *cobra.Command {
	var text string
	var imageURL string
	var videoURL string

	cmd := &cobra.Command{
		Use:   "quote [post-id]",
		Short: "Create a quote post",
		Long:  "Quote an existing post with optional text, image, or video.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			quotedPostID := args[0]

			client, err := getClient(cmd.Context())
			if err != nil {
				return err
			}

			var content interface{}
			switch {
			case videoURL != "":
				content = &threads.VideoPostContent{
					VideoURL: videoURL,
					Text:     text,
				}
			case imageURL != "":
				content = &threads.ImagePostContent{
					ImageURL: imageURL,
					Text:     text,
				}
			default:
				content = &threads.TextPostContent{
					Text: text,
				}
			}

			post, err := client.CreateQuotePost(cmd.Context(), content, quotedPostID)
			if err != nil {
				return err
			}

			f := outfmt.FromContext(cmd.Context())
			return f.Output(post)
		},
	}

	cmd.Flags().StringVar(&text, "text", "", "Quote text")
	cmd.Flags().StringVar(&imageURL, "image", "", "Image URL to include")
	cmd.Flags().StringVar(&videoURL, "video", "", "Video URL to include")

	return cmd
}
```

**Step 2: Add repost subcommand**

```go
func newPostsRepostCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "repost [post-id]",
		Short: "Repost an existing post",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			postID := args[0]

			client, err := getClient(cmd.Context())
			if err != nil {
				return err
			}

			post, err := client.RepostPost(cmd.Context(), threads.PostID(postID))
			if err != nil {
				return err
			}

			f := outfmt.FromContext(cmd.Context())
			return f.Output(post)
		},
	}
	return cmd
}
```

**Step 3: Add to posts command**

```go
cmd.AddCommand(newPostsQuoteCmd())
cmd.AddCommand(newPostsRepostCmd())
```

**Step 4: Commit**

```bash
git add internal/cmd/posts.go
git commit -m "feat(posts): add quote and repost commands"
```

---

### Task 9: Add Location Commands

**Files:**
- Create: `internal/cmd/locations.go`
- Modify: `internal/cmd/root.go`

**Step 1: Write implementation**

```go
// internal/cmd/locations.go
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	threads "github.com/salmonumbrella/threads-go"
	"github.com/salmonumbrella/threads-go/internal/outfmt"
)

func newLocationsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "locations",
		Aliases: []string{"location", "loc"},
		Short:   "Location search and details",
	}

	cmd.AddCommand(newLocationsSearchCmd())
	cmd.AddCommand(newLocationsGetCmd())

	return cmd
}

func newLocationsSearchCmd() *cobra.Command {
	var lat, lng float64
	var limit int

	cmd := &cobra.Command{
		Use:   "search [query]",
		Short: "Search for locations",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var query string
			if len(args) > 0 {
				query = args[0]
			}

			if query == "" && lat == 0 && lng == 0 {
				return fmt.Errorf("provide either a search query or --lat/--lng coordinates")
			}

			client, err := getClient(cmd.Context())
			if err != nil {
				return err
			}

			var latPtr, lngPtr *float64
			if lat != 0 || lng != 0 {
				latPtr = &lat
				lngPtr = &lng
			}

			result, err := client.SearchLocations(cmd.Context(), query, latPtr, lngPtr)
			if err != nil {
				return err
			}

			f := outfmt.FromContext(cmd.Context())

			if outfmt.IsJSON(cmd.Context()) {
				return f.Output(result)
			}

			if len(result.Data) == 0 {
				f.Empty("No locations found")
				return nil
			}

			headers := []string{"ID", "NAME", "ADDRESS"}
			rows := make([][]string, len(result.Data))
			for i, loc := range result.Data {
				rows[i] = []string{
					string(loc.ID),
					loc.Name,
					loc.Address,
				}
			}

			return f.Table(headers, rows, nil)
		},
	}

	cmd.Flags().Float64Var(&lat, "lat", 0, "Latitude for coordinate search")
	cmd.Flags().Float64Var(&lng", "lng", 0, "Longitude for coordinate search")
	cmd.Flags().IntVar(&limit, "limit", 10, "Maximum results")

	return cmd
}

func newLocationsGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get [location-id]",
		Short: "Get location details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			locationID := args[0]

			client, err := getClient(cmd.Context())
			if err != nil {
				return err
			}

			location, err := client.GetLocation(cmd.Context(), threads.LocationID(locationID))
			if err != nil {
				return err
			}

			f := outfmt.FromContext(cmd.Context())
			return f.Output(location)
		},
	}
	return cmd
}
```

**Step 2: Add to root.go**

```go
rootCmd.AddCommand(newLocationsCmd())
```

**Step 3: Commit**

```bash
git add internal/cmd/locations.go internal/cmd/root.go
git commit -m "feat(locations): add location search and get commands"
```

---

### Task 10: Add Rate Limit Command

**Files:**
- Create: `internal/cmd/ratelimit.go`
- Modify: `internal/cmd/root.go`

**Step 1: Write implementation**

```go
// internal/cmd/ratelimit.go
package cmd

import (
	"github.com/spf13/cobra"

	"github.com/salmonumbrella/threads-go/internal/outfmt"
)

func newRateLimitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ratelimit",
		Aliases: []string{"rate", "limits"},
		Short:   "View rate limit status",
	}

	cmd.AddCommand(newRateLimitStatusCmd())
	cmd.AddCommand(newRateLimitPublishingCmd())

	return cmd
}

func newRateLimitStatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show current rate limit status",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := getClient(cmd.Context())
			if err != nil {
				return err
			}

			status := client.GetRateLimitStatus()

			f := outfmt.FromContext(cmd.Context())

			if outfmt.IsJSON(cmd.Context()) {
				return f.Output(map[string]interface{}{
					"is_limited":       status.IsLimited,
					"remaining":        status.Remaining,
					"limit":            status.Limit,
					"reset_at":         status.ResetAt,
					"near_limit":       client.IsNearRateLimit(0.8),
				})
			}

			// Text output
			if status.IsLimited {
				fmt.Fprintf(f.ErrOut(), "Rate limited until %s\n", status.ResetAt.Format(time.RFC3339))
			} else {
				fmt.Fprintf(f.Out(), "Remaining: %d/%d\n", status.Remaining, status.Limit)
				if client.IsNearRateLimit(0.8) {
					fmt.Fprintln(f.ErrOut(), "Warning: Near rate limit threshold")
				}
			}

			return nil
		},
	}
	return cmd
}

func newRateLimitPublishingCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "publishing",
		Short: "Show publishing limits (API quota)",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := getClient(cmd.Context())
			if err != nil {
				return err
			}

			limits, err := client.GetPublishingLimits(cmd.Context())
			if err != nil {
				return err
			}

			f := outfmt.FromContext(cmd.Context())
			return f.Output(limits)
		},
	}
	return cmd
}
```

**Step 2: Add to root.go**

```go
rootCmd.AddCommand(newRateLimitCmd())
```

**Step 3: Commit**

```bash
git add internal/cmd/ratelimit.go internal/cmd/root.go
git commit -m "feat(ratelimit): add rate limit and publishing limits commands"
```

---

### Task 11: Add User Mentions Command

**Files:**
- Modify: `internal/cmd/users.go`

**Step 1: Add mentions subcommand**

Add to `internal/cmd/users.go`:

```go
func newUsersMentionsCmd() *cobra.Command {
	var limit int
	var cursor string

	cmd := &cobra.Command{
		Use:   "mentions",
		Short: "List posts mentioning you",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := getClient(cmd.Context())
			if err != nil {
				return err
			}

			// Get authenticated user
			me, err := client.GetMe(cmd.Context())
			if err != nil {
				return err
			}

			opts := &threads.PaginationOptions{
				Limit:  limit,
				Cursor: cursor,
			}

			result, err := client.GetUserMentions(cmd.Context(), me.ID, opts)
			if err != nil {
				return err
			}

			f := outfmt.FromContext(cmd.Context())

			if outfmt.IsJSON(cmd.Context()) {
				return f.Output(result)
			}

			if len(result.Data) == 0 {
				f.Empty("No mentions found")
				return nil
			}

			headers := []string{"ID", "FROM", "TEXT", "TIMESTAMP"}
			rows := make([][]string, len(result.Data))
			for i, post := range result.Data {
				text := post.Text
				if len(text) > 50 {
					text = text[:47] + "..."
				}
				rows[i] = []string{
					string(post.ID),
					post.Username,
					text,
					post.Timestamp.Format("2006-01-02 15:04"),
				}
			}

			return f.Table(headers, rows, []outfmt.ColumnType{
				outfmt.ColumnID,
				outfmt.ColumnPlain,
				outfmt.ColumnPlain,
				outfmt.ColumnDate,
			})
		},
	}

	cmd.Flags().IntVar(&limit, "limit", 25, "Maximum results")
	cmd.Flags().StringVar(&cursor, "cursor", "", "Pagination cursor")

	return cmd
}
```

**Step 2: Add to users command**

```go
cmd.AddCommand(newUsersMentionsCmd())
```

**Step 3: Commit**

```bash
git add internal/cmd/users.go
git commit -m "feat(users): add mentions command"
```

---

### Task 12: Enhance Search with Advanced Options

**Files:**
- Modify: `internal/cmd/search.go`

**Step 1: Add advanced search flags**

Update the search command to include:

```go
func newSearchCmd() *cobra.Command {
	var limit int
	var cursor string
	var mediaType string
	var since string
	var until string

	cmd := &cobra.Command{
		Use:   "search [query]",
		Short: "Search for posts",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			query := args[0]

			client, err := getClient(cmd.Context())
			if err != nil {
				return err
			}

			opts := &threads.SearchOptions{
				Limit:  limit,
				Cursor: cursor,
			}

			if mediaType != "" {
				opts.MediaType = mediaType
			}

			if since != "" {
				t, err := time.Parse("2006-01-02", since)
				if err != nil {
					return fmt.Errorf("invalid --since date: %w", err)
				}
				opts.Since = &t
			}

			if until != "" {
				t, err := time.Parse("2006-01-02", until)
				if err != nil {
					return fmt.Errorf("invalid --until date: %w", err)
				}
				opts.Until = &t
			}

			result, err := client.KeywordSearch(cmd.Context(), query, opts)
			if err != nil {
				return err
			}

			// ... output handling
		},
	}

	cmd.Flags().IntVar(&limit, "limit", 25, "Maximum results")
	cmd.Flags().StringVar(&cursor, "cursor", "", "Pagination cursor")
	cmd.Flags().StringVar(&mediaType, "media-type", "", "Filter by media type (TEXT, IMAGE, VIDEO)")
	cmd.Flags().StringVar(&since, "since", "", "Posts after date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&until, "until", "", "Posts before date (YYYY-MM-DD)")

	return cmd
}
```

**Step 2: Commit**

```bash
git add internal/cmd/search.go
git commit -m "feat(search): add advanced search filters"
```

---

## Phase 3: Polish & Testing

### Task 13: Add Command Tests

**Files:**
- Create: `internal/cmd/auth_test.go`
- Create: `internal/cmd/posts_test.go` (if not exists)
- Create: `internal/cmd/users_test.go`

**Step 1: Write auth tests**

```go
// internal/cmd/auth_test.go
package cmd

import (
	"testing"
)

func TestAuthLoginCmd_HasRequiredFlags(t *testing.T) {
	cmd := newAuthLoginCmd()

	// Should have scopes flag
	if cmd.Flag("scopes") == nil {
		t.Error("missing --scopes flag")
	}
}

func TestAuthTokenCmd_RequiresToken(t *testing.T) {
	cmd := newAuthTokenCmd()

	// Should require exactly 1 arg
	if cmd.Args == nil {
		t.Error("expected Args validator")
	}
}

func TestAuthStatusCmd_Works(t *testing.T) {
	cmd := newAuthStatusCmd()
	if cmd.Use != "status" {
		t.Errorf("expected Use=status, got %s", cmd.Use)
	}
}
```

**Step 2: Write posts tests**

```go
// internal/cmd/posts_test.go
package cmd

import (
	"testing"
)

func TestPostsCreateCmd_Flags(t *testing.T) {
	cmd := newPostsCreateCmd()

	flags := []string{"text", "image", "video", "alt-text", "reply-to"}
	for _, flag := range flags {
		if cmd.Flag(flag) == nil {
			t.Errorf("missing flag: %s", flag)
		}
	}
}

func TestPostsListCmd_HasPagination(t *testing.T) {
	cmd := newPostsListCmd()

	if cmd.Flag("limit") == nil {
		t.Error("missing --limit flag")
	}
	if cmd.Flag("cursor") == nil {
		t.Error("missing --cursor flag")
	}
}

func TestPostsCarouselCmd_RequiresItems(t *testing.T) {
	cmd := newPostsCarouselCmd()

	itemsFlag := cmd.Flag("items")
	if itemsFlag == nil {
		t.Error("missing --items flag")
	}
}
```

**Step 3: Run all tests**

Run: `go test ./internal/... -v`
Expected: All PASS

**Step 4: Commit**

```bash
git add internal/cmd/*_test.go
git commit -m "test(cmd): add unit tests for commands"
```

---

### Task 14: Add Integration Tests

**Files:**
- Modify: `tests/integration/integration_test.go`

**Step 1: Add CLI integration tests**

```go
// tests/integration/cli_test.go
//go:build integration

package integration

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"testing"
)

func TestCLI_Version(t *testing.T) {
	cmd := exec.Command("go", "run", "../../cmd/threads", "version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("version command failed: %v\n%s", err, output)
	}

	if !bytes.Contains(output, []byte("threads")) {
		t.Errorf("expected version output, got: %s", output)
	}
}

func TestCLI_Help(t *testing.T) {
	cmd := exec.Command("go", "run", "../../cmd/threads", "--help")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("help command failed: %v\n%s", err, output)
	}

	expectedCommands := []string{"auth", "posts", "users", "replies", "insights", "search"}
	for _, expected := range expectedCommands {
		if !bytes.Contains(output, []byte(expected)) {
			t.Errorf("expected command %q in help output", expected)
		}
	}
}

func TestCLI_Completion(t *testing.T) {
	shells := []string{"bash", "zsh", "fish", "powershell"}
	for _, shell := range shells {
		t.Run(shell, func(t *testing.T) {
			cmd := exec.Command("go", "run", "../../cmd/threads", "completion", shell)
			output, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("completion %s failed: %v\n%s", shell, err, output)
			}
			if len(output) == 0 {
				t.Error("expected completion output")
			}
		})
	}
}
```

**Step 2: Commit**

```bash
git add tests/integration/cli_test.go
git commit -m "test(integration): add CLI integration tests"
```

---

### Task 15: Update Context Pipeline in Root

**Files:**
- Modify: `internal/cmd/root.go`

**Step 1: Enhance PersistentPreRunE**

Update root.go to use the new iocontext and improved context pipeline:

```go
PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	// Inject IO context (if not already set for testing)
	if !iocontext.HasIO(ctx) {
		ctx = iocontext.WithIO(ctx, iocontext.DefaultIO())
	}

	// Inject output format
	ctx = outfmt.WithFormat(ctx, flags.output)
	ctx = outfmt.WithQuery(ctx, flags.query)
	ctx = outfmt.WithYes(ctx, flags.yes)
	ctx = outfmt.WithLimit(ctx, flags.limit)

	// Inject UI for colorization
	ctx = ui.WithUI(ctx, ui.New(flags.color))

	// Validate flag combinations
	if flags.sortBy != "" && !outfmt.IsJSON(ctx) {
		// Sort only works with JSON output for now
	}

	cmd.SetContext(ctx)
	return nil
},
```

**Step 2: Commit**

```bash
git add internal/cmd/root.go
git commit -m "refactor(root): enhance context pipeline with iocontext"
```

---

### Task 16: Update README with New Commands

**Files:**
- Modify: `README.md`

**Step 1: Update documentation**

Add sections for new commands:

```markdown
### Posts (Extended)

```bash
threads posts carousel --items url1,url2,url3    # Carousel post
threads posts quote POST_ID --text "My take"     # Quote post
threads posts repost POST_ID                      # Repost
```

### Locations

```bash
threads locations search "San Francisco"         # Search by name
threads locations search --lat 37.7 --lng -122.4 # Search by coords
threads locations get LOCATION_ID                # Get details
```

### Rate Limits

```bash
threads ratelimit status                          # Current rate limit
threads ratelimit publishing                      # Publishing quota
```

### User Mentions

```bash
threads users mentions                            # Posts mentioning you
```

### Shell Completion

```bash
threads completion bash > /etc/bash_completion.d/threads
threads completion zsh > "${fpath[1]}/_threads"
threads completion fish > ~/.config/fish/completions/threads.fish
```
```

**Step 2: Commit**

```bash
git add README.md
git commit -m "docs: update README with new commands"
```

---

## Summary

| Phase | Tasks | Estimated Complexity |
|-------|-------|---------------------|
| Phase 1: Foundation | 6 tasks | Infrastructure & patterns |
| Phase 2: Commands | 6 tasks | New API commands |
| Phase 3: Polish | 4 tasks | Tests & documentation |
| **Total** | **16 tasks** | |

---

**Plan complete and saved to `docs/plans/2025-01-01-threads-cli-v2.md`. Two execution options:**

**1. Subagent-Driven (this session)** - I dispatch fresh subagent per task, review between tasks, fast iteration

**2. Parallel Session (separate)** - Open new session with executing-plans, batch execution with checkpoints

**Which approach?**
