eephttpd - Static file server automatically forwarded to i2p
============================================================

usage:
------

eephttpd requires the use of an I2P router with an enabled SAM API.
On the default Java I2P router this is enabled by going to 
[http://localhost:7657/configclients](http://localhost:7657/configclients).

eephttpd recommends the Java I2P router from [geti2p.net](https://geti2p.net)
as the I2P router to use, as that is the router I use to test it. Other options
are [I2P-Zero](https://github.com/i2p-zero/i2p-zero), a minimized distribution
of the Java I2P router with all required dependencies bundled-in, or
[i2pd](https://i2pd.website), a C++ implementation of the I2P network which is
preferred by some power users and sysadmins for it's performance and minimal,
lightweight interface.

If you are using a Linux distribution it is recommended that you use your
package manager to install or uninstall I2P. Java I2P maintains a Debian-style
repository and a PPA for Ubuntu users and up-to-date packages are available in
Debian Sid.

eephttpd is a static http server which automatically runs on i2p with
the help of the SAM bridge. By default it will only be available from
the localhost and it's i2p tunnel. It can be masked from the localhost
using a container.
```
Usage of ./eephttpd/eephttpd:
  -a string
    	hostname to serve on (default "127.0.0.1")
  -b string
    	URL of a git repository to build populate the static directory with(optional)
  -c	Use an encrypted leaseset(true or false)
  -d string
    	the directory of static files to host(default./www) (default "./www")
  -f string
    	Use an ini file for configuration (default "none")
  -g	Uze gzip(true or false) (default true)
  -i	save i2p keys(and thus destinations) across reboots (default true)
  -ib int
    	Set inbound tunnel backup quantity(0 to 5) (default 1)
  -il int
    	Set inbound tunnel length(0 to 7) (default 3)
  -iq int
    	Set inbound tunnel quantity(0 to 15) (default 2)
  -iv int
    	Set inbound tunnel length variance(-7 to 7)
  -l string
    	Type of access list to use, can be "whitelist" "blacklist" or "none". (default "none")
  -m string
    	Certificate name to use (default "cert")
  -n string
    	name to give the tunnel(default eephttpd) (default "eephttpd")
  -ob int
    	Set outbound tunnel backup quantity(0 to 5) (default 1)
  -ol int
    	Set outbound tunnel length(0 to 7) (default 3)
  -oq int
    	Set outbound tunnel quantity(0 to 15) (default 2)
  -ov int
    	Set outbound tunnel length variance(-7 to 7)
  -p string
    	port to serve locally on (default "7880")
  -r	Reduce tunnel quantity when idle(true or false)
  -rc int
    	Reduce idle tunnel quantity to X (0 to 5) (default 3)
  -rt int
    	Reduce tunnel quantity after X (milliseconds) (default 600000)
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
[http://tvndxxkxcstbtqfxg7iigco6bj22ff2y6jxikmk7wqkyadkhrd4a.b32.i2p](http://tvndxxkxcstbtqfxg7iigco6bj22ff2y6jxikmk7wqkyadkhrd4a.b32.i2p)

