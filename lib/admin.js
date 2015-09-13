var crypto = require('crypto');
var db = require('../lib/database');
var logger = require('../lib/logger');

function Admin() {

	Admin.prototype.validatePassword = function(enteredPassword, callback) {
		db.query("SELECT value FROM settings WHERE name = 'salt';", function(err, rows) {
			if (err) {
				logger.error("Unable to get password-salt: " + err.code);
				callback(false, err.code);
			}
			var salt = rows[0].value;

			db.query("SELECT value FROM settings WHERE name = 'password';", function(err, rows) {
				if (err) {
					logger.error("Unable to get password-hash: " + err.code);
					callback(false, err.code);
				}
				var password = rows[0].value;

				var calculatedPassword = crypto.pbkdf2Sync(enteredPassword, salt, 4096, 512, 'sha256').toString('hex');
				if (calculatedPassword === password) {
					logger.info("Login successful.");
					callback(true);
				} else {
					logger.info("Login failed.");
					callback(false, "Wrong Password.");
				}
			});
		});
	};

	Admin.prototype.exists = function(callback) {
		db.query("SELECT COUNT(name) as total FROM settings WHERE name = 'password';", function(err, rows) {
			if (err) {
				logger.error("Unable to check for admin-password in database: " + err.code);
				callback(false);
			} else {
				if (rows[0].total === 1) {
					db.query("SELECT COUNT(name) as total FROM settings WHERE name = 'salt';", function(err, rows) {
						if (err) {
							logger.error("Unable to check for admin-salt in database: " + err.code);
							callback(false);
						} else {
							if (rows[0].total === 1) {
								callback(true);
							} else {
								callback(false);
							}
						}
					});
				} else {
					callback(false);
				}
			}
		});
	};

	Admin.prototype.changePassword = function(newPassword, callback) {
		var salt = crypto.randomBytes(256).toString('hex');
		var password = crypto.pbkdf2Sync(newPassword, salt, 4096, 512, 'sha256').toString('hex');

		db.query("UPDATE settings SET value = ? WHERE name = 'salt';", [salt], function(err) {
			if (err) {
				logger.error("Unable to change password-salt: " + err.code);
				callback(false, err.code);
			} else {
				db.query("UPDATE settings SET value = ? WHERE name = 'password';", [password], function(err) {
					if (err) {
						logger.error("Unable to change password: " + err.code);
						callback(false, err.code);
					} else {
						callback(true);
					}
				});
			}
		});
	};

	Admin.prototype.addAdmin = function(newPassword, callback) {
		var salt = crypto.randomBytes(256).toString('hex');
		var password = crypto.pbkdf2Sync(newPassword, salt, 4096, 512, 'sha256').toString('hex');

		var insertData = {name: "salt", value: salt};
		db.query("INSERT INTO settings SET ?;", insertData, function(err) {
			if (err) {
				logger.error("Unable to add password-salt: " + err.code);
				callback(false, err.code);
			} else {
				var insertData = {name: "password", value: password};
				db.query("INSERT INTO settings SET ?;", insertData, function(err) {
					if (err) {
						logger.error("Unable to add password: " + err.code);
						callback(false, err.code);
					} else {
						callback(true);
					}
				});
			}
		});
	};
}

module.exports = Admin;