package main

import (
	"database/sql"
	"fmt"
	"github.com/MarvinMenzerath/UpAndRunning2/routes"
	"github.com/MarvinMenzerath/UpAndRunning2/tools"
	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"os"
	"strconv"
)

const VERSION = "1.0"

var Config tools.Configuration
var Database sql.DB

func main() {
	fmt.Printf("Welcome to UpAndRunning2 v%s!\n\n", VERSION)

	readConfig()
	prepareDatabase()

	go startChecking()
	serveRequests()
}

func readConfig() {
	fmt.Println("Reading Configuration...")
	Config = tools.GetConfig("config/local.json")
}

func prepareDatabase() {
	fmt.Println("Preaparing Database...")

	// username:password@protocol(address)/dbname
	Database, err := sql.Open("mysql", Config.Database.User+":"+Config.Database.Password+"@tcp("+Config.Database.Host+":"+strconv.Itoa(Config.Database.Port)+")/"+Config.Database.Database)
	if err != nil {
		fmt.Println("Unable to open Database-Connection: ", err)
		os.Exit(1)
	}
	Database.SetMaxOpenConns(Config.Database.ConnectionLimit)
	defer Database.Close()

	err = Database.Ping()
	if err != nil {
		fmt.Println("Unable to reach Database: ", err)
		os.Exit(1)
	}

	_, err = Database.Query("CREATE TABLE IF NOT EXISTS `website` (`id` int(11) NOT NULL AUTO_INCREMENT, `name` varchar(50) NOT NULL, `enabled` int(1) NOT NULL DEFAULT '1', `visible` int(1) NOT NULL DEFAULT '1', `protocol` varchar(8) NOT NULL, `url` varchar(100) NOT NULL, `status` varchar(50) NOT NULL DEFAULT 'unknown', `time` datetime NOT NULL DEFAULT '0000-00-00 00:00:00', `lastFailStatus` varchar(50) NOT NULL DEFAULT 'unknown', `lastFailTime` datetime NOT NULL DEFAULT '0000-00-00 00:00:00', `ups` int(11) NOT NULL DEFAULT '0', `downs` int(11) NOT NULL DEFAULT '0', `totalChecks` int(11) NOT NULL DEFAULT '0', `avgAvail` float NOT NULL DEFAULT '100', PRIMARY KEY (`id`), UNIQUE KEY `url` (`url`)) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;")
	if err != nil {
		fmt.Println("Unable to create table 'website': ", err)
		os.Exit(1)
	}

	_, err = Database.Query("CREATE TABLE IF NOT EXISTS `settings` (`id` int(11) NOT NULL AUTO_INCREMENT, `name` varchar(20) NOT NULL, `value` varchar(1024) NOT NULL, PRIMARY KEY (`id`), UNIQUE KEY `name` (`name`)) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;")
	if err != nil {
		fmt.Println("Unable to create table 'settings': ", err)
		os.Exit(1)
	}

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

	fmt.Println("Listening on Port " + strconv.Itoa(Config.Port) + "...")
	http.ListenAndServe(":"+strconv.Itoa(Config.Port), router)
}

func startChecking() {
	CheckAllSites()
}

func CheckAllSites() {
	fmt.Println("Checking X active Websites...")
}
