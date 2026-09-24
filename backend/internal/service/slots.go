package service

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"meetico/internal/repo"
)

type SlotsService struct {
	schedulesRepo *repo.SchedulesRepo
	bookingsRepo  *repo.BookingsRepo
	cache         sync.Map
}

func NewSlotsService(schedulesRepo *repo.SchedulesRepo, bookingsRepo *repo.BookingsRepo) *SlotsService {
	return &SlotsService{
		schedulesRepo: schedulesRepo,
		bookingsRepo:  bookingsRepo,
	}
}

type SlotCacheEntry struct {
	slots  []string
	expiry time.Time
}

func (s *SlotsService) GetAvailableSlots(ctx context.Context, userID string, date time.Time, durationMin int, clientTimezone string) ([]string, error) {
	cacheKey := fmt.Sprintf("%s:%s", userID, date.Format("2006-01-02"))

	if entry, ok := s.cache.Load(cacheKey); ok {
		cached := entry.(SlotCacheEntry)
		if time.Now().Before(cached.expiry) {
			return cached.slots, nil
		}
		s.cache.Delete(cacheKey)
	}

	schedule, err := s.schedulesRepo.GetByUserID(ctx, userID)
	if err != nil {
		return []string{}, nil
	}

	recruiterLoc, err := time.LoadLocation(schedule.Timezone)
	if err != nil {
		return nil, err
	}

	clientLoc, err := time.LoadLocation(clientTimezone)
	if err != nil {
		return nil, err
	}

	dayOfWeek := date.Weekday().String()[:3]
	dayOfWeek = fmt.Sprintf("%s%s", string(dayOfWeek[0]-'A'+'a'), dayOfWeek[1:])

	intervals, ok := schedule.WeeklyHours[dayOfWeek]
	if !ok {
		return []string{}, nil
	}

	var slots []time.Time
	for _, interval := range intervals {
		if len(interval) < 2 {
			slog.Error("invalid interval format", "interval", interval)
			continue
		}
		startTime, err := time.Parse("15:04", interval[0])
		if err != nil {
			slog.Error("failed to parse interval start", "error", err)
			continue
		}
		endTime, err := time.Parse("15:04", interval[1])
		if err != nil {
			slog.Error("failed to parse interval end", "error", err)
			continue
		}

		for t := startTime; t.Add(time.Duration(durationMin)*time.Minute).Before(endTime) || t.Add(time.Duration(durationMin)*time.Minute).Equal(endTime); t = t.Add(time.Duration(durationMin) * time.Minute) {
			slotTime := time.Date(date.Year(), date.Month(), date.Day(), t.Hour(), t.Minute(), 0, 0, recruiterLoc)
			slots = append(slots, slotTime)
		}
	}

	bookings, err := s.bookingsRepo.GetByUserIDAndDate(ctx, userID, date)
	if err != nil {
		return nil, err
	}

	var availableSlots []time.Time
	for _, slot := range slots {
		occupied := false
		for _, booking := range bookings {
			if (slot.Equal(booking.StartsAt) || slot.After(booking.StartsAt)) && slot.Before(booking.EndsAt) {
				occupied = true
				break
			}
		}
		if !occupied {
			availableSlots = append(availableSlots, slot)
		}
	}

	var result []string
	for _, slot := range availableSlots {
		inClientTime := slot.In(clientLoc)
		result = append(result, inClientTime.Format("15:04"))
	}

	s.cache.Store(cacheKey, SlotCacheEntry{
		slots:  result,
		expiry: time.Now().Add(60 * time.Second),
	})

	return result, nil
}

func (s *SlotsService) InvalidateCache(userID string, date time.Time) {
	cacheKey := fmt.Sprintf("%s:%s", userID, date.Format("2006-01-02"))
	s.cache.Delete(cacheKey)
}
