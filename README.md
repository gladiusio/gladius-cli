## Gladius CLI

This is an all-in-one Command Line Interface for the Gladius Network

## Installation

`git clone https://github.com/gladiusio/gladius-cli.git`

Our workspace structure: `$GOPATH/src/github.com/gladiusio/gladius-cli`

## Usage

This CLI includes commands to run different modules of the Gladius Network. To actually use it you need one or more of those modules to be running so that the CLI can communicate with the Gladius Network. As of now this CLI supports communication with the [gladius-controld](https://github.com/gladiusio/gladius-controld) and the [gladius-networkd](https://github.com/gladiusio/gladius-networkd). **You need to be running at least one of these modules in order for the CLI to be able to do anything.**

1. [Install dep](https://github.com/golang/dep#installation)
2. `$ cd gladius-cli`
3. `$ make dependencies`
4. `$ make`
5. Run one or both of these modules: [gladius-controld](https://github.com/gladiusio/gladius-controld) or [gladius-networkd](https://github.com/gladiusio/gladius-networkd)
6. `$ ./build/gladius --help`

Optionally, you can install and run linting tools:
```sh
go get gopkg.in/alecthomas/gometalinter.v2
gometalinter.v2 --install
make lint
```

### Example

Use the base command `./build/gladius` to see usage example

**Terminal Window 1:**
```
my-computer: ~
$ gladius-control
Starting API at http://localhost:3001
```

**Terminal Window 2:**
```
my-computer: ~/go/src/github.com/gladiusio/gladius-cli
$ cd build
```
```
my-computer: ~/go/src/github.com/gladiusio/gladius-cli/build
$ ./gladius --help

Gladius CLI. This can be used to interact with various components of the Gladius Network.

Usage:
  gladius [flags]
  gladius [command]

Available Commands:
  apply       Apply to a Gladius Pool
  check       Check status of your submitted pool application
  help        Help about any command
  node        See the status of your node's networking server
  profile     See your profile information
  version     See the version of the Gladius Network

Flags:
  -h, --help          help for gladius
  -l, --level int     set the logging level (default 2)
  -t, --timeout int   set the timeout for requests in seconds (default 10)

Use "gladius [command] --help" for more information about a command.
```
```
my-computer: ~/go/src/github.com/gladiusio/gladius-cli/build
$ ./gladius apply

[Gladius] Pool Address:  0x3BbEbCe4e6E3E6DFBe70415102e457e4EE2903e3
[Gladius] What is your name? Marcelo
[Gladius] What is your email? test@test.com
[Gladius] What country are you in? USA
[Gladius] How much bandwidth do you have? (Mbps) 50
[Gladius] Why do you want to join this pool? To contribute to the Gladius Network
[Gladius] Please type your passphrase:  ****
Your application has been sent! Use gladius check to check on the status of your application!
```

### Full list of commands (in order of usage)
Use `--help` on the base command to see the help menu. Use `--help` any other command for a description of that command.

**base**
```
$ gladius

Welcome to the Gladius CLI!

Here are the commands to create a node and apply to a pool:

$ gladius apply
$ gladius check

After you are accepted into a pool, you become an edge node

Use the -h flag to see the help menu
```

**apply**

Submits your data to the pool you applied for, allowing them to accept or reject you to become a part of the pool
```
$ gladius apply

[Gladius] Pool Address:  0xC88a29cf8F0Baf07fc822DEaA24b383Fc30f27e4 // not a real pool address!
[Gladius] What is your name? Marcelo
[Gladius] What is your email? test@test.com
[Gladius] What country are you in? USA
[Gladius] How much bandwidth do you have? (Mbps) 50
[Gladius] Why do you want to join this pool? To contribute to the Gladius Network
[Gladius] Please type your passphrase:  ****
Your application has been sent! Use gladius check to check on the status of your application!
```

**check**

Check your application status to a specific pool
```
$ gladius check

[Gladius] Pool Address:  0xC88a29cf8F0Baf07fc822DEaA24b383Fc30f27e4 // not a real pool address!
Pool: 0xC88a29cf8F0Baf07fc822DEaA24b383Fc30f27e4	 Status: Pending

Once your application is approved you will automatically become an edge node!
```

**node status**

See the status of the node networking software

```
$ gladius node status

Network Daemon: Online
```

**profile**

See information regarding your node
```
$ gladius profile

Account Address: 0x8C3650F01aA308e0B56F12530378748190c6b454
```

**version**

See the versions of each module
```
$ gladius version

CLI: 0.5.5
CONTROLD: 0.5.3
NETWORKD: 0.5.2
```