package sloghuman

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"runtime"
	"strings"
	"sync"

	"github.com/fatih/color"
)

type Handler struct {
	opts           slog.HandlerOptions
	out            io.Writer
	mu             *sync.Mutex
	redactKeys     []string // Ключи для скрытия (Feature 4)
	highlightKeys  []string // Ключи для выделения (Feature 3)
}

func NewHandler(out io.Writer, opts *slog.HandlerOptions) *Handler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}
	return &Handler{
		out:           out,
		opts:          *opts,
		mu:            &sync.Mutex{},
		redactKeys:    []string{"password", "token", "secret", "cookie", "api_key"},
		highlightKeys: []string{"trace_id", "request_id", "user_id"},
	}
}

func (h *Handler) Enabled(_ context.Context, level slog.Level) bool {
	minLevel := slog.LevelInfo
	if h.opts.Level != nil {
		minLevel = h.opts.Level.Level()
	}
	return level >= minLevel
}

func (h *Handler) Handle(_ context.Context, r slog.Record) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	// 1. Time (Black)
	timeCol := color.New(color.FgHiBlack).Sprint(r.Time.Format("15:04:05"))

	// 2. Level (Colored)
	levelStr := fmt.Sprintf("[%s]", r.Level.String())
	var levelCol string
	switch {
	case r.Level >= slog.LevelError:
		levelCol = color.New(color.FgRed, color.Bold).Sprint(levelStr)
	case r.Level >= slog.LevelWarn:
		levelCol = color.New(color.FgYellow).Sprint(levelStr)
	default:
		levelCol = color.New(color.FgCyan).Sprint(levelStr)
	}

	// 3. Source Link (Feature 1: Absolute path for better clickability)
    var sourceCol string
    if h.opts.AddSource && r.PC != 0 {
        fs := runtime.CallersFrames([]uintptr{r.PC})
        f, _ := fs.Next()
        // f.File содержит полный путь, например /Users/almaz/slog-human/example/main.go
        sourceCol = color.New(color.FgHiBlack, color.Italic).Sprintf("%s:%d", f.File, f.Line)
    }

	// 4. Message (Bold White)
	msgCol := color.New(color.FgWhite, color.Bold).Sprint(r.Message)

	// 5. Attributes Processing (Features 2, 3, 4)
	attrsStr := ""
	r.Attrs(func(a slog.Attr) bool {
		key := strings.ToLower(a.Key)
		val := a.Value.Any()

		// Feature 4: Secret Redaction
		for _, rk := range h.redactKeys {
			if key == rk {
				val = color.New(color.BgRed, color.FgWhite).Sprint("[REDACTED]")
				break
			}
		}

		// Feature 3: ID Highlighting
		isID := false
		for _, hk := range h.highlightKeys {
			if key == hk {
				isID = true
				break
			}
		}

		kCol := color.New(color.FgHiCyan).Sprint(a.Key)
		vCol := fmt.Sprintf("%v", val)

		if isID {
			vCol = color.New(color.FgHiMagenta, color.Underline).Sprint(val)
		}

		// Feature 2: Multi-line for Errors
		if err, ok := val.(error); ok {
			vCol = color.New(color.FgRed).Sprintf("\n  └─> %v", err)
		}

		attrsStr += fmt.Sprintf(" %s=%s", kCol, vCol)
		return true
	})

	// Final Output
	fmt.Fprintf(h.out, "%s %s %s %s%s\n", timeCol, levelCol, sourceCol, msgCol, attrsStr)

	return nil
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler { return h }
func (h *Handler) WithGroup(name string) slog.Handler      { return h }