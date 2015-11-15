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
	Title         string
	Interval      int
	PushbulletKey string
	CheckNow      bool
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

	var title string
	var interval int
	var pushbulletKey string

	// Title
	err := db.QueryRow("SELECT value FROM settings where name = 'title';").Scan(&title)
	if err != nil {
		stmt, err := db.Prepare("INSERT INTO settings (name, value) VALUES ('title', 'UpAndRunning2');")
		if err != nil {
			logging.MustGetLogger("logger").Fatal("Unable to insert 'title'-setting: ", err)
		}
		_, err = stmt.Exec()
		if err != nil {
			logging.MustGetLogger("logger").Fatal("Unable to insert 'title'-setting: ", err)
		}
		title = "UpAndRunning"
	}

	// Interval
	err = db.QueryRow("SELECT value FROM settings where name = 'interval';").Scan(&interval)
	if err != nil {
		stmt, err := db.Prepare("INSERT INTO settings (name, value) VALUES ('interval', 5);")
		if err != nil {
			logging.MustGetLogger("logger").Fatal("Unable to insert 'interval'-setting: ", err)
		}
		_, err = stmt.Exec()
		if err != nil {
			logging.MustGetLogger("logger").Fatal("Unable to insert 'interval'-setting: ", err)
		}
		interval = 5
	}

	// Pushbullet-Key
	err = db.QueryRow("SELECT value FROM settings where name = 'pushbullet_key';").Scan(&pushbulletKey)
	if err != nil {
		stmt, err := db.Prepare("INSERT INTO settings (name, value) VALUES ('pushbullet_key', '');")
		if err != nil {
			logging.MustGetLogger("logger").Fatal("Unable to insert 'pushbullet_key'-setting: ", err)
		}
		_, err = stmt.Exec()
		if err != nil {
			logging.MustGetLogger("logger").Fatal("Unable to insert 'pushbullet_key'-setting: ", err)
		}
		pushbulletKey = ""
	}

	config.Dynamic.Title = title
	config.Dynamic.Interval = interval
	config.Dynamic.PushbulletKey = pushbulletKey

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
