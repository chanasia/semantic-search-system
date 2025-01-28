package repositories

import (
	// "time"

	"fmt"

	"github.com/chanasia/semantic-search-system/internal/core/domain"
	"github.com/pgvector/pgvector-go"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type topicRepository struct {
	db *gorm.DB
}

func NewTopicRepository(db *gorm.DB) domain.TopicRepository {
	return &topicRepository{db: db}
}

func (r *topicRepository) SearchBySimilarity(embedding []float32, limit int) ([]domain.TopicWithSimilarity, error) {
	var results []domain.TopicWithSimilarity

	// ใช้ GORM Joins และ Clauses สำหรับ similarity search
	err := r.db.Model(&domain.Topic{}).
		// เลือก fields ที่ต้องการและคำนวณ similarity
		Select("topics.*, 1 - (topic_embeddings.embedding <=> ?) as similarity", pgvector.NewVector(embedding)).
		Joins("JOIN topic_embeddings ON topics.id = topic_embeddings.topic_id").
		// เรียงลำดับตาม similarity จากมากไปน้อย
		Clauses(clause.OrderBy{
			Expression: clause.Expr{SQL: "similarity DESC"},
		}).
		Limit(limit).
		Find(&results).Error

	if err != nil {
		return nil, fmt.Errorf("failed to search by similarity: %w", err)
	}

	for i := range results {
		if err := r.db.Model(&domain.TopicImage{}).
			Where("topic_id = ?", results[i].ID).
			Find(&results[i].TopicImages).Error; err != nil {
			return nil, fmt.Errorf("failed to fetch topic images: %w", err)
		}
	}

	return results, nil
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
