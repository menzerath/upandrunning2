# Installation
Looking for the Docker Guide? Click [here](docker.md).

* Download and extract all the files in a directory.
* Prepare your MySQL-Server: create a new user and a new database.
* Copy `config/default.json` to `config/local.json` and change this file to your needs.
* Run the application using a systemd-script or just type `./UpAndRunning2`.
* Visit `http://localhost:8080/admin` and use `admin` to authenticate. You should change the password immediately.
* Done!

## systemd
Create a new service-file (e.g. at `/etc/systemd/system/upandrunning2.service`) and insert the following data:

```bash
[Unit]
Description=UpAndRunning2
After=syslog.target network.target mysql.service

[Service]
Type=simple
User=uar2
Group=uar2
WorkingDirectory=/home/uar2/UpAndRunning2
ExecStart=/home/uar2/UpAndRunning2/UpAndRunning2
Restart=always
Environment=USER=uar2 HOME=/home/uar2

[Install]
WantedBy=multi-user.target
```

Adjust all the values to your specific installation and next enable and start the service:
```bash
systemctl enable upandrunning2.service
systemctl start upandrunning2.service
```