package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/saintecroix/diplom/cmd/webUI/internal/web" // Импортируем наш пакет
)

func main() {
	// Настройка логгера
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logger := log.With().Logger()

	// Чтение конфигурации
	port := getEnvString("PORT", "8080")
	uploadDir := getEnvString("UPLOAD_DIR", "./uploads")
	inputConvertAddr := getEnvString("INPUT_CONVERT_ADDRESS", "input-convert:50051")

	// Создаем папку для загрузок
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		logger.Fatal().Err(err).Msg("Failed to create upload directory")
	}

	// Создаем мультиплексор
	mux := http.NewServeMux()

	// Обслуживаем статические файлы
	fs := http.FileServer(http.Dir("./cmd/webUI/internal/web/static"))

	// Регистрируем корневой путь для статики
	mux.Handle("/", fs)
	logger.Info().Str("path", "./cmd/webUI/internal/web/static").Msg("Serving static files")

	// Регистрируем обработчики API
	web.RegisterHandlers(mux, &logger, uploadDir)

	// Настройка HTTP сервера
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Info().
		Str("port", port).
		Str("upload_dir", uploadDir).
		Str("grpc_addr", inputConvertAddr).
		Msg("Server starting")

	// Запускаем сервер
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatal().Err(err).Msg("Server failed")
	}
}

// Вспомогательные функции для чтения переменных окружения
func getEnvString(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
