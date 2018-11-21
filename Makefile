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

# source of edged
CLI_SRC=$(SRC_DIR)

# destination of compiled edged
CLI_DEST=$(DST_DIR)/gladius$(BINARY_SUFFIX)

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

# test and compile the edged
executable:
	$(GOBUILD) -o $(CLI_DEST) $(CLI_SRC)
