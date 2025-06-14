package grpc

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"

	"github.com/rs/zerolog/log" // Import zerolog
	"net"

	"github.com/saintecroix/diplom/cmd/inputConvert/internal/app"
	pb "github.com/saintecroix/diplom/internal/api"
	db "github.com/saintecroix/diplom/internal/db"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedInputConvertServiceServer
	dbConn *pgxpool.Pool
}

func (s *Server) UploadAndConvertExcelData(ctx context.Context, req *pb.UploadAndConvertExcelDataRequest) (*pb.UploadAndConvertExcelDataResponse, error) {
	log.Info().Msg("Received request to convert Excel data")

	connString := "postgres://keril:pass@db:5432/my_db" //todo change to env
	dbpool, err := db.ConnectDB(connString)             // Используем db.ConnectDB
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to db: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	// 1. Сгенерируйте уникальный job_id
	jobID := uuid.New().String()

	// 2. Чтение данных из Excel
	data, err := app.ReadExcelFromBytes(req.FileData)
	if err != nil {
		log.Error().Msgf("Error reading Excel file from bytes: %v", err)
		return &pb.UploadAndConvertExcelDataResponse{
			JobId: jobID, // Still return jobID even on error
			Error: err.Error(),
		}, nil
	}

	// 3. Получение первой строки (ключи)
	firstRow := data[0]
	keys := make([]string, 0, len(firstRow))
	for k := range firstRow {
		keys = append(keys, k)
	}

	// 4. Сопоставление колонок Excel с колонками БД
	mappings, err := app.MapColumns(s.dbConn, keys) // Передаем dbConn
	if err != nil {
		log.Error().Msgf("Error mapping columns: %v", err)
		return &pb.UploadAndConvertExcelDataResponse{
			JobId: jobID,
			Error: err.Error(),
		}, nil
	}

	// 5. Обработка и сохранение данных (пока что пример)
	var results []string
	for _, item := range data {
		mappedItem := make(map[string]interface{})
		for excelCol, dbCol := range mappings {
			mappedItem[dbCol] = item[excelCol]
		}

		// 6. Реализация логики сохранения данных в БД
		err = db.SaveData(s.dbConn, mappedItem)
		if err != nil {
			log.Error().Msgf("Error saving data to DB: %v", err)
			return &pb.UploadAndConvertExcelDataResponse{
				JobId: jobID,
				Error: fmt.Sprintf("Error saving data to DB: %v", err.Error()),
			}, nil
		}

		results = append(results, "OK")
	}

	// 7.  Возвращаем UploadAndConvertExcelDataResponse с job_id и сообщением
	log.Info().Msgf("File processed successfully. Job ID: %s", jobID)
	return &pb.UploadAndConvertExcelDataResponse{
		JobId:   jobID,
		Message: "File processed successfully",
	}, nil
}

// StartGRPCServer запускает gRPC сервер
func StartGRPCServer(dbConn *pgxpool.Pool) error {
	listen, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal().Msgf("failed to listen: %v", err)
		return err
	}

	// Create gRPC server
	grpcServer := grpc.NewServer()
	pb.RegisterInputConvertServiceServer(grpcServer, &Server{dbConn: dbConn}) // Регистрируем gRPC сервис

	log.Info().Msg("gRPC server listening at 50051")
	return grpcServer.Serve(listen)
}
