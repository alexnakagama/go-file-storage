package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/alexnakagama/go-file-storage/internal/database"
)

type FileResponse struct {
	ID       int    `json:"id"`
	OwnerID  int    `json:"owner_id"`
	FileName string `json:"file_name"`
	Size     int64  `json:"size"`
	MimeType string `json:"mime_type"`
	Message  string `json:"message"`
}

func AddFileHandler(db *sql.DB, uploadDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			http.Error(w, "Could not parse multiform part", http.StatusBadRequest)
			return
		}

		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "No file uploaded", http.StatusBadRequest)
			return
		}
		defer file.Close()

		ownerID := 1

		storedFileName := fmt.Sprintf("%d_%d_%s", ownerID, time.Now().UnixNano(), header.Filename)
		dstPath := filepath.Join(uploadDir, storedFileName)
		dst, err := os.Create(dstPath)
		if err != nil {
			http.Error(w, "Could not store file", http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		size, err := io.Copy(dst, file)
		if err != nil {
			http.Error(w, "Could not save file", http.StatusInternalServerError)
			return
		}

		mimeType := header.Header.Get("Content-Type")

		dbFile, err := database.CreateFile(db, ownerID, header.Filename, size, mimeType)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		resp := FileResponse{
			ID:       dbFile.ID,
			OwnerID:  dbFile.OwnerID,
			FileName: dbFile.FileName,
			Size:     dbFile.Size,
			MimeType: dbFile.MimeType,
			Message:  "File uploaded successfully",
		}

		w.Header().Set("Content-Type", "application-json")
		json.NewEncoder(w).Encode(resp)
	}
}
