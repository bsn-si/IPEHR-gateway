### IPEHR Oracle
A subproject that contains all contracts and scripts for publishing and interacting with the oracle.

### About
We've written a set of contracts for providing access to statistics on the blockchain, based on Ether and Chainlink, along with storage contracts and example contracts.

We have several types of statistical data delivery: direct delivery on request and request of data on a schedule.

- Direct delivery is a task in Chainlink that accepts requests from outside for a small fee in LINK tokens, listens to the operator's contract and returns the result from the statistics server.
- Scheduled delivery is a task in Chainlink that updates the storage contract with fresh data according to a schedule. Other contracts can make shareware requests to contract data. (An example of a consumer contract is also available).

Also, a set of scripts was written for the provided contracts to simplify interaction and testing, with which you can publish and call contracts as well as view some of the Chainlink statuses and replenish its balance.

### How To
Several manuals are available for working with contracts and oracle.

- [Setup chainlink node](chainlink/README.md)
- [Contracts with examples](contracts/README.md)
- [How to use scripts for interaction](scripts/README.md)
