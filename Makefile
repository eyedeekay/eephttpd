
GOPATH=$(PWD)/.go

args=-r
samhost=sam-host
samport=7656
network=si
eephttpd=eephttpd

echo:
	echo $(GOPATH)

deps:
	go get -u "github.com/eyedeekay/sam-forwarder"
	go get -u "github.com/eyedeekay/sam-forwarder/config"

build:
	go build -a -tags netgo -ldflags '-w -extldflags "-static"'

release: deps build

install:
	install -m755 $(eephttpd) /usr/bin/$(eephttpd)

docker:
	docker build --no-cache -f Dockerfile -t eyedeekay/$(eephttpd) .

run-volume:
	docker run -i -t -d \
		--name $(eephttpd)-volume \
		--volume $(eephttpd):/opt/$(eephttpd)/ \
		eyedeekay/$(eephttpd); true

copy-volume:
	docker cp www/index.html $(eephttpd)-volume:/opt/$(eephttpd)/www

copy:
	docker cp www/index.html $(eephttpd):/opt/$(eephttpd)/www

stop-volume:
	docker stop $(eephttpd)-volume; true

volume: run-volume copy-volume stop-volume

run: volume
	docker rm -f $(eephttpd); true
	docker run -i -t -d \
		--network $(network) \
		--env samhost=$(samhost) \
		--env samport=$(samport) \
		--env args=$(args) \
		--network-alias $(eephttpd) \
		--hostname $(eephttpd) \
		--name $(eephttpd) \
		--restart always \
		--volumes-from $(eephttpd)-volume \
		eyedeekay/$(eephttpd)

clean:
	docker rm -f $(eephttpd) $(eephttpd)-volume

clobber: clean
	docker system prune -f

usage:
	@echo "$(eephttpd) - Static file server automatically forwarded to i2p" > USAGE.md
	@echo "============================================================" >> USAGE.md
	@echo "" >> USAGE.md
	@echo "usage:" >> USAGE.md
	@echo "------" >> USAGE.md
	@echo "" >> USAGE.md
	@echo "$(eephttpd) is a static http server which automatically runs on i2p with" >> USAGE.md
	@echo "the help of the SAM bridge. By default it will only be available from" >> USAGE.md
	@echo "the localhost and it's i2p tunnel. It can be masked from the localhost" >> USAGE.md
	@echo "using a container." >> USAGE.md
	@echo "" >> USAGE.md
	@echo '```' >> USAGE.md
	./$(eephttpd) -h  2>> USAGE.md; true
	@echo '```' >> USAGE.md
	@echo "" >> USAGE.md
	make docker-cmd
	@echo "" >> USAGE.md
	@echo "instance" >> USAGE.md
	@echo "--------" >> USAGE.md
	@echo "" >> USAGE.md
	@echo "a running instance of eephttpd with the example index file is availble on" >> USAGE.md
	@echo '[http://566niximlxdzpanmn4qouucvua3k7neniwss47li5r6ugoertzuq.b32.i2p](http://566niximlxdzpanmn4qouucvua3k7neniwss47li5r6ugoertzuq.b32.i2p)' >> USAGE.md
	@echo "" >> USAGE.md
	@cat USAGE.md

docker-cmd:
	@echo "### build in docker" >> USAGE.md
	@echo "" >> USAGE.md
	@echo '```' >> USAGE.md
	@echo "docker build --build-arg user=$(eephttpd) \\" >> USAGE.md
	@echo "    --build-arg path=example/www \\" >> USAGE.md
	@echo "    -f Dockerfile -t \\" >> USAGE.md
	@echo "    eyedeekay/$(eephttpd) ." >> USAGE.md
	@echo '```' >> USAGE.md
	@echo "" >> USAGE.md
	@echo "### Run in docker" >> USAGE.md
	@echo "" >> USAGE.md
	@echo '```' >> USAGE.md
	@echo "docker run -i -t -d \\" >> USAGE.md
	@echo "    --name $(eephttpd)-volume \\" >> USAGE.md
	@echo "    --volume $(eephttpd):/opt/$(eephttpd)/ \\" >> USAGE.md
	@echo "    eyedeekay/$(eephttpd)" >> USAGE.md
	@echo '```' >> USAGE.md
	@echo "" >> USAGE.md
	@echo '```' >> USAGE.md
	@echo "docker run -i -t -d \\" >> USAGE.md
	@echo "    --network $(network) \\" >> USAGE.md
	@echo "    --env samhost=$(samhost) \\" >> USAGE.md
	@echo "    --env samport=$(samport) \\" >> USAGE.md
	@echo "    --env args=$(args) # Additional arguments to pass to eephttpd\\" >> USAGE.md
	@echo "    --network-alias $(eephttpd) \\" >> USAGE.md
	@echo "    --hostname $(eephttpd) \\" >> USAGE.md
	@echo "    --name $(eephttpd) \\" >> USAGE.md
	@echo "    --restart always \\" >> USAGE.md
	@echo "    --volumes-from $(eephttpd)-volume \\" >> USAGE.md
	@echo "    eyedeekay/$(eephttpd)" >> USAGE.md
	@echo '```' >> USAGE.md

index:
	pandoc README.md USAGE.md -o www/index.html
