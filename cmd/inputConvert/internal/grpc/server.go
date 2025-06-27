package grpc

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"net"

	"github.com/saintecroix/diplom/cmd/inputConvert/internal/app"
	pb "github.com/saintecroix/diplom/internal/api"
	"github.com/saintecroix/diplom/internal/db"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedInputConvertServiceServer
	dbPool *pgxpool.Pool
}

func (s *Server) UploadAndConvertExcelData(ctx context.Context, req *pb.UploadAndConvertExcelDataRequest) (*pb.UploadAndConvertExcelDataResponse, error) {
	log.Info().Msg("===== REQUEST RECEIVED =====")
	log.Info().Str("filename", req.Filename).Int("size", len(req.FileData)).Msg("File received")

	defer func() {
		if r := recover(); r != nil {
			log.Error().Interface("recover", r).Msg("Panic recovered")
		}
	}()
	startTime := time.Now()
	log.Info().Msg("Received request to convert Excel data")

	// 1. Генерируем уникальный job_id
	jobID := uuid.New().String()

	// 2. Чтение данных из Excel
	data, err := app.ReadExcelFromBytes(req.FileData)
	if err != nil {
		log.Error().Err(err).Msg("Error reading Excel file from bytes")
		return &pb.UploadAndConvertExcelDataResponse{
			JobId: jobID,
			Error: err.Error(),
		}, nil
	}

	if len(data) == 0 {
		return &pb.UploadAndConvertExcelDataResponse{
			JobId: jobID,
			Error: "no data in excel file",
		}, nil
	}

	// 3. Получаем ключи из первой записи
	firstRow := data[0]
	keys := make([]string, 0, len(firstRow))
	for k := range firstRow {
		keys = append(keys, k)
	}

	// 4. Сопоставление колонок Excel с колонками БД
	mappings, err := app.MapColumns(s.dbPool, keys)
	if err != nil {
		log.Error().Err(err).Msg("Error mapping columns")
		return &pb.UploadAndConvertExcelDataResponse{
			JobId: jobID,
			Error: err.Error(),
		}, nil
	}

	// 5. Конвертация данных в модели Trip
	trips := make([]db.Trip, 0, len(data))
	for _, item := range data {
		trip := db.Trip{
			ВремяЗагрузкиДанных: time.Now(),
		}

		// Маппинг значений
		for excelCol, dbCol := range mappings {
			value, ok := item[excelCol].(string)
			if !ok {
				continue
			}

			switch dbCol {
			case "Дата и время начала рейса":
				if t, err := time.Parse("2006-01-02 15:04:05", value); err == nil {
					trip.ДатаИВремяНачалаРейса = t
				}
			case "Номер вагона":
				trip.НомерВагона = value
			case "Дорога отправления":
				trip.ДорогаОтправления = value
			case "Дорога назначения":
				trip.ДорогаНазначения = value
			case "Номер накладной":
				trip.НомерНакладной = value
			case "Станция отправления":
				trip.СтанцияОтправления = value
			case "Станция назначения":
				trip.СтанцияНазначения = value
			case "Наименование груза":
				trip.НаименованиеГруза = value
			case "Грузоотправитель":
				trip.Грузоотправитель = value
			case "Грузополучатель":
				trip.Грузополучатель = value
			case "Тип парка (М/Т)":
				trip.ТипПаркаМТ = value
			case "Тип парка (П/Г)":
				trip.ТипПаркаПГ = value
			}
		}
		trips = append(trips, trip)
	}

	// 6. Сохранение данных в БД
	if err := db.BulkCreateTrips(s.dbPool, trips); err != nil {
		log.Error().Err(err).Msg("Error saving trips to DB")
		return &pb.UploadAndConvertExcelDataResponse{
			JobId: jobID,
			Error: fmt.Sprintf("Error saving data to DB: %v", err.Error()),
		}, nil
	}

	// 7. Формируем ответ
	duration := time.Since(startTime)
	log.Info().
		Str("job_id", jobID).
		Int("trips_processed", len(trips)).
		Str("duration", duration.String()).
		Msg("File processed successfully")

	return &pb.UploadAndConvertExcelDataResponse{
		JobId:   jobID,
		Message: fmt.Sprintf("Processed %d trips successfully", len(trips)),
	}, nil
}

func (s *Server) GetJobStatus(ctx context.Context, req *pb.GetJobStatusRequest) (*pb.GetJobStatusResponse, error) {
	// TODO: Реализовать проверку статуса задачи
	return &pb.GetJobStatusResponse{
		Status:  pb.GetJobStatusResponse_COMPLETED,
		Message: "Job status retrieval not implemented yet",
	}, nil
}

// StartGRPCServer запускает gRPC сервер
func StartGRPCServer(dbPool *pgxpool.Pool) error {
	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "50051"
	}

	listen, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		log.Fatal().Msgf("failed to listen: %v", err)
		return err
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(loggingInterceptor),
	)
	pb.RegisterInputConvertServiceServer(grpcServer, &Server{dbPool: dbPool})

	log.Info().Str("port", port).Msg("gRPC server started")
	return grpcServer.Serve(listen)
}

func loggingInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	duration := time.Since(start)

	log.Info().
		Str("method", info.FullMethod).
		Str("duration", duration.String()).
		Err(err).
		Msg("gRPC request processed")

	return resp, err
}
