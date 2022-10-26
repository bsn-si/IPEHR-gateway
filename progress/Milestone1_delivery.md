# Milestone Delivery :mailbox:

* **Application Document:** [ipEHR application](https://github.com/filecoin-project/devgrants/issues/418)
* **Milestone Number:** 1

**Context**

In this milestone we've designed and developed an MH-ORM and medical data storage structure.

**Deliverables**

| Number                | Deliverable                              | Link                                                                                                                                  | Notes                                                                                                                                                                                                                                                                                                                  |
|-----------------------|------------------------------------------|---------------------------------------------------------------------------------------------------------------------------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| :heavy_check_mark:    | License                                  | [LICENSE](https://github.com/bsn-si/IPEHR-gateway/blob/develop/LICENSE)                                                               | Apache 2.0 license                                                                                                                                                                                                                                                                                                     |
| :heavy_check_mark:    | Testing Guide                            | [Readme.md](https://github.com/bsn-si/IPEHR-gateway/blob/develop/README.md#how-to)                                                    | "How To" guide                                                                                                                                                                                                                                                                                                         |
| 1. :heavy_check_mark: | Design of the medical data storage       | [See IPEHR-gateway/Storage](https://github.com/bsn-si/IPEHR-gateway/blob/develop/progress/Milestone_1/1_Storage/README.md)            | This project will use Filecoin as the EHR document repository. The implementation of saving documents in Filecoin will be implemented in the next milestone.                                                                                                                                                           | 
| 2. :heavy_check_mark: | Design of the index storage              | [See IPEHR-gateway/Index_design](https://github.com/bsn-si/IPEHR-gateway/tree/develop/progress/Milestone_1/2_Index_design#readme)     | To store the indexes, it is planned to create a special smart contract on the blockchain. This will ensure security, data immutability and fast access. Each document that is added to the system is encrypted using the ChaCha20-Poly1305 algorithm. A separate key is generated for each document.                   | 
| 3. :heavy_check_mark: | Data encryption methods                  | [See IPEHR-gateway/Encryption](https://github.com/bsn-si/IPEHR-gateway/blob/develop/progress/Milestone_1/3_Encryption/README.md)      | When stored in the repository, EHR documents are pre-encrypted using the ChaCha20-Poly1305 streaming algorithm with message authentication. To encrypt each document a unique key is generated - a random sequence of 256 bits (32 bytes) + a unique 96 bits (12 bytes). Document ID is used as an authentication tag. |
| 4. :heavy_check_mark: | Filtering functionality                  | [See IPEHR-gateway/Filtering](https://github.com/bsn-si/IPEHR-gateway/tree/develop/progress/Milestone_1/4_Filtering)                  | When new EHR documents are created, the homomorphically encrypted data they contain is placed in a special DataSearch index tree structure, so that later selections can be made from this data using AQL queries.                                                                                                     |
| 5. :heavy_check_mark: | Record creation and update functionality | [See IPEHR-gateway/EHR_creation](https://github.com/bsn-si/IPEHR-gateway/tree/develop/progress/Milestone_1/5_EHR_creation_and_update) | IPEHR-gateway implements functions for creating and updating documents according to openEHR standards. The EHR Information Model version 1.1.0, the latest stable version at the time of development, was used.                                                                                                        | 
| 6. :heavy_check_mark: | API                                      | [See IPEHR-gateway/API](https://github.com/bsn-si/IPEHR-gateway/tree/develop/progress/Milestone_1/6_API)                              | The minimum basic version of the REST API has been implemented to support document handling according to the latest stable version of the openEHR specification. In the next milestones the API will be supplemented with other methods to fully comply with the openEHR specifications.                               | 
| 7. :heavy_check_mark: | Extensive testing                        | [ipEHR wiki](https://github.com/bsn-si/IPEHR-gateway/wiki/IPEHR-project#workflow-example)                                             | Milestone deliverables testing examples with video guides                                                                                                                                                                                                                                                              |

### Project introduction:

[![Watch the video](https://media.bsn.si/ipehr/logo_intro.jpg)](https://www.youtube.com/watch?v=nJFA5W4qoEw)


# Workflow example

The following methods are based on the [standard specification of OpenEHR](https://specifications.openehr.org/releases/ITS-REST/Release-1.0.2/ehr.html)

The following workflow is showing how to create, update and find EHR information using swagger

Precondition: go to [Swagger](http://gateway.iprhr.org/swagger/index.html)
<img width="1440" alt="Screen Shot 2022-06-16 at 22 01 37" src="https://user-images.githubusercontent.com/15906120/174154411-7052f363-684d-4796-ad4d-8864d9b0d255.png">

P.S. For creating a UUID you can use a generator (e.g. [uuid generator](https://www.uuidgenerator.net/version1))

## Create an EHR

### Create an EHR
1. Click `POST /ehr` [Create EHR](https://gateway.ipehr.org/swagger/index.html#/EHR/post_ehr) method
1. Click `Try it out`
1 .Fill in `Authorization` (Bearer <JWT>), `AuthUserId` (UUID format e.g. `46f3df9f-817c-4910-825f-92d5ea595c73`), `EhrSystemID` (e.g. `openEHRSys.example.com`)
   and Prefer `return=representation`
1. Click `Execute`

Result: in the response, you see a created EHR with structured data in JSON format. Also in the response, you see other 
important information (e.g. fields like `ehr_id`, `ehr_status`, `ehr_access`, `time_created` that help you to work with EHR in
the future requests)

[Watch video instruction](https://media.bsn.si/ipehr/Create_basic_EHR.mp4) 📹 

### Create EHR with subject_id and subject_namespace params
1. Click `PUT /ehr/{ehr_id}` [Create EHR](https://gateway.ipehr.org/swagger/index.html#/EHR/post_ehr) method
2. Click `Try it out`
3. Fill in `Authorization` (Bearer <JWT>), `AuthUserId` (UUID format e.g. `46f3df9f-817c-4910-825f-92d5ea595c73`),
   `EhrSystemID` (e.g. `openEHRSys.example.com`) and specify `ehr_id` (also UUID format)
4. Change params in `Request` field
    - "subject" -> "external_ref" -> "id" ->  "value" -> put here id (e.g. num123)
    - "subject" -> "external_ref" -> "id" ->  "namespace" -> put here namespace (e.g. test namespace)
5. Click `Execute`

Result: in the response, you see a created EHR with structured data in JSON format. Also in the response, you see other 
important information (e.g. fields like ehr_id, ehr_status, ehr_access, time_created that help you to work with EHR in future 
requests). There are no subject_id and subject_space fields, they can be only requested with GET request

[Watch video instruction](https://media.bsn.si/ipehr/create_with_different_parameters.mp4) 📹 

### Create EHR with id
1. Click `PUT /ehr/{ehr_id}` [Create EHR with id](https://gateway.ipehr.org/swagger/index.html#/EHR/put_ehr__ehr_id_) method
2. Click `Try it out`
3. Fill in `Authorization` (Bearer <JWT>), `AuthUserId` (UUID format e.g. 46f3df9f-817c-4910-825f-92d5ea595c73), `Prefer` (return=representation),
   `EhrSystemID` (e.g. `openEHRSys.example.com`) and specify `ehr_id` (also UUID format)
4. Click `Execute`

Result: in the response, you see a created EHR with structured data in JSON format. Also in the response, you see other 
important information (e.g. fields like `ehr_id`, `ehr_status`, `ehr_access`, `time_created` that help you to work with EHR in 
the future requests)

## Get EHR with different parameters

### Getting info on created summary EHR by the subject id

Precondition: [Create EHR with subject_id and subject_namespace](https://github.com/bsn-si/IPEHR-gateway/wiki/IPEHR-project/_edit#create-ehr-with-subject_id-and-subject_name-params). Copy params subject_id and subject_namespace to a buffer

1. Click `/GET ehr/{ehr_id}` [Get EHR summary by subject id](https://gateway.ipehr.org/swagger/index.html#/EHR/get_ehr)
1. Click `Try it out`
1. Paste `subject_id`, `subject_namespace`, `Authorization` (Bearer <JWT>) and `AuthUserId` from previously created EHR
1. Click `Execute`

Result: in the response, you see created before EHR with the requested `subject_id` and `subject_namespace`. Here you 
can find exact EHR without `ehr_id`.

[Watch video instruction](https://media.bsn.si/ipehr/get_ehr_by_subject_id.mp4) 📹 

### Getting info on created EHR summary by id

Precondition: [EHR was created](https://github.com/bsn-si/IPEHR-gateway/wiki/IPEHR-project/_edit#create-an-ehr-1) and 
`ehr_id` is copied to a buffer

1. Click `/GET her GET` [Get EHR summary by summary id](https://gateway.ipehr.org/swagger/index.html#/EHR/get_ehr__ehr_id_)
1. Click `Try it out`
1. Fill in (or paste) `ehr_id` field
1. Click `Execute`

Result: in the response, you see created before EHR which was found with `ehr_id`

[Watch video instruction](https://media.bsn.si/ipehr/get_by_summary_id.mp4) 📹 

### Getting info on EHR status version by time

Precondition: [EHR was created](https://github.com/bsn-si/IPEHR-gateway/wiki/IPEHR-project/_edit#create-an-ehr-1) and 
`ehr_id`, `Authorization`, `AuthUserId`,  `Prefer`, and `time_created` (`return=representation`) are copied to a buffer

1. Click `/GET ehr/{ehr_id}/ehr_status` [Get EHR_STATUS version by time](https://gateway.ipehr.org/swagger/index.html#/EHR_STATUS/get_ehr__ehr_id__ehr_status)
2. Click `Try it out`
3. Fill in `ehr_id`, `Authorization`, `AuthUserId`, `Prefer`, and `time_created` (or any other time between creation EHR and current time
   in format "2022-06-22T13:26:39.042+00:00") with copied information 
4. Click `Execute`

Result: in the response, you see the EHR status version by time

[Watch video instruction](https://media.bsn.si/ipehr/get_version_by_time.mp4) 📹 

### Getting info on EHR_STATUS by version id

Precondition: [EHR was created](https://github.com/bsn-si/IPEHR-gateway/wiki/IPEHR-project/_edit#create-an-ehr-1) and 
ehr_id, Authorization, AuthUserId, version_id (look at "ehr_status" -> "value") are copied to a buffer

1. Click `GET /ehr/{ehr_id}/ehr_status/{version_uid}` [Get EHR_STATUS by version id](https://gateway.ipehr.org/swagger/index.html#/EHR_STATUS/get_ehr__ehr_id__ehr_status__version_uid_)
2. Click `Try it out`
3. Fill in `ehr_id`, `Authorization`, `AuthUserId`, `EhrSystemID`, `version_id` (e.g. `8cf1779d-f050-4be3-a671-579bf277f294::openEHRSys.example.com::1`
   where `::1` is version of EHR, and `openEHRSys.example.com` is your `EhrSystemID`) with copied information 
4. Click `Execute`

Result: in the response, you see the EHR status EHR_STATUS by version id (at the example it is `::1` at the end which
means requested the first version)

[Watch video instruction](https://media.bsn.si/ipehr/get_by_summary_id.mp4) 📹 

## Update EHR

### Update EHR status

Precondition: [EHR was created](https://github.com/bsn-si/IPEHR-gateway/wiki/IPEHR-project/_edit#create-an-ehr-1) and 
ehr_id, Authorization, AuthUserId are copied to a buffer 

1. Click `/PUT /ehr/{ehr_id}/ehr_status` [Update EHR_STATUS](https://gateway.ipehr.org/swagger/index.html#/EHR_STATUS/put_ehr__ehr_id__ehr_status)
1. Click `Try it out`
1. Fill in `Authorization`, `AuthUserId` and `EhrSystemID` with copied information 
1. Fill in a new ehr_id according to the UUID format
1. Fill in `If-Match` with `version_id` (e.g. the current version is `8cf1779d-f050-4be3-a671-579bf277f294::openEHRSys.example.com::1`)
1. Fill in Request with a body template (also can find it [here](https://specifications.openehr.org/releases/ITS-REST/Release-1.0.2/ehr.html#ehr_status-ehr_status-put)):
```json
{
  "_type": "EHR_STATUS",
  "archetype_node_id": "openEHR-EHR-EHR_STATUS.generic.v1",
  "name": {
    "value": "EHR Status"
  },
  "uid": {
    "_type": "OBJECT_VERSION_ID",
    "value": "8849182c-82ad-4088-a07f-48ead4180515::openEHRSys.example.com::2"
  },
  "subject": {
    "external_ref": {
      "id": {
        "_type": "HIER_OBJECT_ID",
        "value": "324a4b23-623d-4213-cc1c-23f233b24234"
      },
      "namespace": "DEMOGRAPHIC",
      "type": "PERSON"
    }
  },
  "other_details": {
    "_type": "ITEM_TREE",
    "archetype_node_id": "at0001",
    "name": {
      "value": "Details"
    },
    "items": []
  },
  "is_modifiable": true,
  "is_queryable": true
}
```

1. In the body change data (e.g. `"is_queryable": false` instead of `"is_queryable": true`)
1. Click `Execute`

Result: in the response, you see an updated EHR. Now you can find EHR with a new version and in new version, you will 
see that the changed param has a new value

[Watch video instruction](https://media.bsn.si/ipehr/update_ehr.mp4) 📹 

## Create a composition

### Creates the first version of a new COMPOSITION in the EHR identified by ehr_id

Precondition: [EHR was created](https://github.com/bsn-si/IPEHR-gateway/wiki/IPEHR-project/_edit#create-an-ehr-1) and 
`ehr_id`, `Authorization`, `AuthUserId`, and `Prefer` (`return=representation`) are copied to a buffer

1. Click `/POST /ehr/{ehr_id}/composition` [Create COMPOSITION](https://gateway.ipehr.org/swagger/index.html#/COMPOSITION/post_ehr__ehr_id__composition)
2. Click `Try it out`
3. Fill in `EhrSystemID`, for instance if your composition ID will be `8849182c-82ad-4088-a07f-48ead4180515::openEHRSys.example.com::1`
   then your `EhrSystemID` is `openEHRSys.example.com`
4. Fill in `ehr_id`, `Authorization`, `AuthUserId` and `Prefer` with copied information 
5. In the `Request` put information (body of JSON) from [this](https://media.bsn.si/ipehr/composition.json) file 
6. Click `Execute`

Result: in the response, you see a created composition with a lot of information inside (e.g. archetype_node_id, composer,
health_care_facility etc.)

[Watch video instruction](https://media.bsn.si/ipehr/create_composition.mp4) 📹 

## Execute ad-hoc (non-stored) AQL query

### Execute AQL query
1. Click `/POST query/aql` [Execute ad-hoc (non-stored) AQL query](https://gateway.ipehr.org/swagger/index.html#/QUERY/post_query_aql)
1. Click `Try it out`
1. Fill in `Authorization`, `AuthUserId`
1. Fill in `Request` with AQL request
```json
{
  "q": "SELECT c/uid/value, obs/data[at0001]/events[at0006]/data[at0003]/items[at0004]/value/magnitude AS systolic FROM EHR e[ehr_id/value=$ehr_id] CONTAINS COMPOSITION c[openEHR-EHR-COMPOSITION.health_summary.v1] CONTAINS OBSERVATION obs[openEHR-EHR-OBSERVATION.blood_pressure.v2] WHERE obs/data[at0001]/events[at0006]/data[at0003]/items[at0004]/value/magnitude >= $systolic_bp",
  "offset": 0,
  "fetch": 10,
  "query_parameters": {
    "ehr_id": "07fa4e05-df30-4a2b-9612-6bf1d28ff80c",
    "systolic_bp": 140
  }
}
```
1. Click `Execute`
Result: In the response, you see a requested information

[Watch video instruction](https://media.bsn.si/ipehr/ad-hoc.mp4) 📹 
