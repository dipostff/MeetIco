package handler

import (
	"net/http"

	"meetico/internal/model"
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
	if err != nil || bookings == nil {
		writeJSON(w, http.StatusOK, []*model.Booking{})
		return
	}

	writeJSON(w, http.StatusOK, bookings)
}
