package documentstore

import (
	"github.com/google/btree"
)

type indexedItem struct {
	key  string
	docs []*Document
}

func (a indexedItem) Less(b indexedItem) bool {
	return a.key < b.key
}

type Index struct {
	tree *btree.BTreeG[indexedItem]
}

func NewIndex() *Index {
	return &Index{
		tree: btree.NewG[indexedItem](2, func(a, b indexedItem) bool {
			return a.Less(b)
		}),
	}
}

func (idx *Index) Insert(key string, doc *Document) {
	item := indexedItem{key: key}
	existing, found := idx.tree.Get(item)
	if found {
		existing.docs = append(existing.docs, doc)
		idx.tree.ReplaceOrInsert(existing)
	} else {
		idx.tree.ReplaceOrInsert(indexedItem{key: key, docs: []*Document{doc}})
	}
}

func (idx *Index) Delete(key string) {
	idx.tree.Delete(indexedItem{key: key})
}

func (idx *Index) RangeQuery(min, max *string, desc bool) []*Document {
	var result []*Document

	start := indexedItem{}
	if min != nil {
		start.key = *min
	}

	end := indexedItem{}
	if max != nil {
		end.key = *max
	}

	if desc {
		idx.tree.DescendLessOrEqual(end, func(i indexedItem) bool {
			if min != nil && i.key < *min {
				return false
			}
			result = append(result, i.docs...)
			return true
		})
	} else {
		idx.tree.AscendGreaterOrEqual(start, func(i indexedItem) bool {
			if max != nil && i.key > *max {
				return false
			}
			result = append(result, i.docs...)
			return true
		})
	}

	return result
}
