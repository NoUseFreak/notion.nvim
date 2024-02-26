package cli

import (
	"sync"

	"github.com/nousefreak/notion.nvim/internal/app/cache"
)

var mutex = &sync.Mutex{}
var diskCache *cache.Cache

const (
	CACHE_KEY = "notion.nvim"
)

func loadCache() *cache.Cache {
	mutex.Lock()
	defer mutex.Unlock()
	if diskCache == nil {
		diskCache = cache.New(CACHE_KEY)
	}
	return diskCache
}
