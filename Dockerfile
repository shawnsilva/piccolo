FROM golang:1.10.3-alpine-3.8 as builder
WORKDIR /go/src/github.com/shawnsilva/piccolo/
COPY . .
RUN apk add --update --no-cache opus-dev git make pkgconfig build-base && \
    make deps build

FROM alpine:3.8
ENV APP_USER=piccolo \
    APP_NAME=piccolo

RUN apk add --no-cache --update ffmpeg opus bash ca-certificates && \
    # Setup a base user for running applications
    adduser -u 1337 -D -h /home/${APP_USER} ${APP_USER} && \
    update-ca-certificates && \
    mkdir -p /opt/${APP_NAME}/conf && mkdir -p /opt/${APP_NAME}/video_cache && \
    chown -R ${APP_USER}:${APP_USER} /opt/${APP_NAME} && \
    rm -rf /usr/share/man /tmp/* /var/tmp/* /var/cache/apk/*

COPY --from=builder /go/src/github.com/shawnsilva/piccolo/build/piccolo /opt/${APP_NAME}/.

USER ${APP_USER}
VOLUME /opt/${APP_NAME}/conf
VOLUME /opt/${APP_NAME}/video_cache
WORKDIR /opt/${APP_NAME}
ENTRYPOINT ["./piccolo"]
CMD ["--config","conf/config.json"]
