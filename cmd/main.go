package main

import (
	"fmt"
	"gatherbot-backend/config"
	"gatherbot-backend/internal/handlers"
	"gatherbot-backend/internal/repository"
	"gatherbot-backend/internal/services"
	"gatherbot-backend/internal/telegram"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"log"
	"net/http"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}

	config.InitDB()
	defer config.DB.Close()

	bot := telegram.NewBot()

	repo := repository.NewEventRepository()
	svc := services.NewEventService(repo, bot)
	handler := handlers.NewEventHandler(svc)

	port := ":8080"
	r := chi.NewRouter()

	// TODO: подключение маршрутов
	r.Post("/events", handler.CreateEvent)
	r.Get("/events/{id}", handler.GetEventByID)
	r.Get("/events", handler.GetEvents)
	r.Patch("/events/{id}/join", handler.JoinEvent)
	
	fmt.Println("Сервер слушает на", port)
	http.ListenAndServe(port, r)
}
