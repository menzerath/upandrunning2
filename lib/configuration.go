package lib

import (
	"database/sql"
	"encoding/json"
	"github.com/op/go-logging"
	"os"
)

// This the one and only Configuration-object.
var config *Configuration

// The whole configuration.
// Contains all other configuration-data.
type Configuration struct {
	Port     int
	Database databaseConfiguration
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

// A dynamic Configuration.
// Used to store data, which may be changed through the API.
type dynamicConfiguration struct {
	Title     string
	Interval  int
	Redirects int
	CheckNow  bool
}

// Static data about (e.g.) the application's version.
type StaticConfiguration struct {
	Version   string
	GoVersion string
	GoArch    string
}

// Reads a configuration-file from a specified path.
func ReadConfigurationFromFile(filePath string) {
	logging.MustGetLogger("logger").Info("Reading Configuration from File...")

	file, _ := os.Open(filePath)
	decoder := json.NewDecoder(file)

	err := decoder.Decode(&config)

	if err != nil {
		logging.MustGetLogger("logger").Fatal("Unable to read Configuration: ", err)
	}
}

// Reads all configuration-data from the database.
func ReadConfigurationFromDatabase(db *sql.DB) {
	logging.MustGetLogger("logger").Info("Reading Configuration from Database...")

	var (
		title     string
		interval  int
		redirects int
	)

	// Title
	err := db.QueryRow("SELECT value FROM settings where name = 'title';").Scan(&title)
	if err != nil {
		_, err = db.Exec("INSERT INTO settings (name, value) VALUES ('title', 'UpAndRunning2');")
		if err != nil {
			logging.MustGetLogger("logger").Fatal("Unable to insert 'title'-setting: ", err)
		}
		title = "UpAndRunning"
	}

	// Interval
	err = db.QueryRow("SELECT value FROM settings where name = 'interval';").Scan(&interval)
	if err != nil {
		_, err = db.Exec("INSERT INTO settings (name, value) VALUES ('interval', 30);")
		if err != nil {
			logging.MustGetLogger("logger").Fatal("Unable to insert 'interval'-setting: ", err)
		}
		interval = 5
	}

	// Redirects
	err = db.QueryRow("SELECT value FROM settings where name = 'redirects';").Scan(&redirects)
	if err != nil {
		_, err = db.Exec("INSERT INTO settings (name, value) VALUES ('redirects', 0);")
		if err != nil {
			logging.MustGetLogger("logger").Fatal("Unable to insert 'redirects'-setting: ", err)
		}
		redirects = 0
	}

	config.Dynamic.Title = title
	config.Dynamic.Interval = interval
	config.Dynamic.Redirects = redirects

	config.Dynamic.CheckNow = true
}

// Allows to replace the current StaticConfiguration.
func SetStaticConfiguration(c StaticConfiguration) {
	config.Static = c
}

// Returns the current Configuration-object.
func GetConfiguration() *Configuration {
	if config == nil {
		logging.MustGetLogger("logger").Fatal("No active Configuration found.")
	} else {
		return config
	}
	return nil
}
