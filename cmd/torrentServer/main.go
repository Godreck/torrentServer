// cmd/torrentServer
package main

import (
	"log/slog"
	"net/http"
	"os"
	"torrentServer/http_server/handlers/search"
	"torrentServer/internal/config"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

// func init() {
// 	// loads values from .env into the system
// 	if err := godotenv.Load(); err != nil {
// 		log.Print("No .env file found")
// 	}
// }

func main() {
	// загружаем конфиг
	cfg := config.MustLoad()

	// if err := godotenv.Load("../../.env"); err != nil {
	// 	log.Printf("Error loading .env file: %v", err)
	// }

	// Теперь можно безопасно читать переменные окружения
	// apiURL := os.Getenv("JACKETT_API_URL")
	// apiKey := os.Getenv("JACKETT_API_KEY")

	log := setupLogger(cfg.Env)
	log = log.With(slog.String("env", cfg.Env)) // к каждому сообщению будет добавляться поле с информацией о текущем окружении

	log.Info("initializing server", slog.String("address", cfg.Address)) // Помимо сообщения выведем параметр с адресом
	log.Debug("logger debug mode enabled")

	// router := chi.NewRouter()

	// router.Use(middleware.Logger)

	http.HandleFunc("/search", search.SearchHandler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Error(err.Error())
	}
}

func setupLogger(env string) *slog.Logger {
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
