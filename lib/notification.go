package lib

import (
	"github.com/mitsuse/pushbullet-go"
	"github.com/mitsuse/pushbullet-go/requests"
	"github.com/op/go-logging"
	"gopkg.in/gomail.v2"
)

// Sends a Pushbullet-Push containing the given data to the saved API-key.
func sendPush(apiKey string, name string, url string, newStatus string, oldStatus string) {
	logging.MustGetLogger("").Debug("Sending Push about \"" + url + "\"...")

	pb := pushbullet.New(apiKey)

	push := requests.NewNote()
	push.Title = GetConfiguration().Dynamic.Title + " - Status Change"
	push.Body = name + " (" + url + ") went from \"" + oldStatus + "\" to \"" + newStatus + "\"."

	_, err := pb.PostPushesNote(push)
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
	m.SetBody("text/html", "Hello,<br /><br /><b>"+name+"</b>"+" ("+url+") just went from \""+oldStatus+"\" to <b>\""+newStatus+"\"</b>.<br /><br />Sincerely,<br />"+GetConfiguration().Dynamic.Title+"<br /><br /><small>This email was sent automatically, please do not respond to it.</small>")

	d := gomail.NewPlainDialer(mConf.Host, mConf.Port, mConf.User, mConf.Password)

	if err := d.DialAndSend(m); err != nil {
		logging.MustGetLogger("").Error("Unable to send email: ", err)
	}
}
