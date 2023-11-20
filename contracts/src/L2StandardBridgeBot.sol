pragma solidity 0.8.20;

import { Ownable } from "openzeppelin-contracts/contracts/access/Ownable.sol";

interface IL2StandardBridge {
    function withdrawTo(
        address _l2Token,
        address _to,
        uint256 _amount,
        uint32 _minGasLimit,
        bytes calldata _extraData
    ) external payable;
}

contract L2StandardBridgeBot is Ownable {
    address public constant L2_STANDARD_BRIDGE_ADDRESS = 0x4200000000000000000000000000000000000010;
    IL2StandardBridge public L2_STANDARD_BRIDGE = IL2StandardBridge(payable(L2_STANDARD_BRIDGE_ADDRESS));

    address internal constant LEGACY_ERC20_ETH = 0xDeadDeAddeAddEAddeadDEaDDEAdDeaDDeAD0000;

    uint256 public delegationFee;

    event WithdrawTo(address indexed from, address l2Token, address to, uint256 amount, uint32 minGasLimit, bytes extraData);

    receive() external payable { }

    fallback() payable external { }

    constructor(address payable _owner, uint256 _delegationFee) Ownable(_owner) {
        delegationFee = _delegationFee;
    }

    // withdrawTo withdraws the _amount of _l2Token to _to address.
    function withdrawTo(
        address _l2Token,
        address _to,
        uint256 _amount,
        uint32 _minGasLimit,
        bytes calldata _extraData
    ) public payable {
        if (_l2Token == LEGACY_ERC20_ETH) {
            require(msg.value == delegationFee + _amount, "BNB withdrawal: msg.value does not equal to delegationFee + amount");

            L2_STANDARD_BRIDGE.withdrawTo{value: _amount}(_l2Token, _to, _amount, _minGasLimit, _extraData);
        } else {
            require(msg.value == delegationFee, "BEP20 withdrawal: msg.value does not equal to delegationFee");

            L2_STANDARD_BRIDGE.withdrawTo{value: 0}(_l2Token, _to, _amount, _minGasLimit, _extraData);
        }

        emit WithdrawTo(msg.sender, _l2Token, _to, _amount, _minGasLimit, _extraData);
    }

    function withdraw(
        address _l2Token,
        uint256 _amount,
        uint32 _minGasLimit,
        bytes calldata _extraData
    ) public payable {
        withdrawTo(_l2Token, msg.sender, _amount, _minGasLimit, _extraData);
    }

    // withdrawFee withdraw the delegation fee vault to _recipient address, only owner can call this function.
    function withdrawFee(address _recipient) external onlyOwner {
        payable(_recipient).transfer(address(this).balance);
    }

    // setDelegationFee set the delegation fee, only owner can call this function.
    function setDelegationFee(uint256 _delegationFee) external onlyOwner {
        delegationFee = _delegationFee;
    }
}
