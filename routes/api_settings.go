package routes

import (
	"github.com/MarvinMenzerath/UpAndRunning2/lib"
	"github.com/julienschmidt/httprouter"
	"github.com/op/go-logging"
	"net/http"
	"strconv"
)

// Updates the application's title in the database.
func ApiSettingsTitle(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !lib.IsLoggedIn(r) {
		SendJsonMessage(w, http.StatusUnauthorized, false, "Unauthorized.")
		return
	}

	// Get data from Request
	r.ParseForm()
	value := r.Form.Get("title")

	// Simple Validation
	if value == "" {
		SendJsonMessage(w, http.StatusBadRequest, false, "Unable to process your Request: Submit a valid value.")
		return
	}

	// Update Database-Row
	db := lib.GetDatabase()
	_, err := db.Exec("UPDATE settings SET value = ? WHERE name = 'title';", value)
	if err != nil {
		logging.MustGetLogger("").Error("Unable to change Application-Title: ", err)
		SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to process your Request: "+err.Error())
		return
	}

	// Update Configuration
	lib.GetConfiguration().Dynamic.Title = value
	SendJsonMessage(w, http.StatusOK, true, "")
}

// Updates the user's password in the database.
func ApiSettingsPassword(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !lib.IsLoggedIn(r) {
		SendJsonMessage(w, http.StatusUnauthorized, false, "Unauthorized.")
		return
	}

	// Get data from Request
	r.ParseForm()
	value := r.Form.Get("password")

	// Simple Validation
	if value == "" {
		SendJsonMessage(w, http.StatusBadRequest, false, "Unable to process your Request: Submit a valid value.")
		return
	}

	// Update Password
	admin := lib.Admin{}
	err := admin.ChangePassword(value)
	if err != nil {
		SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to process your Request: "+err.Error())
		return
	}
	SendJsonMessage(w, http.StatusOK, true, "")
}

// Updates the application's check-interval in the database.
func ApiSettingsInterval(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !lib.IsLoggedIn(r) {
		SendJsonMessage(w, http.StatusUnauthorized, false, "Unauthorized.")
		return
	}

	// Get data from Request
	r.ParseForm()
	temp := r.Form.Get("interval")
	value, err := strconv.Atoi(temp)

	// Simple Validation
	if err != nil || value < 10 || value > 600 {
		SendJsonMessage(w, http.StatusBadRequest, false, "Unable to process your Request: Submit a valid value between 10 and 600 seconds.")
		return
	}

	// Update Database-Row
	db := lib.GetDatabase()
	_, err = db.Exec("UPDATE settings SET value = ? WHERE name = 'interval';", value)
	if err != nil {
		logging.MustGetLogger("").Error("Unable to change Interval: ", err)
		SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to process your Request: "+err.Error())
		return
	}

	// Update Configuration
	lib.GetConfiguration().Dynamic.Interval = value
	SendJsonMessage(w, http.StatusOK, true, "")
}

// Updates the application's maximum amount of redirects in the database.
func ApiSettingsRedirects(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !lib.IsLoggedIn(r) {
		SendJsonMessage(w, http.StatusUnauthorized, false, "Unauthorized.")
		return
	}

	// Get data from Request
	r.ParseForm()
	temp := r.Form.Get("redirects")
	value, err := strconv.Atoi(temp)

	// Simple Validation
	if err != nil || value < 0 || value > 10 {
		SendJsonMessage(w, http.StatusBadRequest, false, "Unable to process your Request: Submit a valid value between 0 and 10 redirects.")
		return
	}

	// Update Database-Row
	db := lib.GetDatabase()
	_, err = db.Exec("UPDATE settings SET value = ? WHERE name = 'redirects';", value)
	if err != nil {
		logging.MustGetLogger("").Error("Unable to change Redirects: ", err)
		SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to process your Request: "+err.Error())
		return
	}

	// Update Configuration
	lib.GetConfiguration().Dynamic.Redirects = value
	SendJsonMessage(w, http.StatusOK, true, "")
}

// Updates whether the application should run checks when there is not internet-connection or not.
func ApiSettingsCheckWhenOffline(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !lib.IsLoggedIn(r) {
		SendJsonMessage(w, http.StatusUnauthorized, false, "Unauthorized.")
		return
	}

	// Get data from Request
	r.ParseForm()
	value := r.Form.Get("checkWhenOffline")

	// Simple Validation
	if value == "" {
		SendJsonMessage(w, http.StatusBadRequest, false, "Unable to process your Request: Submit a valid value.")
		return
	}

	var enabledValue int
	if value == "true" {
		enabledValue = 1
	} else if value == "false" {
		enabledValue = 0
	} else {
		SendJsonMessage(w, http.StatusBadRequest, false, "Unable to process your Request: Submit a valid value (true or false).")
		return
	}

	// Update Database-Row
	db := lib.GetDatabase()
	_, err := db.Exec("UPDATE settings SET value = ? WHERE name = 'check_when_offline';", enabledValue)
	if err != nil {
		logging.MustGetLogger("").Error("Unable to change Offline-Checking: ", err)
		SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to process your Request: "+err.Error())
		return
	}

	// Update Configuration
	lib.GetConfiguration().Dynamic.RunChecksWhenOffline = enabledValue
	SendJsonMessage(w, http.StatusOK, true, "")
}
