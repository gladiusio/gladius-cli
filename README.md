# Gladius CLI

Command line interface to control the node daemon.

**This is a work in progress, not all of the below methods will work**

## Installation
- First install nodejs and npm from [here](https://nodejs.org/en/)
- Clone the latest version of `gladius-cli`

## Setup

- Set up a local static IP for the machine you will be running the Gladius node on
- Forward port 443 and 80 on your router to that machine

_TODO: Talk about configuring a hostname for SSL in future_

- Create a [new PGP key](https://pgpkeygen.com/) aswell as a [new Ethereum wallet](https://medium.com/benebit/how-to-create-a-wallet-on-myetherwallet-and-metamask-e84da095d888)
- Paste your **new** private keys in the `pgpKey.txt` and `pvtKey.txt` files located at `/gladius-cli/keys`
- Acquire 1 Ether on the [Ropsten testnet](http://faucet.ropsten.be:3001/) (or go [here](https://blog.bankex.org/how-to-buy-ethereum-using-metamask-ccea0703daec) if you're using Metamask)
-
- Run `gladius-node init` and fill out the requested
information (use the same email that you applied for the beta with)

Once you have filled out your information (and there are no errors) you can run
the command `gladius-node join-pool` to inform the Gladius daemon that this
information is correct and you would like to join the beta pool. Your data will
then be encrypted using the pool's public key, and put on the Ethereum
blockchain. This may take some time to complete, so check your status with the
command `gladius-node check-join` this can return a few different states
(pending, denied, or accepted).

Once you have been accepted to the pool, you can run
the command `gladius-node start` to inform the Gladius daemon (it must be
running and installed) that you would like to start accepting requests for
content.

To inform the daemon that you would like to stop accepting requests, run
`gladius-node stop`. This will leave the daemon running, but will stop
networking and transactions from taking place.

## Configuration
To configure features of the Gladius node software that would normally not be
used by the average user, run `gladius-node config-location` to get the location
of the "config.js" file that dictates parameters like ports and the IP address
of the daemon.   

## Notes and warnings
You can see a full list of the commands available (some may not be fully
 functional) by running `gladius-node --help`


**Warning!** This is the beta implementation of the Gladius node, and as such it
means there are some differences between how this functions and how future
versions will. **Do not add private keys here that store more Ether or GLA (or
any other cryptocurrency) than you are willing to lose.**

In this version there is no way to specify the pool you would
like to join, as the only pool available is run by Gladius. This version may
also be less optimized than future versions, so run this on machines that you
don't normally demand too much from.
