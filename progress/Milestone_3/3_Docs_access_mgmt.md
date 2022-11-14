## Design of document access rights management

### Documents

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

### Document groups

Documents can be grouped by arbitrary criteria. From medical classification to geographical location.

Each group has unique access key which encrypts data within the group and which users who have access to this group of documents must keep in secret.

```
struct DocumentGroup {
    mapping(bytes32 => bool)   CIDHashes;
    mapping(bytes32 => bytes)  params;
    bytes[]                    CIDEncrs;   // CIDs encrypted with the group access key 
    bytes32[]                  userGroups;
}
```

Document groups are stored in [IPEHR smart contract](https://github.com/bsn-si/IPEHR-blockchain-indexes)

```
mapping (bytes32 => DocumentGroup) docGroups;   // groupIdHash => DocumentGroup
```

### Access level

At this point we will distinguish three levels of access: **Owner**, **Admin**, **Read**

Access to documents is managed according to the following access matrix:

|  Who \ Whom  | Owner |       Admin       |       Read        |
|     :---:    | :---: |       :---:       |      :---:        |  
|     Owner    |   no  | grant<br>restrict | grant<br>restrict |
|     Admin    |   no  |        grant      | grant<br>restrict |
|     Read     |   no  |        no         |        no         |

### Access control

List of methods:  

- userGroupCreate - Creates a group of users
- groupAddUser - Adds a user to a group
- groupRemoveUser - Removes a user from a group
- docGroupCreate - Creates a group of documents
- docGroupAddDoc - Add a document to a group
- docGroupGetDocs - Get a list of documents included in the group
- setDocAccess - Sets the level of user access to the specified document
- getUserAccessList - Get a list of documents to which the user has access


