# Docker Guide
[![Docker Layers](https://images.microbadger.com/badges/image/marvinmenzerath/upandrunning2.svg)](http://microbadger.com/images/marvinmenzerath/upandrunning2)

You can use the small and always up-to-date Docker-image from [Docker Hub](https://hub.docker.com/r/marvinmenzerath/upandrunning2/) to deploy UpAndRunning2 quickly and easily.
Internally this image uses a custom user and group (both called `uar2` with an id of `1777`).

## Installation

### Database
Setup a MySQL-compatible database (for example by using the [MariaDB image](https://hub.docker.com/_/mariadb/)) and make a note of the database credentials.

### UpAndRunning2
Use this command to create and start a new UpAndRunning2-container called `uar2`, which exposes the web-interface and API on the host's port `80`.
Additionally it stores the configuration-files inside the host's `/srv/uar2/config/` directory, which allows you to change certain parameters.
```
docker run -d --name uar2 -v /srv/uar2/config/:/app/upandrunning2/config/ -p 80:8080 marvinmenzerath/upandrunning2
```

To make use of the mounted config-directory, you need to:
* pull a recent copy of the [`default.json`-file](../../config/default.json)
* create a copy of it called `local.json`
* and make your custom changes in this file

After those changes you need to restart your container using `docker restart uar2`.

## Upgrading
Just remove the old container, pull the new image and deploy a new container.
```
docker stop uar2
docker rm uar2
docker pull marvinmenzerath/upandrunning2
docker run -d --name uar2 -v /srv/uar2/config/:/app/upandrunning2/config/ -p 80:8080 marvinmenzerath/upandrunning2
```