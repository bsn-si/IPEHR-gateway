![HMS gateway](https://user-images.githubusercontent.com/8058268/171837801-269660f7-b274-4d97-ab17-30de41e3f962.png)

### User/patient

- user\_id - unique identifier  
- key pair - pub\_key + priv\_key, generated according to the Curve25519 elliptic curve-based cryptosystem

### EHR (electronic health record)
Is a data structure described in the openEHR standard

- ehr_id
- includes the documents: EHR\_STATUS, COMPOSITION, FOLDER, CONTRIBUTION

Each patient can have only one EHR.

## Document encryption

Each document that is added to the system is encrypted using the ChaCha20-Poly1305 algorithm.  
A separate key is generated for each document.

### EHR index

Allows you to find `ehr_id` from a specified `user_id`

user\_id -> ehr_id

### Document index

Allows you to find the `doc_storage_id` of all documents by the specified `ehr_id`

```
DocumentMeta: {
    DocType         uint8
    StorageId       [32]byte
    DocIdEncrypted  []byte
    Timestamp       uint32
}
```

ehr_id -> []DocumentMeta // here is an array

### EHRsubject index

Allows to find `ehr_id` using the specified `subject_id` and `subject_namespace`

```
subjectKey = sha3(subjectId+namespace)

subjectKey -> ehr_id
```

### docAccess index

Allows to find an encrypted `doc_key` for the specified `doc_storage_id` and `user_id`

```
key = sha3(doc_storage_id+user_id)
value = user_pub_key.Encrypt(doc_key)

key -> value
```

### dataSearch index

Allows to find the `doc_storage_id` of documents that contain data with the specified values

```
DataEntry: {
    index_id [16]byte
    value    []byte     // open values
}

pathKey = sha3(path)    // path Example: '/data/events[at0006]/data/items[at0004]/value/magnitude'

pathKey  -> DataEntry

index_id -> user_pub_key.Encrypt(doc_storage_id)
```

##Creating EHR

1. An EHR document is created as a json file
2. The document is encrypted using the unique `doc_key`
3. Document is saved in the document storage, `doc_storage_id` is returned
4. An entry is added to the EHR Index to link the patient and their EHR:
5. An entry is added to the Document Index to allow the document in the repository to be linked to the EHR
6. An entry is added to the EHRsubject Index to search for `ehr_id` by its subject
7. An entry is added to the docAccess Index to find the encrypted `doc_key` from the document

## Getting EHR by ehr_id

1. Get `ehr_id` using `user_id` in EHR Index
2. Get `doc_storage_id` of a document with DocType = EHR
3. Get the encrypted `doc_key` from the docAccess index
4. Decrypt `doc_key` with user's `priv_key`
5. Download document with `doc_storage_id` from storage
6. Decrypt the document using the decrypted `doc_key` in step 5

## Creating COMPOSITION

1. The document is created as a json file
2. The document is encrypted with the unique key `doc_key`
3. The document is saved in the document storage, `doc_storage_id` is returned
4. An entry is added to the document index, allowing you to link the document in the repository with the EHR
5. Paths are built for all values within the document
6. An entry is added to the dataSearch index for all values, containing paths, values in open form, and the encrypted `doc_storage_id` to link the entries to the documents in the storage

## Search for documents using queries

1. Paths to the values of interest and conditions are generated from the query
2. Search for records in dataSearch index that match the specified conditions, encrypted `doc_storage_id` of documents containing these values is returned
3. Decrypt `doc_storage_id` with `user_priv_key`
4. Search for `doc_keys` encrypted documents in the docAccess index by `doc_storage_id`
5. Decrypt `doc_keys`
6. Downloading found documents from the storage
7. Documents are decrypted using `doc_keys`

## Implementation

At this stage, the implementation of the indexes is done with their storage in the local file storage plus caching in memory.

The next stage involves the development of a special blockchain-based smart contract for working with indexes.

The packages that implement the indexing functionality are located in the `pkg/indexer` project directory.

Running tests to demonstrate the work of the packages:

```
go test -v ./pkg/indexer/...
```