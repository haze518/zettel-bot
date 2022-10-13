package zettel_bot

import (
	"strings"
	"log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func CommandHandler(
	bot *tgbotapi.BotAPI,
	update *tgbotapi.Update,
	dbAuth *DropboxAuth,
	app *App,
) {
	switch update.Message.Command() {
	case "auth":
		msg := dbAuth.getAuthorizationURLMessage()
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, msg))
	case "access_code":
		msg := strings.Split(update.Message.Text, " ")
		if len(msg) > 1 {
			resp, err := dbAuth.getAccessToken(msg[1])
			if err != nil {
				log.Panic(err)
			}
			app.Token = resp.AccessToken
			*app.Client = NewDropboxClient(app.Token)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "You have succesfully activate token")
			bot.Send(msg)
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Add access code")
			bot.Send(msg)
		}
	case "zero_links":
		zLinks, err := app.Storage.RedisClient.SMembers("zero_links").Result()
		if err != nil {
			log.Println("Error while loading zero links")
		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "ZeroLinks:")
		msg.ReplyMarkup = getKeyboardButton(zLinks)
		bot.Send(msg)
	}
}

func getKeyboardButton(rows []string) *tgbotapi.InlineKeyboardMarkup{
	var result [][]tgbotapi.InlineKeyboardButton
	for _, r := range rows {
		signle_row := []tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData(r, r)}
		result = append(result, signle_row)
	}
	numberKeyboard := tgbotapi.NewInlineKeyboardMarkup(result...)
	return &numberKeyboard
}
