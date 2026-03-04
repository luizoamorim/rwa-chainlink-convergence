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

    constructor(
        address forwarder,
        address nftAddress
    ) ReceiverTemplate(forwarder) {
        vehicleNFT = VehicleNFT(nftAddress);
    }

    /**
     * Dummy function for ABI binding in CRE.
     */
    function mintVehicle(VehicleReport memory data) public {}

    /**
     * Called by Chainlink forwarder with encoded report.
     */
    function _processReport(bytes calldata report) internal override {
        VehicleReport memory data =
            abi.decode(report, (VehicleReport));

        vehicleNFT.mintVehicle(
            data.owner,
            data.plate,
            data.renavam,
            data.value,
            data.uri
        );
    }
}