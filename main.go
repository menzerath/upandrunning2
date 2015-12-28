package main

import (
	"github.com/MarvinMenzerath/UpAndRunning2/lib"
	"github.com/MarvinMenzerath/UpAndRunning2/routes"
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
	logging.MustGetLogger("logger").Info("Welcome to UpAndRunning2 v%s [%s@%s]!", VERSION, goVersion, goArch)

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
	router.GET("/", routes.Index)
	router.GET("/status/:url", routes.Index)

	// Admin
	router.GET("/admin", routes.AdminIndex)
	router.GET("/admin/login", routes.AdminLogin)

	// **********
	// * API v1 *
	// **********

	// Public
	router.GET("/api", routes.ApiIndex)
	router.GET("/api/v1", routes.ApiIndexV1)
	router.GET("/api/v1/websites", routes.ApiWebsites)
	router.GET("/api/v1/websites/:url/status", routes.ApiStatus)
	router.GET("/api/v1/websites/:url/results", routes.ApiResults)

	// Private
	router.PUT("/api/v1/settings/title", routes.ApiAdminSettingTitle)
	router.PUT("/api/v1/settings/password", routes.ApiAdminSettingPassword)
	router.PUT("/api/v1/settings/interval", routes.ApiAdminSettingInterval)
	router.PUT("/api/v1/settings/redirects", routes.ApiAdminSettingRedirects)

	router.GET("/api/v1/check", routes.ApiAdminActionCheck)
	router.POST("/api/v1/auth/login", routes.ApiAdminActionLogin)
	router.GET("/api/v1/auth/logout", routes.ApiAdminActionLogout)

	router.GET("/api/v1/admin/websites", routes.ApiAdminWebsites)
	router.POST("/api/v1/admin/websites/:url", routes.ApiAdminWebsiteAdd)
	router.PUT("/api/v1/admin/websites/:url", routes.ApiAdminWebsiteEdit)
	router.DELETE("/api/v1/admin/websites/:url", routes.ApiAdminWebsiteDelete)

	router.PUT("/api/v1/admin/websites/:url/enabled", routes.ApiAdminWebsiteEnabled)
	router.PUT("/api/v1/admin/websites/:url/visibility", routes.ApiAdminWebsiteVisibility)
	router.GET("/api/v1/admin/websites/:url/notifications", routes.ApiAdminWebsiteGetNotifications)
	router.PUT("/api/v1/admin/websites/:url/notifications", routes.ApiAdminWebsiteUpdateNotifications)

	// Static Files
	router.ServeFiles("/public/*filepath", http.Dir("public"))

	// 404 Handler
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Error 404: Not Found", 404)
	})

	logging.MustGetLogger("logger").Debug("Listening on Port " + strconv.Itoa(lib.GetConfiguration().Port) + "...")
	logging.MustGetLogger("logger").Fatal(http.ListenAndServe(":"+strconv.Itoa(lib.GetConfiguration().Port), router))
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
		logging.MustGetLogger("logger").Error("Unable to fetch Websites: ", err)
		return
	}
	defer rows.Close()

	// Check every Website
	count := 0
	for rows.Next() {
		var website lib.Website
		err = rows.Scan(&website.Id, &website.Protocol, &website.Url, &website.CheckMethod)
		if err != nil {
			logging.MustGetLogger("logger").Error("Unable to read Website-Row: ", err)
			return
		}
		go website.RunCheck(false)
		count++
		time.Sleep(time.Millisecond * 200)
	}

	// Check for Errors
	err = rows.Err()
	if err != nil {
		logging.MustGetLogger("logger").Error("Unable to read Website-Rows: ", err)
		return
	}

	logging.MustGetLogger("logger").Info("Checked " + strconv.Itoa(count) + " active Websites.")
}
