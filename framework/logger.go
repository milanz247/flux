package framework

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"flux/config"
)

// NewLogger builds the application slog.Logger. In local development it
// writes human-readable text to stdout; in production it writes JSON to both
// stdout and storage/logs/flux.log.
func NewLogger(cfg *config.Config) *slog.Logger {
	level := slog.LevelInfo
	if cfg.App.Debug {
		level = slog.LevelDebug
	}

	var writer io.Writer = os.Stdout
	if file := openLogFile(); file != nil {
		writer = io.MultiWriter(os.Stdout, file)
	}

	var handler slog.Handler
	if cfg.App.Env == "production" {
		handler = slog.NewJSONHandler(writer, &slog.HandlerOptions{Level: level})
	} else {
		handler = slog.NewTextHandler(writer, &slog.HandlerOptions{Level: level})
	}

	return slog.New(handler).With(slog.String("app", cfg.App.Name))
}

func openLogFile() *os.File {
	dir := filepath.Join("storage", "logs")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil
	}
	file, err := os.OpenFile(filepath.Join(dir, "flux.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return nil
	}
	return file
}
