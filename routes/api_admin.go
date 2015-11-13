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

}

func ApiAdminWebsiteEnable(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

func ApiAdminWebsiteDisable(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

func ApiAdminWebsiteVisible(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

func ApiAdminWebsiteInvisible(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

func ApiAdminWebsiteEdit(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

func ApiAdminWebsiteDelete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

func ApiAdminSettingTitle(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
		logging.MustGetLogger("logger").Fatal("Unable to change Application-Title: ", err)
		SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to process your Request: "+err.Error())
	}
	_, err = stmt.Exec(value)
	if err != nil {
		logging.MustGetLogger("logger").Fatal("Unable to change Application-Title: ", err)
	}

	lib.GetConfiguration().Dynamic.Title = value
	SendJsonMessage(w, http.StatusOK, true, "")
}

func ApiAdminSettingInterval(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
	stmt, err := db.Prepare("UPDATE settings SET value = ? WHERE name = 'title';")
	if err != nil {
		logging.MustGetLogger("logger").Fatal("Unable to change Interval: ", err)
		SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to process your Request: "+err.Error())
	}
	_, err = stmt.Exec(value)
	if err != nil {
		logging.MustGetLogger("logger").Fatal("Unable to change Interval: ", err)
	}

	lib.GetConfiguration().Dynamic.Interval = value
	SendJsonMessage(w, http.StatusOK, true, "")
}

func ApiAdminSettingPushbulletKey(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
		logging.MustGetLogger("logger").Fatal("Unable to change PushBullet-API-Key: ", err)
		SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to process your Request: "+err.Error())
	}
	_, err = stmt.Exec(value)
	if err != nil {
		logging.MustGetLogger("logger").Fatal("Unable to change PushBullet-API-Key: ", err)
	}

	lib.GetConfiguration().Dynamic.PushbulletKey = value
	SendJsonMessage(w, http.StatusOK, true, "")
}

func ApiAdminSettingPassword(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
	lib.GetConfiguration().Dynamic.CheckNow = true
	SendJsonMessage(w, http.StatusOK, true, "")
}

func ApiAdminActionLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

func ApiAdminActionLogout(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}
