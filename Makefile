REPO=malice
NAME=fileinfo
VERSION=$(shell cat VERSION)

all: build size test

build:
	docker build -t $(REPO)/$(NAME):$(VERSION) .

size: build
	sed -i.bu 's/docker%20image-.*-blue/docker%20image-$(shell docker images --format "{{.Size}}" $(REPO)/$(NAME):$(BUILD)| cut -d' ' -f1)%20MB-blue/' README.md

test:
	docker run $(REPO)/$(NAME):$(VERSION) /bin/cat | jq .

.PHONY: build release test
