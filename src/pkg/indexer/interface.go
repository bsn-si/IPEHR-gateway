package indexer

type Indexer interface {
	Add(id string, item interface{}) error
	Replace(id string, item interface{}) error
	GetById(id string, dst interface{}) error
	Delete(id string) error
}
