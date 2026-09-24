package handler

import (
	"encoding/json"
	"net/http"

	"meetico/internal/repo"
)

type ScheduleHandler struct {
	schedulesRepo *repo.SchedulesRepo
}

func NewScheduleHandler(schedulesRepo *repo.SchedulesRepo) *ScheduleHandler {
	return &ScheduleHandler{schedulesRepo: schedulesRepo}
}

type ScheduleResponse struct {
	Timezone    string                `json:"timezone"`
	WeeklyHours map[string][][]string `json:"weekly_hours"`
}

func (h *ScheduleHandler) GetSchedule(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	schedule, err := h.schedulesRepo.GetByUserID(r.Context(), userID)
	if err != nil {
		writeJSON(w, http.StatusOK, ScheduleResponse{
			Timezone:    "Europe/Moscow",
			WeeklyHours: map[string][][]string{},
		})
		return
	}

	weeklyHours := schedule.WeeklyHours
	if weeklyHours == nil {
		weeklyHours = map[string][][]string{}
	}

	writeJSON(w, http.StatusOK, ScheduleResponse{
		Timezone:    schedule.Timezone,
		WeeklyHours: weeklyHours,
	})
}

type UpdateScheduleRequest struct {
	Timezone    *string                `json:"timezone"`
	WeeklyHours *map[string][][]string `json:"weekly_hours"`
}

func (h *ScheduleHandler) UpdateSchedule(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	var req UpdateScheduleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	err := h.schedulesRepo.Update(r.Context(), userID, req.Timezone, req.WeeklyHours)
	if err != nil {
		http.Error(w, "Failed to update schedule", http.StatusInternalServerError)
		return
	}

	updatedSchedule, err := h.schedulesRepo.GetByUserID(r.Context(), userID)
	if err != nil {
		http.Error(w, "Failed to get updated schedule", http.StatusInternalServerError)
		return
	}

	weeklyHours := updatedSchedule.WeeklyHours
	if weeklyHours == nil {
		weeklyHours = map[string][][]string{}
	}

	writeJSON(w, http.StatusOK, ScheduleResponse{
		Timezone:    updatedSchedule.Timezone,
		WeeklyHours: weeklyHours,
	})
}
