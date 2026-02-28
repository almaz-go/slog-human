package main

import (
	"errors"
	"log/slog"
	"os"
	"github.com/almaz-go/slog-human"
)

func main() {
	handler := sloghuman.NewHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true, // Включаем кликабельные ссылки
	})

	logger := slog.New(handler)
	slog.SetDefault(logger)

	// Тест Feature 1 & 3: Ссылка на код + Подсветка ID
	slog.Info("Checking user session", "user_id", "user_99", "trace_id", "abc-123")

	// Тест Feature 4: Скрытие секретов
	slog.Warn("Login attempt", "user", "admin", "password", "123456")

	// Тест Feature 2: Красивый вывод ошибок
	err := errors.New("connection reset by peer")
	slog.Error("Network failure", "error", err, "retry", true)
}