package lib

import (
	"database/sql"
	"encoding/json"
	"github.com/op/go-logging"
	"os"
)

var config *Configuration

type Configuration struct {
	Port     int
	Database databaseConfiguration
	Dynamic  dynamicConfiguration
	Static   StaticConfiguration
}

type databaseConfiguration struct {
	Host            string
	Port            int
	User            string
	Password        string
	Database        string
	ConnectionLimit int
}

type dynamicConfiguration struct {
	Title         string
	Interval      int
	PushbulletKey string
}

type StaticConfiguration struct {
	Version   string
	GoVersion string
	GoArch    string
}

func ReadConfigurationFromFile(filePath string) {
	logging.MustGetLogger("logger").Info("Reading Configuration from File...")

	file, _ := os.Open(filePath)
	decoder := json.NewDecoder(file)

	err := decoder.Decode(&config)
	if err != nil {
		logging.MustGetLogger("logger").Fatal("Unable to read Configuration: ", err)
	}
}

func ReadConfigurationFromDatabase(db *sql.DB) {
	logging.MustGetLogger("logger").Info("Reading Configuration from Database...")

	var title string
	var interval int
	var pushbulletKey string

	err := db.QueryRow("SELECT value FROM settings where name = 'title';").Scan(&title)
	if err != nil {
		stmt, err := db.Prepare("INSERT INTO settings (name, value) VALUES ('title', 'UpAndRunning');")
		if err != nil {
			logging.MustGetLogger("logger").Fatal("Unable to insert 'title'-setting: ", err)
		}
		_, err = stmt.Exec()
		if err != nil {
			logging.MustGetLogger("logger").Fatal("Unable to insert 'title'-setting: ", err)
		}
		title = "UpAndRunning"
	}

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
}

func SetStaticConfiguration(c StaticConfiguration) {
	config.Static = c
}

func GetConfiguration() *Configuration {
	if config == nil {
		logging.MustGetLogger("logger").Fatal("No active Configuration found.")
	} else {
		return config
	}
	return nil
}
