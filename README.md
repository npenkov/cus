# Cassandra unique storage

Golang library for storing unique data in Cassandra database, by maintaining checksums of the objects that are stored.
It uses LWT to maintain consistency across all data.

## Interface

[create_test.go](https://github.com/npenkov/cus/blob/main/interface.go)

```go
type UniqueStore interface {
	// Create unique object and store it in cassandra
	// In case of the object data is already persisted with other object
	// the function returns the `duplicateId` of the object along with err = `ErrDuplicateObject`
	Create(id string, data []byte) (duplicateId *string, err error)
	// Get returns the object data
	// In case the object is not found - then err = `ErrNotFound` and the data will be `nil`
	Get(id string) (data []byte, err error)
	// CreateSchema creates the tables if they do not exist in the default cluster keyspace
	// If `force` flag is set to true, then the table will be first droped if it exists
	CreateSchema(force bool) error
	// Delete removes the data from the db
	// In case the object is not found - then err = `ErrNotFound` and the data will be `nil`
	// In any other case the original error will be propagated
	Delete(id string) (err error)
}
```

## Blog post

For more details see the [blog post](https://npenkov.copm/2020-11-11-golang-cassandra-lwt/)

## Usage

Import the module:

```sh
go get github.com/npenkov/cus
```

Initialize the storage with cassandra session

```go
cluster := gocql.NewCluster(hostAndPort)
cluster.Keyspace = "cus"
cluster.Timeout = time.Millisecond * 2000
cluster.DisableInitialHostLookup = true
session, err := cluster.CreateSession()
store = cus.NewCassadnraUniqueStore(session)
```

### Store objects

[create_test.go](https://github.com/npenkov/cus/blob/main/create_test.go)

```go
myKey := "my-store-key"
data := []byte{0x64, 0x65, 0x66}

dupId, err := store.Create(myKey, data)
if err != nil || dupId != nil {
  if err == store.ErrDuplicateObject {
    fmt.Printf("The data stored is duplicate with id : %s", *dupId)
  } else {
    fmt.Printf("Failed to fetch data : %v", err)
  }
}
```

### Fetch objects

[get_test.go](https://github.com/npenkov/cus/blob/main/get_test.go)

```go
fetchedData, err := store.Get("55")
if err != nil || fetchedData == nil {
  if err == store.ErrNotFound {
    fmt.Printf("No record found for id 55")
  } else {
    fmt.Printf("Error fetching data from database : %v", err)
  }
}
```

## License

[MIT](https://github.com/npenkov/cus/blob/main/LICENSE)