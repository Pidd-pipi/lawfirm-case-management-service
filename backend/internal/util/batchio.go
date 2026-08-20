package util

import (
	"fmt"
	"os"
)

// DeleteFiles 批量删除文件。中途某个文件删除失败时仍会继续删除剩余文件，
// 返回遇到的第一个错误。直接调用 os.Remove 而非 defer，避免大批量删除时
// 操作堆积到函数返回、迟迟不释放资源。
func DeleteFiles(paths []string) error {
	var first error
	for _, p := range paths {
		if p == "" {
			continue
		}
		if e := os.Remove(p); e != nil && first == nil {
			first = fmt.Errorf("delete %s: %w", p, e)
		}
	}
	return first
}
