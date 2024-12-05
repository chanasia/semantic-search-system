package domain

import "time"

type Topic struct {
	ID        int64
	Title     string
	Context   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type TopicRepository interface {
	Create(topic *Topic) error
	GetByID(id int64) (*Topic, error)
	Update(topic *Topic) error
	Delete(id int64) error
	List() ([]Topic, error)
}
