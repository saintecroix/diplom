package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"time"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

func loadDBConfig() (*DBConfig, error) {
	port, err := strconv.Atoi(getEnv("DB_PORT", "5432"))
	if err != nil {
		return nil, fmt.Errorf("invalid DB_PORT: %w", err)
	}

	return &DBConfig{
		Host:     getEnv("DB_HOST", "db"),
		Port:     port,
		User:     getEnv("DB_USER", "keril"),
		Password: getEnv("DB_PASSWORD", "pass"),
		DBName:   getEnv("DB_NAME", "my_db"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// ConnectDB устанавливает подключение к БД с бесконечными попытками
func ConnectDB(ctx context.Context) (*sql.DB, error) {
	config, err := loadDBConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load DB config: %w", err)
	}

	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode,
	)

	var db *sql.DB
	attempt := 1
	maxWait := 30 * time.Second

	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("context canceled before DB connection established")
		default:
			db, err = sql.Open("postgres", connStr)
			if err != nil {
				log.Warn().Err(err).Msg("Failed to open DB connection")
				time.Sleep(backoff(attempt, maxWait))
				attempt++
				continue
			}

			if err = db.PingContext(ctx); err != nil {
				log.Warn().Err(err).Int("attempt", attempt).Msg("DB ping failed")
				db.Close()
				time.Sleep(backoff(attempt, maxWait))
				attempt++
				continue
			}

			configureConnectionPool(db)
			log.Info().
				Str("db", config.DBName).
				Str("host", config.Host).
				Msg("Successfully connected to database")
			return db, nil
		}
	}
}

// Экспоненциальная задержка с ограничением максимального времени
func backoff(attempt int, max time.Duration) time.Duration {
	wait := time.Duration(attempt*attempt) * time.Second
	if wait > max {
		return max
	}
	return wait
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

// CloseDB безопасно закрывает подключение с учетом контекста
func CloseDB(ctx context.Context, db *sql.DB) error {
	if db != nil {
		// Сначала пытаемся закрыть все соединения грациозно
		done := make(chan error, 1) // Канал для передачи ошибки
		go func() {
			err := db.Close()
			if err != nil {
				log.Error().Err(err).Msg("Failed to close database connection")
				done <- err // Отправляем ошибку в канал
			} else {
				log.Info().Msg("Database connection closed")
				done <- nil // Отправляем nil в канал, чтобы сообщить об успехе
			}
			close(done)
		}()

		// Ожидаем завершения или отмены контекста
		select {
		case err := <-done: // Получаем ошибку из канала
			if err != nil {
				return err // Возвращаем ошибку
			}
			return nil // Возвращаем nil, если закрытие прошло успешно
		case <-ctx.Done():
			log.Warn().Msg("Forced DB connection closure due to context timeout")
			return ctx.Err() // Возвращаем ошибку из контекста
		}
	}
	return nil // Возвращаем nil, если db == nil
}

// NewPostgresTripRepository создает новый экземпляр PostgresTripRepository
func NewPostgresTripRepository(db *sql.DB) TripRepository {
	return &PostgresTripRepository{db: db}
}

// (Определения структур и функций loadDBConfig, getEnv, ConnectDB, CloseDB, backoff, configureConnectionPool остаются без изменений)

// BulkCreateTrips реализует пакетную вставку данных в таблицу trips
func (r *PostgresTripRepository) BulkCreateTrips(ctx context.Context, trips []Trip) error {
	// Проверяем, что у нас есть данные для вставки
	if len(trips) == 0 {
		return nil
	}

	// 1. Строим SQL-запрос на основе структуры Trip
	tripType := reflect.TypeOf(Trip{})
	columnNames := make([]string, tripType.NumField())
	columns := ""
	placeholders := ""
	for i := 0; i < tripType.NumField(); i++ {
		field := tripType.Field(i)
		columnNames[i] = field.Tag.Get("db")

		columns += fmt.Sprintf("\"%s\"", columnNames[i])
		placeholders += fmt.Sprintf("$%d", i+1)

		if i < tripType.NumField()-1 {
			columns += ", "
			placeholders += ", "
		}
	}

	// Формируем SQL-запрос для пакетной вставки
	sqlStatement := fmt.Sprintf(`
		INSERT INTO diplom_raw.trips (%s)
		VALUES (%s)
	`, columns, placeholders)

	// 2. Готовим statement для пакетной вставки
	stmt, err := r.db.PrepareContext(ctx, sqlStatement)
	if err != nil {
		log.Error().Err(err).Msg("Error preparing statement")
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	// 3. Выполняем транзакцию для пакетной вставки
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		log.Error().Err(err).Msg("Error starting transaction")
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // re-throw panic after rollback
		} else if err != nil {
			tx.Rollback() // err is non-nil; don't change it
			log.Error().Err(err).Msg("Error during transaction, rollback")
		} else {
			err = tx.Commit() // err is nil; if Commit returns error update err
			log.Info().Msg("Transaction committed successfully")
		}
	}()

	// 4. Вставляем данные в цикле, используя подготовленный statement
	for _, trip := range trips {
		v := reflect.ValueOf(trip)
		values := make([]interface{}, v.NumField())
		for i := 0; i < v.NumField(); i++ {
			values[i] = v.Field(i).Interface()
		}

		_, err = tx.StmtContext(ctx, stmt).ExecContext(ctx, values...)
		if err != nil {
			log.Error().Err(err).Msg("Error executing statement")
			return fmt.Errorf("error executing statement: %w", err)
		}
	}

	return nil
}

// GetTripByID - получает запись trips по ID (пример, нужно реализовать)
func (r *PostgresTripRepository) GetTripByID(ctx context.Context, id int64) (*Trip, error) {
	// TODO: Реализовать
	return nil, fmt.Errorf("not implemented")
}

// CreateTrip - создает новую запись trips (пример, нужно реализовать)
func (r *PostgresTripRepository) CreateTrip(ctx context.Context, trip *Trip) (int64, error) {
	// TODO: Реализовать
	return 0, fmt.Errorf("not implemented")
}

// UpdateTrip - обновляет запись trips (пример, нужно реализовать)
func (r *PostgresTripRepository) UpdateTrip(ctx context.Context, trip *Trip) error {
	// TODO: Реализовать
	return fmt.Errorf("not implemented")
}

// DeleteTrip - удаляет запись trips (пример, нужно реализовать)
func (r *PostgresTripRepository) DeleteTrip(ctx context.Context, id int64) error {
	// TODO: Реализовать
	return fmt.Errorf("not implemented")
}

// mapToStruct преобразует map[string]interface{} в структуру Trip
func MapToStruct(m map[string]interface{}, s interface{}) error {
	v := reflect.ValueOf(s).Elem()
	typeOfS := v.Type()

	for i := 0; i < typeOfS.NumField(); i++ {
		field := typeOfS.Field(i)
		fieldName := field.Tag.Get("db") // Используем тег `db` для получения имени поля
		fieldValue, ok := m[fieldName]
		if !ok {
			continue // Пропускаем поля, отсутствующие в map
		}

		fieldValueReflect := reflect.ValueOf(fieldValue)

		//Проверяем, что поле можно установить
		if v.Field(i).CanSet() {
			// Преобразуем тип, если необходимо
			convertedValue, err := convertType(fieldValueReflect, field.Type)
			if err != nil {
				return fmt.Errorf("error converting type for field %s: %w", fieldName, err)
			}

			// Устанавливаем значение поля
			v.Field(i).Set(convertedValue)
		} else {
			return fmt.Errorf("cannot set field %s", fieldName)
		}

	}

	return nil
}

// convertType преобразует типы данных при необходимости
func convertType(src reflect.Value, target reflect.Type) (reflect.Value, error) {
	if src.Type() == target {
		return src, nil // Типы совпадают, преобразование не требуется
	}

	// Преобразование string в time.Time
	if src.Type().Kind() == reflect.String && target == reflect.TypeOf(time.Time{}) {
		timeValue, err := time.Parse(time.RFC3339, src.String())
		if err != nil {
			// Пробуем другие форматы
			timeValue, err = time.Parse("02.01.2006 15:04:05", src.String())
			if err != nil {
				return reflect.Value{}, fmt.Errorf("cannot convert string to time.Time: %w", err)
			}
		}
		return reflect.ValueOf(timeValue), nil
	}
	//Добавьте другие преобразования по мере необходимости

	return reflect.Value{}, fmt.Errorf("unsupported type conversion from %s to %s", src.Type(), target)
}
