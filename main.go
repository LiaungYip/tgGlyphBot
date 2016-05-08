package main

import (
	"github.com/LiaungYip/tgGlyphBot/input"
	"github.com/boltdb/bolt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"os"
	"log"
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
	if m.Text == "" { // No action if text field is blank. (i.e. if user sends photos or stickers or voice or something weird.)
		return
	}
	t := m.Text
	glyphNames, _, err := input.ProcessString(t)

	if err != nil {
		log.Printf("ERROR: User: %s %s (@%s), Text: %s, Error: %s", m.From.FirstName, m.From.LastName, m.From.UserName, m.Text, err.Error())
		msg := tgbotapi.NewMessage(u.Message.Chat.ID, err.Error())
		bot.Send(msg)
		return
	}

	//reply_text := fmt.Sprint(glyphNames, edgeLists)
	//msg := tgbotapi.NewMessage(u.Message.Chat.ID, reply_text)
	//bot.Send(msg)

	log.Printf("User: %s %s (@%s), Text: %s", m.From.FirstName, m.From.LastName, m.From.UserName, m.Text)

	fileID := checkCache(glyphNames, db)
	if fileID == nil {
		// Create and send new
		log.Printf("Creating new image for glyphs %s", glyphNames)
		img := makeImage(glyphNames)
		webp := encodeWebp(img, tempdir)
		msg := tgbotapi.NewStickerUpload(u.Message.Chat.ID, webp)
		response, _ := bot.Send(msg)
		//spew.Dump(response)
		fileID := response.Sticker.FileID
		log.Printf("Added to cache: %s -> %s. File size: %d", glyphNames, fileID, response.Sticker.FileSize)
		addToCache(glyphNames, fileID, db)
	} else {
		log.Printf("Hitting cache! %s -> %s", glyphNames, fileID)
		msg := tgbotapi.NewStickerShare(u.Message.Chat.ID, string(fileID))
		bot.Send(msg)
	}
}
