# IPEHR
## Disclamer
The project is under active development and will gradually be supplemented.

## Description
Today common HMS applications store patients’ data in a local or a cloud DB which creates significant security, reliability and operational risks. Centralized storage and access rights administration of sensitive medical data creates additional challenges:

- Administrative overheard due to the rights provisioning on per patients/per record level
- Patients lack control and visibility over who has access to their data which goes against natural data subject rights announced in GDPR (General Data Protection Regulation, chapter 3)
- Super user access for DB and LDAP (access rights catalogue) create additional security risks
- In case of a data breach full registry will be compromised

The IPEHR (InterPlanetary EHR) project is held to propose an alternative way of storing the data registry. Patients’ data will be stored in Filecoin network and will be accessed directly by stakeholders in case they have proper rights. Access rights and documents’ indexes will be stored on a blockchain in a smart-contract. Every data subject will have full unalienable control over his data and manage access rights on a personal level.

![image](https://user-images.githubusercontent.com/98888366/170699014-2ff3cec6-913b-4b4f-85f0-63899382ff24.png)

### Key features of the IPEHR solution:
- all data is stored in a decentralized storage;
- data encryption;
- self management of user’s access rights;
- data integrity and authenticity is guaranteed by a smart contract.

## Development roadmap
This work is being done under the FileCoin development grant program RFP. See our proposal [here](https://github.com/filecoin-project/devgrants/issues/418)

On this stage (Milestone 1) we develop the IPEHR-gateway to provide benefits of decentralized architecture to common HMS solutions using standard APIs.
 
![image](https://user-images.githubusercontent.com/98888366/170698968-56ee7efe-e882-4236-b170-e9680ea12135.png)

### IPEHR-gateway features:
- generates user’s cryptographic keys;
- exchanges medical data with HMS in openEHR format;
- provides filtering and search functions by indexing received openEHR documents;
- encrypts openEHR docs and indexes;
- stores encrypted openEHR documents in the FileCoin decentralized network;
- stores encrypted documents’ indexes in a smart contract on a blockchain;
- sends decrypted documents back to HMS;
- supports AQL queries without decrypting data.



## How to
### Install Prerequisites
Please follow installation instructions provided [here](https://go.dev/doc/install).

### Clone this repo
```
git clone https://github.com/bsn-si/IPEHR-gateway
```

### Run Tests
```
cd IPEHR-gateway/src
go test -v ./...
```
