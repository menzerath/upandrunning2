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

}

func ApiAdminSettingInterval(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

func ApiAdminSettingPushbulletKey(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

func ApiAdminSettingPassword(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

func ApiAdminActionCheck(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	lib.GetConfiguration().Dynamic.CheckNow = true
	SendJsonMessage(w, http.StatusOK, true, "")
}

func ApiAdminActionLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

func ApiAdminActionLogout(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}
