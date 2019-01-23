## Gladius CLI

This is an all-in-one Command Line Interface for the Gladius Network


### Setup

This CLI includes commands to run and monitor a Gladius Node. To use it you need to have the Gladius services running.

**Windows and macOS**

The `gladius-guardian` is installed as a system service by default. This means you should be able to jump into the first command with no additional setup.

**Linux**

The `gladius-guardian` is NOT installed by default. You must run `gladius-guardian` (which should already be in your path) in another terminal. Optionally you could also install the guardian by running `gladius-guardian install; gladius-guardian start`


### Commands

Use `--help` on any command for more information

**gladius** (base command)
```
$ gladius

Welcome to the Gladius CLI!

Here are the commands to setup a node (in order):

$ gladius start
$ gladius apply
$ gladius check

After you are accepted into a pool you will automatically become an edge node

To unlock your wallet after it has been created run:

$ gladius unlock

Use the -h flag to see the help menu
```

**start**

This command will start the Gladius modules needed to create a node
```
$ gladius start
Network Gateway: Started modules
Edge Daemon: Started modules
```

**apply**

Apply and submit data to a pool; allowing them to accept or reject you
```
$ gladius apply

[Gladius] Pool Address:  0xC88a29cf8F0Baf07fc822DEaA24b383Fc30f27e4 // not a real pool address!
[Gladius] What is your name? Marcelo
[Gladius] What is your email? test@test.com
[Gladius] What country are you in? USA
[Gladius] How much bandwidth do you have? (Mbps) 50
[Gladius] Why do you want to join this pool? To contribute to the Gladius Network
[Gladius] Please type your passphrase:  *******

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

**status**

See the status of the various modules

```
$ gladius status

EDGE DAEMON:	 ONLINE
NETWORK GATEWAY: ONLINE
GUARDIAN:        ONLINE
```

**unlock**

Unlock your wallet after it has been created

```
$ gladius unlock

[Gladius] Please type your passphrase:  ********
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

CLI: 0.7.0
EDGED: 0.7.0
NETWORKD: 0.7.0
GUARDIAN: 0.7.0
```

### Developer

- Use `make` to make an executable in the  `./build` folder