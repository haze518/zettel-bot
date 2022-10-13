package zettel_bot

import (
	"log"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)


type App struct {
	Client *DBClient
	Token string
	Storage *Storage
}


func Serve(token string, app *App) {
	os.Setenv("APP_KEY", "0hhdn7ckg30joxw")
	os.Setenv("APP_SECRET", "rp4khokiuqkohl9")
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	setCommands := tgbotapi.NewSetMyCommands(
		tgbotapi.BotCommand{
			Command:     "auth",
			Description: "Get access token",
		},
		tgbotapi.BotCommand{
			Command:     "access_code",
			Description: "Register token",
		},
		tgbotapi.BotCommand{
			Command:     "zero_links",
			Description: "Get zero links",
		},
	)
	if _, err := bot.Request(setCommands); err != nil {
		log.Fatal("Error while try to show commands")
	}
	var dbAuth *DropboxAuth = New()
	go SaveNotes(app.Client, app.Storage)
	go SaveZeroLinks(app.Client, app.Storage)
	for update := range updates {
		if update.Message != nil {
			CommandHandler(bot, &update, dbAuth, app)
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		}
		if update.CallbackQuery != nil {
			zLink := update.CallbackQuery.Data
			files := getNotesByZeroLink(zLink, app.Storage)
			for _, file := range files {
				if len(file) > 0 {
					msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, file)
					msg.ParseMode = "Markdown"
					bot.Send(msg)
				}
			}
		}
	}
}

func getNotesByZeroLink(zLink string, storage *Storage) []string {
	allFiles, err := storage.RedisClient.SMembers("notes").Result()
	if err != nil {
		log.Println("No data in redis")
	}
	files := make([]string, 0, len(allFiles))
	for _, file := range allFiles {
		if strings.Contains(file, zLink) {
			files = append(files, file)
		}
	}
	return files
}
