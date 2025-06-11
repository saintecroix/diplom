package grpc

import (
	"context"
	"database/sql"

	"github.com/rs/zerolog/log" // Import zerolog
	"net"

	"github.com/saintecroix/diplom/cmd/inputConvert/internal/app"
	pb "github.com/saintecroix/diplom/cmd/inputConvert/proto" // Замените на правильный путь
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	pb.UnimplementedInputConvertServiceServer // Важно для совместимости gRPC
	dbConn                                    *sql.DB
}

func NewServer(dbConn *sql.DB) *Server {
	return &Server{dbConn: dbConn}
}

func (s *Server) ConvertExcelData(ctx context.Context, req *pb.ConvertExcelDataRequest) (*pb.ConvertExcelDataResponse, error) {
	log.Info().Msg("Received request to convert Excel data")

	data, err := app.ReadExcelFromBytes(req.Filedata)
	if err != nil {
		log.Error().Msgf("Error reading Excel file from bytes: %v", err)
		return &pb.ConvertExcelDataResponse{Error: err.Error()}, nil // Обратите внимание на обработку ошибок
	}

	firstRow := data[0]

	keys := make([]string, 0, len(firstRow))
	for k := range firstRow {
		keys = append(keys, k)
	}

	mappings, err := app.MapColumns(s.dbConn, keys) // Передаем dbConn
	if err != nil {
		log.Error().Msgf("Error mapping columns: %v", err)
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

// StartGRPCServer запускает gRPC-сервер.
func StartGRPCServer(dbConn *sql.DB, port string) error {
	lis, err := net.Listen("tcp", port) // Порт, на котором будет слушать сервер
	if err != nil {
		log.Fatal().Msgf("Failed to listen: %v", err)
		return err
	}

	s := NewServer(dbConn)
	grpcServer := grpc.NewServer()
	pb.RegisterInputConvertServiceServer(grpcServer, s)

	// Регистрация reflection API (для отладки)
	reflection.Register(grpcServer)

	log.Info().Msgf("Starting gRPC server on port %s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal().Msgf("Failed to serve: %v", err)
		return err
	}

	return nil
}
