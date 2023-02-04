package main

import (
	"log"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	counter := 0
	reset := time.Tick(time.Minute)
	for update := range updates {
		if update.Message != nil {
			if update.Message.Text == "ola" {
				counter++
				if counter > 5 {
					counter = 0
					<-reset
				} else {
					log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "ola")
					msg.ReplyToMessageID = update.Message.MessageID

					bot.Send(msg)
				}

			}
		}
	}
}
