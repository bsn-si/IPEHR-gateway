pragma solidity ^0.8.0;

import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/utils/Multicall.sol";

contract EhrIndexer is Ownable, Multicall {
  /**
    Error codes:
    ADL - already deleted
    WTP - wrong type passed
    LST - new version of the EHR document must be the latest
    NFD - not found
  */

  enum DocType { Ehr, EhrAccess, EhrStatus , Composition }
  enum DocStatus { Active, Deleted }
  enum Role { Patient, Doctor }

  struct User {
    bytes32   id;
    bytes32   systemID;
    Role      role;
    bytes32[] groups;
    bytes     pwdHash;
    bool      isUser;
  }

  struct DocumentMeta {
    DocType docType;
    DocStatus status;
    bytes   CID;
    bytes   dealCID;
    bytes   minerAddress;
    bytes   docUIDEncrypted;
    bytes32 docBaseUIDHash;
    bytes32 version;
    bool    isLast;
    uint32  timestamp;
  }

  struct DataEntry {
    uint128 groupID;
    mapping (string => bytes) valueSet;
    bytes docStorIDEncr;
  }

  struct Element {
    bytes32 itemType;
    bytes32 elementType;
    bytes32 nodeID;
    bytes32 name;
    DataEntry[] dataEntries;
  }

  struct Node {
    bytes32 nodeType;
    bytes32 nodeID;
    mapping (bytes32 => Node) next;
    mapping (bytes32 => Element) items;
  }

  Node public dataSearch;
  mapping (bytes32  => mapping(DocType => DocumentMeta[])) public ehrDocs; // ehr_id -> docType -> DocumentMeta[]
  mapping (bytes32  => bytes32) public ehrUsers; // userId -> EHRid
  mapping (bytes32  => bytes32) public ehrSubject;  // subjectKey -> ehr_id
  mapping (bytes32  => bytes) public docAccess;
  mapping (bytes32  => bytes) public groupAccess;
  mapping (address => bool) public allowedChange;
  mapping (address => User) users;

  event EhrSubjectSet(bytes32 subjectKey, bytes32  ehrId);
  event EhrDocAdded(bytes32 ehrId, bytes CID);
  event DocAccessChanged(bytes32 key, bytes access);
  event GroupAccessChanged(bytes32 key, bytes access);

  modifier onlyAllowed(address _addr) {
    require(allowedChange[_addr] == true, "Not allowed");
    _;
  }

  function setAllowed(address addr, bool allowed) external onlyOwner() {
    allowedChange[addr] = allowed;
  }

  function setEhrUser(bytes32 userId, bytes32 ehrId) external onlyAllowed(msg.sender) {
    ehrUsers[userId] = ehrId;
  }

  function addEhrDoc(bytes32 ehrId, DocumentMeta calldata docMeta) external onlyAllowed(msg.sender) {
      require(docMeta.isLast == true, "LST");

      uint i;
      if (docMeta.docType == DocType.Ehr || docMeta.docType == DocType.EhrStatus) {
        for (i = 0; i < ehrDocs[ehrId][docMeta.docType].length; i++) {
            ehrDocs[ehrId][docMeta.docType][i].isLast = false;
        }
      }

      if (docMeta.docType == DocType.Composition) {
        for (i = 0; i < ehrDocs[ehrId][DocType.Composition].length; i++) {
          if (ehrDocs[ehrId][DocType.Composition][i].docBaseUIDHash == docMeta.docBaseUIDHash) {
            ehrDocs[ehrId][DocType.Composition][i].isLast = false;
          }
        }
      }

      ehrDocs[ehrId][docMeta.docType].push(docMeta);
      emit EhrDocAdded(ehrId, docMeta.CID);
  }

  function getEhrDocs(bytes32 ehrId, DocType docType) public view returns(DocumentMeta[] memory) {
    return ehrDocs[ehrId][docType];
  }

  function setEhrSubject(bytes32 subjectKey, bytes32 ehrId) external onlyAllowed(msg.sender) {
    ehrSubject[subjectKey] = ehrId;
    emit EhrSubjectSet(subjectKey, ehrId);
  }

  function setDocAccess(bytes32 key, bytes calldata access) external onlyAllowed(msg.sender) {
    docAccess[key] = access;
    emit DocAccessChanged(key, access);
  }

  function setGroupAccess(bytes32 key, bytes calldata access) external onlyAllowed(msg.sender) {
    groupAccess[key] = access;
    emit GroupAccessChanged(key, access);
  }

  function getLastEhrDocByType(bytes32 ehrId, DocType docType) public view returns(DocumentMeta memory) {
    for (uint i = 0; i < ehrDocs[ehrId][docType].length; i++) {
      if (ehrDocs[ehrId][docType][i].isLast == true) {
        return ehrDocs[ehrId][docType][i];
      }
    }

    revert("NFD");
  }

  function deleteDoc(bytes32 ehrId, DocType docType, bytes32 docBaseUIDHash, bytes32 version) external onlyAllowed(msg.sender) {
    require(docType == DocType.Composition, "WTP");
    for (uint i = 0; i < ehrDocs[ehrId][docType].length; i++) {
      if (ehrDocs[ehrId][docType][i].docBaseUIDHash == docBaseUIDHash && ehrDocs[ehrId][docType][i].version == version) {
        require (ehrDocs[ehrId][docType][i].status != DocStatus.Deleted, "ADL");
        ehrDocs[ehrId][docType][i].status = DocStatus.Deleted;
        return;
      }
    }

    revert("NFD");
  }

  function getDocByVersion(bytes32 ehrId, DocType docType, bytes32 docBaseUIDHash, bytes32 version) public view returns (DocumentMeta memory) {
    for (uint i = 0; i < ehrDocs[ehrId][docType].length; i++) {
      if (ehrDocs[ehrId][docType][i].docBaseUIDHash == docBaseUIDHash && ehrDocs[ehrId][docType][i].version == version) {
        return ehrDocs[ehrId][docType][i];
      }
    }

    revert("NFD");
  }

  function getDocLastByBaseID(bytes32 ehrId, DocType docType, bytes32 docBaseUIDHash) public view returns (DocumentMeta memory) {
    for (uint i = 0; i < ehrDocs[ehrId][docType].length; i++) {
      if (ehrDocs[ehrId][docType][i].docBaseUIDHash == docBaseUIDHash) {
        return ehrDocs[ehrId][docType][i];
      }
    }

    revert("NFD");
  }

  function getDocByTime(bytes32 ehrId, DocType docType, uint32 timestamp) public view returns (DocumentMeta memory) {
    DocumentMeta memory docMeta;
    for (uint i = 0; i < ehrDocs[ehrId][docType].length; i++) {
      if (ehrDocs[ehrId][docType][i].timestamp <= timestamp) {
        docMeta = ehrDocs[ehrId][docType][i];
      } else {
        break;
      }
    }

    require(docMeta.timestamp != 0, "NFD");

    return docMeta;
  }

  function userAdd(address userAddr, bytes32 id, Role role, bytes calldata pwdHash) external onlyAllowed(msg.sender) {
    users[userAddr].id = id;
    users[userAddr].pwdHash = pwdHash;
    users[userAddr].role = role;
    users[userAddr].isUser = true;
  }

  function getUserPasswordHash(address userAddr) public view returns (bytes memory) {
    if (!users[userAddr].isUser) revert("NFD");
    return users[userAddr].pwdHash;
  }
}
