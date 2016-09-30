package lib

import (
	"database/sql"
	"encoding/json"
	"github.com/op/go-logging"
	"os"
	"strconv"
)

// This the one and only Configuration-object.
var config *Configuration

// The whole configuration.
// Contains all other configuration-data.
type Configuration struct {
	Address      string
	Port         int
	Database     databaseConfiguration
	Application  applicationConfiguration
	Notification notificationConfiguration
	Dynamic      dynamicConfiguration
	Static       StaticConfiguration
}

// The database configuration.
type databaseConfiguration struct {
	Host            string
	Port            int
	User            string
	Password        string
	Database        string
	ConnectionLimit int
}

type applicationConfiguration struct {
	Title             string
	RedirectsToFollow int
	RunCheckIfOffline bool
	CheckLifetime     int
	UseWebFrontend    bool
}

// The notification configuration.
type notificationConfiguration struct {
	Mailer            mailerConfiguration
	TelegramBotApiKey string
}

// The mailer configuration.
type mailerConfiguration struct {
	Host     string
	Port     int
	User     string
	Password string
	From     string
}

// A dynamic Configuration.
// Used to store data, which may be changed through the API.
type dynamicConfiguration struct {
	Interval int
}

// Static data about (e.g.) the application's version.
type StaticConfiguration struct {
	Version   string
	GoVersion string
	GoArch    string
}

// Reads a configuration-file from a specified path.
func ReadConfigurationFromFile(filePath string) {
	ReadDefaultConfiguration("config/default.json")

	if os.Getenv("UAR2_IS_DOCKER") == "true" {
		logging.MustGetLogger("").Info("Reading Configuration from Environment Variables...")

		config.Database.Host = os.Getenv("MYSQL_PORT_3306_TCP_ADDR")
		config.Database.User = "root"
		config.Database.Password = os.Getenv("MYSQL_ENV_MYSQL_ROOT_PASSWORD")
		config.Database.Database = "upandrunning"

		config.Application.Title = os.Getenv("UAR2_APPLICATION_TITLE")

		i, err := strconv.Atoi(os.Getenv("UAR2_REDIRECTSTOFOLLOW"))
		if err == nil {
			config.Application.RedirectsToFollow = i
		}

		b, err := strconv.ParseBool(os.Getenv("UAR2_CHECKIFOFFLINE"))
		if err == nil {
			config.Application.RunCheckIfOffline = b
		}

		i, err = strconv.Atoi(os.Getenv("UAR2_CHECKLIFETIME"))
		if err == nil {
			config.Application.CheckLifetime = i
		}

		b, err = strconv.ParseBool(os.Getenv("UAR2_USEWEBFRONTEND"))
		if err == nil {
			config.Application.UseWebFrontend = b
		}

		config.Notification.Mailer.Host = os.Getenv("UAR2_MAILER_HOST")
		i, err = strconv.Atoi(os.Getenv("UAR2_MAILER_PORT"))
		if err == nil {
			config.Notification.Mailer.Port = i
		}
		config.Notification.Mailer.User = os.Getenv("UAR2_MAILER_USER")
		config.Notification.Mailer.Password = os.Getenv("UAR2_MAILER_PASSWORD")
		config.Notification.Mailer.From = os.Getenv("UAR2_MAILER_FROM")

		config.Notification.TelegramBotApiKey = os.Getenv("UAR2_TELEGRAMBOTAPIKEY")

		return
	}

	logging.MustGetLogger("").Info("Reading Configuration from File (" + filePath + ")...")

	file, _ := os.Open(filePath)
	decoder := json.NewDecoder(file)

	err := decoder.Decode(&config)
	if err != nil {
		logging.MustGetLogger("").Fatal("Unable to read Configuration. Make sure the File exists and is valid: ", err)
	}
}

// Reads the default configuration-file from a specified path.
func ReadDefaultConfiguration(filePath string) {
	logging.MustGetLogger("").Info("Reading Default-Configuration from File (" + filePath + ")...")

	file, _ := os.Open(filePath)
	decoder := json.NewDecoder(file)

	err := decoder.Decode(&config)
	if err != nil {
		logging.MustGetLogger("").Fatal("Unable to read Configuration. Make sure the File exists and is valid: ", err)
	}
}

// Reads all configuration-data from the database.
func ReadConfigurationFromDatabase(db *sql.DB) {
	logging.MustGetLogger("").Info("Reading Configuration from Database...")

	var (
		interval int
	)

	// Interval
	err := db.QueryRow("SELECT value FROM settings where name = 'interval';").Scan(&interval)
	if err != nil {
		_, err = db.Exec("INSERT INTO settings (name, value) VALUES ('interval', 60);")
		if err != nil {
			logging.MustGetLogger("").Fatal("Unable to insert 'interval'-setting: ", err)
		}
		interval = 60
	}

	config.Dynamic.Interval = interval
}

// Allows to replace the current StaticConfiguration.
func SetStaticConfiguration(c StaticConfiguration) {
	config.Static = c
}

// Returns the current Configuration-object.
func GetConfiguration() *Configuration {
	if config == nil {
		logging.MustGetLogger("").Fatal("No active Configuration found.")
	} else {
		return config
	}
	return nil
}
