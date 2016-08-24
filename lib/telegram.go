package lib

import (
	"github.com/op/go-logging"
	"github.com/tucnak/telebot"
	"strconv"
	"time"
)

var Bot *telebot.Bot
var SendOptions = telebot.SendOptions{DisableWebPagePreview: true, ParseMode: telebot.ModeMarkdown}

func RunTelegramBot() {
	if GetConfiguration().Notification.TelegramBotApiKey == "" {
		return
	}

	bot, err := telebot.NewBot(GetConfiguration().Notification.TelegramBotApiKey)
	if err != nil {
		logging.MustGetLogger("").Error("Unable to start Telegram-Bot: ", err)
		return
	}
	logging.MustGetLogger("").Info("Telgram-Bot started.")
	Bot = bot

	messages := make(chan telebot.Message)
	bot.Listen(messages, 1*time.Second)

	for message := range messages {
		if message.Text == "/start" {
			bot.SendMessage(message.Chat, "Welcome to the UpAndRunning2 Telegram-Bot! \U0001F44B\n\nPlease use your User-ID (`"+strconv.Itoa(message.Sender.ID)+"`) as notification-target in UpAndRunning2.", &SendOptions)
		} else if message.Text == "/id" {
			bot.SendMessage(message.Chat, "Your User-ID: `"+strconv.Itoa(message.Sender.ID)+"`", &SendOptions)
		} else if message.Text == "/server" {
			bot.SendMessage(message.Chat, "This server is running *UpAndRunning2 v"+GetConfiguration().Static.Version+"* (`"+GetConfiguration().Static.GoVersion+"@"+GetConfiguration().Static.GoArch+"`).\n"+
				"You can find more information about this application [here](https://github.com/MarvinMenzerath/UpAndRunning2).\n\n"+
				"*Server-Settings:*\n    - Title: `"+GetConfiguration().Application.Title+"`\n    - Interval: `"+strconv.Itoa(GetConfiguration().Dynamic.Interval)+"`\n"+
				"    - Redirects: `"+strconv.Itoa(GetConfiguration().Application.RedirectsToFollow)+"`\n    - Offline-Checks: `"+strconv.FormatBool(GetConfiguration().Application.RunCheckIfOffline)+"`", &SendOptions)
		}
	}
}
