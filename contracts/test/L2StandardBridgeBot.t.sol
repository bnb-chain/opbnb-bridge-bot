pragma solidity 0.8.20;

import "forge-std/Test.sol";
import { IERC20 } from "openzeppelin-contracts/contracts/token/ERC20/IERC20.sol";
import { Ownable } from "openzeppelin-contracts/contracts/access/Ownable.sol";
import "../src/L2StandardBridgeBot.sol";

contract L2StandardBridgeBotTest is Test {
    L2StandardBridgeBot bot;
    uint256 opbnbMainnetFork;
    address deployer = address(0x1234);
    address usdt = 0xCF712f20c85421d00EAa1B6F6545AaEEb4492B75;
    address user = 0x3977f9B1F4912a783B44aBa813dA388AC73a1428;
    uint withdrawFee = 10000;
    address constant LEGACY_ERC20_ETH = 0xDeadDeAddeAddEAddeadDEaDDEAdDeaDDeAD0000;
    address constant L2_STANDARD_BRIDGE = 0x4200000000000000000000000000000000000010;

    event WithdrawTo(address indexed from, address indexed l2Token, address to, uint256 amount, uint32 minGasLimit, bytes extraData);

    event SentMessageExtension1(address indexed sender , uint256 value);

    function setUp() public {
        opbnbMainnetFork = vm.createFork("https://opbnb-testnet.nodereal.io/v1/e9a36765eb8a40b9bd12e680a1fd2bc5");
        vm.selectFork(opbnbMainnetFork);
        vm.rollFork(opbnbMainnetFork, 7125821);
        bot = new L2StandardBridgeBot(payable(deployer), withdrawFee);
    }

    function test_RevertSetDelegationFeeNotOwner() public {
        vm.expectRevert();
        bot.setDelegationFee(100);
    }

    function test_setDelegationFee() public {
        vm.prank(deployer);
        bot.setDelegationFee(100);
        assertEq(bot.delegationFee(), 100);
    }

    function test_withdrawBNB() public {
        vm.prank(user);
        uint amount = 100;
        vm.expectEmit(true, false, false, true);
        emit WithdrawTo(user, LEGACY_ERC20_ETH, user, amount, 200000, "");
        bot.withdrawTo{value: withdrawFee + amount}(LEGACY_ERC20_ETH, user, amount, 200000, "");
    }

    function test_ERC20() public {
        uint amount = 100;
        uint balanceBefore = IERC20(usdt).balanceOf(user);
        console2.log("balanceBefore", balanceBefore);
        vm.prank(user);
        IERC20(usdt).approve(address(bot), amount);
        vm.prank(user);
        vm.expectEmit(true, false, false, true);
        emit WithdrawTo(user, usdt, user, amount, 200000, "");
        bot.withdrawTo{value: withdrawFee}(usdt, user, amount, 200000, "");
    }

    function test_RevertWithdrawFeeNotOwner() public {
        vm.prank(user);
        vm.expectRevert();
        bot.withdrawFee(user);
    }

    function test_WithdrawFee() public {
        vm.prank(deployer);
        bot.withdrawFee(user);
    }

    function test_RevertWithdrawFeeToL1NotOwner() public {
        vm.prank(user);
        vm.expectRevert();
        bot.withdrawFeeToL1(user, 0, "");
    }

    function test_RevertWithdrawFeeToL1InsufficientBalance() public {
        vm.prank(deployer);
        vm.expectRevert();
        bot.withdrawFeeToL1(user, 0, "");
    }

    function test_WithdrawFeeToL1() public {
        // Ensure the vault has sufficient balance
        uint256 prevBalance = address(bot).balance;

        vm.prank(user);
        uint amount = 0;
        uint round = 10;
        for (uint i = 0; i < round; i++) {
            bot.withdrawTo{value: withdrawFee + amount}(LEGACY_ERC20_ETH, user, amount, 200000, "");
        }

        uint256 postBalance = address(bot).balance;
        assertEq(postBalance, prevBalance + withdrawFee * round);

        // WithdrawFeeToL1
        vm.prank(deployer);
        vm.expectEmit(true, true, true, true);
        emit SentMessageExtension1(L2_STANDARD_BRIDGE, postBalance - withdrawFee);
        bot.withdrawFeeToL1(user, 0, "");
    }
}
