package handler

import (
	"errors"
	"net/http"

	"meetico/internal/model"
	"meetico/internal/repo"
	"meetico/internal/service"

	"github.com/go-chi/chi/v5"
)

type BookingsHandler struct {
	bookingsRepo   *repo.BookingsRepo
	bookingService *service.BookingService
}

func NewBookingsHandler(bookingsRepo *repo.BookingsRepo, bookingService *service.BookingService) *BookingsHandler {
	return &BookingsHandler{
		bookingsRepo:   bookingsRepo,
		bookingService: bookingService,
	}
}

func (h *BookingsHandler) GetBookings(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	bookings, err := h.bookingsRepo.GetByUserID(r.Context(), userID)
	if err != nil || bookings == nil {
		writeJSON(w, http.StatusOK, []*model.Booking{})
		return
	}

	writeJSON(w, http.StatusOK, bookings)
}

func (h *BookingsHandler) CancelBooking(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, `{"error":"missing booking id"}`, http.StatusBadRequest)
		return
	}

	err := h.bookingService.CancelBookingByOwner(r.Context(), userID, id)
	if errors.Is(err, service.ErrForbidden) {
		http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
		return
	}
	if err != nil {
		http.Error(w, `{"error":"failed to cancel booking"}`, http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}
