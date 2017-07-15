FROM alpine:3.6
MAINTAINER Marvin Menzerath <github@marvin-menzerath.de>

WORKDIR /app/upandrunning2/
COPY . /app/upandrunning2/
RUN export GOPATH=/tmp/go && \
    export PATH=${PATH}:${GOPATH}/bin && \
    apk -U --no-progress add git go build-base ca-certificates && \
    mkdir -p ${GOPATH}/src/github.com/MarvinMenzerath/ && \
    ln -s /app/upandrunning2/ ${GOPATH}/src/github.com/MarvinMenzerath/UpAndRunning2 && \
    cd ${GOPATH}/src/github.com/MarvinMenzerath/UpAndRunning2 && \
    go get && \
    go build && \
    rm -rf $GOPATH && \
    apk --no-progress del git go build-base && \
    addgroup -g 1777 uar2 && adduser -h /app/upandrunning2/ -H -D -G uar2 -u 1777 uar2 && \
    chown -R uar2:uar2 /app/upandrunning2/

USER uar2
EXPOSE 8080
VOLUME /app/upandrunning2/config/
ENTRYPOINT ["./UpAndRunning2"]