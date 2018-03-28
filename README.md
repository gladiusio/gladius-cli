# Gladius CLI

Command line interface to control the node daemon.

**This is a work in progress, not all of the below methods will work**

## Installation
- First install nodejs and npm from [here](https://nodejs.org/en/)
- Install the `gladius-cli` through npm with `npm install gladius-cli -g`

## Setup
- Set up a local static IP for the machine you will be running the Gladius node on
- Forward port 443 and 80 on your router to that machine
- Create a [new Ethereum wallet](https://medium.com/benebit/how-to-create-a-wallet-on-myetherwallet-and-metamask-e84da095d888)
- Acquire 1 Ether on the [Ropsten testnet](http://faucet.ropsten.be:3001/) (or go [here](https://blog.bankex.org/how-to-buy-ethereum-using-metamask-ccea0703daec) if you're using Metamask)
- Run `gladius-node init` and fill out the requested
information (use the same email that you applied for the beta with)

## Commands (`gladius-node <option>`)
### `init`
Only needs to be run once after installation. Saves user information **locally**. If you want to change your local user information you can run this command and it will take you through the on-boarding process again.

### `create`
Create and deploy a Node smart contract. You only need 1 per computer. If you create a new Node it will disconnect you from your previous one.

### `apply`
Apply to a pool. Enter the pool address and an application with all of your data will be sent to them. This information includes your name, email, bio, ip address, and node contract address. **Do not apply to pools that you don't trust.**

### `check`
Check the status of your application to a particular pool.

### `status`
WIP

### `start`
Starts the edge node networking server. You can call this to become an edge node **after** you've been accepted to a pool.

### `stop`
Stops the edge node networking server. You can call this to stop serving content.

### `gen-keys`
Generate a new pair of PGP keys and a new passphrase. This happens during the `init` process but you can run this again if you forget your passphrase or want to generate new keys for any reason.

### `settings`
Displays the information that the `gladius-control-daemon` has received from your Node contract.

### `reset`
Manually wipe your local user profile. If you want to purge your user profile from your local machine run this.

### `--help`
Brings up the help menu.

## Notes and warnings
**Warning!** This is the beta implementation of the Gladius node. **Please** create a new ETH wallet before using this. Pools that you apply to will receive your name, email, bio, ip address, and node contract address. Therefore, **do not** apply to pools that you do not trust. At this time Gladius is running 1 official pool that is closed to the public and will not be advertised. If you choose to join an independent pool not run by Gladius, that is at your own risk. Gladius is not responsible for your data getting into the wrong hands. Use at your own risk.
