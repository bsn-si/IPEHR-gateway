# IPEHR

## Disclamer

The project is under active development and will gradually be supplemented.

## Description

Today common HMS applications store patients’ data in a local or a cloud DB which creates significant security, reliability and operational risks. Centralized storage and access rights administration of sensitive medical data creates additional challenges:

-	Administrative overheard due to the rights provisioning on per patients/per record level.
-	Patients lack control and visibility over who has access to their data which goes against natural data subject rights announced in GDPR (General Data Protection Regulation, chapter 3).
-	Super user access for DB and LDAP (access rights catalogue) create additional security risks.
-	In case of a data breach full registry will be compromised.

The IPEHR (InterPlanetary EHR) project is held to propose an alternative way of storing the data registry. Patients’ data will be stored in Filecoin network and will be accessed directly by stakeholders in case they have proper rights. Access rights and documents’ indexes will be stored on a blockchain in a smart-contract. Every data subject will have full unalienable control over his data and manage access rights on a personal level.
 
<p align="center">
  <img width="75%" src="https://user-images.githubusercontent.com/8058268/174096015-89aad056-d507-4ea5-8d00-bfef29c4a548.svg">
</p>

### Whatch video introduction:

[![Watch the video](https://media.bsn.si/ipehr/logo_intro.jpg)](https://media.bsn.si/ipehr/introduction.mp4)

### Key features of the IPEHR solution:

- all data is stored in a decentralized storage;
- data encryption;
- self management of user’s access rights;
- data integrity and authenticity is guaranteed by a smart contract.

## Development roadmap

This work is being done under the FileCoin development grant program RFP. See our proposal [here](https://github.com/filecoin-project/devgrants/issues/418)

The solution is being developed with 7 milestones:
* Development of MH-ORM and structure of storage of medical data - **we are here**
* The functionality of encryption and saving/reading personal data to/from Filecoin
* Access rights management on a blockchain
* BsnGateway. Implementation of OpenEHR API, integration with MH-ORM
* Public data publishing features using the Chainlink network
* Application to manage your own medical data and access
* Testing, documentation and deployment

On Milestone 1 we develop the IPEHR-gateway to provide benefits of decentralized architecture to common HMS solutions using standard APIs.
 
<p align="center">
  <img width="75%" src="https://user-images.githubusercontent.com/98888366/170698968-56ee7efe-e882-4236-b170-e9680ea12135.png">
</p>

### IPEHR-gateway features:

- generates user’s cryptographic keys;
- exchanges medical data with HMS in openEHR format;
- provides filtering and search functions by indexing received openEHR documents;
- encrypts openEHR docs and indexes;
- stores encrypted openEHR documents in the FileCoin decentralized network;
- stores encrypted documents’ indexes in a smart contract on a blockchain;
- sends decrypted documents back to HMS;
- supports AQL queries without decrypting data.

### How to

## Install Prerequisites

Please follow installation instructions provided [here](https://go.dev/doc/install).

## Clone this repo

```
git clone https://github.com/bsn-si/IPEHR-gateway
```

## Run Tests

```
cd ./src
go test -v ./...
```

## Build IPEHR-gateway

```
cd ./src
go build -o ../bin/ipehr-gateway cmd/ipehrgw/main.go
```

## Run IPEHR-gateway

```
./bin/ipehr-gateway -config=./config.json
```

## Get swagger UI API documentation

[Swagger UI API docs](http://gateway.ipehr.org/swagger/index.html)

The following methods are based on the [standard specification of OpenEHR](https://specifications.openehr.org/releases/ITS-REST/Release-1.0.2/ehr.html)

## Workflow example

You can use [this](https://github.com/bsn-si/IPEHR-gateway/wiki/IPEHR-project#workflow-example) insctuction and try to create, update or get EHR
For now we have the following methods:

- Create an EHR (also with exact id or different parameters)
- Getting info on created summary EHR by subject id 
- Getting info on created summary EHR by summary id 
- Getting info on created EHR status
- Getting info on EHR status version by time
- Update EHR status with id
- Create composition
- Execute AQL request

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


