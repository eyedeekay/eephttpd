# eephttpd
Serving simple static sites directly to i2p via the SAM API.

to build:

        git clone https://github.com/eyedeekay/eephttpd && cd eephttpd
        go get -u "github.com/eyedeekay/sam-forwarder"
        go get -u "github.com/eyedeekay/sam-forwarder/config"
        go build

to run:

        ./eephttpd

will serve the files from ./www, and store i2p keys in the working directory.
