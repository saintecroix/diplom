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
	ВремяЗагрузкиДанных   time.Time `db:"Время загрузки данных"` // Необязательно, ставится автоматически
}

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// TripRepository определяет интерфейс для работы с таблицей trips
type TripRepository interface {
	GetTripByID(ctx context.Context, id int64) (*Trip, error)
	CreateTrip(ctx context.Context, trip *Trip) (int64, error)
	UpdateTrip(ctx context.Context, trip *Trip) error
	DeleteTrip(ctx context.Context, id int64) error
	BulkCreateTrips(ctx context.Context, trips []Trip) error // Добавим метод для пакетной вставки
}

// PostgresTripRepository - реализация TripRepository для PostgreSQL
type PostgresTripRepository struct {
	db *sql.DB
}
