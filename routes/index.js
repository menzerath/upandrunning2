var express = require('express');
var router = express.Router();

router.get('/', function(req, res) {
	res.render('index', {
		title: global.TITLE,
		partials: {styles: 'partials/styles', footer: 'partials/footer', scripts: 'partials/scripts'}
	});
});

router.get('/status', function(req, res) {
	res.render('index', {
		title: global.TITLE,
		partials: {styles: 'partials/styles', footer: 'partials/footer', scripts: 'partials/scripts'}
	});
});

router.get('/status/:url', function(req, res) {
	res.render('index', {
		title: global.TITLE,
		partials: {styles: 'partials/styles', footer: 'partials/footer', scripts: 'partials/scripts'}
	});
});

module.exports = router;