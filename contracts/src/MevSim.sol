pragma solidity >=0.8.0;

contract MevSim {
    error BlockMismatch();
    error SlotValueMismatch();

    function getSlot(uint256 slot) public view returns (uint256) {
        uint256 value;
        assembly {
            value := sload(slot)
        }
        return value;
    }

    function auction(uint256 slot, uint256 value, uint256 target_block) public payable {
        // check target_block
        if (block.number != target_block) {
            revert BlockMismatch();
        }

        // read current slot value
        uint256 current = getSlot(slot);
        if (current != value) {
            revert SlotValueMismatch();
        }

        uint256 new_value = value + 1;
        assembly {
            sstore(slot, new_value)
        }

        // send all eth to coinbase
        address payable coinbase = payable(block.coinbase);
        coinbase.transfer(address(this).balance);
    }
}