package main

import (
	"fmt"
	"log/slog"
	"os"
	"url-shortener/internal/config"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/storage/postgresql"
)

const (
	envLocal = "local"
	envDev = "dev"
	envProd = "prod"
)

func main () {
	fmt.Println("Запущено")

	cfg := config.MustLoad()

	fmt.Println(cfg.Conn_str)

	log := setupLogger(cfg.Env)

	// log.Info("starting url-shortener", slog.String("env", cfg.Env))
	// log.Debug("debug messages are enabled")

	storage, err := postgresql.New(cfg.Conn_str)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	id, err := storage.SaveURL("https://google.com", "google")
	if err != nil {
		log.Error("failed to save url", sl.Err(err))
		os.Exit(1)
	}
	log.Info("saved url", slog.Int64("id", id))

	url, err := storage.GETUrl("https://google.com")

	if err != nil {
		log.Error("failed to get url", sl.Err(err))
		os.Exit(1)
	}

	log.Info("get url", slog.String("url", url))

	err = storage.DeleteURL("https://google.com")

	if err != nil {
		log.Error("failed to delete url", sl.Err(err))
		os.Exit(1)
	}

	log.Info("deleted url")

	_ = storage
}

func setupLogger (env string) *slog.Logger {

	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
		
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}