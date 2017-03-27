# Docker Guide
[![Docker Version](https://images.microbadger.com/badges/version/marvinmenzerath/upandrunning2.svg)](http://microbadger.com/images/marvinmenzerath/upandrunning2)
[![Docker Layers](https://images.microbadger.com/badges/image/marvinmenzerath/upandrunning2.svg)](http://microbadger.com/images/marvinmenzerath/upandrunning2)

You can use the small and always up-to-date Docker-image from [Docker Hub](https://hub.docker.com/r/marvinmenzerath/upandrunning2/) to deploy UpAndRunning2 quickly and easily.

## Installation

### Database
You may skip this step if you want to use your host's database-server (remember to skip the container-link below).

Use this command to create a new MariaDB-container called `uar2-db`, which stores it's data inside the host's `/srv/uar2/db/` directory and uses the root-password `topSecretPassword`:
```
docker run -d --name uar2-db -v /srv/uar2/db/:/var/lib/mysql/ -e MYSQL_ROOT_PASSWORD='topSecretPassword' mariadb
```

### UpAndRunning2
Use this command to create and start a new UpAndRunning2-container called `uar2`, which is linked to the previously created `uar2-db`-container and exposes the web-interface and API on the host's port `80`.
Additionally it stores the configuration-files inside the host's `/srv/uar2/config/` directory, which allows you to change certain parameters.
```
docker run -d --name uar2 --link uar2-db:mysql -v /srv/uar2/config/:/app/upandrunning2/config/ -p 80:8080 marvinmenzerath/upandrunning2
```

If you do not want to use a linked database-container, just skip the `--link uar2-db:mysql` part and remember to add a correct database-configuration in your local config-file.

To make use of the mounted config-directory, you need to:
* pull a recent copy of the [`default.json`-file](../../config/default.json)
* create a copy of it called `local.json`
* and make your custom changes in this file

Additionally you need to restart your container using `docker restart uar2`.

## Upgrading
Just remove the old container, pull the new image and deploy a new container.
```
docker stop uar2
docker rm uar2
docker pull marvinmenzerath/upandrunning2
docker run -d --name uar2 --link uar2-db:mysql -p 80:8080 marvinmenzerath/upandrunning2
```

## Nice to know
* If your database-server is running on the host, you need to set the database-server's bind-address to `0.0.0.0` and may want to setup a firewall policy to protect it
    * Additionally you should make sure to allow the database-user to access the database-server from a different ip-address
* If you want to use an (Apache) proxy on the host, you need to set the exposed port in your `docker run` and make sure to bind the application on `0.0.0.0`
