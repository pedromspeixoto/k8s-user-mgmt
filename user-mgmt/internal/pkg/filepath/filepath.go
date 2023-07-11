package filepath

import (
	"path/filepath"
	"runtime"
)

func ProjectRootDir() string {
	_, b, _, _ := runtime.Caller(0)

	return filepath.Join(filepath.Dir(b), "../../..")
}
