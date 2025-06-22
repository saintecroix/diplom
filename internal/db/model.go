package db

import (
	"context"
	"database/sql"
	"time"
)

type Trip struct {
	ДатаИВремяНачалаРейса time.Time `db:"Дата и время начала рейса"`
	НомерВагона           string    `db:"Номер вагона"`
	ДорогаОтправления     string    `db:"Дорога отправления"`
	ДорогаНазначения      string    `db:"Дорога назначения"`
	НомерНакладной        string    `db:"Номер накладной"`
	СтанцияОтправления    string    `db:"Станция отправления"`
	СтанцияНазначения     string    `db:"Станция назначения"`
	НаименованиеГруза     string    `db:"Наименование груза"`
	Грузоотправитель      string    `db:"Грузоотправитель"`
	Грузополучатель       string    `db:"Грузополучатель"`
	ТипПаркаМТ            string    `db:"Тип парка (М/Т)"`
	ТипПаркаПГ            string    `db:"Тип парка (П/Г)"`
	ВремяЗагрузкиДанных   time.Time `db:"Время загрузки данных"`
}

// TripRepository определяет интерфейс для работы с таблицей trips
type TripRepository interface {
	GetTripByID(ctx context.Context, id int64) (*Trip, error)
	CreateTrip(ctx context.Context, trip *Trip) (int64, error)
	UpdateTrip(ctx context.Context, trip *Trip) error
	DeleteTrip(ctx context.Context, id int64) error
	BulkCreateTrips(ctx context.Context, trips []Trip) error
}

// PostgresTripRepository - реализация TripRepository для PostgreSQL
type PostgresTripRepository struct {
	db *sql.DB
}

// BulkCreateTrips реализация пакетного сохранения
func (r *PostgresTripRepository) BulkCreateTrips(ctx context.Context, trips []Trip) error {
	// Реализация будет аналогична функции в postgres.go
	// Можно использовать общую реализацию
	return nil
}
