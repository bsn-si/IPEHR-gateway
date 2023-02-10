<p align="center">
  <img width="75%" src="https://user-images.githubusercontent.com/8058268/174271312-3eec1fdb-ad70-4492-8f0e-5624d1a7c408.svg">
</p>

## DataSearch index

When new EHR documents are created, the homomorphically encrypted data they contain is placed in a special DataSearch index tree structure, so that later selections can be made from this data using AQL queries.

<p align="center">
  <img width="75%" src="https://user-images.githubusercontent.com/8058268/174270324-1218d6ba-4cf5-497d-b455-cb084b129141.svg">
</p>

The DataSearch index is located in a blockchain (currently [Goerli Testnet](https://goerli.net/) is used, later will be deployed to FEVM). Indexes are searched using a [smart contract](https://github.com/bsn-si/IPEHR-blockchain-indexes).

## AQL

AQL queries are used to search and filter data.

Archetype Query Language (AQL) is a declarative query language developed specifically for expressing queries used for searching and retrieving the clinical data found in archetype-based EHRs.

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

The implementation of a full-fledged AQL query interpreter will be done in the following steps.

## Data access

Access groups are used to limit access to data that can be fetched.  
A detailed description of the homomorphic data encryption methodology in indexes can be found in section [2. Index design](https://github.com/bsn-si/IPEHR-gateway/tree/develop/progress/Milestone_1/2_Index_design)
