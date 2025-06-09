package db

import "time"

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
