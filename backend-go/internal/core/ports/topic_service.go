package ports

import (
	"github.com/chanasia/semantic-search-system/internal/core/domain"
)

type TopicService interface {
	CreateTopic(topic *domain.Topic) error
	GetTopic(id int64) (*domain.Topic, error)
	UpdateTopic(topic **domain.Topic) error
	DeleteTopic(id int64) error
	GetTopics() ([]domain.Topic, error)
}
