package lib

import (
	"github.com/franela/goreq"
	"github.com/mitsuse/pushbullet-go"
	"github.com/mitsuse/pushbullet-go/requests"
	"github.com/op/go-logging"
	"strconv"
	"time"
	"strings"
)

// Represents a single Website-object.
type Website struct {
	Id       int
	Protocol string
	Url      string
}

// Runs a check and saves the result inside the database.
func (w *Website) RunCheck() {
	// Request new Status
	res, err := goreq.Request{Uri: w.Protocol + "://" + w.Url, Method: "HEAD", UserAgent: "UpAndRunning2 (https://github.com/MarvinMenzerath/UpAndRunning2)", MaxRedirects: 10, Timeout: 10 * time.Second}.Do()

	var newStatus string
	var newStatusCode int
	if err != nil {
		logging.MustGetLogger("logger").Warning("Error while requesting Status: ", err)
		newStatus = "Server not found"
		newStatusCode = 0
	} else {
		newStatus = strconv.Itoa(res.StatusCode) + " - " + GetHttpStatus(res.StatusCode)
		newStatusCode = res.StatusCode
		defer res.Body.Close()
	}

	// If Pushbullet-notifications are active: get old status and Website's name and send a Push
	if GetConfiguration().Dynamic.PushbulletKey != "" {
		db := GetDatabase()
		stmt, err := db.Prepare("SELECT name, status FROM website WHERE id = ?")
		if err != nil {
			logging.MustGetLogger("logger").Error("Unable to get Website's data: ", err)
			return
		}

		var (
			name   string
			status string
		)
		err = stmt.QueryRow(w.Id).Scan(&name, &status)
		if err != nil {
			logging.MustGetLogger("logger").Error("Unable to get Website's data: ", err)
			return
		}

		if newStatus != status {
			sendPush(name, w.Url, newStatus, status)
		}
	}

	newStatusCodeString := strconv.Itoa(newStatusCode)

	// Save the new Result
	if strings.HasPrefix(newStatusCodeString, "2") || strings.HasPrefix(newStatusCodeString, "3") {
		// Success
		stmt, err := db.Prepare("UPDATE website SET status = ?, time = NOW(), ups = ups + 1, totalChecks = totalChecks + 1 WHERE id = ?;")
		if err != nil {
			logging.MustGetLogger("logger").Error("Unable to save the new Website-status: ", err)
			return
		}

		_, err = stmt.Exec(newStatus, w.Id)
		if err != nil {
			logging.MustGetLogger("logger").Error("Unable to save the new Website-status: ", err)
			return
		}
	} else {
		// Failure
		stmt, err := db.Prepare("UPDATE website SET status = ?, time = NOW(), lastFailStatus = ?, lastFailTime = NOW(), downs = downs + 1, totalChecks = totalChecks + 1 WHERE id = ?;")
		if err != nil {
			logging.MustGetLogger("logger").Error("Unable to save the new Website-status: ", err)
			return
		}

		_, err = stmt.Exec(newStatus, newStatus, w.Id)
		if err != nil {
			logging.MustGetLogger("logger").Error("Unable to save the new Website-status: ", err)
			return
		}
	}

	w.calcAvgAvailability()
}

// Calculates the average Website availability and stores it inside the database.
func (w *Website) calcAvgAvailability() {
	// Query the Database
	db := GetDatabase()
	stmt, err := db.Prepare("SELECT ((SELECT ups FROM website WHERE id = ?) / (SELECT totalChecks FROM website WHERE id = ?))*100 AS avg")
	if err != nil {
		logging.MustGetLogger("logger").Error("Unable to calculate Website-Availability: ", err)
		return
	}

	// Format the returned value
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

// Sends a Pushbullet-Push containing the given data to the saved API-key.
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
