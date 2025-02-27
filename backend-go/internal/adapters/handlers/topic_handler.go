package handlers

import (
	"mime/multipart"
	"net/http"
	"strconv"
	"time"

	"github.com/chanasia/semantic-search-system/internal/core/domain"
	"github.com/chanasia/semantic-search-system/internal/core/ports"
	"github.com/labstack/echo/v4"
)

type TopicHandler struct {
	topicService ports.TopicService
}

func NewTopicHandler(service ports.TopicService) *TopicHandler {
	return &TopicHandler{topicService: service}
}

func (h *TopicHandler) CreateTopicJSON(c echo.Context) error {
	type CreateTopicJSONRequest struct {
		Title      string   `json:"title"`
		Context    string   `json:"context"`
		Page       string   `json:"page"`
		Tag        string   `json:"tag"`
		ImagePaths []string `json:"image_paths"`
	}

	var req CreateTopicJSONRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid JSON format",
		})
	}

	topic := &domain.Topic{
		Title:     req.Title,
		Context:   req.Context,
		Page:      req.Page,
		Tag:       req.Tag,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := h.topicService.CreateTopic(c.Request().Context(), topic, nil, req.ImagePaths); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":  "Failed to create topic",
			"detail": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Topic created!",
		"data":    topic,
	})
}

func (h *TopicHandler) CreateTopicWithFiles(c echo.Context) error {
	// Parse form fields
	title := c.FormValue("title")
	context := c.FormValue("context")
	page := c.FormValue("page")
	tag := c.FormValue("tag")
	imagePaths := c.Request().MultipartForm.Value["image_paths"]

	// Get uploaded files
	form, err := c.MultipartForm()
	var files []*multipart.FileHeader
	if err == nil && form != nil {
		files = form.File["images"]
	}

	topic := &domain.Topic{
		Title:     title,
		Context:   context,
		Page:      page,
		Tag:       tag,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := h.topicService.CreateTopic(c.Request().Context(), topic, files, imagePaths); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":  "Failed to create topic",
			"detail": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Topic created!",
		"data":    topic,
	})
}

func (h *TopicHandler) Search(c echo.Context) error {
	// Get search parameters from query string
	searchText := c.QueryParam("search_text")
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil || limit <= 0 {
		limit = 10 // Default limit
	}

	// Validate search text
	if searchText == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Search text is required",
		})
	}

	// Search topics
	results, err := h.topicService.SearchTopics(c.Request().Context(), searchText, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":  "Failed to search topics",
			"detail": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"results": results,
	})
}

func (h *TopicHandler) Get(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid ID",
		})
	}

	topic, err := h.topicService.GetTopic(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error": "Topic not found",
		})
	}

	return c.JSON(http.StatusOK, topic)
}

func (h *TopicHandler) GetTopicImage(c echo.Context) error {
	// Get image ID from path and convert to int32
	imageIDStr := c.Param("image_id")
	imageID, err := strconv.ParseInt(imageIDStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid image ID",
		})
	}

	// Get image data and content type from service
	image, contentType, err := h.topicService.GetTopicImage(c.Request().Context(), int32(imageID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":  "Failed to get image",
			"detail": err.Error(),
		})
	}

	// Set headers for image display
	c.Response().Header().Set("Content-Type", contentType)
	c.Response().Header().Set("Content-Disposition", "inline")
	c.Response().Header().Set("Cache-Control", "public, max-age=31536000")

	// Return image data
	return c.Stream(http.StatusOK, contentType, image)
}

func (h *TopicHandler) List(c echo.Context) error {
	topics, err := h.topicService.GetTopics()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to fetch topics",
		})
	}

	return c.JSON(http.StatusOK, topics)
}