// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

// Import Bytes.sol and MerkleTrie.sol and RLPReader.sol
import "./Bytes.sol";
import "./MerkleTrie.sol";
import "./SecureMerkleTrie.sol";
import "./RLPReader.sol";
import "./Types.sol";
import "./Hashing.sol";

// Troubleshoot contract will call `verifyInclusionProof` function of MerkleTrie library
contract Troubleshoot {

    // event MyLog(string message);

    function proveWithdrawalTransaction(
        Types.WithdrawalTransaction memory _tx,
        uint256 _l2OutputIndex,
        Types.OutputRootProof calldata _outputRootProof,
        bytes[] memory _withdrawalProof
    ) public {

        bytes32 withdrawalHash = Hashing.hashWithdrawal(_tx);

        // Compute the storage slot of the withdrawal hash in the L2ToL1MessagePasser contract.
        // Refer to the Solidity documentation for more information on how storage layouts are
        // computed for mappings.
        bytes32 storageKey = keccak256(
            abi.encode(
                withdrawalHash,
                uint256(0) // The withdrawals mapping is at the first slot in the layout.
            )
        );

        // Verify that the hash of this withdrawal was stored in the L2toL1MessagePasser contract
        // on L2. If this is true, under the assumption that the SecureMerkleTrie does not have
        // bugs, then we know that this withdrawal was actually triggered on L2 and can therefore
        // be relayed on L1.
        require(
            SecureMerkleTrie.verifyInclusionProof(
                abi.encode(storageKey),
                hex"01",
                _withdrawalProof,
                _outputRootProof.messagePasserStorageRoot
            ),
            "OptimismPortal: invalid withdrawal inclusion proof"
        );
    }
}
