# 🌈 slog-human

A professional, high-visibility handler for Go's standard `slog` library. Designed for developers who prioritize debugging speed and data security.

## ✨ Why use slog-human?

Standard `slog` output is great for machines (JSON) but painful for humans during development. `slog-human` transforms your logs into a beautiful, structured, and actionable stream.

* **🚀 Clickable Source Links:** Jump from your terminal directly to the exact line of code (Works in VS Code, iTerm2, GoLand).
* **🛡️ Auto-Redaction:** Automatically hides sensitive keys like `password`, `token`, `api_key`, and `secret`.
* **🔍 ID Highlighting:** Instant visual anchors for `trace_id`, `user_id`, and `request_id`.
* **📂 Multi-line Errors:** Readable error formatting that doesn't break your terminal flow.

## 📦 Installation

To add `slog-human` to your Go project, run:

```bash
go get [github.com/almaz-go/slog-human](https://github.com/almaz-go/slog-human)
🛠 Usage
Integrating slog-human takes only two lines of code. It acts as a drop-in replacement for your default slog handler.

Go
package main

import (
	"log/slog"
	"os"
	"[github.com/almaz-go/slog-human](https://github.com/almaz-go/slog-human)"
)

func main() {
	// 1. Initialize the human-friendly handler
	handler := sloghuman.NewHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true, // Required for clickable source links
		Level:     slog.LevelDebug,
	})

	// 2. Set it as global default
	slog.SetDefault(slog.New(handler))

	// Now use standard slog as usual!
	slog.Info("Service started", "port", 8080, "trace_id", "req-123")
	slog.Error("Auth failed", "password", "secret123", "error", fmt.Errorf("invalid token"))
}
🎨 Features in Action
Secrets: Values for keys like password are replaced with [REDACTED] on a red background.

IDs: Trace and Request IDs are highlighted in magenta with an underline.

Errors: Go error types are printed on a new line with a structured arrow └─>.

Created with focus on Developer Experience (DX).