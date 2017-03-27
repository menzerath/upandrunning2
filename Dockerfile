FROM alpine:3.5
MAINTAINER Marvin Menzerath <github@marvin-menzerath.de>

WORKDIR /app/upandrunning2/
COPY . /app/upandrunning2/
RUN chmod +x ./docker/build.sh && sync && ./docker/build.sh

EXPOSE 8080
ENTRYPOINT ["./UpAndRunning2"]
