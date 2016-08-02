package APIv2

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// Sends a simple welcome-message to the user.
func ApiIndex(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	SendJsonMessage(w, http.StatusOK, true, "Welcome to UpAndRunning2's API!")
}

// Sends a simple welcome-message to the user.
func ApiIndexVersion(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	SendJsonMessage(w, http.StatusOK, true, "Welcome to UpAndRunning2's API v2!")
}
