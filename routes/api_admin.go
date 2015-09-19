package routes

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func ApiAdminIndex(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	response := BasicResponse{true, "Welcome to UpAndRunning's Admin-API! Please be aware that most operations need an incoming POST-request."}
	json.NewEncoder(w).Encode(response)
}

func ApiAdminWebsites(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

func ApiAdminWebsiteAdd(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

func ApiAdminWebsiteEnable(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

func ApiAdminWebsiteDisable(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

func ApiAdminWebsiteVisible(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

func ApiAdminWebsiteInvisible(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

func ApiAdminWebsiteEdit(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

func ApiAdminWebsiteDelete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

func ApiAdminSettingTitle(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

func ApiAdminSettingPassword(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

func ApiAdminSettingInterval(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

func ApiAdminSettingPushbulletKey(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

func ApiAdminActionCheck(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

func ApiAdminActionLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

func ApiAdminActionLogout(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}
