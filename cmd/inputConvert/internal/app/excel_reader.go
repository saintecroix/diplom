package app

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/xuri/excelize/v2"
)

func ReadExcelFromBytes(fileData []byte) ([]map[string]interface{}, error) {
	// Create a new reader from the byte slice
	reader := bytes.NewReader(fileData)

	f, err := excelize.OpenReader(reader)
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия файла из байтов: %v", err)
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
