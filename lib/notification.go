package lib

import (
	"github.com/mitsuse/pushbullet-go"
	"github.com/mitsuse/pushbullet-go/requests"
	"github.com/op/go-logging"
	"github.com/tucnak/telebot"
	"gopkg.in/gomail.v2"
	"strings"
	"time"
)

// Sends a Pushbullet-Push containing the given data to the saved API-key.
func sendPush(apiKey string, name string, url string, newStatus string, oldStatus string) {
	logging.MustGetLogger("").Debug("Sending Push about \"" + url + "\"...")

	pb := pushbullet.New(apiKey)

	push := requests.NewLink()
	push.Title = GetConfiguration().Dynamic.Title + " - Status Change"
	push.Body = name + " went from \"" + oldStatus + "\" to \"" + newStatus + "\"."
	push.Url = url

	_, err := pb.PostPushesLink(push)
	if err != nil {
		logging.MustGetLogger("").Error("Unable to send Push: ", err)
	}
}

// Sends an email containing the given data to the saved e-mail-address.
// Needs a configured SMTP-server (in config-file).
func sendMail(recipient string, name string, url string, newStatus string, oldStatus string) {
	if GetConfiguration().Mailer.Host == "" || GetConfiguration().Mailer.Host == "smtp.mymail.com" {
		logging.MustGetLogger("").Warning("Not sending email because of missing configuration.")
		return
	}

	logging.MustGetLogger("").Debug("Sending email about \"" + url + "\"...")

	mConf := GetConfiguration().Mailer

	m := gomail.NewMessage()
	m.SetAddressHeader("From", mConf.From, GetConfiguration().Dynamic.Title)
	m.SetHeader("To", recipient)
	m.SetHeader("Subject", "Status Change: "+name)
	m.SetBody("text/html", "Hello,<br /><br />"+
		"<b>"+name+"</b>"+" ("+url+") went on <b>"+time.Now().Format("02.01.2006")+"</b> at <b>"+time.Now().Format("15:04:05")+"</b> from \""+oldStatus+"\" to \"<b>"+newStatus+"</b>\".<br /><br />"+
		"Sincerely,<br />"+GetConfiguration().Dynamic.Title+"<br /><br />"+
		"<small>This email was sent automatically, please do not respond to it.</small>")

	d := gomail.NewPlainDialer(mConf.Host, mConf.Port, mConf.User, mConf.Password)

	if err := d.DialAndSend(m); err != nil {
		logging.MustGetLogger("").Error("Unable to send email: ", err)
	}
}

func sendTelegramMessage(userId int, name string, url string, newStatus string, oldStatus string) {
	if GetConfiguration().TelegramBotApiKey == "" {
		logging.MustGetLogger("").Warning("Not sending Telegram-message because of missing configuration.")
		return
	}

	logging.MustGetLogger("").Debug("Sending Telegram-message about \"" + url + "\"...")

	statusEmoji := "\U00002757"
	if strings.HasPrefix(newStatus, "2") {
		statusEmoji = "\U00002714"
	} else if strings.HasPrefix(newStatus, "3") {
		statusEmoji = "\U000026A0"
	} else if strings.HasPrefix(newStatus, "4") || strings.HasPrefix(newStatus, "5") {
		statusEmoji = "\U0000274C"
	}

	Bot.SendMessage(telebot.User{ID: userId}, "*Status Change: "+statusEmoji+" "+name+"*\n`"+url+"` went from `"+oldStatus+"` to `"+newStatus+"`.", &SendOptions)
}
