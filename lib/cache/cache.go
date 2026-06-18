package cache

import (
	"fmt"
	"log/slog"
	"strings"

	"sigolang/config"
)

type CacheFactory struct {
	Prefixes []string
	Create   func(*config.CacheConfig) (ICache, error)
}

var cacheFactories []*CacheFactory = []*CacheFactory{}

func RegisterCache(factory *CacheFactory) {
	cacheFactories = append(cacheFactories, factory)
}

func allPrefixes() string {
	prefixes := []string{}
	for _, bunFactory := range cacheFactories {
		prefixes = append(prefixes, bunFactory.Prefixes...)
	}
	return strings.Join(prefixes, "|")
}

func NewCache(c *config.CacheConfig) (cache ICache, err error) {
	dsn := c.CacheUri
	if dsn == "" {
		slog.Info("not using cache")
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
