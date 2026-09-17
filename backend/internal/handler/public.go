package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"meetico/internal/repo"
	"meetico/internal/service"
)

type PublicHandler struct {
	usersRepo      *repo.UsersRepo
	eventTypesRepo *repo.EventTypesRepo
	bookingsRepo   *repo.BookingsRepo
	schedulesRepo  *repo.SchedulesRepo
	slotsService   *service.SlotsService
	bookingService *service.BookingService
}

func NewPublicHandler(
	usersRepo *repo.UsersRepo,
	eventTypesRepo *repo.EventTypesRepo,
	bookingsRepo *repo.BookingsRepo,
	schedulesRepo *repo.SchedulesRepo,
	slotsService *service.SlotsService,
	bookingService *service.BookingService,
) *PublicHandler {
	return &PublicHandler{
		usersRepo:      usersRepo,
		eventTypesRepo: eventTypesRepo,
		bookingsRepo:   bookingsRepo,
		schedulesRepo:  schedulesRepo,
		slotsService:   slotsService,
		bookingService: bookingService,
	}
}

type PublicBookingRequest struct {
	CandidateName  string `json:"candidate_name"`
	CandidateEmail string `json:"candidate_email"`
	SlotTime       string `json:"slot_time"`
	Date           string `json:"date"`
	Timezone       string `json:"timezone"`
}

type PublicBookingResponse struct {
	BookingID   string    `json:"booking_id"`
	StartsAt    time.Time `json:"starts_at"`
	CancelToken string    `json:"cancel_token"`
}

func (h *PublicHandler) GetPublicInfo(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Path[len("/api/public/"):]
	slugStart := len(username) + 1
	slug := r.URL.Path[slugStart:]

	user, err := h.usersRepo.GetByUsername(r.Context(), username)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	eventType, err := h.eventTypesRepo.GetBySlug(r.Context(), slug)
	if err != nil {
		http.Error(w, "Event type not found", http.StatusNotFound)
		return
	}

	response := map[string]interface{}{
		"recruiter": map[string]interface{}{
			"display_name": user.DisplayName,
			"photo_url":    user.PhotoURL,
			"bio":          user.Bio,
		},
		"event_type": map[string]interface{}{
			"title":         eventType.Title,
			"duration_min":  eventType.DurationMin,
			"location_type": eventType.LocationType,
			"category":      eventType.Category,
			"slug":          eventType.Slug,
		},
	}

	json.NewEncoder(w).Encode(response)
}

func (h *PublicHandler) GetPublicSlots(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Path[len("/api/public/"):]
	slugStart := len(username) + 1
	slugEnd := slugStart + len(r.URL.Path[slugStart:])
	slug := r.URL.Path[slugStart:slugEnd]

	dateStr := r.URL.Query().Get("date")
	tz := r.URL.Query().Get("tz")

	if dateStr == "" || tz == "" {
		http.Error(w, "Missing date or timezone parameter", http.StatusBadRequest)
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	user, err := h.usersRepo.GetByUsername(r.Context(), username)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	eventType, err := h.eventTypesRepo.GetBySlug(r.Context(), slug)
	if err != nil {
		http.Error(w, "Event type not found", http.StatusNotFound)
		return
	}

	slots, err := h.slotsService.GetAvailableSlots(r.Context(), user.ID, date, eventType.DurationMin, tz)
	if err != nil {
		http.Error(w, "Failed to get slots", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"slots": slots,
	}

	json.NewEncoder(w).Encode(response)
}

func (h *PublicHandler) CreatePublicBooking(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Path[len("/api/public/"):]
	slugStart := len(username) + 1
	slugEnd := slugStart + len(r.URL.Path[slugStart:])
	slug := r.URL.Path[slugStart:slugEnd]

	var req PublicBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	eventType, err := h.eventTypesRepo.GetBySlug(r.Context(), slug)
	if err != nil {
		http.Error(w, "Event type not found", http.StatusNotFound)
		return
	}

	booking, err := h.bookingService.CreateBooking(r.Context(), eventType.ID, req.CandidateName, req.CandidateEmail, req.SlotTime, req.Date, req.Timezone)
	if err != nil {
		http.Error(w, "Failed to create booking", http.StatusInternalServerError)
		return
	}

	response := PublicBookingResponse{
		BookingID:   booking.ID,
		StartsAt:    booking.StartsAt,
		CancelToken: booking.CancelToken,
	}

	json.NewEncoder(w).Encode(response)
}

func (h *PublicHandler) GetCancelInfo(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Path[len("/api/public/cancel/"):]

	booking, err := h.bookingsRepo.GetByCancelToken(r.Context(), token)
	if err != nil {
		http.Error(w, "Booking not found", http.StatusNotFound)
		return
	}

	eventType, err := h.eventTypesRepo.GetByID(r.Context(), booking.EventTypeID)
	if err != nil {
		http.Error(w, "Event type not found", http.StatusNotFound)
		return
	}

	user, err := h.usersRepo.GetByID(r.Context(), eventType.UserID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	response := map[string]interface{}{
		"candidate_name":         booking.CandidateName,
		"starts_at":              booking.StartsAt,
		"duration_min":           eventType.DurationMin,
		"event_type_title":       eventType.Title,
		"recruiter_display_name": user.DisplayName,
		"status":                 booking.Status,
	}

	json.NewEncoder(w).Encode(response)
}

func (h *PublicHandler) CancelBooking(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Path[len("/api/public/cancel/"):]

	err := h.bookingService.CancelBooking(r.Context(), token)
	if err != nil {
		http.Error(w, "Failed to cancel booking", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"ok": true,
	}

	json.NewEncoder(w).Encode(response)
}
