package sett_test

import (
	"os"
	"strconv"
	"testing"

	"github.com/olliephillips/sett"
)

// instance for testing
var s *sett.Sett

func TestMain(m *testing.M) {
	// set up database for tests
	opts := sett.DefaultOptions
	opts.Dir = "data"
	opts.ValueDir = "data/log"

	s = sett.Open(opts)
	defer s.Close()

	// clean up
	os.RemoveAll("./data")

	os.Exit(m.Run())
}

func TestSet(t *testing.T) {
	// should be able to add a key and value
	if err := s.Set("key", "value"); err != nil {
		t.Error("Set operation failed:", err)
	}
}

func TestGet(t *testing.T) {
	// should be able to retrieve the value of a key
	val, err := s.Get("key")
	if err != nil {
		t.Error("Get operation failed:", err)
	}
	// val should be "value"
	if val != "value" {
		t.Error("Get: Expected value \"value\", got ", val)
	}
}

func TestDelete(t *testing.T) {
	// should be able to delete key
	if err := s.Delete("key"); err != nil {
		t.Error("Delete operation failed:", err)
	}

	// key should be gone
	_, err := s.Get("key")
	if err == nil {
		t.Error("Key \"key\" should not be found, but it was")
	}
}

func TestTableSet(t *testing.T) {
	// should be able to add key value using table prefix
	if err := s.Table("table").Set("tablekey", "tablevalue"); err != nil {
		t.Error("TableSet operation failed:", err)
	}
}

func TestTableGet(t *testing.T) {
	// should be able to retrieve the value of a key
	val, err := s.Table("table").Get("tablekey")
	if err != nil {
		t.Error("TableGet operation failed:", err)
	}
	// val should be "value"
	if val != "tablevalue" {
		t.Error("TableGet: Expected value \"tablevalue\", got ", val)
	}
}

func TestTableDelete(t *testing.T) {
	// should be able to delete key from table
	if err := s.Table("table").Delete("tablekey"); err != nil {
		t.Error("Delete operation failed:", err)
	}

	// key should be gone
	_, err := s.Table("table").Get("tablekey")
	if err == nil {
		t.Error("Key \"tablekey\" in table \"table\" should not be found, but it was")
	}
}

func TestSetEmptyBatch(t *testing.T) {
	// should not be able to set empty batch
	if err := s.SetBatch(); err == nil {
		t.Error("SetEmptyBatch: should not be able to set empty batch")
	}
}
func TestSetBatch(t *testing.T) {
	// should be able to create a batch and the add it
	// Add 5 items
	s.Batchup("key1", "val1")
	s.Batchup("key2", "val2")
	s.Batchup("key3", "val3")
	s.Batchup("key4", "val4")
	s.Batchup("key5", "val5")

	if err := s.SetBatch(); err != nil {
		t.Error("SetBatch unexpected error:", err)
	}
	// check keys added
	for i := 1; i <= 5; i++ {
		k := "key" + strconv.Itoa(i)
		v := "val" + strconv.Itoa(i)
		val, err := s.Get(k)
		if err != nil {
			t.Error("SetBatch error retrieving", k)
		}
		if val != v {
			t.Error("SetBatch error expected", v, "got", val)
		}
	}
}

func TestSetTableBatch(t *testing.T) {
	// should be able to create a batch and the add it when using table prefix
	// Add 5 items
	s.Batchup("key1", "val1")
	s.Batchup("key2", "val2")
	s.Batchup("key3", "val3")
	s.Batchup("key4", "val4")
	s.Batchup("key5", "val5")
	s.Batchup("filterkey", "filterval")

	if err := s.Table("batch").SetBatch(); err != nil {
		t.Error("SetBatch unexpected error:", err)
	}
	// check keys added
	for i := 1; i <= 5; i++ {
		k := "key" + strconv.Itoa(i)
		v := "val" + strconv.Itoa(i)
		val, err := s.Get(k)
		if err != nil {
			t.Error("SetBatch error retrieving", k)
		}
		if val != v {
			t.Error("SetBatch error expected", v, "got", val)
		}
	}
}

func TestScanFilter(t *testing.T) {
	scan, _ := s.Scan("key")
	l := len(scan)
	if l != 5 {
		t.Error("ScanAll expected 5 keys, got", l)
	}

}

func TestTableScanAll(t *testing.T) {
	scan, _ := s.Table("batch").Scan()
	l := len(scan)
	if l != 6 {
		t.Error("TableScanAll expected 6 keys, got", l)
	}
}

func TestTableScanFilter(t *testing.T) {
	scan, _ := s.Table("batch").Scan("filter")
	l := len(scan)
	if l != 1 {
		t.Error("TableScanFilter expected 1 key, got", l)
	}
}

func TestDrop(t *testing.T) {
	// should be able to delete "table"
	if err := s.Table("batch").Drop(); err != nil {
		t.Error("Drop, unexpected error", err)
	}

	// check that a key should be gone
	_, err := s.Table("batch").Get("key1")
	if err == nil {
		t.Error("Key \"key1\" in table \"batch\" should not be found as table droppped, but it was")
	}
}
