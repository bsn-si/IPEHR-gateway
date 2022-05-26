package storage

type Storager interface {
	Add(data []byte) (id *[32]byte, err error)
	AddWithId(id *[32]byte, data []byte) (err error)
	ReplaceWithId(id *[32]byte, data []byte) (err error)
	Get(id *[32]byte) (data []byte, err error)
}
