all: pre build build-linux build-darwin build-windows post

pre:
	autotag write

build:
	go build

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o tiny-unzip-`autotag current`-linux-amd64
	CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -o tiny-unzip-`autotag current`-linux-386

build-darwin:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o tiny-unzip-`autotag current`-darwin-amd64
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o tiny-unzip-`autotag current`-darwin-arm64

build-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o tiny-unzip-`autotag current`-windows-amd64.exe
	CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -o tiny-unzip-`autotag current`-windows-386.exe

post:
	git restore autotag.go

clean:
	rm tiny-unzip
	rm tiny-unzip*.*
