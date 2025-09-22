package storage

import (
	"fmt"
	"strings"

	"sigolang/config"
)

type StorageFactory struct {
	Prefixes []string
	Create   func(*config.StorageConfig) (IStorage, error)
}

var storageFactories []*StorageFactory = []*StorageFactory{}

func RegisterStorage(factory *StorageFactory) {
	storageFactories = append(storageFactories, factory)
}

func allPrefixes() string {
	prefixes := []string{}
	for _, bunFactory := range storageFactories {
		prefixes = append(prefixes, bunFactory.Prefixes...)
	}
	return strings.Join(prefixes, "|")
}

func NewStorage(c *config.StorageConfig) (storage IStorage, err error) {

	uri := c.StorageUri
	if uri == "" {
		return nil, nil
	}

	found := false

	for _, storageFactory := range storageFactories {
		for _, prefix := range storageFactory.Prefixes {
			if found = strings.HasPrefix(uri, prefix); found {
				storage, err = storageFactory.Create(c)
				if err != nil {
					return nil, err
				}

				break
			}
		}

		if found {
			break
		}
	}

	if !found {
		return nil, fmt.Errorf("invalid storage connection string %s, only (%s) is supported", uri, allPrefixes())
	}

	return storage, nil
}
