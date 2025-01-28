package main

import (
	"fmt"

	"github.com/knights-analytics/hugot"
	// "github.com/knights-analytics/hugot/pipelines"
)

func mains() {
	sessionModel, err := hugot.NewORTSession()

	if err != nil {
		panic(err.Error())
	}

	defer func(session *hugot.Session) {
		err := session.Destroy()
		if err != nil {
			panic(err.Error())
		}
	}(sessionModel)

	configDownloadModel := hugot.NewDownloadOptions()
	configDownloadModel.Verbose = true
	configDownloadModel.SkipSha = true

	modelPath, err := hugot.DownloadModel("sentence-transformers/paraphrase-multilingual-mpnet-base-v2", "./lib/model/", configDownloadModel)

	if err != nil {
		panic(err.Error())
	}

	configModel := hugot.FeatureExtractionConfig{
		// ModelPath:    "lib/model/sentence-transformers_paraphrase-multilingual-mpnet-base-v2",
		ModelPath:    modelPath,
		Name:         "embeddingPipeline",
		OnnxFilename: "model_qint8_avx512.onnx",
	}

	sentenceTransformersPipeline, err := hugot.NewPipeline(sessionModel, configModel)

	if err != nil {
		fmt.Println(sentenceTransformersPipeline)
		panic(err.Error())
	}

	resultEmbeddings, err := sentenceTransformersPipeline.RunPipeline([]string{"Hello my friend!"})

	if err != nil && len(resultEmbeddings.Embeddings) == 0 {
		panic(err.Error())
	}

	var embedding = resultEmbeddings.Embeddings[0]
	fmt.Println(embedding)

}
