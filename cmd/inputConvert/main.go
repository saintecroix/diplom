package main

import "log"

func main() {
	// Подключение к БД (реализуйте отдельно)
	db, err := ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Чтение Excel
	data, err := ReadExcel("114_09032025_1706_empty.xlsx")
	if err != nil {
		log.Fatal(err)
	}

	// Маппинг колонок
	mappings, err := MapColumns(db, data[0].Keys())
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
