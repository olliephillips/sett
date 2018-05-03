# Sett

## A golang package which offers a simple abstraction on Badger key/value store

**Work in progress**

Based on Badger v0.9.0 API. 

## API 

Creating or opening a store with Sett is identical to Badger

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

Tables are virtual, simply a prefix on the key, but formalised through the Sett API. The aim being to making organisation, reasoning and usage, a little simpler.

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

Uses concurrent goroutines to split the insert payload, default is 50000 keys per goroutine. To change this set the `sett.BatchSize` per goutine before opening connection.

```
sett.BatchSize = 100000
s := sett.Open(opts)
defer s.Close()
```

Batch insert into "client" table

```
s.Batchup("hello", "world")
s.Batchup("hello-again", "world")
s.Batchup("goodbye", "world")

s.Table("Client").SetBatch()
```