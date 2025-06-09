package grpc

import (
	"context"
	"database/sql"
	"log"

	"github.com/saintecroix/diplom/cmd/inputConvert/internal/app"
	pb "github.com/saintecroix/diplom/cmd/inputConvert/proto" // Замените на правильный путь
	_ "github.com/saintecroix/diplom/internal/db"
)

type Server struct {
	pb.UnimplementedInputConvertServiceServer // Важно для совместимости gRPC
	dbConn                                    *sql.DB
}

func NewServer(dbConn *sql.DB) *Server {
	return &Server{dbConn: dbConn}
}

func (s *Server) ConvertExcelData(ctx context.Context, req *pb.ConvertExcelDataRequest) (*pb.ConvertExcelDataResponse, error) {
	filePath := req.FilePath
	log.Printf("Received request to convert Excel data from: %s", filePath)

	data, err := app.ReadExcel(filePath)
	if err != nil {
		log.Printf("Error reading Excel file: %v", err)
		return &pb.ConvertExcelDataResponse{Error: err.Error()}, nil // Обратите внимание на обработку ошибок
	}

	firstRow := data[0]

	keys := make([]string, 0, len(firstRow))
	for k := range firstRow {
		keys = append(keys, k)
	}

	mappings, err := app.MapColumns(s.dbConn, keys) // Передаем dbConn
	if err != nil {
		log.Printf("Error mapping columns: %v", err)
		return &pb.ConvertExcelDataResponse{Error: err.Error()}, nil
	}

	var results []string // Здесь будет логика сохранения (пока для примера)

	for _, item := range data {
		mappedItem := make(map[string]interface{})
		for excelCol, dbCol := range mappings {
			mappedItem[dbCol] = item[excelCol]
		}
		// Далее:  Здесь можно реализовать логику сохранения данных в БД
		// Например:  err = db.SaveData(s.dbConn, mappedItem)
		//results = append(results, "OK")
	}

	return &pb.ConvertExcelDataResponse{Results: results}, nil // Возвращаем ответ
}
