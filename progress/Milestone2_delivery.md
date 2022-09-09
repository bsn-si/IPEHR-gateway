v# Milestone Delivery :mailbox:

* **Application Document:** [ipEHR application](https://github.com/filecoin-project/devgrants/issues/418)
* **Milestone Number:** 2

**Context**

In this milestone we've developed the functionality of access to the data stored in Filecoin.

**Deliverables**

| Number                | Deliverable                              | Link                                                                                                                                  | Notes                                                                                                                                                                                                                                                                                                                  |
|-----------------------|------------------------------------------|---------------------------------------------------------------------------------------------------------------------------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| :heavy_check_mark:    | License                                  | [LICENSE](https://github.com/bsn-si/IPEHR-gateway/blob/develop/LICENSE)                                                               | Apache 2.0 license                                                                                                                                                                                                                                                                                                     |
| :heavy_check_mark:    | Testing Guide                            | [Readme.md](https://github.com/bsn-si/IPEHR-gateway/blob/develop/README.md#how-to)                                                    | The "How To" guide is supplemented with all new features developed in this milestone                                                                                                                                                                                                                                                                                                        |
| 1. :heavy_check_mark: | Integrate Filecoin as a storage for MH-ORM database | See [Filecoin_storage](https://github.com/bsn-si/IPEHR-gateway/blob/develop/progress/Milestone_2/1_Filecoin_integration.md) page | At IPEHR we use a two-tiered document storage system: 1. The IPFS distributed file system for quick access. 2. The Filecoin storage network for long-term guaranteed storage of EHR documents. Filecoin storage deals are made for a fixed period of time. Usually from 180 days. As the deadline approaches, the deal must be extended. This functionality will be added in the next phases. | 
| 2. :heavy_check_mark: | Integrate a storage for MH-ORM indexes | See [Indexes_storage](https://github.com/bsn-si/IPEHR-gateway/blob/develop/progress/Milestone_2/2_Indexes_storage.md) page | An EVM smart contract has been developed to store EHR document indexes. | 
| 3. :heavy_check_mark: | An algorithm of data re-encryption while changing access rights | See [Revoking_access](https://github.com/bsn-si/IPEHR-gateway/blob/develop/progress/Milestone_2/3_Revoking_access.md) page | ВОПРОСЫ!!!! Each document is encrypted with a unique access key. User access to documents is controlled by the docAccess index. On revocation of access to a document, a corresponding change will be made to the docAccess index on the smart contract, and the user will not be able to get an access key to the document from that moment. |
| 4. :heavy_check_mark: | Records chains | See [Chains of records](https://github.com/bsn-si/IPEHR-gateway/blob/develop/progress/Milestone_2/4_Chains_of_records.md) page | When a new EHR document is added, the digital signature of the user creating the document is sent along with the document. The signature is stored with the document. This allows you to authorize a request to create a document. |
| 5. :heavy_check_mark: | Performance tests and optimisation | URL to be here | Milestone deliverables testing examples with video guides | 
| 6. :heavy_check_mark: | Payment logic | URL to be here | Smt to be here | 

# Performance tests

Test instructions and performance tests videos to be here 📹 
