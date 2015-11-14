package routes

import (
	"encoding/json"
	"github.com/MarvinMenzerath/UpAndRunning2/lib"
	"github.com/julienschmidt/httprouter"
	"github.com/op/go-logging"
	"net/http"
	"strconv"
)

func ApiAdminIndex(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	SendJsonMessage(w, http.StatusOK, true, "Welcome to UpAndRunning's Admin-API! Please be aware that most operations need an incoming POST-request.")
}

func ApiAdminWebsites(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !lib.IsLoggedIn(r) {
		SendJsonMessage(w, http.StatusUnauthorized, false, "Unauthorized.")
		return
	}

	// Query the Database
	db := lib.GetDatabase()
	rows, err := db.Query("SELECT id, name, enabled, visible, protocol, url, status, time, avgAvail FROM website;")
	if err != nil {
		logging.MustGetLogger("logger").Error("Unable to fetch Websites: ", err)
		SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to process your Request.")
		return
	}
	defer rows.Close()

	var (
		id       int
		name     string
		enabled  bool
		visible  bool
		protocol string
		url      string
		status   string
		time     string
		average  float64
	)

	// Add every Website
	websites := []AdminWebsite{}
	for rows.Next() {
		err = rows.Scan(&id, &name, &enabled, &visible, &protocol, &url, &status, &time, &average)
		if err != nil {
			logging.MustGetLogger("logger").Error("Unable to read Website-Row: ", err)
			SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to process your Request.")
			return
		}
		websites = append(websites, AdminWebsite{id, name, enabled, visible, protocol, url, status, time, strconv.FormatFloat(average, 'f', 2, 64) + "%"})
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

func ApiAdminWebsiteAdd(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !lib.IsLoggedIn(r) {
		SendJsonMessage(w, http.StatusUnauthorized, false, "Unauthorized.")
		return
	}

	r.ParseForm()
	name := r.Form.Get("name")
	protocol := r.Form.Get("protocol")
	url := r.Form.Get("url")

	// Simple Validation
	if name == "" || protocol == "" || url == "" {
		SendJsonMessage(w, http.StatusBadRequest, false, "Unable to process your Request: Submit valid values.")
		return
	}

	// Insert into Database
	db := lib.GetDatabase()
	stmt, err := db.Prepare("INSERT INTO website (name, protocol, url) VALUES (?, ?, ?);")
	if err != nil {
		logging.MustGetLogger("logger").Error("Unable to add Website: ", err)
		SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to process your Request: "+err.Error())
		return
	}
	_, err = stmt.Exec(name, protocol, url)
	if err != nil {
		logging.MustGetLogger("logger").Error("Unable to add Website: ", err)
		SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to process your Request: "+err.Error())
		return
	}

	SendJsonMessage(w, http.StatusOK, true, "")
}

func ApiAdminWebsiteEnable(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !lib.IsLoggedIn(r) {
		SendJsonMessage(w, http.StatusUnauthorized, false, "Unauthorized.")
		return
	}

	r.ParseForm()
	value := r.Form.Get("url")

	// Simple Validation
	if value == "" {
		SendJsonMessage(w, http.StatusBadRequest, false, "Unable to process your Request: Submit a valid value.")
		return
	}

	// Update Database-Row
	db := lib.GetDatabase()
	stmt, err := db.Prepare("UPDATE website SET enabled = 1 WHERE url = ?;")
	if err != nil {
		logging.MustGetLogger("logger").Error("Unable to enable Website: ", err)
		SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to process your Request: "+err.Error())
		return
	}
	res, err := stmt.Exec(value)
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

func ApiAdminWebsiteDisable(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !lib.IsLoggedIn(r) {
		SendJsonMessage(w, http.StatusUnauthorized, false, "Unauthorized.")
		return
	}

	r.ParseForm()
	value := r.Form.Get("url")

	// Simple Validation
	if value == "" {
		SendJsonMessage(w, http.StatusBadRequest, false, "Unable to process your Request: Submit a valid value.")
		return
	}

	// Update Database-Row
	db := lib.GetDatabase()
	stmt, err := db.Prepare("UPDATE website SET enabled = 0 WHERE url = ?;")
	if err != nil {
		logging.MustGetLogger("logger").Error("Unable to disable Website: ", err)
		SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to process your Request: "+err.Error())
		return
	}
	res, err := stmt.Exec(value)
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

func ApiAdminWebsiteVisible(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !lib.IsLoggedIn(r) {
		SendJsonMessage(w, http.StatusUnauthorized, false, "Unauthorized.")
		return
	}

	r.ParseForm()
	value := r.Form.Get("url")

	// Simple Validation
	if value == "" {
		SendJsonMessage(w, http.StatusBadRequest, false, "Unable to process your Request: Submit a valid value.")
		return
	}

	// Update Database-Row
	db := lib.GetDatabase()
	stmt, err := db.Prepare("UPDATE website SET visible = 1 WHERE url = ?;")
	if err != nil {
		logging.MustGetLogger("logger").Error("Unable to set Website visible: ", err)
		SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to process your Request: "+err.Error())
		return
	}
	res, err := stmt.Exec(value)
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

func ApiAdminWebsiteInvisible(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !lib.IsLoggedIn(r) {
		SendJsonMessage(w, http.StatusUnauthorized, false, "Unauthorized.")
		return
	}

	r.ParseForm()
	value := r.Form.Get("url")

	// Simple Validation
	if value == "" {
		SendJsonMessage(w, http.StatusBadRequest, false, "Unable to process your Request: Submit a valid value.")
		return
	}

	// Update Database-Row
	db := lib.GetDatabase()
	stmt, err := db.Prepare("UPDATE website SET visible = 0 WHERE url = ?;")
	if err != nil {
		logging.MustGetLogger("logger").Error("Unable to set Website invisible: ", err)
		SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to process your Request: "+err.Error())
		return
	}
	res, err := stmt.Exec(value)
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

func ApiAdminWebsiteEdit(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !lib.IsLoggedIn(r) {
		SendJsonMessage(w, http.StatusUnauthorized, false, "Unauthorized.")
		return
	}

	r.ParseForm()
	oldUrl := r.Form.Get("oldUrl")
	name := r.Form.Get("name")
	protocol := r.Form.Get("protocol")
	url := r.Form.Get("url")

	// Simple Validation
	if oldUrl == "" || name == "" || protocol == "" || url == "" {
		SendJsonMessage(w, http.StatusBadRequest, false, "Unable to process your Request: Submit valid values.")
		return
	}

	// Update Database
	db := lib.GetDatabase()
	stmt, err := db.Prepare("UPDATE website SET name = ?, protocol = ?, url = ? WHERE url = ?;")
	if err != nil {
		logging.MustGetLogger("logger").Error("Unable to edit Website: ", err)
		SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to process your Request: "+err.Error())
		return
	}
	res, err := stmt.Exec(name, protocol, url, oldUrl)
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

func ApiAdminWebsiteDelete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !lib.IsLoggedIn(r) {
		SendJsonMessage(w, http.StatusUnauthorized, false, "Unauthorized.")
		return
	}

	r.ParseForm()
	value := r.Form.Get("url")

	// Simple Validation
	if value == "" {
		SendJsonMessage(w, http.StatusBadRequest, false, "Unable to process your Request: Submit a valid value.")
		return
	}

	// Remove from Database
	db := lib.GetDatabase()
	stmt, err := db.Prepare("DELETE FROM website WHERE url = ?;")
	if err != nil {
		logging.MustGetLogger("logger").Error("Unable to delete Website: ", err)
		SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to process your Request: "+err.Error())
		return
	}
	res, err := stmt.Exec(value)
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

func ApiAdminSettingTitle(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !lib.IsLoggedIn(r) {
		SendJsonMessage(w, http.StatusUnauthorized, false, "Unauthorized.")
		return
	}

	r.ParseForm()
	value := r.Form.Get("title")

	// Simple Validation
	if value == "" {
		SendJsonMessage(w, http.StatusBadRequest, false, "Unable to process your Request: Submit a valid value.")
		return
	}

	// Update Database-Row
	db := lib.GetDatabase()
	stmt, err := db.Prepare("UPDATE settings SET value = ? WHERE name = 'title';")
	if err != nil {
		logging.MustGetLogger("logger").Error("Unable to change Application-Title: ", err)
		SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to process your Request: "+err.Error())
		return
	}
	_, err = stmt.Exec(value)
	if err != nil {
		logging.MustGetLogger("logger").Error("Unable to change Application-Title: ", err)
		SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to process your Request: "+err.Error())
		return
	}

	lib.GetConfiguration().Dynamic.Title = value
	SendJsonMessage(w, http.StatusOK, true, "")
}

func ApiAdminSettingInterval(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !lib.IsLoggedIn(r) {
		SendJsonMessage(w, http.StatusUnauthorized, false, "Unauthorized.")
		return
	}

	r.ParseForm()
	temp := r.Form.Get("interval")
	value, err := strconv.Atoi(temp)

	// Simple Validation
	if err != nil || value < 1 || value > 600 {
		SendJsonMessage(w, http.StatusBadRequest, false, "Unable to process your Request: Submit a valid value between 1 and 600 seconds.")
		return
	}

	// Update Database-Row
	db := lib.GetDatabase()
	stmt, err := db.Prepare("UPDATE settings SET value = ? WHERE name = 'interval';")
	if err != nil {
		logging.MustGetLogger("logger").Error("Unable to change Interval: ", err)
		SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to process your Request: "+err.Error())
		return
	}
	_, err = stmt.Exec(value)
	if err != nil {
		logging.MustGetLogger("logger").Error("Unable to change Interval: ", err)
		SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to process your Request: "+err.Error())
		return
	}

	lib.GetConfiguration().Dynamic.Interval = value
	SendJsonMessage(w, http.StatusOK, true, "")
}

func ApiAdminSettingPushbulletKey(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !lib.IsLoggedIn(r) {
		SendJsonMessage(w, http.StatusUnauthorized, false, "Unauthorized.")
		return
	}

	r.ParseForm()
	value := r.Form.Get("key")

	// Simple Validation
	if value == "" {
		SendJsonMessage(w, http.StatusBadRequest, false, "Unable to process your Request: Submit a valid value.")
		return
	}

	// Update Database-Row
	db := lib.GetDatabase()
	stmt, err := db.Prepare("UPDATE settings SET value = ? WHERE name = 'pushbullet_key';")
	if err != nil {
		logging.MustGetLogger("logger").Error("Unable to change PushBullet-API-Key: ", err)
		SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to process your Request: "+err.Error())
		return
	}
	_, err = stmt.Exec(value)
	if err != nil {
		logging.MustGetLogger("logger").Error("Unable to change PushBullet-API-Key: ", err)
		SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to process your Request: "+err.Error())
		return
	}

	lib.GetConfiguration().Dynamic.PushbulletKey = value
	SendJsonMessage(w, http.StatusOK, true, "")
}

func ApiAdminSettingPassword(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !lib.IsLoggedIn(r) {
		SendJsonMessage(w, http.StatusUnauthorized, false, "Unauthorized.")
		return
	}

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

func ApiAdminActionCheck(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !lib.IsLoggedIn(r) {
		SendJsonMessage(w, http.StatusUnauthorized, false, "Unauthorized.")
		return
	}

	lib.GetConfiguration().Dynamic.CheckNow = true
	SendJsonMessage(w, http.StatusOK, true, "")
}

func ApiAdminActionLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if lib.IsLoggedIn(r) {
		SendJsonMessage(w, http.StatusBadRequest, false, "Already logged in.")
		return
	}

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

func ApiAdminActionLogout(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !lib.IsLoggedIn(r) {
		SendJsonMessage(w, http.StatusUnauthorized, false, "Unauthorized.")
		return
	}

	cookie := lib.LogoutAndDestroyCookie(r)
	http.SetCookie(w, &cookie)
	SendJsonMessage(w, http.StatusOK, true, "")
}
