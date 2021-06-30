package cmd

import (
	"fmt"
	"os"
)

func CacheClear() (int, error) {
	cacheDir := DefaultCacheDir()
	if cacheDir == "/" || cacheDir == "" {
		panic(fmt.Sprintf("attempted to rm -rf %s", cacheDir))
	}
	err := os.RemoveAll(cacheDir)
	return 0, err
}
