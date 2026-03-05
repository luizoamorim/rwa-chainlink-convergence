// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

contract VehicleNFT is ERC721, Ownable {

    uint256 private _nextTokenId = 1;

    string private baseURI;

    event VehicleTokenized(
        uint256 indexed tokenId,
        string plate,
        string renavam,
        uint256 value
    );

    mapping(string => bool) public vehicleExists;

    constructor(address initialOwner, string memory baseURI_)
        ERC721("AutoLock RWA", "ALRWA")
        Ownable(initialOwner)
    {
        baseURI = baseURI_;
    }

    function mintVehicle(
        address to,
        string memory plate,
        string memory renavam,
        uint256 value
    ) external onlyOwner returns (uint256) {

        require(to != address(0), "Invalid owner");
        require(!vehicleExists[plate], "Vehicle already tokenized");

        uint256 tokenId = _nextTokenId++;

        vehicleExists[plate] = true;

        _safeMint(to, tokenId);

        emit VehicleTokenized(tokenId, plate, renavam, value);

        return tokenId;
    }

    function _baseURI() internal view override returns (string memory) {
        return baseURI;
    }

    /**
     * Soulbound
     */
    function _update(
        address to,
        uint256 tokenId,
        address auth
    ) internal override returns (address) {

        address from = _ownerOf(tokenId);

        // allow mint (from == address(0))
        // block transfers
        if (from != address(0) && to != address(0)) {
            revert("Vehicle title is non-transferable");
        }

        return super._update(to, tokenId, auth);
    }
}