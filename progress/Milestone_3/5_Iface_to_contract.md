## Design and development of an interface to interact with smart contract. Development of the dedicated library.

To interact with a smart contract containing a repository of users, documents and access rights, the `indexer` package was developed using the [Go Ethereum](https://geth.ethereum.org/) library, which is popular among dApps developers.

To prevent the execution of repeated transactions, the nonce mechanism is used in combination with ECDSA signature verification.

Source code of package: [https://github.com/bsn-si/IPEHR-gateway/tree/develop/src/pkg/indexer](https://github.com/bsn-si/IPEHR-gateway/tree/develop/src/pkg/indexer)
