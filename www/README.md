# eephttpd

So much for a simple static file server.

eepHttpd is a web hosting tool for content in I2P, I2P sites, and
I2P torrents. On it's own, it's basically a static directory server
with limited scripting abilities written in pure-Go. 

However, it responds differently to different kinds of interaction.

 1. If a Git client attempts to access it, then they will be able to clone
the entire site, up from the document root(So **Use this for things**
**You want to *SHARE ANONYMOUSLY* with a large audience**, not for
things you want to keep secret.) This allows people to clone the site
in order to mirror it.
 2. When any file is changed in the docroot, eephttpd generates a
multi-file torrent of the site and places it in the docroot under the file
name `$SITEHOSTNAME.torrent`. This allows people to mirror the site's exact
content, and participate in keeping the site's content up.
 2. When a browser with I2P in Private Browsing connects to it, it creates
a magnet link and replies with it as an `X-I2P-TORRENTLOCATION` header. In
this way, the browser can help the user download the whole web site using
Bittorrent and substitute HTTP resources for Bittorrent resources when
they are ready.
 4. If a Torrent client attempts to access the `/a` URL, it is forwarded
to an Open Torrent Tracker. **Every single eephttpd site is also an open**
**torrent tracker.** Moreover, every single eephttpd site treats itself as
the primary tracker for the whole-site torrent it generates. **This is**
**intended to encourage the distribution of *open trackers* on I2P.**

So... more to come on why this is cool.

In order to build a .deb file, either use `checkinstall` or run:

        go mod vendor
        make orig
        debuild -us -uc

or just run:

        make deb

eephttpd - Static file server automatically forwarded to i2p
============================================================

usage:
------

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
    	the directory of static files to host(default ./www) (default "./www")
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

Soon, you should also be able to mirror the site with bittorrent as well:
[http://tvndxxkxcstbtqfxg7iigco6bj22ff2y6jxikmk7wqkyadkhrd4a.b32.i2p/tvndxxkxcstbtqfxg7iigco6bj22ff2y6jxikmk7wqkyadkhrd4a.b32.i2p.torrent](http://tvndxxkxcstbtqfxg7iigco6bj22ff2y6jxikmk7wqkyadkhrd4a.b32.i2p/tvndxxkxcstbtqfxg7iigco6bj22ff2y6jxikmk7wqkyadkhrd4a.b32.i2p.torrent)

