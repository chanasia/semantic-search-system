package ports

import (
	"github.com/chanasia/semantic-search-system/internal/core/domain"
)

type EmbeddingService interface {
	GenerateEmbedding(text string) (*[]float32, error)
	GenerateEmbeddings(texts []string) (*[][]float32, error)
}

type ModelDownloader interface {
	DownloadModel(modelName string, destination string) (string, error)
}

type ModelManager interface {
	EnsureModelExists(config domain.EmbeddingModelConfig) error
	GetModelPath(config domain.EmbeddingModelConfig) string
}

type SessionManager interface {
	Initialize() error
	Destroy() error
	GetSession() interface{}
}

type PipelineFactory interface {
	CreatePipeline(session interface{}, modelPath string) (interface{}, error)
}
