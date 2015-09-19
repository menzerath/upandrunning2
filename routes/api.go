package routes

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func ApiIndex(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	response := BasicResponse{true, "Welcome to UpAndRunning's API!"}
	json.NewEncoder(w).Encode(response)
}

func ApiStatus(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

func ApiWebsites(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}
