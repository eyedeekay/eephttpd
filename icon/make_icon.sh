#/bin/sh

if [ -z "$GOPATH" ]; then
    echo GOPATH environment variable not set
    exit
fi

if [ ! -e "$GOPATH/bin/2goarray" ]; then
    echo "Installing 2goarray..."
    go get github.com/cratonica/2goarray
    if [ $? -ne 0 ]; then
        echo Failure executing go get github.com/cratonica/2goarray
        exit
    fi
fi

#if [ -z "$1" ]; then
#    echo Please specify a PNG or ICO file
#    exit
#fi

for png in $(ls *.png); do
  convert $png $(echo $png | sed 's|png|ico|')
  echo '//+build linux darwin' > $(echo $png | sed 's|.png|png.go|')
  echo '//+build windows' > $(echo $png | sed 's|.png|ico.go|')
  echo '' >> $(echo $png | sed 's|.png|png.go|')
  echo '' >> $(echo $png | sed 's|.png|ico.go|')
  cat $png | $GOPATH/bin/2goarray $(echo $png | sed 's|.png||' | sed -e "s/\b\(.\)/\u\1/g") icon >> $(echo $png | sed 's|.png|png.go|')
  cat $(echo $png | sed 's|png|ico|') | $GOPATH/bin/2goarray $(echo $png | sed 's|.png||' | sed -e "s/\b\(.\)/\u\1/g") icon >> $(echo $png | sed 's|.png|ico.go|')
done

