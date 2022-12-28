# Milestone Delivery :mailbox:

* **Application Document:** [ipEHR application](https://github.com/filecoin-project/devgrants/issues/418)
* **Milestone Number:** 5

**Context**

In this milestone we've developed a data feed of EHR stats data to blockchain through Chainlink.

**Deliverables**

| Number                | Deliverable                              | Link                                                                                                                                  | Notes                                                                                                                                                                                                                                                                                                                  |
|-----------------------|------------------------------------------|---------------------------------------------------------------------------------------------------------------------------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| :heavy_check_mark:    | License | [LICENSE](https://github.com/bsn-si/IPEHR-stat/blob/main/LICENSE) | Apache 2.0 license |
| :heavy_check_mark:    | Testing Guide | [Readme.md](https://github.com/bsn-si/IPEHR-stat/tree/main#local-deployment) | The "How To" guide is supplemented with usage examples videos |
| 1. :heavy_check_mark: | Develop public data gatherer | See [IPEHR stat API](https://stat.ipehr.org/swagger/index.html) | EHR statistics can be collected by the service through transaction analysis of contracts IPEHR-blockchain-indexes, periodic direct invocation of contract methods IPEHR-blockchain-indexes and making AQL queries to IPEHR-gateway. | 
| 2. :heavy_check_mark: | Data feed (data provider) node development | See [Oracle/Chainlink](https://github.com/bsn-si/IPEHR-stat/blob/main/oracle/README.md) page | A simple quick guide on how to deploy a local chainlink node for development and testing. After deploying ipEHR smart contracts to FEVM testnet we will correct these instructions accordingly (Jan-Feb '23). | 
| 3. :heavy_check_mark: | Develop smart contract to push public data | See [Oracle/Contracts](https://github.com/bsn-si/IPEHR-stat/tree/main/oracle/contracts) page | Smart contracts for direct & scheduled delivery. |
| 4. :heavy_check_mark: | Design methods to access data in a gateway | See [ipEHR stats](https://github.com/bsn-si/IPEHR-stat#ipehr-stat) page | For demonstration purposes, the following metrics are implemented: number of patients registered in the system over all time; number of patients logged in the system for a specified period; number of EHR documents registered in the system; number of EHR documents registered in the system for a given period. |

### Project introduction:

[![Watch the video](https://media.bsn.si/ipehr/logo_intro.jpg)](https://www.youtube.com/watch?v=nJFA5W4qoEw)

# Workflow example

We have developed a service that allows making statistical data available to be collected and processed by the IPEHR system. The service implements an open API with specified metrics. The data is collected and processed by accessing [IPEHR-blockchain-indexes](https://github.com/bsn-si/IPEHR-blockchain-indexes) smart contracts.

For demonstration purposes, the following metrics are implemented:

-   number of patients registered in the system overall time;
-   number of patients logged in the system for a specified month;
-   number of EHR documents registered in the system for all time;
-   number of EHR documents registered in the system for a given month.

Please note that the following metrics are samples that are made to demonstrate the possibilities of an ipEHR MVP. At this point, we are providing access only to statistical data. With the future development of the project when it will be monetized with its own tokens, we will be able to provide access to a much wider range of data and contracts.

Two types of stats delivery:

-   Direct delivery. Is implemented as a task in Chainlink that receives requests from outside by listening to the Oracle contract. When the Consumer contract sends a request for the statistical data to Oracle, the job collects a small fee in Chainlink tokens and returns the result from the statistics server. For this case, we have an open API for statistics and documented schema of job for Chainlink. The contract for this request is not a library but a sample contract in which we request statistical data via Oracle from the certain job of Chainlink

📹 A video example of direct deivery:
[![video-m5-1](https://user-images.githubusercontent.com/98888366/209851585-3ecf965f-0f71-49fe-a35e-25b4e3641c8b.png)](https://media.bsn.si/ipehr/video-m5-1.mp4)

-   Scheduled delivery. It consists of two contracts and the schema of a Chainlink job. This job automatically requests statistical data within the specified interval and sends it to a storage contract. All other external contracts can request statistical data from the storage contract. The implementation includes two contracts. The first is the storage itself. We publish it and pay for its updates. The second is the contract of a Consumer. It is a sample of a simple contract that requests statistical data from storage.

📹 A video example of a scheduled delivery:
[![video-m5-2](https://user-images.githubusercontent.com/98888366/209851873-ffe97a94-bc75-43fe-baa2-eba73a36744c.png)](https://media.bsn.si/ipehr/video-m5-2.mp4).
