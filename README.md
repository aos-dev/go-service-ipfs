# go-service-ipfs

[InterPlanetary File System(IPFS)](https://ipfs.io/) support for [go-storage](https://github.com/beyondstorage/go-storage).

## Notes

**This package has been moved to [go-storage](https://github.com/beyondstorage/go-storage/tree/master/services/ipfs).**

```shell
go get go.beyondstorage.io/services/ipfs
```

## Install

```go
go get github.com/beyondstorage/go-service-ipfs
```

## Usage

```go
import (
	"log"

	_ "github.com/beyondstorage/go-service-ipfs"
	"github.com/beyondstorage/go-storage/v4/services"
)

func main() {
	store, err := services.NewStoragerFromString("ipfs:///path/to/workdir?endpoint=<ipfs_http_api_endpoint>&gateway=<ipfs_http_gateway>")
	if err != nil {
		log.Fatal(err)
	}
	
	// Write data from io.Reader into hello.txt
	n, err := store.Write("hello.txt", r, length)
}
```

- See more examples in [go-storage-example](https://github.com/beyondstorage/go-storage-example).
- Read [more docs](https://beyondstorage.io/docs/go-storage/services/ipfs) about go-service-ipfs.
