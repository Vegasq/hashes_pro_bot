package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// https://github.com/noraj/haiti/blob/master/data/prototypes.json
//go:embed prototypes.json
var prototypes string

type Mode struct {
	Regex string `json:"regex"`
	Modes []struct {
		John    interface{} `json:"john"`
		Hashcat int         `json:"hashcat"`
		Name    string      `json:"name"`
	} `json:"modes"`

	RegexProcessed *regexp.Regexp
}

type Modes []Mode

func (m Modes) Stringify() string {
	var s string
	var names []string
	for i := range m {
		for j := range m[i].Modes {
			s += m[i].Modes[j].Name + " "
			names = append(names, fmt.Sprintf("%s (%d)", m[i].Modes[j].Name, m[i].Modes[j].Hashcat))
		}
	}

	return strings.Join(names, "\n")
}

func (m Modes) Find(s []byte) Modes {
	filteredModes := Modes{}
	for i := range m {
		if m[i].RegexProcessed != nil && m[i].RegexProcessed.Match(s) {
			filteredModes = append(filteredModes, m[i])
		}
	}
	return filteredModes
}

func prepareModes() Modes {
	var modes Modes
	err := json.Unmarshal([]byte(prototypes), &modes)
	if err != nil {
		log.Println("Failed to parse prototypes file.")
		log.Fatal(err)
	}

	for i := range modes {
		re, err := regexp.Compile("(?i)" + modes[i].Regex)
		if err != nil {
			log.Printf("Failed to compile regex for %s, %s, %e", modes[i].Modes[0].Name, modes[i].Regex, err)
		}
		modes[i].RegexProcessed = re
	}

	return modes
}

func processCommand(update tgbotapi.Update) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	switch update.Message.Command() {
	default:
		msg.Text = "Hashes.pro\n\n" +
			"API: https://github.com/go-telegram-bot-api/telegram-bot-api\n" +
			"Regexes: https://github.com/noraj/haiti/"
		msg.DisableWebPagePreview = true
	}
	return msg
}

func processHashRequest(modes Modes, update tgbotapi.Update) tgbotapi.MessageConfig {
	found := modes.Find([]byte(update.Message.Text)).Stringify()
	if found == "" {
		found = "No matches found."
	}
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, found)
	msg.ReplyToMessageID = update.Message.MessageID

	return msg
}

func main() {
	token := os.Getenv("TELEGRAM_TOKEN")

	modes := prepareModes()
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			msg := processCommand(update)
			bot.Send(msg)
		} else {
			msg := processHashRequest(modes, update)
			bot.Send(msg)
		}
	}
}
