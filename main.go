package main

import (
	"fmt"
	"github.com/MarvinMenzerath/UpAndRunning2/lib"
	"github.com/MarvinMenzerath/UpAndRunning2/routes"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

const VERSION = "1.0"

func main() {
	fmt.Printf("Welcome to UpAndRunning2 v%s!\n\n", VERSION)

	lib.ReadConfigFromFile("config/local.json")
	lib.PrepareDatabase()
	lib.ReadConfigFromDatabase()

	go startChecking()
	serveRequests()
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

	fmt.Println("Listening on Port " + strconv.Itoa(lib.Config.Port) + "...")
	http.ListenAndServe(":"+strconv.Itoa(lib.Config.Port), router)
}

func startChecking() {
	CheckAllSites()
}

func CheckAllSites() {
	fmt.Println("Checking X active Websites...")
}
