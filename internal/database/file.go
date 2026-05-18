package database

import "time"

type Files struct {
	ID        int       `json:"id"`
	OwnerID   int       `json:"owner_id"`
	FileName  string    `json:"file_name"`
	Size      int64     `json:"size"`
	MimeType  string    `json:"mime_type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
