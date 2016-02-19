# Features
* Easy [Installation](installation/install.md) and Configuration
	* Use a small [Docker](installation/docker.md)-image for quick and easy deployment
* Checks all of your websites regularly
	* Use `HEAD` or `GET` requests
	* Set an interval of 10 seconds up to 10 minutes
	* Set a maximum amount of redirects to follow
	* Detects nearly all HTTP-status-codes, timeouts and unknown hosts
	* Control whether the application checks without an internet-connection or not
* Simple, but powerful [JSON-REST-API](api/index.md)
* Build your own client or use the fully-featured web-interface
* bcrypt ensures your password is stored safely
* Get notifications via
	* Email
	* Pushbullet

## Some Details on the Check-Algorithm
UpAndRunning2 checks the response it gets from a simple HTTP-HEAD-request to the specified url.  
This table shows how the different responses are handled:

| Response Code | Category |
|---------------|----------|
| 1xx           | OK       |
| 2xx           | OK       |
| 3xx           | Warning  |
| 4xx           | Error    |
| 5xx           | Error    |

Next to those HTTP status codes the application is also able to recognize a request timeout (allows a second check) or unknown hosts.

**Notice**: Some websites or applications may not respond correctly to a HEAD-request.  
In this case you need to adjust the used Check-Method to a GET-request.