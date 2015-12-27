package lib

import (
	"github.com/franela/goreq"
	"github.com/mitsuse/pushbullet-go"
	"github.com/mitsuse/pushbullet-go/requests"
	"github.com/op/go-logging"
	"strconv"
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
	var requestStartTime = time.Now()
	res, err := goreq.Request{Uri: w.Protocol + "://" + w.Url, Method: w.CheckMethod, UserAgent: "UpAndRunning2 (https://github.com/MarvinMenzerath/UpAndRunning2)", MaxRedirects: GetConfiguration().Dynamic.Redirects, Timeout: 5 * time.Second}.Do()
	var requestDuration = time.Now().Sub(requestStartTime).Nanoseconds() / 1000000

	var newStatusText string
	var newStatusCode int
	if err != nil {
		if secondTry {
			newStatusText = "Host not found"
			newStatusCode = 0

			// On Timeout: allow second try
			if serr, ok := err.(*goreq.Error); ok {
				if serr.Timeout() {
					newStatusText = "Timeout"
					newStatusCode = 0
				}
			}
		} else {
			time.Sleep(time.Millisecond * 1000)
			w.RunCheck(true)
			return
		}
	} else {
		newStatusText = GetHttpStatus(res.StatusCode)
		newStatusCode = res.StatusCode
		defer res.Body.Close()
	}

	// If Pushbullet-notifications are active: get old status and Website's name and send a Push
	if GetConfiguration().Dynamic.PushbulletKey != "" {
		var (
			name          string
			oldStatusCode string
			oldStatusText string
		)

		db := GetDatabase()
		noError := true
		err = db.QueryRow("SELECT name FROM websites WHERE id = ?", w.Id).Scan(&name)
		if err != nil {
			logging.MustGetLogger("logger").Error("Unable to get Website's data: ", err)
			noError = false
		}
		err = db.QueryRow("SELECT statusCode, statusText FROM checks WHERE websiteId = ? ORDER BY id DESC LIMIT 1", w.Id).Scan(&oldStatusCode, &oldStatusText)
		if err != nil {
			logging.MustGetLogger("logger").Warning("Unable to get Website's data: ", err)
			logging.MustGetLogger("logger").Warning("This is totally normal if this is the first result inserted into the database.")
			noError = false
		}

		oldStatus := oldStatusCode + " - " + oldStatusText
		newStatus := strconv.Itoa(newStatusCode) + " - " + newStatusText
		if oldStatus != newStatus && noError {
			sendPush(name, w.Url, newStatus, oldStatus)
		}
	}

	// Save the new Result
	_, err = db.Exec("INSERT INTO checks (websiteId, statusCode, statusText, responseTime, time) VALUES (?, ?, ?, ?, NOW());", w.Id, newStatusCode, newStatusText, requestDuration)
	if err != nil {
		logging.MustGetLogger("logger").Error("Unable to save the new Website-status: ", err)
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
