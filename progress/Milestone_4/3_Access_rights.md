## Integration with access rights management system.

In stages [MS3.3](https://github.com/bsn-si/IPEHR-gateway/blob/develop/progress/Milestone_3/3_Docs_access_mgmt.md) and [MS3.4](https://github.com/bsn-si/IPEHR-gateway/blob/develop/progress/Milestone_3/4_Access_store.md) we developed a document access rights management system.  

At this stage we have implemented [API](https://gateway.ipehr.org/swagger/index.html#/) methods for:

- creating user groups
- adding/removing users from groups
- getting information about user groups
- getting information about user access
- getting information about access of groups
- delegation of access to documents

The mechanism of creating groups of documents and delegating access rights to these groups was also implemented.

In particular - when registering a patient, the `All documents` group is automatically created, and all documents related to this patient's EHR are automatically added to it.
