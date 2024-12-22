package services

import (
	"github.com/chanasia/semantic-search-system/internal/core/domain"
	"github.com/chanasia/semantic-search-system/internal/core/ports"
)

type ModelService struct {
	modelManager    ports.ModelManager
	sessionManager  ports.SessionManager
	pipelineFactory ports.PipelineFactory
}

func NewEmbeddingModelService(
	modelManager ports.ModelManager,
	sessionManager ports.SessionManager,
	pipelineFactory ports.PipelineFactory,
) *ModelService {
	return &ModelService{
		modelManager:    modelManager,
		sessionManager:  sessionManager,
		pipelineFactory: pipelineFactory,
	}
}

func (s *ModelService) InitializeModelAndPipeline(modelConfig domain.EmbeddingModelConfig) (interface{}, error) {
	// Initialize session
	if err := s.sessionManager.Initialize(); err != nil {
		return nil, err
	}

	// Ensure model exists
	if err := s.modelManager.EnsureModelExists(modelConfig); err != nil {
		return nil, err
	}

	// Get model path
	modelPath := s.modelManager.GetModelPath(modelConfig)

	// Create pipeline
	pipeline, err := s.pipelineFactory.CreatePipeline(
		s.sessionManager.GetSession(),
		modelPath,
	)
	if err != nil {
		return nil, err
	}

	return pipeline, nil
}
