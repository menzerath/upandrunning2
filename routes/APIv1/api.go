package APIv1

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/op/go-logging"
	"net/http"
)

// Sends a simple welcome-message to the user.
func ApiIndexVersion(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	SendJsonMessage(w, http.StatusOK, true, "APIv1 has been removed in favor of APIv2. Please update your application!")
}

// Contains a success-bool and a message, which may be empty.
type BasicResponse struct {
	Success bool   `json:"requestSuccess"`
	Message string `json:"message"`
}

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
