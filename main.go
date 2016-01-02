package main

import (
	"github.com/MarvinMenzerath/UpAndRunning2/lib"
	"github.com/MarvinMenzerath/UpAndRunning2/routes"
	"github.com/franela/goreq"
	"github.com/julienschmidt/httprouter"
	"github.com/op/go-logging"
	"net/http"
	"runtime"
	"strconv"
	"time"
)

const VERSION = "2.0.2"

var goVersion = runtime.Version()
var goArch = runtime.GOOS + "_" + runtime.GOARCH

// UpAndRunning2 Main - The application's entrance-point
func main() {
	// Logger
	lib.SetupLogger()

	// Welcome
	logging.MustGetLogger("").Info("Welcome to UpAndRunning2 v%s [%s@%s]!", VERSION, goVersion, goArch)

	// Config
	lib.ReadConfigurationFromFile("config/local.json")
	lib.SetStaticConfiguration(lib.StaticConfiguration{VERSION, goVersion, goArch})

	// Database
	lib.OpenDatabase(lib.GetConfiguration().Database)

	// Config (again)
	lib.ReadConfigurationFromDatabase(lib.GetDatabase())

	// Admin-User
	admin := lib.Admin{}
	admin.Init()

	// Session-Management
	lib.InitSessionManagement()

	// Additional Libraries
	goreq.SetConnectTimeout(5 * time.Second)
	lib.InitHttpStatusCodeMap()

	// Start Checking and Serving
	startCheckTimer()
	startCheckNowTimer()
	serveRequests()

	lib.GetDatabase().Close()
}

// Create all routes and start the HTTP-server
func serveRequests() {
	router := httprouter.New()

	// Index
	router.GET("/", routes.ViewIndex)
	router.GET("/status/:url", routes.ViewIndex)
	router.GET("/results/:url", routes.ViewIndex)

	// Admin
	router.GET("/admin", routes.ViewAdmin)
	router.GET("/admin/login", routes.ViewLogin)

	// **********
	// * API v1 *
	// **********

	// Just Text
	router.GET("/api", routes.ApiIndex)
	router.GET("/api/v1", routes.ApiIndexV1)

	// Public Statistics
	router.GET("/api/v1/websites", routes.ApiWebsites)
	router.GET("/api/v1/websites/:url/status", routes.ApiWebsitesStatus)
	router.GET("/api/v1/websites/:url/results", routes.ApiWebsitesResults)

	// Actions
	router.GET("/api/v1/action/check", routes.ApiActionCheck)

	// Authentication
	router.POST("/api/v1/auth/login", routes.ApiAuthLogin)
	router.GET("/api/v1/auth/logout", routes.ApiAuthLogout)

	// Settings
	router.PUT("/api/v1/settings/title", routes.ApiSettingsTitle)
	router.PUT("/api/v1/settings/password", routes.ApiSettingsPassword)
	router.PUT("/api/v1/settings/interval", routes.ApiSettingsInterval)
	router.PUT("/api/v1/settings/redirects", routes.ApiSettingsRedirects)

	// Website Management
	router.POST("/api/v1/websites/:url", routes.ApiWebsitesAdd)
	router.PUT("/api/v1/websites/:url", routes.ApiWebsitesEdit)
	router.DELETE("/api/v1/websites/:url", routes.ApiWebsitesDelete)
	router.PUT("/api/v1/websites/:url/enabled", routes.ApiWebsitesEnabled)
	router.PUT("/api/v1/websites/:url/visibility", routes.ApiWebsitesVisibility)
	router.GET("/api/v1/websites/:url/notifications", routes.ApiWebsitesGetNotifications)
	router.PUT("/api/v1/websites/:url/notifications", routes.ApiWebsitePutNotifications)

	// Static Files
	router.ServeFiles("/public/*filepath", http.Dir("public"))

	// 404 Handler
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Error 404: Not Found", 404)
	})

	logging.MustGetLogger("").Debug("Listening on Port " + strconv.Itoa(lib.GetConfiguration().Port) + "...")
	logging.MustGetLogger("").Fatal(http.ListenAndServe(":"+strconv.Itoa(lib.GetConfiguration().Port), router))
}

// Creates a timer to regularly check all Websites
func startCheckTimer() {
	timer := time.NewTimer(time.Second * time.Duration(lib.GetConfiguration().Dynamic.Interval))
	go func() {
		<-timer.C
		checkAllSites()
		startCheckTimer()
	}()
}

// Creates a timer to check all Websites when triggered through the API
func startCheckNowTimer() {
	timer := time.NewTimer(time.Second * 1)
	go func() {
		<-timer.C
		if lib.GetConfiguration().Dynamic.CheckNow {
			checkAllSites()
			lib.GetConfiguration().Dynamic.CheckNow = false
		}
		startCheckNowTimer()
	}()
}

// Checks all enabled Websites
func checkAllSites() {
	// Query the Database
	db := lib.GetDatabase()
	rows, err := db.Query("SELECT id, protocol, url, checkMethod FROM websites WHERE enabled = 1;")
	if err != nil {
		logging.MustGetLogger("").Error("Unable to fetch Websites: ", err)
		return
	}
	defer rows.Close()

	// Check every Website
	count := 0
	for rows.Next() {
		var website lib.Website
		err = rows.Scan(&website.Id, &website.Protocol, &website.Url, &website.CheckMethod)
		if err != nil {
			logging.MustGetLogger("").Error("Unable to read Website-Row: ", err)
			return
		}
		go website.RunCheck(false)
		count++
		time.Sleep(time.Millisecond * 200)
	}

	// Check for Errors
	err = rows.Err()
	if err != nil {
		logging.MustGetLogger("").Error("Unable to read Website-Rows: ", err)
		return
	}

	logging.MustGetLogger("").Info("Checked " + strconv.Itoa(count) + " active Websites.")
}
