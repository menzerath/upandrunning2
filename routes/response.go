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

// ******************
// * VIEW-RESOURCES *
// ******************

// Contains the application's data, which will be used on publicly visible pages.
type SiteData struct {
	Title string
}

// Contains the application's data, which will be used on admin-pages.
type AdminSiteData struct {
	Title            string
	Interval         int
	AppVersion       string
	GoVersion        string
	GoArch           string
}

// *************
// * FUNCTIONS *
// *************

// Sends a simple Json-message.
// It contains a success-bool and a message, which may be empty.
func SendJsonMessage(w http.ResponseWriter, code int, success bool, message string) {
	responseBytes, err := json.Marshal(BasicResponse{success, message})
	if err != nil {
		logging.MustGetLogger("").Error("Unable to send JSON-Message: ", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(responseBytes)
}
