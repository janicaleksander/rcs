package utils

import (
	"log/slog"
	"os"
	"time"
)

var Logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

const (
	WaitTime = 3 * time.Second
)
