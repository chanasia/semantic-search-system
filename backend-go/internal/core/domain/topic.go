package domain

import (
	"time"

	"github.com/pgvector/pgvector-go"
	"gorm.io/gorm"
)

type Topic struct {
	ID              int32            `json:"id" gorm:"primaryKey"`
	Title           string           `json:"title" gorm:"notNull"`
	Context         string           `json:"context"  gorm:"notNull"`
	Page            string           `json:"page"`
	Tag             string           `json:"tag"`
	CreatedAt       time.Time        `json:"created_at" gorm:"notNull"`
	UpdatedAt       time.Time        `json:"updated_at" gorm:"notNull"`
	TopicImages     []TopicImage     `gorm:"foreignKey:TopicID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	TopicEmbeddings []TopicEmbedding `gorm:"foreignKey:TopicID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type TopicImage struct {
	ID        int32     `json:"id" gorm:"primaryKey"`
	TopicID   string    `json:"topic_id" gorm:"notNull"`
	ImagePath string    `json:"image_path" gorm:"notNull"`
	CreatedAt time.Time `json:"created_at" gorm:"notNull"`
	UpdatedAt time.Time `json:"updated_at" gorm:"notNull"`
}

type TopicEmbedding struct {
	ID        int32           `json:"id" gorm:"primaryKey"`
	TopicID   string          `json:"topic_id" gorm:"notNull"`
	Embedding pgvector.Vector `json:"embedding" gorm:"type:vector(768);notNull"`
	CreatedAt time.Time       `json:"created_at" gorm:"notNull"`
	UpdatedAt time.Time       `json:"updated_at" gorm:"notNull"`
}

type TopicWithSimilarity struct {
	Topic
	Similarity float64 `json:"similarity"`
}

type TopicRepository interface {
	Create(tx *gorm.DB, topic *Topic) error
	GetByID(id int64) (*Topic, error)
	GetImagePathsByTopicID(id int64) ([]TopicImage, error)
	Update(topic *Topic) error
	Delete(id int64) error
	List() ([]Topic, error)
	SearchBySimilarity(embedding []float32, limit int) ([]TopicWithSimilarity, error)
}

type TopicImageRepository interface {
	Create(tx *gorm.DB, topicImage *TopicImage) error
	GetByTopicID(id int64) ([]TopicImage, error)
}

type TopicEmbeddingRepository interface {
	Create(tx *gorm.DB, embedding *TopicEmbedding) error
}
