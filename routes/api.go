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
		id                   int
		name                 string
		protocol             string
		url                  string
		statusCode           string
		statusText           string
		responseTime         string
		time                 string
		lastFailStatusCode   string
		lastFailStatusText   string
		lastFailResponseTime string
		lastFailTime         string
		ups                  int
		totalChecks          int
	)

	// Query the Database for basic data and the last successful check
	db := lib.GetDatabase()
	err := db.QueryRow("SELECT websites.id, websites.name, websites.protocol, websites.url, checks.statusCode, checks.statusText, checks.responseTime, checks.time FROM checks, websites WHERE checks.websiteId = websites.id AND websites.url = ? AND websites.enabled = 1 ORDER BY checks.id DESC LIMIT 1;", ps.ByName("url")).Scan(&id, &name, &protocol, &url, &statusCode, &statusText, &responseTime, &time)
	switch {
	case err == sql.ErrNoRows:
		SendJsonMessage(w, http.StatusNotFound, false, "Unable to find any data matching the given url.")
		return
	case err != nil:
		logging.MustGetLogger("logger").Error("Unable to fetch Website-Status: ", err)
		SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to process your Request.")
		return
	}

	// Query the Database for the last unsuccessful check
	err = db.QueryRow("SELECT checks.statusCode, checks.statusText, checks.responseTime, checks.time FROM checks, websites WHERE checks.websiteId = websites.id AND (statusCode NOT LIKE '2%' AND statusCode NOT LIKE '3%') AND websites.url = ? AND websites.enabled = 1 ORDER BY checks.id DESC LIMIT 1;", ps.ByName("url")).Scan(&lastFailStatusCode, &lastFailStatusText, &lastFailResponseTime, &lastFailTime)
	switch {
	case err == sql.ErrNoRows:
		lastFailStatusCode = "0"
		lastFailStatusText = "unknown"
		lastFailResponseTime = "0"
		lastFailTime = "0000-00-00 00:00:00"
	case err != nil:
		logging.MustGetLogger("logger").Error("Unable to fetch Website-Status: ", err)
		SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to process your Request.")
		return
	}

	// Query the Database for the amount of (successful / total) checks
	err = db.QueryRow("SELECT (SELECT COUNT(checks.id) FROM checks, websites WHERE checks.websiteId = websites.id AND (statusCode LIKE '2%' OR statusCode LIKE '3%') AND websites.url = ?) AS ups, (SELECT COUNT(checks.id) FROM checks, websites WHERE checks.websiteId = websites.id AND websites.url = ?) AS total FROM checks LIMIT 1;", ps.ByName("url"), ps.ByName("url")).Scan(&ups, &totalChecks)
	switch {
	case err == sql.ErrNoRows:
		logging.MustGetLogger("logger").Error("Unable to fetch Website-Status: ", err)
		SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to process your Request.")
		return
	case err != nil:
		logging.MustGetLogger("logger").Error("Unable to fetch Website-Status: ", err)
		SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to process your Request.")
		return
	}

	// Build Response
	responseJson := DetailedWebsiteResponse{true, WebsiteData{id, name, protocol + "://" + url}, WebsiteAvailability{ups, totalChecks - ups, totalChecks, strconv.FormatFloat((float64(ups)/float64(totalChecks))*100, 'f', 2, 64) + "%"}, WebsiteCheckResult{statusCode + " - " + statusText, responseTime + " ms", time}, WebsiteCheckResult{lastFailStatusCode + " - " + lastFailStatusText, lastFailResponseTime + " ms", lastFailTime}}

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
	rows, err := db.Query("SELECT websites.name, websites.protocol, websites.url, checks.statusCode, checks.statusText FROM checks, websites WHERE checks.websiteId = websites.id AND NOT EXISTS (SELECT id FROM checks c2 WHERE checks.websiteId = c2.websiteId AND checks.id < c2.id) AND websites.enabled = 1 AND websites.visible = 1 ORDER BY websites.id;")
	if err != nil {
		logging.MustGetLogger("logger").Error("Unable to fetch Websites: ", err)
		SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to process your Request.")
		return
	}
	defer rows.Close()

	// Add every Website
	websites := []BasicWebsite{}
	var (
		name       string
		protocol   string
		url        string
		statusCode string
		statusText string
	)
	for rows.Next() {
		err = rows.Scan(&name, &protocol, &url, &statusCode, &statusText)
		if err != nil {
			logging.MustGetLogger("logger").Error("Unable to read Website-Row: ", err)
			SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to process your Request.")
			return
		}

		websites = append(websites, BasicWebsite{name, protocol, url, statusCode + " - " + statusText})
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
