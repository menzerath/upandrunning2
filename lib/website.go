package lib

import (
	"github.com/franela/goreq"
	"github.com/mitsuse/pushbullet-go"
	"github.com/mitsuse/pushbullet-go/requests"
	"github.com/op/go-logging"
	"strconv"
	"strings"
	"time"
)

// Represents a single Website-object.
type Website struct {
	Id          int
	Protocol    string
	Url         string
	CheckMethod string
}

// Runs a check and saves the result inside the database.
func (w *Website) RunCheck(secondTry bool) {
	// Request new Status
	res, err := goreq.Request{Uri: w.Protocol + "://" + w.Url, Method: w.CheckMethod, UserAgent: "UpAndRunning2 (https://github.com/MarvinMenzerath/UpAndRunning2)", MaxRedirects: 10, Timeout: 5 * time.Second}.Do()

	var newStatus string
	var newStatusCode int
	if err != nil {
		if secondTry {
			newStatus = "Host not found"
			newStatusCode = 0

			// On Timeout: allow second try
			if serr, ok := err.(*goreq.Error); ok {
				if serr.Timeout() {
					newStatus = "Timeout"
					newStatusCode = 0
				}
			}
		} else {
			time.Sleep(time.Millisecond * 1000)
			w.RunCheck(true)
			return
		}
	} else {
		newStatus = strconv.Itoa(res.StatusCode) + " - " + GetHttpStatus(res.StatusCode)
		newStatusCode = res.StatusCode
		defer res.Body.Close()
	}

	// If Pushbullet-notifications are active: get old status and Website's name and send a Push
	if GetConfiguration().Dynamic.PushbulletKey != "" {
		var (
			name   string
			status string
		)

		db := GetDatabase()
		err = db.QueryRow("SELECT name, status FROM website WHERE id = ?", w.Id).Scan(&name, &status)
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
		_, err = db.Exec("UPDATE website SET status = ?, time = NOW(), ups = ups + 1, totalChecks = totalChecks + 1 WHERE id = ?;", newStatus, w.Id)
		if err != nil {
			logging.MustGetLogger("logger").Error("Unable to save the new Website-status: ", err)
			return
		}
	} else {
		// Failure
		_, err = db.Exec("UPDATE website SET status = ?, time = NOW(), lastFailStatus = ?, lastFailTime = NOW(), downs = downs + 1, totalChecks = totalChecks + 1 WHERE id = ?;", newStatus, newStatus, w.Id)
		if err != nil {
			logging.MustGetLogger("logger").Error("Unable to save the new Website-status: ", err)
			return
		}
	}

	w.calcAvgAvailability()
}

// Calculates the average Website availability and stores it inside the database.
func (w *Website) calcAvgAvailability() {
	// Query the Database and format the returned value
	db := GetDatabase()
	var avg float64
	err := db.QueryRow("SELECT ((SELECT ups FROM website WHERE id = ?) / (SELECT totalChecks FROM website WHERE id = ?))*100 AS avg", w.Id, w.Id).Scan(&avg)
	if err != nil {
		logging.MustGetLogger("logger").Error("Unable to calculate Website-Availability: ", err)
		return
	}
	strconv.FormatFloat(avg, 'f', 2, 64)

	// Save the new value
	_, err = db.Exec("UPDATE website SET avgAvail = ? WHERE id = ?;", avg, w.Id)
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
