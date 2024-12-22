package repositories

import (
	"github.com/chanasia/semantic-search-system/internal/core/domain"
	"gorm.io/gorm"
)

type topicEmbeddingRepository struct {
	db *gorm.DB
}

func NewTopicEmbeddingRepository(db *gorm.DB) domain.TopicEmbeddingRepository {
	return &topicEmbeddingRepository{db: db}
}

func (r *topicEmbeddingRepository) Create(tx *gorm.DB, embedding *domain.TopicEmbedding) error {
	return tx.Create(embedding).Error
}
