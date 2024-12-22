package main

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/chanasia/semantic-search-system/internal/adapters/handlers"
	"github.com/chanasia/semantic-search-system/internal/adapters/hugot"
	"github.com/chanasia/semantic-search-system/internal/adapters/minio"
	"github.com/chanasia/semantic-search-system/internal/adapters/repositories"
	"github.com/chanasia/semantic-search-system/internal/config"
	"github.com/chanasia/semantic-search-system/internal/core/domain"
	"github.com/chanasia/semantic-search-system/internal/core/services"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	// Connnect to db
	DBHost := config.GetEnv("DB_HOST", "localhost")
	DBPort := config.GetEnv("DB_PORT", "5432")
	DBUser := config.GetEnv("DB_USER", "postgres")
	DBPassword := config.GetEnv("DB_PASSWORD", "postgres")
	DBName := config.GetEnv("DB_NAME", "projectdb")
	DBSchema := config.GetEnv("DB_SCHEMA", "public")
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable search_path=%s", DBHost, DBPort, DBUser, DBPassword, DBName, DBSchema)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// Initialize pgvector
	db.Exec("CREATE EXTENSION IF NOT EXISTS vector")

	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	modelPath := path.Join(currentDir, "lib", "models")
	modelConfig := domain.EmbeddingModelConfig{
		ModelName:         "embeddingPipeline",
		ModelBasePath:     modelPath,
		HuggingFaceName:   "sentence-transformers/paraphrase-multilingual-mpnet-base-v2",
		OnnxModelFilename: "model_qint8_avx512.onnx",
	}

	// Initialize Hugot
	sessionManager := hugot.NewHugotSessionManager()
	modelManager := hugot.NewHugotModelManager(modelPath)
	pipelineFactory := hugot.NewHugotPipelineFactory()

	// Create service
	modelService := services.NewEmbeddingModelService(modelManager, sessionManager, pipelineFactory)

	// Initialize pipeline
	pipeline, err := modelService.InitializeModelAndPipeline(modelConfig)
	if err != nil {
		log.Fatal(err)
	}

	// Create embedding adapter
	embeddingAdapter, err := hugot.NewHugotEmbeddingAdapter(pipeline)
	if err != nil {
		log.Fatal(err)
	}

	defer sessionManager.Destroy()

	fmt.Println(embeddingAdapter)

	db.AutoMigrate(&domain.Topic{}, &domain.TopicImage{}, &domain.TopicEmbedding{})
	fmt.Println("Migrate is successfully!")

	//Initial Minio
	minioAdapter, err := minio.NewMinioAdapter(
		config.GetEnv("MINIO_ENDPOINT", "localhost:9000"),
		config.GetEnv("MINIO_ACCESS", "minioadmin"),
		config.GetEnv("MINIO_SECRET", "minioadmin"),
		false,
		config.GetEnv("MINIO_BUCKET", "topics"),
	)
	if err != nil {
		log.Fatal(err)
	}

	topicRepo := repositories.NewTopicRepository(db)
	topicImageRepo := repositories.NewTopicImageRepository(db)
	topicEmbeddingRepo := repositories.NewTopicEmbeddingRepository(db)
	topicService := services.NewTopicService(
		db,
		topicRepo,
		topicImageRepo,
		topicEmbeddingRepo,
		minioAdapter,
		embeddingAdapter,
	)
	topicHandler := handlers.NewTopicHandler(topicService)

	app := fiber.New()
	api := app.Group("/api")
	v1 := api.Group("/v1")
	v1.Post("/topics", topicHandler.Create)
	v1.Get("/topics", topicHandler.List)
	v1.Get("/topics/:id", topicHandler.Get)

	log.Fatal(app.Listen(":3000"))
}
