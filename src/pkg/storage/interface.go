package storage

type Storager interface {
	Add(data []byte) (id *[32]byte, err error)
	AddWithID(id *[32]byte, data []byte) (err error)
	ReplaceWithID(id *[32]byte, data []byte) (err error)
	Get(id *[32]byte) (data []byte, err error)
	Clean() (err error)
}
