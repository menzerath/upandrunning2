package lib

import (
	"github.com/mitsuse/pushbullet-go"
	"github.com/mitsuse/pushbullet-go/requests"
	"github.com/op/go-logging"
	"strconv"
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
	// Query the Database
	db := GetDatabase()
	stmt, err := db.Prepare("SELECT ((SELECT ups FROM website WHERE id = ?) / (SELECT totalChecks FROM website WHERE id = ?))*100 AS avg")
	if err != nil {
		logging.MustGetLogger("logger").Error("Unable to calculate Website-Availability: ", err)
		return
	}

	var avg float64
	err = stmt.QueryRow(w.Id, w.Id).Scan(&avg)
	if err != nil {
		logging.MustGetLogger("logger").Error("Unable to calculate Website-Availability: ", err)
		return
	}

	strconv.FormatFloat(avg, 'f', 2, 64)

	// Save the new value
	stmt, err = db.Prepare("UPDATE website SET avgAvail = ? WHERE id = ?;")
	if err != nil {
		logging.MustGetLogger("logger").Error("Unable to calculate Website-Availability: ", err)
		return
	}

	_, err = stmt.Exec(avg, w.Id)
	if err != nil {
		logging.MustGetLogger("logger").Error("Unable to save Website-Availability: ", err)
		return
	}
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

	_, err := pb.PostPushesNote(push)
	if err != nil {
		logging.MustGetLogger("logger").Error("Unable to send Push: ", err)
	}
}
