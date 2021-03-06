help:	## to show this help.
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

gnutar:= $(shell gtar --version 2> /dev/null)

package: ## make tarball package out of artifacts
	cp src/plugin.yaml ./build
ifdef gnutar
	gtar -cvf build/pusha.tar build
else
	tar -cvf build/pusha.tar build
endif

linux: ## build linux
	docker run -v `pwd`/:/go/src/github.com/app -w /go/src/github.com/app haaartland/golang-glide-builder bash -c "glide install && env GOOS=linux GOARCH=amd64 go build -o build/linux/pusher src/pusher.go"

osx: ## build linux
	docker run -v `pwd`/:/go/src/github.com/app -w /go/src/github.com/app haaartland/golang-glide-builder bash -c "glide install && env GOOS=darwin GOARCH=amd64 go build -o build/osx/pusher src/pusher.go"

clean: ## clean working dir
	test -n build && rm -rf build
	test -n vendor && rm -rf vendor

all: clean osx linux package
