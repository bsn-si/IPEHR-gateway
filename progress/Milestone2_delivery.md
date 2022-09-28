# Milestone Delivery :mailbox:

* **Application Document:** [ipEHR application](https://github.com/filecoin-project/devgrants/issues/418)
* **Milestone Number:** 2

**Context**

In this milestone we've developed the functionality of access to the data stored in Filecoin.

**Deliverables**

| Number                | Deliverable                              | Link                                                                                                                                  | Notes                                                                                                                                                                                                                                                                                                                  |
|-----------------------|------------------------------------------|---------------------------------------------------------------------------------------------------------------------------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| :heavy_check_mark:    | License                                  | [LICENSE](https://github.com/bsn-si/IPEHR-gateway/blob/develop/LICENSE)                                                               | Apache 2.0 license                                                                                                                                                                                                                                                                                                     |
| :heavy_check_mark:    | Testing Guide                            | [Readme.md](https://github.com/bsn-si/IPEHR-gateway/blob/develop/README.md#how-to)                                                    | The "How To" guide is supplemented with all new features developed in this milestone                                                                                                                                                                                                                                                                                                        |
| 1. :heavy_check_mark: | Integrate Filecoin as a storage for MH-ORM database | See [Filecoin_storage](https://github.com/bsn-si/IPEHR-gateway/blob/develop/progress/Milestone_2/1_Filecoin_integration.md) page | At IPEHR we use a two-tiered document storage system: 1. The IPFS distributed file system for quick access. 2. The Filecoin storage network for long-term guaranteed storage of EHR documents. Filecoin storage deals are made for a fixed period of time. Usually from 180 days. As the deadline approaches, the deal must be extended. This functionality will be added in the next phases. | 
| 2. :heavy_check_mark: | Integrate a storage for MH-ORM indexes | See [Indexes_storage](https://github.com/bsn-si/IPEHR-gateway/blob/develop/progress/Milestone_2/2_Indexes_storage.md) page | An EVM smart contract has been developed to store EHR document indexes. | 
| 3. :heavy_check_mark: | An algorithm of data re-encryption while changing access rights | See [Revoking_access](https://github.com/bsn-si/IPEHR-gateway/blob/develop/progress/Milestone_2/3_Revoking_access.md) page | To grant access to a document, the document access key is asymmetrically encrypted with the public key of the user (or group) being granted access and added to the IPEHR smart contract table. To revoke access to a document, a user with the rights of the document owner just deletes the appropriate record in the table. |
| 4. :heavy_check_mark: | Records chains | See [Chains of records](https://github.com/bsn-si/IPEHR-gateway/blob/develop/progress/Milestone_2/4_Chains_of_records.md) page | When a new EHR document is added, the digital signature of the user creating the document is sent along with the document. The signature is stored with the document. This allows you to authorize a request to create a document. |
| 5. :heavy_check_mark: | Performance tests and optimisation | See [Workflow examples](https://github.com/bsn-si/IPEHR-gateway/blob/develop/progress/Milestone2_delivery.md#workflow-example) | Milestone deliverables testing examples with video guides | 
| 6. :heavy_check_mark: | Payment logic | See [payment logic](https://github.com/bsn-si/IPEHR-gateway/blob/develop/progress/Milestone_2/6_Payment_logic.md) | when working with the IPEHR system, payment for transactions is made from the organization's account within the available balance | 

### Project introduction:

[![Watch the video](https://media.bsn.si/ipehr/logo_intro.jpg)](https://www.youtube.com/watch?v=nJFA5W4qoEw)

# Workflow example

The following methods are based on the [standard specification of OpenEHR](https://specifications.openehr.org/releases/ITS-REST/Release-1.0.2/ehr.html)

<img width="800" src="https://user-images.githubusercontent.com/15906120/192385584-0d50ac5b-ae3c-40f2-bce5-8a0ab3a8da1f.png">

P.S. For creating a UUID you can use a generator (e.g [uuidgenerator](https://www.uuidgenerator.net/version1))

## How to create EHR

### Create a general EHR
Precondition: Open [swagger](https://gateway.ipehr.org/swagger/index.html)
1. Choose [Create EHR](https://gateway.ipehr.org/swagger/index.html#/EHR/post_ehr) method and click **Try it out**
1. Fill in **AuthUserId** (UUID format e.g 46f3df9f-817c-4910-825f-92d5ea595c73), **EhrSystemId** (any system name or UUID) and **Prefer** (return=representation) and click **Execute**

**Result**: While the request is running, the document is saved to IPFS, the indexes and meta-information about the document are written to a smart contract on the blockchain, and a deal is started to store the document in Filecoin.
You will get a **requestid** in the **response header**, which can be used to get more information about the status of the request.

[Watch video instruction](https://media.bsn.si/ipehr/v2/create_general_ehr.mp4) 📹 

#### Check if the file is available on the IPFS network
Precondition: Execute [Create EHR](https://github.com/bsn-si/IPEHR-gateway/wiki/IPEHR-project#create-a-general-ehr) and [Get created EHR by ID](https://github.com/bsn-si/IPEHR-gateway/wiki/IPEHR-project#get-created-ehr-by-request-id) steps
1. Copy the **CID** value from the **EhrCreate** section 
1. Check if the file is available on the IPFS network using the link [ipfs.io/ipfs/XXX](https://ipfs.io/ipfs/) where **XXX** is your **CID** value

**Result**: The document will be downloaded to your computer, but you will not be able to read the content because it is encrypted

[Watch video instruction](https://media.bsn.si/ipehr/v2/check_ipfs_file.mp4) 📹 

#### Check if the execution of the transaction is available on the goerli.etherscan.io
Precondition: Execute [Create EHR](https://github.com/bsn-si/IPEHR-gateway/wiki/IPEHR-project#create-a-general-ehr) and [Get created EHR by ID](https://github.com/bsn-si/IPEHR-gateway/wiki/IPEHR-project#get-created-ehr-by-id) steps
1. Copy the **hash** of a transaction in sections such as **SetEhrDocs** or **SetDocAccess**
1. Check the execution of the transaction using the [goerli.etherscan.io](https://goerli.etherscan.io) explorer

**Result**: You see the execution of the transaction

[Watch video instruction](https://media.bsn.si/ipehr/v2/check_goerli.mp4) 📹 

#### Check Filecoin document storage transaction is complete
Precondition: Execute [Create EHR](https://github.com/bsn-si/IPEHR-gateway/wiki/IPEHR-project#create-a-general-ehr) and [Get created EHR by ID](https://github.com/bsn-si/IPEHR-gateway/wiki/IPEHR-project#get-created-ehr-by-id) steps. Wait until **FilecoinStartDeal** status change from 'Pending' to 'Success' (usually takes about 24 hours)
1. Copy the **dealId**
1. Go to [filfox.info/en/deal/XXX](https://filfox.info/en/deal/) and paste to the search your dealId instead of XXX

**Result**: You see file and data stored on filecoin

[Watch video instruction](https://media.bsn.si/ipehr/v2/check_filecoin.mp4) 📹 

### Create EHR with EHR_ID
1. Choose [Create EHR with id](https://gateway.ipehr.org/swagger/index.html#/EHR/put_ehr__ehr_id_) method and click **Try it out**
1. Fill in **AuthUserId** (UUID format e.g 46f3df9f-817c-4910-825f-92d5ea595c73), **EhrSystemId** (any system name or UUID), **Prefer** (return=representation) and **ehr_id** (UUID as a user-specified EHR identifier) and click **Execute**

**Result**: in the response, you see a created EHR with structured data in JSON format and exact ehr_id that you put before

[Watch video instruction](https://media.bsn.si/ipehr/v2/create_ehr_with_id.mp4) 📹

### Create EHR with different parameters, such as subject id and id namespace
1. Choose [PUT Create EHR](https://gateway.ipehr.org/swagger/index.html#/EHR/post_ehr) method click **Try it out**
1. Fill in **AuthUserId** (UUID format e.g 46f3df9f-817c-4910-825f-92d5ea595c73) and specify **ehr_id** (also UUID format)
1. Change params in the Request field:
** "subject" -> "external_ref" -> "id" -> "value" -> put here id (e.g num123)
** "subject" -> "external_ref" -> "id" -> "namespace" -> put here namespace (e.g testnamespace)
1. Click **Execute**

**Result**: You see a created EHR with structured data in JSON format in the response. There are no subject_id and subject_space fields, they can be only requested with GET request

[Watch video instruction](https://media.bsn.si/ipehr/v2/create_ehr_with_subjectnamespace_and_subjectid.mp4) 📹

## How to get EHR with different parameters

### Get created EHR by request id
Precondition: Execute [Create EHR](https://github.com/bsn-si/IPEHR-gateway/wiki/IPEHR-project#create-a-general-ehr) steps and copy **AuthUserId** and **requestid**
1. Choose [Get request by ID](https://gateway.ipehr.org/swagger/index.html#/REQUEST/get_requests__request_id_) method
1. Paste **AuthUserId** and **request_id** values and click **Execute**

**Result**: In the response, you will get additional information about the progress of this request, including the transaction hashes in the blockchain, the file cid in IPFS, and the dealCid - the transaction id in Filecoin

[Watch video instruction](https://media.bsn.si/ipehr/v2/get_by_id.mp4) 📹 

### Get info on created summary EHR by the subject id

Precondition: [Create EHR with subject_id and subject_namespace](https://github.com/bsn-si/IPEHR-gateway/wiki/IPEHR-project#create-ehr-with-different-parameters-such-as-subject-id-and-id-namespace). Copy params subject_id and subject_namespace to a buffer

1. Choose [Get EHR summary by subject id](https://github.com/bsn-si/IPEHR-gateway/wiki/IPEHR-project#get-created-ehr-by-id) method and click **Try it out**
1. Paste **subject_id**, **subject_namespace** and **AuthUserId** from previously created EHR and click **Execute**

**Result**: in the response, you see created before EHR with the requested **subject_id** and **subject_namespace**. Here you can find exact EHR without ehr_id

[Watch video instruction](https://media.bsn.si/ipehr/get_ehr_by_subject_id.mp4) 📹 

### Get info on created EHR summary by id
Precondition: [EHR was created](https://github.com/bsn-si/IPEHR-gateway/wiki/IPEHR-project/_edit#create-an-ehr-1) and **ehr_id** is copied to a buffer
1. Choose [Get EHR summary by summary id](https://gateway.ipehr.org/swagger/index.html#/EHR/get_ehr__ehr_id_) method and click **Try it out**
1. Fill in (or paste) **ehr_id**, **AuthUserId** and **EhrSystemId** field and click **Execute**

**Result**: in the response, you see created before EHR which was found with **ehr_id**

[Watch video instruction](https://media.bsn.si/ipehr/get_by_summary_id.mp4) 📹 

### Get info on EHR status version by time
Precondition: [EHR was created](https://github.com/bsn-si/IPEHR-gateway/wiki/IPEHR-project#create-a-general-ehr) and **ehr_id**, **AuthUserId**, **Prefer**, and **time_created** are copied to a buffer
1. Choose [Get EHR_STATUS version by time](https://gateway.ipehr.org/swagger/index.html#/EHR_STATUS/get_ehr__ehr_id__ehr_status) method and click **Try it out**
1. Fill in **ehr_id**, **AuthUserId**, **Prefer**, and **time_created** (or any other time between creation EHR and current time in format "2022-06-22T13:26:39.042+00:00") with copied information and click **Execute**

**Result**: in the response, you see the EHR status version by the exact time

[Watch video instruction](https://media.bsn.si/ipehr/get_version_by_time.mp4) 📹 

### Getting info on EHR_STATUS by version id

Precondition: [EHR was created](https://github.com/bsn-si/IPEHR-gateway/wiki/IPEHR-project#create-a-general-ehr) and make sure that **ehr_id**, **AuthUserId**, **version_id** (look at "ehr_status" -> "value") are copied to a buffer
1. Choose [Get EHR_STATUS by version id](https://gateway.ipehr.org/swagger/index.html#/EHR_STATUS/get_ehr__ehr_id__ehr_status__version_uid_) method and  click **Try it out**
1. Fill in **ehr_id**, **AuthUserId**, **EhrSystemId** and **version_id** (e.g 8cf1779d-f050-4be3-a671-579bf277f294::openEHRSys.example.com::1 where ::1 is version of EHR) with copied information and click **Execute**

**Result**: in the response, you see the EHR status EHR_STATUS by version id (in the example it is ::1 at the end which means requested the first version)

[Watch video instruction](https://media.bsn.si/ipehr/get_version_id.mp4) 📹 

## How to update EHR status

Precondition: [EHR was created](https://github.com/bsn-si/IPEHR-gateway/wiki/IPEHR-project#create-a-general-ehr) and **ehr_id**, **AuthUserId** are copied to a buffer 
1. Choose [Update EHR_STATUS](https://gateway.ipehr.org/swagger/index.html#/EHR_STATUS/put_ehr__ehr_id__ehr_status) method and click **Try it out**
1. Fill in **AuthUserId** and ** EhrSystemId** with copied information 
1. Fill in a new **ehr_id** according to the UUID format
1. Fill in **If-Match** with **current version_id** (e.g the current version is 8cf1779d-f050-4be3-a671-579bf277f294::openEHRSys.example.com::1)
1. Fill in **Request** with a body template (also can find it [here](https://specifications.openehr.org/releases/ITS-REST/Release-1.0.2/ehr.html#ehr_status-ehr_status-put)):
    > {
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
    
1. In the body change data (e.g. "is_queryable": false instead of "is_queryable": true) and  click **Execute**

**Result**: in the response, you see an updated EHR. Now you can find EHR with a new version and in new version, you will see that the changed param has a new value

[Watch video instruction](https://media.bsn.si/ipehr/update_ehr.mp4) 📹 

## How to work with composition

### Create composition

Precondition: [EHR was created](https://github.com/bsn-si/IPEHR-gateway/wiki/IPEHR-project#create-a-general-ehr) and **ehr_id**, **AuthUserId**, and **Prefer** (return=representation) are copied to a buffer
1. Choose [Create composition](https://gateway.ipehr.org/swagger/index.html#/COMPOSITION/post_ehr__ehr_id__composition) method and  click **Try it out**
1. Fill in **ehr_id**, **AuthUserId**, **EhrSystemId** and **Prefer** with copied information 
1. In the **Request** put information (body of JSON) from [this](https://media.bsn.si/ipehr/v2/create_composition.json) file and click **Execute**

**Result**: in the response, you see a created composition with a lot of information inside (e.g archetype_node_id, composer, health_care_facility etc)

[Watch video instruction](https://media.bsn.si/ipehr/v2/create_composition.mp4) 📹 

### Get composition by version id

Precondition: [Composition is created](https://github.com/bsn-si/IPEHR-gateway/wiki/IPEHR-project#create-composition)
1. Choose [Get composition by version id](https://gateway.ipehr.org/swagger/index.html#/COMPOSITION/get_ehr__ehr_id__composition__version_uid_) method and click **Try it out**
1. Fill in **ehr_id**, **AuthUserId**, **versioned_object_uid**, **EhrSystemId** with copied information 
Click **Execute** 

**Result**: in the response, you see composition by version id

[Watch video instruction](https://media.bsn.si/ipehr/v2/get_composition.mp4) 📹 

### Update composition by version id

Precondition: [EHR was created](https://github.com/bsn-si/IPEHR-gateway/wiki/IPEHR-project#create-a-general-ehr) and **ehr_id**, **AuthUserId**, and **Prefer** (return=representation) are copied to a buffer and [Composition is created](https://github.com/bsn-si/IPEHR-gateway/wiki/IPEHR-project#update-composition-by-version-id)
1. Choose [Update composition by version id](https://gateway.ipehr.org/swagger/index.html#/COMPOSITION/put_ehr__ehr_id__composition__versioned_object_uid_) method and  click **Try it out**
1. Fill in **ehr_id**, **AuthUserId**, **versioned_object_uid**, **EhrSystemId**, **If-Match** and **Prefer** with copied information 
1. In the **Request** change information and click **Execute** 

**Result**: in the response, you see an updated composition

[Watch video instruction](https://media.bsn.si/ipehr/v2/update_composition.mp4) 📹 

### Delete composition by version id

Precondition: [EHR was created](https://github.com/bsn-si/IPEHR-gateway/wiki/IPEHR-project#create-a-general-ehr) and **ehr_id**, **AuthUserId**, and **Prefer** (return=representation) are copied to a buffer and [Composition is created](https://github.com/bsn-si/IPEHR-gateway/wiki/IPEHR-project#update-composition-by-version-id)
1. Choose [Delete composition by version id](https://gateway.ipehr.org/swagger/index.html#/COMPOSITION/delete_ehr__ehr_id__composition__preceding_version_uid_) method and  click **Try it out**
1. Fill in **ehr_id**, **AuthUserId**, **preceding_version_uid**, **EhrSystemId** with copied information 
1. Click **Execute** 

**Result**: composition by version id is deleted. You can additionally check it by trying to [Get composition by id](https://github.com/bsn-si/IPEHR-gateway/wiki/IPEHR-project#get-composition-by-version-id)

[Watch video instruction](https://media.bsn.si/ipehr/v2/delete_composition.mp4) 📹 

## Execute ad-hoc (non-stored) AQL query

### Execute AQL query
1. Choose [Execute ad-hoc (non-stored) AQL query](https://gateway.ipehr.org/swagger/index.html#/QUERY/post_query_aql) method and click **Try it out**
1. Fill in **AuthUserId**
1. Fill in **Request** with AQL request
    > {
  "q": "SELECT c/uid/value, obs/data[at0001]/events[at0006]/data[at0003]/items[at0004]/value/magnitude AS systolic FROM EHR e[ehr_id/value=$ehr_id] CONTAINS COMPOSITION c[openEHR-EHR-COMPOSITION.health_summary.v1] CONTAINS OBSERVATION obs[openEHR-EHR-OBSERVATION.blood_pressure.v2] WHERE obs/data[at0001]/events[at0006]/data[at0003]/items[at0004]/value/magnitude >= $systolic_bp",
  "offset": 0,
  "fetch": 10,
  "query_parameters": {
    "ehr_id": "07fa4e05-df30-4a2b-9612-6bf1d28ff80c",
    "systolic_bp": 140
  }
}
1. Click **Execute**
Result: In the response, you see a requested information

[Watch video instruction](https://media.bsn.si/ipehr/ad-hoc.mp4) 📹 
