package treestorage

import "hms/gateway/pkg/docs/model"

type TreeStorage struct {
	root node
}

func NewTreeStorage() *TreeStorage {
	return &TreeStorage{
		root: node{},
	}
}

type node struct {
	value    interface{}
	children []*node
}

func (ts *TreeStorage) StoreCompositionObject(comp *model.Composition) error {
	return nil
}

func (ts *TreeStorage) FindData() (interface{}, error) {
	return nil, nil
}
