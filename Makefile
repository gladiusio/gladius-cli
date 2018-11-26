##
## Makefile to test and build the gladius binaries
##

##
# GLOBAL VARIABLES
##

# commands for go
GOMOD=GO111MODULE=on
GOBUILD=$(GOMOD) go build
GOTEST=$(GOMOD) go test
GOCLEAN=$(GOMOD) go clean

# if we are running on a windows machine
# we need to append a .exe to the
# compiled binary
BINARY_SUFFIX=
ifeq ($(OS),Windows_NT)
	BINARY_SUFFIX=.exe
endif

ifeq ($(GOOS),windows)
	BINARY_SUFFIX=.exe
endif

# code source and build directories
SRC_DIR=./cmd
DST_DIR=./build

BINARY=gladius$(BINARY_SUFFIX)

# source of edged
CLI_SRC=$(SRC_DIR)

# destination of compiled edged
CLI_DEST=$(DST_DIR)/$(BINARY)

##
# MAKE TARGETS
##

# default, will be called if no arguments supplied
all: clean test executable

# delete anything in the build dir and clean
clean:
	rm -rf ./build/*
	$(GOCLEAN) $(CLI_SRC)

# test edged
test: $(CLI_SRC)
	$(GOTEST) $(CLI_SRC)

# Made for macOS at the moment
# Install gcc cross compilers for macOS
# `brew install mingw-w64` - windows
# `brew install FiloSottile/musl-cross/musl-cross` - linux
release: clean release-win release-linux release-mac

release-win:
	CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(DST_DIR)/release/windows/$(BINARY).exe $(CLI_SRC)
release-linux:
	CGO_ENABLED=1 CC=x86_64-linux-musl-gcc GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(DST_DIR)/release/linux/$(BINARY) $(CLI_SRC)
release-mac:
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(DST_DIR)/release/macos/$(BINARY) $(CLI_SRC)

# test and compile the edged
executable:
	$(GOBUILD) -o $(CLI_DEST) $(CLI_SRC)
