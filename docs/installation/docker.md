# Docker Guide
[![Docker Version](https://images.microbadger.com/badges/version/marvinmenzerath/upandrunning2.svg)](http://microbadger.com/images/marvinmenzerath/upandrunning2)
[![Docker Layers](https://images.microbadger.com/badges/image/marvinmenzerath/upandrunning2.svg)](http://microbadger.com/images/marvinmenzerath/upandrunning2)

You can use the small and always up-to-date Docker-image from [Docker Hub](https://hub.docker.com/r/marvinmenzerath/upandrunning2/) to deploy UpAndRunning2 quickly and easily.

## Installation

### Database
Use these commands to create a new volume called `uar2-db-data` and start a new MariaDB-container called `uar2-db`, which stores it's data inside the previously created volume and uses the root-password `topSecretPassword`:
```
docker volume create --name uar2-db-data
docker run -d --name uar2-db -v uar2-db-data:/var/lib/mysql/ -e MYSQL_ROOT_PASSWORD='topSecretPassword' mariadb
```

If you do not want to use a volume to store the databases's data, simply use this command:
```
docker run -d --name uar2-db -e MYSQL_ROOT_PASSWORD='topSecretPassword' mariadb
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

### Configurable Settings
* `UAR2_APPLICATION_TITLE` (e.g. `UpAndRunning2`)
* `UAR2_CHECKLIFETIME` (e.g. `31`)
* `UAR2_USEWEBFRONTEND` (e.g. `true`)
* `UAR2_TELEGRAMBOTAPIKEY` (e.g. `123456`)

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
