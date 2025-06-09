package main

import (
	"context"
	"github.com/saintecroix/diplom/cmd/inputConvert/internal/app"
	"github.com/saintecroix/diplom/internal/db"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Recovery from panic
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic: %v", r)
			os.Exit(1) // Exit with error code if panic occurs
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Graceful shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	// Подключение к БД с бесконечными попытками
	dbConn, err := db.ConnectDB(ctx)
	if err != nil {
		log.Printf("Critical DB connection error: %v", err)
		os.Exit(1) // Exit with error code on critical DB connection error
	}
	defer db.CloseDB(ctx, dbConn)

	// Чтение Excel
	data, err := app.ReadExcel("114_09032025_1706_empty.xlsx")
	if err != nil {
		log.Printf("Error reading Excel: %v", err)
		//  Здесь можно добавить логику повторных попыток чтения, если это необходимо.
		os.Exit(1) // Exit with error code on Excel read error
	}

	firstRow := data[0]

	keys := make([]string, 0, len(firstRow))
	for k := range firstRow {
		keys = append(keys, k)
	}

	// Маппинг колонок
	mappings, err := app.MapColumns(dbConn, keys)
	if err != nil {
		log.Printf("Error mapping columns: %v", err)
		os.Exit(1) // Exit with error code on mapping error
	}

	// Преобразование данных (пункт 2.3)
	for _, item := range data {
		mappedItem := make(map[string]interface{})
		for excelCol, dbCol := range mappings {
			mappedItem[dbCol] = item[excelCol]
		}
		// Далее: сохранение в БД (пункт 2.4)
		// Если здесь возникает ошибка, обработайте её аналогично (логирование, возможно, retry)
	}

	// Ожидаем сигнал завершения
	select {
	case <-shutdown:
		log.Println("Shutdown signal received")
	case <-ctx.Done():
		log.Println("Context canceled")
	}

	log.Println("Service stopped gracefully")
	os.Exit(0) // Exit with success code
}
