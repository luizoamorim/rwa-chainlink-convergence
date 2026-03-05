// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "forge-std/Script.sol";
import "../src/evm/VehicleNFT.sol";

contract DeployVehicleNFT is Script {

    function run() external {

        vm.startBroadcast();

        new VehicleNFT(
            msg.sender,
            "https://api.autolock.xyz/vehicle/"
        );

        vm.stopBroadcast();
    }
}