package sloghuman

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"sync"

	"github.com/fatih/color"
)

// Handler implements slog.Handler to provide human-readable, colored output.
type Handler struct {
	opts slog.HandlerOptions
	out  io.Writer
	mu   *sync.Mutex
}

func NewHandler(out io.Writer, opts *slog.HandlerOptions) *Handler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}
	return &Handler{
		out:  out,
		opts: *opts,
		mu:   &sync.Mutex{},
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

	// 1. Format Time
	timeStr := r.Time.Format("15:04:05")
	timeCol := color.New(color.FgHiBlack).Sprint(timeStr)

	// 2. Format Level with colors
	levelStr := "[" + r.Level.String() + "]"
	var levelCol string
	switch {
	case r.Level >= slog.LevelError:
		levelCol = color.New(color.FgRed, color.Bold).Sprint(levelStr)
	case r.Level >= slog.LevelWarn:
		levelCol = color.New(color.FgYellow).Sprint(levelStr)
	default:
		levelCol = color.New(color.FgCyan).Sprint(levelStr)
	}

	// 3. Format Message
	msgCol := color.New(color.FgWhite, color.Bold).Sprint(r.Message)

	// 4. Format Attributes (Key-Value pairs)
	attrs := ""
	r.Attrs(func(a slog.Attr) bool {
		k := color.New(color.FgHiCyan).Sprint(a.Key)
		v := fmt.Sprintf("%v", a.Value.Any())
		attrs += fmt.Sprintf(" %s=%s", k, v)
		return true
	})

	// Print the final line
	fmt.Fprintf(h.out, "%s %s %s%s\n", timeCol, levelCol, msgCol, attrs)

	return nil
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	// For simplicity in MVP, we return the same handler.
	// Production-grade would handle nested attributes.
	return h
}

func (h *Handler) WithGroup(name string) slog.Handler {
	return h
}