# Docker Guide
[![Docker Version](https://images.microbadger.com/badges/version/marvinmenzerath/upandrunning2.svg)](http://microbadger.com/images/marvinmenzerath/upandrunning2)
[![Docker Layers](https://images.microbadger.com/badges/image/marvinmenzerath/upandrunning2.svg)](http://microbadger.com/images/marvinmenzerath/upandrunning2)

You can use the small and always up-to-date Docker-image from [Docker Hub](https://hub.docker.com/r/marvinmenzerath/upandrunning2/) to deploy UpAndRunning2 quickly and easily.

## Installation

### Database
Use this command to create a new MariaDB-container called `uar2-db`, which stores it's data inside the host's `/srv/uar2/db/` directory and uses the root-password `topSecretPassword`:
```
docker run -d --name uar2-db -v /srv/uar2/db/:/var/lib/mysql/ -e MYSQL_ROOT_PASSWORD='topSecretPassword' mariadb
```

### UpAndRunning2
Use this command to create and start a new UpAndRunning2-container called `uar2`, which is linked to the previously created `uar2-db`-container and exposes the web-interface and API on the host's port `80`.
Additionally it stores the configuration-files inside the host's `/srv/uar2/config/` directory and allows you to change certain parameters.
```
docker run -d --name uar2 --link uar2-db:mysql -v /srv/uar2/config/:/app/upandrunning2/config/ -p 80:8080 marvinmenzerath/upandrunning2
```

If you do not want to use a linked database-container, just skip the `--link uar2-db:mysql` part and remember to add a correct database-configuration in your local config-file.

## Upgrading
Just remove the old container, pull the new image and deploy a new container.
```
docker stop uar2
docker rm uar2
docker pull marvinmenzerath/upandrunning2
docker run -d --name uar2 --link uar2-db:mysql -p 80:8080 marvinmenzerath/upandrunning2
```