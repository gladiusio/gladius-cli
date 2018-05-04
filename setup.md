# How to setup your config file

In order to make some transactions on the blockchain you must have ether in a wallet that is then used to pay for the transaction. This is the case with some aspects of the Gladius Network such as deploying a new Node smart contract. This will show you how to import your wallet and encryption keys in order to make secure transactions for the Gladius Network.

Make a new file called `env.toml` in the root of `gladius-cli`. Then copy this and fill in the missing parameters. Feel to keep the default addresses for the smart contracts.

```
# Example: env.toml

[environment]
  marketAddress = "0x9f5bd0fd43f2ee2e1ef406cec5936768ab8c587d"
  nodeFactoryAddress = "0xb4d68308e81687698cbeb4817cc91a7ede9d8a54"
  passphrase = "PGP PASSPHRASE HERE" # fill in
  pgpKey = """
  PRIVATE PGP KEY HERE""" # fill in
  poolAddress = "0xC88a29cf8F0Baf07fc822DEaA24b383Fc30f27e4"
  privateKey = "ETH WALLET PVT KEY. PLEASE MAKE A NEW WALLET" # fill in
  provider = "https://ropsten.infura.io/APIKEY" # fill in

[node]
```

Here are some resources for generating these fields:

- [TOML spec](https://github.com/toml-lang/toml)
- [Generate PGP keys](https://pgpkeygen.com/)
- [Create new ETH wallet](https://kb.myetherwallet.com/getting-started/creating-a-new-wallet-on-myetherwallet.html)
- [Infura](https://infura.io/signup)
