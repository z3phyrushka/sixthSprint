package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/service"
)

func HandleMain(w http.ResponseWriter, r *http.Request) {
	possiblePaths := []string{
		"index.html",
		"../index.html",
	}

	var data []byte
	var err error

	for _, path := range possiblePaths {
		data, err = os.ReadFile(path)
		if err == nil {
			log.Printf("Успешно загружен index.html из: %s", path)
			break
		}
	}

	if err != nil {
		http.Error(w, "index.html не найден", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(data)
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Ошибка при парсинге формы", http.StatusInternalServerError)
		return
	}

	file, header, err := r.FormFile("myFile")
	if err != nil {
		http.Error(w, "Ошибка при получении файла", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Ошибка при чтении файла", http.StatusInternalServerError)
		return
	}

	original := string(data)
	converted, err := service.AutoConvert(string(data))
	if err != nil {
		http.Error(w, "Ошибка при конвертации", http.StatusInternalServerError)
		return
	}

	newFileName := fmt.Sprintf("%s%s", time.Now().UTC().Format("20060102_150405"), filepath.Ext(header.Filename))
	newFile, err := os.Create(newFileName)
	if err != nil {
		http.Error(w, "Ошибка при создании файла", http.StatusInternalServerError)
		return
	}
	defer newFile.Close()

	if _, err := newFile.WriteString(converted); err != nil {
		http.Error(w, "Ошибка при записи файла", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	response := "Исходный текст:\n" + original + "\n\n" +
		"Результат:\n" + converted
	_, _ = w.Write([]byte(response))
}
