package fs

import (
	"io"
	"os"
)

// 复制文件，仅文件，文件夹不可复制
func CopyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// 创建目标文件
	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	// 将源文件内容复制到目标文件
	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return err
	}

	return nil
}

// TODO: 复制文件夹到文件夹
func CopyDir(src, dst string) error {
	return nil
}
