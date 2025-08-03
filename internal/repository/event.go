package repository

import (
	"context"
	"fmt"

	"gatherbot-backend/config"
	"gatherbot-backend/internal/models"
	"github.com/jackc/pgx/v5"
)

type EventRepository struct{}

func NewEventRepository() *EventRepository {
	return &EventRepository{}
}

func (r *EventRepository) GetByID(ctx context.Context, id string) (*models.Event, error) {
	query := `SELECT id, type, title, description, date, time, place, start_month, end_month, amount, owner_id, max_participants, participants, created_at
	FROM events WHERE id = $1`

	row := config.DB.QueryRow(ctx, query, id)

	var event models.Event
	err := row.Scan(
		&event.ID,
		&event.Type,
		&event.Title,
		&event.Description,
		&event.Date,
		&event.Time,
		&event.Place,
		&event.StartMonth,
		&event.EndMonth,
		&event.Amount,
		&event.OwnerID,
		&event.MaxParticipants,
		&event.Participants,
		&event.CreatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("ошибка получения события: %w", err)
	}

	return &event, nil
}

func (r *EventRepository) Insert(ctx context.Context, e *models.Event) (string, error) {
	query := `INSERT INTO events (type, title, description, date, time, place, start_month, end_month, amount, owner_id, max_participants, participants, created_at)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13) RETURNING id`

	var id string
	err := config.DB.QueryRow(ctx, query,
		e.Type, e.Title, e.Description, e.Date, e.Time, e.Place,
		e.StartMonth, e.EndMonth, e.Amount, e.OwnerID,
		e.MaxParticipants, e.Participants, e.CreatedAt,
	).Scan(&id)

	return id, err
}

func (r *EventRepository) GetAllByFilter(ctx context.Context, ownerId *int64, participantId *int64) ([]*models.Event, error) {
	var (
		rows pgx.Rows
		err  error
	)

	if ownerId != nil {
		rows, err = config.DB.Query(ctx, `
			SELECT id, type, title, description, date, time, place, start_month, end_month, amount,
			       owner_id, max_participants, participants, created_at
			FROM events
			WHERE owner_id = $1
			ORDER BY created_at DESC
		`, *ownerId)
	} else if participantId != nil {
		rows, err = config.DB.Query(ctx, `
			SELECT id, type, title, description, date, time, place, start_month, end_month, amount,
			       owner_id, max_participants, participants, created_at
			FROM events
			WHERE $1 = ANY(participants)
			ORDER BY created_at DESC
		`, *participantId)
	} else {
		rows, err = config.DB.Query(ctx, `
			SELECT id, type, title, description, date, time, place, start_month, end_month, amount,
			       owner_id, max_participants, participants, created_at
			FROM events
			ORDER BY created_at DESC
		`)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*models.Event

	for rows.Next() {
		var e models.Event
		err := rows.Scan(
			&e.ID, &e.Type, &e.Title, &e.Description,
			&e.Date, &e.Time, &e.Place,
			&e.StartMonth, &e.EndMonth, &e.Amount,
			&e.OwnerID, &e.MaxParticipants, &e.Participants,
			&e.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		events = append(events, &e)
	}

	return events, nil
}

func (r *EventRepository) UpdateParticipants(ctx context.Context, eventID string, participants []int64) error {
	query := `UPDATE events SET participants = $1 WHERE id = $2`
	_, err := config.DB.Exec(ctx, query, participants, eventID)
	return err
}
