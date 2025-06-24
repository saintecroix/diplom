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
	var fileData []byte
	var fileName string
	var err error

	// Определяем тип контента
	contentType := r.Header.Get("Content-Type")

	// Обработка multipart/form-data (новый фронтенд)
	if strings.HasPrefix(contentType, "multipart/form-data") {
		// Парсим форму с ограничением размера 10MB
		err = r.ParseMultipartForm(10 << 20)
		if err != nil {
			w.logger.Error().Err(err).Msg("Error parsing multipart form")
			http.Error(wr, "Error parsing form", http.StatusBadRequest)
			return
		}

		// Получаем файл из формы
		file, handler, err := r.FormFile("file")
		if err != nil {
			w.logger.Error().Err(err).Msg("Error retrieving file from form")
			http.Error(wr, "Error retrieving file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Читаем содержимое файла
		fileData, err = io.ReadAll(file)
		if err != nil {
			w.logger.Error().Err(err).Msg("Error reading file")
			http.Error(wr, "Error reading file", http.StatusInternalServerError)
			return
		}

		fileName = handler.Filename
	} else if strings.HasPrefix(contentType, "application/json") {
		var req struct {
			Filename string `json:"filename"`
			Data     string `json:"data"` // base64 строка
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.logger.Error().Err(err).Msg("Failed to decode JSON")
			http.Error(wr, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Декодируем base64
		fileData, err = base64.StdEncoding.DecodeString(req.Data)
		if err != nil {
			w.logger.Error().Err(err).Msg("Invalid base64 data")
			http.Error(wr, "Invalid file data", http.StatusBadRequest)
			return
		}

		fileName = req.Filename
	} else {
		http.Error(wr, "Unsupported media type", http.StatusUnsupportedMediaType)
		return
	}

	// Создаем уникальное имя файла
	tempFileName := filepath.Join(w.uploadDir, time.Now().Format("20060102-150405")+"-"+fileName)

	// Сохраняем файл на диск
	if err := os.WriteFile(tempFileName, fileData, 0644); err != nil {
		w.logger.Error().Err(err).Str("file", tempFileName).Msg("Failed to save file")
		http.Error(wr, "Server error", http.StatusInternalServerError)
		return
	}

	// Гарантируем удаление файла
	defer func() {
		if err := os.Remove(tempFileName); err != nil {
			w.logger.Error().Err(err).Str("file", tempFileName).Msg("Failed to delete temp file")
		}
	}()

	// Создаем gRPC запрос
	grpcReq := &api.UploadAndConvertExcelDataRequest{
		FileData: fileData,
		Filename: fileName,
	}

	// Отправляем запрос в inputConvert сервис
	ctx := context.Background()
	resp, err := w.client.UploadAndConvertExcelData(ctx, grpcReq)
	if err != nil {
		w.logger.Error().Err(err).Msg("Error calling ConvertExcelData")
		http.Error(wr, "Failed to convert data", http.StatusInternalServerError)
		return
	}

	// Обрабатываем ответ от inputConvert сервиса
	w.logger.Info().Interface("response", resp).Msg("Received response from gRPC service")

	// Формируем JSON-ответ
	wr.Header().Set("Content-Type", "application/json")
	wr.WriteHeader(http.StatusOK)

	// Для нового фронтенда возвращаем упрощенный ответ
	if strings.HasPrefix(contentType, "multipart/form-data") {
		json.NewEncoder(wr).Encode(map[string]interface{}{
			"status":  "success",
			"message": fmt.Sprintf("Файл %s успешно обработан", fileName),
		})
	} else {
		// Для старого формата возвращаем полный ответ
		json.NewEncoder(wr).Encode(map[string]interface{}{
			"status":  "success",
			"message": "File processed and converted",
			"details": resp,
		})
	}
}

// Функции-заглушки для совместимости
func homeHandler(w http.ResponseWriter, r *http.Request)                   {}
func LoadTemplates() (*template.Template, error)                           { return nil, nil }
func RenderTemplate(w io.Writer, tmpl string, data TemplateData, web *Web) {}
