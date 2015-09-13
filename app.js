var express = require('express');
var session = require('express-session');
var bodyParser = require('body-parser');
var path = require('path');

var db = require('./lib/database');
var logger = require('./lib/logger');
var admin = require('./lib/admin');
var website = require('./lib/website');

var app = express();

// welcome!
logger.info("Welcome to UpAndRunning v" + require('./package.json').version + "!");

// create database
db.query("CREATE TABLE IF NOT EXISTS `website` (`id` int(11) NOT NULL AUTO_INCREMENT, `name` varchar(50) NOT NULL, `enabled` int(1) NOT NULL DEFAULT '1', `visible` int(1) NOT NULL DEFAULT '1', `protocol` varchar(8) NOT NULL, `url` varchar(100) NOT NULL, `status` varchar(50) NOT NULL DEFAULT 'unknown', `time` datetime NOT NULL DEFAULT '0000-00-00 00:00:00', `lastFailStatus` varchar(50) NOT NULL DEFAULT 'unknown', `lastFailTime` datetime NOT NULL DEFAULT '0000-00-00 00:00:00', `ups` int(11) NOT NULL DEFAULT '0', `downs` int(11) NOT NULL DEFAULT '0', `totalChecks` int(11) NOT NULL DEFAULT '0', `avgAvail` float NOT NULL DEFAULT '100', PRIMARY KEY (`id`), UNIQUE KEY `url` (`url`)) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;", function(err) {
	if (err) {
		logger.error(err);
		process.exit(1);
	} else {
		logger.info("Website-Database successfully prepared.");
	}
});

db.query("CREATE TABLE IF NOT EXISTS `settings` (`id` int(11) NOT NULL AUTO_INCREMENT, `name` varchar(20) NOT NULL, `value` varchar(1024) NOT NULL, PRIMARY KEY (`id`), UNIQUE KEY `name` (`name`)) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;", function(err) {
	if (err) {
		logger.error(err);
		process.exit(1);
	} else {
		logger.info("Settings-Database successfully prepared");

		db.query("SELECT value FROM settings where name = 'title';", function(err, rows) {
			if (err) {
				logger.error("Unable to check for title: " + err.code);
				return;
			}

			if (typeof rows[0] != 'undefined') {
				global.TITLE = rows[0].value;
			} else {
				db.query("INSERT INTO settings (name, value) VALUES ('title', 'UpAndRunning');", function(err) {
					if (err) {
						logger.error("Unable to add title: " + err.code);
						return;
					}
					global.TITLE = "UpAndRunning";
					logger.info("Set title to default-value of \"UpAndRunning\".");
				});
			}
		});

		new admin().exists(function(status) {
			if (status === false) {
				new admin().addAdmin("admin", function(status) {
					if (status === true) {
						logger.info("Admin-User [Password: admin] created.");
					}
				});
			}
		});

		db.query("SELECT value FROM settings where name = 'interval';", function(err, rows) {
			if (err) {
				logger.error("Unable to check for interval: " + err.code);
				return;
			}

			if (typeof rows[0] != 'undefined') {
				global.INTERVAL = rows[0].value;
				logger.info("Set interval to " + global.INTERVAL + " minutes");
			} else {
				db.query("INSERT INTO settings (name, value) VALUES ('interval', 5);", function(err) {
					if (err) {
						logger.error("Unable to add check-interval: " + err.code);
						return;
					}
					global.INTERVAL = 5;
					logger.info("Set interval to default-value of 5 minutes.");
				});
			}
		});

		db.query("SELECT value FROM settings where name = 'pushbullet_key';", function(err, rows) {
			if (err) {
				logger.error("Unable to check for PushBullet-API-Key: " + err.code);
				return;
			}

			if (typeof rows[0] != 'undefined') {
				global.PBAPI = rows[0].value;
				logger.info("Set PushBullet-API-Key to \"" + global.PBAPI.slice(0, 5) + " ...\".");
			} else {
				db.query("INSERT INTO settings (name, value) VALUES ('pushbullet_key', '');", function(err) {
					if (err) {
						logger.error("Unable to add PushBullet-API-Key: " + err.code);
						return;
					}
					global.PBAPI = "";
					logger.info("Set PushBullet-API-Key to default-value of \"\".");
				});
			}
		});
	}
});

// view engine setup
app.set('views', path.join(__dirname, 'views'));
app.set('view engine', 'hjs');

app.use(session({secret: 'mySecret', resave: false, saveUninitialized: false}));
app.use(require('morgan')('dev'));
app.use(bodyParser.json());
app.use(bodyParser.urlencoded({extended: false}));
app.use(express.static(path.join(__dirname, 'public')));

// add our custom header
app.use(function(req, res, next) {
	res.setHeader("X-Powered-By", "UpAndRunning");
	next();
});

// the most important routes
app.use('/', require('./routes/index'));
app.use('/admin', require('./routes/admin'));
app.use('/api', require('./routes/api'));
app.use('/api/admin', require('./routes/api-admin'));

// catch 404 and forward to error handler
app.use(function(req, res, next) {
	var err = new Error('Not Found');
	err.status = 404;
	next(err);
});

// error handler [note for IDE: param "next" has to stay!]
app.use(function(err, req, res, next) {
	res.status(err.status || 500).render('error', {code: err.status, message: err.message});
});

// check all the websites now (after a 3 second init-delay)
setTimeout(function() {
	checkAllWebsites();
	startTimer();
}, 3 * 1000);

// checks if the "check now"-button was clicked
setInterval(function() {
	if (global.CHECK_NOW) {
		global.CHECK_NOW = false;
		checkAllWebsites();
	}
}, 1000);

// checks all websites according to the interval
function startTimer() {
	setTimeout(function() {
		checkAllWebsites();
		startTimer();
	}, global.INTERVAL * 60 * 1000);
}

// "check all websites"-function
function checkAllWebsites() {
	db.query("SELECT id, protocol, url FROM website WHERE enabled = 1;", function(err, rows) {
		if (err) {
			logger.error("Unable to search for websites in my database: " + err.code);
		} else {
			logger.info("Checking " + rows.length + " active websites...");

			for (var i in rows) {
				if (rows.hasOwnProperty(i)) {
					new website(rows[i].id, rows[i].protocol, rows[i].url).runCheck(true);
				}
			}
		}
	});
}

module.exports = app;