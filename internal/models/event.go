package models

import (
	"time"
)

type EventType string

const (
	EventMeeting EventType = "meeting"
	EventKotyol  EventType = "kotyol"
)

type Event struct {
	ID              string     `json:"id"`
	Type            EventType  `json:"type"`
	Title           string     `json:"title"`
	Description     string     `json:"description"`
	Date            *time.Time `json:"date,omitempty"`
	Time            *time.Time `json:"time,omitempty"`
	Place           *string    `json:"place,omitempty"`
	StartMonth      *time.Time `json:"startMonth,omitempty"`
	EndMonth        *time.Time `json:"endMonth,omitempty"`
	Amount          *int       `json:"amount,omitempty"`
	OwnerID         int64      `json:"ownerId"`
	MaxParticipants *int       `json:"maxParticipants,omitempty"`
	Participants    []int64    `json:"participants"`
	CreatedAt       time.Time  `json:"createdAt"`
}
