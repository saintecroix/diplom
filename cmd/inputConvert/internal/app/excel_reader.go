package app

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"github.com/xuri/excelize/v2"
)

func ReadExcelFromBytes(fileData []byte) ([]map[string]interface{}, error) {
	reader := bytes.NewReader(fileData)

	f, err := excelize.OpenReader(reader)
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия файла из байтов: %v", err)
	}
	defer f.Close()

	sheetName := f.GetSheetName(0)
	if sheetName == "" {
		return nil, fmt.Errorf("excel файл не содержит листов")
	}

	// Получаем все строки
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения строк: %v", err)
	}

	if len(rows) < 2 {
		return nil, fmt.Errorf("файл должен содержать хотя бы одну строку данных после заголовка")
	}

	// Читаем заголовки
	headers := rows[0]
	for i, h := range headers {
		headers[i] = strings.TrimSpace(h)
	}

	// Обрабатываем данные
	data := make([]map[string]interface{}, 0)
	for _, row := range rows[1:] {
		item := make(map[string]interface{})
		for colIdx, colName := range headers {
			if colIdx < len(row) {
				// Пробуем преобразовать в дату, если возможно
				if value, err := f.GetCellValue(sheetName, fmt.Sprintf("%s%d", toAlpha(colIdx+1), colIdx+2)); err == nil {
					if t, err := time.Parse("2006-01-02 15:04:05", value); err == nil {
						item[colName] = t
						continue
					}
				}
				item[colName] = row[colIdx]
			}
		}
		data = append(data, item)
	}

	return data, nil
}

// Вспомогательная функция для преобразования индекса в букву столбца
func toAlpha(n int) string {
	if n < 1 {
		return ""
	}
	return toAlpha((n-1)/26) + string(rune('A'+(n-1)%26))
}

func MapColumns(db *pgxpool.Pool, excelHeaders []string) (map[string]string, error) {
	mappings := make(map[string]string)
	for _, header := range excelHeaders {
		var originalName string
		query := `
            SELECT "Оригинальное наименование" 
            FROM diplom_raw.naming 
            WHERE $1 = ANY(string_to_array("Альтернативные имена", ';')) 
            OR "Оригинальное наименование" = $1`

		err := db.QueryRow(context.Background(), query, header).Scan(&originalName)
		if err != nil {
			if err == pgx.ErrNoRows {
				log.Warn().Str("header", header).Msg("Column not found in naming table")
				mappings[header] = header // Используем оригинальное название
				continue
			}
			return nil, fmt.Errorf("ошибка запроса: %v", err)
		}
		mappings[header] = originalName
	}
	return mappings, nil
}
