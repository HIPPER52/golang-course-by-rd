package documentstore

type IndexNode struct {
	Key   string
	Docs  []*Document
	Left  *IndexNode
	Right *IndexNode
}

func NewIndexNode() *IndexNode {
	return &IndexNode{
		Key:   "",
		Docs:  make([]*Document, 0),
		Left:  nil,
		Right: nil,
	}
}

type Index struct {
	Root *IndexNode
}

func (idx *Index) Insert(key string, doc *Document) {
	idx.Root = idx.insertRecursive(idx.Root, key, doc)
}

func (idx *Index) Delete(key string) {
	idx.Root = idx.deleteRecursive(idx.Root, key)
}

func (idx *Index) RangeQuery(min, max *string, desc bool) []*Document {
	var result []*Document
	idx.rangeQueryRecursive(idx.Root, min, max, desc, &result)
	return result
}

func (idx *Index) insertRecursive(node *IndexNode, key string, doc *Document) *IndexNode {
	if idx.Root == nil {
		newNode := &IndexNode{
			Key:  key,
			Docs: []*Document{doc},
		}
		idx.Root = newNode
		return newNode
	}

	if node == nil {
		newNode := &IndexNode{
			Key:  key,
			Docs: []*Document{doc},
		}
		return newNode
	}

	if key < node.Key {
		newLeft := idx.insertRecursive(node.Left, key, doc)
		node.Left = newLeft
	} else if key > node.Key {
		newRight := idx.insertRecursive(node.Right, key, doc)
		node.Right = newRight
	} else {
		node.Docs = append(node.Docs, doc)
	}

	return node
}

func (idx *Index) deleteRecursive(node *IndexNode, key string) *IndexNode {
	if node == nil {
		return nil
	}

	if key < node.Key {
		node.Left = idx.deleteRecursive(node.Left, key)
	} else if key > node.Key {
		node.Right = idx.deleteRecursive(node.Right, key)
	} else {
		if node.Left == nil {
			return node.Right
		}
		if node.Right == nil {
			return node.Left
		}

		minLargerNode := idx.findMin(node.Right)
		node.Key = minLargerNode.Key
		node.Docs = minLargerNode.Docs
		node.Right = idx.deleteRecursive(node.Right, minLargerNode.Key)
	}

	return node
}

func (idx *Index) findMin(node *IndexNode) *IndexNode {
	current := node
	for current.Left != nil {
		current = current.Left
	}
	return current
}

func (idx *Index) rangeQueryRecursive(node *IndexNode, min, max *string, desc bool, result *[]*Document) {
	if node == nil {
		return
	}

	if min != nil && node.Key < *min {
		idx.rangeQueryRecursive(node.Right, min, max, desc, result)
		return
	}

	if max != nil && node.Key > *max {
		idx.rangeQueryRecursive(node.Left, min, max, desc, result)
		return
	}

	if desc {
		idx.rangeQueryRecursive(node.Right, min, max, desc, result)
		if idx.keyInRange(node.Key, min, max) {
			*result = append(*result, node.Docs...)
		}
		idx.rangeQueryRecursive(node.Left, min, max, desc, result)
	} else {
		idx.rangeQueryRecursive(node.Left, min, max, desc, result)
		if idx.keyInRange(node.Key, min, max) {
			*result = append(*result, node.Docs...)
		}
		idx.rangeQueryRecursive(node.Right, min, max, desc, result)
	}
}

func (idx *Index) keyInRange(key string, min, max *string) bool {
	if min != nil && key < *min {
		return false
	}
	if max != nil && key > *max {
		return false
	}
	return true
}
