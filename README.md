# IPEHR

[![golangci-lint](https://github.com/bsn-si/IPEHR-gateway/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/bsn-si/IPEHR-gateway/actions/workflows/golangci-lint.yml)

## Disclaimer

The project is under active development and will gradually be supplemented.

## Description

Today common HMS applications store patients’ data in a local or a cloud DB which creates significant security, reliability and operational risks. Centralized storage and access rights administration of sensitive medical data create additional challenges:

-	Administrative overhead due to the rights provisioning on per patients/per record level.
-	Patients lack control and visibility over who has access to their data which goes against natural data subject rights announced in GDPR (General Data Protection Regulation, chapter 3).
-	Superuser access for DB and LDAP (access rights catalogue) create additional security risks.
-	In case of a data breach full registry will be compromised.

The IPEHR (InterPlanetary EHR) project is held to propose an alternative way of storing the data registry. Patients’ data will be stored in Filecoin network and will be accessed directly by stakeholders in case they have proper rights. Access rights and documents’ indexes will be stored on a blockchain in a smart-contract. Every data subject will have full unalienable control over his data and manage access rights on a personal level.
 
<p align="center">
  <img width="75%" src="https://user-images.githubusercontent.com/8058268/174096015-89aad056-d507-4ea5-8d00-bfef29c4a548.svg">
</p>

### Watch video introduction:

[![Watch the video](https://media.bsn.si/ipehr/logo_intro.jpg)](https://media.bsn.si/ipehr/introduction.mp4)

### Key features of the IPEHR solution:

- all data is stored in a decentralized storage;
- data encryption;
- self-management of user’s access rights;
- data integrity and authenticity is guaranteed by a smart contract.

## Development roadmap

This work is being done under the FileCoin development grant program RFP. See our proposal [here](https://github.com/filecoin-project/devgrants/issues/418)

The solution is being developed with 7 milestones:
* Development of MH-ORM and structure of storage of medical data - **completed**
* The functionality of encryption and saving/reading personal data to/from Filecoin - **completed**
* Access rights management on a blockchain - **completed**
* BsnGateway. Implementation of OpenEHR API, integration with MH-ORM  - **delivery**
* Public data publishing features using the Chainlink network - **completed**
* Application to manage your own medical data and access - **delivery**
* Testing, documentation and deployment

### Milestone 1

On Milestone 1 we've developed the IPEHR-gateway to provide benefits of decentralized architecture to common HMS solutions using standard APIs.
 
<p align="center">
  <img width="75%" src="https://user-images.githubusercontent.com/98888366/170698968-56ee7efe-e882-4236-b170-e9680ea12135.png">
</p>

#### IPEHR-gateway features:

- generates user’s cryptographic keys;
- exchanges medical data with HMS in openEHR format;
- provides filtering and search functions by indexing received openEHR documents;
- encrypts openEHR docs and indexes;
- stores encrypted openEHR documents in the FileCoin decentralized network;
- stores encrypted documents’ indexes in a smart contract on a blockchain;
- sends decrypted documents back to HMS;
- supports AQL queries without decrypting data.

### Milestone 2

On Milestone 2 we have developed indexes storage, access revocation algorithm and integrated the IPEHR-gateway with the Filecoin network.

![smart-contract](https://user-images.githubusercontent.com/8058268/190085702-6edf9437-1273-4db3-a7c9-414f66afe823.svg)

#### The main functions of a smart contract are

1. Search for `ehr_id` by user ID
2. Obtaining a list of documents related to the specified `ehr_id`.
3. Getting meta-information of the document by `ehr_id` and `document_uid`.
4. Getting an access key to the document.
5. Management of access to the document.
6. Document search using [AQL](https://specifications.openehr.org/releases/QUERY/latest/AQL.html) queries.

User and EHR document data in the contract is stored encrypted and prevents unauthorized persons from accessing private information.

For more information see [here](https://github.com/bsn-si/IPEHR-gateway/blob/develop/progress/Milestone_2/2_Indexes_storage.md).

#### Access rights management

![doc_access](https://user-images.githubusercontent.com/8058268/190620811-fd433f0b-44b7-4e04-a425-d77f62b55835.svg)

To grant access to a document, the document access key is asymmetrically encrypted with the public key of the user (or group) being granted access and added to the IPEHR smart contract table.

For more information see [here](https://github.com/bsn-si/IPEHR-gateway/blob/develop/progress/Milestone_2/3_Revoking_access.md).

### Milestone 3

From the point of view of a smart contract, a document is a structure that contains a certain set of attributes:

```
struct DocumentMeta {
    DocType   docType;
    DocStatus status;
    bytes     CID;
    bytes     dealCID;
    bytes     minerAddress;
    bytes     docUIDEncrypted;
    bytes32   docBaseUIDHash;
    bytes32   version;
    bool      isLast;
    uint32    timestamp;
}
```

At this point we will distinguish three levels of access: **Owner**, **Admin**, **Read**

Access to documents is managed according to the following access matrix:

|  Who \ Whom  | Owner |       Admin       |       Read        |
|     :---:    | :---: |       :---:       |      :---:        |  
|     Owner    |   no  | grant<br>restrict | grant<br>restrict |
|     Admin    |   no  |        grant      | grant<br>restrict |
|     Read     |   no  |        no         |        no         |

List of methods:  

- userGroupCreate - Creates a group of users
- groupAddUser - Adds a user to a group
- groupRemoveUser - Removes a user from a group
- docGroupCreate - Creates a group of documents
- docGroupAddDoc - Add a document to a group
- docGroupGetDocs - Get a list of documents included in the group
- setDocAccess - Sets the level of user access to the specified document
- getUserAccessList - Get a list of documents to which the user has access

For more information see [Milestone 3 repository](https://github.com/bsn-si/IPEHR-gateway/tree/develop/progress/Milestone_3)

### Milestone 4

AQL queries are used to search and filter data.

Archetype Query Language (AQL) is a declarative query language developed specifically for expressing queries used for searching and retrieving the clinical data found in archetype-based EHRs.

When new EHR documents are created, the homomorphically encrypted data they contain is placed in a special DataSearch index tree structure located in a blockchain. Indexes are searched using a smart contract. Later selections can be made from this data using AQL queries.

<p align="center">
  <img width="75%" src="https://user-images.githubusercontent.com/8058268/174270324-1218d6ba-4cf5-497d-b455-cb084b129141.svg">
</p>

You can find a detailed description of the AQL specification on the openEHR website: <https://specifications.openehr.org/releases/QUERY/latest/AQL.html>

Client applications use AQL data queries, this way platform can be used by the 3rd party applications supporting the standard.

In the process of implementing the ability to execute AQL queries, the following tasks were performed:

1. Index data repository developed: [https://github.com/bsn-si/IPEHR-gateway/tree/develop/src/pkg/storage/treeindex](https://github.com/bsn-si/IPEHR-gateway/tree/develop/src/pkg/storage/treeindex)
2. Smart contract for storing index data: [https://github.com/bsn-si/IPEHR-blockchain-indexes/blob/develop/contracts/DataStore.sol](https://github.com/bsn-si/IPEHR-blockchain-indexes/blob/develop/contracts/DataStore.sol)
3. AQL querier(*) - library for parsing an AQL query: [https://github.com/bsn-si/IPEHR-gateway/tree/develop/src/pkg/aqlprocessor](https://github.com/bsn-si/IPEHR-gateway/tree/develop/src/pkg/aqlprocessor)
4. AQL executor(*) - service that searches the index tree for data according to an AQL query: [https://github.com/bsn-si/IPEHR-stat/tree/main/internal/aqlquerier](https://github.com/bsn-si/IPEHR-stat/tree/main/internal/aqlquerier)

On receipt of a request, the IPEHR gateway interprets the request into a set of conditions, which is used to search the DataSearch index structure and returns the result to the requestor as specific values or as links to documents containing the requested data. 

For more information see [Milestone 4 repository](https://github.com/bsn-si/IPEHR-gateway/tree/develop/progress/Milestone_4)

### Milestone 5

On Milestone 5 we've developed a service that pushes public statistical data collected by the ipEHR system from stored EHRs to Chainlink network. The service implements an open API with specified metrics. The data is collected and processed by [IPEHR-blockchain-indexes](https://github.com/bsn-si/IPEHR-blockchain-indexes) smart contracts.

EHR stats can be collected by the service through:

-   transaction analysis of contracts [IPEHR-blockchain-indexes](https://github.com/bsn-si/IPEHR-blockchain-indexes)
-   periodic direct invocation of contract methods [IPEHR-blockchain-indexes](https://github.com/bsn-si/IPEHR-blockchain-indexes)
-   making AQL queries to IPEHR-gateway

We have implemented two types of statistical data delivery:

-   Direct delivery. Is implemented as a task in Chainlink that receives requests from outside by listening to the Oracle contract. When the Consumer contract sends a request for the statistical data to Oracle, the job collects a small fee in Chainlink tokens and returns the result from the statistics server. For this case, we have an open API for statistics and documented schema of job for Chainlink. The contract for this request is not a library but a sample contract in which we request statistical data via Oracle from the certain job of Chainlink

-   Scheduled delivery. It consists of two contracts and the schema of a Chainlink job. This job automatically requests statistical data within the specified interval and sends it to a storage contract. All other external contracts can request statistical data from the storage contract. The implementation includes two contracts. The first is the storage itself. We publish it and pay for its updates. The second is the contract of a Consumer. It is a sample of a simple contract that requests statistical data from storage.

For more information see [Milestone 5 repository](https://github.com/bsn-si/IPEHR-gateway/tree/develop/progress/Milestone_5)

### Milestone 6

We've developed an application (mobile/web) helping users to manage their personal data and its’ access rights. To control access to EHR Documents the following smart contracts are used: [EhrIndexer](https://github.com/bsn-si/IPEHR-blockchain-indexes/blob/develop/contracts/EhrIndexer.sol), [Users](https://github.com/bsn-si/IPEHR-blockchain-indexes/blob/develop/contracts/Users.sol), [AccessStore](https://github.com/bsn-si/IPEHR-blockchain-indexes/blob/develop/contracts/AccessStore.sol)

For each Patient during registration, a group of documents `All documents` and a group of users `Doctors` are created.


-   All new documents of the Patient are assigned to the group `All documents`.
    
-   All members of the `Doctors` group have access to the documents assigned to the `All documents' group.
    
-   The Patient can add Doctors to the `Doctors` group and they automatically get access to all Patient's Documents.
    
-   When a Doctor is removed from the 'Doctors' group, their access to document keys is terminated.
    

To read more info about access rights management, please visit the ["docs access management"](https://github.com/bsn-si/IPEHR-gateway/blob/develop/progress/Milestone_3/3_Docs_access_mgmt.md) section.


To ensure the quality of the application we developed and showed a test case:
[![Дизайн без названия (1)](https://user-images.githubusercontent.com/98888366/214616759-e0c84f22-b524-4879-acdd-68b81e775676.png)](https://media.bsn.si/ipehr/v2/how_to_add_doctor_into_app.mp4)

## How to

### Install Prerequisites

### Go 
Please follow installation instructions provided [here](https://go.dev/doc/install).

### IPFS
The ipEHR gateway requires a connection to the IPFS network. You can use a third-party service or install your own node. Installation instructions can be found [here](https://github.com/ipfs/kubo#install).

### Filecoin
The ipEHR gateway requires a connection to the Filecoin network. It is necessary to install your own Lotus instance. You can run it in either full node mode or [light](https://lotus.filecoin.io/tutorials/lotus/store-and-retrieve/set-up/#install-a-lite-node) mode. Installation instructions can be found [here](https://lotus.filecoin.io/lotus/install/prerequisites/).

Enable support for fetching data from IPFS before launching in the config `~/.lotus/config.toml`
```
[Client]
UseIpfs = true
```

Creating a directory to retrieve files and then transfer them to the IPEHR-gateway
```
mkdir -p $LOTUS_DIR/files
```

Install Nginx to retrieve files on the IPEHR-gateway.

/etc/nginx/conf.d/lotus.conf:
```
server {
        server_name <HOSTNAME>;

        location /files {
                alias <LOTUS_DIR/files;
                add_header Content-disposition "attachment; filename=$1";
                default_type application/octet-stream;
        }

        location / {
                proxy_pass      http://127.0.0.1:1234;
        }
}
```

Replace \<HOSTNAME\> and <LOTUS_DIR> with your values.


### Clone this repo

```
git clone https://github.com/bsn-si/IPEHR-gateway
```

### Run Tests

```
cd ./src
go test -v ./...
```

### Build IPEHR-gateway

```
cd ./src
go build -o ../bin/ipehr-gateway cmd/ipehrgw/main.go
```

### Run IPEHR-gateway

```
./bin/ipehr-gateway -config=./config.json
```

### Get swagger UI API documentation

[Swagger UI API docs](http://gateway.ipehr.org/swagger/index.html)

The following methods are based on the [standard specification of OpenEHR](https://specifications.openehr.org/releases/ITS-REST/Release-1.0.2/ehr.html)

## Workflow example

You can use [this](https://github.com/bsn-si/IPEHR-gateway/wiki/IPEHR-project#workflow-example) instruction and try to create, update or get EHR
For now we have the following methods:

- Register a user
- Log in under the or log out
- Refresh JWT token
- Create an EHR (also with exact id or different parameters)
- Getting info on created summary EHR by subject id 
- Getting info on created summary EHR by summary id 
- Getting info on created EHR status
- Getting info on EHR status version by time
- Getting info on EHR by request id
- Update EHR status with id
- Create composition
- Get composition
- Update composition
- Delete composition
- Create group access
- Get group access
- Execute AQL request
- Get a document access list
- Set user access to the document

## Docker
You can start a project in Docker

Before building the image, you need to create config.json (take config.json.example as a basis)

The repository contains a Dockerfile
At the root of the project run the command 
```
docker build -t ipehr:gtw .
```

After building the image, start the container
```
docker run -d --restart always -p 8080:8080 --name ipehr-gateway ipehr:gtw
```

### Related repositories
A [mobile app](https://github.com/bsn-si/IPEHR-access-control-app) for access management to user's (patient's) EHRs.
A [Chainlink publisher](https://github.com/bsn-si/IPEHR-stat) for public EHR stats.
An [ipEHR smart contract](https://github.com/bsn-si/IPEHR-blockchain-indexes) for access management, indexes storage and obtaining EHR stats.
