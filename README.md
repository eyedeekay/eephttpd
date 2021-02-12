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
name `eephttpd.torrent`. This allows people to mirror the site's exact
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

