FROM alpine:latest
ENV samhost=localhost samport=7656 args="-r"
RUN apk update -U
RUN apk add go git make musl-dev
RUN mkdir -p /opt/eephttpd
RUN adduser -h /opt/eephttpd -D -g "eephttpd,,,," eephttpd
COPY . /usr/src/eephttpd
WORKDIR /usr/src/eephttpd
RUN go mod vendor && make build && install -m755 eephttpd /usr/bin/eephttpd
WORKDIR /opt/eephttpd
USER eephttpd
COPY www /opt/eephttpd/www
VOLUME /opt/eephttpd/
CMD eephttpd -f /usr/src/eephttpd/etc/eephttpd/eephttpd.conf \
    -i \
    -d /opt/eephttpd/www \
    -sh=$samhost \
    -sp=$samport $args
