var db = require('../lib/database');
var logger = require('../lib/logger');
var httpcodes = require('../lib/http-status-codes');
var pushBullet = require('pushbullet');
var request = require('request');

function Website(id, protocol, url) {
	this.id = id;
	this.protocol = protocol;
	this.url = url;

	Website.prototype.runCheck = function(allowSecondTry) {
		request({
			url: protocol + '://' + url,
			headers: {'User-Agent': 'UpAndRunning (https://github.com/MarvinMenzerath/UpAndRunning)'}
		}, function(error, response) {
			var status;
			if (response === undefined) {
				status = 'Server not found';
			} else {
				status = response.statusCode + ' - ' + httpcodes[response.statusCode];
			}

			// get website's name and previous status
			var name;
			var oldStatus;
			var allowPush = true;
			if (global.PBAPI != "") {
				db.query("SELECT name, status FROM website WHERE id = ?", [id], function(err, rows) {
					if (err) {
						logger.error("Unable get website's name and old status: " + err.code);
						allowPush = false;
					} else {
						name = rows[0].name;
						oldStatus = rows[0].status;
					}
				});
			}

			if (status.indexOf(200) > -1 || status.indexOf(301) > -1 || status.indexOf(302) > -1) {
				// success
				db.query("UPDATE website SET status = ?, time = NOW(), ups = ups + 1, totalChecks = totalChecks + 1 WHERE id = ?;", [status, id], function(err) {
					if (err) {
						logger.error("Unable to save new website-status: " + err.code);
					} else {
						calcAvgAvailability();
					}
				});
			} else {
				// failure
				if (allowSecondTry) {
					setTimeout(function() {
						new Website(id, protocol, url).runCheck(false);
					}, 1000);
					allowPush = false;
				} else {
					db.query("UPDATE website SET status = ?, time = NOW(), lastFailStatus = ?, lastFailTime = NOW(), downs = downs + 1, totalChecks = totalChecks + 1 WHERE id = ?;", [status, status, id], function(err) {
						if (err) {
							logger.error("Unable to save new website-status: " + err.code);
						} else {
							calcAvgAvailability();
						}
					});
				}
			}

			// compare previous and new status and send a push if there is a change
			if (global.PBAPI != "" && allowPush && name != undefined && oldStatus != undefined) {
				setTimeout(function() {
					if (oldStatus != status) {
						var pusher = new pushBullet(global.PBAPI);
						pusher.note({}, global.TITLE + " - Status Change", name + " (" + url + ") went from \"" + oldStatus + "\" to \"" + status + "\".", function(err) {
							if (err) {
								logger.error("Unable to send PushBullet-notification: " + err);
							}
						});
					}
				}, 2000); // wait until you have all the data (dirty fix)
			}
		});

		function calcAvgAvailability() {
			db.query("SELECT ((SELECT ups FROM website WHERE id = ?) / (SELECT totalChecks FROM website WHERE id = ?))*100 AS avg", [id, id], function(err, rows) {
				if (err) {
					logger.error("Unable to calculate new website-availability: " + err.code);
					return;
				}
				var avgAvail = rows[0].avg.toFixed(2);

				db.query("UPDATE website SET avgAvail = ? WHERE id = ?;", [avgAvail, id], function(err) {
					if (err) {
						logger.error("Unable to save new website-availability: " + err.code);
					}
				});
			});
		}
	};
}

module.exports = Website;