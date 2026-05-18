package database

import "time"

type Files struct {
	ID        int
	OwnerID   int
	FileName  string
	Size      int64
	MimeType  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
