package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// fill the bucket with 5 tokens. This applies to all chats where bot is active for simplicity
var rateLimit = 5

var People = make(map[string]time.Time)

type MythicPlusScore struct {
	Score float64 `json:"score"`
	Color string  `json:"color"`
}

type MythicPlusScoresBySeason struct {
	Season  string                `json:"season"`
	Scores  map[string]float64    `json:"scores"`
	Segments map[string]MythicPlusScore `json:"segments"`
}

type Data struct {
	MythicPlusScoresBySeason []MythicPlusScoresBySeason `json:"mythic_plus_scores_by_season"`
}

func IsValentinesDay() bool {
	currentTime := time.Now()
	currentMonth := currentTime.Month()
	currentDay := currentTime.Day()
	if (currentMonth == time.February && currentDay == 14) {
		return true
	}
	return false
}

func HandleUpdate(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	log.Printf("received update and message is %+v\n ", update.Message)
	if rateLimit > 0 {
		if strings.ToLower(update.Message.Text) == "ola" {
			rateLimit--
			replyMsg := "ola"
			if IsValentinesDay() {
				replyMsg = "ola guapo, feliz san valentin :)"
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, replyMsg)
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
		}

		if strings.ToLower(update.Message.From.FirstName) == "sergi" {
			re := regexp.MustCompile(`(http|ftp|https):\/\/([\w\-_]+(?:(?:\.[\w\-_]+)+))([\w\-\.,@?^=%&:/~\+#]*[\w\-\@?^=%&/~\+#])?`)
			matches := re.MatchString(strings.ToLower(update.Message.Text))
			if matches == true {
				rateLimit--
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "OJUT!!! Puede ser deep web!")
				msg.ReplyToMessageID = update.Message.MessageID
				bot.Send(msg)
			}
		}

		if strings.HasPrefix(update.Message.Text, "/getrioscore") {
			rateLimit--
			score := GetRioStats(update)
			additionalMessage := ""
			roundedScore := math.Round(score)
			if (roundedScore < 2000) {
				additionalMessage = ". eres BASURA"
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "rio score is "+fmt.Sprintf("%.2f", roundedScore) + additionalMessage)
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
		}
		if strings.ToLower(update.Message.Text) == "que figura" || strings.ToLower(update.Message.Text) == "q figura" ||  strings.ToLower(update.Message.Text) == "figura" {
			rateLimit--
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "el marcelino")
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
		}
	} else {
		log.Print("rate limit exceeded")
	}
	
}

func HandleLastSeen(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	currentTime := time.Now().UTC()
	lastMessage := People[update.Message.From.FirstName]
	diff := currentTime.Sub(lastMessage)
	oneDay := 24 * time.Hour
	log.Print("diff is ", diff)
	if diff > oneDay {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "ola. no te vemos por aqui desde " + lastMessage.UTC().String() + ". estas mas rate limited que yo")
		msg.ReplyToMessageID = update.Message.MessageID
		bot.Send(msg)
	}
	People[update.Message.From.FirstName] = currentTime

}

func parseRealmAndName(update tgbotapi.Update) (string, string) {
	params := strings.Split(update.Message.Text, " ")
	var character, realm string
	for _, param := range params {
		if strings.HasPrefix(param, "character=") {
			character = strings.TrimPrefix(param, "character=")
		}
		if strings.HasPrefix(param, "realm=") {
			realm = strings.TrimPrefix(param, "realm=")
		}
	}
	return character, realm
}

func GetRioStats(update tgbotapi.Update) float64 {
	name, realm := parseRealmAndName(update)
	resp, err := http.Get(fmt.Sprintf("https://raider.io/api/v1/characters/profile?region=eu&realm=%s&name=%s&fields=mythic_plus_scores_by_season:current", realm, name))
	log.Print("requesting to ->" + fmt.Sprintf("https://raider.io/api/v1/characters/profile?region=eu&realm=%s&name=%s&fields=mythic_plus_scores_by_season:current", realm, name))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var unmarshalledData Data
	json.Unmarshal(body, &unmarshalledData)
	return unmarshalledData.MythicPlusScoresBySeason[0].Scores["all"]
}

func main() {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	currentTime := time.Now().UTC()
	People["Sean"] = currentTime
	People["Sergi"] = currentTime
	People["Adam"] = currentTime
	People["Aleix"] = currentTime
	People["\u00c1lvaro"] = currentTime


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
				HandleUpdate(update, bot)
				HandleLastSeen(update, bot)
			}
		case <-ticker.C:
			log.Println("tick")
			rateLimit = 5
		}
	}
}
