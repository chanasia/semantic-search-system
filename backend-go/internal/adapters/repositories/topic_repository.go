package repositories

import (
	"time"

	"github.com/chanasia/semantic-search-system/internal/core/domain"
	"gorm.io/gorm"
)

type topicRepository struct {
	db *gorm.DB
}

type TopicModel struct {
	ID        int64 `gorm:"primaryKey"`
	Title     string
	Context   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewTopicRepository(db *gorm.DB) domain.TopicRepository {
	return &topicRepository{db: db}
}

func (r *topicRepository) Create(topic *domain.Topic) error {
	model := toModel(topic)
	return r.db.Create(&model).Error
}

func (r *topicRepository) GetByID(id int64) (*domain.Topic, error) {
	var model TopicModel
	if err := r.db.First(&model, id).Error; err != nil {
		return nil, err
	}
	return toDomain(&model), nil
}

func (r *topicRepository) Update(topic *domain.Topic) error {
	model := toModel(topic)
	return r.db.Save(&model).Error
}

func (r *topicRepository) Delete(id int64) error {
	return r.db.Delete(&TopicModel{}, id).Error
}

func (r *topicRepository) List() ([]domain.Topic, error) {
	var models []TopicModel
	if err := r.db.Find(&models).Error; err != nil {
		return nil, err
	}
	return toDomainList(models), nil
}

func toModel(topic *domain.Topic) *TopicModel {
	return &TopicModel{
		ID:        topic.ID,
		Title:     topic.Title,
		Context:   topic.Context,
		CreatedAt: topic.CreatedAt,
		UpdatedAt: topic.UpdatedAt,
	}
}

func toDomain(model *TopicModel) *domain.Topic {
	return &domain.Topic{
		ID:        model.ID,
		Title:     model.Title,
		Context:   model.Context,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}

func toDomainList(models []TopicModel) []domain.Topic {
	topics := make([]domain.Topic, len(models))
	for i, model := range models {
		topics[i] = *toDomain(&model)
	}
	return topics
}
