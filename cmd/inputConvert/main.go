package main

import (
	"github.com/saintecroix/diplom/cmd/inputConvert/internal/app"
	"github.com/saintecroix/diplom/internal/db"
	"log"
)

func main() {
	// Подключение к БД (реализуйте отдельно)
	dbConn, err := db.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.CloseDB(dbConn)

	// Чтение Excel
	data, err := app.ReadExcel("114_09032025_1706_empty.xlsx")
	if err != nil {
		log.Fatal(err)
	}

	firstRow := data[0]

	keys := make([]string, 0, len(firstRow))
	for k := range firstRow {
		keys = append(keys, k)
	}

	// Маппинг колонок
	mappings, err := app.MapColumns(dbConn, keys)
	if err != nil {
		log.Fatal(err)
	}

	// Преобразование данных (пункт 2.3)
	for _, item := range data {
		mappedItem := make(map[string]interface{})
		for excelCol, dbCol := range mappings {
			mappedItem[dbCol] = item[excelCol]
		}
		// Далее: сохранение в БД (пункт 2.4)
	}
}
