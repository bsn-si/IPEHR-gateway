# Design and development of an application to manage personal medical data and its’ access rights

## Application design compatible with web and mobile platforms


As described in the ["users identity"](https://github.com/bsn-si/IPEHR-gateway/blob/develop/progress/Milestone_3/2_Users_identity.md) section there are two roles used in the current implementation of ipEHR: a `Patient` and a `Doctor`.

The `user` is the following structure:

```
struct User {
    bytes32   id;
    bytes32   systemID;
    Role      role;
    bytes     pwdHash;
  }
```

User groups:

```
struct UserGroup {
    mapping(bytes32 => bytes) params;
    mapping(address => AccessLevel) members;
    uint membersCount;
}
```

Only a member with `Owner` or `Admin` access rights can add users to a group.


Users and user groups are stored in [IPEHR smart contract](https://github.com/bsn-si/IPEHR-blockchain-indexes)

```
mapping (address => User)      users;
mapping (bytes32 => UserGroup) userGroups;
```

Pre-registration is required to work with the IPEHR system.
  
 
 
To enable Patients to control access to their medical data we designed web and mobile applications.
  

The Patient should be able to:

-   Login into the application
    
-   See the list of available EHR Documents
    
-   See the list of Doctors who have access to Patient’s EHR Documents
    
-   Manage access of Doctors to personal EHR Documents

  

To ensure cross-platform compatibility we used PWA (progressive web app) technology.

  
## UI/UX design
  

We developed a set of [wireframes](https://miro.com/app/board/uXjVPwaRjjY=/?share_link_id=25084805149) and [specification](https://docs.google.com/document/d/1aIDZMmukk8Y0d_b_e3eD0CrBj9kCnacKaGxKX1iDer4/edit?usp=sharing) according to the planned functionality. Based on these artifacts a [design layout](https://www.figma.com/file/TSKmYUCG3pHjDGtkCoB62Y/%D0%BF%D1%80%D0%B8%D0%BB%D0%BE%D0%B6%D0%B5%D0%BD%D1%8C%D0%BA%D0%B0-ipehr?node-id=0%3A1&t=cDnGRxraxnS0eXoM-0) was developed.

  

## Frontend development.

To ensure cross-platform compatibility we used PWA (progressive web app) technology.
  

The only user of the application is the Patient.

Available functionality:

-   Log in / log out
    
-   Browsing the list of EHR Documents associated with user’s ipehr-id
    
-   Adding Doctors via pin-code
    
-   Adding Doctors via QR-code
    
-   Removing Doctors from the list of Doctors
    

For now all blockchain operations are carried out on the Backend side. The stack of the application does not include any Web3 technologies. It is built on Next.js, Nextauth. In Nextjs we use out-of-the-box React, Typescript, SCSS, server-side-rendering, and others.


## Blockchain integration to manage rights.

To control access to EHR Documents the following smart contracts are used: [EhrIndexer](https://github.com/bsn-si/IPEHR-blockchain-indexes/blob/develop/contracts/EhrIndexer.sol), [Users](https://github.com/bsn-si/IPEHR-blockchain-indexes/blob/develop/contracts/Users.sol), [AccessStore](https://github.com/bsn-si/IPEHR-blockchain-indexes/blob/develop/contracts/AccessStore.sol)

### Algorithm of changing document access rights

![doc_access](https://user-images.githubusercontent.com/8058268/190620811-fd433f0b-44b7-4e04-a425-d77f62b55835.svg)

In case of revoking access to an EHR, re-encrypting (and hence re-writing) all documents would be quite an inexpedient challenge, taking into account how much time and financial resources it will take. In addition, deleting files from FC is not currently implemented. The most reasonable method for revoking access is re-encrypting (hence re-writing) access keys!

Each EHR document is symmetrically encrypted with a unique access key. User access to documents is controlled by the `docAccess` smart-contract index.


You can get more information in the ["revoking access"](https://github.com/bsn-si/IPEHR-gateway/blob/develop/progress/Milestone_2/3_Revoking_access.md) section.

For each Patient during registration, a group of documents `All documents` and a group of users `Doctors` are created.


-   All new documents of the Patient are assigned to the group `All documents`.
    
-   All members of the `Doctors` group have access to the documents assigned to the `All documents' group.
    
-   The Patient can add Doctors to the `Doctors` group and they automatically get access to all Patient's Documents.
    
-   When a Doctor is removed from the 'Doctors' group, their access to document keys is terminated.
    

To read more info about access rights management, please visit the ["docs access management"](https://github.com/bsn-si/IPEHR-gateway/blob/develop/progress/Milestone_3/3_Docs_access_mgmt.md) section.

## Test case design and development + Testing


To ensure the quality of the application we developed a test case:

  
Precondition:

-   A Patient is registered (swagger POST /user/register)
    
-   A Doctor is registered (swagger POST /user/register)
    
-   Doctor’s code is acquired (swagger GET /user/:user_id)
    
-   Doctor’s QR-code is generated
    
  

Main flow:

-   Patient logs into the application
    
-   Patient checks the list to Documents
    
-   List of Documents is Empty
    
-   We create a new EHR in the Swagger (swagger POST /ehr)
    
-   Patient refreshes the list of Documents
    
-   List of Documents shows two Documents (ehr, ehr_status)
    
-   Patient checks the list of Doctors
    
-   List of Doctors is empty
    
-   Patient adds a Doctor ( via QR-code in the mobile version or pin code in the web version)
    
-   Patient refreshes the list of Doctors
    
-   List of Doctors shows the added Doctor
    
-   We use Swagger to check that the Doctor has access to the Patient’s Documents (GET /access/document/)
    
-   Patient removes the Doctor from the list of Doctors
    
-   Patient refreshes the list of Doctors
    
-   List of Doctors is empty
    
-   We use Swagger to check if the Doctor is stripped of access to the Patient’s Documents
    
  
  
Here is the demonstration video:
[![Дизайн без названия (1)](https://user-images.githubusercontent.com/98888366/214616759-e0c84f22-b524-4879-acdd-68b81e775676.png)](https://media.bsn.si/ipehr/v2/how_to_add_doctor_into_app.mp4)

