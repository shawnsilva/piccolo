FROM alpine:3.4
MAINTAINER shawnsilva

ENV APP_USER=piccolo \
    APP_NAME=piccolo \
    APP_CREATOR=shawnsilva \
    GOROOT=/usr/lib/go \
    GOPATH=/opt/golang \
    GOBIN=/opt/golang/bin
ENV PATH=${PATH}:${GOROOT}/bin:${GOPATH}/bin

RUN apk add --update ffmpeg opus bash ca-certificates && \
    # Setup a base user for running applications
    adduser -u 1337 -D -h /home/${APP_USER} ${APP_USER} && \
    update-ca-certificates && \
    rm -rf /usr/share/man /tmp/* /var/tmp/* /var/cache/apk/*

RUN apk add --update curl python gnupg && \
    gpg --keyserver pool.sks-keyservers.net --recv-keys \
    7D33D762FD6C35130481347FDB4B54CBA4826A18 \
    ED7F5BF46B3BBED81C87368E2C393E0F18A9236D && \

    curl -fSL https://yt-dl.org/downloads/2017.02.07/youtube-dl -o youtube-dl && \
    curl -fSL https://yt-dl.org/downloads/2017.02.07/youtube-dl.sig -o youtube-dl.sig && \
    gpg --verify youtube-dl.sig && \

    chmod +x youtube-dl && \
    mv youtube-dl /bin/. && rm -f youtube-dl.sig && \

    apk del curl gnupg && \
    rm -rf /usr/share/man /tmp/* /var/tmp/* /var/cache/apk/*

COPY . ${GOPATH}/src/github.com/${APP_CREATOR}/${APP_NAME}

RUN apk add --update opus-dev git make pkgconfig build-base && \
    apk add go --update-cache --repository http://dl-cdn.alpinelinux.org/alpine/edge/community/ --allow-untrusted && \

    cd ${GOPATH}/src/github.com/${APP_CREATOR}/${APP_NAME} && \
    make build && \
    mkdir -p /opt/${APP_NAME}/conf && mkdir -p /opt/${APP_NAME}/video_cache && \
    mv build/${APP_NAME} /opt/${APP_NAME}/. && \
    chown -R ${APP_USER}:${APP_USER} /opt/${APP_NAME} && \

    apk del go opus-dev git make pkgconfig build-base && \
    rm -rf /usr/share/man /tmp/* /var/tmp/* /var/cache/apk/* ${GOPATH}

USER ${APP_USER}
VOLUME /opt/${APP_NAME}/conf
VOLUME /opt/${APP_NAME}/video_cache
WORKDIR /opt/${APP_NAME}
ENTRYPOINT ["./piccolo"]
CMD ["--config","conf/config.json"]
