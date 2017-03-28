FROM alpine:3.5
MAINTAINER Marvin Menzerath <github@marvin-menzerath.de>

WORKDIR /app/upandrunning2/
COPY . /app/upandrunning2/
RUN chmod +x ./docker/build.sh && sync && ./docker/build.sh

USER uar2
EXPOSE 8080
VOLUME /app/upandrunning2/config/
ENTRYPOINT ["./UpAndRunning2"]