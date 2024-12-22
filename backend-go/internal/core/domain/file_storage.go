package domain

import "time"

type FileMetadata struct {
	Filename    string
	Size        int64
	ContentType string
	Bucket      string
	Path        string
	URL         string
	UploadedAt  time.Time
}
