package services

import (
	"github.com/chanasia/semantic-search-system/internal/core/domain"
)

type topicService struct {
	repo domain.TopicRepository
}

func NewTopicService(repo domain.TopicRepository) *topicService {
	return &topicService{repo: repo}
}

func (s *topicService) CreateTopic(topic *domain.Topic) error {
	return s.repo.Create(topic)
}

func (s *topicService) GetTopic(id int64) (*domain.Topic, error) {
	return s.repo.GetByID(id)
}

func (s *topicService) UpdateTopic(topic **domain.Topic) error {
	return s.repo.Update(*topic)
}

func (s *topicService) DeleteTopic(id int64) error {
	return s.repo.Delete(id)
}

func (s *topicService) GetTopics() ([]domain.Topic, error) {
	return s.repo.List()
}
