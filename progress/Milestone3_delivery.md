# Milestone Delivery :mailbox:

* **Application Document:** [ipEHR application](https://github.com/filecoin-project/devgrants/issues/418)
* **Milestone Number:** 3

**Context**

In this milestone we've developed the functionality to manage access rights on a blockchain level.

**Deliverables**

| Number                | Deliverable                              | Link                                                                                                                                  | Notes                                                                                                                                                                                                                                                                                                                  |
|-----------------------|------------------------------------------|---------------------------------------------------------------------------------------------------------------------------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| :heavy_check_mark:    | License | [LICENSE](https://github.com/bsn-si/IPEHR-gateway/blob/develop/LICENSE) | Apache 2.0 license |
| :heavy_check_mark:    | Testing Guide | [Readme.md](https://github.com/bsn-si/IPEHR-gateway/blob/develop/README.md#how-to) | The "How To" guide is supplemented with all new features developed in this milestone |
| 1. :heavy_check_mark: | Research of available blockchains supporting EVM-based smart contracts | See [IPEHR smart contract](https://github.com/bsn-si/IPEHR-blockchain-indexes) page | The test contract is deployed to the Goerli testnet. For the contract fork to the FVM network a golang client is needed to correctly interact with FEVM JSON-RPC API. | 
| 2. :heavy_check_mark: | Development of the user account catalogue and the user identification mechanism | See [Users_identity](https://github.com/bsn-si/IPEHR-gateway/blob/develop/progress/Milestone_3/2_Users_identity.md) page | Users and user groups are stored in the [IPEHR smart contract](https://github.com/bsn-si/IPEHR-blockchain-indexes). The authorization of requests to the IPEHR gateway API is done via the JWT token. | 
| 3. :heavy_check_mark: | Embedding access rights management | See [Docs_access_management](https://github.com/bsn-si/IPEHR-gateway/blob/develop/progress/Milestone_3/3_Docs_access_mgmt.md) page | Documents can be grouped by arbitrary criteria. From medical classification to geographical location. Access to documents is managed according to a particular access matrix. |
| 4. :heavy_check_mark: | Design and development of the access keys storage | See [Access_store](https://github.com/bsn-si/IPEHR-gateway/blob/develop/progress/Milestone_3/4_Access_store.md) page | Access Key Storage is a hash table that is located in the smart contract. The key of the table is a special 32 byte identification number. |
| 5. :heavy_check_mark: | Design and development of a smart contract API | See [Iface](https://github.com/bsn-si/IPEHR-gateway/blob/develop/progress/Milestone_3/5_Iface_to_contract.md) | To interact with a smart contract containing a repository of users, documents and access rights, the indexer package was developed using the Go Ethereum library. To prevent the execution of repeated transactions, the nonce mechanism is used in combination with ECDSA signature verification. | 

### Project introduction:

[![Watch the video](https://media.bsn.si/ipehr/logo_intro.jpg)](https://www.youtube.com/watch?v=nJFA5W4qoEw)

# Workflow example

