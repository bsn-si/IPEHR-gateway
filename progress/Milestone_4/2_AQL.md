## Implementation AQL according to the openEHR protocol.

Client applications use for data queries, this way platform can be used by the 3rd party applications supporting the standard.

In the process of implementing the ability to execute AQL queries, the following tasks were performed:

1. Index data repository developed: [https://github.com/bsn-si/IPEHR-gateway/tree/develop/src/pkg/storage/treeindex](https://github.com/bsn-si/IPEHR-gateway/tree/develop/src/pkg/storage/treeindex)
2. Smart contract for storing index data: [https://github.com/bsn-si/IPEHR-blockchain-indexes/blob/develop/contracts/DataStore.sol](https://github.com/bsn-si/IPEHR-blockchain-indexes/blob/develop/contracts/DataStore.sol)
3. AQL querier(*) - library for parsing an AQL query: [https://github.com/bsn-si/IPEHR-gateway/tree/develop/src/pkg/aqlprocessor](https://github.com/bsn-si/IPEHR-gateway/tree/develop/src/pkg/aqlprocessor)
4. AQL executor(*) - service that searches the index tree for data according to an AQL query: [https://github.com/bsn-si/IPEHR-stat/tree/main/internal/aqlquerier](https://github.com/bsn-si/IPEHR-stat/tree/main/internal/aqlquerier)

\* - in this stage of the project we have implemented basic AQL functionality as in the list below:
|         Feature         | Implementation |
| ----------------------- |      :---:     |
| EHR documents data parsing           | + |
| EHR data Tree index design           | + |
| EHR data Tree index in-memory        | + |
| AQL requests parsing                 | + |
| AQL processing on IPEHR-gateway      | + |
|||
| Select primitives      | + |
| Select values      | + |
| Select values with WHERE      | + |
| Selevt values with WHERE EXISTS      | + |
| Select values with WHERE AND      | + |
| Select values with WHERE OR      | + |
| Select values with WHERE NOT      | + |
| Select values with WHERE AND (OR)      | + |
| Select values and int value      | + |
| Select multiple columns      | + |
| Select with filter by EHR id      | + |
| Select with filter by EHR id and obs. version      | + |
| Select with $parametes      | + |
| 		|
| PROCESSOR		|
| Select field      | + |
| Select field with path_predicate      | + |
| Select field[...]/value      | + |
| Select o/field[...]/value1/value2 with alisas      | + |
| Select primitive      | + |
| LIMIT		
| limit     | + |
| offset     | + |
