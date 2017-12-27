NAME=akamai-gtm
VERSION=0.0.3
TAG=v$(VERSION)
ARCH=$(shell uname -m)
PREFIX=/usr/local

all: lint vet build

updatedeps:
	go get -u github.com/golang/lint/golint
	go get -u github.com/kardianos/govendor
	govendor sync

install: build
	mkdir -p $(PREFIX)/bin
	cp -v bin/$(NAME) $(PREFIX)/bin/$(NAME)

build: updatedeps
	go build -ldflags "-X main.version=$(VERSION)" -o bin/$(NAME)

build_releases: updatedeps
	mkdir -p build/Linux  && GOOS=linux  go build -ldflags "-X main.version=$(VERSION)" -o build/Linux/$(NAME)
	mkdir -p build/Darwin && GOOS=darwin go build -ldflags "-X main.version=$(VERSION)" -o build/Darwin/$(NAME)
	rm -rf release && mkdir release
	tar -zcf release/$(NAME)_$(VERSION)_linux_$(ARCH).tgz -C build/Linux $(NAME)
	tar -zcf release/$(NAME)_$(VERSION)_darwin_$(ARCH).tgz -C build/Darwin $(NAME)

release: build_releases
	go get github.com/progrium/gh-release
	gh-release create Comcast/$(NAME) $(VERSION) $(shell git rev-parse --abbrev-ref HEAD)

lint: updatedeps
	golint -set_exit_status

# vet runs the Go source code static analysis tool `vet` to find
# any common errors.
vet:
	@go tool vet 2>/dev/null ; if [ $$? -eq 3 ]; then \
		go get golang.org/x/tools/cmd/vet; \
	fi
	@echo "go tool vet $(VETARGS)"
	@go tool vet $(VETARGS) . ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi
