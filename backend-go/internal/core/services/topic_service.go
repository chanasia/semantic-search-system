package services

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"path"
	"path/filepath"
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

func (s *topicService) CreateTopic(ctx context.Context, topic *domain.Topic, files []*multipart.FileHeader, imagePaths []string) error {
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

	// 4. Handle image paths
	if len(imagePaths) > 0 {
		for _, imagePath := range imagePaths {
			minioPath := fmt.Sprintf("images/%s", imagePath)

			topicImage := &domain.TopicImage{
				TopicID:   fmt.Sprintf("%d", topic.ID),
				ImagePath: minioPath,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			if err := s.topicImageRepo.Create(tx, topicImage); err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to create topic image record: %w", err)
			}

			topic.TopicImages = append(topic.TopicImages, *topicImage)
		}
	} else if len(files) > 0 {
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

func (s *topicService) SearchTopics(ctx context.Context, searchText string, limit int) ([]domain.TopicWithSimilarity, error) {
	// Generate embedding สำหรับ search text
	embedding, err := s.embeddingAdapter.GenerateEmbedding(searchText)
	if err != nil {
		return nil, fmt.Errorf("failed to generate embedding: %w", err)
	}

	// ค้นหา topics ด้วย similarity
	results, err := s.topicRepo.SearchBySimilarity(*embedding, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search topics: %w", err)
	}

	// แปลง MinIO path เป็น URL สำหรับแต่ละ topic
	for i := range results {
		for j := range results[i].TopicImages {
			// เปลี่ยน ImagePath เป็น URL
			results[i].TopicImages[j].ImagePath = fmt.Sprintf("/api/v1/images/%d", results[i].TopicImages[j].ID)
		}
	}

	return results, nil
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

func (s *topicService) GetTopicImage(ctx context.Context, imageID int32) (io.Reader, string, error) {
	// Find the image record
	var topicImage domain.TopicImage
	if err := s.db.Where("id = ?", imageID).First(&topicImage).Error; err != nil {
		return nil, "", fmt.Errorf("image not found: %w", err)
	}

	// Get content type based on file extension
	contentType := getContentTypeFromPath(topicImage.ImagePath)

	// Get object from MinIO
	obj, err := s.minioAdapter.GetObject(ctx, topicImage.ImagePath)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get image from storage: %w", err)
	}

	return obj, contentType, nil
}

// Add helper function to determine content type
func getContentTypeFromPath(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	case ".svg":
		return "image/svg+xml"
	default:
		return "image/jpeg" // Default to JPEG if unknown
	}
}
