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

    event WithdrawTo(address indexed from, address l2Token, address to, uint256 amount, uint32 minGasLimit, bytes extraData);

    function setUp() public {
        opbnbMainnetFork = vm.createFork("https://opbnb-testnet-rpc.bnbchain.org");
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
}
