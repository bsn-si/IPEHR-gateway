When new EHR documents are created, the homomorphically encrypted data they contain is placed in a special DataSearch index tree structure, so that later selections can be made from this data using AQL queries.

Схема

The DataSearch index will be located in a blockchain. The index will be searched using a smart contract.

AQL queries are used to search and filter data.

Archetype Query Language (AQL) is a declarative query language developed specifically for expressing queries used for searching and retrieving the clinical data found in archetype-based EHRs.

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

Access groups are used to limit access to data that can be fetched.