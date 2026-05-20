package database

import (
	"database/sql"
	"time"
)

type File struct {
	ID        int       `json:"id"`
	OwnerID   int       `json:"owner_id"`
	FileName  string    `json:"file_name"`
	Size      int64     `json:"size"`
	MimeType  string    `json:"mime_type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type FileResponse struct {
	ID        int       `json:"id"`
	OwnerID   int       `json:"owner_id"`
	FileName  string    `json:"file_name"`
	Size      int64     `json:"size"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (f *File) ToResponse() *FileResponse {
	return &FileResponse{
		ID:        f.ID,
		OwnerID:   f.OwnerID,
		FileName:  f.FileName,
		Size:      f.Size,
		CreatedAt: f.CreatedAt,
		UpdatedAt: f.UpdatedAt,
	}
}

func CreateFile(db *sql.DB, ownerID int, fileName string, size int64, mimeType string) (*File, error) {
	now := time.Now()
	var file File

	query := `
		INSERT INTO files (owner_id, file_name, size, mime_type, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, owner_id, file_name, size, mime_type, created_at, updated_at
	`
	err := db.QueryRow(query, ownerID, fileName, size, mimeType, now, now).
		Scan(&file.ID, &file.OwnerID, &file.FileName, &file.Size, &file.MimeType, &file.CreatedAt, &file.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &file, nil
}

func GetFileByOwner(db *sql.DB, ownerID int) ([]*File, error) {
	query := `
		SELECT id, owner_id, file_name, size, mime_type, created_at, updated_at 
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

		err := rows.Scan(&file.ID, &file.OwnerID, &file.FileName, &file.Size, &file.MimeType, &file.CreatedAt, &file.UpdatedAt)
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
		SELECT id, owner_id, file_name, size, mime_type, created_at, updated_at
		FROM files WHERE id = $1
	`
	err := db.QueryRow(query, fileID).Scan(
		&file.ID, &file.OwnerID, &file.FileName, &file.Size, &file.MimeType, &file.CreatedAt, &file.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &file, nil
}

func DeleteFile() {

}
