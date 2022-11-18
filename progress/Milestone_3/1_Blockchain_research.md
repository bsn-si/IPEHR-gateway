## Research of available blockchains supporting EVM-based smart contracts. EVM will be at the basis of our access management system


### Current state

To deploy smart contracts our solution uses [EVM](https://ethereum.org/en/developers/docs/evm/). To interact with the EVM contract we use the [Go-Ethereum library](https://github.com/ethereum/go-ethereum).
We are currently using [Goerli](https://goerli.net/), the testnet of the Ethereum network, in the development of IPEHR. The native token of the network is test Ether.

Current contracts may be found [here](https://goerli.etherscan.io/address/0x3fcEa11C70A205CF2610807b0F0cdA774079fAf3).

Pros:

- it works with [Solidity](https://docs.soliditylang.org/en/v0.8.17/),
- it is a popular solution with a big community,
- our team has a profound experience working with it.

Cons:

- One of the notable drawbacks of placing the ipEHR smart-contract on the Ethereum network is the need to use assets in separate networks to ensure interaction. Since we are working with Filecoin and Ethereum networks, it imposes additional costs and management complications.

### Possible solution

To eliminate described drawbacks, our team is looking forward to deploying on [FVM](https://fvm.filecoin.io/). This task will require, first of all, creation of a library similar to Go-Ethereum in order to ensure correct interaction with the contract in the FEVM.

Our team is considering possibility of developing tools that will help us and other developers to deploy smart contracts to FVM.
