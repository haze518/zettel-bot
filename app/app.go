package zettel_bot

import (
	// "bytes"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Serve(token string) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
			// buf := &bytes.Buffer{}
			as := GetZeroLinks()
			// for _, val := range as {
			// 	buf.WriteString(val)
			// 	buf.WriteString("\n")
			// }
			new_msg := tgbotapi.NewMessage(update.Message.Chat.ID, as)
			bot.Send(new_msg)

			md := DownloadFile()
			md_message := tgbotapi.NewMessage(update.Message.Chat.ID, md)
			md_message.ParseMode = "Markdown"
			bot.Send(md_message)
		}
	}
}