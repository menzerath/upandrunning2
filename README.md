# UpAndRunning2 [![Build Status](https://drone.io/github.com/MarvinMenzerath/UpAndRunning2/status.png)](https://drone.io/github.com/MarvinMenzerath/UpAndRunning2/latest)
UpAndRunning2 is a lightweight Go application which **monitors all of your websites**, offers a simple **JSON-REST-API** and user-defined **notifications**.

## Features
* Easy [Installation](#installation) and Configuration
* Checks all of your websites regularly
	* Use `HEAD` or `GET` requests
	* Set an interval of 10 seconds up to 10 minutes
	* Set a maximum amount of redirects to follow
	* Detects nearly all HTTP-status-codes, timeouts and unknown hosts
* Simple, but powerful [JSON-REST-API](#api)
* Build your own client or use the fully-featured web-interface
* bcrypt ensures your password is stored safely
* Get notifications via
	* Email
	* Pushbullet

### Some Details on the Check-Algorithm
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

## Installation
* Download and extract all the files in a directory
* Prepare your MySQL-Server: create a new user and a new database
* Copy `config/default.json` to `config/local.json` and change this file to your needs
* Visit `http://localhost:8080/admin` and use `admin` to authenticate. You should change the password immediately.
* Done!

### Upgrading
When upgrading the application, you need to update to the latest release of your currently installed minor version (e.g. v2.0.2) before upgrading to the next minor or major version (e.g. v2.1).  
Otherwise you may loose access to your data and need to wipe your database.

Also make sure that your custom applications using UpAndRunning's API are up-to-date and are not using a deprecated API-version.

## API

### User
Notice: Everyone is able to access those APIs.

#### Status
`GET` `/api/v1/websites/:url/status`:

```json
{
	"requestSuccess": true,
	"websiteData": {
		"id": 1,
		"name": "My Website",
		"url": "https://website.com"
	},
	"availability": {
		"ups": 99,
		"downs": 1,
		"total": 100,
		"average": "99.00%"
	},
	"lastCheckResult": {
		"status": "200 - OK",
		"responseTime": "150 ms",
		"time": "2015-01-01 00:00:00"
	},
	"lastFailedCheckResult": {
		"status": "500 - Internal Server Error",
		"responseTime": "0 ms",
		"time": "2014-12-31 20:15:00"
	}
}
```

#### Results
`GET` `/api/v1/websites/:url/results`:

Optional parameter: `?limit=100`  
Optional parameter: `?offset=50`

```json
{
	"requestSuccess": true,
	"results": [
		{
			"status": "200 - OK",
            "responseTime": "150 ms",
            "time": "2015-01-01 00:00:00"
		}
	]
}
```

#### List
`GET` `/api/v1/websites`:

```json
{
	"requestSuccess": true,
	"websites": [
		{
			"name": "My Website",
			"protocol": "https",
			"url": "website.com",
			"status": "200 - OK"
		}
	]
}
```

### Admin
Notice: You have to login before you are able to use those APIs.

#### List all Websites
`GET` `/api/v1/admin/websites`:

```json
{
	"requestSuccess": true,
	"websites": [
		{
			"id": 1,
			"name": "My Website",
			"enabled": true,
			"visible": true,
			"protocol": "https",
			"url": "website.com",
			"checkMethod": "HEAD",
			"status": "200 - OK",
			"time": "2015-01-01 00:00:00"
		}
	]
}
```

#### Add a Website
`POST` `/api/v1/admin/websites/:url`:

```
POST-parameters: name, protocol, checkMethod
```

#### Edit a Website
`PUT` `/api/v1/admin/websites/:url`:

```
PUT-parameters: name, protocol, url, checkMethod
```

#### Delete a Website
`DELETE` `/api/v1/admin/websites/:url`:

#### Enable / Disable a Website
`PUT` `/api/v1/admin/websites/:url/enabled`:

```
PUT-parameters: enabled {true / false}
```

#### Set a Website's visibility
`PUT` `/api/v1/admin/websites/:url/visibility`:

```
PUT-parameters: visible {true / false}
```

#### Get a Website's notification settings
`GET` `/api/v1/admin/websites/:url/notifications`:

```
{
	"requestSuccess": true,
	"notifications": {
		"pushbulletKey": "abcdef123456",
		"email": "me@mymail.com"
	}
}
```

#### Set a Website's notification settings
`PUT` `/api/v1/admin/websites/:url/notifications`:

```
PUT-parameters: pushbulletKey, email {"" to disable}
```

#### Change Application-Title
`PUT` `/api/v1/settings/title`:

```
PUT-parameters: title
```

#### Change Check-Interval
`PUT` `/api/v1/settings/interval`:

```
PUT-parameters: interval
```

#### Change Admin-Password
`PUT` `/api/v1/settings/password`:

```
PUT-parameters: password
```

#### Change amount of Redirects
`PUT` `/api/v1/settings/redirects`:

```
PUT-parameters: redirects
```

#### Trigger a Check
`GET` `/api/v1/check`:

#### Login
`POST` `/api/v1/auth/login`:

```
POST-parameters: password
```

#### Logout
`GET` `/api/v1/auth/logout`:

## Screenshots
![User-Interface](doc/Screenshot1.jpg)
![API](doc/Screenshot2.jpg)
![Admin-Backend](doc/Screenshot3.jpg)

## Credits

### Application Icon
[Icon](https://www.iconfinder.com/icons/328014/back_on_top_top_up_upload_icon) created by [Aha-Soft Team](http://www.aha-soft.com) - [CC BY 2.5 License](http://creativecommons.org/licenses/by/2.5/)

## License
The MIT License (MIT)

Copyright (c) 2015 Marvin Menzerath

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
