package lib

import (
	"github.com/mitsuse/pushbullet-go"
	"github.com/mitsuse/pushbullet-go/requests"
	"github.com/op/go-logging"
)

// Sends a Pushbullet-Push containing the given data to the saved API-key.
func sendPush(apiKey string, name string, url string, newStatus string, oldStatus string) {
	logging.MustGetLogger("logger").Debug("Sending Push about \"" + url + "\"...")

	pb := pushbullet.New(apiKey)

	push := requests.NewNote()
	push.Title = GetConfiguration().Dynamic.Title + " - Status Change"
	push.Body = name + " (" + url + ") went from \"" + oldStatus + "\" to \"" + newStatus + "\"."

	_, err := pb.PostPushesNote(push)
	if err != nil {
		logging.MustGetLogger("logger").Error("Unable to send Push: ", err)
	}
}

// Sends an e-mail containing the given data to the saved e-mail-address.
// Needs a configured SMTP-server (in config-file).
func sendMail() {
	// TODO
}
