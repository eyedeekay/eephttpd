
GOPATH=$(PWD)/.go

args=-r
samhost=sam-host
samport=7656
network=si

echo:
	echo $(GOPATH)

deps:
	go get -u "github.com/eyedeekay/sam-forwarder"
	go get -u "github.com/eyedeekay/sam-forwarder/config"

build:
	go build -a -tags netgo -ldflags '-w -extldflags "-static"'

release: deps build

install:
	install -m755 eephttpd /usr/bin/eephttpd

docker:
	docker build -f Dockerfile -t eyedeekay/eephttpd .

volume:
	docker run -i -t -d \
		--name eephttpd-volume \
		--volume eephttpd:/opt/eephttpd/ \
		eyedeekay/eephttpd; true
	docker cp www/* eephttpd-volume:/opt/eephttpd/www
	docker stop eephttpd-volume; true

run: volume
	docker run -i -t -d \
		--network $(network) \
		--env samhost=$(samhost) \
		--env samport=$(samport) \
		--env args=$(args) \
		--network-alias eephttpd \
		--hostname eephttpd \
		--name eephttpd \
		--restart always \
		--volumes-from eephttpd-volume \
		eyedeekay/eephttpd
