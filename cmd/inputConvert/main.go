package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
	"github.com/saintecroix/diplom/cmd/inputConvert/internal/grpc" // Импортируем пакет grpc
	"github.com/saintecroix/diplom/internal/db"
)

func main() {
	// Recovery from panic
	defer func() {
		if r := recover(); r != nil {
			log.Error().Msgf("Recovered from panic: %v", r)
			os.Exit(1)
		}
	}()

	// Graceful shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	// 1. Подключение к БД с бесконечными попытками
	connString := "postgres://keril:pass@db:5432/my_db" // TODO: change to env
	dbConn, err := db.ConnectDB(connString)
	if err != nil {
		log.Error().Msgf("Critical DB connection error: %v", err)
		os.Exit(1)
	}
	defer func() {
		db.CloseDB(dbConn)
		log.Info().Msg("Database connection closed")
	}()

	// 2. Запуск gRPC-сервера
	go func() {
		if err := grpc.StartGRPCServer(dbConn); err != nil { // Исправленный вызов
			log.Error().Msgf("Failed to start gRPC server: %v", err)
			os.Exit(1)
		}
	}()

	// 3. Graceful shutdown
	<-shutdown
	log.Info().Msg("Shutting down...")
	fmt.Println("Gracefully stopping the service")
	os.Exit(0)
}
