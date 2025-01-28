package hugot

import (
	"os"
	"path"
	"strings"

	"github.com/chanasia/semantic-search-system/internal/core/domain"
	"github.com/knights-analytics/hugot"
)

type HugotModelManager struct {
	basePath string
}

func NewHugotModelManager(basePath string) *HugotModelManager {
	return &HugotModelManager{
		basePath: basePath,
	}
}

func (m *HugotModelManager) EnsureModelExists(modelConfig domain.EmbeddingModelConfig) error {
	// Check if model exists
	folders, err := os.ReadDir(m.basePath)
	if err != nil {
		return err
	}

	modelDirName := strings.Replace(modelConfig.HuggingFaceName, "/", "_", -1)

	for _, folder := range folders {
		if folder.Name() == modelDirName {
			return nil // Model already exists
		}
	}

	// Download model if it doesn't exist
	options := hugot.NewDownloadOptions()
	options.Verbose = true
	options.SkipSha = true

	_, err = hugot.DownloadModel(modelConfig.HuggingFaceName, m.basePath, options)
	return err
}

func (m *HugotModelManager) GetModelPath(modelConfig domain.EmbeddingModelConfig) string {
	modelDirName := strings.Replace(modelConfig.HuggingFaceName, "/", "_", -1)
	return path.Join(m.basePath, modelDirName)
}
