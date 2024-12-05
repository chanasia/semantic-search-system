package handlers

import (
	"strconv"

	"github.com/chanasia/semantic-search-system/internal/core/domain"
	"github.com/chanasia/semantic-search-system/internal/core/ports"
	"github.com/gofiber/fiber/v3"
)

type TopicHandler struct {
	topicService ports.TopicService
}

func NewTopicHandler(service ports.TopicService) *TopicHandler {
	return &TopicHandler{topicService: service}
}

func (h *TopicHandler) Create(c fiber.Ctx) error {
	var topic domain.Topic

	if err := c.Bind().Body(&topic); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if err := h.topicService.CreateTopic(&topic); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.Status(fiber.StatusCreated).JSON(topic)
}

func (h *TopicHandler) Get(c fiber.Ctx) error {
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

func (h *TopicHandler) List(c fiber.Ctx) error {
	topics, err := h.topicService.GetTopics()
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(topics)
}
