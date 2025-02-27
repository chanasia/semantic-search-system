package domain

import "time"

type TopicImages struct {
	ID        int32     `db:"id"`
	TopicId   string    `db:"topic_id"`
	ImagePath string    `db:"image_path"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
