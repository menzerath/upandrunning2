package lib

import (
	"database/sql"
	"github.com/franela/goreq"
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
	res, err := goreq.Request{Uri: w.Protocol + "://" + w.Url, Method: w.CheckMethod, UserAgent: "UpAndRunning2/" + GetConfiguration().Static.Version + " (https://github.com/MarvinMenzerath/UpAndRunning2)", MaxRedirects: GetConfiguration().Dynamic.Redirects, Timeout: 5 * time.Second}.Do()
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

	w.sendNotifications(newStatusCode, newStatusText)

	// Save the new Result
	_, err = db.Exec("INSERT INTO checks (websiteId, statusCode, statusText, responseTime, time) VALUES (?, ?, ?, ?, NOW());", w.Id, newStatusCode, newStatusText, requestDuration)
	if err != nil {
		logging.MustGetLogger("").Error("Unable to save the new Website-status: ", err)
		return
	}
}

// Gets the notification-settings and sends a notification (if necessary and requested)
func (w *Website) sendNotifications(newStatusCode int, newStatusText string) {
	var (
		pushbulletKey string
		email         string
		name          string
		oldStatusCode string
		oldStatusText string
	)

	db := GetDatabase()
	err := db.QueryRow("SELECT pushbulletKey, email FROM notifications WHERE websiteId = ?", w.Id).Scan(&pushbulletKey, &email)
	if err != nil {
		if err == sql.ErrNoRows {
			return
		}
		logging.MustGetLogger("").Error("Unable to get Website's notification-settings: ", err)
	}

	// Check for empty result
	if pushbulletKey == "" && email == "" {
		return
	}

	err = db.QueryRow("SELECT name FROM websites WHERE id = ?", w.Id).Scan(&name)
	if err != nil {
		logging.MustGetLogger("").Error("Unable to get Website's data: ", err)
		return
	}
	err = db.QueryRow("SELECT statusCode, statusText FROM checks WHERE websiteId = ? ORDER BY id DESC LIMIT 1", w.Id).Scan(&oldStatusCode, &oldStatusText)
	switch {
	case err == sql.ErrNoRows:
		return
	case err != nil:
		logging.MustGetLogger("").Error("Unable to get Website's data: ", err)
		return
	}

	oldStatus := oldStatusCode + " - " + oldStatusText
	newStatus := strconv.Itoa(newStatusCode) + " - " + newStatusText
	if oldStatus != newStatus {
		if pushbulletKey != "" {
			sendPush(pushbulletKey, name, w.Url, newStatus, oldStatus)
		}
		if email != "" {
			sendMail(email, name, w.Url, newStatus, oldStatus)
		}
	}
}
