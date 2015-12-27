package lib

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/op/go-logging"
	"strconv"
)

// This is the one and only Database-object.
var db *sql.DB

// Opens a new connection-pool to the database using the given databaseConfiguration.
func OpenDatabase(config databaseConfiguration) {
	logging.MustGetLogger("logger").Info("Opening Database...")
	var err error = nil

	// username:password@protocol(address)/dbname
	db, err = sql.Open("mysql", config.User+":"+config.Password+"@tcp("+config.Host+":"+strconv.Itoa(config.Port)+")/"+config.Database)
	if err != nil {
		logging.MustGetLogger("logger").Fatal("Unable to open Database-Connection: ", err)
	}
	db.SetMaxOpenConns(config.ConnectionLimit)

	err = db.Ping()
	if err != nil {
		logging.MustGetLogger("logger").Fatal("Unable to reach Database: ", err)
	}

	prepareDatabase()
}

// Creates the needed tables in the database.
func prepareDatabase() {
	logging.MustGetLogger("logger").Debug("Preparing Database...")

	// v2.1.0
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS `checks` (`id` INT(11) NOT NULL AUTO_INCREMENT, `websiteId` INT(11) NOT NULL, `statusCode` INT(3) NOT NULL, `statusText` VARCHAR(50) NOT NULL, `responseTime` VARCHAR(50) NOT NULL, `time` DATETIME NOT NULL, PRIMARY KEY (`id`), FOREIGN KEY (`websiteId`) REFERENCES websites(`id`)) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;")
	if err != nil {
		logging.MustGetLogger("logger").Fatal("Unable to create table 'checks': ", err)
	}

	// v2.1.0
	_, err = db.Exec("ALTER TABLE `websites` DROP `status`, DROP `responseTime`, DROP `time`, DROP `lastFailStatus`, DROP `lastFailTime`, DROP `ups`, DROP `downs`, DROP `totalChecks`, DROP `avgAvail`;")
	if mysqlError, ok := err.(*mysql.MySQLError); ok {
		if mysqlError.Number != 1091 { // Columns do not exists: no need to remove them
			logging.MustGetLogger("logger").Fatal("Unable to drop unneccessary columns: ", err)
		}
	}

	// Default Setup
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS `websites` (`id` INT(11) NOT NULL AUTO_INCREMENT, `name` VARCHAR(50) NOT NULL, `enabled` INT(1) NOT NULL DEFAULT '1', `visible` INT(1) NOT NULL DEFAULT '1', `protocol` VARCHAR(8) NOT NULL DEFAULT 'https', `url` VARCHAR(100) NOT NULL, `checkMethod` VARCHAR(10) NOT NULL DEFAULT 'HEAD', PRIMARY KEY (`id`), UNIQUE KEY `url` (`url`)) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;")
	if err != nil {
		logging.MustGetLogger("logger").Fatal("Unable to create table 'websites': ", err)
	}

	// Default Setup
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS `settings` (`id` int(11) NOT NULL AUTO_INCREMENT, `name` varchar(20) NOT NULL, `value` varchar(1024) NOT NULL, PRIMARY KEY (`id`), UNIQUE KEY `name` (`name`)) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;")
	if err != nil {
		logging.MustGetLogger("logger").Fatal("Unable to create table 'settings': ", err)
	}
}

// Returns the current Database-object.
func GetDatabase() *sql.DB {
	if db == nil {
		logging.MustGetLogger("logger").Fatal("No active Database found.")
	} else {
		return db
	}
	return nil
}
