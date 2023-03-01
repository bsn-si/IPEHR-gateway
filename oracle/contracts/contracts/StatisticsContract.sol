// SPDX-License-Identifier: MIT
pragma solidity ^0.8.12;

import "@chainlink/contracts/src/v0.8/ConfirmedOwner.sol";

// @HOWTO: Please read docs about this with example of chainlink Job config.

interface IStatisticMethods {
    function getTotal() external view returns (
        uint64,
        uint64,
        uint256
    );

    function setTotal(
        uint64 documents,
        uint64 patients,
        uint256 time
    ) external;

    function getByPeriod(
        uint256 period
    ) external view returns (uint64, uint64, uint256);

    function setPeriod(
        uint256 period,
        uint64 documents,
        uint64 patients,
        uint256 time
    ) external;

    function setOracle(
        address _oracle
    ) external;
}

/**
 * @title The StatisticsContract
 * @notice Simple contract for store usage stats of IPEHR.
 * This contract is intended to be updated externally using the chainlink cron service.
 */
contract StatisticsContract is ConfirmedOwner, IStatisticMethods {
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

    // event pushed when update stats by period
    event PeriodStatisticUpdated(uint256 indexed period);

    // event pushed when total stats
    event TotalStatisticUpdated();

    // Data mapped in uint256 is formatted month period with format (YYYYMM)
    // like 202201 where 2022 is year, and 01 is month
    mapping(uint256 => PeriodStatistic) public periods;
    // variables of latest stats
    TotalStatistic public total;
    // chainlink oracle address who can update values
    address public oracle;

    /**
     * @notice Constructor receive chainlink oracle address, who can updates inner values.
     * A contract can only allow one external user for  set updates.
     */
    constructor(address _oracle) ConfirmedOwner(msg.sender) {
        total = TotalStatistic(0, 0, block.timestamp);
        oracle = _oracle;
    }

    /**
     * @notice Method for get latest total stats, all can call this method
     * @return (u64 documents, u64 partients, u256 time)
     */
    function getTotal() override external view returns (uint64, uint64, uint256) {
        return (total.documents, total.patients, total.time);
    }

    /**
     * @notice Set latest total stats with time
     */
    function setTotal(
        uint64 documents,
        uint64 patients,
        uint256 time
    ) override external checkRightsForUpdate {
        total = TotalStatistic(documents, patients, time);
        emit TotalStatisticUpdated();
    }

    /**
     * @notice Method for get stats by period, all can call this method
     * @return (u64 documents, u64 partients, uint256 time)
     */
    function getByPeriod(
        uint256 period
    ) override external view returns (uint64, uint64, uint256) {
        return (
            periods[period].documents,
            periods[period].patients,
            periods[period].time
        );
    }

    /**
     * @notice Set stats for target period
     */
    function setPeriod(
        uint256 period,
        uint64 documents,
        uint64 patients,
        uint256 time
    ) override external checkRightsForUpdate {
        periods[period] = PeriodStatistic(documents, patients, time);
        emit PeriodStatisticUpdated(period);
    }

    /**
     * @notice Set new external updater, who can update stats
     */
    function setOracle(address _oracle) override external onlyOwner {
        oracle = _oracle;
    }

    /**
     * @notice Validate access to update data
     */
    function _validateRightsForUpdate() internal view {
        require(
            msg.sender == oracle || msg.sender == owner(),
            "Only owner and oracle can update stats"
        );
    }

    /**
     * @notice Reverts if called by anyone other than the contract owner or oracle.
     */
    modifier checkRightsForUpdate() {
        _validateRightsForUpdate();
        _;
    }
}
