// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import "./VehicleNFT.sol";
import "./ReceiverTemplate.sol";

contract VehicleTokenConsumer is ReceiverTemplate {

    struct VehicleReport {
        address owner;
        string plate;
        string renavam;
        uint256 value;
        string uri;
    }

    VehicleNFT public immutable vehicleNFT;

    event ReportDecoded(
        address owner,
        string plate,
        string renavam,
        uint256 value,
        string uri
    );

    event RawReport(bytes report);

    constructor(
        address forwarder,
        address nftAddress
    ) ReceiverTemplate(forwarder) {
        vehicleNFT = VehicleNFT(nftAddress);
    }

    ////////////////////////////////////////////////////////////
    // CRE ABI BINDING
    ////////////////////////////////////////////////////////////

    function mintVehicle(VehicleReport memory data) public {
        _processVehicleReport(data);
    }

    ////////////////////////////////////////////////////////////
    // CHAINLINK FORWARDER ENTRYPOINT
    ////////////////////////////////////////////////////////////

    function _processReport(bytes calldata report) internal override {

        emit RawReport(report);

        VehicleReport memory data =
            abi.decode(report, (VehicleReport));

        _processVehicleReport(data);
    }

    ////////////////////////////////////////////////////////////
    // INTERNAL PROCESSOR
    ////////////////////////////////////////////////////////////

    function _processVehicleReport(VehicleReport memory data) internal {

        require(data.owner != address(0), "Invalid owner");

        emit ReportDecoded(
            data.owner,
            data.plate,
            data.renavam,
            data.value,
            data.uri
        );

        vehicleNFT.mintVehicle(
            data.owner,
            data.plate,
            data.renavam,
            data.value
        );
    }
}