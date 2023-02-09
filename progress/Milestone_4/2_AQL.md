## AQL

AQL queries are used to search and filter data.

Archetype Query Language (AQL) is a declarative query language developed specifically for expressing queries used for searching and retrieving the clinical data found in archetype-based EHRs.

When new EHR documents are created, the homomorphically encrypted data they contain is placed in a special DataSearch index tree structure located in a blockchain. Indexes are searched using a smart contract. Later selections can be made from this data using AQL queries.

<p align="center">
  <img width="75%" src="https://user-images.githubusercontent.com/8058268/174270324-1218d6ba-4cf5-497d-b455-cb084b129141.svg">
</p>

You can find a detailed description of the AQL specification on the openEHR website: <https://specifications.openehr.org/releases/QUERY/latest/AQL.html>

Query example: Get the latest 5 abnormal blood pressure values that were recorded in a health encounter for a specific patient.

```
SELECT
    obs/data[at0001]/events[at0006]/data[at0003]/items[at0004]/value/magnitude AS systolic,
    obs/data[at0001]/events[at0006]/data[at0003]/items[at0005]/value/magnitude AS diastlic,
    c/context/start_time AS date_time
FROM
    EHR [ehr_id/value=$ehrUid]
        CONTAINS COMPOSITION c [openEHR-EHR-COMPOSITION.encounter.v1]
            CONTAINS OBSERVATION obs [openEHR-EHR-OBSERVATION.blood_pressure.v1]
WHERE
    obs/data[at0001]/events[at0006]/data[at0003]/items[at0004]/value/magnitude >= 140 OR
    obs/data[at0001]/events[at0006]/data[at0003]/items[at0005]/value/magnitude >= 90
ORDER BY
    c/context/start_time DESC
LIMIT 5
```

On receipt of a request, the IPEHR gateway interprets the request into a set of conditions, which is used to search the DataSearch index structure and returns the result to the requestor as specific values or as links to documents containing the requested data. 


## Implementation AQL according to the openEHR protocol.

Client applications use AQL data queries, this way platform can be used by the 3rd party applications supporting the standard.

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
