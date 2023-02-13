# Milestone Delivery :mailbox:

* **Application Document:** [ipEHR application](https://github.com/filecoin-project/devgrants/issues/418)
* **Milestone Number:** 4

**Context**

In this milestone we've developed OpenEHR API and integrated it with [MH-ORM](https://github.com/bsn-si/IPEHR-gateway/blob/develop/progress/Milestone_4/4_MH_ORM.md).

**Deliverables**

| Number                | Deliverable                              | Link                                                                                                                                  | Notes                                                                                                                                                                                                                                                                                                                  |
|-----------------------|------------------------------------------|---------------------------------------------------------------------------------------------------------------------------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| :heavy_check_mark:    | License | [LICENSE](https://github.com/bsn-si/IPEHR-gateway/blob/develop/LICENSE) | Apache 2.0 license |
| :heavy_check_mark:    | Testing Guide | [Readme.md](https://github.com/bsn-si/IPEHR-gateway/blob/develop/README.md#how-to) | The "How To" guide is supplemented with all new features developed in this milestone |
| 1. :heavy_check_mark: | EHR system core | See [system description](https://github.com/bsn-si/IPEHR-gateway/blob/develop/progress/Milestone_4/1_System.md) page | The IPEHR system is a dApp consisting of the following elements: the IPEHR Gateway; Smart-contracts; IPEHR Stats. | 
| 2. :heavy_check_mark: | AQL implementation | See [implemented AQL functionality](https://github.com/bsn-si/IPEHR-gateway/blob/develop/progress/Milestone_4/2_AQL.md) page | We have developed an AQL querier (library for parsing an AQL query) and an AQL executor (a service that searches the index tree for data according to an AQL query). | 
| 3. :heavy_check_mark: | Integration with access rights | See the [access rights](https://github.com/bsn-si/IPEHR-gateway/blob/develop/progress/Milestone_4/3_Access_rights.md) description | We have implemented API methods for creating user groups, adding/removing users from groups, getting information about user groups, getting information about user access, getting information about access of groups, delegation of access to documents. |
| 4. :heavy_check_mark: | Integration with MH-ORM | See [MH-ORM integration](https://github.com/bsn-si/IPEHR-gateway/blob/develop/progress/Milestone_4/4_MH_ORM.md) description | Some information models of the OpenEHR entities were implemented according to the [specifications](https://specifications.openehr.org/releases/RM/Release-1.1.0). |

### Project introduction:

[![Watch the video](https://media.bsn.si/ipehr/logo_intro.jpg)](https://www.youtube.com/watch?v=nJFA5W4qoEw)

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

![AQL flow](https://user-images.githubusercontent.com/98888366/218052758-f98c5f20-5d1c-4bcc-8350-3c9796654da2.svg)

On receipt of a request, the IPEHR gateway interprets the request into a set of conditions, which is used to search the DataSearch index structure and returns the result to the requestor as specific values or as links to documents containing the requested data.

# Workflow example

[![image](https://user-images.githubusercontent.com/98888366/217956242-8f7596a0-13ca-4141-a214-781f2c860bcf.png)](https://media.bsn.si/ipehr/v2/how_aql_works.mp4)
