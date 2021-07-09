//SPDX-License-Identifier: Unlicensed
pragma solidity ^0.8.0;

import "./EnumerableStringMap.sol";

contract AgentRegistry {
    using EnumerableStringMap for EnumerableStringMap.MapData;

    event PoolAdded(bytes32 poolId, address by);
    event AgentAdded(bytes32 poolId, bytes32 agentId, string ref, address by);
    event AgentUpdated(bytes32 poolId, bytes32 agentId, string ref, address by);
    event AgentRemoved(bytes32 poolId, bytes32 agentId, address by);
    event PoolOwnershipTransfered(bytes32 poolId, address from, address to);
    event PoolAdminAdded(bytes32 poolId, address addr);
    event PoolAdminRemoved(bytes32 poolId, address addr);
    event AgentAdminAdded(
        bytes32 poolId,
        bytes32 agentId,
        address admin,
        address by
    );
    event AgentAdminRemoved(
        bytes32 poolId,
        bytes32 agentId,
        address admin,
        address by
    );

    // map[poolId] = agent map
    mapping(bytes32 => EnumerableStringMap.MapData) agents;

    // map[poolId] = boolean (true if exists)
    mapping(bytes32 => bool) public poolExistsMap;

    // pool admins are the same as owners except they can't create admins
    // map[poolId] = owner => bool (true if is an owner of poolId)
    mapping(bytes32 => mapping(address => bool)) public poolAdmins;

    // pool owners can create pool admins and everything else for a pool
    mapping(bytes32 => address) public poolOwners;

    // agent admins can modify specific agents
    // poolId -> agentId -> address -> bool
    mapping(bytes32 => mapping(bytes32 => mapping(address => bool)))
        public agentAdmins;

    modifier onlyPoolOwner(bytes32 _poolId) {
        require(
            poolOwners[_poolId] == msg.sender,
            "Caller is not owner of pool"
        );
        _;
    }

    modifier onlyPoolAdmin(bytes32 _poolId) {
        require(
            poolOwners[_poolId] == msg.sender ||
                poolAdmins[_poolId][msg.sender],
            "Only pool owner or pool admin can perform this operation"
        );
        _;
    }

    modifier onlyAgentAdmin(bytes32 _poolId, bytes32 _agentId) {
        require(
            poolOwners[_poolId] == msg.sender ||
                poolAdmins[_poolId][msg.sender] ||
                agentAdmins[_poolId][_agentId][msg.sender],
            "Only pool owner, pool admin, or agent owner can update agent"
        );
        _;
    }

    // add a pool so that agents can be added to it, sets owner to msg.sender
    function addPool(bytes32 _poolId) public {
        require(!poolExistsMap[_poolId], "Pool already exists");
        poolOwners[_poolId] = msg.sender;
        poolExistsMap[_poolId] = true;
        emit PoolAdded(_poolId, msg.sender);
    }

    // transfers a pool to another address so that it can manage agents
    function transferPoolOwnership(bytes32 _poolId, address _to)
        public
        onlyPoolOwner(_poolId)
    {
        require(_to != address(0), "address(0) is not allowed");
        require(poolOwners[_poolId] != _to, "Address is already owner");
        poolOwners[_poolId] = _to;
        emit PoolOwnershipTransfered(_poolId, msg.sender, _to);
    }

    // adds a new pool admin
    function addPoolAdmin(bytes32 _poolId, address admin)
        public
        onlyPoolOwner(_poolId)
    {
        require(admin != address(0), "Address(0) is not allowed");
        require(!poolAdmins[_poolId][admin], "Address is already an admin");
        poolAdmins[_poolId][admin] = true;
        emit PoolAdminAdded(_poolId, admin);
    }

    // removes a pool admin
    function removePoolAdmin(bytes32 _poolId, address admin)
        public
        onlyPoolOwner(_poolId)
    {
        require(poolAdmins[_poolId][admin], "Address is not an admin");
        poolAdmins[_poolId][admin] = false;
        emit PoolAdminRemoved(_poolId, admin);
    }

    // adds an agent admin
    function addAgentAdmin(
        bytes32 _poolId,
        bytes32 _agentId,
        address admin
    ) public onlyPoolAdmin(_poolId) {
        require(admin != address(0), "Address(0) is not allowed");
        require(
            !agentAdmins[_poolId][_agentId][admin],
            "Address is already an agent owner"
        );

        agentAdmins[_poolId][_agentId][admin] = true;
        emit AgentAdminAdded(_poolId, _agentId, admin, msg.sender);
    }

    // removes an agent admin
    function removeAgentAdmin(
        bytes32 _poolId,
        bytes32 _agentId,
        address admin
    ) public onlyPoolAdmin(_poolId) {
        require(admin != address(0), "Address(0) is not allowed");
        require(
            agentAdmins[_poolId][_agentId][admin],
            "Address is not an agent owner"
        );

        agentAdmins[_poolId][_agentId][admin] = false;
        emit AgentAdminRemoved(_poolId, _agentId, admin, msg.sender);
    }

    // adds an agent to a pool
    function addAgent(
        bytes32 _poolId,
        bytes32 _agentId,
        string memory _ref
    ) public onlyAgentAdmin(_poolId, _agentId) {
        require(
            !agents[_poolId].contains(_agentId),
            "Agent already exists on pool"
        );

        agents[_poolId].set(_agentId, _ref);
        emit AgentAdded(_poolId, _agentId, _ref, msg.sender);
    }

    // update an agent
    function updateAgent(
        bytes32 _poolId,
        bytes32 _agentId,
        string memory _ref
    ) public onlyAgentAdmin(_poolId, _agentId) {
        require(
            agents[_poolId].contains(_agentId),
            "Agent must exist to be updated"
        );
        require(
            keccak256(bytes(agents[_poolId].get(_agentId))) !=
                keccak256(bytes(_ref)),
            "New reference must be different than old reference"
        );

        agents[_poolId].set(_agentId, _ref);
        emit AgentUpdated(_poolId, _agentId, _ref, msg.sender);
    }

    // removes an agent from a pool
    function removeAgent(bytes32 _poolId, bytes32 _agentId)
        public
        onlyAgentAdmin(_poolId, _agentId)
    {
        require(
            agents[_poolId].contains(_agentId),
            "Agent does not exist on pool"
        );
        agents[_poolId].remove(_agentId);
        emit AgentRemoved(_poolId, _agentId, msg.sender);
    }

    function agentLength(bytes32 _poolId) public view returns (uint256) {
        return agents[_poolId].length();
    }

    function agentAt(bytes32 _poolId, uint256 index)
        public
        view
        returns (bytes32, string memory)
    {
        return agents[_poolId].entryAt(index);
    }
}
