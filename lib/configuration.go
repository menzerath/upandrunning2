package lib

import (
	"encoding/json"
	"fmt"
	"os"
)

var Config Configuration

type Configuration struct {
	Port     int
	Database databaseConfiguration
	Dynamic  dynamicConfiguration
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

func ReadConfigFromFile(filePath string) {
	fmt.Println("Reading Configuration from File...")

	file, _ := os.Open(filePath)
	decoder := json.NewDecoder(file)

	err := decoder.Decode(&Config)
	if err != nil {
		fmt.Println("Unable to read Configuration: ", err)
	}
}

func ReadConfigFromDatabase() {
	fmt.Println("Reading Configuration from Database...")

	var title string
	var interval int
	var pushbulletKey string

	err := Database.QueryRow("SELECT value FROM settings where name = 'title';").Scan(&title)
	if err != nil {
		stmt, err := Database.Prepare("INSERT INTO settings (name, value) VALUES ('title', 'UpAndRunning');")
		if err != nil {
			fmt.Println("Unable to insert 'title'-setting: ", err)
			os.Exit(1)
		}
		_, err = stmt.Exec()
		if err != nil {
			fmt.Println("Unable to insert 'title'-setting: ", err)
			os.Exit(1)
		}
		title = "UpAndRunning"
	}

	err = Database.QueryRow("SELECT value FROM settings where name = 'interval';").Scan(&interval)
	if err != nil {
		stmt, err := Database.Prepare("INSERT INTO settings (name, value) VALUES ('interval', 5);")
		if err != nil {
			fmt.Println("Unable to insert 'interval'-setting: ", err)
			os.Exit(1)
		}
		_, err = stmt.Exec()
		if err != nil {
			fmt.Println("Unable to insert 'interval'-setting: ", err)
			os.Exit(1)
		}
		interval = 5
	}

	err = Database.QueryRow("SELECT value FROM settings where name = 'pushbullet_key';").Scan(&pushbulletKey)
	if err != nil {
		stmt, err := Database.Prepare("INSERT INTO settings (name, value) VALUES ('pushbullet_key', '');")
		if err != nil {
			fmt.Println("Unable to insert 'pushbullet_key'-setting: ", err)
			os.Exit(1)
		}
		_, err = stmt.Exec()
		if err != nil {
			fmt.Println("Unable to insert 'pushbullet_key'-setting: ", err)
			os.Exit(1)
		}
		pushbulletKey = ""
	}
}
