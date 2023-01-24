package utils

import (
	"path/filepath"
	"runtime"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
)

// ProjectRootDir Get path to this project root
// Got here: https://stackoverflow.com/questions/31873396/is-it-possible-to-get-the-current-root-of-package-structure-as-a-string-in-golan
func ProjectRootDir() (string, error) {
	// Depth of this file from project root
	currentDirDepth := 6

	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", errors.ErrGetRuntimeCaller
	}

	return filepath.Dir(file + pathUp(currentDirDepth)), nil
}

// Return chain of "../" to jump depth levels up
func pathUp(depth int) string {
	postfix := ""
	for i := 0; i < depth; i++ {
		postfix += "../"
	}

	return postfix
}
