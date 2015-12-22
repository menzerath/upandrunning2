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

	// v2.0.2
	_, err := db.Exec("ALTER TABLE `website` RENAME `websites`;")
	if mysqlError, ok := err.(*mysql.MySQLError); ok {
		if mysqlError.Number != 1051 && mysqlError.Number != 1146 && mysqlError.Number != 1050 { // Table does not exist: no need to change it
			logging.MustGetLogger("logger").Fatal("Unable to rename 'website'-table to 'websites': ", err)
		}
	}

	// v2.0.2
	_, err = db.Exec("ALTER TABLE `websites` ADD `responseTime` VARCHAR(50) NOT NULL DEFAULT 'unknown' AFTER `status`;")
	if mysqlError, ok := err.(*mysql.MySQLError); ok {
		if mysqlError.Number != 1060 { // Column exists: no need to add it
			logging.MustGetLogger("logger").Fatal("Unable to add 'responseTime'-column: ", err)
		}
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS `websites` (`id` int(11) NOT NULL AUTO_INCREMENT, `name` varchar(50) NOT NULL, `enabled` int(1) NOT NULL DEFAULT '1', `visible` int(1) NOT NULL DEFAULT '1', `protocol` varchar(8) NOT NULL DEFAULT 'https', `url` varchar(100) NOT NULL, `checkMethod` VARCHAR(10) NOT NULL DEFAULT 'HEAD', `status` varchar(50) NOT NULL DEFAULT 'unknown', `time` datetime NOT NULL DEFAULT '0000-00-00 00:00:00', `lastFailStatus` varchar(50) NOT NULL DEFAULT 'unknown', `lastFailTime` datetime NOT NULL DEFAULT '0000-00-00 00:00:00', `ups` int(11) NOT NULL DEFAULT '0', `downs` int(11) NOT NULL DEFAULT '0', `totalChecks` int(11) NOT NULL DEFAULT '0', `avgAvail` float NOT NULL DEFAULT '100', PRIMARY KEY (`id`), UNIQUE KEY `url` (`url`)) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;")
	if err != nil {
		logging.MustGetLogger("logger").Fatal("Unable to create table 'websites': ", err)
	}

	// v2.0.0
	_, err = db.Exec("ALTER TABLE `websites` ADD `checkMethod` VARCHAR(10) NOT NULL DEFAULT 'HEAD' AFTER `url`;")
	if mysqlError, ok := err.(*mysql.MySQLError); ok {
		if mysqlError.Number != 1060 { // Column exists: no need to add it
			logging.MustGetLogger("logger").Fatal("Unable to add 'checkMethod'-column: ", err)
		}
	}

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
