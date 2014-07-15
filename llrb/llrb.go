// Copyright 2010 Petar Maymounkov. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// A Left-Leaning Red-Black (LLRB) implementation of 2-3 balanced binary search trees,
// based on the following work:
//
//   http://www.cs.princeton.edu/~rs/talks/LLRB/08Penn.pdf
//   http://www.cs.princeton.edu/~rs/talks/LLRB/LLRB.pdf
//   http://www.cs.princeton.edu/~rs/talks/LLRB/Java/RedBlackBST.java
//
//  2-3 trees (and the run-time equivalent 2-3-4 trees) are the de facto standard BST
//  algoritms found in implementations of Python, Java, and other libraries. The LLRB
//  implementation of 2-3 trees is a recent improvement on the traditional implementation,
//  observed and documented by Robert Sedgewick.
//
package llrb

// Tree is a Left-Leaning Red-Black (LLRB) implementation of 2-3 trees
type LLRB struct {
	root *Node
}

type Node struct {
	Item
	Left, Right *Node // Pointers to left and right child nodes
	count       int
	Black       bool // If set, the color of the link (incoming from the parent) is black
	// In the LLRB, new nodes are always red, hence the zero-value for node
}

type Item interface {
	Less(than Item) bool
}

//
func less(x, y Item) bool {
	if x == pinf || y == ninf {
		return false
	}
	if x == ninf || y == pinf {
		return true
	}
	return x.Less(y)
}

// Inf returns an Item that is "bigger than" any other item, if sign is positive.
// Otherwise  it returns an Item that is "smaller than" any other item.
func Inf(sign int) Item {
	if sign == 0 {
		panic("sign")
	}
	if sign > 0 {
		return pinf
	}
	return ninf
}

var (
	ninf = nInf{}
	pinf = pInf{}
)

type nInf struct{}

func (nInf) Less(Item) bool {
	return true
}

type pInf struct{}

func (pInf) Less(Item) bool {
	return false
}

// New() allocates a new tree
func New() *LLRB {
	return &LLRB{}
}

// SetRoot sets the root node of the tree.
// It is intended to be used by functions that deserialize the tree.
func (t *LLRB) SetRoot(r *Node) {
	t.root = r
}

// Root returns the root node of the tree.
// It is intended to be used by functions that serialize the tree.
func (t *LLRB) Root() *Node {
	return t.root
}

// Len returns the number of nodes in the tree.
func (t *LLRB) Len() int {
	if t.root != nil {
		return t.root.count
	}

	return 0
}

// Has returns true if the tree contains an element whose order is the same as that of key.
func (t *LLRB) Has(key Item) bool {
	return t.Get(key) != nil
}

// Get retrieves an element from the tree whose order is the same as that of key.
func (t *LLRB) Get(key Item) Item {
	h := t.root
	for h != nil {
		switch {
		case less(key, h.Item):
			h = h.Left
		case less(h.Item, key):
			h = h.Right
		default:
			return h.Item
		}
	}
	return nil
}

// Min returns the minimum element in the tree.
func (t *LLRB) Min() Item {
	h := t.root
	if h == nil {
		return nil
	}
	for h.Left != nil {
		h = h.Left
	}
	return h.Item
}

// Max returns the maximum element in the tree.
func (t *LLRB) Max() Item {
	h := t.root
	if h == nil {
		return nil
	}
	for h.Right != nil {
		h = h.Right
	}
	return h.Item
}

func (t *LLRB) ReplaceOrInsertBulk(items ...Item) {
	for _, i := range items {
		t.ReplaceOrInsert(i)
	}
}

func (t *LLRB) InsertNoReplaceBulk(items ...Item) {
	for _, i := range items {
		t.InsertNoReplace(i)
	}
}

// ReplaceOrInsert inserts item into the tree. If an existing
// element has the same order, it is removed from the tree and returned.
func (t *LLRB) ReplaceOrInsert(item Item) Item {
	if item == nil {
		panic("inserting nil item")
	}
	var replaced Item
	t.root, replaced = t.replaceOrInsert(t.root, item)
	t.root.Black = true
	return replaced
}

func (t *LLRB) replaceOrInsert(h *Node, item Item) (*Node, Item) {
	if h == nil {
		return newNode(item), nil
	}

	h = h.walkDownRot23()

	var replaced Item
	if less(item, h.Item) { // BUG
		h.Left, replaced = t.replaceOrInsert(h.Left, item)
		if replaced == nil {
			h.count++
		}
	} else if less(h.Item, item) {
		h.Right, replaced = t.replaceOrInsert(h.Right, item)
		if replaced == nil {
			h.count++
		}
	} else {
		replaced, h.Item = h.Item, item
	}

	h = h.walkUpRot23()

	return h, replaced
}

// InsertNoReplace inserts item into the tree. If an existing
// element has the same order, both elements remain in the tree.
func (t *LLRB) InsertNoReplace(item Item) {
	if item == nil {
		panic("inserting nil item")
	}
	t.root = t.insertNoReplace(t.root, item)
	t.root.Black = true
}

func (t *LLRB) insertNoReplace(h *Node, item Item) *Node {
	if h == nil {
		return newNode(item)
	}

	h = h.walkDownRot23()

	if less(item, h.Item) {
		h.Left = t.insertNoReplace(h.Left, item)
	} else {
		h.Right = t.insertNoReplace(h.Right, item)
	}
	h.count++

	return h.walkUpRot23()
}

// Rotation driver routines for 2-3 algorithm

func (h *Node) walkDownRot23() *Node { return h }

func (h *Node) walkUpRot23() *Node {
	if h.Right.isRed() && !h.Left.isRed() {
		h = h.rotateLeft()
	}

	if h.Left.isRed() && h.Left.Left.isRed() {
		h = h.rotateRight()
	}

	if h.Left.isRed() && h.Right.isRed() {
		h.flip()
	}

	return h
}

// Rotation driver routines for 2-3-4 algorithm

func (h *Node) walkDownRot234() *Node {
	if h.Left.isRed() && h.Right.isRed() {
		h.flip()
	}

	return h
}

func (h *Node) walkUpRot234() *Node {
	if h.Right.isRed() && !h.Left.isRed() {
		h = h.rotateLeft()
	}

	if h.Left.isRed() && h.Left.Left.isRed() {
		h = h.rotateRight()
	}

	return h
}

// DeleteMin deletes the minimum element in the tree and returns the
// deleted item or nil otherwise.
func (t *LLRB) DeleteMin() Item {
	var deleted Item
	t.root, deleted = t.root.deleteMin()
	if t.root != nil {
		t.root.Black = true
	}
	return deleted
}

// deleteMin code for LLRB 2-3 trees
func (h *Node) deleteMin() (*Node, Item) {
	if h == nil {
		return nil, nil
	}
	if h.Left == nil {
		return nil, h.Item
	}

	if !h.Left.isRed() && !h.Left.Left.isRed() {
		h = h.moveRedLeft()
	}

	var deleted Item
	h.Left, deleted = h.Left.deleteMin()
	if deleted != nil {
		h.count--
	}

	return h.fixUp(), deleted
}

// DeleteMax deletes the maximum element in the tree and returns
// the deleted item or nil otherwise
func (t *LLRB) DeleteMax() Item {
	var deleted Item
	t.root, deleted = t.root.deleteMax()
	if t.root != nil {
		t.root.Black = true
	}
	return deleted
}

func (h *Node) deleteMax() (*Node, Item) {
	if h == nil {
		return nil, nil
	}
	if h.Left.isRed() {
		h = h.rotateRight()
	}
	if h.Right == nil {
		return nil, h.Item
	}
	if !h.Right.isRed() && !h.Right.Left.isRed() {
		h = h.moveRedRight()
	}
	var deleted Item
	h.Right, deleted = h.Right.deleteMax()
	if deleted != nil {
		h.count--
	}

	return h.fixUp(), deleted
}

// Delete deletes an item from the tree whose key equals key.
// The deleted item is return, otherwise nil is returned.
func (t *LLRB) Delete(key Item) Item {
	var deleted Item
	t.root, deleted = t.delete(t.root, key)
	if t.root != nil {
		t.root.Black = true
	}
	return deleted
}

func (t *LLRB) delete(h *Node, item Item) (*Node, Item) {
	var deleted Item
	if h == nil {
		return nil, nil
	}
	if less(item, h.Item) {
		if h.Left == nil { // item not present. Nothing to delete
			return h, nil
		}
		if !h.Left.isRed() && !h.Left.Left.isRed() {
			h = h.moveRedLeft()
		}
		h.Left, deleted = t.delete(h.Left, item)
	} else {
		if h.Left.isRed() {
			h = h.rotateRight()
		}
		// If @item equals @h.Item and no right children at @h
		if !less(h.Item, item) && h.Right == nil {
			return nil, h.Item
		}
		// PETAR: Added 'h.Right != nil' below
		// foobaz: Added 'h.Left != nil' below
		//if h.Left != nil && h.Right != nil && !h.Right.isRed() && !h.Right.Left.isRed() {
		if h.Right != nil {
			if !h.Right.isRed() && !h.Right.Left.isRed() {
				h = h.moveRedRight()
			}
			// If @item equals @h.Item, and (from above) 'h.Right != nil'
			if !less(h.Item, item) {
				var subDeleted Item
				h.Right, subDeleted = h.Right.deleteMin()
				if subDeleted == nil {
					panic("logic")
				}
				deleted, h.Item = h.Item, subDeleted
			} else { // Else, @item is bigger than @h.Item
				h.Right, deleted = t.delete(h.Right, item)
			}
		}
	}
	if deleted != nil {
		h.count--
	}

	return h.fixUp(), deleted
}

// Internal node manipulation routines

func newNode(item Item) *Node { return &Node{Item: item, count: 1} }

func (h *Node) isRed() bool {
	if h == nil {
		return false
	}
	return !h.Black
}

func (h *Node) rotateLeft() *Node {
	x := h.Right
	if x.Black {
		panic("rotating a black link")
	}

	x.count, h.count = h.count, h.count-x.count
	if x.Left != nil {
		h.count += x.Left.count
	}

	h.Right = x.Left
	x.Left = h
	x.Black = h.Black
	h.Black = false
	return x
}

func (h *Node) rotateRight() *Node {
	x := h.Left
	if x.Black {
		panic("rotating a black link")
	}

	x.count, h.count = h.count, h.count-x.count
	if x.Right != nil {
		h.count += x.Right.count
	}

	h.Left = x.Right
	x.Right = h
	x.Black = h.Black
	h.Black = false
	return x
}

// REQUIRE: Left and Right children must be present
func (h *Node) flip() {
	h.Black = !h.Black
	if h.Left != nil {
		h.Left.Black = !h.Left.Black
	}
	if h.Right != nil {
		h.Right.Black = !h.Right.Black
	}
}

// REQUIRE: Left and Right children must be present
func (h *Node) moveRedLeft() *Node {
	h.flip()
	if h.Right != nil && h.Right.Left.isRed() {
		h.Right = h.Right.rotateRight()
		h = h.rotateLeft()
		h.flip()
	}
	return h
}

// REQUIRE: Left and Right children must be present
func (h *Node) moveRedRight() *Node {
	h.flip()
	if h.Left != nil && h.Left.Left.isRed() {
		h = h.rotateRight()
		h.flip()
	}
	return h
}

func (h *Node) fixUp() *Node {
	if h.Right.isRed() {
		h = h.rotateLeft()
	}

	if h.Left.isRed() && h.Left.Left.isRed() {
		h = h.rotateRight()
	}

	if h.Left.isRed() && h.Right.isRed() {
		h.flip()
	}

	return h
}
