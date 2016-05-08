package main

import (
	"github.com/LiaungYip/tgGlyphBot/input"
	"github.com/boltdb/bolt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"os"
	"strings"
)

func main() {
	bot, updates := initTelegram()
	db := initDatabase(databaseFilename)
	defer db.Close()

	tempdir := tempImageDir()
	defer os.RemoveAll(tempdir)

	for u := range updates {
		handleUpdate(bot, u, db, tempdir)
	}
}

func handleUpdate(bot *tgbotapi.BotAPI, u tgbotapi.Update, db *bolt.DB, tempdir string) {
	m := u.Message
	t := m.Text
	if m.Text == "" { // No action if text field is blank. (i.e. if user sends photos or stickers or voice or something weird.)
		return
	}

	if strings.HasPrefix(t, "/help") || strings.HasPrefix(t, "/start") {
		sendHelp(bot, u)
		return
	}

	if strings.HasPrefix(t, "/glyphs@IngressGlyphBot") {
		t = t[len("/glyphs@IngressGlyphBot"):]
		print(t)
	} else if strings.HasPrefix(t, "/glyphs") {
		t = t[len("/glyphs"):]
		print(t)
	}

	glyphNames, _, err := input.ProcessString(t)
	if err != nil {
		sendError(bot, u, err)
		return
	}

	log.Printf("User: %s %s (@%s), Text: %s", m.From.FirstName, m.From.LastName, m.From.UserName, m.Text)

	fileID := checkCache(glyphNames, db)
	if fileID == nil {
		sendNewSticker(bot, u, db, tempdir, glyphNames)
	} else {
		sendFromCache(bot, u, glyphNames, fileID)
	}
}

func sendError(bot *tgbotapi.BotAPI, u tgbotapi.Update, err error) {
	m := u.Message
	log.Printf("ERROR: User: %s %s (@%s), Text: %s, Error: %s", m.From.FirstName, m.From.LastName, m.From.UserName, m.Text, err.Error())
	msg := tgbotapi.NewMessage(u.Message.Chat.ID, err.Error())
	bot.Send(msg)
}

func sendNewSticker(bot *tgbotapi.BotAPI, u tgbotapi.Update, db *bolt.DB, tempdir string, glyphNames []string) {
	log.Printf("Creating new image for glyphs %s", glyphNames)

	img := makeImage(glyphNames)
	webp := encodeWebp(img, tempdir)

	msg := tgbotapi.NewStickerUpload(u.Message.Chat.ID, webp)
	response, _ := bot.Send(msg)

	fileID := response.Sticker.FileID
	addToCache(glyphNames, fileID, db)
	log.Printf("Added to cache: %s -> %s. File size: %d", glyphNames, fileID, response.Sticker.FileSize)
}

func sendFromCache(bot *tgbotapi.BotAPI, u tgbotapi.Update, glyphNames []string, fileID []byte) {
	log.Printf("Hitting cache! %s -> %s", glyphNames, fileID)
	msg := tgbotapi.NewStickerShare(u.Message.Chat.ID, string(fileID))
	bot.Send(msg)
}

func sendHelp(bot *tgbotapi.BotAPI, u tgbotapi.Update) {
	msg := tgbotapi.NewMessage(u.Message.Chat.ID, helpString)
	msg.ParseMode = "HTML"
	bot.Send(msg)
}
