package cache

import (
	"fmt"
	"strings"

	"sigolang/config"
)

type Cache struct {
	Impl    any
	Version string
}

var _ ICache = (*Cache)(nil)

type CacheFactory struct {
	Prefixes []string
	Create   func(*config.CacheConfig) (*Cache, error)
}

var cacheFactories []*CacheFactory = []*CacheFactory{}

func allPrefixes() string {
	prefixes := []string{}
	for _, bunFactory := range cacheFactories {
		prefixes = append(prefixes, bunFactory.Prefixes...)
	}
	return strings.Join(prefixes, "|")
}

func NewCache(c *config.CacheConfig) (cache *Cache, err error) {
	dsn := c.CacheUri
	if dsn == "" {
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
