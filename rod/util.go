package rodutil

import (
	"path/filepath"

	"github.com/duke-git/lancet/v2/datetime"
	"github.com/duke-git/lancet/v2/fileutil"
)

// GetUnixTime
func GetUnixTime() int64 {
	return datetime.NewUnixNow().ToUnix()
}

// AppendPath append path
func AppendPath(parent string, path ...string) string {
	p := fileutil.CurrentPath()
	p = filepath.Join(p, parent)

	if !fileutil.IsExist(p) {
		_ = fileutil.CreateDir(p)
	}
	return filepath.Join(p, filepath.Join(path...))
}

// AddTmp add res dir
func AddRes(path ...string) string {
	return AppendPath("res", path...)
}

// AddTmp add tmp dir
func AddTmp(path ...string) string {
	return AppendPath("tmp", path...)
}
