## Design of the chains of records within the storage to ensure the data integrity and authenticity.

When a new EHR document is added, the digital signature of the user creating the document is sent along with the document. The signature is stored with the document. This allows you to authorize a request to create a document.

Before saving, the EHR document is encrypted with [ChaCha20-Poly1305](https://en.wikipedia.org/wiki/ChaCha20-Poly1305), which includes a message authentication code (MAC). This mechanism ensures the integrity of the document when it is decrypted.

When the document is saved in the IPFS network, the [CID](https://docs.ipfs.tech/concepts/content-addressing/#what-is-a-cid) of the file is calculated, which is essentially the hash sum of the document and provides the assurance that the file has not been changed.

The Filecoin decentralized file storage stores the document using the CID from IPFS.

The [IPEHR smart contract](https://github.com/bsn-si/IPEHR-blockchain-indexes/blob/develop/contracts/EhrIndexer.sol) records meta-information of the EHR document containing the CID and an encrypted identifier. This allows us to record the fact and time of creation of a particular document.

The smart contract also implements document versioning logic, allowing to trace the chain from creation to the final version of the document.
