##
## Makefile to test and build the gladius binaries
##

##
# GLOBAL VARIABLES
##

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

CLI_SRC=$(SRC_DIR)
CLI_DEST=$(DST_DIR)/gladius$(BINARY_SUFFIX)

# commands for go
GOBUILD=go build
GOTEST=go test
##
# MAKE TARGETS
##

# general make targets
all: dependencies cli

clean:
	rm -rf ./build/*
	go clean

# dependency management
dependencies:
	# installing dependencies
	dep ensure

# build steps
test: $(CLI_SRC)
	$(GOTEST) $(CLI_SRC)

cli: test
	# compiling binary
	$(GOBUILD) -o $(CLI_DEST) $(CLI_SRC)
