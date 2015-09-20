package main

import (
	"database/sql"
	"github.com/MarvinMenzerath/UpAndRunning2/lib"
	"github.com/MarvinMenzerath/UpAndRunning2/routes"
	"github.com/julienschmidt/httprouter"
	"github.com/op/go-logging"
	"net/http"
	"strconv"
)

const VERSION = "1.0"

var Admin lib.Admin
var Config *lib.Configuration
var Database *sql.DB

func main() {
	// Logger
	lib.SetupLogger()

	// Welcome
	logging.MustGetLogger("logger").Info("Welcome to UpAndRunning2 v%s!", VERSION)

	// Config
	lib.ReadConfigurationFromFile("config/local.json")
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
	go startChecking()
	serveRequests()

	lib.GetDatabase().Close()
}

func serveRequests() {
	router := httprouter.New()

	// Index
	router.GET("/", routes.IndexIndex)
	router.GET("/status/:url", routes.IndexStatus)

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

func startChecking() {
	CheckAllSites()
}

func CheckAllSites() {
	logging.MustGetLogger("logger").Info("Checking X active Websites...")
}
