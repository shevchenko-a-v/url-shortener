package app

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"url-shortener/internal/app/endpoint"
	"url-shortener/internal/app/middleware"
	"url-shortener/internal/app/service"
	"url-shortener/internal/config"
	"url-shortener/internal/pkg/repository/sqlite"
)

type App struct {
	config   *config.Config
	mux      *http.ServeMux
	endpoint *endpoint.Endpoint
	service  *service.Service
}

const (
	configPath = "./configs/local.yaml"
	formatJson = "json"
	logDebug   = "debug"
	logWarn    = "warn"
	logError   = "error"
)

func New() *App {
	c := config.MustLoad(configPath)
	initLogger(c.LogFormat, c.LogLevel)
	slog.Debug(fmt.Sprintf("Loaded config from: %s %+v", configPath, *c))

	storage, err := sqlite.New(c.StoragePath)
	if err != nil {
		slog.Error(fmt.Sprintf("couldn't open repository: %s", err.Error()))
		return nil
	}
	s := service.New(storage)
	e := endpoint.New(slog.Default(), s, c.AliasLength)
	mux := http.NewServeMux()

	mux.Handle("/save", middleware.RequestId(middleware.Logger(http.HandlerFunc(e.SaveUrl))))
	mux.Handle("/", middleware.RequestId(middleware.Logger(http.HandlerFunc(e.Redirect))))
	return &App{config: c, mux: mux, endpoint: e, service: s}
}

func (a *App) Run() error {
	return http.ListenAndServe(a.config.HttpServer.Address, a.mux)
}

func initLogger(format string, level string) {
	var handler slog.Handler
	var logLevel slog.Level
	switch strings.ToLower(level) {
	case logDebug:
		logLevel = slog.LevelDebug
	case logWarn:
		logLevel = slog.LevelWarn
	case logError:
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	switch format {
	case formatJson:
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})
	default:
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})
	}

	slog.SetDefault(slog.New(handler))
}
