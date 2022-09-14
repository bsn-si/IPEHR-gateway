![1](https://user-images.githubusercontent.com/8058268/189339273-3e49a10f-f7b1-46ef-8fa3-514e0d4d9edf.svg)

### At IPEHR we use a two-tiered document storage system.

1. The IPFS distributed file system is used for quick access to documents.
2. For long-term guaranteed storage of EHR documents, we use a decentralized Filecoin storage network.

### Stages of processing of newly stored documents:

- validation
- compression and encryption
- saving to IPFS
- saving in a [smart contract](https://github.com/bsn-si/IPEHR-blockchain-indexes) on the blockchain the meta-information about the document 
- saving meaningful information from a document into a [smart contract](https://github.com/bsn-si/IPEHR-blockchain-indexes) on the blockchain using homomorphic encryption for search and filtering purposes
- saving the document to Filecoin:
	- searching for a suitable miner to conclude a transaction using the [FilRep ](https://filrep.io) service, taking into account the criteria: overall rating, minimum/maximum data size, statistics on previous transactions, price, etc.
	- signing a deal with the miner in the Filecoin network
	- transaction status tracking

Filecoin storage deals are made for a fixed period of time. Usually from 180 days.
As the deadline approaches, the deal must be extended. This functionality will be added in the next phases.

To interact with the Filecoin network, a connection to [Lotus](https://lotus.filecoin.io), which provides JSON-RPC API, is used.

### Stages of EHR document request processing:

- checking access rights to the document
- Search for meta-information about the document in the smart contract
- searching for the document in IPFS by CID
- if the document is not found in IPFS, the recovery procedure from Filecoin is launched
	- loading the document by DealCID and minerAddress
	- saving to IPFS
- the document is decrypted and transmitted to the user 

## Implementation

Client code for connecting to Lotus: [https://github.com/bsn-si/IPEHR-gateway/tree/develop/src/pkg/storage/filecoin](https://github.com/bsn-si/IPEHR-gateway/tree/develop/src/pkg/storage/filecoin)

Running tests to demonstrate the work of the packages:

```
go test pkg/storage/filecoin/*
```
