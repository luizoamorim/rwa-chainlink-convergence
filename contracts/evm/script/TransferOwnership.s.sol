// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "forge-std/Script.sol";
import "../src/evm/VehicleNFT.sol";

contract TransferOwnership is Script {

    function run() external {
        address nftAddress = vm.envAddress("VEHICLE_NFT_ADDRESS");
        address consumer = vm.envAddress("CONSUMER_ADDRESS");

        vm.startBroadcast();

        VehicleNFT(nftAddress).transferOwnership(consumer);

        vm.stopBroadcast();
    }
}