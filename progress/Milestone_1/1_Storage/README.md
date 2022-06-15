![HMS gateway](https://user-images.githubusercontent.com/8058268/171821436-ebd013b6-0deb-4f86-8aaa-b5254e913104.png)

## openEHR documents

![EHR Information model](https://specifications.openehr.org/releases/RM/latest/ehr/diagrams/high_level_ehr_structure.svg)

According to the openEHR information model, there are 6 basic types of documents:

- EHR: the root object, identified by a globally unique EHR identifier;
- EHR_access: an object containing access control settings for the record;
- EHR_status: an object containing various status and control information, optionally including the identifier of the subject (i.e. patient) currently associated with the record;
- Folders: optional hierarchical folder structures that can be used to logically index Compositions;
- Compositions: the containers of all clinical and administrative content of the record;
- Contributions: the change-set records for every change made to the health record; each Contribution references a set of one or more Versions of any of the versioned items in the record that were committed or attested together by a user to an EHR system.

## Document storage

When you create new documents, they are saved in the IPEHR system as separate files in encrypted form.  
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
