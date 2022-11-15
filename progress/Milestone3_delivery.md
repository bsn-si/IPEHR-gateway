# Milestone Delivery :mailbox:

* **Application Document:** [ipEHR application](https://github.com/filecoin-project/devgrants/issues/418)
* **Milestone Number:** 3

**Context**

In this milestone we've developed the functionality to manage access rights on a blockchain level.

**Deliverables**

| Number                | Deliverable                              | Link                                                                                                                                  | Notes                                                                                                                                                                                                                                                                                                                  |
|-----------------------|------------------------------------------|---------------------------------------------------------------------------------------------------------------------------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| :heavy_check_mark:    | License | [LICENSE](https://github.com/bsn-si/IPEHR-gateway/blob/develop/LICENSE) | Apache 2.0 license |
| :heavy_check_mark:    | Testing Guide | [Readme.md](https://github.com/bsn-si/IPEHR-gateway/blob/develop/README.md#how-to) | The "How To" guide is supplemented with all new features developed in this milestone |
| 1. :heavy_check_mark: | Research of available blockchains supporting EVM-based smart contracts | See [IPEHR smart contract](https://github.com/bsn-si/IPEHR-blockchain-indexes) page | The test contract is deployed to the Goerli testnet. For the contract fork to the FVM network a golang client is needed to correctly interact with FEVM JSON-RPC API. | 
| 2. :heavy_check_mark: | Development of the user account catalogue and the user identification mechanism | See [Users_identity](https://github.com/bsn-si/IPEHR-gateway/blob/develop/progress/Milestone_3/2_Users_identity.md) page | Users and user groups are stored in the [IPEHR smart contract](https://github.com/bsn-si/IPEHR-blockchain-indexes). The authorization of requests to the IPEHR gateway API is done via the JWT token. | 
| 3. :heavy_check_mark: | Embedding access rights management | See [Docs_access_management](https://github.com/bsn-si/IPEHR-gateway/blob/develop/progress/Milestone_3/3_Docs_access_mgmt.md) page | Documents can be grouped by arbitrary criteria. From medical classification to geographical location. Access to documents is managed according to a particular access matrix. |
| 4. :heavy_check_mark: | Design and development of the access keys storage | See [Access_store](https://github.com/bsn-si/IPEHR-gateway/blob/develop/progress/Milestone_3/4_Access_store.md) page | Access Key Storage is a hash table that is located in the smart contract. The key of the table is a special 32 byte identification number. |
| 5. :heavy_check_mark: | Design and development of a smart contract API | See [Iface](https://github.com/bsn-si/IPEHR-gateway/blob/develop/progress/Milestone_3/5_Iface_to_contract.md) | To interact with a smart contract containing a repository of users, documents and access rights, the indexer package was developed using the Go Ethereum library. To prevent the execution of repeated transactions, the nonce mechanism is used in combination with ECDSA signature verification. | 

### Project introduction:

[![Watch the video](https://media.bsn.si/ipehr/logo_intro.jpg)](https://www.youtube.com/watch?v=nJFA5W4qoEw)

# Workflow example

The following methods are based on the [standard specification of OpenEHR](https://specifications.openehr.org/releases/ITS-REST/Release-1.0.2/ehr.html)

<img width="800" src="https://user-images.githubusercontent.com/15906120/192385584-0d50ac5b-ae3c-40f2-bce5-8a0ab3a8da1f.png">

P.S. For creating a UUID you can use a generator (e.g [uuidgenerator](https://www.uuidgenerator.net/version1))

## How to create and work with a user
`Disclaimer: you have to register a user and must log in under your credentials before doing any actions. We recommend saving all essential information such as userID, tokens, EhrSystem, and the other data in a separate file for easy access to copy-paste it`

### How to register a user
Precondition: Open [swagger](https://gateway.ipehr.org/swagger/index.html)
1. Choose [Register user](https://gateway.ipehr.org/swagger/index.html#/USER/post_user_register) method and click **Try it out**
1. Fill in **EhrSystemId** (any system name or UUID format). In **Request** form put your userID (e.g user1), password and click **Execute**

**Result**: you see that the user is successfully created. You can use your userID and password to log in

[Watch video instruction](https://media.bsn.si/ipehr/v2/register_user.mp4) 📹 

### How to login under the user
Precondition: Open [swagger](https://gateway.ipehr.org/swagger/index.html) and [register](https://github.com/bsn-si/IPEHR-gateway/wiki/IPEHR-project#how-to-register-a-user) under registered user OR use exciting user
1. Choose [Login user](https://gateway.ipehr.org/swagger/index.html#/USER/post_user_login) method and click **Try it out**
1. Fill in **AuthUserId** (userID that you got before), **EhrSystemId** (that you put when register). In **Request** form put your userID and password that you registered before and click **Execute**

**Result**: you've got access and a refresh token which can be used instead password for future actions

[Watch video instruction](https://media.bsn.si/ipehr/v2/login.mp4) 📹 

### How to refresh JWT token
Precondition: Open [swagger](https://gateway.ipehr.org/swagger/index.html) and [login](https://gateway.ipehr.org/swagger/index.html#/USER/post_user_login) under registered user
1. Choose [Refresh JWT](https://gateway.ipehr.org/swagger/index.html#/USER/get_user_refresh_) method and click **Try it out**
1. Fill in **AuthUserId** (it is your REFRESH token in the format "Bearer XXX" where XXX is your refresh token), **AuthUserId** (that you are logged in) and **EhrSystemId**, click **Execute**

**Result**: your access and refresh tokens are updated. If you failed it means that you were logged out, try to log in again.

[Watch video instruction](https://media.bsn.si/ipehr/v2/refresh_jwt.mp4) 📹 

### How to logout
Precondition: Open [swagger](https://gateway.ipehr.org/swagger/index.html) and [login](https://github.com/bsn-si/IPEHR-gateway/wiki/IPEHR-project#hou-to-login-under-the-user) under registered user
1. Choose [Logout User](https://gateway.ipehr.org/swagger/index.html#/USER/post_user_logout) method and click **Try it out**
1. Fill in **AuthUserId** (it is your ACCESS token in the format "Bearer XXX" where XXX is your access token), **AuthUserId** (that you are logged in) and **EhrSystemId**, click **Execute**

**Result**: you are successfully logged out

[Watch video instruction](https://media.bsn.si/ipehr/v2/logout.mp4) 📹

## How to create EHR

### Create a general EHR
Precondition: Open [swagger](https://gateway.ipehr.org/swagger/index.html) and [login](https://github.com/bsn-si/IPEHR-gateway/wiki/IPEHR-project#hou-to-login-under-the-user) under registered user
1. Choose [Create EHR](https://gateway.ipehr.org/swagger/index.html#/EHR/post_ehr) method and click **Try it out**
1. Fill in **Authorization** (you Bearer and access token), **AuthUserId** (your userID), **EhrSystemId** (any system name or UUID) and **Prefer** (return=representation) and click **Execute**

**Result**: While the request is running, the document is saved to IPFS, the indexes and meta-information about the document are written to a smart contract on the blockchain, and a deal is started to store the document in Filecoin.
You will get a **requestid** in the **response header**, which can be used to get more information about the request's status.

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
1. Fill in **Authorization** (you Bearer and access token), **AuthUserId** (your userID), **EhrSystemId** (that you paste when registered the user), **Prefer** (return=representation) and **ehr_id** (UUID as a user-specified EHR identifier) and click **Execute**

**Result**: in the response, you see a created EHR with structured data in JSON format and exact ehr_id that you put before

[Watch video instruction](https://media.bsn.si/ipehr/v2/create_ehr_with_id.mp4) 📹

### Create EHR with different parameters, such as subject id and id namespace
1. Choose [PUT Create EHR](https://gateway.ipehr.org/swagger/index.html#/EHR/post_ehr) method click **Try it out**
1. Fill in **Authorization** (your Bearer and access token), **AuthUserId** (your userID), **EhrSystemId** (that you paste when registered the user) and specify **ehr_id** (also UUID format)
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
1. Paste **AuthUserId** and **request_id** values, fill in **Authorization** (you Bearer and access token) and click **Execute**

**Result**: In the response, you will get additional information about the progress of this request, including the transaction hashes in the blockchain, the file cid in IPFS, and the dealCid - the transaction id in Filecoin

[Watch video instruction](https://media.bsn.si/ipehr/v2/get_by_id.mp4) 📹 

### Get info on created summary EHR by the subject id

Precondition: [Create EHR with subject_id and subject_namespace](https://github.com/bsn-si/IPEHR-gateway/wiki/IPEHR-project#create-ehr-with-different-parameters-such-as-subject-id-and-id-namespace). Copy params subject_id and subject_namespace to a buffer

1. Choose [Get EHR summary by subject id](https://github.com/bsn-si/IPEHR-gateway/wiki/IPEHR-project#get-created-ehr-by-id) method and click **Try it out**
1. Paste **subject_id**, **subject_namespace** and **Authorization** (you Bearer and access token), **AuthUserId** from previously created EHR and click **Execute**

**Result**: in the response, you see created before EHR with the requested **subject_id** and **subject_namespace**. Here you can find exact EHR without ehr_id

[Watch video instruction](https://media.bsn.si/ipehr/get_ehr_by_subject_id.mp4) 📹 

### Get info on created EHR summary by id
Precondition: [EHR was created](https://github.com/bsn-si/IPEHR-gateway/wiki/IPEHR-project/_edit#create-an-ehr-1) and **ehr_id** is copied to a buffer
1. Choose [Get EHR summary by summary id](https://gateway.ipehr.org/swagger/index.html#/EHR/get_ehr__ehr_id_) method and click **Try it out**
1. Fill in (or paste) **ehr_id**, **Authorization** (you Bearer and access token), **AuthUserId** (your userID) and **EhrSystemId** (that you paste when registered the user) field and click **Execute**

**Result**: in the response, you see created before EHR which was found with **ehr_id**

[Watch video instruction](https://media.bsn.si/ipehr/get_by_summary_id.mp4) 📹 

### Get info on EHR status version by time
Precondition: [EHR was created](https://github.com/bsn-si/IPEHR-gateway/wiki/IPEHR-project#create-a-general-ehr) and **ehr_id**, **AuthUserId**, **Prefer**, and **time_created** are copied to a buffer
1. Choose [Get EHR_STATUS version by time](https://gateway.ipehr.org/swagger/index.html#/EHR_STATUS/get_ehr__ehr_id__ehr_status) method and click **Try it out**
1. Fill in **ehr_id**, **AuthUserId**, **Authorization** (you Bearer and access token), **Prefer**, and **time_created** (or any other time between creation EHR and current time in format "2022-06-22T13:26:39.042+00:00") with copied information and click **Execute**

**Result**: in the response, you see the EHR status version by the exact time

[Watch video instruction](https://media.bsn.si/ipehr/get_version_by_time.mp4) 📹 

### Getting info on EHR_STATUS by version id

Precondition: [EHR was created](https://github.com/bsn-si/IPEHR-gateway/wiki/IPEHR-project#create-a-general-ehr) and make sure that **ehr_id**, **AuthUserId**, **version_id** (look at "ehr_status" -> "value") are copied to a buffer
1. Choose [Get EHR_STATUS by version id](https://gateway.ipehr.org/swagger/index.html#/EHR_STATUS/get_ehr__ehr_id__ehr_status__version_uid_) method and  click **Try it out**
1. Fill in **ehr_id**, **Authorization** (you Bearer and access token), **AuthUserId** (your userID), **EhrSystemId** and **version_id** (e.g 8cf1779d-f050-4be3-a671-579bf277f294::openEHRSys.example.com::1 where ::1 is version of EHR) with copied information and click **Execute**

**Result**: in the response, you see the EHR status EHR_STATUS by version id (in the example it is ::1 at the end which means requested the first version)

[Watch video instruction](https://media.bsn.si/ipehr/get_version_id.mp4) 📹 

## How to update EHR status

Precondition: [EHR was created](https://github.com/bsn-si/IPEHR-gateway/wiki/IPEHR-project#create-a-general-ehr) and **ehr_id**, **AuthUserId** are copied to a buffer 
1. Choose [Update EHR_STATUS](https://gateway.ipehr.org/swagger/index.html#/EHR_STATUS/put_ehr__ehr_id__ehr_status) method and click **Try it out**
1. Fill in **Authorization** (you Bearer and access token), **AuthUserId** (your userID) and **EhrSystemId** (that you paste when registered the user) with copied information 
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
1. Fill in **ehr_id**,**Authorization** (you Bearer and access token), **AuthUserId** (your userID), **EhrSystemId** (that you paste when registered the user) and **Prefer** with copied information 
1. In the **Request** put information (body of JSON) from [this](https://media.bsn.si/ipehr/v2/create_composition.json) file and click **Execute**

**Result**: in the response, you see a created composition with a lot of information inside (e.g archetype_node_id, composer, health_care_facility etc)

[Watch video instruction](https://media.bsn.si/ipehr/v2/create_composition.mp4) 📹 

### Get composition by version id

Precondition: [Composition is created](https://github.com/bsn-si/IPEHR-gateway/wiki/IPEHR-project#create-composition)
1. Choose [Get composition by version id](https://gateway.ipehr.org/swagger/index.html#/COMPOSITION/get_ehr__ehr_id__composition__version_uid_) method and click **Try it out**
1. Fill in **ehr_id**, **Authorization** (you Bearer and access token), **AuthUserId** (your userID), **versioned_object_uid**, **EhrSystemId** (that you paste when registering the user) with copied information 
Click **Execute** 

**Result**: in the response, you see composition by version id

[Watch video instruction](https://media.bsn.si/ipehr/v2/get_composition.mp4) 📹 

### Update composition by version id

Precondition: [EHR was created](https://github.com/bsn-si/IPEHR-gateway/wiki/IPEHR-project#create-a-general-ehr) and **ehr_id**, **AuthUserId**, and **Prefer** (return=representation) are copied to a buffer and [Composition is created](https://github.com/bsn-si/IPEHR-gateway/wiki/IPEHR-project#update-composition-by-version-id)
1. Choose [Update composition by version id](https://gateway.ipehr.org/swagger/index.html#/COMPOSITION/put_ehr__ehr_id__composition__versioned_object_uid_) method and  click **Try it out**
1. Fill in **ehr_id**, **Authorization** (you Bearer and access token), **AuthUserId** (your userID), **versioned_object_uid**, **EhrSystemId** (that you paste when registered the user), **If-Match** and **Prefer** with copied information 
1. In the **Request** change information and click **Execute** 

**Result**: in the response, you see an updated composition

[Watch video instruction](https://media.bsn.si/ipehr/v2/update_composition.mp4) 📹 

### Delete composition by version id

Precondition: [EHR was created](https://github.com/bsn-si/IPEHR-gateway/wiki/IPEHR-project#create-a-general-ehr) and **ehr_id**, **AuthUserId**, and **Prefer** (return=representation) are copied to a buffer and [Composition is created](https://github.com/bsn-si/IPEHR-gateway/wiki/IPEHR-project#update-composition-by-version-id)
1. Choose [Delete composition by version id](https://gateway.ipehr.org/swagger/index.html#/COMPOSITION/delete_ehr__ehr_id__composition__preceding_version_uid_) method and  click **Try it out**
1. Fill in **ehr_id**, **Authorization** (you Bearer and access token), **AuthUserId** (your userID), **preceding_version_uid**, **EhrSystemId** (that you paste when registered the user) with copied information 
1. Click **Execute** 

**Result**: composition by version id is deleted. You can additionally check it by trying to [Get composition by id](https://github.com/bsn-si/IPEHR-gateway/wiki/IPEHR-project#get-composition-by-version-id)

[Watch video instruction](https://media.bsn.si/ipehr/v2/delete_composition.mp4) 📹 

## Execute ad-hoc (non-stored) AQL query

### Execute AQL query
1. Choose [Execute ad-hoc (non-stored) AQL query](https://gateway.ipehr.org/swagger/index.html#/QUERY/post_query_aql) method and click **Try it out**
1. Fill in **Authorization** (you Bearer and access token), **AuthUserId** (your userID)
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

### How to set and get access to the document

## How to get access to a document
Precondition: Open [swagger](https://gateway.ipehr.org/swagger/index.html) and [login](https://github.com/bsn-si/IPEHR-gateway/wiki/IPEHR-project#hou-to-login-under-the-user) under registered user
1. Choose [Get a document access list](https://gateway.ipehr.org/swagger/index.html#/ACCESS/get_access_document_) method and click **Try it out**
1. Fill in **Authorization** (you Bearer and access token), **AuthUserId** (your userID), and click **Execute**

**Result**: in the response you see list of documents

[Watch video instruction](https://media.bsn.si/ipehr/v2/get_document.mp4) 📹 

## How to set user access to the document
Precondition: Open [swagger](https://gateway.ipehr.org/swagger/index.html) and [login](https://github.com/bsn-si/IPEHR-gateway/wiki/IPEHR-project#hou-to-login-under-the-user) under registered user
1. Choose [Set user access to the document](https://gateway.ipehr.org/swagger/index.html#/ACCESS/post_access_document) method and click **Try it out**
1. Fill in **Authorization** (you Bearer and access token), **AuthUserId** (your userID). In Request fill in the access level (e.g owner, admin, read, noAccess), cid (paste from created ehr) and userID (user that you want to set access, then click **Execute**

**Result**: the request to change the level of access to the document was successfully created

[Watch video instruction](https://media.bsn.si/ipehr/v2/set_document.mp4) 📹 

