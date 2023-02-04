package main

import (
	"log"
	"os"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// fill the bucket with 5 tokens. This applies to all chats where bot is active for simplicity
var rateLimit = 5

func handleUpdate(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	log.Printf("received update and message is %+v\n ", update.Message)
	if strings.ToLower(update.Message.Text) == "ola" {
		if rateLimit > 0 {
			rateLimit--
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "ola")
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
		} else {
			log.Print("rate limit exceeded")
		}	
	}
}

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
	ticker := time.NewTicker(time.Minute)

	for {
		select {
		case update := <-updates:
			if update.Message != nil {
				handleUpdate(update, bot)
			}
		case <-ticker.C:
			log.Println("tick")
			rateLimit = 5
		}
	}
}
