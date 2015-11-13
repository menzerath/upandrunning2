package routes

import (
	"encoding/json"
	"net/http"
)

// Basic Response
type BasicResponse struct {
	Success bool   `json:"requestSuccess"`
	Message string `json:"message"`
}

// Website Response
type WebsiteResponse struct {
	Success  bool           `json:"requestSuccess"`
	Websites []BasicWebsite `json:"websites"`
}

type BasicWebsite struct {
	Name     string `json:"name"`
	Protocol string `json:"protocol"`
	Url      string `json:"url"`
	Status   string `json:"status"`
}

// Detailed Website-Status Response
type DetailedWebsiteResponse struct {
	Success               bool                `json:"requestSuccess"`
	WebsiteData           WebsiteData         `json:"websiteData"`
	Availability          WebsiteAvailability `json:"availability"`
	LastCheckResult       WebsiteCheckResult  `json:"lastCheckResult"`
	LastFailedCheckResult WebsiteCheckResult  `json:"lastFailedCheckResult"`
}

type WebsiteData struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Url  string `json:"url"`
}

type WebsiteAvailability struct {
	Ups     int    `json:"ups"`
	Downs   int    `json:"downs"`
	Total   int    `json:"total"`
	Average string `json:"average"`
}

type WebsiteCheckResult struct {
	Status string `json:"status"`
	Time   string `json:"time"`
}

// Admin Website Response
type AdminWebsiteResponse struct {
	Success  bool           `json:"requestSuccess"`
	Websites []AdminWebsite `json:"websites"`
}

type AdminWebsite struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Enabled  bool   `json:"enabled"`
	Visible  bool   `json:"visible"`
	Protocol string `json:"protocol"`
	Url      string `json:"url"`
	Status   string `json:"status"`
	Time     string `json:"time"`
	Average  string `json:"avgAvail"`
}

// Site-Data
type SiteData struct {
	Title string
}

type AdminSiteData struct {
	Title      string
	Interval   int
	PbKey      string
	AppVersion string
	GoVersion  string
	GoArch     string
}

// Functions
func SendJsonMessage(w http.ResponseWriter, code int, success bool, message string) {
	responseBytes, err := json.Marshal(BasicResponse{success, message})
	if err != nil {
		SendJsonMessage(w, http.StatusInternalServerError, false, "Unable to send JSON-Message.")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(responseBytes)
}
