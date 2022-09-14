![smart-contract](https://user-images.githubusercontent.com/8058268/190085702-6edf9437-1273-4db3-a7c9-414f66afe823.svg)

A smart contract has been developed to store EHR document indexes.

The smart contract is written in Solidity language and can be deployed in any [EVM](https://ethereum.org/en/developers/docs/evm/) compatible blockchain.

### Document meta-information
- CID for retrieving a document in IPFS. Is the hash sum of the document and allows to check the integrity.
- dealCid and minerAddress for getting the document in Filecoin
- `document_uid` (encrypted)
- document status
- version
- creation time

### The main functions of a smart contract are

1. Search for `ehr_id` by user ID
2. Obtaining a list of documents related to the specified `ehr_id`.
3. Getting meta-information of the document by `ehr_id` and `document_uid`.
4. Getting an access key to the document. (Access keys are encrypted with the keys of users who have access)
5. Management of access to the document (*work in progress*)
6. Document search using [AQL](https://specifications.openehr.org/releases/QUERY/latest/AQL.html) queries. (*work in progress*)

User and EHR document data in the contract is stored encrypted and prevents unauthorized persons from accessing private information.

For development and testing purposes, [Goerli Testnet](https://goerli.net/) is used

## Implementation

The contract code is located in the repository: [https://github.com/bsn-si/IPEHR-blockchain-indexes](https://github.com/bsn-si/IPEHR-blockchain-indexes)

Running tests:

```
npx hardhat test
```

Address of the current version of the deployed contract: [https://goerli.etherscan.io/address/0x90346f14e3d22bff62415707928486a42282b19f](https://goerli.etherscan.io/address/0x90346f14e3d22bff62415707928486a42282b19f)
