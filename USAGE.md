eephttpd - Static file server automatically forwarded to i2p
============================================================

usage:
------

eephttpd is a static http server which automatically runs on i2p with
the help of the SAM bridge. By default it will only be available from
the localhost and it's i2p tunnel. It can be masked from the localhost
using a container.

```
/bin/sh: 1: ./eephttpd: Permission denied
```

### build in docker

```
docker build --build-arg user=eephttpd \
    --build-arg path=example/www \
    -f Dockerfile -t \
    eyedeekay/eephttpd .
```

### Run in docker

```
docker run -i -t -d \
    --name eephttpd-volume \
    --volume eephttpd:/opt/eephttpd/ \
    eyedeekay/eephttpd
```

```
docker run -i -t -d \
    --network si \
    --env samhost=sam-host \
    --env samport=7656 \
    --env args=-r # Additional arguments to pass to eephttpd\
    --network-alias eephttpd \
    --hostname eephttpd \
    --name eephttpd \
    --restart always \
    --volumes-from eephttpd-volume \
    eyedeekay/eephttpd
```

instance
--------

a running instance of eephttpd with the example index file is availble on
[http://tvndxxkxcstbtqfxg7iigco6bj22ff2y6jxikmk7wqkyadkhrd4a.b32.i2p](http://tvndxxkxcstbtqfxg7iigco6bj22ff2y6jxikmk7wqkyadkhrd4a.b32.i2p)

