// SPDX-License-Identifier: MIT
pragma solidity ^0.8.12;

import "./StatisticsContract.sol";

// @HOWTO: Please read docs about this with example.

/**
 * @title The StatConsumer
 * @notice Simple contract for example how to call data from IpehrStatContract
 */
contract StatisticsConsumerContract is ConfirmedOwner {
    address statisticsAddress;

    // variables of latest stats
    uint64 documents;
    uint64 patients;

    // Time when data latest update
    uint256 time;

    /**
     * @notice Constructor receive address to StatisticsContract oracle address, and request values. 
     */
    constructor(address _statisticsAddress) ConfirmedOwner(msg.sender) {
        statisticsAddress = _statisticsAddress;
    }

    /**
     * @notice Method for get latest total stats, all can call this method
     * @return (u64 documents, u64 partients, u256 time)
     */
    function getTotal() public view returns (uint64, uint64, uint256) {
        return (documents, patients, time);
    }

    /**
     * @notice Change requested contract address
     */
    function setStatsContract(address _statisticsAddress) public onlyOwner {
        statisticsAddress = _statisticsAddress;
    }

    function requestLatestStats() public onlyOwner {
        (uint64 _documents, uint64 _patients, uint256 _time) = IStatisticMethods(statisticsAddress).getTotal();
        documents = _documents;
        patients = _patients;
        time = _time;
    }
}