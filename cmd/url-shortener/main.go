package main

import (
	"log/slog"
	"url-shortener/internal/pkg/app"
)

func main() {
	instance := app.New()
	if instance == nil || instance.Run() != nil {
		slog.Error("error during running app")
	}
}
