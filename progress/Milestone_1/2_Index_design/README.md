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
    groupId        [16]byte   // UUID
    value          []byte     // encrypted values
    docStorIdEncr  []byte     // user_pub_key.Encrypt(doc_storage_id)
}

pathKey = sha3(path)    // path Example: '/data/events[at0006]/data/items[at0004]/value/magnitude'
pathKey  -> DataEntry
```

### dataAccess index

`access_group` - access group, within which encrypted data is searched

```
key = sha3(user_id+access_group_id)
value = user_pub_key.Encrypt(access_group_key)

key -> value
```

## Index storage

To store the indexes, it is planned to create a special smart contract on the blockchain. This will ensure security, data immutability and fast access.

## Encrypting data in indexes

Searching through the encrypted data set requires that the data is encrypted with one key. The notion of access_group is introduced.  
In practice, any group of people based on some administrative boundaries (department, department, clinic, etc.) can be combined into a group. For the access group, the encryption key is generated - 32 bytes + unique access_group_nonce - 12 bytes, which are kept secret within the group.

### Encrypting text values

Text values are encrypted using the ChaCha20-Poly1305 algorithm with a fixed access_group_nonce before being saved into a public index

Example:

```
access_group_key: 1d07add12296f142bf730c4871d5c4eaf4652d5b0b213a1c700e2a3437a2c1e0
access_group_nonce: 76b1092c73a1d28a7c9eb7b2
message: "Hello"
auth_data: "Hello"
encrypted: 2802e22fd71907ecfca353fee2890f45fd326210d9
```

Using the same `access_group_nonce` makes it possible to search ciphertexts without decrypting them.  
**Limitation**: you must know the exact case-sensitive value to search.  
To do a keyword search, you can break the phrase into separate words, encode them individually, and put them in an index.

### Encryption of numeric values 

To convert (encrypt) numeric values we will use the function `F(x) = aX + b`, where a and b are positive integers > 0. 

Warning! When generating the `access_group_key`, you have to make sure that the first 8 bytes are not zero.  

As the a and b coefficients, we take the first 8 bytes of the `access_group_key` (4 bytes each)

Example:
```
access_group_key: 1d07add12296f142bf730c4871d5c4eaf4652d5b0b213a1c700e2a3437a2c1e0
a = 0x1d07add1 // 487042513
b = 0x2296f142 // 580317506
value = 36.6
encrypted = 18406073481.8 // 36.6 * 487042513 + 580317506
```

This conversion scheme makes it possible to compare numbers, search through ranges, and perform arithmetic operations on the converted values, while only a member of the access group who knows the key can access the real data.

Example:
You need to search for records where the temperature is > 39 C
For the query we first "encrypt" the required value 39.0 * 487042513 + 580317506 = 19574975513

```
SELECT value FROM data_index 
WHERE type = 'temperature' AND value > 19574975513
```

## Creating EHR

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
5. Path construction for all field values within the document
6. Determining the `access_group` for searching the values in the dataSearch index
7. Records containing paths, values encrypted with `access_group_key` and `encrypted doc_storage_id` to associate records with documents in the repository and `access_group_id` are added to the dataSearch Index for all values

## Search for documents using queries

1. Paths to the values of interest and conditions are generated from the query
2. Search for records in dataSearch index, which satisfy the specified conditions, the encrypted `doc_storage_id` of documents containing these values is returned, or the encrypted values
3. Search for an `access_group_key` in the dataAccess index
4. Decrypt `access_group_key` with `user_priv_key`
5. Decrypting field values with `access_group_key`
6. Decrypt `doc_storage_id` with `user_priv_key`
7. Search for `doc_keys` encrypted documents in the docAccess index by `doc_storage_id`
8. Decrypt `doc_keys`
9. Downloading found documents from the storage
10. Documents are decrypted using `doc_keys`

## Implementation

At this stage, the implementation of the indexes is done with their storage in the local file storage plus caching in memory.

The next stage involves the development of a special blockchain-based smart contract for working with indexes.

The packages that implement the indexing functionality are located in the `pkg/indexer` project directory.

Running tests to demonstrate the work of the packages:

```
go test -v ./pkg/indexer/...
```
