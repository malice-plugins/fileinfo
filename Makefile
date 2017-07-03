REPO=malice-plugins/fileinfo
ORG=malice
NAME=fileinfo
VERSION=$(shell cat VERSION)

all: build size test

build:
	docker build -t $(ORG)/$(NAME):$(VERSION) .

size:
	sed -i.bu 's/docker%20image-.*-blue/docker%20image-$(shell docker images --format "{{.Size}}" $(ORG)/$(NAME):$(VERSION)| cut -d' ' -f1)-blue/' README.md

tar: build
	docker save $(ORG)/$(NAME):$(VERSION) -o wdef.tar

go-test:
	go get
	go test -v

fi-test: test
	@echo "===> FileInfo sample Test"
	@docker run --init --rm -v $(PWD):/malware --entrypoint=bash $(ORG)/$(NAME):$(VERSION) -c "ssdeep sample" > test/ssdeep.out || true
	@docker run --init --rm -v $(PWD):/malware --entrypoint=bash $(ORG)/$(NAME):$(VERSION) -c "trid sample" > test/trid.out || true
	@docker run --init --rm -v $(PWD):/malware --entrypoint=bash $(ORG)/$(NAME):$(VERSION) -c "exiftool sample" > test/exiftool.out || true

test:
	docker run --init --rm $(ORG)/$(NAME):$(VERSION) --help
	test -f sample || wget https://github.com/maliceio/malice-av/raw/master/samples/befb88b89c2eb401900a68e9f5b78764203f2b48264fcc3f7121bf04a57fd408 -O sample
	docker run --init --rm -v $(PWD):/malware $(ORG)/$(NAME):$(VERSION) -t sample > SAMPLE.md
	docker run --init --rm -v $(PWD):/malware $(ORG)/$(NAME):$(VERSION) -V sample > results.json
	cat results.json | jq .

circle:
	http https://circleci.com/api/v1.1/project/github/${REPO} | jq '.[0].build_num' > .circleci/build_num
	http "$(shell http https://circleci.com/api/v1.1/project/github/${REPO}/$(shell cat .circleci/build_num)/artifacts${CIRCLE_TOKEN} | jq '.[].url')" > .circleci/SIZE
	sed -i.bu 's/docker%20image-.*-blue/docker%20image-$(shell cat .circleci/SIZE)-blue/' README.md

clean:
	rm sample
	docker-clean stop
	docker rmi $(ORG)/$(NAME)
	docker rmi $(ORG)/$(NAME):$(BUILD)

.PHONY: build size tags test fi-test go-test tar circle clean
