### An algorithm for re-encrypting data in case of revocation or change of rights for a single user or a group.

During the development of the document storage system, we rejected the idea of re-encrypting documents for two reasons:

1. Filecoin does not currently support the function of deleting previously saved data.
2. In case a user got access to a document, it is impossible to guarantee the fact that the decrypted document is not stored in his local archive.

### Solution

As described in [Milestone 1.2](https://github.com/bsn-si/IPEHR-gateway/tree/develop/progress/Milestone_1/2_Index_design), each document is encrypted with a unique access key.  
User access to documents is controlled by the `docAccess` index.  
On revocation of access to a document, a corresponding change will be made to the `docAccess` index on the smart contract, and the user will not be able to get an access key to the document from that moment.

The described functionality will be implemented in Milestone 3.
