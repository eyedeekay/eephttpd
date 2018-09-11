FROM alpine:3.8
ENV samhost=sam-host
ENV samport=7656
ENV args="-r"
RUN apk update -U
RUN apk add go git make musl-dev
RUN mkdir -p /opt/eephttpd
RUN adduser -h /opt/eephttpd -D -g "eephttpd,,,," eephttpd
COPY . /usr/src/eephttpd
WORKDIR /usr/src/eephttpd
RUN make release install
USER eephttpd
CMD eephttpd -f /usr/src/eephttpd/etc/eephttpd/eephttpd.conf \
    -i \
    -d /opt/eephttpd/www \
    -s /opt/eephttpd/ \
    -sh=$samhost \
    -sp=$samport $args
