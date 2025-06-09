package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"time"
)

const (
	dbHost     = "db" // Имя сервиса в Docker Compose
	dbPort     = 5432
	dbUser     = "keril"
	dbPassword = "pass"
	dbName     = "my_db"
	dbSSLMode  = "disable" // "require" для продакшена
)

func ConnectDB() (*sql.DB, error) {
	// Формируем строку подключения из констант
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode,
	)

	// Открываем соединение
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Проверяем подключение
	if err = pingDatabase(db); err != nil {
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	// Настраиваем пул соединений
	configureConnectionPool(db)

	log.Info().
		Str("db", dbName).
		Str("host", dbHost).
		Msg("Successfully connected to database")
	return db, nil
}

func pingDatabase(db *sql.DB) error {
	var err error
	maxAttempts := 5
	for i := 1; i <= maxAttempts; i++ {
		err = db.Ping()
		if err == nil {
			return nil
		}

		wait := time.Duration(i*i) * time.Second
		log.Warn().Err(err).
			Int("attempt", i).
			Dur("wait", wait).
			Msg("Database connection failed, retrying...")

		time.Sleep(wait)
	}
	return fmt.Errorf("failed to connect after %d attempts: %w", maxAttempts, err)
}

func configureConnectionPool(db *sql.DB) {
	maxOpenConns := 25
	maxIdleConns := 5
	maxLifetime := 30 * time.Minute
	idleTimeout := 5 * time.Minute

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(maxLifetime)
	db.SetConnMaxIdleTime(idleTimeout)

	log.Info().
		Int("max_open_conns", maxOpenConns).
		Int("max_idle_conns", maxIdleConns).
		Dur("max_lifetime", maxLifetime).
		Dur("idle_timeout", idleTimeout).
		Msg("Configured database connection pool")
}

func CloseDB(db *sql.DB) {
	if db != nil {
		if err := db.Close(); err != nil {
			log.Error().Err(err).Msg("Failed to close database connection")
		} else {
			log.Info().Msg("Database connection closed")
		}
	}
}
