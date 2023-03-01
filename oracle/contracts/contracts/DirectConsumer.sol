// SPDX-License-Identifier: MIT
pragma solidity ^0.8.12;

import "@chainlink/contracts/src/v0.8/ChainlinkClient.sol";
import "@chainlink/contracts/src/v0.8/ConfirmedOwner.sol";

// @HOWTO: Please read docs about this with example of chainlink Job config.

/**
 * @title The DirectConsumerContract
 * @notice Simple contract for store usage stats of IPEHR.
 * This contract is intended to be updated externally via direct request to chainlink.
 */
contract DirectConsumerContract is ChainlinkClient, ConfirmedOwner {
    using Chainlink for Chainlink.Request;

    // This stats by MONTH period
    struct PeriodStatistic {
        uint64 documents;
        uint64 patients;
        // Time when data latest update
        uint256 time;
    }

    // Latest total stats
    struct TotalStatistic {
        uint64 documents;
        uint64 patients;
        // Time when data latest update
        uint256 time;
    }

    // cost for pay oracle with link-token
    uint256 private fee;

    // target job id in chainlink
    // @NOTE: accepted only id without dashed
    //        if chainlink return uuid with dashes - remove it from string
    string jobId;
    // Target host for oracle request stats
    string statsHost;

    // Data mapped in int32 is formatted month period with format (YYYYMM)
    // like 202201 where 2022 is year, and 01 is month
    mapping(uint256 => PeriodStatistic) public periods;
    // variables of latest stats
    TotalStatistic public total;

    // @NOTE: this is sample contract, in production please make oracle address & job-id as constant
    constructor(
        address _linkToken,
        address _oracle,
        string memory _jobId,
        string memory _statsHost
    ) ConfirmedOwner(msg.sender) {
        setChainlinkToken(_linkToken);
        setChainlinkOracle(_oracle);

        fee = (1 * LINK_DIVISIBILITY) / 10;
        total = TotalStatistic(0, 0, 0);
        statsHost = _statsHost;
        jobId = _jobId;
    }

    /**
     * @notice Method for request latest total stats from oracle
     */
    function requestLatestStatistics() public onlyOwner {
        Chainlink.Request memory req = buildChainlinkRequest(
            stringToBytes32(jobId),
            address(this),
            this.fulfillLatestStatistics.selector
        );

        // get latest total data
        req.add("url", statsHost);

        sendChainlinkRequest(req, fee);
    }

    /**
     * @notice Method for request stats by period from oracle
     */
    function requestStatisticsByPeriod(string memory period) public onlyOwner {
        Chainlink.Request memory req = buildChainlinkRequest(
            stringToBytes32(jobId),
            address(this),
            this.fulfillStatisticsByPeriod.selector
        );

        // get latest total data
        req.add("url", string.concat(statsHost, "/", period));

        sendChainlinkRequest(req, fee);
    }

    /**
     * @notice Callback for oracle response result
     */
    function fulfillLatestStatistics(
        bytes32 _requestId,
        uint64 documents,
        uint64 patients,
        uint256 time
    ) public recordChainlinkFulfillment(_requestId) {
        total = TotalStatistic(documents, patients, time);
    }

    /**
     * @notice Callback for oracle response result
     */
    function fulfillStatisticsByPeriod(
        bytes32 _requestId,
        uint64 documents,
        uint64 patients,
        uint256 time
    ) public recordChainlinkFulfillment(_requestId) {
        periods[time] = PeriodStatistic(documents, patients, time);
    }

    /**
     * @notice Method for get stats by period, all can call this method
     * @return (u64 documents, u64 partients, uint256 time)
     */
    function getByPeriod(
        uint32 period
    ) public view returns (uint64, uint64, uint256) {
        return (
            periods[period].documents,
            periods[period].patients,
            periods[period].time
        );
    }

    /**
     * @notice Method for get latest total stats, all can call this method
     * @return (u64 documents, u64 partients, u256 time)
     */
    function getTotal() public view returns (uint64, uint64, uint256) {
        return (total.documents, total.patients, total.time);
    }

    /**
     * @notice Method for chainlink token address
     */
    function getChainlinkToken() public view returns (address) {
        return chainlinkTokenAddress();
    }

    /**
     * @notice Method for change chainlink oracle address
     */
    function setOracle(address _oracle) public onlyOwner {
        setChainlinkOracle(_oracle);
    }

    function withdrawLink() public onlyOwner {
        LinkTokenInterface link = LinkTokenInterface(chainlinkTokenAddress());
        require(
            link.transfer(msg.sender, link.balanceOf(address(this))),
            "Unable to transfer"
        );
    }

    function cancelRequest(
        bytes32 _requestId,
        uint256 _payment,
        bytes4 _callbackFunctionId,
        uint256 _expiration
    ) public onlyOwner {
        cancelChainlinkRequest(
            _requestId,
            _payment,
            _callbackFunctionId,
            _expiration
        );
    }

    function stringToBytes32(
        string memory source
    ) private pure returns (bytes32 result) {
        bytes memory tempEmptyStringTest = bytes(source);
        if (tempEmptyStringTest.length == 0) {
            return 0x0;
        }

        assembly {
            // solhint-disable-line no-inline-assembly
            result := mload(add(source, 32))
        }
    }
}
