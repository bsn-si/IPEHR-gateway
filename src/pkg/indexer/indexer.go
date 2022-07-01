package indexer

import (
	"fmt"
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

func Init(name string) *Index {
	if name == "" {
		log.Fatal("name is empty")
	}

	id := sha3.Sum256([]byte(name))

	stor := storage.Storage()

	data, err := stor.Get(&id)
	if err != nil && !errors.Is(err, errors.ErrIsNotExist) {
		log.Fatal(err)
	}

	var cache map[string][]byte
	if errors.Is(err, errors.ErrIsNotExist) {
		cache = make(map[string][]byte)
	} else {
		err = msgpack.Unmarshal(data, &cache)
		if err != nil {
			log.Fatal(err)
		}
	}

	return &Index{
		id:      &id,
		name:    name,
		cache:   cache,
		storage: stor,
	}
}

func (i *Index) Add(itemID string, item interface{}) (err error) {
	i.Lock()
	defer func() {
		if err != nil {
			delete(i.cache, itemID)
		}
		i.Unlock()
	}()

	if _, ok := i.cache[itemID]; ok {
		return errors.ErrAlreadyExist
	}

	data, err := msgpack.Marshal(item)
	if err != nil {
		return fmt.Errorf("item marshal error: %w", err)
	}

	i.cache[itemID] = data

	data, err = msgpack.Marshal(i.cache)
	if err != nil {
		return fmt.Errorf("cache marshal error: %w", err)
	}

	if err = i.storage.ReplaceWithID(i.id, data); err != nil {
		return fmt.Errorf("storage.ReplaceWithID error: %w", err)
	}

	return nil
}

func (i *Index) Replace(itemID string, item interface{}) (err error) {
	i.Lock()
	defer func() {
		if err != nil {
			delete(i.cache, itemID)
		}
		i.Unlock()
	}()

	data, err := msgpack.Marshal(item)
	if err != nil {
		return fmt.Errorf("item marshal error: %w", err)
	}

	i.cache[itemID] = data

	data, err = msgpack.Marshal(i.cache)
	if err != nil {
		return err
	}

	err = i.storage.ReplaceWithID(i.id, data)
	if err != nil {
		return fmt.Errorf("storage.ReplaceWithID error: %w", err)
	}

	return nil
}

func (i *Index) GetByID(itemID string, dst interface{}) error {
	i.RLock()
	item, ok := i.cache[itemID]
	i.RUnlock()

	if !ok {
		return errors.ErrIsNotExist
	}

	if err := msgpack.Unmarshal(item, dst); err != nil {
		return fmt.Errorf("item unmarshal error: %w", err)
	}

	return nil
}

func (i *Index) Delete(itemID string) error {
	i.Lock()
	defer i.Unlock()

	item, ok := i.cache[itemID]
	if !ok {
		return errors.ErrIsNotExist
	}

	delete(i.cache, itemID)

	data, err := msgpack.Marshal(i.cache)
	if err != nil {
		i.cache[itemID] = item
		return err
	}

	err = i.storage.ReplaceWithID(i.id, data)
	if err != nil {
		i.cache[itemID] = item
		return err
	}

	return nil
}
