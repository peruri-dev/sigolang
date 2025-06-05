package cache

import (
	"fmt"
	"strings"

	"sigolang/config"
)

type Cache struct {
	Impl interface{}
}

type CacheFactory struct {
	Prefixes []string
	Create   func(*config.Config) (*Cache, error)
}

var cacheFactories []*CacheFactory = []*CacheFactory{}

func allPrefixes() string {
	prefixes := []string{}
	for _, bunFactory := range cacheFactories {
		prefixes = append(prefixes, bunFactory.Prefixes...)
	}
	return strings.Join(prefixes, "|")
}

func NewCache(c *config.Config) (cache *Cache, err error) {
	dsn := c.Cache.CacheUri
	if dsn == "" {
		fmt.Println("not using cache")
		return
	}

	found := false

	for _, cacheFactory := range cacheFactories {
		for _, prefix := range cacheFactory.Prefixes {
			if found = strings.HasPrefix(dsn, prefix); found {
				cache, err = cacheFactory.Create(c)
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
		return nil, fmt.Errorf("invalid cache connection string %s, only (%s) is supported", dsn, allPrefixes())
	}

	return cache, nil
}
