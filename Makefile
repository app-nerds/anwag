.DEFAULT_GOAL := help
.PHONY: help

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

VERSION := $(shell cat ./VERSION)
BUILDFLAGS=-s -w -X 'main.Version=${VERSION}'
PROJECTNAME=anwag
GOENV=CGO_ENABLED=0 GOPRIVATE="github.com/app-nerds/*" GONOPROXY="github.com/app-nerds/*"
GC=${GOENV} go build -ldflags="${BUILDFLAGS}" -mod=mod -o ${PROJECTNAME}

ifeq ($(OS),Windows_NT)
	GOOS=windows
	ifeq ($(PROCESSOR_ARCHITEW6432),AMD64)
		GOARCH=amd64
	else
		ifeq ($(PROCESSOR_ARCHITECTURE),AMD64)
			GOARCH=amd64
		endif
		ifeq ($(PROCESSOR_ARCHITECTURE),x86)
			GOARCH=x86
		endif
	endif
else
	UNAME_S := $(shell uname -s)
	ifeq ($(UNAME_S),Linux)
		GOOS=linux
	endif
	ifeq ($(UNAME_S),Darwin)
		GOOS=darwin
	endif
	UNAME_P := $(shell uname -p)
	ifeq ($(UNAME_P),x86_64)
		GOARCH=amd64
	endif
	ifneq ($(filter %86,$(UNAME_P)),)
		GOARCH=x86        
   endif
   ifneq ($(filter arm%,$(UNAME_P)),)
      GOARCH=arm64
   endif	
endif

#
# Build tasks
#

run: ## Run the app
	go run .

build-windows: ## Create a compiled Windows binary
	GOOS=windows GOARCH=amd64 ${GC}.exe

build-mac: ## Create a compiled MacOS binary (arm64)
	GOOS=darwin GOARCH=arm64 ${GC}

build-mac-amd: ## Create a compiled MacOS binary (amd64)
	GOOS=darwin GOARCH=amd64 ${GC}

build-linux: ## Create a compiled Linux binary (amd64)
	GOOS=linux GOARCH=amd64 ${GC}

build: ## Automatically determine OS and Architecture, and build an executable
	GOOS=${GOOS} GOARCH=${GOARCH} ${GC}

build-all: ## Build Windows, Linux, Mac 
	GOOS=windows GOARCH=amd64 ${GC}-windows-${VERSION}.exe
	GOOS=darwin GOARCH=arm64 ${GC}-darwin-${VERSION}
	GOOS=darwin GOARCH=amd64 ${GC}-darwin-amd64-${VERSION}
	GOOS=linux GOARCH=amd64 ${GC}-linux-${VERSION}

package: ## Package executables into a ZIP file
	zip  ./${PROJECTNAME}-linux-amd64-${VERSION}.zip ./*linux*
	zip  ./${PROJECTNAME}-windows-amd64-${VERSION}.zip ./*windows*
	zip ./${PROJECTNAME}-darwin-${VERSION}.zip ./*darwin*

