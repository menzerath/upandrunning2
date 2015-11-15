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

const VERSION = "2.0.0 Beta"

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

	// API
	router.GET("/api", routes.ApiIndex)
	router.GET("/api/status/:url", routes.ApiStatus)
	router.GET("/api/websites", routes.ApiWebsites)

	// Admin-API
	router.GET("/api/admin", routes.ApiAdminIndex)

	router.POST("/api/admin/settings/title", routes.ApiAdminSettingTitle)
	router.POST("/api/admin/settings/password", routes.ApiAdminSettingPassword)
	router.POST("/api/admin/settings/interval", routes.ApiAdminSettingInterval)
	router.POST("/api/admin/settings/pbkey", routes.ApiAdminSettingPushbulletKey)

	router.GET("/api/admin/websites", routes.ApiAdminWebsites)
	router.POST("/api/admin/websites/add", routes.ApiAdminWebsiteAdd)
	router.POST("/api/admin/websites/enable", routes.ApiAdminWebsiteEnable)
	router.POST("/api/admin/websites/disable", routes.ApiAdminWebsiteDisable)
	router.POST("/api/admin/websites/visible", routes.ApiAdminWebsiteVisible)
	router.POST("/api/admin/websites/invisible", routes.ApiAdminWebsiteInvisible)
	router.POST("/api/admin/websites/edit", routes.ApiAdminWebsiteEdit)
	router.POST("/api/admin/websites/delete", routes.ApiAdminWebsiteDelete)

	router.POST("/api/admin/check", routes.ApiAdminActionCheck)
	router.POST("/api/admin/login", routes.ApiAdminActionLogin)
	router.POST("/api/admin/logout", routes.ApiAdminActionLogout)

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
	rows, err := db.Query("SELECT id, protocol, url, checkMethod FROM website WHERE enabled = 1;")
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
		go website.RunCheck()
		count++
	}

	// Check for Errors
	err = rows.Err()
	if err != nil {
		logging.MustGetLogger("logger").Error("Unable to read Website-Rows: ", err)
		return
	}

	logging.MustGetLogger("logger").Info("Checking " + strconv.Itoa(count) + " active Websites...")
}
