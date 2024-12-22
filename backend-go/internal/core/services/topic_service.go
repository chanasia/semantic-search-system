package services

import (
	"context"
	"fmt"
	"mime/multipart"
	"path"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/chanasia/semantic-search-system/internal/adapters/hugot"
	"github.com/chanasia/semantic-search-system/internal/adapters/minio"
	"github.com/chanasia/semantic-search-system/internal/core/domain"
	"github.com/pgvector/pgvector-go"
)

type topicService struct {
	db                 *gorm.DB
	topicRepo          domain.TopicRepository
	topicImageRepo     domain.TopicImageRepository
	topicEmbeddingRepo domain.TopicEmbeddingRepository
	minioAdapter       *minio.MinioAdapter
	embeddingAdapter   *hugot.HugotEmbeddingAdapter
}

func NewTopicService(
	db *gorm.DB,
	topicRepo domain.TopicRepository,
	topicImageRepo domain.TopicImageRepository,
	topicEmbeddingRepo domain.TopicEmbeddingRepository,
	minioAdapter *minio.MinioAdapter,
	embeddingAdapter *hugot.HugotEmbeddingAdapter,
) *topicService {
	return &topicService{
		db:                 db,
		topicRepo:          topicRepo,
		topicImageRepo:     topicImageRepo,
		topicEmbeddingRepo: topicEmbeddingRepo,
		minioAdapter:       minioAdapter,
		embeddingAdapter:   embeddingAdapter,
	}
}

func (s *topicService) CreateTopic(ctx context.Context, topic *domain.Topic, files []*multipart.FileHeader) error {
	tx := s.db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. Create topic
	if err := s.topicRepo.Create(tx, topic); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create topic: %w", err)
	}

	// 2. Generate embedding จาก title
	embedding, err := s.embeddingAdapter.GenerateEmbedding(topic.Title)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to generate embedding: %w", err)
	}

	// 3. Create topic_embedding
	topicEmbedding := &domain.TopicEmbedding{
		TopicID:   fmt.Sprintf("%d", topic.ID),
		Embedding: pgvector.NewVector(*embedding),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.topicEmbeddingRepo.Create(tx, topicEmbedding); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create topic embedding: %w", err)
	}

	// 4. Handle images
	for _, file := range files {
		src, err := file.Open()
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to open file: %w", err)
		}
		defer src.Close()

		filename := fmt.Sprintf("%d-%d-%s", topic.ID, time.Now().UnixNano(), file.Filename)
		objectPath := fmt.Sprintf("images/%s", filename)

		contentType := file.Header.Get("Content-Type")
		if contentType == "" {
			contentType = getContentTypeFromFileName(file.Filename)
		}

		metadata, err := s.minioAdapter.UploadFile(
			ctx,
			objectPath,
			src,
			file.Size,
			contentType,
		)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to upload image: %w", err)
		}

		topicImage := &domain.TopicImage{
			TopicID:   fmt.Sprintf("%d", topic.ID),
			ImagePath: metadata.Path,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := s.topicImageRepo.Create(tx, topicImage); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to create topic image record: %w", err)
		}

		topic.TopicImages = append(topic.TopicImages, *topicImage)
	}

	// 5. Commit transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func getContentTypeFromFileName(filename string) string {
	ext := strings.ToLower(path.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	default:
		return "application/octet-stream"
	}
}
func (s *topicService) GetTopic(id int64) (*domain.Topic, error) {
	return s.topicRepo.GetByID(id)
}

func (s *topicService) UpdateTopic(topic **domain.Topic) error {
	return s.topicRepo.Update(*topic)
}

func (s *topicService) DeleteTopic(id int64) error {
	return s.topicRepo.Delete(id)
}

func (s *topicService) GetTopics() ([]domain.Topic, error) {
	return s.topicRepo.List()
}

func (s *topicService) GetTopicImages(id int64) ([]domain.TopicImage, error) {
	return s.topicRepo.GetImagePathsByTopicID(id)
}
