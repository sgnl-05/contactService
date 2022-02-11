package main

import (
	"flag"
	"fmt"
)

func parseFlags(h *contactHandler) {
	storageType = flag.String("d", "file", "Choose data storage type")
	flag.Parse()

	switch *storageType {
	case "memory":
		fmt.Println("Store in memory")
		h.storage = MemoryStorage{
			contactBook: make(map[string]*Contact),
		}
	case "file":
		fmt.Println("Store in \"storage.json\" file")
		h.storage = FileStorage{}
	case "elastic":
		fmt.Println("Store in Elastic")
	default:
		fmt.Println("Fuck are you writing?")
		return
	}
}
