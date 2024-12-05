package main

import (
	"fmt"
	"log"

	"github.com/chanasia/semantic-search-system/internal/adapters/handlers"
	"github.com/chanasia/semantic-search-system/internal/adapters/repositories"
	"github.com/chanasia/semantic-search-system/internal/config"
	"github.com/chanasia/semantic-search-system/internal/core/domain"
	"github.com/chanasia/semantic-search-system/internal/core/services"
	"github.com/gofiber/fiber/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dbConfig, err := config.LoadDBConfig()
	if err != nil {
		panic(err)
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable search_path=%s", dbConfig.DBHost, dbConfig.DBPort, dbConfig.DBUser, dbConfig.DBPassword, dbConfig.DBName, dbConfig.DBSchema)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	db.AutoMigrate(&domain.Topic{})
	println("Migrate is successfully!")

	topicRepo := repositories.NewTopicRepository(db)
	topicService := services.NewTopicService(topicRepo)
	topicHandler := handlers.NewTopicHandler(topicService)

	app := fiber.New()
	api := app.Group("/api")
	v1 := api.Group("/v1")
	v1.Post("/topics", topicHandler.Create)
	v1.Get("/topics", topicHandler.List)
	v1.Get("/topics/:id", topicHandler.Get)

	log.Fatal(app.Listen(":3000"))
}
