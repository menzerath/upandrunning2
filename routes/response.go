package routes

import (
	"encoding/json"
	"github.com/op/go-logging"
	"net/http"
)

// ******************
// * USER-RESPONSES *
// ******************

// Contains a success-bool and a message, which may be empty.
type BasicResponse struct {
	Success bool   `json:"requestSuccess"`
	Message string `json:"message"`
}

// Contains a success-bool and an array of BasicWebsites.
type WebsiteResponse struct {
	Success  bool           `json:"requestSuccess"`
	Websites []BasicWebsite `json:"websites"`
}

// Contains a success-bool and an array of WebsiteCheckResults.
type ResultsResponse struct {
	Success  bool           `json:"requestSuccess"`
	Websites []WebsiteCheckResult `json:"results"`
}

// Contains a success-bool and the Website's details.
type DetailedWebsiteResponse struct {
	Success               bool                `json:"requestSuccess"`
	WebsiteData           WebsiteData         `json:"websiteData"`
	Availability          WebsiteAvailability `json:"availability"`
	LastCheckResult       WebsiteCheckResult  `json:"lastCheckResult"`
	LastFailedCheckResult WebsiteCheckResult  `json:"lastFailedCheckResult"`
}

// Contains the Website's basic data such as name, protocol, url and current status.
type BasicWebsite struct {
	Name     string `json:"name"`
	Protocol string `json:"protocol"`
	Url      string `json:"url"`
	Status   string `json:"status"`
}

// Contains the Website's basic data such as id, name and url.
type WebsiteData struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Url  string `json:"url"`
}

// Contains the Website's availability data like the amount of ups / downs, total checks and the average availability.
type WebsiteAvailability struct {
	Ups     int    `json:"ups"`
	Downs   int    `json:"downs"`
	Total   int    `json:"total"`
	Average string `json:"average"`
}

// Contains the Website's latest check result.
type WebsiteCheckResult struct {
	Status       string `json:"status"`
	ResponseTime string `json:"responseTime"`
	Time         string `json:"time"`
}

// *******************
// * ADMIN-RESPONSES *
// *******************

// Contains a success-bool and an array of AdminWebsites.
type AdminWebsiteResponse struct {
	Success  bool           `json:"requestSuccess"`
	Websites []AdminWebsite `json:"websites"`
}

// Contains the Website's data, which will be shown inside the admin-backend.
type AdminWebsite struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Enabled     bool   `json:"enabled"`
	Visible     bool   `json:"visible"`
	Protocol    string `json:"protocol"`
	Url         string `json:"url"`
	CheckMethod string `json:"checkMethod"`
	Status      string `json:"status"`
	Time        string `json:"time"`
}

// Contains the application's data, which will be used on publicly visible pages.
type SiteData struct {
	Title string
}

// Contains the application's data, which will be used on admin-pages.
type AdminSiteData struct {
	Title      string
	Interval   int
	Redirects  int
	PbKey      string
	AppVersion string
	GoVersion  string
	GoArch     string
}

// *************
// * FUNCTIONS *
// *************

// Sends a simple Json-message.
// It contains a success-bool and a message, which may be empty.
func SendJsonMessage(w http.ResponseWriter, code int, success bool, message string) {
	responseBytes, err := json.Marshal(BasicResponse{success, message})
	if err != nil {
		logging.MustGetLogger("logger").Error("Unable to send JSON-Message: ", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(responseBytes)
}
