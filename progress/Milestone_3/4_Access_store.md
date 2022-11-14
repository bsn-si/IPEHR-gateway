### Design and development of the encrypted storage to store data access keys

Methods to verify the access rights and acquire the access keys.

```
enum AccessLevel { NoAccess, Owner, Admin, Read }
enum AccessKind { Doc, DocGroup, UserGroup }

struct Access {
	bytes32 idHash;
	bytes idEncr;
	bytes keyEncr;
	AccessLevel level;
}


mapping(bytes32 => Access[]) accessStore; // accessID => Object[]
```

Upon the creation of a document, a unique symmetrical access key is generated. It is used to encrypt the document.

You can find more about the encryption of documents [here](https://github.com/bsn-si/IPEHR-gateway/tree/develop/progress/Milestone\_1/3\_Encryption#readme)


### Access key storage

Access Key Storage is a hash table that is located in the smart contract.  
The key of the table is a special 32 byte identification number known as `accessID`.  
The value is an array of `Access` type objects.

### Adding an access key to a document in the repository

Adding an access key to the repository is the same as granting access to the document.

If the document already exists in the system, only a user with the `Owner` or `Admin` access level rights can add an access key to this document.

To add an access key to the repository, one needs to:

- calculate accessID = keccak256(abi.encode(userID, AccessKind.Doc)),
- add an Access object to the repository at accessID: address,
- idHash - document ID hash,
- idEncr - document ID encrypted with the access key,
- keyEncr - object access key encrypted with the public key of the user who created the document,
- level - document access level value, one of (Owner, Admin, Read).

### Gaining a document access key

To gain an access key to a document, one needs to:

- calculate accessID = keccak256(abi.encode(userID, AccessKind.Doc)),
- find an array of user documents in the repository,
- find Access object by document idHash,
- for each document decrypt Access.keyEncr (encrypted access key) with the user's private key.

If the user does not have access to the document, there will be no entry corresponding to the Access object in the repository.

Or if the user does not have a private key, he will not be able to decrypt the access key.

### Gaining a list of documents to which the user has access

To get a list of document IDs:

- calculate accessID = keccak256(abi.encode(userID, AccessKind.Doc))
- find an array of user documents in the repository,
- for each document decrypt Access.keyEncr (encrypted access key) with the user's private key,
- then use the access key to decrypt Access.idHash,
- add docID to the list.
