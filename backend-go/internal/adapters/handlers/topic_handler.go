package handlers

import (
	"strconv"
	"time"

	"github.com/chanasia/semantic-search-system/internal/core/domain"
	"github.com/chanasia/semantic-search-system/internal/core/ports"
	"github.com/gofiber/fiber/v2"
)

type TopicHandler struct {
	topicService ports.TopicService
}

func NewTopicHandler(service ports.TopicService) *TopicHandler {
	return &TopicHandler{topicService: service}
}

func (h *TopicHandler) Create(c *fiber.Ctx) error {
	type CreateTopicRequest struct {
		Title   string `form:"title"`
		Context string `form:"context"`
	}

	var req CreateTopicRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid form data",
		})
	}

	files := form.File["images"]

	topic := &domain.Topic{
		Title:     req.Title,
		Context:   req.Context,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	ctx := c.Context()
	if err := h.topicService.CreateTopic(ctx, topic, files); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":  "Failed to create topic",
			"detail": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(
		fiber.Map{
			"message": "Topic created!",
		})
}

func (h *TopicHandler) Get(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	topic, err := h.topicService.GetTopic(int64(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Topic not found"})
	}

	return c.JSON(topic)
}

func (h *TopicHandler) List(c *fiber.Ctx) error {
	topics, err := h.topicService.GetTopics()
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(topics)
}
