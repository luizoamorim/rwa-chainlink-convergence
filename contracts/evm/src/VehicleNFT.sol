// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "@openzeppelin/contracts/token/ERC721/extensions/ERC721URIStorage.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

contract VehicleNFT is ERC721URIStorage, Ownable {
    uint256 private _nextTokenId;

    event VehicleTokenized(uint256 indexed tokenId, string plate, string renavam, uint256 value);

    constructor() ERC721("AutoLock RWA", "ALRWA") Ownable(msg.sender) {}
 
    /**
     * @dev Mint a new Vehicle NFT. 
     * In production, only the Chainlink Forwarder or the Workflow owner should call this.
     */
    function mintVehicle(
        address to,
        string memory plate,
        string memory renavam,
        uint256 value,
        string memory uri
    ) public onlyOwner returns (uint256) {
        uint256 tokenId = _nextTokenId++;
        _safeMint(to, tokenId);
        _setTokenURI(tokenId, uri);

        emit VehicleTokenized(tokenId, plate, renavam, value);

        return tokenId;
    }
}