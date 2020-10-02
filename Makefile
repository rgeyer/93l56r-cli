# ########################################################## #
# Makefile for Golang Project
# Includes cross-compiling, installation, cleanup
# ########################################################## #

# Check for required command tools to build or stop immediately
EXECUTABLES = git go find pwd
K := $(foreach exec,$(EXECUTABLES),\
        $(if $(shell which $(exec)),some string,$(error "No $(exec) in PATH")))

ROOT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

BINARY=93l56r-cli
VERSION=0.0.1
BUILD=`git rev-parse HEAD`
PLATFORMS=darwin linux windows
ARCHITECTURES=amd64

# Setup linker flags option for build that interoperate with variable names in src code
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD}"

default: build

all: clean test build_all install

build:
	cd ${ROOT_DIR} && go build ${LDFLAGS} -o ${ROOT_DIR}/bin/${BINARY}

build_all:
	$(foreach GOOS, $(PLATFORMS),\
	$(foreach GOARCH, $(ARCHITECTURES), $(shell cd $(ROOT_DIR) &&  export GOOS=$(GOOS); export GOARCH=$(GOARCH); go build -v -o $(ROOT_DIR)/bin/$(BINARY)-$(GOOS)-$(GOARCH)$(if $(findstring windows,$(GOOS)),.exe,))))

install:
	go install ${LDFLAGS}

test:
	echo "No Tests yet, badboy!"
	# ginkgo -r -cover -skipPackage vendor

# Remove only what we've created
clean:
	find ${ROOT_DIR}/bin -name '${BINARY}[-?][a-zA-Z0-9]*[-?][a-zA-Z0-9]*' -delete

.PHONY: check clean install build_all all test
