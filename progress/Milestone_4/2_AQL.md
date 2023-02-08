### Implementation AQL according to the openEHR protocol.

Client applications use for data queries, this way platform can be used by the 3rd party applications supporting the standard.

In the process of implementing the ability to execute AQL queries, the following tasks were performed:

1. Index data repository developed: [https://github.com/bsn-si/IPEHR-gateway/tree/develop/src/pkg/storage/treeindex](https://github.com/bsn-si/IPEHR-gateway/tree/develop/src/pkg/storage/treeindex)
2. Smart contract for storing index data: [https://github.com/bsn-si/IPEHR-blockchain-indexes/blob/develop/contracts/DataStore.sol](https://github.com/bsn-si/IPEHR-blockchain-indexes/blob/develop/contracts/DataStore.sol)
3. AQL querier(*) - library for parsing an AQL query: [https://github.com/bsn-si/IPEHR-gateway/tree/develop/src/pkg/aqlprocessor](https://github.com/bsn-si/IPEHR-gateway/tree/develop/src/pkg/aqlprocessor)
4. AQL executor(*) - service that searches the index tree for data according to an AQL query: [https://github.com/bsn-si/IPEHR-gateway/tree/develop/src/pkg/aqlquerier](https://github.com/bsn-si/IPEHR-gateway/tree/develop/src/pkg/aqlquerier)

\* - as the project has MVP status, only basic AQL functionality has been implemented
