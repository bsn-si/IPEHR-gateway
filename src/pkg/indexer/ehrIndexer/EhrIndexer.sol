pragma solidity ^0.8.0;

import "hardhat/console.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/access/AccessControl.sol";

contract EhrIndexer is Ownable {
  struct DocumentMeta {
    uint8 docType;
    uint8 status;
    uint256 storageId;
    bytes docIdEncrypted;
    uint32 timestamp;
  }

  mapping (uint256 => DocumentMeta[]) public ehrDocs; // ehr_id -> DocumentMeta[]
  mapping (uint256 => uint256) public ehrUsers; // userId -> EHRid
  mapping (uint256 => uint256) public ehrSubject;  // subjectKey -> ehr_id
  mapping (uint256 => bytes) public docAccess;
  mapping (uint256 => bytes) public dataAccess;
  mapping (address => bool) public allowedChange;

  event EhrSubjectSet(uint256 subjectKey, uint256 ehrId);
  event EhrDocAdded(uint256 ehrId, uint256 storageId);
  event DocAccessChanged(uint256 userId, bytes access);
  event DataAccessChanged(uint256 userId, bytes access);

  modifier onlyAllowed(address _addr) {
    require(allowedChange[_addr] == true, "Not allowed");
    _;
  }

  function setAllowed(address addr, bool allowed) external onlyOwner() returns (bool) {
    allowedChange[addr] = allowed;
    return true;
  }

  function setEhrUser(uint256 userId, uint256 ehrId) external onlyAllowed(msg.sender) returns (uint256) {
    ehrUsers[userId] = ehrId;
    return ehrId;
  }

  function addEhrDoc(uint256 ehrId, DocumentMeta calldata docMeta) external onlyAllowed(msg.sender) {
      ehrDocs[ehrId].push(docMeta);
      emit EhrDocAdded(ehrId, docMeta.storageId);
  }

  function getEhrDocs(uint256 ehrId) public view returns(DocumentMeta[] memory) {
    return ehrDocs[ehrId];
  }

  function setEhrSubject(uint256 subjectKey, uint256 _ehrId) external onlyAllowed(msg.sender) returns (uint256) {
    ehrSubject[subjectKey] = _ehrId;
    emit EhrSubjectSet(subjectKey, _ehrId);
    return _ehrId;
  }

  function setDocAccess(uint256 userId, bytes memory _access) external onlyAllowed(msg.sender) returns (uint256) {
    docAccess[userId] = _access;
    emit DocAccessChanged(userId, _access);
    return userId;
  }

  function setDataAccess(uint256 userId, bytes memory _access) external onlyAllowed(msg.sender) returns (uint256) {
    dataAccess[userId] = _access;
    emit DataAccessChanged(userId, _access);
    return userId;
  }
}

