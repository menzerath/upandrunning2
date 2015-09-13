var express = require('express');
var router = express.Router();

var sanitizer = require('sanitizer');
var db = require('../lib/database');
var logger = require('../lib/logger');
var admin = require('../lib/admin');

router.get('/', function(req, res) {
	res.send({
		requestSuccess: true,
		message: 'Welcome to UpAndRunning\'s Admin-API! Please be aware that most operations need an incoming POST-request.'
	});
});

router.get('/websites', function(req, res) {
	if (!req.session.loggedin) {
		res.status(401).send({requestSuccess: false, message: 'Unauthorized'});
		return;
	}
	db.query("SELECT * FROM website;", function(err, rows) {
		if (err) {
			logger.error("Unable to fetch websites: " + err.code);
			res.status(500).send({requestSuccess: false, message: 'Unable to process your request.'});
		} else {
			if (rows[0] === undefined) {
				res.status(404).send({requestSuccess: false, message: 'Unable to find any data.'});
			} else {
				var content = {requestSuccess: true, websites: []};
				for (var i = 0; i < rows.length; i++) {
					content.websites.push({
						id: rows[i].id,
						name: rows[i].name,
						enabled: rows[i].enabled ? true : false,
						visible: rows[i].visible ? true : false,
						protocol: rows[i].protocol,
						url: rows[i].url,
						status: rows[i].status,
						time: rows[i].time,
						avgAvail: rows[i].avgAvail + '%'
					});
				}
				res.send(content);
			}
		}
	});
});

router.post('/websites/add', function(req, res) {
	if (!req.session.loggedin) {
		res.status(401).send({requestSuccess: false, message: 'Unauthorized'});
		return;
	}
	var insertData = {name: req.body.name, protocol: req.body.protocol, url: req.body.url};
	db.query("INSERT INTO website SET ?;", insertData, function(err) {
		if (err) {
			logger.error("Unable to add new website: " + err.code);
			res.status(400).send({requestSuccess: false, message: 'Unable to process your request: ' + err.code});
		} else {
			res.send({requestSuccess: true});
		}
	});
});

router.post('/websites/enable', function(req, res) {
	if (!req.session.loggedin) {
		res.status(401).send({requestSuccess: false, message: 'Unauthorized'});
		return;
	}
	db.query("UPDATE website SET enabled = 1 WHERE id = ?;", [req.body.id], function(err) {
		if (err) {
			logger.error("Unable to enable website: " + err.code);
			res.status(400).send({requestSuccess: false, message: 'Unable to process your request: ' + err.code});
		} else {
			res.send({requestSuccess: true});
		}
	});
});

router.post('/websites/disable', function(req, res) {
	if (!req.session.loggedin) {
		res.status(401).send({requestSuccess: false, message: 'Unauthorized'});
		return;
	}
	db.query("UPDATE website SET enabled = 0 WHERE id = ?;", [req.body.id], function(err) {
		if (err) {
			logger.error("Unable to disable website: " + err.code);
			res.status(400).send({requestSuccess: false, message: 'Unable to process your request: ' + err.code});
		} else {
			res.send({requestSuccess: true});
		}
	});
});

router.post('/websites/visible', function(req, res) {
	if (!req.session.loggedin) {
		res.status(401).send({requestSuccess: false, message: 'Unauthorized'});
		return;
	}
	db.query("UPDATE website SET visible = 1 WHERE id = ?;", [req.body.id], function(err) {
		if (err) {
			logger.error("Unable to enable website-visibility: " + err.code);
			res.status(400).send({requestSuccess: false, message: 'Unable to process your request: ' + err.code});
		} else {
			res.send({requestSuccess: true});
		}
	});
});

router.post('/websites/invisible', function(req, res) {
	if (!req.session.loggedin) {
		res.status(401).send({requestSuccess: false, message: 'Unauthorized'});
		return;
	}
	db.query("UPDATE website SET visible = 0 WHERE id = ?;", [req.body.id], function(err) {
		if (err) {
			logger.error("Unable to disable website-visibility: " + err.code);
			res.status(400).send({requestSuccess: false, message: 'Unable to process your request: ' + err.code});
		} else {
			res.send({requestSuccess: true});
		}
	});
});

router.post('/websites/edit', function(req, res) {
	if (!req.session.loggedin) {
		res.status(401).send({requestSuccess: false, message: 'Unauthorized'});
		return;
	}
	db.query("UPDATE website SET name = ?, protocol = ?, url = ? WHERE id = ?;", [req.body.name, req.body.protocol, req.body.url, req.body.id], function(err) {
		if (err) {
			logger.error("Unable to edit website: " + err.code);
			res.status(400).send({requestSuccess: false, message: 'Unable to process your request: ' + err.code});
		} else {
			res.send({requestSuccess: true});
		}
	});
});

router.post('/websites/delete', function(req, res) {
	if (!req.session.loggedin) {
		res.status(401).send({requestSuccess: false, message: 'Unauthorized'});
		return;
	}
	db.query("DELETE FROM website WHERE id = ?;", [req.body.id], function(err) {
		if (err) {
			logger.error("Unable to remove website: " + err.code);
			res.status(400).send({requestSuccess: false, message: 'Unable to process your request: ' + err.code});
		} else {
			res.send({requestSuccess: true});
		}
	});
});

router.post('/settings/title', function(req, res) {
	if (!req.session.loggedin) {
		res.status(401).send({requestSuccess: false, message: 'Unauthorized'});
		return;
	}
	var newTitle = sanitizer.escape(req.body.title);
	db.query("UPDATE settings SET value = ? WHERE name = 'title';", [newTitle], function(err) {
		if (err) {
			logger.error("Unable to change title: " + err.code);
			res.status(400).send({requestSuccess: false, message: 'Unable to process your request: ' + err.code});
		} else {
			res.send({requestSuccess: true});
			global.TITLE = newTitle;
		}
	});
});

router.post('/settings/password', function(req, res) {
	if (!req.session.loggedin) {
		res.status(401).send({requestSuccess: false, message: 'Unauthorized'});
		return;
	}
	new admin().changePassword(req.body.password, function(status, error) {
		if (status === false) {
			logger.error("Unable to change password: " + error);
			res.status(400).send({requestSuccess: false, message: 'Unable to process your request: ' + error});
		} else {
			res.send({requestSuccess: true});
		}
	});
});

router.post('/settings/interval', function(req, res) {
	if (!req.session.loggedin) {
		res.status(401).send({requestSuccess: false, message: 'Unauthorized'});
		return;
	}
	var newInterval = req.body.interval;
	if (isNaN(newInterval) || newInterval < 1 || newInterval > 60) {
		res.status(400).send({
			requestSuccess: false,
			message: 'Unable to process your request: Interval has to be a number ranging between 1 and 60 minutes.'
		});
		return;
	}
	db.query("UPDATE settings SET value = ? WHERE name = 'interval';", [newInterval], function(err) {
		if (err) {
			logger.error("Unable to change interval: " + err.code);
			res.status(400).send({requestSuccess: false, message: 'Unable to process your request: ' + err.code});
		} else {
			res.send({requestSuccess: true});
			global.INTERVAL = newInterval;
		}
	});
});

router.post('/settings/pbkey', function(req, res) {
	if (!req.session.loggedin) {
		res.status(401).send({requestSuccess: false, message: 'Unauthorized'});
		return;
	}
	var newKey = sanitizer.escape(req.body.key);
	db.query("UPDATE settings SET value = ? WHERE name = 'pushbullet_key';", [newKey], function(err) {
		if (err) {
			logger.error("Unable to change PushBullet-API-Key: " + err.code);
			res.status(400).send({requestSuccess: false, message: 'Unable to process your request: ' + err.code});
		} else {
			res.send({requestSuccess: true});
			global.PBAPI = newKey;
		}
	});
});

router.post('/check', function(req, res) {
	if (!req.session.loggedin) {
		res.status(401).send({requestSuccess: false, message: 'Unauthorized'});
		return;
	}
	global.CHECK_NOW = true;
	res.send({requestSuccess: true});
});

router.post('/login', function(req, res) {
	new admin().validatePassword(req.body.password, function(status, error) {
		if (status === false) {
			logger.error("Unable to login: " + error);
			res.status(400).send({requestSuccess: false, message: 'Unable to process your request: ' + error});
		} else {
			req.session.loggedin = true;
			res.send({requestSuccess: true});
		}
	});
});

router.post('/logout', function(req, res) {
	req.session.destroy(function(err) {
		if (err) {
			logger.error("Unable to logout: " + err.code);
			res.status(400).send({requestSuccess: false, message: 'Unable to process your request: ' + err.code});
		} else {
			res.send({requestSuccess: true});
		}
	});
});

module.exports = router;