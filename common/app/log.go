package app

import (
	"context"
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

const (
	configKeyLogLevel  = "log.level"
	configKeyLogFormat = "log.format"
)

func stringToSlogLevel(
	level string,
) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

type logHandler struct {
	slog.Handler
}

func (h *logHandler) Handle(ctx context.Context, r slog.Record) error {
	if cmd != "" {
		r.AddAttrs(slog.String("cmd", cmd))
	}
	if hostname != "" {
		r.AddAttrs(slog.String("hostname", hostname))
	}
	return h.Handler.Handle(ctx, r)
}

func initLog() {
	logLevel := stringToSlogLevel(config.GetString(configKeyLogLevel))

	var handler slog.Handler
	switch config.GetString(configKeyLogFormat) {
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})
	case "plain":
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})
	default:
		handler = tint.NewHandler(os.Stdout, &tint.Options{Level: logLevel})
	}
	handler = &logHandler{handler}
	slog.SetDefault(slog.New(handler))
}
