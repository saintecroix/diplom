package web

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rs/zerolog"
	api "github.com/saintecroix/diplom/internal/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// TemplateData структура для передачи данных в шаблон
type TemplateData struct {
	Title   string
	Header  string
	Year    int
	Content interface{} // Можно использовать interface{} для разных типов контента
}

// Web структура для хранения зависимостей обработчиков
type Web struct {
	client    api.InputConvertServiceClient
	templates *template.Template
	logger    *zerolog.Logger
	uploadDir string
}

var w *Web

// RegisterHandlers регистрирует обработчики для веб-интерфейса
func RegisterHandlers(mux *http.ServeMux, logger *zerolog.Logger, uploadDir string) {
	var err error
	templates, err := LoadTemplates()
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to load templates")
	}

	client, err := NewInputConvertClient(logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to create gRPC client")
	}

	w = &Web{
		client:    client,
		templates: templates,
		logger:    logger,
		uploadDir: uploadDir,
	}

	// Создаем папку для загрузок
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		logger.Fatal().Err(err).Str("dir", uploadDir).Msg("Failed to create upload directory")
	}

	mux.HandleFunc("POST /api/upload", w.uploadHandler) // Обработчик для /trips/upload
}

// NewInputConvertClient создает gRPC клиент для inputConvert сервиса
func NewInputConvertClient(logger *zerolog.Logger) (api.InputConvertServiceClient, error) {
	inputConvertAddress := os.Getenv("INPUT_CONVERT_ADDRESS")
	if inputConvertAddress == "" {
		inputConvertAddress = "bus_log:50051"
		logger.Info().Str("address", inputConvertAddress).Msg("INPUT_CONVERT_ADDRESS not set, using default")
	}

	conn, err := grpc.Dial(inputConvertAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to dial inputConvert service: %w", err)
	}

	return api.NewInputConvertServiceClient(conn), nil
}

func (w *Web) uploadHandler(wr http.ResponseWriter, r *http.Request) {
	wr.Header().Set("Content-Type", "application/json")

	var fileData []byte
	var fileName string
	var err error

	contentType := r.Header.Get("Content-Type")

	if strings.HasPrefix(contentType, "multipart/form-data") {
		err = r.ParseMultipartForm(10 << 20)
		if err != nil {
			w.logger.Error().Err(err).Msg("Error parsing multipart form")
			wr.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(wr).Encode(map[string]string{
				"status":  "error",
				"message": "Ошибка при разборе формы: " + err.Error(),
			})
			return
		}

		file, handler, err := r.FormFile("file")
		if err != nil {
			w.logger.Error().Err(err).Msg("Error retrieving file from form")
			wr.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(wr).Encode(map[string]string{
				"status":  "error",
				"message": "Ошибка при получении файла: " + err.Error(),
			})
			return
		}
		defer file.Close()

		fileData, err = io.ReadAll(file)
		if err != nil {
			w.logger.Error().Err(err).Msg("Error reading file")
			wr.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(wr).Encode(map[string]string{
				"status":  "error",
				"message": "Ошибка чтения файла: " + err.Error(),
			})
			return
		}

		fileName = handler.Filename
	} else if strings.HasPrefix(contentType, "application/json") {
		var req struct {
			Filename string `json:"filename"`
			Data     string `json:"data"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.logger.Error().Err(err).Msg("Failed to decode JSON")
			wr.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(wr).Encode(map[string]string{
				"status":  "error",
				"message": "Неверный формат JSON: " + err.Error(),
			})
			return
		}

		fileData, err = base64.StdEncoding.DecodeString(req.Data)
		if err != nil {
			w.logger.Error().Err(err).Msg("Invalid base64 data")
			wr.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(wr).Encode(map[string]string{
				"status":  "error",
				"message": "Неверные данные Base64: " + err.Error(),
			})
			return
		}

		fileName = req.Filename
	} else {
		wr.WriteHeader(http.StatusUnsupportedMediaType)
		json.NewEncoder(wr).Encode(map[string]string{
			"status":  "error",
			"message": "Неподдерживаемый тип контента: " + contentType,
		})
		return
	}

	tempFileName := filepath.Join(w.uploadDir, time.Now().Format("20060102-150405")+"-"+fileName)

	if err := os.WriteFile(tempFileName, fileData, 0644); err != nil {
		w.logger.Error().Err(err).Str("file", tempFileName).Msg("Failed to save file")
		wr.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(wr).Encode(map[string]string{
			"status":  "error",
			"message": "Ошибка сохранения файла: " + err.Error(),
		})
		return
	}

	defer func() {
		if err := os.Remove(tempFileName); err != nil {
			w.logger.Error().Err(err).Str("file", tempFileName).Msg("Failed to delete temp file")
		}
	}()

	grpcReq := &api.UploadAndConvertExcelDataRequest{
		FileData: fileData,
		Filename: fileName,
	}

	ctx := context.Background()
	w.logger.Info().Msg("Calling gRPC service...")
	resp, err := w.client.UploadAndConvertExcelData(ctx, grpcReq)
	if err != nil {
		w.logger.Error().Err(err).Msg("gRPC call failed")
		wr.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(wr).Encode(map[string]string{
			"status":  "error",
			"message": "Ошибка обработки файла: " + err.Error(),
		})
		return
	}
	w.logger.Info().Str("job_id", resp.JobId).Msg("gRPC call successful")

	json.NewEncoder(wr).Encode(map[string]interface{}{
		"status":  "success",
		"message": fmt.Sprintf("Файл %s успешно обработан", fileName),
		"job_id":  resp.JobId,
	})
}

// Функции-заглушки для совместимости
func homeHandler(w http.ResponseWriter, r *http.Request)                   {}
func LoadTemplates() (*template.Template, error)                           { return nil, nil }
func RenderTemplate(w io.Writer, tmpl string, data TemplateData, web *Web) {}
