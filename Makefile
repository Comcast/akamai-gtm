NAME=akamai-gtm
VERSION=0.0.3
TAG=v$(VERSION)
ARCH=$(shell uname -m)
PREFIX=/usr/local

all: build

updatedeps:
	go get -u golang.org/x/lint/golint
	go get -u github.com/kardianos/govendor
	govendor sync

install: build
	mkdir -p $(PREFIX)/bin
	cp -v bin/$(NAME) $(PREFIX)/bin/$(NAME)

build: updatedeps vet
	go build -ldflags "-X main.version=$(VERSION)" -o bin/$(NAME)

build_releases: updatedeps
	mkdir -p build/Linux  && GOOS=linux  go build -ldflags "-X main.version=$(VERSION)" -o build/Linux/$(NAME)
	mkdir -p build/Darwin && GOOS=darwin go build -ldflags "-X main.version=$(VERSION)" -o build/Darwin/$(NAME)
	rm -rf release && mkdir release
	tar -zcf release/$(NAME)_$(VERSION)_linux_$(ARCH).tgz -C build/Linux $(NAME)
	tar -zcf release/$(NAME)_$(VERSION)_darwin_$(ARCH).tgz -C build/Darwin $(NAME)

release: build_releases
	go get github.com/aktau/github-release
	github-release release \
		--user comcast \
		--repo akamai-gtm \
		--tag $(TAG) \
		--name "$(TAG)" \
		--description "akamai-gtm version $(VERSION)"
	ls release/*.tgz | xargs -I FILE github-release upload \
		--user comcast \
		--repo akamai-gtm \
		--tag $(TAG) \
		--name FILE \
		--file FILE

# NOTE: TravisCI will auto-deploy a GitHub release when a tag is pushed
tag:
	git tag $(TAG)
	git push origin $(TAG)

lint: updatedeps
	golint -set_exit_status

vet:
	go vet
