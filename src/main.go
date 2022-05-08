package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"regexp"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

var fuckRegex = regexp.MustCompile("f+u+c+k*")
var fukcRegex = regexp.MustCompile("f+u+k+c*")

var гudeScreams = [5]string{"WHAT THE FUCK", "FUCK YOU", "SHUT THE FUCK UP", "EAT SHIT", "LOST THE GAME?"}
var veryGood = []string{"very good", "very very good", "very very very good", "super cool", "very very well", "well done", "nice", "very nice", "not bad", "quite not bad", "okay", "very very okay", "about to lose the GAME"}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	bot, err := tgbotapi.NewBotAPI(os.Getenv("API_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	sorrySirLog := func(chatId int64, err error) {
		msg := tgbotapi.NewMessage(chatId, fmt.Sprintf("im so sory sir, but i got an error: %s", err))
		if _, err := bot.Send(msg); err != nil {
			log.Println(err)
		}
	}

	fuckManager := NewFuckManager()
	if err := fuckManager.LoadData(); err != nil {
		log.Fatal(err)
	}

	log.Println("bot started")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil {
			continue
		}

		userId := update.Message.From.ID
		userName := update.Message.From.UserName

		if !update.Message.IsCommand() {
			found1 := fuckRegex.Find([]byte(update.Message.Text))
			found2 := fukcRegex.Find([]byte(update.Message.Text))

			if found1 != nil || found2 != nil {
				fuckManager.AddMessage(userId, userName, 1)
				log.Printf("got fuck message from user %d\n", userId)
			} else {
				fuckManager.AddMessage(userId, userName, 0)
				log.Printf("got non-fuck message from user %d\n", userId)
			}

			if err := fuckManager.SaveData(); err != nil {
				sorrySirLog(update.Message.Chat.ID, err)
			}

			continue
		}

		command := update.Message.Command()
		log.Printf("got command %q from user %s\n", command, update.Message.From.UserName)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		msg.ParseMode = tgbotapi.ModeMarkdown

		if command != "nigga" {
			msg.Text = гudeScreams[rand.Intn(len(гudeScreams))]
		} else {
			for _, v := range fuckManager.Entries {
				percents := v.Fuck * 100 / v.Total
				resultMessage := veryGood[rand.Intn(len(veryGood))]
				line := fmt.Sprintf("_%s_ fuc k is *%d%%* (%d/%d) - %s result\n", v.Name, percents, v.Fuck, v.Total, resultMessage)

				msg.Text = msg.Text + line
			}
		}

		if _, err := bot.Send(msg); err != nil {
			log.Println(err)
		}
	}
}
