package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"io/ioutil"
	"log"
)

func check(e error) {
	if e != nil {
		log.Panic(e)
	}
}

func initTelegram() (*tgbotapi.BotAPI, <-chan tgbotapi.Update) {
	bot_key, err := ioutil.ReadFile("tg_bot_api_key.txt")
	check(err)

	bot, err := tgbotapi.NewBotAPI(string(bot_key))
	check(err)

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)
	check(err)
	return bot, updates
}
