package lib

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
)

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

func (c *Configuration) ReadFromFile(filePath string) {
	fmt.Println("Reading Configuration from File...")

	file, _ := os.Open(filePath)
	decoder := json.NewDecoder(file)

	err := decoder.Decode(&c)
	if err != nil {
		fmt.Println("Unable to read Configuration: ", err)
	}
}

func (c *dynamicConfiguration) ReadFromDatabase(db *sql.DB) {
	fmt.Println("Reading Configuration from Database...")

	var title string
	var interval int
	var pushbulletKey string

	err := db.QueryRow("SELECT value FROM settings where name = 'title';").Scan(&title)
	if err != nil {
		stmt, err := db.Prepare("INSERT INTO settings (name, value) VALUES ('title', 'UpAndRunning');")
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

	err = db.QueryRow("SELECT value FROM settings where name = 'interval';").Scan(&interval)
	if err != nil {
		stmt, err := db.Prepare("INSERT INTO settings (name, value) VALUES ('interval', 5);")
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

	err = db.QueryRow("SELECT value FROM settings where name = 'pushbullet_key';").Scan(&pushbulletKey)
	if err != nil {
		stmt, err := db.Prepare("INSERT INTO settings (name, value) VALUES ('pushbullet_key', '');")
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

	c.Title = title
	c.Interval = interval
	c.PushbulletKey = pushbulletKey
}
