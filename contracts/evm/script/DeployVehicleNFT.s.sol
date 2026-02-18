// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "forge-std/Script.sol";
import "../src/VehicleNFT.sol";

contract DeployVehicleNFT is Script {
    function run() external {
        uint256 deployerPrivateKey = vm.envUint("PRIVATE_KEY");
        vm.startBroadcast(deployerPrivateKey);

        VehicleNFT nft = new VehicleNFT();
        console.log("Deployed to:", address(nft));

        vm.stopBroadcast();
    }
}