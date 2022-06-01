package indexer

import (
	commonerr "errors"
	"log"
	"sync"

	"github.com/vmihailenco/msgpack/v5"
	"golang.org/x/crypto/sha3"

	"hms/gateway/pkg/errors"
	"hms/gateway/pkg/storage"
)

type Index struct {
	sync.RWMutex
	id      *[32]byte
	name    string
	cache   map[string][]byte
	storage storage.Storager
}

var instances = make(map[string]*Index)

func Init(name string) *Index {
	if name == "" {
		log.Fatal("name is empty")
		return nil
	}

	if nil == instances[name] {
		id := sha3.Sum256([]byte(name))

		stor := storage.Init()

		data, err := stor.Get(&id)
		if err != nil && !commonerr.Is(err, errors.IsNotExist) {
			log.Fatal(err)
			return nil
		}

		var cache map[string][]byte
		if commonerr.Is(err, errors.IsNotExist) {
			cache = make(map[string][]byte)
		} else {
			err = msgpack.Unmarshal(data, &cache)
			if err != nil {
				log.Fatal(err)
				return nil
			}
		}

		instances[name] = &Index{
			id:      &id,
			name:    name,
			cache:   cache,
			storage: stor,
		}
	}
	return instances[name]
}

func (i *Index) Add(itemId string, item interface{}) (err error) {
	i.Lock()
	defer func() {
		if err != nil {
			delete(i.cache, itemId)
		}
		i.Unlock()
	}()

	_, ok := i.cache[itemId]
	if ok {
		return errors.AlreadyExist
	}

	data, err := msgpack.Marshal(item)
	if err != nil {
		return err
	}

	i.cache[itemId] = data

	data, err = msgpack.Marshal(i.cache)
	if err != nil {
		return err
	}

	err = i.storage.ReplaceWithId(i.id, data)
	if err != nil {
		return err
	}

	return nil
}

func (i *Index) Replace(itemId string, item interface{}) (err error) {
	i.Lock()
	defer func() {
		if err != nil {
			delete(i.cache, itemId)
		}
		i.Unlock()
	}()

	data, err := msgpack.Marshal(item)
	if err != nil {
		return err
	}

	i.cache[itemId] = data

	data, err = msgpack.Marshal(i.cache)
	if err != nil {
		return err
	}

	err = i.storage.ReplaceWithId(i.id, data)
	if err != nil {
		return err
	}

	return nil
}

func (i *Index) GetById(itemId string, dst interface{}) error {
	i.RLock()
	item, ok := i.cache[itemId]
	i.RUnlock()

	if !ok {
		return errors.IsNotExist
	}

	return msgpack.Unmarshal(item, dst)
}

func (i *Index) Delete(itemId string) error {
	i.Lock()
	defer i.Unlock()

	item, ok := i.cache[itemId]
	if !ok {
		return errors.IsNotExist
	}

	delete(i.cache, itemId)
	data, err := msgpack.Marshal(i.cache)
	if err != nil {
		i.cache[itemId] = item
		return err
	}

	err = i.storage.ReplaceWithId(i.id, data)
	if err != nil {
		i.cache[itemId] = item
		return err
	}

	return nil
}
