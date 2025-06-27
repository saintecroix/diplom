package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
	"github.com/saintecroix/diplom/cmd/inputConvert/internal/grpc" // Импортируем пакет grpc
	"github.com/saintecroix/diplom/internal/db"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Error().Interface("panic", r).Msg("Recovered from panic")
			panic(r)
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
	CheckTableSchema(dbConn)
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

func CheckTableSchema(dbPool *pgxpool.Pool) {
	query := `SELECT column_name, data_type 
             FROM information_schema.columns 
             WHERE table_schema = 'diplom_raw' 
             AND table_name = 'trips'`

	rows, err := dbPool.Query(context.Background(), query)
	if err != nil {
		log.Error().Err(err).Msg("Failed to check table schema")
		return
	}

	for rows.Next() {
		var colName, dataType string
		rows.Scan(&colName, &dataType)
		log.Info().Str("column", colName).Str("type", dataType).Msg("Table column")
	}
}
