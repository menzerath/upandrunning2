# UpAndRunning2
[![Build Status](https://drone.io/github.com/MarvinMenzerath/UpAndRunning2/status.png)](https://drone.io/github.com/MarvinMenzerath/UpAndRunning2/latest)
[![Docker Layers](https://badge.imagelayers.io/marvinmenzerath/upandrunning2:latest.svg)](https://imagelayers.io/?images=marvinmenzerath/upandrunning2:latest)

UpAndRunning2 is a lightweight Go application which **monitors all of your websites**, offers a simple **JSON-REST-API** and user-defined **notifications**.

## Features
* Easy [Installation](docs/installation/install.md) and Configuration
	* Use a small [Docker](docs/installation/docker.md)-image for quick and easy deployment
* Checks all of your websites regularly
	* Use `HEAD` or `GET` requests
	* Set an interval of 10 seconds up to 10 minutes
	* Set a maximum amount of redirects to follow
	* Detects nearly all HTTP-status-codes, timeouts and unknown hosts
	* Control whether the application checks without an internet-connection or not
* Simple, but powerful [JSON-REST-API](docs/api/index.md)
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

## Documentation
* [Installation](docs/installation/install.md)
	* [Upgrading](docs/installation/upgrade.md)
	* [Docker Guide](docs/installation/docker.md)
* [API](docs/api/index.md)
	* [v1](docs/api/v1.md)
* [Screenshots](docs/screenshots/index.md)

## Credits

### Used Software
* [Bootstrap](https://github.com/twbs/bootstrap)
	* [Bootswatch Theme: Paper](https://github.com/thomaspark/bootswatch)
* [Font Awesome](http://fontawesome.io)
* [jQuery](https://jquery.com)
	* [Chart.js](https://github.com/nnnick/Chart.js)
	* [SweetAlert](https://github.com/t4t5/sweetalert)
* Golang Libraries
	* [Go-MySQL-Driver](https://github.com/go-sql-driver/mysql)
	* [Golang logging library](https://github.com/op/go-logging)
	* [gomail](http://gopkg.in/gomail.v2)
	* [GoReq](https://github.com/franela/goreq)
	* [HttpRouter](https://github.com/julienschmidt/httprouter)
	* [pushbullet-go](https://github.com/mitsuse/pushbullet-go)

### Application Icon
[Icon](https://www.iconfinder.com/icons/328014/back_on_top_top_up_upload_icon) created by [Aha-Soft Team](http://www.aha-soft.com) - [CC BY 2.5 License](http://creativecommons.org/licenses/by/2.5/)

## License
The MIT License (MIT)

Copyright (c) 2015-2016 Marvin Menzerath

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
