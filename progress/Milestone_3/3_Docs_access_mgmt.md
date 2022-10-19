## Design of document access rights management

### Documents

From the point of view of a smart contract, a document is a structure that contains a certain set of attributes:

```
struct DocumentMeta {
    CID         []byte
    ...
    docGroups   [][32]byte
    userGroups  [][32]byte
}
```

`docGroups` is a list of identifiers of document groups that include this document.  
`userGroups` is a list of identifiers of user groups that have access to this document.

### Document groups

Documents can be grouped by arbitrary criteria. From medical classification to geographical location.

```
struct DocumentGroup {
	ID                [32]byte
	owner             [32]byte
	description        string
	documents          [][]byte
	userGroupsAccess   [][32]byte
}
```

`documents ` is a list of documents contained in this group.  
`userGroups` is a list of user groups that have acces to this document group.

Document groups are stored in [IPEHR smart contract](https://github.com/bsn-si/IPEHR-blockchain-indexes)

```
mapping (bytes32 => DocumentGroup) docGroups;
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

- userGroupAddUser - Add a user to a group
- userGroupGetUsers - Get a list of users included to the group
- docGroupAddDoc - Add a document to a group
- docGroupGetDocs - Get a list of documents included in the group
- grantAccessUserToDoc - Allow the user to access the document
- grantAccessUserToDocGroup - Allow user access to a group of documents
- grantAccessUserGroupToDoc - Allow a group of users to access the document
- grantAccessUserGroupToDocGroup - Allow a group of users access to a group of documents
- getDocPermissions - Get a list of users and groups that have access to the document
- getUserPermissions - Get a list of documents to which the user has access
