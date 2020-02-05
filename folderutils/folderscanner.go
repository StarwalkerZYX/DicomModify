package folderutils

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

var childPaths []string

func ScanFolder(dicomFolder string) ([]string, error) {

	err := filepath.Walk(dicomFolder, collectFile)
	if err != nil {
		fmt.Printf("filepath.Walk() error: %v\n", err)
	}
	return childPaths, nil
}

func collectFile(path string, info os.FileInfo, err error) error {
	if info == nil {
		// 文件名称超过限定长度等其他问题也会导致info == nil
		// 如果此时return err 就会显示找不到路径，并停止查找。
		println("can't find:(" + path + ")")
		return errors.New("can not find")
	}

	if info.IsDir() {
		println("This is folder:(" + path + ")")
		return nil
	} else {

		println("This is file:(" + path + ")")
		childPaths = append(childPaths, path)
		return nil
	}
}
