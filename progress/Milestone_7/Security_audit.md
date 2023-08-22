<h1 align="center">
    SMART CONTRACTS CODE REVIEW AND SECURITY ANALYSIS REPORT
</h1>

> Customer: Bela Supernova  
> Project: ipEHR  
> Language: Solidity  
> Performed by: BelSoft Dev DOO, Beograd, Srbija  
> 06/06/2023

## Executive Summary

This document deals with the issues of information security of the implementation of smart contracts of the ipEHR project.
In the process of analysis, BelSoft Dev DOO specialists identified potential security vulnerabilities and threats, and provided recommendations for their elimination and additional suggestions and recommendations for improvement.

## Audit Methodology  
### Steps
1. Contract code scanning to detect known vulnerabilities and errors in the code using specialized automated testing and static analysis utilities. For testing we used currently available utilities: Slither, Mythril, Remix IDE analyzer.  
2. Manually audit contract logic for security issues. Contracts are manually analyzed for searches to identify any potential issues.  
3. Additional recommendations. Contract code analysis in terms of best practices, gas consumption optimization, etc.

### Standardization
To standardize the evaluation, we define the following terminology based on OWASP Risk Rating Methodology:

- <b>Likelihood</b> represents how likely a particular vulnerability is to be uncovered and exploited in the wild;  
- <b>Impact</b> measures the technical loss and business damage of a successful attack;  
- <b>Severity</b> demonstrates the overall criticality of the risk.  

Likelihood and impact are categorized into three ratings: high, medium and low respectively. Severity is determined by likelihood and impact and can be classified into four categories accordingly, i.e., Critical, High, Medium, Low shown in the Table below  
<img width="366" alt="SecAud pic 1" src="https://github.com/bsn-si/IPEHR-gateway/assets/98888366/3cbc5dad-ed91-4510-80d3-e3abd4a454b4">

### Project Background
ipEHR is a project that provides the ability to store medical records in Filecoin's decentralized storage and manage access to them via smart contracts.

Project website: https://ipehr.org

Smart contracts in the ipEHR project are used to:
- storing information about users and groups in the system;
- storage of medical records metadata;
- storing medical record data for indexing and processing;
- document access management;

### Contract structure
The project consists of four interconnected smart contracts: EhrIndexer, Users, AccessStore, DataStore

<b>EhrIndexer</b> - is a repository of meta-information about documents, groups of documents and establishes correspondence between users and document identifiers.  
<b>Users</b> - the contract allows to register users and user groups of the system.  
<b>AccessStore</b> - contract-storage of access rights to system objects.  
<b>DataStore</b> - repository of meaningful EHR document data in encrypted form that is used to build a data index and execute AQL queries.  

### Contracts deployments
| Contract              | Network     | Address                                                                                                                            |
|-----------------------|-------------|------------------------------------------------------------------------------------------------------------------------------------|
| EhrIndexer            | Sepolia     | [0x28668f2C40b7FB018F109cC5b2A14C951ff78c91](https://sepolia.etherscan.io/address/0x28668f2C40b7FB018F109cC5b2A14C951ff78c91#code) |
| Users                 | Sepolia     | [0x1B89Ec21E1A0e3E85038daEdf5cA3feF4F087957](https://sepolia.etherscan.io/address/0x1B89Ec21E1A0e3E85038daEdf5cA3feF4F087957#code) |
| AccessStore           | Sepolia     | [0x380ce0529022b4DeEE67A601D1FA1be3B3d2D781](https://sepolia.etherscan.io/address/0x380ce0529022b4DeEE67A601D1FA1be3B3d2D781#code) |
| DataStore             | Sepolia     | [0x906e66Cb9937c1891aC63576B8C082eedB37DF1e](https://sepolia.etherscan.io/address/0x906e66Cb9937c1891aC63576B8C082eedB37DF1e#code) |

## Code Overview
### Libraries
The project uses one library - Attributes. This library includes the structure:

```
struct Attribute {
	Code  code;
	bytes value;
}
```  
which is used to pass attribute sets as function arguments. A set of predefined attribute codes is also specified.  

### Function visibility analysis
#### EhrIndexer
| Function Name                  | Type     | Visibility     |
|--------------------------------|----------|----------------|
| allowedChange                  | read     | public         |
| ehrSubject                     | read     | public         |
| docGroupGetAttrs               | read     | external       |
| docGroupGetDocs                | read     | external       |
| getDocByTime                   | read     | public 🔴      |
| getDocByVersion                | read     | public 🔴      |
| getDocLastByBaseID             | read     | public 🔴      |
| getEhrDocs                     | read     | public 🔴      |
| getEhrUser                     | read     | public 🔴      |
| getLastEhrDocByType            | read     | public 🔴      |
| transferOwnership              | write    | public 🔴      |
| setAllowed                     | write    | external       |
| addEhrDoc                      | write    | external       |
| deleteDoc                      | write    | external       |
| docGroupCreate                 | write    | external       |
| docGroupAddDoc                 | write    | external       |
| setEhrSubject                  | write    | external       |
| setEhrUser                     | write    | external       |
| setAccess                      | write    | external       |
| setEhrDocAttr                  | write    | private        |
| multicall                      | write    | external       |

#### Users
| Function Name                  | Type     | Visibility     |
|--------------------------------|----------|----------------|
| allowedChange                  | read     | public         |
| getUser                        | read     | external       |
| getUserByCode                  | read     | external       |
| userGroupGetByID               | read     | external       |
| transferOwnership              | write    | public 🔴      |
| setAllowed                     | write    | external       |
| userNew                        | write    | external       |
| userGroupCreate                | write    | external       |
| groupAddUser                   | write    | external       |
| groupRemoveUser                | write    | external       |
| setAccess                      | write    | external       |
| multicall                      | write    | external       |

#### AccessStore
| Function Name                  | Type     | Visibility     |
|--------------------------------|----------|----------------|
| allowedChange                  | read     | public 🔴      |
| getAccess                      | read     | external       |
| getAccessByIdHash              | read     | public 🔴      |
| userAccess                     | read     | external       |
| transferOwnership              | write    | public 🔴      |
| setAllowed                     | write    | external       |
| setUsersContractAddress        | write    | external       |
| setAccess                      | write    | external       |

#### DataStore
| Function Name                  | Type     | Visibility     |
|--------------------------------|----------|----------------|
| allowedChange                  | read     | public         |
| transferOwnership              | write    | public 🔴      |
| setAllowed                     | write    | external       |
| dataUpdate                     | write    | external       |

## Code Audit
### Static analysis
#### Restrictable
| Code line             | Issue                                                                       | Severity      |
|-----------------------|-----------------------------------------------------------------------------|---------------|
| Restrictable.sol#43   | Block timestamp comparasion                                                 | low 🟢       |
| Restrictable.sol#21   | Boolean equality                                                            | info 🔵      |
| Restrictable.sol#51   | Boolean equality                                                            | info 🔵      |

#### Common
| Code line             | Issue                                                                       | Severity      |
|-----------------------|-----------------------------------------------------------------------------|---------------|
|                       | solc-0.8.17 is not recommended for deployment. Recommended version: 0.8.18+ | info 🔵      |
| Restrictable.sol#21   | Unused functions from @openzeppelin libraries                               | info 🔵      |

#### EhrIndexer contract
| Code line             | Issue                                                                       | Severity      |
|-----------------------|-----------------------------------------------------------------------------|---------------|
| Docs.sol#11           | Explicitly mark visibility of state                                         | info 🔵      |
| Docs.sol#12           | Explicitly mark visibility of state                                         | info 🔵      |
| Docs.sol#185          | Tautology or contradiction                                                  | medium 🟠    |
| Docs.sol#251          | Unused return                                                               | medium 🟠    |
| Docs.sol#155          | Boolean equality                                                            | info 🔵      |
| Docs.sol#83           | Boolean equality                                                            | info 🔵      |
| DocGroups.sol#15      | Explicitly mark visibility of state                                         | info 🔵      |
| DocGroups.sol#70      | Boolean equality                                                            | info 🔵      |

#### Users contract
| Code line             | Issue                                                                       | Severity      |
|-----------------------|-----------------------------------------------------------------------------|---------------|
| Users.sol#14          | Explicitly mark visibility of state                                         | info 🔵      |
| Users.sol#15          | Explicitly mark visibility of state                                         | info 🔵      |
| Users.sol#16          | Explicitly mark visibility of state                                         | info 🔵      |
| Users.sol#162         | Unused return                                                               | medium 🟠    |

#### AccessStore contract
| Code line             | Issue                                                                       | Severity      |
|-----------------------|-----------------------------------------------------------------------------|---------------|
| AccessStore.sol#14    | Explicitly mark visibility of state                                         | info 🔵      |
| AccessStore.sol#22    | Missing zero address validation                                             | low 🟢       |

#### DataStore contract
| Code line             | Issue                                                                       | Severity      |
|-----------------------|-----------------------------------------------------------------------------|---------------|
| DataStore.sol#34      | Boolean equality                                                            | info 🔵      |

### Manual audit
#### Potential vulnerability of user privacy

Code location: Docs.sol#32, Docs.sol#144, Docs.sol#153, Docs.sol#163, Docs.sol#179, Docs.sol#195

The following public functions:  
`getEhrUser(bytes32 userIDHash), getEhrDocs(bytes32 userIDHash, IDocs.Type docType)`, `getLastEhrDocByType(bytes32 ehrId, IDocs.Type docType)`, `getDocByVersion(bytes32 ehrId, IDocs.Type docType, bytes32 docBaseUIDHash, bytes32 version)`, `getDocByTime(bytes32 ehrID, IDocs.Type docType, uint32 timestamp)`, `getDocLastByBaseID(bytes32 userIDHash, IDocs.Type docType, bytes32 UIDHash)`
could be used to retrieve sensitive information. It is recommended to add permission check before retrieving information.

Severity level: Low 🟢  

#### Potential for a replay attack

Code location: Restrictable.sol#42

Description: The `signCheck(address signer, uint deadline, bytes calldata signature)` function uses the `deadline` argument as a time limit for transaction execution. This measure seems to be aimed at preventing replay attacks. However, transactions may be re-sent during the deadline period, which can lead to hard-to-predict consequences. It is recommended to add an additional hash check from `calldata` to prevent replay attacks.

Severity level: Medium 🟠  

#### Over-checking

Code location: DataStore.sol#33

Description: the function `dataUpdate(bytes32 groupID, bytes32 dataID, bytes32 ehrID, bytes calldata data, address signer, uint deadline, bytes calldata signature)` calculates `dataHash` from the function arguments. Probably this action aimed at preventing data duplication is unnecessary since `dataID` can be used, or if the hash check described above in <b>Potential for a replay attack</b> is implemented, there will be no need for `dataHash`.

Severity level: Info 🔵

#### Unused code
Code location: Docs.sol#51-66

Description: The `setEhrDocAttr` function has private visibility, but is not used anywhere.

Severity level: Info 🔵

#### Garbage
Code location: Docs.sol#128-140

Severity level: Info 🔵

#### Additional recommendations

- EhrIndexer contract: Docs.sol#95 - using `Attributes.get` instead of a `‘for’ loop`
- Users contract: Users.sol#37, Users#55 - using the same error code. If these errors occur, it will be impossible to distinguish between them and may make debugging difficult. It is recommended to set different codes for these situations. 
- Solidity Gas Optimizer uses 100 optimization cycles. It is recommend to increase the value to 1000 - 10000 in order to optimize gas flow rate for function calls
- Change the compiler version to solc-0.8.20
- Change the visibility level from <b>public</b> to <b>external</b> for functions with which are no plans to use internal calls. (They are marked in red in the <b>Function visibility analysis</b> table).

## Audit Result
The audit process reviewed the code for smart contracts used in the ipEHR project.
#### Issues detected
- 0 critical
- 0 high
- 4 medium
- 3 low
- 17 informational

Recommendations on code optimization and improvement given.  
Recommended to resolve the identified problems before launching the project into production use.

## Disclamer

This report is not, nor should be considered, an “endorsement” or “disapproval” of any particular project or team. This report is not, nor should be considered, an indication of the economics or value of any “product” or “asset” created by any team or project that contracts BelSoft Dev DOO to perform a security assessment. This report does not provide any warranty or guarantee regarding the absolute bug-free nature of the code analyzed, nor do they provide any indication of the technologies proprietors, business, business model or legal compliance.  
This report should not be used in any way to make decisions around investment or involvement with any particular project. This report represents an extensive assessing process intending to help the ipEHR project increase the quality of the code while reducing a high level of risk presented by cryptographic tokens and blockchain technology.

The goal of this audit is to help reduce the attack vectors and the high level of variance associated with utilizing new and consistently changing technologies, and in no way claims any guarantee of security or functionality of the technology analyzed.
