package lib

var httpCodes map[int]string

// Init the httpCodes-map.
func InitHttpStatusCodeMap() {
	httpCodes = map[int]string{
		100: "Continue",
		101: "Switching Protocols",
		200: "OK",
		201: "Created",
		202: "Accepted",
		203: "Non-Authoritative Information",
		204: "No Content",
		205: "Reset Content",
		206: "Partial Content",
		300: "Multiple Choices",
		301: "Moved Permanently",
		302: "Found",
		303: "See Other",
		304: "Not Modified",
		305: "Use Proxy",
		307: "Temporary Redirect",
		400: "Bad Request",
		401: "Unauthorized",
		402: "Payment Required",
		403: "Forbidden",
		404: "Not Found",
		405: "Method Not Allowed",
		406: "Not Acceptable",
		407: "Proxy Authentication Required",
		408: "Request Time-out",
		409: "Conflict",
		410: "Gone",
		411: "Length Required",
		412: "Precondition Failed",
		413: "Request Entity Too Large",
		414: "Request-URI Too Large",
		415: "Unsupported Media Type",
		416: "Requested Range not Satisfiable",
		417: "Expectation Failed",
		422: "Unprocessable Entity",
		429: "Too Many Requests",
		500: "Internal Server Error",
		501: "Not Implemented",
		502: "Bad Gateway",
		503: "Service Unavailable",
		504: "Gateway Time-out",
		505: "HTTP Version not Supported",
		520: "Unknown error (CloudFlare)",
		521: "Connection refused (CloudFlare)",
		522: "Connection timed out (CloudFlare)",
		523: "Origin is unreachable (CloudFlare)",
		524: "A timeout occurred (CloudFlare)",
		525: "SSL handshake failed (CloudFlare)",
		526: "Invalid SSL certificate",
	}
}

// Returns a HTTP-status-code-string representing the given HTTP-status-code.
func GetHttpStatus(code int) string {
	return httpCodes[code]
}
