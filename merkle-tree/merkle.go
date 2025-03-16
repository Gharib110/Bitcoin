package merkletree

/*
receive file => 100G, binary installer => md5 hash code => file is real
md5(file) => md5 hash code

how proof a list of objects is real?

16 objects, md5(obj1 + obj2+ ... obj16) => too slow

too large to compute


obj1, obj2 .... obj16

hash(obj1)=> hash1, hash(obj2)=>hash2, ... hash(obj16)=> hash16
100G => 256 bytes

hash(hash1, hash2)=> hash11 is the parent of hash1 and hash2, hash(hash3, hash4) ...., hash(hash15, hash16), => 8 hashes

hash(hash11,                hash22 )....                    hash(hash77 hash88) => 4 hashes

  hash( hash111,   hash222, )                           hash(hash333, hash444) => 2 hashes

  hash(hash1111, hash2222) => hash11111 => merkle root

  merkle tree

  if we have even numbers of hashes => put them in a pair => do hash256 on the pair,
  how about odd number?
  if odd number, duplicate the last one => even number => put them in pair=> do hash256 on the pair
*/

import (
	"fmt"
	ecc "github.com/Gharib110/Bitcoin/elliptic_curve"
	"math"
	"strings"
)

// MerkleParent combine two hashes and compute their hash result
func MerkleParent(hash1 []byte, hash2 []byte) []byte {
	buf := make([]byte, 0)
	buf = append(buf, hash1...)
	buf = append(buf, hash2...)
	return ecc.Hash256(string(buf))
}

/*
MerkleParentLevel

	if there are even number of hashes, put them into pairs, and compute merkle parent for each pair
	if there is odd number, duplicate the last one, put them into pairs, and compute merkle parent for each pair,
*/
func MerkleParentLevel(hashes [][]byte) [][]byte {

	if len(hashes) <= 1 {
		panic("can't take parent level with no more than 1 item")
	}

	if len(hashes)%2 == 1 {
		//odd number, duplicate last one
		hashes = append(hashes, hashes[len(hashes)-1])
	}

	parentLevel := make([][]byte, 0)
	for i := 0; i < len(hashes); i += 2 {
		parent := MerkleParent(hashes[i], hashes[i+1])
		parentLevel = append(parentLevel, parent)
	}

	return parentLevel
}

func MerkleRoot(hashes [][]byte) []byte {
	curLevel := hashes
	for len(curLevel) > 1 {
		curLevel = MerkleParentLevel(curLevel)
	}

	return curLevel[0]
}

/*
task: => list of nodes => build up merkle tree => get value of the root,

height of the tree, give N nodes as the lowest layer, the height of the tree
is lg(N)+1, N = 27, lg(27)+1  lg(27) is about more than 4 but less than 5
int(lg(27)+1) = 5
*/

func ConstructTree(n int32) [][][]byte {
	maxDepth := math.Ceil(math.Log2(float64(n)) + 1)
	merkleTree := make([][][]byte, int(maxDepth))
	nodesInLayer := int(n)
	for depth := maxDepth; depth > 0; depth-- {
		layer := make([][]byte, 0)
		for i := 0; i < nodesInLayer; i++ {
			layer = append(layer, []byte{})
		}
		merkleTree[int(depth-1)] = layer
		//go up to one layer, the number of nodes reduces by half
		if nodesInLayer%2 == 0 {
			nodesInLayer /= 2
		} else {
			nodesInLayer = (nodesInLayer + 1) / 2
		}
	}

	return merkleTree
}

type MerkleTree struct {
	total        int //number of nodes in bottom
	nodes        [][][]byte
	currentDepth int32
	currentIndex int32
	maxDepth     int32
}

func InitEmptyMerkleTree(total int) *MerkleTree {
	merkleTree := &MerkleTree{
		total:        total,
		currentDepth: 0,
		currentIndex: 0,
		maxDepth:     int32(math.Ceil(math.Log2(float64(total)))),
	}

	merkleTree.nodes = ConstructTree(int32(total))
	return merkleTree
}

func (m *MerkleTree) PopulateTree(flagBits string, hashes [][]byte) {
	for len(m.Root()) == 0 {
		/*
			get the current first bit, if it is 0, then we can get the hash value from
			hashes, if it is 1 and the current node is an internal node, then we need to
			get the value from its children, if it is 1 and the node is leaf, then we can
			get the value from hashes
		*/
		if m.IsLeaf() {
			//for leaf, we always have its value in hashes
			flagBits = flagBits[1:]
			//set the value for current node from hashes
			m.SetCurrentNode(hashes[0])
			hashes = hashes[1:]
			m.Up()
		} else {
			leftHash := m.GetLeftNode()
			/*
				if the left child is empty, means we are visiting the node at first time,
				then we can remove the current first bit, if it is not the first time we
				visit the node, then we can't remove the bit
			*/
			if len(leftHash) == 0 {
				//we are visit the node the first time
				if flagBits[0] == '0' {
					//we have current node's value in hashes
					m.SetCurrentNode(hashes[0])
					hashes = hashes[1:]
					//we don't need to go to its children anymore
					m.Up()
				} else {
					m.Left()
				}
				//we only remove the current bit when we first visit the node
				flagBits = flagBits[1:]
			} else if m.RightExist() {
				rightHash := m.GetRightNode()
				if len(rightHash) == 0 {
					m.Right()
				} else {
					//both have the left and right child
					m.SetCurrentNode(MerkleParent(leftHash, rightHash))
					m.Up()
				}
			} else {
				//duplicate the left child
				m.SetCurrentNode(MerkleParent(leftHash, leftHash))
				m.Up()
			}
		}
	}

	if len(hashes) != 0 {
		panic("hashes not all consumed")
	}

	for _, bit := range flagBits {
		if bit != '0' {
			panic("flag bits not all consumed")
		}
	}
}

func NewMerkleTree(hashes [][]byte) *MerkleTree {
	merkleTree := &MerkleTree{
		total:        len(hashes),
		currentDepth: 0,
		currentIndex: 0,
		maxDepth:     int32(math.Ceil(math.Log2(float64(len(hashes))))),
	}

	merkleTree.nodes = ConstructTree(int32(len(hashes)))
	//set up the value for the lowest layer
	for idx, hash := range hashes {
		merkleTree.nodes[merkleTree.maxDepth][idx] = hash
	}

	//set up nodes in up layer
	for len(merkleTree.Root()) == 0 {
		if merkleTree.IsLeaf() {
			merkleTree.Up()
		} else {
			/*
				in order to compute the hash value of the current node, we need to get the value of
				left child and right child, if the left child is empty, then we get the value of left
				child at first, then we check the right child, if  the value right child is empty,
				then we go to get the value of right child, then we MerkleParent(left, right) to
				set the value of  current node,
				whole process is in-order travel of a binary tree
			*/
			leftHash := merkleTree.GetLeftNode()
			rightHash := merkleTree.GetRightNode()
			if len(leftHash) == 0 {
				merkleTree.Left()
			} else if len(rightHash) == 0 {
				merkleTree.Right()
			} else {
				merkleTree.SetCurrentNode(MerkleParent(leftHash, rightHash))
				merkleTree.Up()
			}
		}
	}

	return merkleTree
}

func (m *MerkleTree) String() string {
	result := make([]string, 0)
	for depth, level := range m.nodes {
		items := make([]string, 0)
		short := "nil"
		for index, h := range level {
			if len(h) != 0 {
				short = fmt.Sprintf("%x...", h[:4])
			}
			if depth == int(m.currentDepth) && index == int(m.currentIndex) {
				//mark the current node
				items = append(items, fmt.Sprintf("*%x*...", h[:3]))
			} else {
				items = append(items, short)
			}
		}

		result = append(result, strings.Join(items, ","))
	}

	return strings.Join(result, "\n")
}

func (m *MerkleTree) Up() {
	//point to current node's parent
	if m.currentDepth > 0 {
		m.currentDepth -= 1
	}
	m.currentIndex /= 2
}

func (m *MerkleTree) Left() {
	//go to the left child of the current node
	m.currentDepth += 1
	m.currentIndex *= 2
}

func (m *MerkleTree) Right() {
	//go to the right child of the current node
	m.currentDepth += 1
	m.currentIndex = m.currentIndex*2 + 1
}

func (m *MerkleTree) Root() []byte {
	return m.nodes[0][0]
}

func (m *MerkleTree) SetCurrentNode(value []byte) {
	m.nodes[m.currentDepth][m.currentIndex] = value
}

func (m *MerkleTree) GetCurrentNode() []byte {
	return m.nodes[m.currentDepth][m.currentIndex]
}

func (m *MerkleTree) GetLeftNode() []byte {
	//get the value of left child
	return m.nodes[m.currentDepth+1][m.currentIndex*2]
}

func (m *MerkleTree) GetRightNode() []byte {
	//get the value of the right child
	return m.nodes[m.currentDepth+1][m.currentIndex*2+1]
}

func (m *MerkleTree) IsLeaf() bool {
	return m.currentDepth == m.maxDepth
}

func (m *MerkleTree) RightExist() bool {
	/*
	   if the number of nodes in the list is not power of 2, then some nodes may not have
	   a right child
	*/
	//bug fix
	return len(m.nodes[m.currentDepth+1]) > int(m.currentIndex)*2+1
}
