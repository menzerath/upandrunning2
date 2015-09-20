package routes

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func ApiIndex(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	responseJson := BasicResponse{true, "Welcome to UpAndRunning's API!"}

	responseBytes, err := json.Marshal(responseJson)
	if err != nil {
		http.Error(w, "Error 500: Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBytes)
}

func ApiStatus(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

func ApiWebsites(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}
