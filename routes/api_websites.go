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

// Returns a WebsiteResponse containing all publicly visible Websites as BasicWebsite.
func ApiWebsites(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if lib.IsLoggedIn(r) {
		// Send a more detailed version if the user is logged in
		ApiWebsitesDetailed(w, r, ps)
		return
	}

	// Query the Database for basic data
	db := lib.GetDatabase()
	rows, err := db.Query("SELECT id, name, protocol, url FROM websites WHERE enabled = 1 AND visible = 1 ORDER BY name;")
	if err != nil {
		logging.MustGetLogger("").Error("Unable to fetch Websites: ", err)
		SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to process your Request.")
		return
	}
	defer rows.Close()

	// Add every Website
	websites := []BasicWebsite{}
	var (
		id         int
		name       string
		protocol   string
		url        string
		statusCode string
		statusText string
	)

	for rows.Next() {
		err = rows.Scan(&id, &name, &protocol, &url)
		if err != nil {
			logging.MustGetLogger("").Error("Unable to read Website-Data-Row: ", err)
			SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to process your Request.")
			return
		}

		// Query the database for status data
		err = db.QueryRow("SELECT statusCode, statusText FROM checks WHERE websiteId = ? ORDER BY id DESC LIMIT 1;", id).Scan(&statusCode, &statusText)
		switch {
		case err == sql.ErrNoRows:
			statusCode = "0"
			statusText = "unknown"
		case err != nil:
			logging.MustGetLogger("").Error("Unable to fetch Website-Status: ", err)
			SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to process your Request.")
			return
		}

		websites = append(websites, BasicWebsite{name, protocol, url, statusCode + " - " + statusText})
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

// Returns a AdminWebsiteResponse containing all Websites as AdminWebsite.
func ApiWebsitesDetailed(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !lib.IsLoggedIn(r) {
		SendJsonMessage(w, http.StatusUnauthorized, false, "Unauthorized.")
		return
	}

	// Query the Database for basic data
	db := lib.GetDatabase()
	rows, err := db.Query("SELECT id, name, enabled, visible, protocol, url, checkMethod FROM websites ORDER BY name;")
	if err != nil {
		logging.MustGetLogger("").Error("Unable to fetch Websites: ", err)
		SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to process your Request.")
		return
	}
	defer rows.Close()

	// Add every Website
	websites := []DetailedWebsite{}
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
			logging.MustGetLogger("").Error("Unable to read Website-Data-Row: ", err)
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
			logging.MustGetLogger("").Error("Unable to fetch Website's status: ", err)
			SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to process your Request.")
			return
		}

		websites = append(websites, DetailedWebsite{id, name, enabled, visible, protocol, url, checkMethod, statusCode + " - " + statusText, time})
	}

	// Send Response
	responseBytes, err := json.Marshal(DetailedWebsiteResponse{true, websites})
	if err != nil {
		SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to process your Request.")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBytes)
}

// Returns a StatusResponse containing all the Website's important data if the Website is enabled.
func ApiWebsitesStatus(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
		logging.MustGetLogger("").Error("Unable to fetch Website-Status: ", err)
		SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to process your Request.")
		return
	}

	// Query the Database for the last unsuccessful check
	err = db.QueryRow("SELECT checks.statusCode, checks.statusText, checks.responseTime, checks.time FROM checks, websites WHERE checks.websiteId = websites.id AND (checks.statusCode NOT LIKE '2%' AND checks.statusCode NOT LIKE '3%') AND websites.url = ? AND websites.enabled = 1 ORDER BY checks.id DESC LIMIT 1;", ps.ByName("url")).Scan(&lastFailStatusCode, &lastFailStatusText, &lastFailResponseTime, &lastFailTime)
	switch {
	case err == sql.ErrNoRows:
		lastFailStatusCode = "0"
		lastFailStatusText = "unknown"
		lastFailResponseTime = "0"
		lastFailTime = "0000-00-00 00:00:00"
	case err != nil:
		logging.MustGetLogger("").Error("Unable to fetch Website-Status: ", err)
		SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to process your Request.")
		return
	}

	// Query the Database for the amount of (successful / total) checks
	err = db.QueryRow("SELECT (SELECT COUNT(checks.id) FROM checks, websites WHERE checks.websiteId = websites.id AND (checks.statusCode LIKE '2%' OR checks.statusCode LIKE '3%') AND websites.url = ?) AS ups, (SELECT COUNT(checks.id) FROM checks, websites WHERE checks.websiteId = websites.id AND websites.url = ?) AS total FROM checks LIMIT 1;", ps.ByName("url"), ps.ByName("url")).Scan(&ups, &totalChecks)
	switch {
	case err == sql.ErrNoRows:
		logging.MustGetLogger("").Error("Unable to fetch Website-Status: ", err)
		SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to process your Request.")
		return
	case err != nil:
		logging.MustGetLogger("").Error("Unable to fetch Website-Status: ", err)
		SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to process your Request.")
		return
	}

	// Build Response
	responseJson := StatusResponse{true, WebsiteData{id, name, protocol + "://" + url}, WebsiteAvailability{ups, totalChecks - ups, totalChecks, strconv.FormatFloat((float64(ups)/float64(totalChecks))*100, 'f', 2, 64) + "%"}, WebsiteCheckResult{statusCode + " - " + statusText, responseTime + " ms", time}, WebsiteCheckResult{lastFailStatusCode + " - " + lastFailStatusText, lastFailResponseTime + " ms", lastFailTime}}

	// Send Response
	responseBytes, err := json.Marshal(responseJson)
	if err != nil {
		SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to process your Request.")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBytes)
}

// Returns a ResultsResponse containing an array of WebsiteCheckResults.
func ApiWebsitesResults(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Get limit-parameter from Request
	limit := 100
	limitString := r.URL.Query().Get("limit")
	if len(limitString) != 0 {
		parsedLimit, err := strconv.Atoi(limitString)
		if err != nil {
			SendJsonMessage(w, http.StatusBadRequest, false, "Unable to parse given limit-parameter.")
			return
		}
		if parsedLimit > 9999 {
			SendJsonMessage(w, http.StatusBadRequest, false, "Unable to process your Request: Limit has to be less than 10000.")
			return
		}
		limit = parsedLimit
	}

	// Get offset-parameter from Request
	offset := 0
	offsetString := r.URL.Query().Get("offset")
	if len(offsetString) != 0 {
		parsedOffset, err := strconv.Atoi(offsetString)
		if err != nil {
			SendJsonMessage(w, http.StatusBadRequest, false, "Unable to parse given offset-parameter.")
			return
		}
		if parsedOffset > 9999 {
			SendJsonMessage(w, http.StatusBadRequest, false, "Unable to process your Request: Offset has to be less than 10000.")
			return
		}
		offset = parsedOffset
	}

	// Query the Database
	db := lib.GetDatabase()
	rows, err := db.Query("SELECT statusCode, statusText, responseTime, time FROM checks, websites WHERE checks.websiteId = websites.id AND websites.url = ? AND websites.enabled = 1 ORDER BY time DESC LIMIT ? OFFSET ?;", ps.ByName("url"), limit, offset)
	if err != nil {
		logging.MustGetLogger("").Error("Unable to fetch Results: ", err)
		SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to process your Request.")
		return
	}
	defer rows.Close()

	// Add every Result
	results := []WebsiteCheckResult{}
	var (
		statusCode   string
		statusText   string
		responseTime string
		time         string
	)
	for rows.Next() {
		err = rows.Scan(&statusCode, &statusText, &responseTime, &time)
		if err != nil {
			logging.MustGetLogger("").Error("Unable to read Result-Row: ", err)
			SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to process your Request.")
			return
		}

		results = append(results, WebsiteCheckResult{statusCode + " - " + statusText, responseTime, time})
	}

	// Check for Errors
	err = rows.Err()
	if err != nil {
		logging.MustGetLogger("").Error("Unable to read Result-Rows: ", err)
		SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to process your Request.")
		return
	}

	// Send Response
	responseBytes, err := json.Marshal(ResultsResponse{true, results})
	if err != nil {
		SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to process your Request.")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBytes)
}
