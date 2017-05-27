.SUFFIXES:
.SUFFIXES: .go

.PHONY: default
default: bin/rest-cms

# -------------------------------------------- #

DOCKER_IMAGE_BUILD  = pocka/rest-cms-build
DOCKER_IMAGE_DEPLOY = pocka/rest-cms

WORKDIR  = /go/src/github.com/pocka/rest-cms

GO_FILES = $(shell find . -type f -name "*.go")
GO_PKGS  = $(shell find src -type d -not -path "*/\.*")

GO_BUILDOPTS = -a \
			   -installsuffix netgo \
			   -tags netgo \
			   --ldflags '-extldflags "-static"'

GO_BUILDENV  = CGO_ENABLED=0 \
			   GOOS=linux

DOCKER_RUN = docker run \
			 --rm \
			 -v $(shell pwd):$(WORKDIR) \
			 -w $(WORKDIR) \
			 -u $(shell id -u):$(shell id -g) \
			 $(addprefix -e ,$(GO_BUILDENV)) \
			 $(DOCKER_IMAGE_BUILD)

GO = $(DOCKER_RUN) go

GO_BUILD = $(GO) build $(GO_BUILDOPTS)
GO_TEST  = $(GO) test
GOFMT    = $(DOCKER_RUN) gofmt

GOMETALINTER = $(DOCKER_RUN) gometalinter.v1

# -------------------------------------------- #

bin/rest-cms: bin $(GO_FILES) container/build
	$(GOFMT) -s -w $(GO_FILES)
	$(GO_BUILD) -o $@ src/main.go


bin:
	mkdir -p $@


coverage.txt: $(GO_FILES) container/build
	echo "mode: count" > $@
	$(foreach PKG,$(GO_PKGS), \
		$(GO_TEST) \
			-coverprofile $(PKG).out \
			-covermode count ./$(PKG) && \
		tail -n +2 $(PKG).out >> $@ && \
		rm $(PKG).out;\
	)


.PHONY: test
test: $(GO_FILES) container/build
	$(GOFMT) -s -w $(GO_FILES)
	$(GO_TEST) -v ./...


.PHONY: lint
lint: container/build
	$(GOFMT) -s -w $(GO_FILES)
	$(GOMETALINTER) ./...


.PHONY: container/build
container/build:  dockerfiles/build/Dockerfile
	docker build -t $(DOCKER_IMAGE_BUILD) -f $< .


.PHONY: container/deploy
container/deploy: dockerfiles/deploy/Dockerfile bin/rest-cms
	docker build -t $(DOCKER_IMAGE_DEPLOY) -f $< .
