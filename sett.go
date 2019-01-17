package sett

import (
	"errors"
	"log"
	"strings"
	"sync"

	"github.com/dgraph-io/badger"
)

var (
	DefaultOptions         = badger.DefaultOptions
	DefaultIteratorOptions = badger.DefaultIteratorOptions
	BatchSize              = 500
)

type Sett struct {
	batchSize int
	db        *badger.DB
	table     string
	batch     []batchItem
}

type batchItem struct {
	Key string
	Val string
}

// Open is constructor function to create badger instance,
// configure defaults and return struct instance
func Open(opts badger.Options) *Sett {
	s := Sett{}
	// defaults
	s.batchSize = BatchSize

	db, err := badger.Open(opts)
	if err != nil {
		log.Fatal("Open: create or open failed")
	}
	s.db = db
	return &s
}

// Table selects the table, operations are to be performed
// on. Used as a prefix on the keys passed to badger
func (s *Sett) Table(table string) *Sett {
	s.table = strings.ToLower(table)
	return s
}

// Set passes a key & value to badger. Expects string for both
// key and value for convenience, unlike badger itself
func (s *Sett) Set(key string, val string) error {
	var err error
	err = s.db.Update(func(txn *badger.Txn) error {
		err = txn.Set([]byte(s.makeKey(key)), []byte(val))
		return err
	})
	return err
}

// Get returns value of queried key from badger
func (s *Sett) Get(key string) (string, error) {
	var val []byte
	var err error
	err = s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(s.makeKey(key)))
		if err != nil {
			return err
		}
		val, err = item.Value()
		if err != nil {
			return err
		}
		return nil
	})
	return string(val), err
}

// Batchup is a helper for creating batches for use in SetBatch
func (s *Sett) Batchup(key string, val string) {
	item := batchItem{key, val}
	s.batch = append(s.batch, item)
}

// SetBatch creates/updates a batch using mutiple goroutines
func (s *Sett) SetBatch() error {
	var wg sync.WaitGroup
	var err error
	itemCount := len(s.batch)
	if itemCount == 0 {
		return errors.New("no batch ready")
	}
	batchCount := int(itemCount / s.batchSize)
	mod := itemCount % s.batchSize
	if mod != 0 {
		batchCount++
	}

	wg.Add(batchCount)
	for i := 0; i < batchCount; i++ {
		var insert []batchItem
		start := i * s.batchSize
		end := ((i + 1) * s.batchSize)
		if i == (batchCount - 1) {
			end = itemCount
		}
		insert = s.batch[start:end]
		go s.batchSetter(insert, &wg)
	}
	wg.Wait()
	s.batch = nil
	return err
}

func (s *Sett) batchSetter(batch []batchItem, wg *sync.WaitGroup) error {
	// goroutine assigned a slice of batchItem, needs more thought
	// on error handling
	var err error
	defer wg.Done()
	err = s.db.Update(func(txn *badger.Txn) error {
		for _, el := range batch {
			_ = txn.Set([]byte(s.makeKey(el.Key)), []byte(el.Val))
		}
		return err
	})
	return err
}

// Scan returns all key/values from a (virtual) table. An
// optional filter allows the table prefix on the key search
// to be expanded
func (s *Sett) Scan(filter ...string) (map[string]string, error) {
	var result = make(map[string]string)
	var err error
	err = s.db.View(func(txn *badger.Txn) error {
		var fullFilter string
		it := txn.NewIterator(DefaultIteratorOptions)
		defer it.Close()

		fullFilter = s.table
		if len(filter) == 1 {
			fullFilter += filter[0]
		}

		for it.Seek([]byte(fullFilter)); it.ValidForPrefix([]byte(fullFilter)); it.Next() {
			item := it.Item()
			k := string(item.Key())
			k = strings.TrimLeft(k, s.table)
			v, _ := item.Value()
			result[k] = string(v)
		}
		return err
	})
	return result, err
}

// Delete removes a key and its value from badger instance
func (s *Sett) Delete(key string) error {
	var err error
	err = s.db.Update(func(txn *badger.Txn) error {
		err = txn.Delete([]byte(s.makeKey(key)))
		return err
	})
	return err
}

// Drop removes all keys with table prefix from badger,
// the effect is as if a table was deleted
func (s *Sett) Drop() error {
	var err error
	var deleteKey []string
	err = s.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(DefaultIteratorOptions)
		prefix := []byte(s.table)
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			key := string(item.Key())
			deleteKey = append(deleteKey, key)
		}
		it.Close()
		return nil
	})
	err = s.db.Update(func(txn *badger.Txn) error {
		for _, d := range deleteKey {
			err = txn.Delete([]byte(d))
			if err != nil {
				break
			}
		}
		return err
	})
	return err
}

// Close wraps badger Close method for defer
func (s *Sett) Close() error {
	return s.db.Close()
}

/* DEPRECATED */
// Purge wraps badger PurgeOlderVersions method
func (s *Sett) Purge() error {
	// return s.db.PurgeOlderVersions()
	return nil
}

func (s *Sett) makeKey(key string) string {
	// makes the real key to be stored which
	// comprises table name and key set
	return s.table + key
}
