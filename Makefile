
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
VERSION=0.0.9993

tag:
	gothub release -s $(GITHUB_TOKEN) -u $(USER_GH) -r $(packagename) -t v$(VERSION) -d "I2P Tunnel Management tool for Go applications"

upload-linux:
	gothub upload -R -u $(USER_GH) -r "$(packagename)" -t v$(VERSION) \
		-l "Linux Terminal -`sha256sum eephttpd/$(packagename)`" -n "$(packagename)" -f "eephttpd/$(packagename)"
	gothub upload -R -u $(USER_GH) -r "$(packagename)" -t v$(VERSION) \
		-l "Linux GUI -`sha256sum eephttpd/$(packagename)-gui`" -n "$(packagename)-gui" -f "eephttpd/$(packagename)-gui"

upload-osx:
	gothub upload -R -u $(USER_GH) -r "$(packagename)" -t v$(VERSION) \
		-l "OSX Terminal -`sha256sum eephttpd/$(packagename)-osx`" -n "$(packagename)-osx" -f "eephttpd/$(packagename)-osx"
	gothub upload -R -u $(USER_GH) -r "$(packagename)" -t v$(VERSION) \
		-l "OSX GUI -`sha256sum eephttpd/$(packagename)-osx-gui`" -n "$(packagename)-osx-gui" -f "eephttpd/$(packagename)-osx-gui"

upload-windows:
	gothub upload -R -u $(USER_GH) -r "$(packagename)" -t v$(VERSION) \
		-l "Windows -`sha256sum eephttpd/$(packagename).exe`" -n "$(packagename).exe" -f "eephttpd/$(packagename).exe"

upload-orig:
	gothub upload -R -u $(USER_GH) -r "$(packagename)" -t v$(VERSION) \
		-l "Debian orig.tar.gz -`sha256sum eephttpd/$(packagename)`" -n "$(packagename)_$(VERSION).orig.tar.gz" -f "../$(packagename)_$(VERSION).orig.tar.gz"

upload: upload-linux upload-osx upload-windows upload-orig

upload-deb:
	gothub upload -R -u $(USER_GH) -r "$(packagename)" -t v$(VERSION) \
		-l "Debian($(distro) only) -`sha256sum deb/$(distro)/$(packagename)_$(VERSION).dsc`" \
		-n "$(packagename)_$(VERSION)_$(distro).dsc" \
		-f "deb/$(distro)/$(packagename)_$(VERSION).dsc"
	gothub upload -R -u $(USER_GH) -r "$(packagename)" -t v$(VERSION) \
		-l "Debian($(distro) only) -`sha256sum deb/$(distro)/$(packagename)_$(VERSION).tar.xz`" \
		-n "$(packagename)_$(VERSION)_$(distro).tar.xz" \
		-f "deb/$(distro)/$(packagename)_$(VERSION).tar.xz"
	gothub upload -R -u $(USER_GH) -r "$(packagename)" -t v$(VERSION) \
		-l "Debian($(distro) only) -`sha256sum deb/$(distro)/$(packagename)_$(VERSION)_amd64.buildinfo`" \
		-n "$(packagename)_$(VERSION)_$(distro)_amd64.buildinfo" \
		-f "deb/$(distro)/$(packagename)_$(VERSION)_amd64.buildinfo"
	gothub upload -R -u $(USER_GH) -r "$(packagename)" -t v$(VERSION) \
		-l "Debian($(distro) only) -`sha256sum deb/$(distro)/$(packagename)_$(VERSION)_amd64.changes`" \
		-n "$(packagename)_$(VERSION)_$(distro)_amd64.changes" \
		-f "deb/$(distro)/$(packagename)_$(VERSION)_amd64.changes"
	gothub upload -R -u $(USER_GH) -r "$(packagename)" -t v$(VERSION) \
		-l "Debian($(distro) only) -`sha256sum deb/$(distro)/$(packagename)_$(VERSION)_amd64.deb`" \
		-n "$(packagename)_$(VERSION)_$(distro)_amd64.deb" \
		-f "deb/$(distro)/$(packagename)_$(VERSION)_amd64.deb"
	gothub upload -R -u $(USER_GH) -r "$(packagename)" -t v$(VERSION) \
		-l "Debian($(distro) only) -`sha256sum deb/$(distro)/$(packagename)_$(VERSION)_source.changes`" \
		-n "$(packagename)_$(VERSION)_$(distro)_source.changes" \
		-f "deb/$(distro)/$(packagename)_$(VERSION)_source.changes"

mod:
	go get -u github.com/$(USER_GH)/$(packagename)@v$(VERSION)


fmt:
	find . -name '*.go' -exec gofmt -w -s {} \;

deb:
	go mod vendor
	make orig
	debuild

orig:
	tar --exclude=.git --exclude=debian -czvf ../eephttpd_$(VERSION).orig.tar.gz .

deps:
	go get -u ./...

build:
	cd eephttpd && go build -a -tags netgo -ldflags '-w -extldflags "-static"'

build-gui:
	cd eephttpd && \
		GOOS=linux \
		GOARCH=amd64 \
		CGO_ENABLED=1 \
		go build -a -tags "netgo gui" -o eephttpd-gui

 #		CGO_ENABLED=1 \

build-osx:
	cd eephttpd && \
		GOOS=darwin \
		GOARCH=amd64 \
		go build -a -tags netgo -o eephttpd-osx

build-osx-gui:
	##TODO: Cross-Compile for OSX from Linux
	#cd eephttpd && \
	#	GOOS=darwin \
	#	GOARCH=amd64 \
	#	go build -a -tags "netgo gui" -o eephttpd-osx-gui; true
	#-a -tags netgo -ldflags '-w -extldflags "-static"'
	scp admin@192.168.1.104:~/go/src/github.com/eyedeekay/eephttpd/eephttpd/eephttpd-osx-gui eephttpd/eephttpd-osx-gui


build-windows:
	cd eephttpd && \
		GOOS=windows \
		GOARCH=amd64 \
		CGO_ENABLED=1 \
		CXX=x86_64-w64-mingw32-g++ \
		CC=x86_64-w64-mingw32-gcc \
		go build -a -tags "netgo gui" -ldflags '-w -extldflags "-static"' -o eephttpd.exe

all: deps build build-gui build-osx build-osx-gui build-windows

release: deps all tag upload

install:
	install -m755 $(eephttpd)/$(eephttpd) /usr/bin/$(eephttpd)
	install -m755 $(eephttpd)/$(eephttpd)-gui /usr/bin/$(eephttpd)-gui

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
	@cat .USAGE.md > USAGE.md
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
	@echo "You can mirror the site with bittorrent as well:" >> USAGE.md
	@echo "[http://tvndxxkxcstbtqfxg7iigco6bj22ff2y6jxikmk7wqkyadkhrd4a.b32.i2p/eephttpd.torrent](http://tvndxxkxcstbtqfxg7iigco6bj22ff2y6jxikmk7wqkyadkhrd4a.b32.i2p/eephttpd.torrent)" >> USAGE.md
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
	pandoc README.md PROPAGANDA.md USAGE.md -o www/pre-index.html
	@cat www/pre-index.html >> www/index.html
	rm -f www/pre-index.html
	@echo '  </body>' >> www/index.html
	@echo '</html>' >> www/index.html
	cat README.md USAGE.md | tee www/README.md
	cp www/index.html index.html

export PKG=$(package)

release-debs: orig
	release-pdeb stable eephttpd
	distro=stable make upload-deb
	release-pdeb testing eephttpd
	distro=testing make upload-deb
	release-pdeb unstable eephttpd
	distro=unstable make upload-deb
	release-pdeb focal eephttpd
	distro=focal make upload-deb

APPNAME=$(packagename)
APPBUNDLE=$(APPNAME).app
APPBUNDLECONTENTS=$(APPBUNDLE)/Contents
APPBUNDLEEXE=$(APPBUNDLECONTENTS)/MacOS
APPBUNDLERESOURCES=$(APPBUNDLECONTENTS)/Resources
APPBUNDLEICON=$(APPBUNDLECONTENTS)/Resources
OUTFILE=eephttpd/eephttpd-osx-gui

appbundle: macosx/$(APPNAME).icns
	rm -rf $(APPBUNDLE)
	mkdir $(APPBUNDLE)
	mkdir $(APPBUNDLE)/Contents
	mkdir $(APPBUNDLE)/Contents/MacOS
	mkdir $(APPBUNDLE)/Contents/Resources
	cp macosx/Info.plist $(APPBUNDLECONTENTS)/
	cp macosx/PkgInfo $(APPBUNDLECONTENTS)/
	cp macosx/$(APPNAME).icns $(APPBUNDLEICON)/
	cp -r www $(APPBUNDLERESOURCES)/www
	cp -r etc $(APPBUNDLERESOURCES)/etc
	cp $(OUTFILE) $(APPBUNDLEEXE)/$(APPNAME)

macosx/$(APPNAME).icns: macosx/$(APPNAME)Icon.png
	rm -rf macosx/$(APPNAME).iconset
	mkdir macosx/$(APPNAME).iconset
	mogrify -resize 16x16	  -write macosx/$(APPNAME).iconset/icon_16x16.png macosx/$(APPNAME)Icon.png
	mogrify -resize 32x32	  -write macosx/$(APPNAME).iconset/icon_16x16@2x.png macosx/$(APPNAME)Icon.png
	mogrify -resize 32x32	  -write macosx/$(APPNAME).iconset/icon_32x32.png macosx/$(APPNAME)Icon.png
	mogrify -resize 64x64	  -write macosx/$(APPNAME).iconset/icon_64x64.png macosx/$(APPNAME)Icon.png
	mogrify -resize 128x128    -write macosx/$(APPNAME).iconset/icon_128x128.png macosx/$(APPNAME)Icon.png
	mogrify -resize 256x256    -write macosx/$(APPNAME).iconset/icon_128x128@2x.png macosx/$(APPNAME)Icon.png
	mogrify -resize 256x256    -write macosx/$(APPNAME).iconset/icon_256x256.png macosx/$(APPNAME)Icon.png
	mogrify -resize 512x512    -write macosx/$(APPNAME).iconset/icon_256x256@2x.png macosx/$(APPNAME)Icon.png
	mogrify -resize 512x512    -write macosx/$(APPNAME).iconset/icon_512x512.png macosx/$(APPNAME)Icon.png
	cp macosx/$(APPNAME)Icon.png macosx/$(APPNAME).iconset/icon_512x512@2x.png
	png2icns macosx/$(APPNAME).icns \
		macosx/$(APPNAME).iconset/icon_16x16.png \
		macosx/$(APPNAME).iconset/icon_32x32.png \
		macosx/$(APPNAME).iconset/icon_128x128.png \
		macosx/$(APPNAME).iconset/icon_256x256.png \
		macosx/$(APPNAME).iconset/icon_512x512.png
	#/*-c icns -o macosx/$(APPNAME).icns macosx/$(APPNAME).iconset*/
	rm -r macosx/$(APPNAME).iconset


