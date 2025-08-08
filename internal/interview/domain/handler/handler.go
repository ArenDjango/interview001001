package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"online-registration/internal/interview/domain/dto"
	"online-registration/internal/interview/domain/usecase"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/rs/zerolog/log"
)

type CreateEventRequest struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
}

type Handler struct {
	createEventUseCase *usecase.CreateEventUseCase
}

// NewHandler creates a new HTTP handler
func NewHandler(proxyUseCase *usecase.CreateEventUseCase) *Handler {
	return &Handler{
		createEventUseCase: proxyUseCase,
	}
}

func (h *Handler) CreateEvent(c *gin.Context) {
	var req CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ErrorData": "Uncorrected data"})
		return
	}

	// validating that starttime less then endtime
	if req.StartTime.After(req.EndTime) {
		log.Error().Msgf("Start time %s cannot be after end time %s", req.StartTime, req.EndTime)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Start time cannot be after end time"})
		return
	}

	if req.Title == "" || req.Description == "" || req.StartTime.IsZero() || req.EndTime.IsZero() {
		log.Error().Msg("Data cannot be empty")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data  cannot be empty"})
		return
	}

	requestDTO := dto.CreateEventRequestDTO{
		Title:       req.Title,
		Description: req.Description,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
	}

	reqBody, err := json.Marshal(requestDTO)
	if err != nil {
		log.Error().Msgf("Failed to marshal CreateEventRequestDTO: %v", err)
		c.JSON(http.StatusExpectationFailed, "result")
	}

	log.Info().Msgf("Handler creating event request: %s", string(reqBody))

	event, err := h.createEventUseCase.CreateEvent(context.Background(), &requestDTO)
	if err != nil {
		log.Error().Msgf("Failed to create event: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create event"})
		return
	}

	c.JSON(http.StatusCreated, event)
	return
}

//
//func (h *Handler) GetEvents(c *gin.Context) {
//
//	c.JSON(http.StatusOK, gin.H{"message": "GetEvents not implemented yet"})
//
//	events, err := h.getEventsUseCase.CreateEvent(context.Background(), &requestDTO)
//	if err != nil {
//		log.Error().Msgf("Failed to create event: %v", err)
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create event"})
//		return
//	}
//
//	c.JSON(http.StatusCreated, event)
//	return
//}
