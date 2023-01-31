# Milestone Delivery :mailbox:

* **Application Document:** [ipEHR application](https://github.com/filecoin-project/devgrants/issues/418)
* **Milestone Number:** 6

**Context**

In this milestone we've developed an application to manage access rights to personal medical data stored in the FIL network.

**Deliverables**

| Number                | Deliverable                              | Link                                                                                                                                  | Notes                                                                                                                                                                                                                                                                                                                  |
|-----------------------|------------------------------------------|---------------------------------------------------------------------------------------------------------------------------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| :heavy_check_mark:    | License | [LICENSE](https://github.com/bsn-si/IPEHR-stat/blob/main/LICENSE) | Apache 2.0 license |
| :heavy_check_mark:    | Testing Guide | [Readme.md](https://github.com/bsn-si/IPEHR-stat/tree/main#local-deployment) | The "How To" guide is supplemented with usage examples videos |
| 1.-3. :heavy_check_mark: | Application design (UI/UX, frontend development) compatible with with web and mobile platforms | See [design description](https://github.com/bsn-si/IPEHR-gateway/tree/develop/progress/Milestone_6#design-and-development-of-an-application-to-manage-personal-medical-data-and-its-access-rights) | Here you can find wireframes and specifications as well as layout design. | 
| 4. :heavy_check_mark: | Blockchain integration to manage rights | See [readme.md](https://github.com/bsn-si/IPEHR-gateway/tree/develop/progress/Milestone_6#blockchain-integration-to-manage-rights) | To control access to EHR Documents the following smart contracts are used: [EhrIndexer](https://github.com/bsn-si/IPEHR-blockchain-indexes/blob/develop/contracts/EhrIndexer.sol), [Users](https://github.com/bsn-si/IPEHR-blockchain-indexes/blob/develop/contracts/Users.sol), [AccessStore](https://github.com/bsn-si/IPEHR-blockchain-indexes/blob/develop/contracts/AccessStore.sol). For each Patient during registration, a group of documents `All documents` and a group of users `Doctors` are created. | 
| 5. :heavy_check_mark: | Test case design and development | See [readme.md](https://github.com/bsn-si/IPEHR-gateway/tree/develop/progress/Milestone_6#test-case-design-and-development--testing) | To ensure the quality of the application we've developed a test case described. |
| 6. :heavy_check_mark: | Testing | See [demonstration video](https://media.bsn.si/ipehr/v2/how_to_add_doctor_into_app.mp4) | Preconditions in the demonstration video: a Patient is registered; a Doctor is registered; doctor’s code is acquired; doctor’s QR-code is generated. |

## Project introduction:

[![Watch the video](https://media.bsn.si/ipehr/logo_intro.jpg)](https://www.youtube.com/watch?v=nJFA5W4qoEw)


## Workflow example

We have developed an application compatible with web and mobile platforms designed to manage access rights to personal EHR records/documents вещкув шт the FIL network.


## The test case:
  
  
### Preconditions:

-   A Patient is registered (swagger POST /user/register)
    
-   A Doctor is registered (swagger POST /user/register)
    
-   Doctor’s code is acquired (swagger GET /user/:user_id)
    
-   Doctor’s QR-code is generated
  

### Main flow:

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
    
  
Here is a demonstration video:

[![mobile app](https://user-images.githubusercontent.com/98888366/215516946-c8f37970-0c1b-47ca-b356-797ea8149da2.png)
](https://media.bsn.si/ipehr/v2/how_to_add_doctor_into_app.mp4)
