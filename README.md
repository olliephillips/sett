# Sett

## A golang package which offers a simple abstraction on Badger key/value store

**Work in progress**

Based on badger v0.9.0 api. 

Note. Tables are virtual - they are just a prefix on the key, but something I have found useful.

Strings used in prefence to byte slices. Again something I find more useful, with fewer conversion operations.

## API 

Creating or opening a store with Sett is identical to badger

```
opts := sett.DefaultOptions
opts.Dir = "data"
opts.ValueDir = "data/log"

s := sett.Open(opts)
defer s.Close()
```

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


### Optional syntax

Targeting a table for subsequent commands

```
client := s.Table("client")
client.Set("hello", "world")
client.Get("hello")
client.Delete("hello")
```