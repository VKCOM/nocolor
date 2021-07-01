package cmd

import (
	"os"

	"github.com/i582/cfmt/cmd/cfmt"
)

func CacheClear() (int, error) {
	cacheDir := DefaultCacheDir()
	if cacheDir == "/" || cacheDir == "" {
		cfmt.Printf("Cache clearing error: attempted to {{rm -rf %s}}::red", cacheDir)
		return 2, nil
	}
	err := os.RemoveAll(cacheDir)
	cfmt.Println("The cache has been {{successfully}}::green deleted.")
	return 0, err
}
