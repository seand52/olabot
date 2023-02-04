package main

import (
	"log"
	"os"
	"strings"
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
	u.AllowedUpdates = []string{"message"}
	// don't make the bot inline bot

	updates := bot.GetUpdatesChan(u)
	counter := 0
	reset := time.Tick(time.Minute)

	for update := range updates {
		if update.Message != nil {
				select {
				case <-reset:
					counter = 0
				default:
					log.Printf("counter is %d", counter)
					if counter < 5 {
						if 	strings.ToLower(update.Message.Text) == "ola" {
							counter++
							if counter > 5 {
								counter = 0
								<-reset
							} else {
								log.Printf("hello [%s] %s", update.Message.From.UserName, update.Message.Text)
			
								msg := tgbotapi.NewMessage(update.Message.Chat.ID, "ola")
								msg.ReplyToMessageID = update.Message.MessageID
			
								bot.Send(msg)
							}
			
						}
					} else {
						log.Printf("Rate limit reached. Try again in %s", time.Until(time.Now().Truncate(time.Minute).Add(time.Minute)))
					}
				}
			
		}
	}
}
