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
	Address  string
	Port     int
	Database databaseConfiguration
	Mailer   mailerConfiguration
	Dynamic  dynamicConfiguration
	Static   StaticConfiguration
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
	Title                string
	Interval             int
	Redirects            int
	RunChecksWhenOffline int
	CleanDatabase        int
	CheckNow             bool
}

// Static data about (e.g.) the application's version.
type StaticConfiguration struct {
	Version   string
	GoVersion string
	GoArch    string
}

// Reads a configuration-file from a specified path.
func ReadConfigurationFromFile(filePath string) {
	if os.Getenv("UAR2_IS_DOCKER") == "true" {
		ReadDefaultConfiguration("config/default.json")
		logging.MustGetLogger("").Info("Reading Configuration from Environment Variables...")

		config.Database.Host = os.Getenv("MYSQL_PORT_3306_TCP_ADDR")
		config.Database.User = "root"
		config.Database.Password = os.Getenv("MYSQL_ENV_MYSQL_ROOT_PASSWORD")
		config.Database.Database = "upandrunning"

		config.Mailer.Host = os.Getenv("UAR2_MAILER_HOST")
		i, err := strconv.Atoi(os.Getenv("UAR2_MAILER_PORT"))
		if err == nil {
			config.Mailer.Port = i
		}
		config.Mailer.User = os.Getenv("UAR2_MAILER_USER")
		config.Mailer.Password = os.Getenv("UAR2_MAILER_PASSWORD")
		config.Mailer.From = os.Getenv("UAR2_MAILER_FROM")

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
		title                string
		interval             int
		redirects            int
		runChecksWhenOffline int
		cleanDatabase        int
	)

	// Title
	err := db.QueryRow("SELECT value FROM settings where name = 'title';").Scan(&title)
	if err != nil {
		_, err = db.Exec("INSERT INTO settings (name, value) VALUES ('title', 'UpAndRunning2');")
		if err != nil {
			logging.MustGetLogger("").Fatal("Unable to insert 'title'-setting: ", err)
		}
		title = "UpAndRunning"
	}

	// Interval
	err = db.QueryRow("SELECT value FROM settings where name = 'interval';").Scan(&interval)
	if err != nil {
		_, err = db.Exec("INSERT INTO settings (name, value) VALUES ('interval', 60);")
		if err != nil {
			logging.MustGetLogger("").Fatal("Unable to insert 'interval'-setting: ", err)
		}
		interval = 60
	}

	// Redirects
	err = db.QueryRow("SELECT value FROM settings where name = 'redirects';").Scan(&redirects)
	if err != nil {
		_, err = db.Exec("INSERT INTO settings (name, value) VALUES ('redirects', 0);")
		if err != nil {
			logging.MustGetLogger("").Fatal("Unable to insert 'redirects'-setting: ", err)
		}
		redirects = 0
	}

	// Run Checks when offline
	err = db.QueryRow("SELECT value FROM settings where name = 'check_when_offline';").Scan(&runChecksWhenOffline)
	if err != nil {
		_, err = db.Exec("INSERT INTO settings (name, value) VALUES ('check_when_offline', 1);")
		if err != nil {
			logging.MustGetLogger("").Fatal("Unable to insert 'check_when_offline'-setting: ", err)
		}
		runChecksWhenOffline = 1
	}

	// Regularly clean old check-results from database
	err = db.QueryRow("SELECT value FROM settings where name = 'clean_database';").Scan(&cleanDatabase)
	if err != nil {
		_, err = db.Exec("INSERT INTO settings (name, value) VALUES ('clean_database', 1);")
		if err != nil {
			logging.MustGetLogger("").Fatal("Unable to insert 'clean_database'-setting: ", err)
		}
		cleanDatabase = 1
	}

	config.Dynamic.Title = title
	config.Dynamic.Interval = interval
	config.Dynamic.Redirects = redirects
	config.Dynamic.RunChecksWhenOffline = runChecksWhenOffline
	config.Dynamic.CleanDatabase = cleanDatabase

	config.Dynamic.CheckNow = true
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
