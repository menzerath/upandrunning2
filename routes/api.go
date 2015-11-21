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

// Sends a simple welcome-message to the user.
func ApiIndex(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	SendJsonMessage(w, http.StatusOK, true, "Welcome to UpAndRunning2's API!")
}

// Returns a DetailedWebsiteResponse containing all the Website's important data if the Website is enabled.
func ApiStatus(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var (
		id             int
		name           string
		protocol       string
		url            string
		status         string
		time           string
		lastFailStatus string
		lastFailTime   string
		ups            int
		downs          int
		totalChecks    int
		avgAvail       float64
	)

	// Query the Database
	db := lib.GetDatabase()
	err := db.QueryRow("SELECT id, name, protocol, url, status, time, lastFailStatus, lastFailTime, ups, downs, totalChecks, avgAvail FROM websites WHERE url = ? AND enabled = 1;", ps.ByName("url")).Scan(&id, &name, &protocol, &url, &status, &time, &lastFailStatus, &lastFailTime, &ups, &downs, &totalChecks, &avgAvail)
	switch {
	case err == sql.ErrNoRows:
		SendJsonMessage(w, http.StatusNotFound, false, "Unable to find any data matching the given url.")
		return
	case err != nil:
		logging.MustGetLogger("logger").Error("Unable to fetch Website-Status: ", err)
		SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to process your Request.")
		return
	}

	// Build Response
	responseJson := DetailedWebsiteResponse{true, WebsiteData{id, name, protocol + "://" + url}, WebsiteAvailability{ups, downs, totalChecks, strconv.FormatFloat(avgAvail, 'f', 2, 64) + "%"}, WebsiteCheckResult{status, time}, WebsiteCheckResult{lastFailStatus, lastFailTime}}

	// Send Response
	responseBytes, err := json.Marshal(responseJson)
	if err != nil {
		SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to process your Request.")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBytes)
}

// Returns a WebsiteResponse containing all publicly visible Websites as BasicWebsite.
func ApiWebsites(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Query the Database
	db := lib.GetDatabase()
	rows, err := db.Query("SELECT name, protocol, url, status FROM websites WHERE enabled = 1 AND visible = 1;")
	if err != nil {
		logging.MustGetLogger("logger").Error("Unable to fetch Websites: ", err)
		SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to process your Request.")
		return
	}
	defer rows.Close()

	// Add every Website
	websites := []BasicWebsite{}
	for rows.Next() {
		var row BasicWebsite
		err = rows.Scan(&row.Name, &row.Protocol, &row.Url, &row.Status)
		if err != nil {
			logging.MustGetLogger("logger").Error("Unable to read Website-Row: ", err)
			SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to process your Request.")
			return
		}
		websites = append(websites, row)
	}

	// Check for Errors
	err = rows.Err()
	if err != nil {
		logging.MustGetLogger("logger").Error("Unable to read Website-Rows: ", err)
		SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to process your Request.")
		return
	}

	// Send Response
	responseBytes, err := json.Marshal(WebsiteResponse{true, websites})
	if err != nil {
		SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to process your Request.")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBytes)
}
