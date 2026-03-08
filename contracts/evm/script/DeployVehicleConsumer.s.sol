// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "forge-std/Script.sol";
import "../src/evm/VehicleTokenConsumer.sol";

contract DeployVehicleConsumer is Script {

    function run() external {
        address nftAddress = vm.envAddress("VEHICLE_NFT_ADDRESS");
        address forwarder = vm.envAddress("CHAINLINK_FORWARDER");

        vm.startBroadcast();

        new VehicleTokenConsumer(forwarder, nftAddress);

        vm.stopBroadcast();
    }
}