package handler

import (
	"encoding/json"
	"net/http"

	"meetico/internal/model"
	"meetico/internal/repo"
)

type EventTypesHandler struct {
	eventTypesRepo *repo.EventTypesRepo
}

func NewEventTypesHandler(eventTypesRepo *repo.EventTypesRepo) *EventTypesHandler {
	return &EventTypesHandler{eventTypesRepo: eventTypesRepo}
}

func (h *EventTypesHandler) GetEventTypes(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	eventTypes, err := h.eventTypesRepo.GetByUserID(r.Context(), userID)
	if err != nil || eventTypes == nil {
		writeJSON(w, http.StatusOK, []*model.EventType{})
		return
	}

	writeJSON(w, http.StatusOK, eventTypes)
}

type CreateEventTypeRequest struct {
	Slug          string `json:"slug"`
	Title         string `json:"title"`
	DurationMin   int    `json:"duration_min"`
	LocationType  string `json:"location_type"`
	LocationValue string `json:"location_value"`
	Category      string `json:"category"`
}

func (h *EventTypesHandler) CreateEventType(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	var req CreateEventTypeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	eventType := &model.EventType{
		UserID:        userID,
		Slug:          req.Slug,
		Title:         req.Title,
		DurationMin:   req.DurationMin,
		LocationType:  req.LocationType,
		LocationValue: req.LocationValue,
		Category:      req.Category,
		IsActive:      true,
	}

	err := h.eventTypesRepo.Create(r.Context(), eventType)
	if err != nil {
		http.Error(w, "Failed to create event type", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, eventType)
}

type UpdateEventTypeRequest struct {
	Title         *string `json:"title"`
	DurationMin   *int    `json:"duration_min"`
	LocationType  *string `json:"location_type"`
	LocationValue *string `json:"location_value"`
	IsActive      *bool   `json:"is_active"`
}

func (h *EventTypesHandler) UpdateEventType(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/api/event-types/"):]
	userID := r.Header.Get("X-User-ID")
	var req UpdateEventTypeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	eventType, err := h.eventTypesRepo.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Event type not found", http.StatusNotFound)
		return
	}

	if eventType.UserID != userID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	err = h.eventTypesRepo.Update(r.Context(), id, req.Title, req.DurationMin, req.LocationType, req.LocationValue, req.IsActive)
	if err != nil {
		http.Error(w, "Failed to update event type", http.StatusInternalServerError)
		return
	}

	updatedEventType, err := h.eventTypesRepo.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Failed to get updated event type", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, updatedEventType)
}

func (h *EventTypesHandler) DeleteEventType(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/api/event-types/"):]
	userID := r.Header.Get("X-User-ID")

	eventType, err := h.eventTypesRepo.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Event type not found", http.StatusNotFound)
		return
	}

	if eventType.UserID != userID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	err = h.eventTypesRepo.Delete(r.Context(), id)
	if err != nil {
		http.Error(w, "Failed to delete event type", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
