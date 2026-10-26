package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

//go:embed index.html
var htmlContent embed.FS

type User struct {
	Username string
	Password string
	Role     string // "admin", "viewer", "uploader"
}

var users = map[string]User{
	"admin": {Username: "admin", Password: "admin123", Role: "admin"},
	"user2": {Username: "user2", Password: "user2pass", Role: "viewer"},
	"user3": {Username: "user3", Password: "user3pass", Role: "uploader"},
}

const uploadDir = "./uploads"

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok {
			http.Error(w, "Требуется авторизация", http.StatusUnauthorized)
			return
		}

		user, exists := users[username]
		if !exists || user.Password != password {
			http.Error(w, "Неверные учетные данные", http.StatusUnauthorized)
			return
		}

		r.Header.Set("X-User-Role", user.Role)
		r.Header.Set("X-Username", user.Username)
		next(w, r)
	}
}

func meHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}
	role := r.Header.Get("X-User-Role")
	username := r.Header.Get("X-Username")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"username": username,
		"role":     role,
	})
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	role := r.Header.Get("X-User-Role")
	if role != "admin" && role != "uploader" {
		http.Error(w, "Доступ запрещен: у вас нет прав на загрузку", http.StatusForbidden)
		return
	}

	err := r.ParseMultipartForm(10 << 20) // Лимит 10 MB
	if err != nil {
		http.Error(w, "Ошибка парсинга формы", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Ошибка получения файла", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// 1. Определяем тип файла по содержимому (первые 512 байт)
	buffer := make([]byte, 512)
	n, _ := file.Read(buffer)
	file.Seek(0, 0) // Сброс указателя в начало

	mimeType := http.DetectContentType(buffer[:n])

	// 2. Получаем расширение (с фоллбэком на исходное имя файла, если MIME не определен)
	ext := getExtension(header.Filename, mimeType)
	folderName := strings.TrimPrefix(ext, ".")
	if folderName == "" {
		folderName = "unknown"
	}

	targetDir := filepath.Join(uploadDir, folderName)

	// 3. Создаем папку, если её нет
	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		if err := os.MkdirAll(targetDir, os.ModePerm); err != nil {
			http.Error(w, "Ошибка создания директории", http.StatusInternalServerError)
			return
		}
	}

	// 4. Сохраняем файл с уникальным именем
	cleanName := strings.ReplaceAll(header.Filename, " ", "_")
	filename := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), cleanName, ext)
	targetPath := filepath.Join(targetDir, filename)

	dst, err := os.Create(targetPath)
	if err != nil {
		http.Error(w, "Ошибка создания файла на диске", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Ошибка сохранения файла", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Файл успешно загружен в папку: %s/%s", folderName, filename)
}

func getExtension(originalFilename string, mimeType string) string {
	// Пробуем получить расширение из MIME типа
	exts, err := mime.ExtensionsByType(mimeType)
	if err == nil && len(exts) > 0 {
		ext := exts[0]
		if ext == ".jpeg" {
			return ".jpg" // Нормализация для единообразия
		}
		return ext
	}

	// Фоллбэк: берем расширение из исходного имени файла, если MIME = application/octet-stream
	ext := filepath.Ext(originalFilename)
	if ext != "" {
		return ext
	}
	return ".unknown"
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	role := r.Header.Get("X-User-Role")
	if role != "admin" && role != "viewer" {
		http.Error(w, "Доступ запрещен: у вас нет прав на просмотр", http.StatusForbidden)
		return
	}

	entries, err := os.ReadDir(uploadDir)
	if err != nil {
		os.MkdirAll(uploadDir, os.ModePerm)
		entries, _ = os.ReadDir(uploadDir)
	}

	var result []map[string]interface{}
	for _, entry := range entries {
		if entry.IsDir() {
			dirPath := filepath.Join(uploadDir, entry.Name())
			files, _ := os.ReadDir(dirPath)

			var fileList []string
			for _, f := range files {
				if !f.IsDir() {
					fileList = append(fileList, f.Name()) // Возвращаем ТОЛЬКО имена, без путей для скачивания
				}
			}

			item := map[string]interface{}{
				"folder": entry.Name(),
				"files":  fileList,
			}
			result = append(result, item)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func main() {
	os.MkdirAll(uploadDir, os.ModePerm)

	// Раздаем встроенный HTML
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		file, _ := htmlContent.ReadFile("index.html")
		w.Write(file)
	})

	http.HandleFunc("/me", authMiddleware(meHandler))
	http.HandleFunc("/upload", authMiddleware(uploadHandler))
	http.HandleFunc("/list", authMiddleware(listHandler))

	fmt.Println("✅ Сервер запущен на порту 8080")
	fmt.Println("🌐 Откройте в браузере: http://localhost:8080")
	fmt.Println("👤 admin / admin123 (Загрузка + Просмотр)")
	fmt.Println("👤 user2 / user2pass (Только просмотр)")
	fmt.Println("👤 user3 / user3pass (Только загрузка)")

	log.Fatal(http.ListenAndServe(":8080", nil))
}