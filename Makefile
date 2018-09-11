
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
	docker build --no-cache -f Dockerfile -t eyedeekay/eephttpd .

run-volume:
	docker run -i -t -d \
		--name eephttpd-volume \
		--volume eephttpd:/opt/eephttpd/ \
		eyedeekay/eephttpd; true

copy-volume:
	docker cp www/index.html eephttpd-volume:/opt/eephttpd/www

copy:
	docker cp www/index.html eephttpd:/opt/eephttpd/www

stop-volume:
	docker stop eephttpd-volume; true

volume: run-volume copy-volume stop-volume

run:
	docker rm -f eephttpd; true
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

clean:
	docker rm -f eephttpd eephttpd-volume

clobber: clean
	docker system prune -f
