FROM alpine:3.3
MAINTAINER Marvin Menzerath <github@marvin-menzerath.de>

ENV UAR2_IS_DOCKER true

WORKDIR /app/upandrunning2/
COPY . /app/upandrunning2/
RUN chmod +x ./docker/build.sh && ./docker/build.sh

EXPOSE 8080
ENTRYPOINT ["./UpAndRunning2"]
