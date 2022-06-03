![HMS gateway](https://user-images.githubusercontent.com/8058268/170997404-c1a20845-a7c5-4663-a291-9f088c0d05ae.png)

IPEHR-gateway implements functions for creating and updating documents according to openEHR standards.

The EHR Information Model version 1.1.0, the latest stable version at the time of development, was used.

<https://specifications.openehr.org/releases/RM/Release-1.1.0/ehr.html>

To create and update EHR documents, we implemented methods from the openEHR REST API specification

<https://specifications.openehr.org/releases/ITS-REST/Release-1.0.2>

The following methods are currently implemented:

- Create EHR
- Get EHR summary by id
- Get EHR_STATUS by version id
- Update EHR_STATUS

## How it works
1. IPEHR gateway receives a request via HTTP API to create a new EHR document or update the document status. The request header specifies the ID of the calling user
2. IPEHR gateway generates a key pair for the user - public and private, if they were not generated earlier and stores them in the key store
3. the unique key is generated for encryption of the document
4. the document is encrypted using a symmetric algorithm with the key from the previous step
5. the encrypted document is stored in the storage and the document storage id is generated
6. document index is created, which contains information about the document;
7. now the user can get the created EHR document using the corresponding HTTP API request
8. user can also update EHR_STATUS

## Implementation

The document structures are located in the directory `pkg/docs/model`

Running tests that simulate creating and receiving documents:

```
go test -v ./pkg/api/...
```

