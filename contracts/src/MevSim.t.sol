pragma solidity >=0.8.0;

import "forge-std/Test.sol";
import "./MevSim.sol";

contract MevSimTest is Test {
    MevSim mevSim;

    function setUp() public {
        // deploy MevSim
        mevSim = new MevSim();
    }

    function testSimulate() public {
        uint value = mevSim.getSlot(1);
        mevSim.auction{ value: 1}(1, value, block.number);
        assertEq(mevSim.getSlot(1), value + 1);
    }
}

