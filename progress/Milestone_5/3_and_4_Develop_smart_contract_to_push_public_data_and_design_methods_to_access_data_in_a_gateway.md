We developed a set of contracts for providing access to statistics on the blockchain, based on ether and chainlink, along with storage contracts, and example contracts.

We have several types of statistical data delivery:

-   Direct delivery. Is implemented as a task in Chainlink that receives requests from outside by listening to the Oracle contract. When the Consumer contract sends a request for the statistical data to Oracle, the job collects a small fee in Chainlink tokens and returns the result from the statistics server. For this case, we have an open API for statistics and documented schema of job for Chainlink. The contract for this request is not a library but a sample contract in which we request statistical data via Oracle from the certain job of Chainlink

A video example of direct deivery:
[![video-m5-1](https://user-images.githubusercontent.com/98888366/209851585-3ecf965f-0f71-49fe-a35e-25b4e3641c8b.png)](https://media.bsn.si/ipehr/video-m5-1.mp4)
    
-   Scheduled delivery. It consists of two contracts and the schema of a Chainlink job. This job automatically requests statistical data within the specified interval and sends it to a storage contract. All other external contracts can request statistical data from the storage contract. The implementation includes two contracts. The first is the storage itself. We publish it and pay for its updates. The second is the contract of a Consumer. It is a sample of a simple contract that requests statistical data from storage.
    
A video example of a scheduled delivery:
[![video-m5-2](https://user-images.githubusercontent.com/98888366/209851873-ffe97a94-bc75-43fe-baa2-eba73a36744c.png)](https://media.bsn.si/ipehr/video-m5-2.mp4)

Also, we developed a set of cli scripts that allow automating the process. They allow to:

-   Publish a Chainlink token
    
-   Publish Oracle contract
    
-   Publish and call a contract for a direct request from Chainlink via an operator
    
-   Publish and call via Cron contracts for storage and Consumer
    
-   View Chainlink statuses
    
-   Replenish Chainlink balance
    

A [link](https://github.com/bsn-si/IPEHR-stat/tree/main/oracle) for more information and manuals. 
