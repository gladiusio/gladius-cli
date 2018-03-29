# Gladius CLI

Command line interface to control the node daemon.

## Installation

#### Node.js

Node.js provides a general installation guide [here](https://nodejs.org/en/download/package-manager/) but we will walk through the installation for Windows, Ubuntu, and macOS.

We based this application off of the latest branch (9.9.0) at the time of this writing.

Here are some shortcuts to commands

* Windows
  * Download Installer, [here](https://nodejs.org/en/#download)
  * Select the latest, 9.9.0+
* Ubuntu
  * `curl -sL https://deb.nodesource.com/setup_9.x | sudo -E bash -`
  * `sudo apt-get install -y nodejs`
  * Change Global Installation Directory
    * Our packages requires some dependencies that require superuser access if installed in the default Ubuntu paths. We recommend changing the default installation of global node modules to `~/.npm-global` as stated in the [npm.js docs](https://docs.npmjs.com/getting-started/fixing-npm-permissions#option-two-change-npms-default-directory). We included the commands below:
      * Run `mkdir ~/.npm-global`
      * Run `npm config set prefix '~/.npm-global'`
      * Add `export PATH=~/.npm-global/bin:$PATH` to your `.profile` of `.zshrc` file
      * Run `source ~/.profile`

    * Another option is to use [NVM](https://docs.npmjs.com/getting-started/fixing-npm-permissions#option-one-reinstall-with-a-node-version-manager) to handle permissions.
* macOS
  * Install Homebrew, [instructions](https://brew.sh/)
    * `/usr/bin/ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"`
    * `brew install node`

#### Git

* Windows
  * https://gitforwindows.org
* Ubuntu
  * `apt-get install git`
* macOS
  * Comes default with mac but can also be installed via [Homebrew](https://brew.sh/) (`brew install git`)

#### Gladius CLI

  * Run `npm install -g gladius-cli`

#### Gladius Control Daemon

* Run `npm install -g gladius-control-daemon`

#### Gladius Edge Daemon

  * Run `npm install -g gladius-edge-daemon`


## Setup

#### Gladius Control Daemon

* Run `gladius-control` to start the server
  * Expected Output:
    ```
    $ gladius-control                                                                       
    Running at http://localhost:3000
    ```
  * **Leave this running in a new window for the CLI to communicate**

#### Gladius Edge Daemon

  * Run `gladius-edge` to start the server
    * Expected Output:
      ```
      $ gladius-edge                                                                       
      Running - Use "gladius-node start" to start it
      ```
    * **Leave this running in a new window for the CLI to communicate**

#### Gladius CLI

- Set up a local static IP for the machine you will be running the Gladius node on
- Forward port 8080 on your router to that machine
- Create a [new Ethereum wallet](https://medium.com/benebit/how-to-create-a-wallet-on-myetherwallet-and-metamask-e84da095d888)
- Acquire 1 Ether on the [Ropsten testnet](http://faucet.ropsten.be:3001/) (or go [here](https://blog.bankex.org/how-to-buy-ethereum-using-metamask-ccea0703daec) if you're using Metamask)
- Run `gladius-node init` and fill out the requested
information (use the same email that you applied for the beta with)

After you execute a command it will suggest the next logical command. For example, after `init` you can run `gladius-node create` to create a new Node. As of now the Node manager only supports 1 Node per user therefore if you run `gladius-node create` multiple times you will keep overwriting your current node.

The `gladius-cli` acts as an interface for a user to interact with the `gladius-control-daemon` and the `gladius-edge-daemon`. Therefore, if you want to use the CLI you must have **both** of the daemons running either in the background or on seperate terminal windows. Both daemons run servers on your machine once you start them. If you no longer want them to be running simply exit the window or stop the processes. This will stop the servers and if you want to use the CLI you'll have to start them again.


## Commands
`gladius-node <option>`

#### **init**
Only needs to be run once after installation. Saves user information **locally**. If you want to change your local user information you can run this command and it will take you through the on-boarding process again.

```
$ gladius-node init
[Gladius-Node] What's your email? (So we can contact you about the beta):  test@mail.com
[Gladius-Node] What's your first name?  Marcelo
[Gladius-Node] Short bio about why you're interested in Gladius:  I want to contribute my bandwidth to the Gladius Network!
[Gladius-Node] Please make a new ETH wallet and paste your private key (include 0x at the beginning):   
[Gladius-Node] Please enter a passphrase for your new PGP keys:   
[Gladius-Node] User profile created! You may create a node with gladius-node create
[Gladius-Node] If you'd like to change your information run gladius-node init again
```

#### **create**
Create and deploy a Node smart contract. You only need 1 per computer. If you create a new Node it will disconnect you from your previous one.

```
$ gladius-node create
[Gladius-Node] Please enter the passphrase for your PGP private key:  
[Gladius-Node] Creating Node contract, please wait for tx to complete (this might take a couple of minutes)
[Gladius-Node] Transaction: 0xe313b53a099addd8619f645ace76af2ddf9b4dae3e9c5ab307f2999cceb861a6	[Success]
[Gladius-Node] Setting Node data, please wait for tx to complete (this might take a couple of minutes)
[Gladius-Node] Transaction: 0x1f50a8fb77974c3543b977e42b96d0e5aa7d257ef6b7d5d98a38f7fcc3b95ffc	[Success]
[Gladius-Node] Node successfully created and ready to use
[Gladius-Node] Use gladius-node apply to apply to a pool
```

#### **apply**
Apply to a pool. Enter the pool address and an application with all of your data will be sent to them. This information includes your name, email, bio, ip address, and node contract address. **Do not apply to pools that you don't trust.**

```
$ gladius-node apply
[Gladius-Node] Please enter the passphrase for your PGP private key:  
[Gladius-Node] Please enter the address of the pool you want to join:   0x2E27eE682C14f02daa106E8f659Ed8235ef11332
[Gladius-Node] Transaction: 0x6f33237e886b1d89f69df375fc541cfe0f2b12dc9f22a9df65047c820555f8b9	[Success]
[Gladius-Node] Application sent to Pool!
[Gladius-Node] Use gladius-node check to check your application status
```

#### **check**
Check the status of your application to a particular pool.

```
$ gladius-node check
[Gladius-Node] Please enter the address of the pool you want to check on:   0x2E27eE682C14f02daa106E8f659Ed8235ef11332
[Gladius-Node] Pool: 0x2E27eE682C14f02daa106E8f659Ed8235ef11332	[Application Status: Pending]
[Gladius-Node] Wait until the pool manager accepts your application in order to become an edge node
```

#### **status**
Check the status of the Gladius Control Daemon and Gladius Edge Daemon
```
$ gladius-node status
[Gladius-Node] Gladius Control Daemon server is running!
[Gladius-Node] Gladius Edge Daemon running, you are an edge node!
[Gladius-Node] If you'd like to stop, run gladius-node stop
```

#### **start**
Starts the edge node networking. You can call this to become an edge node **after** you've been accepted to a pool. If you call this without being accepted to a pool first your server will be running but no one will be able to connect to your machine since you are not part of a pool.

```
$ gladius-node start
[Gladius-Node] Gladius Edge Daemon running, you are now an edge node!
[Gladius-Node] If you'd like to stop, run gladius-node stop
```

#### **stop**
Stops the edge node networking. You can call this to stop serving content.

```
$ gladius-node stop
[Gladius-Node] Gladius Edge Daemon is not running
[Gladius-Node] If you'd like to start, run gladius-node start
```

#### **gen-keys**
Generate a new pair of PGP keys and a new passphrase. This happens during the **init** process but you can run this again if you forget your passphrase or want to generate new keys for any reason. After this you should do `gladius-node update-node` to update the information on your node contract with your new PGP keys

```
$ gladius-node gen-keys   
[Gladius-Node] Please enter a passphrase for your new PGP keys:   
[Gladius-Node] New PGP keys generated
[Gladius-Node] Please run gladius-node update-node to update the information on your node contract

```
#### **update-node**
Overwrites your current node information (that was set upon initial creation of your node contract) with your current user data

```
$ gladius-node update-node
[Gladius-Node] Please enter the passphrase for your PGP private key:  
[Gladius-Node] Transaction: 0xdeee16ac02f0e2080d91fdb3dd982dd78350b664eddf553d0ba902dbb54c0178	[Success]
[Gladius-Node] Node information successfully update

```

#### **settings**
Displays the information that the **gladius-control-daemon** is using

```
{ running: true,
  privateKey: '0x1234567890123456789012345678901234567890123456789012345678901234',
  address: '0xE9F75E329292758c2a77f30967304cD749a88837',
  marketAddress: '0x0cd8d142238554acb52b17c9243baf6938ee3214',
  nodeFactoryAddress: '0xfb834903bcdc3ab0a2409629e3c9303e6c567a40',
  providerUrl: 'https://ropsten.infura.io/tjqLYxxGIUp0NylVCiWw',
  endpoints: { start: 'http://localhost:3000/api/settings/start' }}

```

#### **reset**
Manually wipe your local user profile. If you want to purge your user profile from your local machine run this.

```
$ gladius-node reset
[Gladius-Node] User data has been reset
```

#### **--help**
Brings up the help menu.

## Notes and warnings
**Warning!** This is the beta implementation of the Gladius node. **Please** create a new ETH wallet before using this. Pools that you apply to will receive your name, email, bio, ip address, and node contract address. Therefore, **do not** apply to pools that you do not trust. At this time Gladius is running 1 official pool that is closed to the public and will not be advertised. If you choose to join an independent pool not run by Gladius, that is at your own risk. Gladius is not responsible for your data getting into the wrong hands. Use at your own risk.
