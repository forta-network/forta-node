//SPDX-License-Identifier: Unlicensed
pragma solidity ^0.8.0;

import "./EnumerableSetUpgradeable.sol";

// EnumerableStringSet wraps an Bytes32Set set with a lookup map for the source string data
library EnumerableStringMap {
    using EnumerableSetUpgradeable for EnumerableSetUpgradeable.Bytes32Set;

    struct MapData {
        EnumerableSetUpgradeable.Bytes32Set bytes32Set;
        mapping(bytes32 => string) values;
    }

    /**
     * @dev Add a value to the map. O(1).
     *
     * Returns true if the value was added to the set, that is if it was not
     * already present.
     */
    function set(
        MapData storage self,
        bytes32 key,
        string memory value
    ) internal returns (bool) {
        self.values[key] = value;
        return self.bytes32Set.add(key);
    }

    /**
     * @dev Removes a value from a set. O(1).
     *
     * Returns true if the value was removed from the set, that is if it was
     * present.
     */
    function remove(MapData storage self, bytes32 key) internal returns (bool) {
        delete self.values[key];
        return self.bytes32Set.remove(key);
    }

    /**
     * @dev Returns true if the value is in the map. O(1).
     */
    function contains(MapData storage self, bytes32 key)
        internal
        view
        returns (bool)
    {
        return self.bytes32Set.contains(key);
    }

    /**
     * @dev Returns the number of values in the map. O(1).
     */
    function length(MapData storage self) internal view returns (uint256) {
        return self.bytes32Set.length();
    }

    /**
     * @dev Returns the value stored at position `index` in the map. O(1).
     *
     * Note that there are no guarantees on the ordering of values inside the
     * array, and it may change when more values are added or removed.
     *
     * Requirements:
     *
     * - `index` must be strictly less than {length}.
     */
    function at(MapData storage self, uint256 index)
        internal
        view
        returns (string storage)
    {
        bytes32 key = self.bytes32Set.at(index);
        return self.values[key];
    }

    /**
     * @dev Returns the key and value stored at position `index` in the map. O(1).
     *
     * Note that there are no guarantees on the ordering of values inside the
     * array, and it may change when more values are added or removed.
     *
     * Requirements:
     *
     * - `index` must be strictly less than {length}.
     */
    function entryAt(MapData storage self, uint256 index)
        internal
        view
        returns (bytes32, string storage)
    {
        bytes32 key = self.bytes32Set.at(index);
        return (key, self.values[key]);
    }

    /**
     * @dev Returns the value stored at key `key` in the map. O(1).
     */
    function get(MapData storage self, bytes32 key)
        internal
        view
        returns (string storage)
    {
        return self.values[key];
    }
}
