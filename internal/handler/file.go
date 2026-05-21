package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/alexnakagama/go-file-storage/internal/database"
)

type FileResponse struct {
	ID          int       `json:"id"`
	OwnerID     int       `json:"owner_id"`
	FileName    string    `json:"file_name"`
	StoragePath string    `json:"storage_path"`
	Size        int64     `json:"size"`
	MimeType    string    `json:"mime_type"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Message     string    `json:"message"`
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

		ctx := r.Context()
		userIDVal := ctx.Value("userID")
		if userIDVal == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		ownerID, ok := userIDVal.(int)
		if !ok {
			http.Error(w, "Invalid userID", http.StatusInternalServerError)
			return
		}

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

		dbFile, err := database.CreateFile(db, ownerID, header.Filename, storedFileName, size, mimeType)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		resp := FileResponse{
			ID:          dbFile.ID,
			OwnerID:     dbFile.OwnerID,
			FileName:    dbFile.FileName,
			StoragePath: dbFile.StoragePath,
			Size:        dbFile.Size,
			MimeType:    dbFile.MimeType,
			CreatedAt:   dbFile.CreatedAt,
			UpdatedAt:   dbFile.UpdatedAt,
			Message:     "File uploaded successfully",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

func SearchFileByOwnerHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userIDVal := ctx.Value("userID")
		if userIDVal == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ownerID, ok := userIDVal.(int)
		if !ok {
			http.Error(w, "Invalid userID", http.StatusInternalServerError)
			return
		}

		files, err := database.GetFileByOwner(db, ownerID)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		fileList := make([]FileResponse, 0, len(files))
		for _, f := range files {
			fileList = append(fileList, FileResponse{
				ID:          f.ID,
				OwnerID:     f.OwnerID,
				FileName:    f.FileName,
				StoragePath: f.StoragePath,
				Size:        f.Size,
				MimeType:    f.MimeType,
				CreatedAt:   f.CreatedAt,
				UpdatedAt:   f.UpdatedAt,
				Message:     "",
			})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(fileList)
	}
}

func DownloadFileHandler(db *sql.DB, uploadDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userIDVal := ctx.Value("userID")
		if userIDVal == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ownerID, ok := userIDVal.(int)
		if !ok {
			http.Error(w, "Invalid userID", http.StatusInternalServerError)
			return
		}

		fileIDStr := r.URL.Query().Get("file_id")
		if fileIDStr == "" {
			http.Error(w, "file_id is required", http.StatusBadRequest)
			return
		}

		fileID, err := strconv.Atoi(fileIDStr)
		if err != nil {
			http.Error(w, "Invalid file_id", http.StatusBadRequest)
			return
		}

		dbFile, err := database.GetFileByIDAndOwner(db, fileID, ownerID)
		if err != nil {
			http.Error(w, "fileID is required", http.StatusInternalServerError)
			return
		}

		filePath := filepath.Join(uploadDir, dbFile.StoragePath)
		file, err := os.Open(filePath)
		if err != nil {
			http.Error(w, "Could not open file", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", dbFile.FileName))
		w.Header().Set("Content-Type", dbFile.MimeType)
		w.Header().Set("Content-Lenght", fmt.Sprintf("&d", dbFile.Size))
		io.Copy(w, file)
	}
}
