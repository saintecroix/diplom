package web

import (
	"context"
	"embed"
	"fmt"
	"github.com/rs/zerolog/log"
	api "github.com/saintecroix/diplom/cmd/webUI/internal/api"
	"html/template"
	"io"
	"net/http"
	"os"

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

//go:embed templates/*
var templateFS embed.FS

// web структура для хранения зависимостей обработчиков
type web struct {
	client    api.InputConvertServiceClient
	templates *template.Template
}

var w *web

// LoadTemplates загружает все шаблоны из встроенной файловой системы
func LoadTemplates() (*template.Template, error) {
	tmpl := template.New("") // Создаем новый набор шаблонов

	// Используем ParseFS для загрузки шаблонов из embed.FS
	tmpl, err := tmpl.ParseFS(templateFS, "templates/*.html")
	if err != nil {
		return nil, err
	}

	return tmpl, nil
}

// RenderTemplate выполняет указанный шаблон с данными и отправляет результат в ResponseWriter
func RenderTemplate(wr http.ResponseWriter, name string, data TemplateData, w *web) {
	err := w.templates.ExecuteTemplate(wr, "base.html", data)
	if err != nil {
		log.Printf("Error rendering template: %v", err)
		http.Error(wr, "Internal Server Error", 500)
	}
}

func RegisterHandlers(mux *http.ServeMux) {
	var err error
	templates, err := LoadTemplates()
	if err != nil {
		panic(err)
	}

	client, err := NewInputConvertClient()
	if err != nil {
		panic(err)
	}

	w = &web{
		client:    client,
		templates: templates,
	}

	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/trips/upload", uploadHandler) // Добавляем обработчик для /trips/upload
	// TODO: Add other API endpoints here
}

func homeHandler(wr http.ResponseWriter, r *http.Request) {
	data := TemplateData{
		Title:   "Home Page",
		Header:  "Welcome to the WebUI!",
		Year:    2024,
		Content: "This is the home page content.",
	}
	RenderTemplate(wr, "index.html", data, w)
}

// NewInputConvertClient создает gRPC клиент для inputConvert сервиса.
func NewInputConvertClient() (api.InputConvertServiceClient, error) {
	inputConvertAddress := os.Getenv("INPUT_CONVERT_ADDRESS")
	if inputConvertAddress == "" {
		inputConvertAddress = "localhost:50051" // Default address
		log.Printf("INPUT_CONVERT_ADDRESS not set, using default:", inputConvertAddress)
	}

	conn, err := grpc.Dial(inputConvertAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to dial inputConvert service: %w", err)
	}

	return api.NewInputConvertServiceClient(conn), nil
}

func uploadHandler(wr http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(wr, "Method not allowed", 405)
		return
	}

	// 1. Получаем файл из запроса
	file, header, err := r.FormFile("file")
	if err != nil {
		log.Printf("Error getting file from request: %v", err)
		http.Error(wr, "Invalid file upload", 400)
		return
	}
	defer file.Close()

	log.Printf("Uploaded file: %s, size: %d bytes", header.Filename, header.Size)

	// 2. Читаем содержимое файла в байтовый массив
	fileData, err := io.ReadAll(file)
	if err != nil {
		log.Printf("Error reading file: %v", err)
		http.Error(wr, "Failed to read file", 500)
		return
	}

	// 3. Создаем временный файл
	tmpFile, err := os.CreateTemp("", "upload-*.xlsx")
	if err != nil {
		log.Printf("Error creating temp file: %v", err)
		http.Error(wr, "Failed to create temp file", 500)
		return
	}
	defer os.Remove(tmpFile.Name()) // Clean up the temp file
	defer tmpFile.Close()

	// 4. Записываем содержимое загруженного файла во временный файл
	_, err = tmpFile.Write(fileData)
	if err != nil {
		log.Printf("Error writing to temp file: %v", err)
		http.Error(wr, "Failed to write to temp file", 500)
		return
	}

	// 5. Создаем gRPC запрос
	req := &api.ConvertExcelDataRequest{
		FileData: fileData,
		Filename: header.Filename,
	}

	// 6. Отправляем запрос в inputConvert сервис
	ctx := context.Background()
	resp, err := w.client.ConvertExcelData(ctx, req)
	if err != nil {
		log.Printf("Error calling ConvertExcelData: %v", err)
		http.Error(wr, "Failed to convert data", 500)
		return
	}

	// 7. Обрабатываем ответ от inputConvert сервиса
	log.Printf("Received response: %v", resp.Results)

	// 8. Формируем сообщение для пользователя
	message := fmt.Sprintf("Conversion results: %v, Error: %s", resp.Results, resp.Error)
	if resp.Error != "" {
		log.Printf("Error from inputConvert: %s", resp.Error)
	}

	// 9. Отправляем ответ клиенту
	data := TemplateData{
		Title:   "Upload Result",
		Header:  "Upload Result",
		Year:    2024,
		Content: message, // Отображаем результаты от inputConvert
	}
	RenderTemplate(wr, "index.html", data, w)
}
