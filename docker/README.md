# Docker Guide
[![Docker Layers](https://badge.imagelayers.io/marvinmenzerath/upandrunning2:latest.svg)](https://imagelayers.io/?images=marvinmenzerath/upandrunning2:latest)

You can use the small and always up-to-date Docker-image from [Docker Hub](https://hub.docker.com/r/marvinmenzerath/upandrunning2/) to deploy UpAndRunning2 quickly and easily.

## Installation

### Database
Use this command to create and start a new MariaDB-container called `uar2-db`, which stores it's data inside `/data/uar2-db/` and uses the root-password `topSecretPassword`:
```
docker run -d --name uar2-db -v /data/uar2-db/:/var/lib/mysql/ -e MYSQL_ROOT_PASSWORD='topSecretPassword' mariadb
```

### UpAndRunning2
Use this command to create and start a new UpAndRunning2-container called `uar2`, which is linked to the previously created `uar2-db`-container and exposes the web-interface and API on the host's port `80`.
```
docker run -d --name uar2 --link uar2-db:mysql -p 80:8080 marvinmenzerath/upandrunning2
```

## Upgrading
Just remove the old container, pull the new image and deploy a new container. Make sure to add previously set environment-variables.
```
docker stop uar2
docker rm uar2
docker pull marvinmenzerath/upandrunning2
docker run -d --name uar2 --link uar2-db:mysql -p 80:8080 marvinmenzerath/upandrunning2
```

## Configuration
There are a few things you can configure using environment-variables.  
To do so, just add those environment-variables when creating the container like this:
```
docker run -d --name uar2 --link uar2-db:mysql -p 80:8080 -e UAR2_VARIABLE_NAME='CONTENT' marvinmenzerath/upandrunning2
```

### Mailer
If you want to use the built-in mailer, you will need to set those environment-variables:
* `UAR2_MAILER_HOST` (e.g. `smtp.mymail.com`)
* `UAR2_MAILER_PORT` (e.g. `587`)
* `UAR2_MAILER_USER` (e.g. `myUser@mymail.com`)
* `UAR2_MAILER_PASSWORD` (e.g. `mySecretPassword`)
* `UAR2_MAILER_FROM` (e.g. `upandrunning2@mymail.com`)

#### Example
```
docker run -d --name uar2 --link uar2-db:mysql -p 80:8080 -e UAR2_MAILER_HOST='smtp.mymail.com' -e UAR2_MAILER_PORT='587' -e UAR2_MAILER_USER='myUser@mymail.com' -e UAR2_MAILER_PASSWORD='mySecretPassword' -e UAR2_MAILER_FROM='upandrunning2@mymail.com' marvinmenzerath/upandrunning2
```