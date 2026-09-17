package handler

import (
	"encoding/json"
	"net/http"

	"meetico/internal/repo"
)

type BookingsHandler struct {
	bookingsRepo *repo.BookingsRepo
}

func NewBookingsHandler(bookingsRepo *repo.BookingsRepo) *BookingsHandler {
	return &BookingsHandler{bookingsRepo: bookingsRepo}
}

func (h *BookingsHandler) GetBookings(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	bookings, err := h.bookingsRepo.GetByUserID(r.Context(), userID)
	if err != nil {
		http.Error(w, "Failed to get bookings", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(bookings)
}
