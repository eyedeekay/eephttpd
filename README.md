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

## Docker

First, create a volume helper:

        docker run -i -t -d \
            --name eephttpd-volume \
            --volume eephttpd:/home/eephttpd/ \
            eyedeekay/eephttpd

Then, copy files you wish to serve into the volume folder:

        docker cp www/* eephttpd-volume:/home/eephttpd/www

Stop the volume helper:

        docker stop eephttpd-volume

and run the container:

        docker run -i -t -d \
            --env samhost=$(samhost) \
            --env samport=$(samport) \
            --env args=$(args) \
            --network-alias eephttpd \
            --hostname eephttpd \
            --name eephttpd \
            --restart always \
            --volumes-from eephttpd-volume \
            eyedeekay/eephttpd
