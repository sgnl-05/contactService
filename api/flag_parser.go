package api

import (
	"flag"
	"fmt"
	"github.com/sgnl-05/contactService/storage"
)

var storageType *string

func ParseFlags(h *ContactHandler) {
	storageType = flag.String("d", "file", "Choose data storage type")
	flag.Parse()

	switch *storageType {
	case "memory":
		fmt.Println("Store in memory")
		h.Storage = storage.MemoryStorage{
			ContactBook: make(map[string]*storage.Contact),
		}
	case "file":
		fmt.Println("Store in local file")
		h.Storage = storage.FileStorage{}
	case "elastic":
		fmt.Println("Store in Elastic")
		h.Storage = storage.NewElasticStorage()
	default:
		fmt.Println("Available -d key values: memory|file|elastic")
		return
	}
}
