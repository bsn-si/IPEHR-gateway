## Working with documents in IPEHR

Document storage deals on the Filecoin network require payment with FIL tokens.

The network's native token is used to pay for transactions when writing information to smart contracts in blockchains. 

We are currently using [Goerli](https://goerli.net/), the testnet of the Ethereum network, in the development of IPEHR. The native token of the network is test Ether. 

The source that generates the EHR documents is the HMS systems used by medical organizations.

The `ehrSystemID` field, which corresponds to the [system_id](https://specifications.openehr.org/releases/BASE/latest/architecture_overview.html#_system_identity) field in the openEHR specification, is used to identify medical organizations.

Registration on the IPEHR-gateway is required before starting. Registration creates a dedicated account and generates two addresses for deposit of FIL and Ether tokens. Next, token deposit to the received addresses must be performed.

Further, when working with the IPEHR system, payment for transactions is made from the organization's account within the available balance.

At any time, the organization can check the available token balance using Filecoin and Goerli browsers, or via the IPEHR-gateway API.

## EHR Document Rights Management

The document access rights management system will be implemented as a blockchain-based smart contract.

Users (patients) will be able to manage access to their documents through a special application. To pay for transactions, a wallet with tokens from the blockchain in which the smart contract will be deployed is required.