pragma solidity ^0.8.0;

contract MockRegistry {
    /// dev: all contracts will use default tags because of this
    string public constant version = "0.0.1";

    uint256 constant AGENT_ID = 0x1;

    string public scannerNodeVersion;

    uint256 private _agentsHash;
    uint256 private _agentCount;

    string private _agentManifest;

    struct ScannerNode {
        bool registered;
        bool disabled;
        uint256 scannerPoolId;
        uint256 chainId;
        string metadata;
    }

    constructor (string memory __scannerNodeVersion, string memory __agentManifest) {
        scannerNodeVersion = __scannerNodeVersion;

        _agentsHash = 0;
        _agentCount = 0;
        _agentManifest = __agentManifest;
    }

    function getAgent(uint256 agentId)
    public view
    returns (bool registered, address owner,uint256 agentVersion, string memory metadata, uint256[] memory chainIds) {
        uint256[] memory chains = new uint256[](1);
        chains[0] = 137;
        return (
            true,
            address(0x0),
            1,
            _agentManifest,
            chains
        );
    }

    function getScanner(address scanner) public view returns (ScannerNode memory) {
        ScannerNode memory scannerNode;
        scannerNode.registered = true;
        scannerNode.disabled = true;
        scannerNode.scannerPoolId = 1;
        scannerNode.chainId = 137;
        return scannerNode;
    }

    function getScanner(uint256 scannerId)
        external
        view
        returns (
            bool registered,
            address owner,
            uint256 chainId,
            string memory metadata
        )
    {
        return (true, address(0x0), 137, "");
    }

    /// dev: for both of scanners and agents - anything is enabled
    function isEnabled(uint256 id) public view returns (bool) {
        return true;
    }

    /// dev: does not exist in production contract ABIs
    function linkTestAgent() public {
        _agentsHash = 1;
        _agentCount = 1;
    }

    /// dev: does not exist in production contract ABIs
    function unlinkTestAgent() public {
        _agentsHash = 2;
        _agentCount = 0;
    }

    function numAgentsFor(uint256 scannerId) public view returns (uint256) {
        return _agentCount;
    }

    /// dev: dispatch scanner hash
    function scannerHash(uint256 scannerId) external view returns (uint256 length, bytes32 manifest) {
        return (_agentCount, bytes32(_agentsHash));
    }

    function agentRefAt(uint256 scannerId, uint256 pos)
        external
        view
        returns (
            bool registered,
            address owner,
            uint256 agentId,
            uint256 agentVersion,
            string memory metadata,
            uint256[] memory chainIds,
            bool enabled,
            uint256 disabledFlags
        )
    {
        (registered, owner, agentVersion, metadata, chainIds) = getAgent(AGENT_ID);
        return (registered, owner, agentId, agentVersion, metadata, chainIds, true, 0);
    }
}
