package lib
import (
	"github.com/mitsuse/pushbullet-go"
	"github.com/mitsuse/pushbullet-go/requests"
	"github.com/op/go-logging"
)

type Website struct {
	Id       int
	Protocol string
	Url      string
}

func (w *Website) RunCheck() {
	w.calcAvgAvailability()
}

func (w *Website) calcAvgAvailability() {

}

func sendPush(name string, url string, newStatus string, oldStatus string) {
	if GetConfiguration().Dynamic.PushbulletKey == "" {
		return
	}

	logging.MustGetLogger("logger").Debug("Sending Push about \"" + url + "\"...")

	pb := pushbullet.New(GetConfiguration().Dynamic.PushbulletKey)

	push := requests.NewNote()
	push.Title = GetConfiguration().Dynamic.Title + " - Status Change"
	push.Body = name + " (" + url + ") went from \"" + oldStatus + "\" to \"" + newStatus + "\"."

	_, err := pb.PostPushesNote(push);
	if err != nil {
		logging.MustGetLogger("logger").Error("Unable to send Push: ", err)
	}
}