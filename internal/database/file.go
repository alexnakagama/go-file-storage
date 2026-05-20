package database

import "time"

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
