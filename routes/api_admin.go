package routes

import (
	"encoding/json"
	"github.com/MarvinMenzerath/UpAndRunning2/lib"
	"github.com/julienschmidt/httprouter"
	"github.com/op/go-logging"
	"net/http"
	"strconv"
)

// Sends a simple welcome-message to the user.
func ApiAdminIndex(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	SendJsonMessage(w, http.StatusOK, true, "Welcome to UpAndRunning2's Admin-API! Please be aware that most operations need an incoming POST-request.")
}

// Returns a AdminWebsiteResponse containing all Websites as AdminWebsite.
func ApiAdminWebsites(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !lib.IsLoggedIn(r) {
		SendJsonMessage(w, http.StatusUnauthorized, false, "Unauthorized.")
		return
	}

	// Query the Database for basic data and the last check
	db := lib.GetDatabase()
	rows, err := db.Query("SELECT websites.id, websites.name, websites.enabled, websites.visible, websites.protocol, websites.url, websites.checkMethod, checks.statusCode, checks.statusText, checks.time FROM checks, websites WHERE checks.websiteId = websites.id AND NOT EXISTS (SELECT id FROM checks c2 WHERE checks.websiteId = c2.websiteId AND checks.id < c2.id) ORDER BY websites.id;")
	if err != nil {
		logging.MustGetLogger("logger").Error("Unable to fetch Websites: ", err)
		SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to process your Request.")
		return
	}
	defer rows.Close()

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

	// Add every Website
	websites := []AdminWebsite{}
	for rows.Next() {
		err = rows.Scan(&id, &name, &enabled, &visible, &protocol, &url, &checkMethod, &statusCode, &statusText, &time)
		if err != nil {
			logging.MustGetLogger("logger").Error("Unable to read Website-Row: ", err)
			SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to process your Request.")
			return
		}
		websites = append(websites, AdminWebsite{id, name, enabled, visible, protocol, url, checkMethod, statusCode + " - " + statusText, time})
	}

	// Check for Errors
	err = rows.Err()
	if err != nil {
		logging.MustGetLogger("logger").Error("Unable to read Website-Rows: ", err)
		SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to process your Request.")
		return
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
	url := r.Form.Get("url")
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
func ApiAdminWebsiteEnable(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !lib.IsLoggedIn(r) {
		SendJsonMessage(w, http.StatusUnauthorized, false, "Unauthorized.")
		return
	}

	// Get data from Request
	r.ParseForm()
	value := r.Form.Get("url")

	// Simple Validation
	if value == "" {
		SendJsonMessage(w, http.StatusBadRequest, false, "Unable to process your Request: Submit a valid value.")
		return
	}

	// Update Database-Row
	db := lib.GetDatabase()
	res, err := db.Exec("UPDATE websites SET enabled = 1 WHERE url = ?;", value)
	if err != nil {
		logging.MustGetLogger("logger").Error("Unable to enable Website: ", err)
		SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to process your Request: "+err.Error())
		return
	}

	// Check if exactly one Website has been enabled
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 1 {
		SendJsonMessage(w, http.StatusOK, true, "")
	} else {
		SendJsonMessage(w, http.StatusBadRequest, false, "Unable to process your Request: Could not enable Website.")
	}
}

// Disables an existing Website in the database.
func ApiAdminWebsiteDisable(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !lib.IsLoggedIn(r) {
		SendJsonMessage(w, http.StatusUnauthorized, false, "Unauthorized.")
		return
	}

	// Get data from Request
	r.ParseForm()
	value := r.Form.Get("url")

	// Simple Validation
	if value == "" {
		SendJsonMessage(w, http.StatusBadRequest, false, "Unable to process your Request: Submit a valid value.")
		return
	}

	// Update Database-Row
	db := lib.GetDatabase()
	res, err := db.Exec("UPDATE websites SET enabled = 0 WHERE url = ?;", value)
	if err != nil {
		logging.MustGetLogger("logger").Error("Unable to disable Website: ", err)
		SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to process your Request: "+err.Error())
		return
	}

	// Check if exactly one Website has been disabled
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 1 {
		SendJsonMessage(w, http.StatusOK, true, "")
	} else {
		SendJsonMessage(w, http.StatusBadRequest, false, "Unable to process your Request: Could not disable Website.")
	}
}

// Sets an existing Website to visible in the database.
func ApiAdminWebsiteVisible(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !lib.IsLoggedIn(r) {
		SendJsonMessage(w, http.StatusUnauthorized, false, "Unauthorized.")
		return
	}

	// Get data from Request
	r.ParseForm()
	value := r.Form.Get("url")

	// Simple Validation
	if value == "" {
		SendJsonMessage(w, http.StatusBadRequest, false, "Unable to process your Request: Submit a valid value.")
		return
	}

	// Update Database-Row
	db := lib.GetDatabase()
	res, err := db.Exec("UPDATE websites SET visible = 1 WHERE url = ?;", value)
	if err != nil {
		logging.MustGetLogger("logger").Error("Unable to set Website visible: ", err)
		SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to process your Request: "+err.Error())
		return
	}

	// Check if exactly one Website has been set to visible
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 1 {
		SendJsonMessage(w, http.StatusOK, true, "")
	} else {
		SendJsonMessage(w, http.StatusBadRequest, false, "Unable to process your Request: Could not set Website to visible.")
	}
}

// Sets an existing Website to invisible in the database.
func ApiAdminWebsiteInvisible(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !lib.IsLoggedIn(r) {
		SendJsonMessage(w, http.StatusUnauthorized, false, "Unauthorized.")
		return
	}

	// Get data from Request
	r.ParseForm()
	value := r.Form.Get("url")

	// Simple Validation
	if value == "" {
		SendJsonMessage(w, http.StatusBadRequest, false, "Unable to process your Request: Submit a valid value.")
		return
	}

	// Update Database-Row
	db := lib.GetDatabase()
	res, err := db.Exec("UPDATE websites SET visible = 0 WHERE url = ?;", value)
	if err != nil {
		logging.MustGetLogger("logger").Error("Unable to set Website invisible: ", err)
		SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to process your Request: "+err.Error())
		return
	}

	// Check if exactly one Website has been set to invisible
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 1 {
		SendJsonMessage(w, http.StatusOK, true, "")
	} else {
		SendJsonMessage(w, http.StatusBadRequest, false, "Unable to process your Request: Could not set Website to invisible.")
	}
}

// Edits an existing Website in the database.
func ApiAdminWebsiteEdit(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !lib.IsLoggedIn(r) {
		SendJsonMessage(w, http.StatusUnauthorized, false, "Unauthorized.")
		return
	}

	// Get data from Request
	r.ParseForm()
	oldUrl := r.Form.Get("oldUrl")
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
	r.ParseForm()
	value := r.Form.Get("url")

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

// Updates the application's Pushbullet-key in the database.
func ApiAdminSettingPushbulletKey(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !lib.IsLoggedIn(r) {
		SendJsonMessage(w, http.StatusUnauthorized, false, "Unauthorized.")
		return
	}

	// Get data from Request
	r.ParseForm()
	value := r.Form.Get("key")

	// Simple Validation
	if value == "" {
		SendJsonMessage(w, http.StatusBadRequest, false, "Unable to process your Request: Submit a valid value.")
		return
	}

	// Update Database-Row
	db := lib.GetDatabase()
	_, err := db.Exec("UPDATE settings SET value = ? WHERE name = 'pushbullet_key';", value)
	if err != nil {
		logging.MustGetLogger("logger").Error("Unable to change PushBullet-API-Key: ", err)
		SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to process your Request: "+err.Error())
		return
	}

	// Update Configuration
	lib.GetConfiguration().Dynamic.PushbulletKey = value
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
