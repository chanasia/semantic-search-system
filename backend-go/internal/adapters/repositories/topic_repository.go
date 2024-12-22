package repositories

import (
	// "time"

	"github.com/chanasia/semantic-search-system/internal/core/domain"
	"gorm.io/gorm"
)

type topicRepository struct {
	db *gorm.DB
}

func NewTopicRepository(db *gorm.DB) domain.TopicRepository {
	return &topicRepository{db: db}
}

func (r *topicRepository) Create(tx *gorm.DB, topic *domain.Topic) error {
	return tx.Create(topic).Error
}

func (r *topicRepository) GetByID(id int64) (*domain.Topic, error) {
	var topic domain.Topic
	if err := r.db.First(&topic, id).Error; err != nil {
		return nil, err
	}
	return &topic, nil
}

func (r *topicRepository) Update(topic *domain.Topic) error {
	return r.db.Save(&topic).Error
}

func (r *topicRepository) Delete(id int64) error {
	return r.db.Delete(&domain.Topic{}, id).Error
}

func (r *topicRepository) List() ([]domain.Topic, error) {
	var topics []domain.Topic
	if err := r.db.Find(&topics).Error; err != nil {
		return nil, err
	}
	return topics, nil
}

func (r *topicRepository) GetImagePathsByTopicID(id int64) ([]domain.TopicImage, error) {
	var topicImages []domain.TopicImage

	if err := r.db.Where("id = ?", id).Find(&topicImages).Error; err != nil {
		return nil, err
	}
	return topicImages, nil
}
