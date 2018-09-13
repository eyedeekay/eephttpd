# eephttpd

So, basically everything I put on i2p is a bunch of static files. Until now, I
tried to host them using darkhttpd(a fork of lighttpd from Alpine which
functions as a static Web Server) and by adding tunnel configuration information
to tunnels.conf for i2pd. This is easier than manipulating a web interface, but
still tedious and kind of error-prone. So instead, this serves simple static
sites directly to i2p via the SAM API.

to build:

        git clone https://github.com/eyedeekay/eephttpd && cd eephttpd
        go get -u "github.com/eyedeekay/sam-forwarder"
        go get -u "github.com/eyedeekay/sam-forwarder/config"
        go build

to run:

        ./eephttpd

will serve the files from ./www, and store i2p keys in the working directory.

## [Usage](USAGE.md)
