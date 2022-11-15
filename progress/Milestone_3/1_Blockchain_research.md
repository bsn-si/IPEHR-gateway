## Research of available blockchains supporting EVM-based smart contracts. EVM will be at the basis of our access management system


### Current state

To deploy smart contracts our solution uses [EVM](https://ethereum.org/en/developers/docs/evm/).

Pros:

- it works with [Solidity](https://docs.soliditylang.org/en/v0.8.17/),
- it is a popular solution with a big community,
- our team has a profound experience working with it.

Cons:

- it is less convenient due to the need for two types of coins since we are working with Filecoin and Ethereum networks,
- It imposes additional costs.

### Possible solution

To eliminate these drawbacks, our team is looking forward to using [FVM](https://fvm.filecoin.io/).

It is a promising solution since EVM compatibility is announced in Milestone 2 but it can't be used at the moment because:

- the libraries used in the project to interact with the EVM contract proved to be incompatible with FEVM and require modification.
- due to a list of limitations, our team has failed to deploy a smart contract into FVM right now.

### Possible future steps

Our team is considering the possibility of developing tools that will help us and other developers to deploy smart contracts to FVM.
