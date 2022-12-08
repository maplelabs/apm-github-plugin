BINARY_NAME?=github-audit
SRC_PATH=${PWD}
VERSION?=dev
BuildTime=$(shell date +%T-%D-%Z)
GitCommit=$(shell git rev-parse HEAD)
LDFLAGS=-ldflags="-s -w -X main.version=$(VERSION) -X main.buildTime=$(BuildTime) -X main.commit=$(GitCommit)"
GCFLAGS=-gcflags=-trimpath=$(SRC_PATH)
ASMFLAGS=-asmflags=-trimpath=$(SRC_PATH)

build:
	CGO_ENABLED=1 go build -v -race -trimpath $(GCFLAGS) $(ASMFLAGS) $(LDFLAGS) -o $(SRC_PATH)/$(BINARY_NAME) $(SRC_PATH)/cmd/github-audit.go

static:
	CGO_ENABLED=0 go build -v -a -trimpath $(GCFLAGS) $(ASMFLAGS) $(LDFLAGS) -o $(SRC_PATH)/$(BINARY_NAME) $(SRC_PATH)/cmd/github-audit.go