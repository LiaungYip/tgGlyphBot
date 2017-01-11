package main

import (
	"github.com/LiaungYip/tgGlyphBot/cache"
	"github.com/LiaungYip/tgGlyphBot/config"
	"github.com/LiaungYip/tgGlyphBot/input"
	"github.com/boltdb/bolt"
	"github.com/davecgh/go-spew/spew"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"os"
	"strings"
)

func main() {
	defer LogSetupAndDestruct()()

	bot, updates := initTelegram()
	db := cache.Init(config.DatabaseFilename)
	defer db.Close()

	tempdir := tempImageDir()
	defer os.RemoveAll(tempdir)

	for u := range updates {
		log.Println(spew.Sdump(u))
		if u.Message == nil {
			continue
		}
		m := u.Message
		log.Printf("User: %s %s (@%s), Text: %s", m.From.FirstName, m.From.LastName, m.From.UserName, m.Text)
		handleUpdate(bot, u, db, tempdir)
	}
}

func handleUpdate(bot *tgbotapi.BotAPI, u tgbotapi.Update, db *bolt.DB, tempdir string) {
	log.Println("Handling update")
	m := u.Message
	if m.Text == "" { // No action if text field is blank. (i.e. if user sends photos or stickers or voice or something weird.)
		return
	}

	if isCommand(m, "help") || isCommand(m, "start") {
		log.Println("Sent help message!")
		sendHelp(bot, u)
		return
	}

	if isCommand(m, "glyph") || isCommand(m, "glyphs") || m.Chat.IsPrivate() && isCommand(m, "") {
		glyphsCommand(bot, u, db, tempdir)
		return
	}
}

func isCommand(m *tgbotapi.Message, commandName string) bool {
	if strings.EqualFold(commandName, m.Command()) {
		return true
	} else {
		return false
	}
}

func glyphsCommand(bot *tgbotapi.BotAPI, u tgbotapi.Update, db *bolt.DB, tempdir string) {
	var t string
	if u.Message.Command() == "" {
		t = u.Message.Text
	} else {
		t = u.Message.CommandArguments()
	}

	glyphNames, _, err := input.ProcessString(t)
	if err != nil {
		sendError(bot, u, err)
		return
	}

	//log.Printf("User: %s %s (@%s), Text: %s", m.From.FirstName, m.From.LastName, m.From.UserName, m.Text)

	fileID := cache.Check(glyphNames, db)
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
	cache.Add(glyphNames, fileID, db)
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
