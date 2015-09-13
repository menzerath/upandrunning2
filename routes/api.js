var express = require('express');
var router = express.Router();

var db = require('../lib/database');
var logger = require('../lib/logger');

router.get('/', function(req, res) {
	res.send({requestSuccess: true, message: 'Welcome to UpAndRunning\'s API!'});
});

router.get('/status/:url', function(req, res) {
	db.query("SELECT * FROM website WHERE url = ? AND enabled = 1;", [req.params.url], function(err, rows) {
		if (err) {
			logger.error("Unable to fetch website-status: " + err.code);
			res.status(500).send({requestSuccess: false, message: 'Unable to process your request.'});
		} else {
			if (rows[0] === undefined) {
				res.status(404).send({
					requestSuccess: false,
					message: 'Unable to find any data matching the given url.'
				});
			} else {
				res.send({
					requestSuccess: true,
					websiteData: {id: rows[0].id, name: rows[0].name, url: rows[0].protocol + '://' + rows[0].url},
					availability: {
						ups: rows[0].ups,
						downs: rows[0].downs,
						total: rows[0].totalChecks,
						average: rows[0].avgAvail + '%'
					},
					lastCheckResult: {status: rows[0].status, time: rows[0].time},
					lastFailedCheckResult: {status: rows[0].lastFailStatus, time: rows[0].lastFailTime}
				});
			}
		}
	});
});

router.get('/websites', function(req, res) {
	db.query("SELECT name, protocol, url, status FROM website WHERE enabled = 1 AND visible = 1;", function(err, rows) {
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
						name: rows[i].name,
						protocol: rows[i].protocol,
						url: rows[i].url,
						status: rows[i].status
					});
				}
				res.send(content);
			}
		}
	});
});

module.exports = router;