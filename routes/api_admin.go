package routes

import (
	"database/sql"
	"encoding/json"
	"github.com/MarvinMenzerath/UpAndRunning2/lib"
	"github.com/julienschmidt/httprouter"
	"github.com/op/go-logging"
	"net/http"
	"strconv"
)

// Returns a AdminWebsiteResponse containing all Websites as AdminWebsite.
func ApiAdminWebsites(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !lib.IsLoggedIn(r) {
		SendJsonMessage(w, http.StatusUnauthorized, false, "Unauthorized.")
		return
	}

	// Query the Database for basic data
	db := lib.GetDatabase()
	rows, err := db.Query("SELECT id, name, enabled, visible, protocol, url, checkMethod FROM websites ORDER BY name;")
	if err != nil {
		logging.MustGetLogger("logger").Error("Unable to fetch Websites: ", err)
		SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to process your Request.")
		return
	}
	defer rows.Close()

	// Add every Website
	websites := []AdminWebsite{}
	var (
		id          int
		name        string
		enabled     bool
		visible     bool
		protocol    string
		url         string
		checkMethod string
		statusCode  string
		statusText  string
		time        string
	)

	for rows.Next() {
		err = rows.Scan(&id, &name, &enabled, &visible, &protocol, &url, &checkMethod)
		if err != nil {
			logging.MustGetLogger("logger").Error("Unable to read Website-Data-Row: ", err)
			SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to process your Request.")
			return
		}

		// Query the database for status data
		err = db.QueryRow("SELECT statusCode, statusText, time FROM checks WHERE websiteId = ? ORDER BY id DESC LIMIT 1;", id).Scan(&statusCode, &statusText, &time)
		switch {
		case err == sql.ErrNoRows:
			statusCode = "0"
			statusText = "unknown"
			time = "0000-00-00 00:00:00"
		case err != nil:
			logging.MustGetLogger("logger").Error("Unable to fetch Website's status: ", err)
			SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to process your Request.")
			return
		}

		websites = append(websites, AdminWebsite{id, name, enabled, visible, protocol, url, checkMethod, statusCode + " - " + statusText, time})
	}

	// Send Response
	responseBytes, err := json.Marshal(AdminWebsiteResponse{true, websites})
	if err != nil {
		SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to process your Request.")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBytes)
}

// Inserts a new Website into the database.
func ApiAdminWebsiteAdd(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	// Insert into Database
	db := lib.GetDatabase()
	_, err := db.Exec("INSERT INTO websites (name, protocol, url, checkMethod) VALUES (?, ?, ?, ?);", name, protocol, url, method)
	if err != nil {
		logging.MustGetLogger("logger").Error("Unable to add Website: ", err)
		SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to process your Request: "+err.Error())
		return
	}

	SendJsonMessage(w, http.StatusOK, true, "")
}

// Enables an existing Website in the database.
func ApiAdminWebsiteEnabled(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
	if (enabled == "true") {
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
		logging.MustGetLogger("logger").Error("Unable to enable / disable Website: ", err)
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

// Sets an existing Website to visible in the database.
func ApiAdminWebsiteVisibility(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
	if (visible == "true") {
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
		logging.MustGetLogger("logger").Error("Unable to set Website's visibility: ", err)
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
func ApiAdminWebsiteGetNotifications(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
	var resp AdminWebsiteNotificationsResponse
	err := db.QueryRow("SELECT pushbulletKey, email FROM notifications, websites WHERE notifications.websiteId = websites.id AND url = ?;", value).Scan(&cPushbulletKey, &cEmail)
	if err != nil {
		if err == sql.ErrNoRows {
			resp = AdminWebsiteNotificationsResponse{true, WebsiteNotifications{"", ""}}
		} else {
			logging.MustGetLogger("logger").Error("Unable to get Website's notification settings: ", err)
			SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to process your Request: "+err.Error())
			return
		}
	} else {
		resp = AdminWebsiteNotificationsResponse{true, WebsiteNotifications{cPushbulletKey, cEmail}}
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
func ApiAdminWebsiteUpdateNotifications(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
				logging.MustGetLogger("logger").Error("Unable to get Website's id for notification-insertion: ", err)
				SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to process your Request: "+err.Error())
				return
			}

			_, err = db.Exec("INSERT INTO notifications (websiteId, pushbulletKey, email) VALUES (?, ?, ?);", id, pushbulletKey, email)
			if err != nil {
				logging.MustGetLogger("logger").Error("Unable to insert Website's notification settings: ", err)
				SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to process your Request: "+err.Error())
				return
			}
		} else {
			logging.MustGetLogger("logger").Error("Unable to get Website's notification settings: ", err)
			SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to process your Request: "+err.Error())
			return
		}
	} else {
		// existing settings found --> Update
		_, err = db.Exec("UPDATE notifications, websites SET pushbulletKey = ?, email = ? WHERE notifications.websiteId = websites.id AND url = ?;", pushbulletKey, email, url)
		if err != nil {
			logging.MustGetLogger("logger").Error("Unable to update Website's notification settings: ", err)
			SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to process your Request: "+err.Error())
			return
		}
	}

	SendJsonMessage(w, http.StatusOK, true, "")
}

// Edits an existing Website in the database.
func ApiAdminWebsiteEdit(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
	if method != "HEAD" && method != "GET" {
		SendJsonMessage(w, http.StatusBadRequest, false, "Unable to process your Request: Submit a valid check method.")
		return
	}

	// Update Database
	db := lib.GetDatabase()
	res, err := db.Exec("UPDATE websites SET name = ?, protocol = ?, url = ?, checkMethod = ? WHERE url = ?;", name, protocol, url, method, oldUrl)
	if err != nil {
		logging.MustGetLogger("logger").Error("Unable to edit Website: ", err)
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
func ApiAdminWebsiteDelete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
		logging.MustGetLogger("logger").Error("Unable to delete Check-Results: ", err)
		SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to process your Request: "+err.Error())
		return
	}

	// Remove Notifications from Database
	res, err = db.Exec("DELETE n FROM notifications n INNER JOIN websites w ON n.websiteId = w.id WHERE w.url = ?;", value)
	if err != nil {
		logging.MustGetLogger("logger").Error("Unable to delete Notifications: ", err)
		SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to process your Request: "+err.Error())
		return
	}

	// Remove Website from Database
	db = lib.GetDatabase()
	res, err = db.Exec("DELETE FROM websites WHERE url = ?;", value)
	if err != nil {
		logging.MustGetLogger("logger").Error("Unable to delete Website: ", err)
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

// Updates the application's title in the database.
func ApiAdminSettingTitle(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
		logging.MustGetLogger("logger").Error("Unable to change Application-Title: ", err)
		SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to process your Request: "+err.Error())
		return
	}

	// Update Configuration
	lib.GetConfiguration().Dynamic.Title = value
	SendJsonMessage(w, http.StatusOK, true, "")
}

// Updates the application's check-interval in the database.
func ApiAdminSettingInterval(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
		logging.MustGetLogger("logger").Error("Unable to change Interval: ", err)
		SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to process your Request: "+err.Error())
		return
	}

	// Update Configuration
	lib.GetConfiguration().Dynamic.Interval = value
	SendJsonMessage(w, http.StatusOK, true, "")
}

// Updates the application's maximum amount of redirects in the database.
func ApiAdminSettingRedirects(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
		logging.MustGetLogger("logger").Error("Unable to change Redirects: ", err)
		SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to process your Request: "+err.Error())
		return
	}

	// Update Configuration
	lib.GetConfiguration().Dynamic.Redirects = value
	SendJsonMessage(w, http.StatusOK, true, "")
}

// Updates the user's password in the database.
func ApiAdminSettingPassword(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

// Triggers a check of all enabled Websites.
func ApiAdminActionCheck(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !lib.IsLoggedIn(r) {
		SendJsonMessage(w, http.StatusUnauthorized, false, "Unauthorized.")
		return
	}

	// Trigger a check
	lib.GetConfiguration().Dynamic.CheckNow = true
	SendJsonMessage(w, http.StatusOK, true, "")
}

// Processes a login-request and sends an authentication-cookie to the browser.
func ApiAdminActionLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if lib.IsLoggedIn(r) {
		SendJsonMessage(w, http.StatusBadRequest, false, "Already logged in.")
		return
	}

	// Get data from Request
	r.ParseForm()
	value := r.Form.Get("password")

	// Check Password
	admin := lib.Admin{}
	if admin.ValidatePassword(value) {
		cookie := lib.LoginAndGetCookie("admin")
		http.SetCookie(w, &cookie)
		SendJsonMessage(w, http.StatusOK, true, "")
	} else {
		SendJsonMessage(w, http.StatusBadRequest, false, "Unable to process your Request: Invalid Password.")
	}
}

// Processes a logout-request and sends a termination-cookie to the browser.
func ApiAdminActionLogout(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !lib.IsLoggedIn(r) {
		SendJsonMessage(w, http.StatusUnauthorized, false, "Unauthorized.")
		return
	}

	// Process logout
	cookie := lib.LogoutAndDestroyCookie(r)
	http.SetCookie(w, &cookie)
	SendJsonMessage(w, http.StatusOK, true, "")
}
