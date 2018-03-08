# Gladius CLI

Command line interface to control the node daemon.

**This is a work in progress, not all of the below methods will work**

## Installation
First install nodejs and npm from [here](https://nodejs.org/en/).

Navigate to the project and run the following command to install the Gladius
CLI.

`npm install -g .`

Finally, install the daemon using either [docker](https://gladius.io) or
[directly](https://gladius.io)

## Setup

Set up a local static IP for the machine you will be running the Gladius node on
, and forward port 443 and 80 on your router to this machine.

_TODO: Talk about configuring a hostname for SSL in future_

Run `gladius-node init` to initialize your Gladius node. Fill out the requested
information and make sure you use the same private key as you used when you
registered for the beta. This will ensure that you can sign up without a hitch.

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
of the "config.json" file that dictates parameters like ports and the IP address
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
