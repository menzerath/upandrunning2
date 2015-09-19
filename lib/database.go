package lib

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"strconv"
)

var db *sql.DB

func OpenDatabase(config databaseConfiguration) {
	fmt.Println("Preaparing Database...")
	var err error = nil

	// username:password@protocol(address)/dbname
	db, err = sql.Open("mysql", config.User+":"+config.Password+"@tcp("+config.Host+":"+strconv.Itoa(config.Port)+")/"+config.Database)
	if err != nil {
		fmt.Println("Unable to open Database-Connection: ", err)
		os.Exit(1)
	}
	db.SetMaxOpenConns(config.ConnectionLimit)

	err = db.Ping()
	if err != nil {
		fmt.Println("Unable to reach Database: ", err)
		os.Exit(1)
	}

	prepareDatabase()
}

func prepareDatabase() {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS `website` (`id` int(11) NOT NULL AUTO_INCREMENT, `name` varchar(50) NOT NULL, `enabled` int(1) NOT NULL DEFAULT '1', `visible` int(1) NOT NULL DEFAULT '1', `protocol` varchar(8) NOT NULL, `url` varchar(100) NOT NULL, `status` varchar(50) NOT NULL DEFAULT 'unknown', `time` datetime NOT NULL DEFAULT '0000-00-00 00:00:00', `lastFailStatus` varchar(50) NOT NULL DEFAULT 'unknown', `lastFailTime` datetime NOT NULL DEFAULT '0000-00-00 00:00:00', `ups` int(11) NOT NULL DEFAULT '0', `downs` int(11) NOT NULL DEFAULT '0', `totalChecks` int(11) NOT NULL DEFAULT '0', `avgAvail` float NOT NULL DEFAULT '100', PRIMARY KEY (`id`), UNIQUE KEY `url` (`url`)) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;")
	if err != nil {
		fmt.Println("Unable to create table 'website': ", err)
		os.Exit(1)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS `settings` (`id` int(11) NOT NULL AUTO_INCREMENT, `name` varchar(20) NOT NULL, `value` varchar(1024) NOT NULL, PRIMARY KEY (`id`), UNIQUE KEY `name` (`name`)) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;")
	if err != nil {
		fmt.Println("Unable to create table 'settings': ", err)
		os.Exit(1)
	}
}

func GetDatabase() *sql.DB {
	if db == nil {
		fmt.Println("No active Database found.")
		os.Exit(1)
	} else {
		return db
	}
	return nil
}
