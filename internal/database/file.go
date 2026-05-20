package database

import (
	"database/sql"
	"fmt"
	"time"
)

type File struct {
	ID          int       `json:"id"`
	OwnerID     int       `json:"owner_id"`
	FileName    string    `json:"file_name"`
	StoragePath string    `json:"storage_path"`
	Size        int64     `json:"size"`
	MimeType    string    `json:"mime_type"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func CreateFile(db *sql.DB, ownerID int, fileName string, storagePath string, size int64, mimeType string) (*File, error) {
	now := time.Now()
	var file File

	query := `
		INSERT INTO files (owner_id, file_name, storage_path, size, mime_type, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, owner_id, file_name, storage_path, size, mime_type, created_at, updated_at
	`
	err := db.QueryRow(query, ownerID, fileName, storagePath, size, mimeType, now, now).
		Scan(&file.ID, &file.OwnerID, &file.FileName, &file.StoragePath, &file.Size, &file.MimeType, &file.CreatedAt, &file.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &file, nil
}

func GetFileByOwner(db *sql.DB, ownerID int) ([]*File, error) {
	query := `
		SELECT id, owner_id, file_name, storage_path, size, mime_type, created_at, updated_at 
		FROM files WHERE owner_id = $1
	`
	rows, err := db.Query(query, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []*File
	for rows.Next() {
		var file File

		err := rows.Scan(&file.ID, &file.OwnerID, &file.FileName, &file.StoragePath, &file.Size, &file.MimeType, &file.CreatedAt, &file.UpdatedAt)
		if err != nil {
			return nil, err
		}

		files = append(files, &file)
	}
	return files, nil
}

func GetFileByID(db *sql.DB, fileID int) (*File, error) {
	var file File
	query := `
		SELECT id, owner_id, file_name, storage_path, size, mime_type, created_at, updated_at
		FROM files WHERE id = $1
	`
	err := db.QueryRow(query, fileID).Scan(
		&file.ID, &file.OwnerID, &file.FileName, &file.StoragePath, &file.Size, &file.MimeType, &file.CreatedAt, &file.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &file, nil
}

func DeleteFile(db *sql.DB, fileID int) (*File, error) {
	var file File
	query := `
		DELETE FROM files WHERE id=$1
		RETURNING id, owner_id, file_name, storage_path, size, mime_type, created_at, updated_at
	`

	err := db.QueryRow(query, fileID).
		Scan(&file.ID, &file.OwnerID, &file.FileName, &file.StoragePath, &file.Size, &file.MimeType, &file.CreatedAt, &file.UpdatedAt)

	if err != nil {
		return nil, err
	}

	fmt.Println("File has been deleted")
	return &file, nil
}

func UpdateFileName(db *sql.DB, fileID int, newFileName string) (*File, error) {
	var file File
	query := `
		UPDATE files
		SET file_name = $1, updated_at = $2
		WHERE id = $3
		RETURNING id, owner_id, file_name, storage_path, size, mime_type, created_at, updated_at
	`
	now := time.Now()

	err := db.QueryRow(query, newFileName, now, fileID).Scan(
		&file.ID, &file.OwnerID, &file.FileName, &file.StoragePath, &file.Size, &file.MimeType, &file.CreatedAt, &file.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &file, err
}
