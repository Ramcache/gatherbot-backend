package telegram

import (
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	api *tgbotapi.BotAPI
}

func NewBot() *Bot {
	botToken := os.Getenv("BOT_TOKEN")
	api, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatalf("❌ Не удалось инициализировать Telegram бота: %v", err)
	}

	return &Bot{api: api}
}

func (b *Bot) NotifyGroupFull(title string, userIDs []int64) {
	text := fmt.Sprintf("Группа по событию «%s» набрана. До встречи!", title)

	for _, id := range userIDs {
		msg := tgbotapi.NewMessage(id, text)
		_, err := b.api.Send(msg)
		if err != nil {
			log.Printf("⚠️ Ошибка отправки сообщения пользователю %d: %v", id, err)
		}
	}
}
