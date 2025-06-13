package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"strings"
)

// Postgres структура для работы с базой данных PostgreSQL
type Postgres struct {
	ConnString string
}

// ConnectDB устанавливает соединение с базой данных PostgreSQL, используя pgxpool
func ConnectDB(connString string) (*pgxpool.Pool, error) {
	dbpool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	// Проверка соединения (опционально)
	err = PingDB(dbpool)
	if err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	fmt.Println("Successfully connected to PostgreSQL!")
	return dbpool, nil
}

// PingDB проверяет соединение с базой данных
func PingDB(dbpool *pgxpool.Pool) error {
	ctx := context.Background()
	err := dbpool.Ping(ctx)
	if err != nil {
		return fmt.Errorf("ping failed: %w", err)
	}
	return nil
}

// CloseDB закрывает соединение с базой данных
func CloseDB(dbpool *pgxpool.Pool) {
	dbpool.Close()
	fmt.Println("Connection pool closed")
}

// SaveData сохраняет данные в таблицу diplom_raw.trips
func SaveData(dbPool *pgxpool.Pool, data map[string]interface{}) error {
	// 1. Создайте строку запроса INSERT
	//  Здесь нужно динамически построить запрос, основываясь на ключах в data
	//  и значениях в data.  Это более безопасно, чем просто конкатенация строк.

	// 2. Подготовьте placeholders для значений
	columns := make([]string, 0, len(data))
	placeholders := make([]string, 0, len(data))
	values := make([]interface{}, 0, len(data))
	i := 1 // Placeholder index

	for column, value := range data {
		columns = append(columns, column)
		placeholders = append(placeholders, fmt.Sprintf("$%d", i))
		values = append(values, value)
		i++
	}

	// 3.  Соберите запрос
	query := fmt.Sprintf(
		`INSERT INTO diplom_raw.trips (%s) VALUES (%s)`,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "),
	)

	// 4. Выполните запрос
	ctx := context.Background()
	_, err := dbPool.Exec(ctx, query, values...)
	if err != nil {
		log.Printf("Error inserting data into database: %v", err)
		return fmt.Errorf("error inserting data: %w", err)
	}

	log.Printf("Data inserted successfully")
	return nil
}
