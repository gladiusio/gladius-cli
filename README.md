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
Running at http://localhost:3001
```

**Terminal Window 2:**
```
my-computer: ~/go/src/github.com/gladiusio/gladius-cli
$ cd build

my-computer: ~/go/src/github.com/gladiusio/gladius-cli/build
$ ./gladius --help

Gladius CLI. This can be used to interact with various components of the Gladius Network.

Usage:
  gladius [flags]
  gladius [command]

Available Commands:
  apply       Apply to a Gladius Pool
  check       Check status of your submitted pool application
  create      Deploy a new Node smart contract
  echo        Echo anything to the screen
  edge        Start the edge daemon
  help        Help about any command
  test        Test function

Flags:
  -h, --help   help for gladius

Use "gladius [command] --help" for more information about a command.

my-computer: ~/go/src/github.com/gladiusio/gladius-cli/build
$ ./gladius create

[Gladius] What is your name? Marcelo
[Gladius] What is your email? test@email.com
[Gladius] Enter your password: **********

Tx: 0x12aaa4517e8c0899791de40403d7c0a9a5b44f904e0bfe19c2207d9e338ba68e	 Status: Pending
Tx: 0x12aaa4517e8c0899791de40403d7c0a9a5b44f904e0bfe19c2207d9e338ba68e	 Status: Successful
Node created!

Tx: 0x3e39c6892195cde9dda7944f47030387d752087955f599cb3c2d538204bffd8e	 Status: Pending
Tx: 0x3e39c6892195cde9dda7944f47030387d752087955f599cb3c2d538204bffd8e	 Status: Successful
Node data set!

Node Address: 0x4607210e97eD3e7D43929f0eF324c259d4Fa0690

```

### Full list of commands (in order of usage)
Use `--help` on the base command to see the help menu. Use `--help` any other command for a description of that command.

**base**
```
$ gladius

Welcome to the Gladius CLI!

Here are the commands to create a node and apply to a pool in order:

$ gladius create
$ gladius apply
$ gladius check

After you are accepted into a pool, you can become an edge node:

$ gladius node start

Use the -h flag to see the help menu
```

**create**

Deploys a new Gladius Node smart contract containing the encrypted version of the data you submitted. If you enter in the wrong information you can just run the command again to make a new node.
```
$ gladius create

[Gladius] What is your name? Marcelo Test
[Gladius] What is your email? email@test.com
[Gladius] Please type your password:  ********

Tx: 0xb37a017d2877ab7350e0c7199326bc97bda32e4d8ae46c6aaecc2f9b0cd3b133	 Status: Pending...
Tx: 0xb37a017d2877ab7350e0c7199326bc97bda32e4d8ae46c6aaecc2f9b0cd3b133	 Status: Successful
Node created!

Tx: 0x6931f0394684ebef6c0fa9c83ccf1ae7fa2811b93b4480fcf0ba163e8eb03ff6	 Status: Pending...
Tx: 0x6931f0394684ebef6c0fa9c83ccf1ae7fa2811b93b4480fcf0ba163e8eb03ff6	 Status: Successful
Node data set!

Node Address: 0xb04578990b1cbb515b8764ca8778e5ba7f6eb8e5

Use gladius apply to apply to a pool
```

**apply**

Submits the data to a specific pool, allowing them to accept or reject you to become a part of the pool
```
$ gladius apply

[Gladius] Pool Address:  0xC88a29cf8F0Baf07fc822DEaA24b383Fc30f27e4
[Gladius] Please type your password:  ********

Tx: 0x14e796ce7939c035586ff2b6f26e1ad9db71be7a760715debbad68b4cb9d9496	 Status: Pending
Tx: 0x14e796ce7939c035586ff2b6f26e1ad9db71be7a760715debbad68b4cb9d9496	 Status: Successful

Application sent to pool!
Use gladius check to check your application status
```

**check**

Check your application status to a specific pool
```

$ gladius check

[Gladius] Pool Address:  0xC88a29cf8F0Baf07fc822DEaA24b383Fc30f27e4
Pool: 0xC88a29cf8F0Baf07fc822DEaA24b383Fc30f27e4	 Status: Pending

Use gladius node start to start the node networking software
```

**node [start | stop | status]**

Start/stop or check the status of the node networking software
```

$ gladius node start
Network Daemon:	 Started the server

Use gladius node stop to stop the node networking software
Use gladius node status to check the status of the node networking software
```

```

$ gladius node stop
Network Daemon:	 Stopped the server

Use gladius node start to start the node networking software
Use gladius node status to check the status of the node networking software
```

```

$ gladius node status
Network Daemon:	 Server is Running

Use gladius node start to start the node networking software
Use gladius node stop to stop the node networking software
```

**profile**

See information regarding your node
```
$ gladius profile

Account Address: 0x8C3650F01aA308e0B56F12530378748190c6b454
Node Address: 0xf15aea30341982b117583f36cf516f6cea5ddf91
Node Name: Marcelo
Node Email: marcelo@test.com
Node IP: 12.12.123.12
```