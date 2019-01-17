# Sett

## A golang package which offers a simple abstraction on BadgerDB key/value store

Based on BadgerDB v1 API. 
No support for v2 at this time.

## API 

Creating or opening a store with Sett is identical to BadgerDB

```
opts := sett.DefaultOptions
opts.Dir = "data"
opts.ValueDir = "data/log"

s := sett.Open(opts)
defer s.Close()
```

Simple set, get and delete a key. Strings used in preference to byte slices. 

```
s.Set("hello", "world")
s.Get("hello")
s.Delete("hello")
```

### Tables

Tables are virtual, simply a prefix on the key, but formalised through the Sett API. The aim being to make organisation, reasoning and usage, a little simpler.

Add a key/value to "client" table

```
s.Table("client").Set("hello", "world")
```

Get value of key from "client" table

```
s.Table("client").Get("hello")
```

Delete key and value from "client" table

```
s.Table("client").Delete("hello")
```

Drop "client" table including all keys

```
s.Table("client").Drop()
```

### Batch Set

Uses concurrent goroutines to split the insert payload, default is 500 keys per goroutine. To change this set the `sett.BatchSize` variable before opening connection.

```
sett.BatchSize = 50
s := sett.Open(opts)
defer s.Close()
```

Batch insert into "client" table

```
s.Batchup("hello", "world")
s.Batchup("hello-again", "world")
s.Batchup("goodbye", "world")

s.Table("client").SetBatch()
```

### Get entire table, or subset of table

Use `sett.Scan()` to return contents of virtual table or a subset of that table based on a prefix filter.

Retrieving all key/values from the "client" table

```
scan, _ := s.Table("client").Scan()
for k, v := range scan {
	log.Println(k, v)
}
```

Using a prefix filter to get a subset of a key/values from "client" table. In the below example the key prefix filter is "hello"

```
scan, _ := s.Table("client").Scan("hello")
for k, v := range scan {
	log.Println(k, v)
}
```