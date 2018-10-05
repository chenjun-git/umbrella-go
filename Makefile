export GOPATH := $(CURDIR)/_project
export GOBIN := $(CURDIR)/bin

CURRENT_GIT_GROUP := umbrella-go
CURRENT_GIT_REPO := .
COMMONENVVAR ?= GOOS=linux GOARCH=amd64
BUILDENVVAR ?= cgo_enabled=0

all: deps linux_build

folder_dep:
	mkdir -p $(CURDIR)/_project/src/$(CURRENT_GIT_GROUP)
	test -d $(CURDIR)/_project/src/$(CURRENT_GIT_GROUP)/$(CURRENT_GIT_REPO) || ln -s $(CURDIR) $(CURDIR)/_project/src/$(CURRENT_GIT_GROUP)

deps: folder_dep
	mkdir -p $(CURDIR)/vendor
	glide install

build: folder_dep
	$(BUILDENVVAR) go build -o $(GOBIN)/umbrella -ldflags "-X main.BuildTime=`date '+%Y-%m-%d_%I:%M:%S%p'` -X main.BuildGitHash=`git rev-parse HEAD` -X main.BuildGitTag=`git describe --tags`" $(CURRENT_GIT_GROUP)/$(CURRENT_GIT_REPO)

linux_build: deps
	$(BUILDENVVAR) make build

proto:

test: folder_dep
	go test $(CURRENT_GIT_GROUP)/$(CURRENT_GIT_REPO)/umbrella-common/caller/grpc
	go test $(CURRENT_GIT_GROUP)/$(CURRENT_GIT_REPO)/umbrella-common/caller/http
	go test $(CURRENT_GIT_GROUP)/$(CURRENT_GIT_REPO)/umbrella-common/json
	go test $(CURRENT_GIT_GROUP)/$(CURRENT_GIT_REPO)/umbrella-common/errors
	go test $(CURRENT_GIT_GROUP)/$(CURRENT_GIT_REPO)/umbrella-common/lang
	go test $(CURRENT_GIT_GROUP)/$(CURRENT_GIT_REPO)/umbrella-common/redis

clean:
	@rm -rf bin _project
	@cd ./vender && rm -rf `ls | grep -v 'vender.json'` && cd ..

.PHONY: install test clean proto
