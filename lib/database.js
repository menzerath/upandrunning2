var config = require('config');
var mysql = require('mysql');
var logger = require('../lib/logger');

var pool = mysql.createPool(process.env.CLEARDB_DATABASE_URL || config.get('database'));

pool.on('enqueue', function() {
	logger.warn('Request has to wait for an available connection slot. You should consider adding more connections to the pool.');
});

module.exports = pool;