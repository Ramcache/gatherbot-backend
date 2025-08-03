package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"

	"gatherbot-backend/internal/models"
	"gatherbot-backend/internal/services"
)

type EventHandler struct {
	service *services.EventService
}

func NewEventHandler(s *services.EventService) *EventHandler {
	return &EventHandler{service: s}
}

type JoinRequest struct {
	UserID int64 `json:"userId"`
}

func (h *EventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var event models.Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, "Некорректный JSON", http.StatusBadRequest)
		return
	}

	id, err := h.service.CreateEvent(r.Context(), &event)
	if err != nil {
		http.Error(w, "Ошибка при создании события", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": id})
}

func (h *EventHandler) GetEvents(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	var (
		ownerIdPtr       *int64
		participantIdPtr *int64
	)

	if ownerStr := query.Get("ownerId"); ownerStr != "" {
		id, err := strconv.ParseInt(ownerStr, 10, 64)
		if err != nil {
			http.Error(w, "Неверный ownerId", http.StatusBadRequest)
			return
		}
		ownerIdPtr = &id
	}

	if participantStr := query.Get("participantId"); participantStr != "" {
		id, err := strconv.ParseInt(participantStr, 10, 64)
		if err != nil {
			http.Error(w, "Неверный participantId", http.StatusBadRequest)
			return
		}
		participantIdPtr = &id
	}

	events, err := h.service.GetEvents(r.Context(), ownerIdPtr, participantIdPtr)
	if err != nil {
		http.Error(w, "Ошибка при получении событий", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(events)
}

func (h *EventHandler) JoinEvent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req JoinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный JSON", http.StatusBadRequest)
		return
	}

	err := h.service.JoinEvent(r.Context(), id, req.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Участие подтверждено"})
}

func (h *EventHandler) GetEventByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	event, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Ошибка при получении события", http.StatusInternalServerError)
		return
	}
	if event == nil {
		http.Error(w, "Событие не найдено", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(event)
}
