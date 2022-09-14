![doc_access](https://user-images.githubusercontent.com/8058268/190081743-c2beb7e8-24d2-4fb4-8c7b-31506c2f62b9.svg)

## Algorithm of changing document access rights

As described in [Milestone 1.2](https://github.com/bsn-si/IPEHR-gateway/tree/develop/progress/Milestone_1/2_Index_design), each EHR document is symmetrically encrypted with a unique access key.

User access to documents is controlled by the `docAccess` smart-contract index.

EHR data can be accessed sub-documentally.

To grant access to a document, the document access key is asymmetrically encrypted with the public key of the user (or group) being granted access and added to the IPEHR smart contract table.

Thus, a user with a private key can decrypt the access key to the document and to the document itself.

To revoke access to a document, a user with the rights of the document owner deletes a record in the table. At this point, the user for whom access was restricted will no longer be able to obtain the document's access key.

## Implementation

The described functionality for document access rights management will be implemented as part of the 3rd milestone.
