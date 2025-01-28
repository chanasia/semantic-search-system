package repositories

import (
	"github.com/chanasia/semantic-search-system/internal/core/domain"
	"gorm.io/gorm"
)

type topicImageRepository struct {
	db *gorm.DB
}

func NewTopicImageRepository(db *gorm.DB) domain.TopicImageRepository {
	return &topicImageRepository{db: db}
}

func (r *topicImageRepository) Create(tx *gorm.DB, topicImage *domain.TopicImage) error {
	return tx.Create(topicImage).Error
}

func (r *topicImageRepository) GetByTopicID(id int64) ([]domain.TopicImage, error) {
	var images []domain.TopicImage
	err := r.db.Where("topic_id = ?", id).Find(&images).Error
	return images, err
}
