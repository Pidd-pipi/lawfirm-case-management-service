package util

import (
	"fmt"
	"os"
)

// DeleteFiles 批量删除文件。
func DeleteFiles(paths []string) error {
	var err error
	for _, p := range paths {
		if p == "" {
			continue
		}
		defer func(p string) {
			if e := os.Remove(p); e != nil {
				err = fmt.Errorf("delete %s: %v", p, e)
			}
		}(p)
	}
	return err
}
