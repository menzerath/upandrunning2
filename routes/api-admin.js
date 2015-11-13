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
