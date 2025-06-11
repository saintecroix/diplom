package main

import (
	"context"
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Graceful shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	// 1. Подключение к БД с бесконечными попытками
	dbConn, err := db.ConnectDB(ctx)
	if err != nil {
		log.Error().Msgf("Critical DB connection error: %v", err)
		os.Exit(1)
	}
	defer func() {
		err := db.CloseDB(ctx, dbConn)
		if err != nil {
			log.Error().Msgf("Error closing DB connection: %v", err)
		}
	}()

	// 2. Запуск gRPC-сервера
	port := ":50051" // Порт, на котором будет слушать сервер
	go func() {
		if err := grpc.StartGRPCServer(dbConn, port); err != nil {
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
