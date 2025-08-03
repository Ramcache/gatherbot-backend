package services

import (
	"context"
	"fmt"
	"gatherbot-backend/internal/telegram"
	"time"

	"gatherbot-backend/internal/models"
	"gatherbot-backend/internal/repository"
)

type EventService struct {
	repo *repository.EventRepository
	bot  *telegram.Bot
}

func NewEventService(repo *repository.EventRepository, bot *telegram.Bot) *EventService {
	return &EventService{repo: repo, bot: bot}
}

func (s *EventService) CreateEvent(ctx context.Context, e *models.Event) (string, error) {
	e.CreatedAt = time.Now()
	return s.repo.Insert(ctx, e)
}

func (s *EventService) GetEvents(ctx context.Context, ownerId *int64, participantId *int64) ([]*models.Event, error) {
	return s.repo.GetAllByFilter(ctx, ownerId, participantId)
}

func (s *EventService) JoinEvent(ctx context.Context, eventID string, userID int64) error {
	event, err := s.repo.GetByID(ctx, eventID)
	if err != nil {
		return err
	}
	if event == nil {
		return fmt.Errorf("событие не найдено")
	}

	// проверка, уже ли в участниках
	for _, id := range event.Participants {
		if id == userID {
			return nil // уже участвует
		}
	}

	// проверка на лимит
	if event.MaxParticipants != nil && len(event.Participants) == *event.MaxParticipants {
		// Группа набрана — уведомить участников
		go s.bot.NotifyGroupFull(event.Title, event.Participants)
	}

	event.Participants = append(event.Participants, userID)

	// обновим в базе
	return s.repo.UpdateParticipants(ctx, eventID, event.Participants)
}

func (s *EventService) GetByID(ctx context.Context, id string) (*models.Event, error) {
	return s.repo.GetByID(ctx, id)
}
