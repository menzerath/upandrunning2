package routes

import (
	"github.com/MarvinMenzerath/UpAndRunning2/lib"
	"github.com/julienschmidt/httprouter"
	"github.com/op/go-logging"
	"html/template"
	"net/http"
	"runtime"
)

// Sends a simple welcome-message to the user.
func NoWebFrontendIndex(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	SendJsonMessage(w, http.StatusOK, true, "Welcome to UpAndRunning2! Currently the Web-Frontend is disabled, but you can still use our simple API.")
}

// Renders the main-page
func ViewIndex(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Parse template-files
	data := SiteData{lib.GetConfiguration().Application.Title}
	t, err := template.ParseFiles("views/index.html", "views/partials/styles.html", "views/partials/footer.html", "views/partials/scripts.html")

	if t != nil {
		w.Header().Set("Content-Type", "text/html")
		t.Execute(w, data)
	} else {
		logging.MustGetLogger("").Error("Error while parsing Template: ", err)
		http.Error(w, "Error 500: Internal Server Error", http.StatusInternalServerError)
	}
}

// Renders the login-page if the user is not logged in.
// If the user is logged in, he will be redirected to the admin-backend.
func ViewLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if lib.IsLoggedIn(r) {
		http.Redirect(w, r, "/admin", http.StatusFound)
		return
	}

	// Parse template-files
	data := SiteData{lib.GetConfiguration().Application.Title}
	t, err := template.ParseFiles("views/login.html", "views/partials/styles.html", "views/partials/footer.html", "views/partials/scripts.html")

	if t != nil {
		w.Header().Set("Content-Type", "text/html")
		t.Execute(w, data)
	} else {
		logging.MustGetLogger("").Error("Error while parsing Template: ", err)
		http.Error(w, "Error 500: Internal Server Error", http.StatusInternalServerError)
	}
}

// Renders the admin-backend if the user is logged in.
// If the user is not logged in, he will be redirected to the login-page.
func ViewAdmin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !lib.IsLoggedIn(r) {
		http.Redirect(w, r, "/admin/login", http.StatusFound)
		return
	}

	// Parse template-files
	c := lib.GetConfiguration()

	dynCheckWhenOffline := ""
	if c.Dynamic.RunChecksWhenOffline == 1 {
		dynCheckWhenOffline = "checked"
	}

	data := AdminSiteData{c.Application.Title, c.Dynamic.Interval, c.Dynamic.Redirects, dynCheckWhenOffline, c.Static.Version, runtime.Version(), runtime.GOOS + "_" + runtime.GOARCH}
	t, err := template.ParseFiles("views/admin.html", "views/partials/styles.html", "views/partials/footer.html", "views/partials/scripts.html")

	if t != nil {
		w.Header().Set("Content-Type", "text/html")
		t.Execute(w, data)
	} else {
		logging.MustGetLogger("").Error("Error while parsing Template: ", err)
		http.Error(w, "Error 500: Internal Server Error", http.StatusInternalServerError)
	}
}
