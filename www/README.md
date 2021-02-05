idk's home page
===============

I like to make peer-to-peer things, and think we should structure the future in
a way which builds in privacy by default in a maximally peer-to-peer way.

Projects:
---------

Besides working on [I2P](https://geti2p.net/), I have a bunch of side-projects
centered around the use of I2P, especially in Go and Javascript. Some of these
are also I2P labs projects.

### I2P Webextensions and Browser Research

 * [I2P in Private Browsing Mode for Firefox](I2P-in-Private-Browsing-Mode-Firefox/)
 inspired by Brave, this browser extension enforces a few privacy rules for
 Firefox, then implements a set of "Container Tabs" which can be used to browse
 I2P in a way which is automatic and safe.
 * [I2P Configuration Helper for Chromium](I2P-Configuration-For-Chromium/)
 A fork of the Firefox plugin before it used container tabs, after the user sets
 up an I2P browsing profile this plugin can automatically set up the profile to
 use I2P with the maximum privacy available from Chromium.
 * [I2P Browser Fingerprint Gallery](I2P-Browser-Attackability-Evaluation/)
 This is an ongoing evaluation of the various ways there are to configure
 browsers for I2P and a developing rating system for them.
 
### I2P Git Hosting

I host git services on I2P at [git.idk.i2p](http://git.idk.i2p/), with a non-anonymous
mirror available at [i2pgit.org](https://i2pgit.org). I provide this service freely to
the I2P community and it is open to the public, but do have a Terms of Service which
is predicated on my own threat model. If the TOS is not acceptable to you, I highly
encourage you to run your own gitlab instance using the instructions I wrote, available
on the I2P project [I2P Site](http://i2p-projekt.i2p/en/docs/applications/gitlab) and
on the [Web](https://geti2p.net/en/docs/applications/gitlab).

### What's Weird about this I2P Site?

This I2P site uses an experimental feature of I2P in Private browsing called 
X-I2P-TorrentLocation. If you are using the latest version of the extension, you
may notice that there is a pageAction available in the URL bar(It's the little I2P
logo). If you click that pageAction and follow the magnet link, you will begin to
download a torrent named idk.i2p. As the torrent completes, the extension will begin
to replace on-page resources hosted on my server with exactly the same files, except
shared and downloaded via I2PSnark. The result is a sort of distributed, voluntary
pseudo-CDN which makes it possible to do things like embed videos directly in your
I2P Site and actually have them play completely. All of this is accomplished, of
course, by cheating. If you download the torrent, the file is on your computer, so
of course it's available in a reliable way. Besides that, even if your content
goes down, any of your visitors will be able to reproduce your site on a new hostname,
which may provide a level of resistance to being taken down. How it affects traffic
flows remains to be seen, but it means that some I2P users will be fetching less
content via their HTTP Proxies, and will be doing it less repeatedly.

 * [InfoGraphic Gallery for testing X-I2P-TorrentLocation](infographics.html) This
 page is to test X-I2P-TorrentLocation. It embeds a bunch of large images/infographics
 I collected off of reddit and is intentionally heavy so it may be slow to load. It will
 use torrent-based resources if I2P In Private Browsing mode is installed and the idk.i2p
 torrent is downloaded.
 * [Videos about I2P Gallery for testing X-I2P-TorrentLocation](video.html) This page
 is even heavier, it embeds videos that have to do with I2P and other crypto/privacy/overlay
 networking related topics.
 * [Plugin Archive](plugins.html) This page is my mirror of the plugin archive at
 [stats.i2p](http://stats.i2p/i2p/plugins). I created one here because the *utility* of
 X-I2P-TorrentLocation is that it allows you to mingle the versatile presentation abilities
 of hypertext with redundant, peer-to-peer resources, and by using it for plugins, we can
 make it much harder to take them down by taking down their archives. The same would apply
 for any software, actually, this is just where I started.

### Go(golang) I2P Tools

 * [samcatd](https://github.com/eyedeekay/sam-forwarder) a.k.a. sam-forwarder
 many of the other applications use sam-forwarder as a way of automatically
 configuring i2ptunnels, including:
  - [httptunnel](https://github.com/eyedeekay/httptunnel) is a standalone http
  proxy for I2P that uses SAM and implements an interface like sam-forwarder.
  - [eephttpd](https://github.com/eyedeekay/eephttpd) is a simple static http
  server with markdown parsing support.
  - [gitsam](https://github.com/eyedeekay/gitsam) is a super-simple git
  repository setup built on eephttpd and [gitkit]().
  - [reposam](https://github.com/eyedeekay/reposam) is a binary deb repository
  built on [repogen]().
  - [samtracker](https://github.com/eyedeekay/samtracker) is a simple torrent
  tracker built upon [retracker]().
  - [cowyosam](https://github.com/eyedeekay/cowyosam) is a pastebin-wiki hybrid
  built on [cowyo]()
  - [colluding_sites_attack](https://github.com/eyedeekay/colluding_sites_attack)
  is a tool for fingerprinting browsers as they visit eepSites to determine if
  useful information can be extracted.
  - [outproxy](https://github.com/eyedeekay/outproxy) is a standalone outproxy
  built on SAM. Definitely don't use it if you don't know what you're in for.
  - [libanonvpn](https://github.com/RTradeLtd/libanonvpn) is a VPN library and
  terminal application that uses SAM Datagrams. Sort of like onioncat, but
  cooler.
 * [checki2cp](https://github.com/eyedeekay/checki2cp) is an I2P router presence
 detection tool. Use it to find out if an I2P router is installed.
 * [goSam](https://github.com/eyedeekay/goSam) is a SAM library for Go that
 implements an HTTP Transport.
 * [i2pdig](https://github.com/eyedeekay/i2pdig) is dig, but for I2P. It's been
 a while, I'll update it soon.
 * [iget](https://github.com/eyedeekay/iget) iget is an eepget clone, with some
 extra features and room to grow.
 * [keyto](https://github.com/eyedeekay/keyto) is a text key conversion tool.
 * [sam3](https://github.com/eyedeekay/sam3) is another SAM library for Go, but
 it implements a net.Conn and net.Packetconn making it a near drop-in
 replacement for regular connections.

Blog:
-----

#### Sun Nov 26 03:21:12 EST 2017

Hi. This is the blog where I'm going to document all the wierd stuff I do on my
home network. I'm most passionate about the areas where I am relatively free of
constraints, and for me, that is in hobby computing in my own home. But since
it's not a place with an IT staff and other organizational resources, I
sometimes do wierd, ill-advised things to get my computers just the way I like
them.

Also I'm pretty bad at blogging.

#### Mon Jan 22 12:41:21 EST 2018

Getting nervous, about to flash an up-to-date coreboot port to my netbook via
a ch341a flasher. I'm about 99% sure I'm not going to hurt anything, but who
knows?

#### Tue Mar 31, 15:04:40 EST 2020

See, I told you I was pretty bad at blogging. Over 2 years. Lots of code though.

#### Sun Oct 11, 04:08:56 EDT 2020

Note to self: from now on, build the site with:

`make all && make seed && git commit -am "Example commit message" && git push --all`
