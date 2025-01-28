package ports

import (
	"context"
	"io"
	"mime/multipart"

	"github.com/chanasia/semantic-search-system/internal/core/domain"
)

type TopicService interface {
	CreateTopic(ctx context.Context, topic *domain.Topic, files []*multipart.FileHeader, imagePaths []string) error
	SearchTopics(ctx context.Context, searchText string, limit int) ([]domain.TopicWithSimilarity, error)
	GetTopic(id int64) (*domain.Topic, error)
	UpdateTopic(topic **domain.Topic) error
	DeleteTopic(id int64) error
	GetTopics() ([]domain.Topic, error)
	GetTopicImage(ctx context.Context, imageID int32) (io.Reader, string, error)
}
