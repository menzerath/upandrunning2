package APIv2

import (
	"github.com/MarvinMenzerath/UpAndRunning2/lib"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// Processes a login-request and sends an authentication-cookie to the browser.
func ApiAuthLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if lib.IsLoggedIn(r) {
		SendJsonMessage(w, http.StatusBadRequest, false, "Already logged in.")
		return
	}

	// Get data from Request
	r.ParseForm()
	value := r.Form.Get("password")

	// Check Password
	admin := lib.Admin{}
	if admin.ValidatePassword(value) {
		cookie := lib.LoginAndGetCookie("admin")
		http.SetCookie(w, &cookie)
		SendJsonMessage(w, http.StatusOK, true, "")
	} else {
		SendJsonMessage(w, http.StatusBadRequest, false, "Unable to process your Request: Invalid Password.")
	}
}

// Processes a logout-request and sends a termination-cookie to the browser.
func ApiAuthLogout(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !lib.IsLoggedIn(r) {
		SendJsonMessage(w, http.StatusUnauthorized, false, "Unauthorized.")
		return
	}

	// Process logout
	cookie := lib.LogoutAndDestroyCookie(r)
	http.SetCookie(w, &cookie)
	SendJsonMessage(w, http.StatusOK, true, "")
}
