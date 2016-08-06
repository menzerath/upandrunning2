package APIv2

import (
	"database/sql"
	"encoding/json"
	"github.com/MarvinMenzerath/UpAndRunning2/lib"
	"github.com/asaskevich/govalidator"
	"github.com/julienschmidt/httprouter"
	"github.com/op/go-logging"
	"net/http"
)

// Inserts a new Website into the database.
func ApiWebsitesAdd(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !lib.IsLoggedIn(r) {
		SendJsonMessage(w, http.StatusUnauthorized, false, "Unauthorized.")
		return
	}

	// Get data from Request
	r.ParseForm()
	name := r.Form.Get("name")
	protocol := r.Form.Get("protocol")
	url := ps.ByName("url")
	method := r.Form.Get("checkMethod")

	// Simple Validation
	if name == "" || protocol == "" || url == "" || method == "" {
		SendJsonMessage(w, http.StatusBadRequest, false, "Unable to process your Request: Submit valid values.")
		return
	}
	if protocol != "http" && protocol != "https" {
		SendJsonMessage(w, http.StatusBadRequest, false, "Unable to process your Request: Submit a valid protocol.")
		return
	}
	if !govalidator.IsURL(protocol + "://" + url) {
		SendJsonMessage(w, http.StatusBadRequest, false, "Unable to process your Request: Submit a valid url.")
		return
	}
	if method != "HEAD" && method != "GET" {
		SendJsonMessage(w, http.StatusBadRequest, false, "Unable to process your Request: Submit a valid check method.")
		return
	}

	// Insert into Database
	db := lib.GetDatabase()
	_, err := db.Exec("INSERT INTO websites (name, protocol, url, checkMethod) VALUES (?, ?, ?, ?);", name, protocol, url, method)
	if err != nil {
		logging.MustGetLogger("").Error("Unable to add Website: ", err)
		SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to process your Request: "+err.Error())
		return
	}

	SendJsonMessage(w, http.StatusOK, true, "")
}

// Edits an existing Website in the database.
func ApiWebsitesEdit(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !lib.IsLoggedIn(r) {
		SendJsonMessage(w, http.StatusUnauthorized, false, "Unauthorized.")
		return
	}

	// Get data from Request
	r.ParseForm()
	oldUrl := ps.ByName("url")
	name := r.Form.Get("name")
	protocol := r.Form.Get("protocol")
	url := r.Form.Get("url")
	method := r.Form.Get("checkMethod")

	// Simple Validation
	if oldUrl == "" || name == "" || protocol == "" || url == "" || method == "" {
		SendJsonMessage(w, http.StatusBadRequest, false, "Unable to process your Request: Submit valid values.")
		return
	}
	if protocol != "http" && protocol != "https" {
		SendJsonMessage(w, http.StatusBadRequest, false, "Unable to process your Request: Submit a valid protocol.")
		return
	}
	if !govalidator.IsURL(protocol + "://" + url) {
		SendJsonMessage(w, http.StatusBadRequest, false, "Unable to process your Request: Submit a valid url.")
		return
	}
	if method != "HEAD" && method != "GET" {
		SendJsonMessage(w, http.StatusBadRequest, false, "Unable to process your Request: Submit a valid check method.")
		return
	}

	// Update Database
	db := lib.GetDatabase()
	res, err := db.Exec("UPDATE websites SET name = ?, protocol = ?, url = ?, checkMethod = ? WHERE url = ?;", name, protocol, url, method, oldUrl)
	if err != nil {
		logging.MustGetLogger("").Error("Unable to edit Website: ", err)
		SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to process your Request: "+err.Error())
		return
	}

	// Check if exactly one Website has been edited
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 1 {
		SendJsonMessage(w, http.StatusOK, true, "")
	} else {
		SendJsonMessage(w, http.StatusBadRequest, false, "Unable to process your Request: Could not edit Website.")
	}
}

// Removes an existing Website from the database.
func ApiWebsitesDelete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !lib.IsLoggedIn(r) {
		SendJsonMessage(w, http.StatusUnauthorized, false, "Unauthorized.")
		return
	}

	// Get data from Request
	value := ps.ByName("url")

	// Simple Validation
	if value == "" {
		SendJsonMessage(w, http.StatusBadRequest, false, "Unable to process your Request: Submit a valid value.")
		return
	}

	// Remove Check-Results from Database
	db := lib.GetDatabase()
	res, err := db.Exec("DELETE c FROM checks c INNER JOIN websites w ON c.websiteId = w.id WHERE w.url = ?;", value)
	if err != nil {
		logging.MustGetLogger("").Error("Unable to delete Check-Results: ", err)
		SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to process your Request: "+err.Error())
		return
	}

	// Remove Notifications from Database
	res, err = db.Exec("DELETE n FROM notifications n INNER JOIN websites w ON n.websiteId = w.id WHERE w.url = ?;", value)
	if err != nil {
		logging.MustGetLogger("").Error("Unable to delete Notifications: ", err)
		SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to process your Request: "+err.Error())
		return
	}

	// Remove Website from Database
	db = lib.GetDatabase()
	res, err = db.Exec("DELETE FROM websites WHERE url = ?;", value)
	if err != nil {
		logging.MustGetLogger("").Error("Unable to delete Website: ", err)
		SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to process your Request: "+err.Error())
		return
	}

	// Check if exactly one Website has been deleted
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 1 {
		SendJsonMessage(w, http.StatusOK, true, "")
	} else {
		SendJsonMessage(w, http.StatusBadRequest, false, "Unable to process your Request: Could not delete Website.")
	}
}

// Enables / Disables an existing Website in the database.
func ApiWebsitesEnabled(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !lib.IsLoggedIn(r) {
		SendJsonMessage(w, http.StatusUnauthorized, false, "Unauthorized.")
		return
	}

	// Get data from Request
	r.ParseForm()
	value := ps.ByName("url")
	enabled := r.Form.Get("enabled")

	// Simple Validation
	if value == "" || enabled == "" {
		SendJsonMessage(w, http.StatusBadRequest, false, "Unable to process your Request: Submit valid values.")
		return
	}

	var enabledValue int
	if enabled == "true" {
		enabledValue = 1
	} else if enabled == "false" {
		enabledValue = 0
	} else {
		SendJsonMessage(w, http.StatusBadRequest, false, "Unable to process your Request: Submit valid values.")
		return
	}

	// Update Database-Row
	db := lib.GetDatabase()
	res, err := db.Exec("UPDATE websites SET enabled = ? WHERE url = ?;", enabledValue, value)
	if err != nil {
		logging.MustGetLogger("").Error("Unable to enable / disable Website: ", err)
		SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to process your Request: "+err.Error())
		return
	}

	// Check if exactly one Website is affected
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 1 {
		SendJsonMessage(w, http.StatusOK, true, "")
	} else {
		SendJsonMessage(w, http.StatusBadRequest, false, "Unable to process your Request: Could not enable / disable Website.")
	}
}

// Sets an existing Website to visible / invisible in the database.
func ApiWebsitesVisibility(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !lib.IsLoggedIn(r) {
		SendJsonMessage(w, http.StatusUnauthorized, false, "Unauthorized.")
		return
	}

	// Get data from Request
	r.ParseForm()
	value := ps.ByName("url")
	visible := r.Form.Get("visible")

	// Simple Validation
	if value == "" || visible == "" {
		SendJsonMessage(w, http.StatusBadRequest, false, "Unable to process your Request: Submit valid values.")
		return
	}

	var visibilityValue int
	if visible == "true" {
		visibilityValue = 1
	} else if visible == "false" {
		visibilityValue = 0
	} else {
		SendJsonMessage(w, http.StatusBadRequest, false, "Unable to process your Request: Submit valid values.")
		return
	}

	// Update Database-Row
	db := lib.GetDatabase()
	res, err := db.Exec("UPDATE websites SET visible = ? WHERE url = ?;", visibilityValue, value)
	if err != nil {
		logging.MustGetLogger("").Error("Unable to set Website's visibility: ", err)
		SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to process your Request: "+err.Error())
		return
	}

	// Check if exactly one Website is affected
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 1 {
		SendJsonMessage(w, http.StatusOK, true, "")
	} else {
		SendJsonMessage(w, http.StatusBadRequest, false, "Unable to process your Request: Could not set Website's visibility.")
	}
}

// Get an existing Website's notification-preferences.
func ApiWebsitesGetNotifications(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !lib.IsLoggedIn(r) {
		SendJsonMessage(w, http.StatusUnauthorized, false, "Unauthorized.")
		return
	}

	// Get data from Request
	value := ps.ByName("url")

	// Simple Validation
	if value == "" {
		SendJsonMessage(w, http.StatusBadRequest, false, "Unable to process your Request: Submit a valid value.")
		return
	}

	// Get existing settings
	var (
		cPushbulletKey string
		cEmail         string
	)
	db := lib.GetDatabase()
	var resp WebsiteNotificationsResponse
	err := db.QueryRow("SELECT pushbulletKey, email FROM notifications, websites WHERE notifications.websiteId = websites.id AND url = ?;", value).Scan(&cPushbulletKey, &cEmail)
	if err != nil {
		if err == sql.ErrNoRows {
			// Check if Website exists
			var id int
			err := db.QueryRow("SELECT id FROM websites WHERE url = ?;", value).Scan(&id)
			if err != nil {
				SendJsonMessage(w, http.StatusNotFound, false, "Unable to process your Request: Could not find Website.")
				return
			}
			resp = WebsiteNotificationsResponse{true, Notifications{"", ""}}
		} else {
			logging.MustGetLogger("").Error("Unable to get Website's notification settings: ", err)
			SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to process your Request: "+err.Error())
			return
		}
	} else {
		resp = WebsiteNotificationsResponse{true, Notifications{cPushbulletKey, cEmail}}
	}

	// Send Response
	responseBytes, err := json.Marshal(resp)
	if err != nil {
		SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to process your Request.")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBytes)
}

// Sets an existing Website's notification-preferences.
func ApiWebsitePutNotifications(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !lib.IsLoggedIn(r) {
		SendJsonMessage(w, http.StatusUnauthorized, false, "Unauthorized.")
		return
	}

	// Get data from Request
	r.ParseForm()
	url := ps.ByName("url")
	pushbulletKey := r.Form.Get("pushbulletKey")
	email := r.Form.Get("email")

	// Simple Validation
	if url == "" {
		SendJsonMessage(w, http.StatusBadRequest, false, "Unable to process your Request: Submit a valid value.")
		return
	}

	// Check for existing settings
	var (
		cPushbulletKey string
		cEmail         string
	)
	db := lib.GetDatabase()
	err := db.QueryRow("SELECT pushbulletKey, email FROM notifications, websites WHERE notifications.websiteId = websites.id AND url = ?;", url).Scan(&cPushbulletKey, &cEmail)
	if err != nil {
		if err == sql.ErrNoRows {
			// no settings found --> Insert
			var id int
			err := db.QueryRow("SELECT id FROM websites WHERE url = ?;", url).Scan(&id)
			if err != nil {
				SendJsonMessage(w, http.StatusNotFound, false, "Unable to process your Request: Could not find Website.")
				return
			}

			_, err = db.Exec("INSERT INTO notifications (websiteId, pushbulletKey, email) VALUES (?, ?, ?);", id, pushbulletKey, email)
			if err != nil {
				logging.MustGetLogger("").Error("Unable to insert Website's notification settings: ", err)
				SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to process your Request: "+err.Error())
				return
			}
		} else {
			logging.MustGetLogger("").Error("Unable to get Website's notification settings: ", err)
			SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to process your Request: "+err.Error())
			return
		}
	} else {
		// existing settings found --> Update
		_, err = db.Exec("UPDATE notifications, websites SET pushbulletKey = ?, email = ? WHERE notifications.websiteId = websites.id AND url = ?;", pushbulletKey, email, url)
		if err != nil {
			logging.MustGetLogger("").Error("Unable to update Website's notification settings: ", err)
			SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to process your Request: "+err.Error())
			return
		}
	}

	SendJsonMessage(w, http.StatusOK, true, "")
}

// Triggers a check of all enabled Websites.
func ApiWebsiteCheck(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !lib.IsLoggedIn(r) {
		SendJsonMessage(w, http.StatusUnauthorized, false, "Unauthorized.")
		return
	}

	// Get data from Request
	r.ParseForm()
	url := ps.ByName("url")

	// Simple Validation
	if url == "" {
		SendJsonMessage(w, http.StatusBadRequest, false, "Unable to process your Request: Submit a valid value.")
		return
	}

	// Query the Database
	db := lib.GetDatabase()
	var website lib.Website
	err := db.QueryRow("SELECT id, protocol, url, checkMethod FROM websites WHERE enabled = 1 AND url = ?;", url).Scan(&website.Id, &website.Protocol, &website.Url, &website.CheckMethod)

	if err != nil {
		if err == sql.ErrNoRows {
			SendJsonMessage(w, http.StatusNotFound, false, "Unable to process your Request: Could not find Website.")
			return
		} else {
			logging.MustGetLogger("").Error("Unable to get Website's notification settings: ", err)
			SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to process your Request: "+err.Error())
			return
		}
	} else {
		// Run the requested check
		logging.MustGetLogger("").Info("Checking requested Website (" + website.Url + ").")
		website.RunCheck(false)
		SendJsonMessage(w, http.StatusOK, true, "")
	}
}
