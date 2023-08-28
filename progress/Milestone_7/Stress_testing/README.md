# Stress testing

## Used hardware and software

[opentelemetry.io](https://opentelemetry.io) tools were used to measure the speed of request processing and collect metrics  
[k6.io](https://k6.io) tools was used for load testing   
Были предварительно разработаны [скрипты]() для тестирования API ipEHR Gateway.

Testing was performed on a laptop with the following specifications:

**MacBook Pro (16-inch, 2019)**  
**Processor:**   2.6GHz 6‑core Intel Core i7, with 12MB shared L3 cache  
**Memory:** 32GB of 2666MHz DDR4  
**Storage:** 512GB SSD

## Test scenario

A set of standard operations was used as a typical scenario:

- patient registration
- authentication (obtaining an access token)
- create a new electronic health record (EHR) for the patient
- retrieve patient and EHR information
- request the created EHR document
- logout

## Blockchain

Test runs were performed on 2 EVM-enabled networks: Calibration and Sepolia.

The ipEHR smart contracts have been pre-deployed:

| Contract                                                                                                       | Sepolia | Calibration |
| -------------------------------------------------------------------------------------------------------------- | -------------------------------------------- | -------------------------------------------- |
| [EhrIndexer](https://github.com/bsn-si/IPEHR-blockchain-indexes/blob/develop/contracts/EhrIndexer.sol)         | `0x28668f2C40b7FB018F109cC5b2A14C951ff78c91` | `0x9944D37bFeC481868baad7b6E05b76Db01cA0865` |
| [Users](https://github.com/bsn-si/IPEHR-blockchain-indexes/blob/develop/contracts/Users.sol)                   | `0x1B89Ec21E1A0e3E85038daEdf5cA3feF4F087957` | `0xe95dB24EA185c7a7D4ED6e8D20Caaa4cCb852AF4` |
| [AccessStore](https://github.com/bsn-si/IPEHR-blockchain-indexes/blob/develop/contracts/AccessStore.sol)       | `0x380ce0529022b4DeEE67A601D1FA1be3B3d2D781` | `0x946e8BB742AAe895F9D369Fa44eF41414607A0CA` |
| [DataStore](https://github.com/bsn-si/IPEHR-blockchain-indexes/blob/develop/contracts/DataStore.sol)           | `0x906e66Cb9937c1891aC63576B8C082eedB37DF1e` | `0x0d223203D54c453Fb70A1462C3901AC59CE1F103` |

## Bottlenecks identified

In the process of testing, a number of problems were identified and eliminated

1. The incremental Nonce was used for transaction signature formatting. This prevented parallel calls of smart contract functions. To solve the problem, instead of incremental Nonce, `deadline` - a period of time during which a transaction is valid - is now used;

2.  The ipehr-gateway used the `PendingNonceAt` call when making smart contract calls, but this also did not allow parallel calls. To solve the problem, the control of the `Nonce` value was moved to the gateway;

3. The module for tracking the status of current transactions in the blockchain was redesigned for more efficient operation;

4. The values of timeouts for waiting for HTTP requests, waiting for processing blockchain transactions, and test timeouts were selected based on the peculiarities of the used testnets.

## Results

### Calibration

[ipehr\_gateway\_log]()  
[k6\_test\_report]()  
[requests_statistics]()

### Sepolia

[ipehr\_gateway\_log]()  
[k6\_test\_report]()  
[requests_statistics]()

*Summary table:*


|      Metric         |   Calibration   |  Sepolia  |
|:-------------------:|:---------------:|:---------:|
| block_time          | 30 sec | 12 sec   |
| iterations          | 100    | 100      |
| threads             | 20     | 20       |
| time_total          | 39m14s | 23m53s   |
| http_reqs           | 4180   | 3545     |
| iteration_time      | 14m13s | 6m37s |
| user\_register_time | 2m20s | 3m57s |
| user\_login_time    | 112ms | 139ms |
| ehr\_create_time    | 4m14s | 3m34s |
| ehr\_get_time       | 119ms | 107ms |
| user\_logout_time   | 124ms | 130ms |
| user\_create\_gas\_used | 258304636 | 943537 |
| ehr\_create\_gas\_used | 1094164389 | 3168824 |

## Conclusions

In the process of load testing, optimal modes of application operation were determined when using two testnet's - [Calibration](https://calibration.filfox.info/en) and [Sepolia](https://sepolia.etherscan.io/). Also, some problems were identified and eliminated during the operation of the ipEHR system during parallel processing of requests. 

As a result of our test runs, we conclude that the throughput of the system depends mainly on the performance of the blockchain in which the smart contracts are deployed. 

Using the Calibration testnet, 100 patients were able to register and create a document in 39 minutes, which equates to a throughput of 2.56 users per minute.

For the Sepolia network, the throughput is 4.16 users per minute. 

To increase throughput, smart contracts can be further optimized by increasing the use of multicall transactions, or higher performing blockchains can be used.