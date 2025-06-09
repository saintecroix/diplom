package app

import (
	"database/sql"
	"fmt"
	"github.com/xuri/excelize/v2"
	"log"
)

func ReadExcel(filePath string) ([]map[string]interface{}, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия файла: %v", err)
	}
	defer f.Close()

	sheetName := f.GetSheetName(0)

	headers, err := readHeaders(f, sheetName)
	if err != nil {
		return nil, err
	}

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения строк: %v", err)
	}

	data := make([]map[string]interface{}, 0)
	for _, row := range rows[1:] {
		item := make(map[string]interface{})
		for colIdx, colName := range headers {
			if colIdx < len(row) {
				item[colName] = row[colIdx]
			}
		}
		data = append(data, item)
	}

	return data, nil
}

func readHeaders(f *excelize.File, sheetName string) ([]string, error) {
	rows, err := f.GetRows(sheetName)
	if err != nil || len(rows) == 0 {
		return nil, fmt.Errorf("файл не содержит данных")
	}
	return rows[0], nil
}

func MapColumns(db *sql.DB, excelHeaders []string) (map[string]string, error) {
	mappings := make(map[string]string)
	for _, header := range excelHeaders {
		var originalName string
		query := `
            SELECT "Оригинальное наименование" 
            FROM diplom_raw.naming 
            WHERE $1 = ANY(string_to_array("Альтернативные имена", ';')) 
            OR "Оригинальное наименование" = $1`

		err := db.QueryRow(query, header).Scan(&originalName)
		if err != nil {
			if err == sql.ErrNoRows {
				log.Printf("Предупреждение: колонка '%s' не найдена в таблице naming", header)
				continue
			}
			return nil, fmt.Errorf("ошибка запроса: %v", err)
		}
		mappings[header] = originalName
	}
	return mappings, nil
}
