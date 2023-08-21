# Milestone Delivery :mailbox:

* **Application Document:** [ipEHR application](https://github.com/filecoin-project/devgrants/issues/418)
* **Milestone Number:** 7

**Context**

In this milestone we've developed the functional MVP with the Better Studio HMS and deployed it in the in the Calibration test net.

**Deliverables**

| Number                | Deliverable                              | Link                                                                                                                                  | Notes                                                                                                                                                                                                                                                                                                                  |
|-----------------------|------------------------------------------|---------------------------------------------------------------------------------------------------------------------------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| :heavy_check_mark:    | License | [LICENSE](https://github.com/bsn-si/IPEHR-stat/blob/main/LICENSE) | Apache 2.0 license |
| :heavy_check_mark:    | Testing Guide | [the ipEHR docs](https://ipehr.gitbook.io/docs/guides/install) | The "How To" guide for contracts deployment, ipEHR gateway and the stats. |
| 1., 3. :heavy_check_mark: | Stress testing to find performance bottlenecks | See [-](-) | -. | 
| 2. :heavy_check_mark: | Security audit | See [BelSoft Dev report](https://github.com/bsn-si/IPEHR-gateway/blob/develop/progress/Milestone_7/Security_audit.md) | The goal of this audit is to help reduce the attack vectors and the high level of variance associated with utilizing new and consistently changing technologies. | 
| 4. :heavy_check_mark: | Creation of the documentation | See [the ipEHR docs](https://ipehr.gitbook.io/docs/) | The ipEHR project documentation. |
| 5.-6. :heavy_check_mark: | MVP release | See [demonstration video](https://media.bsn.si/ipehr/v2/ipehr_API_3.mp4) | This video shows creation of an EHR record for a user in the FIL blockchain from the Better HMS playground and adding an EHR record by a doctor to the ipEHR blockchain repository. |

## Workflow example

We have integrated the ipEHR gateway with the Better Studio HMS playground to pass EHR documents from the HMS to the FIL network.

## The test case:
### Preconditions:

-   A Patient is registered (swagger POST /user/register)
-   A Doctor is registered (swagger POST /user/register)

### Main flow:

-   An HMS admin creates a patient in the HMS with users's ipEHR `UserID`;
-   The user's EHR is autmatically created by the HMS in thу blockchain;
-   A doctor creates a new health record (blood pressure) in the HMS;
-   The created record is found in the blockchain with added blood pressure values;
-   We use Swagger to check blockchain operations;

Here is the demonstration video:

[![video preview](https://github.com/bsn-si/IPEHR-gateway/assets/98888366/8f56a41a-fa8d-41fc-a659-0b525e2cb29f)](https://media.bsn.si/ipehr/v2/ipehr_API_3.mp4)
