package db

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	_ "strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

// BulkCreateTrips сохраняет пакет поездок в БД
func BulkCreateTrips(dbPool *pgxpool.Pool, trips []Trip) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := dbPool.Ping(ctx); err != nil {
		log.Error().Err(err).Msg("Database ping failed")
		return err
	}
	log.Info().Int("count", len(trips)).Msg("Starting bulk insert")

	// Временное логирование первой записи
	if len(trips) > 0 {
		tripJson, _ := json.Marshal(trips[0])
		log.Info().RawJSON("first_trip", tripJson).Msg("First trip data")
	}
	if len(trips) == 0 {
		return nil
	}

	tx, err := dbPool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Создаем пакетный обработчик
	batch := &pgx.Batch{}

	// Подготавливаем запрос
	query := `INSERT INTO diplom_raw.trips (
		"Дата и время начала рейса",
		"Номер вагона",
		"Дорога отправления",
		"Дорога назначения",
		"Номер накладной",
		"Станция отправления",
		"Станция назначения",
		"Наименование груза",
		"Грузоотправитель",
		"Грузополучатель",
		"Тип парка (М/Т)",
		"Тип парка (П/Г)",
		"Время загрузки данных"
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`

	// Добавляем все поездки в пакет
	for _, trip := range trips {
		batch.Queue(query,
			trip.ДатаИВремяНачалаРейса,
			trip.НомерВагона,
			trip.ДорогаОтправления,
			trip.ДорогаНазначения,
			trip.НомерНакладной,
			trip.СтанцияОтправления,
			trip.СтанцияНазначения,
			trip.НаименованиеГруза,
			trip.Грузоотправитель,
			trip.Грузополучатель,
			trip.ТипПаркаМТ,
			trip.ТипПаркаПГ,
			trip.ВремяЗагрузкиДанных,
		)
	}

	// Отправляем пакет
	br := tx.SendBatch(ctx, batch)
	defer br.Close()

	// Обрабатываем результаты
	for range trips {
		_, err := br.Exec()
		if err != nil {
			return fmt.Errorf("batch insert error: %w", err)
		}
	}

	// Фиксируем транзакцию
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("transaction commit error: %w", err)
	}

	log.Info().Int("count", len(trips)).Msg("Trips batch inserted successfully")
	return nil
}

// ConnectDB устанавливает соединение с базой данных PostgreSQL
func ConnectDB(connString string) (*pgxpool.Pool, error) {
	log.Info().Str("conn_string", maskPassword(connString)).Msg("Connecting to DB")
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse connection string: %w", err)
	}

	config.MaxConns = 10
	config.MinConns = 2
	config.HealthCheckPeriod = 1 * time.Minute

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Проверяем соединение
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Info().Msg("Database connection established")
	return pool, nil
}

func maskPassword(connString string) string {
	return regexp.MustCompile(`password=\w+`).ReplaceAllString(connString, "password=***")
}

// CloseDB закрывает соединение с базой данных
func CloseDB(pool *pgxpool.Pool) {
	if pool != nil {
		pool.Close()
		log.Info().Msg("Database connection closed")
	}
}
