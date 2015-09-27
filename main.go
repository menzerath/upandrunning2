package main

import (
	"database/sql"
	"github.com/MarvinMenzerath/UpAndRunning2/lib"
	"github.com/MarvinMenzerath/UpAndRunning2/routes"
	"github.com/julienschmidt/httprouter"
	"github.com/op/go-logging"
	"net/http"
	"runtime"
	"strconv"
	"time"
)

const VERSION = "1.0"

var goVersion = runtime.Version()
var goArch = runtime.GOOS + "_" + runtime.GOARCH

var Admin lib.Admin
var Config *lib.Configuration
var Database *sql.DB

func main() {
	// Logger
	lib.SetupLogger()

	// Welcome
	logging.MustGetLogger("logger").Info("Welcome to UpAndRunning2 v%s [%s@%s]!", VERSION, goVersion, goArch)

	// Config
	lib.ReadConfigurationFromFile("config/local.json")
	lib.SetStaticConfiguration(lib.StaticConfiguration{VERSION, goVersion, goArch})
	Config = lib.GetConfiguration()

	// Database
	lib.OpenDatabase(Config.Database)
	Database = lib.GetDatabase()

	// Config (again)
	lib.ReadConfigurationFromDatabase(Database)

	// Admin-User
	Admin = lib.Admin{}
	if !Admin.Exists() {
		Admin.Add()
	}

	// Additional Libraries
	lib.InitHttpStatusCodeMap()

	// Start Checking and Serving
	startCheckTimer()
	startCheckNowTimer()
	serveRequests()

	lib.GetDatabase().Close()
}

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

	router.POST("/api/admin/websites", routes.ApiAdminWebsites)
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

	logging.MustGetLogger("logger").Debug("Listening on Port " + strconv.Itoa(Config.Port) + "...")
	logging.MustGetLogger("logger").Fatal(http.ListenAndServe(":"+strconv.Itoa(Config.Port), router))
}

func startCheckTimer() {
	timer := time.NewTimer(time.Second * time.Duration(Config.Dynamic.Interval))
	go func() {
		<-timer.C
		checkAllSites()
		startCheckTimer()
	}()
}

func startCheckNowTimer() {
	timer := time.NewTimer(time.Second * 1)
	go func() {
		<-timer.C
		if Config.Dynamic.CheckNow {
			checkAllSites()
			Config.Dynamic.CheckNow = false
		}
		startCheckNowTimer()
	}()
}

func checkAllSites() {
	// Query the Database
	db := lib.GetDatabase()
	rows, err := db.Query("SELECT id, protocol, url FROM website WHERE enabled = 1;")
	if err != nil {
		logging.MustGetLogger("logger").Error("Unable to fetch Websites: ", err)
		return
	}
	defer rows.Close()

	// Check every Website
	count := 0
	for rows.Next() {
		var website lib.Website
		err = rows.Scan(&website.Id, &website.Protocol, &website.Url)
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
