
#GOPATH=$(PWD)/.go

args=-r
samhost=sam-host
samport=7656
network=si
eephttpd=eephttpd

echo:
	echo $(GOPATH)

USER_GH=eyedeekay
packagename=eephttpd
VERSION=0.0.998

tag:
	gothub release -s $(GITHUB_TOKEN) -u $(USER_GH) -r $(packagename) -t v$(VERSION) -d "I2P Tunnel Management tool for Go applications"

upload:
	gothub upload -R -u $(USER_GH) -r "$(packagename)" -t v$(VERSION) -l "`sha256sum eephttpd/$(packagename).exe`" -n "$(packagename).exe" -f "eephttpd/$(packagename).exe"
	gothub upload -R -u $(USER_GH) -r "$(packagename)" -t v$(VERSION) -l "`sha256sum eephttpd/$(packagename)`" -n "$(packagename)" -f "eephttpd/$(packagename)"
	gothub upload -R -u $(USER_GH) -r "$(packagename)" -t v$(VERSION) -l "`sha256sum eephttpd/$(packagename)-gui`" -n "$(packagename)-gui" -f "eephttpd/$(packagename)-gui"
	gothub upload -R -u $(USER_GH) -r "$(packagename)" -t v$(VERSION) -l "`sha256sum eephttpd/$(packagename)-osx`" -n "$(packagename)-osx" -f "eephttpd/$(packagename)-osx"
	gothub upload -R -u $(USER_GH) -r "$(packagename)" -t v$(VERSION) -l "`sha256sum eephttpd/$(packagename)-osx-gui`" -n "$(packagename)-osx-gui" -f "eephttpd/$(packagename)-osx-gui"

mod:
	go get -u github.com/$(USER_GH)/$(packagename)@v$(VERSION)


fmt:
	find . -name '*.go' -exec gofmt -w -s {} \;

orig:
	tar --exclude=.git --exclude=debian -czvf ../eephttpd_0.0~git20181031.a4b6058.orig.tar.gz .

deps:
	go get -u ./...

build:
	cd eephttpd && go build -a -tags netgo -ldflags '-w -extldflags "-static"'
	
build-gui:
	cd eephttpd && \
		GOOS=linux && \
		GOARCH=amd64 && \
		CGO_ENABLED=1 && \
		go build -a -tags "netgo gui" -o eephttpd-gui

build-osx:
	cd eephttpd && \
		GOOS=darwin && \
		GOARCH=amd64 && \
		go build -a -tags netgo -o eephttpd-osx

build-osx-gui:
	cd eephttpd && \
		GOOS=darwin && \
		GOARCH=amd64 && \
		CGO_ENABLED=1 && \
		go build -a -tags "netgo gui" -o eephttpd-osx-gui

#-a -tags netgo -ldflags '-w -extldflags "-static"'

build-windows:
	cd eephttpd && \
		GOOS=windows \
		GOARCH=amd64 \
		CGO_ENABLED=1 \
		CXX=x86_64-w64-mingw32-g++ \
		CC=x86_64-w64-mingw32-gcc \
		go build -a -tags "netgo gui" -ldflags '-w -extldflags "-static"' -o eephttpd.exe

all: build build-gui build-osx build-osx-gui build-windows

release: deps all tag upload

install:
	install -m755 $(eephttpd)/$(eephttpd) /usr/bin/$(eephttpd)

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
	./$(eephttpd)/eephttpd -h  2>> USAGE.md; true
	@echo '```' >> USAGE.md
	@echo "" >> USAGE.md
	make docker-cmd
	@echo "" >> USAGE.md
	@echo "instance" >> USAGE.md
	@echo "--------" >> USAGE.md
	@echo "" >> USAGE.md
	@echo "a running instance of eephttpd with the example index file is availble on" >> USAGE.md
	@echo '[http://tvndxxkxcstbtqfxg7iigco6bj22ff2y6jxikmk7wqkyadkhrd4a.b32.i2p](http://tvndxxkxcstbtqfxg7iigco6bj22ff2y6jxikmk7wqkyadkhrd4a.b32.i2p)' >> USAGE.md
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
	@echo '<!DOCTYPE html>' > www/index.html
	@echo '<html lang="en">' >> www/index.html
	@echo '  <head>' >> www/index.html
	@echo '    <meta charset="utf-8">' >> www/index.html
	@echo '    <title>title</title>' >> www/index.html
	@echo '    <link rel="stylesheet" href="style.css">' >> www/index.html
	@echo '    <script src="script.js"></script>' >> www/index.html
	@echo '  </head>' >> www/index.html
	@echo '  <body>' >> www/index.html
	pandoc README.md USAGE.md -o www/pre-index.html
	@cat www/pre-index.html >> www/index.html
	rm -f www/pre-index.html
	@echo '  </body>' >> www/index.html
	@echo '</html>' >> www/index.html
	cat README.md USAGE.md | tee www/README.md

