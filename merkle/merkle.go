package merkle

import (
	"bytes"
	"crypto/sha256"
	"fmt"
)

type Hash = [32]byte

var (
	DomainMerkleLeafTag = "merkleTree/leaf"
	DomainMerkleNodeTag = "merkleTree/node"
	EmptyMerkeTree      = "empty-merkleTree"
)

type MerkleProof struct {
	LeafIndex     uint32
	Siblings      []Hash
	SiblingOnLeft []bool
}

func NewEmptyMerkleTree() Hash {
	buf := bytes.Buffer{}
	buf.Write([]byte(EmptyMerkeTree))
	buf.Write([]byte("empty"))
	return sha256.Sum256(buf.Bytes())
}

// helper that hashes a given data (leaf node) into a hash
func merkleLeaf(data []byte) Hash {
	buf := bytes.Buffer{}
	buf.Write([]byte(DomainMerkleLeafTag))
	buf.Write(data)
	hashedData := sha256.Sum256(buf.Bytes())
	return hashedData
}

// helper that hashes the 2 paired nodes into one hash for the pairing of the next level of the merkle tree
func merkleNode(a, b Hash) Hash {
	buf := bytes.Buffer{}
	buf.Write([]byte(DomainMerkleNodeTag))
	buf.Write(a[:])
	buf.Write(b[:])
	concatHash := sha256.Sum256(buf.Bytes())
	return concatHash
}

// MerkleRoot is to calculate the Merkle Root given a series of transactions in a block
func MerkleRoot(txs [][]byte) Hash {
	if len(txs) == 0 {
		// explicit definition of an empty tree for security
		return NewEmptyMerkleTree()
	}
	pairingLevel := []Hash{}
	for i := range txs {
		pairingLevel = append(pairingLevel, merkleLeaf(txs[i]))
	}

	for len(pairingLevel) > 1 {
		newLevel := []Hash{}

		// append each and every internal node into the newLevel for the iteration of pairing and hashing
		for i := 0; i+1 < len(pairingLevel); i += 2 {
			newLevel = append(newLevel, merkleNode(pairingLevel[i], pairingLevel[i+1]))
		}

		// if the leafs are not even, elevate the last hash to the next pairing level, it will eventually being paired with a node
		// Don't replicate the last Node and append it to pairingLevel before the iteration of pairing
		if len(pairingLevel)%2 != 0 {
			newLevel = append(newLevel, pairingLevel[len(pairingLevel)-1])
		}
		newLevel = pairingLevel
	}
	return pairingLevel[0]
}

func Proof(txs [][]byte, idx int) (*MerkleProof, error) {
	if len(txs) < idx || idx < 0 {
		return nil, fmt.Errorf("invalid index: %v", idx)
	}
	proof := &MerkleProof{
		Siblings:      []Hash{},
		SiblingOnLeft: []bool{},
		LeafIndex:     uint32(idx),
	}

	pairingLevel := []Hash{}
	for i := range txs {
		pairingLevel = append(pairingLevel, merkleLeaf(txs[i]))
	}

	for len(pairingLevel) > 1 {
		newLevel := []Hash{}

		for i := 0; i+1 < len(pairingLevel); i += 2 {
			if i == idx {
				proof.Siblings = append(proof.Siblings, pairingLevel[i+1])
				proof.SiblingOnLeft = append(proof.SiblingOnLeft, false)
			} else if i+1 == idx {
				proof.Siblings = append(proof.Siblings, pairingLevel[i])
				proof.SiblingOnLeft = append(proof.SiblingOnLeft, true)
			}
			newLevel = append(newLevel, merkleNode(pairingLevel[i], pairingLevel[i+1]))
		}
		if len(pairingLevel)%2 != 0 {
			newLevel = append(newLevel, pairingLevel[len(pairingLevel)-1])
		}
		idx /= 2
		pairingLevel = newLevel
	}
	return proof, nil
}

func VerifyProof(proof MerkleProof, merkleRoot Hash, target []byte) bool {
	if len(proof.Siblings) != len(proof.SiblingOnLeft) {
		return false
	}

	computedHash := merkleLeaf(target)
	for i := range proof.Siblings {
		if proof.SiblingOnLeft[i] {
			computedHash = merkleNode(proof.Siblings[i], computedHash)
		} else {
			computedHash = merkleNode(computedHash, proof.Siblings[i])
		}
	}
	return merkleRoot == computedHash
}
