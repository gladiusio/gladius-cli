# Gladius CLI

This is an all-in-one Command Line Interface for the Gladius Network

## Installation

`git clone https://github.com/gladiusio/gladius-cli.git`

Our workspace structure: `$GOPATH/src/github.com/gladiusio/gladius-cli`

## Usage

This CLI includes commands to run different modules of the Gladius Network. To actually use it you need one or more of those modules to be running so that the CLI can communicate with the Gladius Network. As of now this CLI supports communication with the [gladius-control-daemon](https://github.com/gladiusio/gladius-control-daemon) and the [gladius-networkd](https://github.com/gladiusio/gladius-networkd). **You need to be running at least one of these modules in order for the CLI to be able to do anything.**

1. [Install](https://github.com/golang/dep#installation) `dep`
2. `$ cd gladius-cli`
3. `$ make dependencies`
4. `$ make cli`
5. Run one or both of these modules: [gladius-control-daemon](https://github.com/gladiusio/gladius-control-daemon) or [gladius-networkd](https://github.com/gladiusio/gladius-networkd)
6. [Write your environment file](./setup.md)
7. `$ gladius-cli/build`
8. `$ ./gladius-cli --help`
  * This runs the CLI executable and displays all available commands
9. Use `./gladius-cli [command] --help` for usage and help information for specific commands

### Example

Terminal Window 1:
```
my-computer: ~
> gladius-control
Running at http://localhost:3000
```

Terminal Window 2:
```
my-computer: ~/go/src/github.com/gladiusio/gladius-cli
> cd build

my-computer: ~/go/src/github.com/gladiusio/gladius-cli/build
> ./gladius-cli --help
Gladius CLI. This can be used to interact with various components of the Gladius Network.

Usage:
  gladius-cli [flags]
  gladius-cli [command]

Available Commands:
  apply       Apply to a Gladius Pool
  check       Check status of your submitted pool application
  create      Deploy a new Node smart contract
  echo        Echo anything to the screen
  edge        Start the edge daemon
  help        Help about any command
  test        Test function

Flags:
  -h, --help   help for gladius-cli

Use "gladius-cli [command] --help" for more information about a command.

my-computer: ~/go/src/github.com/gladiusio/gladius-cli/build
> ./gladius-cli create
? What is your name? Marcelo
? What is your email? test@email.com
Tx: 0x12aaa4517e8c0899791de40403d7c0a9a5b44f904e0bfe19c2207d9e338ba68e	 Status: Pending
Tx: 0x12aaa4517e8c0899791de40403d7c0a9a5b44f904e0bfe19c2207d9e338ba68e	 Status: Successful
Node created!
Tx: 0x3e39c6892195cde9dda7944f47030387d752087955f599cb3c2d538204bffd8e	 Status: Pending
Tx: 0x3e39c6892195cde9dda7944f47030387d752087955f599cb3c2d538204bffd8e	 Status: Successful
Node data set!

0x4607210e97eD3e7D43929f0eF324c259d4Fa0690 //node address

```
