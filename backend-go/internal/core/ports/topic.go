package ports

import (
	"context"
	"mime/multipart"

	"github.com/chanasia/semantic-search-system/internal/core/domain"
)

type TopicService interface {
	CreateTopic(ctx context.Context, topic *domain.Topic, files []*multipart.FileHeader) error
	GetTopic(id int64) (*domain.Topic, error)
	UpdateTopic(topic **domain.Topic) error
	DeleteTopic(id int64) error
	GetTopics() ([]domain.Topic, error)
}
