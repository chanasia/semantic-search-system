package hugot

import (
	"fmt"

	"github.com/knights-analytics/hugot"
)

type HugotPipelineFactory struct{}

func NewHugotPipelineFactory() *HugotPipelineFactory {
	return &HugotPipelineFactory{}
}

func (f *HugotPipelineFactory) CreatePipeline(session interface{}, modelPath string) (interface{}, error) {
	hugotSession, ok := session.(*hugot.Session)
	if !ok {
		return nil, fmt.Errorf("invalid session type")
	}

	config := hugot.FeatureExtractionConfig{
		ModelPath:    modelPath,
		Name:         "embeddingPipeline",
		OnnxFilename: "model_qint8_avx512.onnx",
	}

	return hugot.NewPipeline(hugotSession, config)
}
