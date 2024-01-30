// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "forge-std/console.sol";
import { Bytes } from "./Bytes.sol";
import { RLPReader } from "./RLPReader.sol";

/**
 * @title MerkleTrie
 * @notice MerkleTrie is a small library for verifying standard Ethereum Merkle-Patricia trie
 *         inclusion proofs. By default, this library assumes a hexary trie. One can change the
 *         trie radix constant to support other trie radixes.
 */
library MerkleTrie {
    /**
     * @notice Struct representing a node in the trie.
     *
     * @custom:field encoded The RLP-encoded node.
     * @custom:field decoded The RLP-decoded node.
     */
    struct TrieNode {
        bytes encoded;
        RLPReader.RLPItem[] decoded;
    }

    /**
     * @notice Determines the number of elements per branch node.
     */
    uint256 internal constant TREE_RADIX = 16;

    /**
     * @notice Branch nodes have TREE_RADIX elements and one value element.
     */
    uint256 internal constant BRANCH_NODE_LENGTH = TREE_RADIX + 1;

    /**
     * @notice Leaf nodes and extension nodes have two elements, a `path` and a `value`.
     */
    uint256 internal constant LEAF_OR_EXTENSION_NODE_LENGTH = 2;

    /**
     * @notice Prefix for even-nibbled extension node paths.
     */
    uint8 internal constant PREFIX_EXTENSION_EVEN = 0;

    /**
     * @notice Prefix for odd-nibbled extension node paths.
     */
    uint8 internal constant PREFIX_EXTENSION_ODD = 1;

    /**
     * @notice Prefix for even-nibbled leaf node paths.
     */
    uint8 internal constant PREFIX_LEAF_EVEN = 2;

    /**
     * @notice Prefix for odd-nibbled leaf node paths.
     */
    uint8 internal constant PREFIX_LEAF_ODD = 3;

    /**
     * @notice Verifies a proof that a given key/value pair is present in the trie.
     *
     * @param _key   Key of the node to search for, as a hex string.
     * @param _value Value of the node to search for, as a hex string.
     * @param _proof Merkle trie inclusion proof for the desired node. Unlike traditional Merkle
     *               trees, this proof is executed top-down and consists of a list of RLP-encoded
     *               nodes that make a path down to the target node.
     * @param _root  Known root of the Merkle trie. Used to verify that the included proof is
     *               correctly constructed.
     *
     * @return Whether or not the proof is valid.
     */
    function verifyInclusionProof(
        bytes memory _key,
        bytes memory _value,
        bytes[] memory _proof,
        bytes32 _root
    ) public returns (bool) {
        return Bytes.equal(_value, get(_key, _proof, _root));
    }

    function toString(uint256 value) public returns(string memory) {
        // 初始分配的内存大小，为10个字节
        bytes memory buffer = new bytes(10);

        uint i = 10;
        //将整数的每一位转换为字符
        do {
            buffer[--i] = bytes1(uint8(48 + uint(value % 10)));
            value /= 10;
        } while (value > 0);

        return string(buffer);
    }


    /**
     * @notice Retrieves the value associated with a given key.
     *
     * @param _key   Key to search for, as hex bytes.
     * @param _proof Merkle trie inclusion proof for the key.
     * @param _root  Known root of the Merkle trie.
     *
     * @return Value of the key if it exists.
     */
    function get(
        bytes memory _key,
        bytes[] memory _proof,
        bytes32 _root
    ) public returns (bytes memory) {

        // _key: 0x5ba6b70f96620d27102d4799d259fe08f690c0a21b80d0a1a903db682417a23f

        //
        // bytes[] memory _proof = new bytes[](8);
        // _proof[0] = bytes(hex"f90211a0feb6c6afca5ea568e0b7183d7292c112e9d28a1c846d889636d5c8822463aa77a0df719a24d8049d04a47e0e07ba493110fed0665d897c0ffdb05a7ce474782026a01bec363c069dc337ce16b4ca5cffa11d20e305fd7dcf6d87547b151c31aabf0ca00583d5a451bfa01f5e51d992e129ee17e65adda9495bda77898187484625aa9ea0bd9fef27cd52c2721e519bcd40c63a818095374e505a3bbd88db29c6c4521982a009a31937313c5b2bc06da6fec4075110f19526a9dbde3629ebe36aa8cd8e0851a0555f89002a34f6570bab438e69c5d8b855f2664e31f5dd59b807d0956577cc81a07e6f59ec99405b7afa436bf62abd098a85b43861fd3a9da3214a5d46b2d12543a03e7e883b0406cb781e5bf2a71cafc02ebbcdb0684ac03439633faf339d1259e2a0fa154fcf0bedd25aacac74e9864e4b96742c3d06f58c170cd320ade69db16c88a07e34154a031da5e50a3886efcb15c6681deeb20832726f083f80b03678b56d33a0d24408a39fb57838fdfe85c68f7eac74daab709ec3c22dfd0a579c1dbd64861aa093c1837ab2907fc51de7e918b81b997d60a6c8a18e3538720ee0d0f369530c26a0fca1af24293865f21d18d5eb823e0ecc007306823864f11d59c1c5c370fb22c2a0de8c894b5067d637df401aa6b963a3fc52b774abbedf6688e25bbf9a29672b7ca021c85d42de67601599e605d9ed1257451d1ce403d2478496d0a99d60db643d5280");
        // _proof[1] = bytes(hex"f90211a02288f14667ef872b449e0f0f5526c2ad3634dcdb7376616d18fd87f6af9293bfa027223b95acda1508dbc42caac8edfb422759eedd13628e4dff0c5bf79d275182a0352a96940c817156a6e33c6f6cb2010ac216296f6e912c335982ba75ffb11898a0c3b7339a174d0a5fd842d3bb946a2ba6d9c1088215d9c3f56449af4dabb59e1ca03b331c72db0a3e0b51167ecb4d0ee76126181e3a54718e7cdefbf89ae0ae0d7ea03f162cadc26daa12d91e9ba459c9e61008ef57314e9cb810f074a878f2b81910a00dc670d0489ec214e2a636274d1e8b2df9db343218503ae2bf7ed868d8ab553fa0bffdd8cf5fe58452a64254f44e28ecd886964898fbab68711749423478d9c94ea0ca722fcaddd052f5a8ddce25c0c4161cd9e5552a0697ceccba540b0520cfcb61a0ffc0036ebf8ed2d9603de8c22d36cdc3fb451dd92a23adefe5593a4d8a12299da0bff43610f0c3d44e73e741d58f34f4d8e25aa3f67e0796241210bf321c2f1117a0341dcdbf2722bd7eacc3af46b6a9d80bb775dace93808a6737a1e135b9b225f3a07b3c3ec4c199c11374ba2dd28227e8c1efdbd8cb065768eff27e3918647eaccfa0e5a903a66eedfe8ea6d28a0f93cc91fc5f00f612a54aeeff100dafe84182a815a0203b138ec5130556a58d06184fdc399cb726954d28fe9fc68131f23d211382f5a057e92c1aec0a5b4e723cbbb3928d1008c9a9b4a26206656fccf70c14136a77de80");
        // _proof[2] = bytes(hex"f90211a0437a65e2661356d938aa5321927d85003770566f25f121ec0dba4ad4acae9a8ca04acf516468076bf35f8c10e6aa904612fe08fcfba5c7b76bd87a5fe38e212c4da08db4fa1d774e1653468078c4c6a26cc6236025d6ab88dd1c02dae7c8bc7b19f7a03fb292d186e2d63e163788b3af313403cd12bc4b839037be4ad15a6fbbee217aa0dc569eab6a0b0a1d1dd7d12a24512bff7d4a6f5dc55a5a2dea31ed1a5415d383a0e4eaab256b73b4abc3402f44b6fa90a0854b733cdd5b884915ae8c07f5c78002a075ccef6818cc6c2870a70a0b41e3904ba2355937e028b02cb54a7ab321492357a0b059014d41d68c6c28b7de8f06b57c015210167d64c8589f9d7801dd4f025a60a07bf03f4912eaba2dd295222c1e01c1ddd7baedad61cf90d4e5987fc7ae1c6beaa0309d5aff3d1194cf259389cbb365e43548291aa9125d0550ab4908d1cc958844a0ff3fa5bba3bd95b08506336101c1c71cb920dd3ee68d8120654dac8f4b05d390a03a7b3de7b7476d9bf665396b718a28f8d746828bbf51a396e43669a4c90ca770a0913bc7b7ef04689a19cebf2c785f4eea68ccf9c97be94bb5bdd1346d4b145147a0cc21ab258bc2d961534438ddbbf185de30f8387acbc40f89857d88a23f90047ea0cf8b6876bd0cfd19ddd39fcfe88470d76757f64ba8496cfc382c92ed65706586a051de24dd2a27ffce89bb1dfcd6aa33469a91b50c4332ff9ef45d3cc917613a5c80");
        // _proof[3] = bytes(hex"f90211a00aaeb3f8caea8f7167815fcdf1470777f7e100c781bc9ed79afefac295bb0d03a04b9f1d4ade00185fc62732b28abfb243afff140ebfa034daa4f139e0015a2195a0f9608d7d4346638e5ad4219ef6ea0769aa9fdb4bc6207831cf0f64854b072c2da0b953f4ee5f8ae17d9da8e3c1efaf9a822d0301dab05cae4b776fbd666fce8158a0bd8460bbbcfa3e02af0909445883176e3b65c85a63c6f91c2994fa2979228726a024f9c1c9952c2b6edec11955b50eaccd1ef44099db309f807f85aa4c4130b31ea056ad08d1e812df59fb54d2b847b1e7fa0872aaa1f3ac8f4e314d6893baece56aa02ca9f7cf8514dcca9401081f94b824f03b4d74dc690dfeed6779ddee4a9ad72ba050dffdadd79fbb8d857ad658864cd9cb19e69c73f8ba9c241fc2b451e2124320a00377a910bb501d1b22b7b03a3fa60aa47facf33515c10b0742fb76bd5a3d7a6ca0abcffd053ac378f2aa4d90e95c64dd7cc450fe6f961f9f55722cc49ef7b47f7ba042a8b3431675d5909bad710d519dd2c5d3287f3af8641122cce53b137087ce6ba0b77542431b232d0aabbbdf27e1c1adab761e144d396d0658ec1a84a712906a73a099ca74decdad4e079c0d1f35c8da403acbe65bf284253dab1be788f1808687eba06d866839d5ce08a75b3c8a2383b7eec6120ab32c7dfbe9b036831fc3101c8250a0fd315d671bf51da89b2db014d0f199be522fbf8e19d0276d0a14a1c66c1df8ac80");
        // _proof[4] = bytes(hex"f901f1a0d6228395f59dea0ead05458d95d54c315c7feddf6416bb41aed62dc165446a6ca051e1ecb5424d948d386c0d291bb674fac4c47349f67e40a77850f54254b612b5a0e3b1d858fc2bb8001052eb3a36c9e2bea2966c795848aff0e5fb3cd779cfb565a0046a3c69a3639aa8bc1b3cf2c052d54424a88d2ea90a6388595b160ef5ef161780a0b2b1a778cbd8a3257ff7114116d892491a3ab148be723e22262b842680a9b29ea026da2284e756fcade15d922c0b8f6282a008e1c842e5030a307b3ea2e9a55726a026c0862e8e9e06ddef9ce37044fe2228a3ba9863faaf0cf1ad0863d41f7d7f13a0bcfbf48375296ad05ade4710ffe7e2b419fe246623369f7a6a4fa6f720bb185fa0d982dbbef7671cf7add6efd8612d47477b77f40d62e11983fe4478e959c0a86da0b3bf440462efb53c7056f5ff3ea134a072363a95a15a6945b83e416e6adee245a00396e89cb82998271f55ef0e892d456a7bbf624c23bf751f1ea2da6a7690be42a06f2588a03a07740fe4d300795096252abeb6b984e41b4f7027df84c1c520ea84a0ce163ce486c3de11ca2f6656f5dcf82e62496cb0d5472e042dede34d0e1aedd8a0fec8a8bfbf3e3418c4730c3c7c047cd73695c4c429f6ec6babc7165b9b27e635a002711d0ce2951376bac13daf1d3bcf359a90f68bc804b6b2b2e9efe09a8157af80");
        // _proof[5] = bytes(hex"f89180808080808080a0326fe6fc4a847e4db1144153b15d0d3811f0317239f2885cc2a383825ccb1d8480a04455f6e0d1f8a2022ea33d5e6a537e384483dfccd19d07d21047e9b264b56e3c808080a06174ffa3ff6accab10abb5d61000db050324c608b10c0344a3eec6fa0107a079a0d494a04c23ba556fc3bdbe234562fe6825d55f3aaaa81e43415b8070a733d0bd8080");
        // _proof[6] = bytes(hex"e482000fa03cb626e2849a157c57ca2e62c3dd139cf803efe281f5b6331b3ba92280dd8c42");
        // _proof[7] = bytes(hex"f84d8080808080de9c332c35a4d03ec6ab9b3ffd06c69652ce8e02ff95537f98b7a0feb29c01808080de9c36620d27102d4799d259fe08f690c0a21b80d0a1a903db682417a23f0180808080808080");

        // _root: 0x00483c4fe5298b89a85a912d6d6899428e51e344785cea8d5d3aaa39295e1e17

        console.log("start");
        console.logBytes32(_root);
        // console.logBytes("root", );


        require(_key.length > 0, "MerkleTrie: empty key");

        TrieNode[] memory proof = _parseProof(_proof);
        bytes memory key = Bytes.toNibbles(_key);
        bytes memory currentNodeID = abi.encodePacked(_root);
        uint256 currentKeyIndex = 0;

        uint256 i = 0;
        // Proof is top-down, so we start at the first element (root).
        for (i = 0; i < proof.length; i++) {
            console.log("loop proof[", toString(i), "] , currentKeyIndex = ", toString(currentKeyIndex));

            TrieNode memory currentNode = proof[i];

            // Key index should never exceed total key length or we'll be out of bounds.
            require(
                currentKeyIndex <= key.length,
                "MerkleTrie: key index exceeds total key length"
            );

            if (currentKeyIndex == 0) {
                // First proof element is always the root node.
                require(
                    Bytes.equal(abi.encodePacked(keccak256(currentNode.encoded)), currentNodeID),
                    "MerkleTrie: invalid root hash"
                );
            } else if (currentNode.encoded.length >= 32) {
                // Nodes 32 bytes or larger are hashed inside branch nodes.
                require(
                    Bytes.equal(abi.encodePacked(keccak256(currentNode.encoded)), currentNodeID),
                    "MerkleTrie: invalid large internal hash"
                );
            } else {
                // Nodes smaller than 32 bytes aren't hashed.
                require(
                    Bytes.equal(currentNode.encoded, currentNodeID),
                    "MerkleTrie: invalid internal node hash"
                );
            }

            if (currentNode.decoded.length == BRANCH_NODE_LENGTH) {
                if (currentKeyIndex == key.length) {
                    console.log("       ==> BRANCH_NODE: currentKeyIndex == key.length --> should return");

                    // Value is the last element of the decoded list (for branch nodes). There's
                    // some ambiguity in the Merkle trie specification because bytes(0) is a
                    // valid value to place into the trie, but for branch nodes bytes(0) can exist
                    // even when the value wasn't explicitly placed there. Geth treats a value of
                    // bytes(0) as "key does not exist" and so we do the same.
                    bytes memory value = RLPReader.readBytes(currentNode.decoded[TREE_RADIX]);
                    require(
                        value.length > 0,
                        "MerkleTrie: value length must be greater than zero (branch)"
                    );

                    // Extra proof elements are not allowed.
                    require(
                        i == proof.length - 1,
                        "MerkleTrie: value node must be last node in proof (branch)"
                    );

                    return value;
                } else {
                    console.log("       ==> BRANCH_NODE: currentKeyIndex != key.length --> next");

                    // We're not at the end of the key yet.
                    // Figure out what the next node ID should be and continue.
                    uint8 branchKey = uint8(key[currentKeyIndex]);
                    RLPReader.RLPItem memory nextNode = currentNode.decoded[branchKey];
                    currentNodeID = _getNodeID(nextNode);
                    currentKeyIndex += 1;

                    if (nextNode.length < 32) {
                        bytes memory value = RLPReader.readBytes( RLPReader.readList(currentNodeID)[1]);
                        require(
                            value.length > 0,
                            "MerkleTrie: value length must be greater than zero (branch)"
                        );

                        // Extra proof elements are not allowed.
                        require(
                            i == proof.length - 1,
                            "MerkleTrie: value node must be last node in proof (branch)"
                        );

                        return value;
                    }
                }
            } else if (currentNode.decoded.length == LEAF_OR_EXTENSION_NODE_LENGTH) {

                bytes memory path = _getNodePath(currentNode);
                uint8 prefix = uint8(path[0]);
                uint8 offset = 2 - (prefix % 2);
                bytes memory pathRemainder = Bytes.slice(path, offset);
                bytes memory keyRemainder = Bytes.slice(key, currentKeyIndex);
                uint256 sharedNibbleLength = _getSharedNibbleLength(pathRemainder, keyRemainder);

                // Whether this is a leaf node or an extension node, the path remainder MUST be a
                // prefix of the key remainder (or be equal to the key remainder) or the proof is
                // considered invalid.
                require(
                    pathRemainder.length == sharedNibbleLength,
                    "MerkleTrie: path remainder must share all nibbles with key"
                );

                if (prefix == PREFIX_LEAF_EVEN || prefix == PREFIX_LEAF_ODD) {
                    console.log("       => PREFIX_LEAF");

                    // Prefix of 2 or 3 means this is a leaf node. For the leaf node to be valid,
                    // the key remainder must be exactly equal to the path remainder. We already
                    // did the necessary byte comparison, so it's more efficient here to check that
                    // the key remainder length equals the shared nibble length, which implies
                    // equality with the path remainder (since we already did the same check with
                    // the path remainder and the shared nibble length).
                    require(
                        keyRemainder.length == sharedNibbleLength,
                        "MerkleTrie: key remainder must be identical to path remainder"
                    );

                    // Our Merkle Trie is designed specifically for the purposes of the Ethereum
                    // state trie. Empty values are not allowed in the state trie, so we can safely
                    // say that if the value is empty, the key should not exist and the proof is
                    // invalid.
                    bytes memory value = RLPReader.readBytes(currentNode.decoded[1]);
                    require(
                        value.length > 0,
                        "MerkleTrie: value length must be greater than zero (leaf)"
                    );

                    // Extra proof elements are not allowed.
                    require(
                        i == proof.length - 1,
                        "MerkleTrie: value node must be last node in proof (leaf)"
                    );

                    return value;
                } else if (prefix == PREFIX_EXTENSION_EVEN || prefix == PREFIX_EXTENSION_ODD) {

                    console.log("       => EXTENSION");

                    // Prefix of 0 or 1 means this is an extension node. We move onto the next node
                    // in the proof and increment the key index by the length of the path remainder
                    // which is equal to the shared nibble length.
                    RLPReader.RLPItem memory nextNode = currentNode.decoded[1];
                    currentNodeID = _getNodeID(nextNode);
                    currentKeyIndex += sharedNibbleLength;

                    if (nextNode.length < 32) {
                        bytes memory value = RLPReader.readBytes( RLPReader.readList(currentNodeID)[1]);
                        console.log("       => PREFIX_LEAF");
                        require(
                            value.length > 0,
                            "MerkleTrie: value length must be greater than zero (branch)"
                        );

                        // Extra proof elements are not allowed.
                        require(
                            i == proof.length - 1,
                            "MerkleTrie: value node must be last node in proof (branch)"
                        );

                        return value;
                    }
                } else {
                    revert("MerkleTrie: received a node with an unknown prefix");
                }
            } else {
                revert("MerkleTrie: received an unparseable node");
            }
        }

        string memory errorMessage = string(abi.encodePacked("MerkleTrie: ran out of proof elements ", toString(currentKeyIndex), "  ", toString(key.length)));
        revert(errorMessage);
    }

    /**
     * @notice Parses an array of proof elements into a new array that contains both the original
     *         encoded element and the RLP-decoded element.
     *
     * @param _proof Array of proof elements to parse.
     *
     * @return Proof parsed into easily accessible structs.
     */
    function _parseProof(bytes[] memory _proof) private pure returns (TrieNode[] memory) {
        uint256 length = _proof.length;
        TrieNode[] memory proof = new TrieNode[](length);
        for (uint256 i = 0; i < length; ) {
            proof[i] = TrieNode({ encoded: _proof[i], decoded: RLPReader.readList(_proof[i]) });
            unchecked {
                ++i;
            }
        }

        return proof;
    }

    /**
     * @notice Picks out the ID for a node. Node ID is referred to as the "hash" within the
     *         specification, but nodes < 32 bytes are not actually hashed.
     *
     * @param _node Node to pull an ID for.
     *
     * @return ID for the node, depending on the size of its contents.
     */
    function _getNodeID(RLPReader.RLPItem memory _node) private pure returns (bytes memory) {
        return _node.length < 32 ? RLPReader.readRawBytes(_node) : RLPReader.readBytes(_node);
    }

    /**
     * @notice Gets the path for a leaf or extension node.
     *
     * @param _node Node to get a path for.
     *
     * @return Node path, converted to an array of nibbles.
     */
    function _getNodePath(TrieNode memory _node) private pure returns (bytes memory) {
        return Bytes.toNibbles(RLPReader.readBytes(_node.decoded[0]));
    }

    /**
     * @notice Utility; determines the number of nibbles shared between two nibble arrays.
     *
     * @param _a First nibble array.
     * @param _b Second nibble array.
     *
     * @return Number of shared nibbles.
     */
    function _getSharedNibbleLength(bytes memory _a, bytes memory _b)
        private
        pure
        returns (uint256)
    {
        uint256 shared;
        uint256 max = (_a.length < _b.length) ? _a.length : _b.length;
        for (; shared < max && _a[shared] == _b[shared]; ) {
            unchecked {
                ++shared;
            }
        }
        return shared;
    }
}

