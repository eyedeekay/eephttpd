eephttpd - Static file server automatically forwarded to i2p
============================================================

usage:
------

eephttpd is a static http server which automatically runs on i2p with
the help of the SAM bridge. By default it will only be available from
the localhost and it's i2p tunnel. It can be masked from the localhost
using a container.

```
Usage of ./eephttpd:
  -a string
    	hostname to serve on (default "127.0.0.1")
  -c	Use an encrypted leaseset(true or false)
  -d string
    	the directory of static files to host(default ./www) (default "./www")
  -f string
    	Use an ini file for configuration (default "none")
  -g	Uze gzip(true or false) (default true)
  -i	save i2p keys(and thus destinations) across reboots (default true)
  -ib int
    	Set inbound tunnel backup quantity(0 to 5) (default 4)
  -il int
    	Set inbound tunnel length(0 to 7) (default 3)
  -iq int
    	Set inbound tunnel quantity(0 to 15) (default 8)
  -iv int
    	Set inbound tunnel length variance(-7 to 7)
  -l string
    	Type of access list to use, can be "whitelist" "blacklist" or "none". (default "none")
  -m string
    	Certificate name to use (default "cert")
  -n string
    	name to give the tunnel(default eephttpd) (default "eephttpd")
  -ob int
    	Set outbound tunnel backup quantity(0 to 5) (default 4)
  -ol int
    	Set outbound tunnel length(0 to 7) (default 3)
  -oq int
    	Set outbound tunnel quantity(0 to 15) (default 8)
  -ov int
    	Set outbound tunnel length variance(-7 to 7)
  -p string
    	port to serve locally on (default "7880")
  -r	Reduce tunnel quantity when idle(true or false)
  -rc int
    	Reduce idle tunnel quantity to X (0 to 5) (default 3)
  -rt int
    	Reduce tunnel quantity after X (milliseconds) (default 600000)
  -s string
    	the directory to save the keys in(default ./) (default ".")
  -sh string
    	sam host to connect to (default "127.0.0.1")
  -sp string
    	sam port to connect to (default "7656")
  -t	Generate or use an existing TLS certificate
  -z	Allow zero-hop, non-anonymous tunnels(true or false)
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
[http://566niximlxdzpanmn4qouucvua3k7neniwss47li5r6ugoertzuq.b32.i2p](http://566niximlxdzpanmn4qouucvua3k7neniwss47li5r6ugoertzuq.b32.i2p)
