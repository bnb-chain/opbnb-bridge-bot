pragma solidity 0.8.20;

import { Ownable } from "openzeppelin-contracts/contracts/access/Ownable.sol";
import { IERC20 } from "openzeppelin-contracts/contracts/token/ERC20/IERC20.sol";

// See also https://github.com/bnb-chain/opbnb/blob/9505ae88d0ec8f593ee036284c9a13672526a232/packages/contracts-bedrock/contracts/L2/L2StandardBridge.sol#L20
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

            IERC20 l2Token = IERC20(_l2Token);
            bool approveSuccess = l2Token.approve(L2_STANDARD_BRIDGE_ADDRESS, _amount + l2Token.allowance(address(this), L2_STANDARD_BRIDGE_ADDRESS));
            require(approveSuccess, "BEP20 withdrawal: approve failed");
            bool transferSuccess = l2Token.transferFrom(msg.sender, address(this), _amount);
            require(transferSuccess, "BEP20 withdrawal: transferFrom failed");

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

    // withdrawFee withdraw the delegation fee vault to _recipient address on L2, only owner can call this function.
    function withdrawFee(address _recipient) external onlyOwner {
        (bool sent, ) = _recipient.call{ value: address(this).balance }("");
        require(sent, "Failed to send Ether");
    }

    // withdrawFeeToL1 withdraw the delegation fee vault to _recipient address on L1, only owner can call this function.
    function withdrawFeeToL1(address _recipient, uint32 _minGasLimit, bytes calldata _extraData) external onlyOwner {
        uint256 _balance = address(this).balance;
        require(_balance > delegationFee, "fee vault balance is insufficient to pay the required delegation fee");

        uint256 _amount = _balance - delegationFee;
        this.withdrawTo{ value: _balance }(
            LEGACY_ERC20_ETH,
            _recipient,
            _amount,
            _minGasLimit,
            _extraData
        );
    }

    // setDelegationFee set the delegation fee, only owner can call this function.
    function setDelegationFee(uint256 _delegationFee) external onlyOwner {
        require(_delegationFee > 0, "_delegationFee cannot be less than or equal to 0 ether");
        require(_delegationFee <= 1e18, "_delegationFee cannot be more than 1 ether");
        delegationFee = _delegationFee;
    }
}
