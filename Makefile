.SUFFIXES:
.SUFFIXES: .go

.PHONY: default
default: bin/rest-cms

# -------------------------------------------- #

DOCKER_IMAGE_BUILD  = pocka/rest-cms-build
DOCKER_IMAGE_DEPLOY = pocka/rest-cms

WORKDIR  = /go/src/github.com/pocka/rest-cms

GO_FILES = $(shell find . -type f -name "*.go")

GO_BUILDOPTS = -a \
			   -installsuffix netgo \
			   -tags netgo \
			   --ldflags '-extldflags "-static"'

GO_BUILDENV  = CGO_ENABLED=0 \
			   GOOS=linux

GO = docker run \
	 --rm \
	 -v $(shell pwd):$(WORKDIR) \
	 -w $(WORKDIR) \
	 -u $(shell id -u):$(shell id -g) \
	 $(addprefix -e ,$(GO_BUILDENV)) \
	 $(DOCKER_IMAGE_BUILD) \
	 go

GO_BUILD = $(GO) build $(GO_BUILDOPTS)
GO_TEST  = $(GO) test -v ./...

GOFMT = docker run \
		--rm \
		-v $(shell pwd):$(WORKDIR) \
		-w $(WORKDIR) \
		-u $(shell id -u):$(shell id -g) \
		$(DOCKER_IMAGE_BUILD) \
		gofmt

# -------------------------------------------- #

bin/rest-cms: bin $(GO_FILES) container/build
	$(GOFMT) -s -w $(GO_FILES)
	$(GO_BUILD) -o $@ src/main.go


bin:
	mkdir -p $@


coverage.txt: $(GO_FILES) container/build
	$(GO_TEST) -coverprofile $@ -covermode atomic


.PHONY: test
test: $(GO_FILES) container/build
	$(GO_TEST)


.PHONY: container/build
container/build:  dockerfiles/build/Dockerfile
	docker build -t $(DOCKER_IMAGE_BUILD) -f $< .


.PHONY: container/deploy
container/deploy: dockerfiles/deploy/Dockerfile bin/rest-cms
	docker build -t $(DOCKER_IMAGE_DEPLOY) -f $< .
