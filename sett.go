package sett

import (
	"log"

	"github.com/dgraph-io/badger"
)

var (
	DefaultOptions         = badger.DefaultOptions
	DefaultIteratorOptions = badger.DefaultIteratorOptions
	BatchSize              = 100000
)

type Sett struct {
	BatchSize int
	db        *badger.DB
	table     string
}

// Open is constructor function to create badger instance,
// configure defaults and return struct instance
func Open(opts badger.Options) *Sett {
	s := Sett{}

	// defaults
	s.BatchSize = BatchSize

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
	s.table = table
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

// SetAll creates/updates a map of key values in one go.
// Dependent on size, multiple go routines will be used
func (s *Sett) SetAll(batch map[string]string) error {
	/* NOT IMPLEMENTED. WORK IN PROGRESS */
	var err error
	items := len(batch)
	log.Println(items)
	return err
}

// GetAll returns all values from a (virtual) table. An
// optional filter allows the table prefix on the key search
// to be expanded
func (s *Sett) GetAll(filter ...string) (map[string]string, error) {
	/* NOT IMPLEMENTED. WORK IN PROGRESS */
	var values = make(map[string]string)
	var err error
	return values, err
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
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		prefix := []byte(s.table)
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			key := string(item.Key())
			deleteKey = append(deleteKey, key)
		}
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

func (s *Sett) makeKey(key string) string {
	// makes the real key to be stored which
	// comprises table name and key set
	return s.table + key
}
