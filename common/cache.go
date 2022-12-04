package common

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"os"
)

func cacheSafeName(fn string) string {
	hash := md5.Sum([]byte(fn))
	return base64.URLEncoding.EncodeToString(hash[:])
}

const cacheDir = "_skolcache"

// CachedASTName returns the path to a AST cache file for a given filename. This
// also ensures the cache directory exists.
func CachedASTName(fn string) string {
	os.Mkdir(cacheDir, os.ModePerm)
	return fmt.Sprintf(
		"%s%c%s.skol_ast",
		cacheDir, os.PathSeparator, cacheSafeName(fn))
}
