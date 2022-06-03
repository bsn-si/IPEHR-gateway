## Document storage

EHR documents are supposed to be stored as separate files.
To store documents containing medical data, a repository with the following characteristics is required:

- fault tolerance
- high availability
- scalability (ability to store large amounts of data)
- high speed of data access
- low cost of data storage

This project will use [Filecoin](https://filecoin.io), which meets the above requirements, as the EHR document repository.

The implementation of saving documents in Filecoin will be implemented in the next milestone.

As part of this phase, to speed up the development process, the saving of documents is done in regular files to the local file storage on the IPEHR gateway, which emulates the work of the decentralized file storage.

When saving a file, the repository returns an identifier that can be used to access the file at a later time.

## Implementation

The packages that implement the localfile storage functionality are located in the `pkg/storage/localfile` project directory.

Running tests to demonstrate the work of the packages:

```
go test -v ./pkg/storage/localfile
```