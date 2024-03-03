package fs

import (
	"encoding/json"
	"io"
	"os"
)

func Exists(path string) bool {
	_, e := os.Stat(path)
	if e != nil {
		return false
	}
	return true
}

// 文件状态
func Stat(path string) (os.FileInfo, error) {
	return os.Stat(path)
}

// 读文件
func ReadFile(fileName string) ([]byte, error) {
	return os.ReadFile(fileName)
}

// 写文件
func WriteFile(fileName string, content string) error {
	var (
		file *os.File
		err  error
	)

	fStat, err := Stat(fileName)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	}

	// 是文件
	if !fStat.IsDir() {
		file, err = os.OpenFile(fileName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	} else {
		// 是文件夹
		file, err = os.Create(fileName)
	}
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.WriteString(file, content)
	if err != nil {
		return err
	}

	return nil
}

// 读取 json 文件，dst 需传入指针
func ReadJson(jsonFilePath string, dst any) error {
	content, err := ReadFile(jsonFilePath)
	if err != nil {
		return err
	}

	return json.Unmarshal(content, dst)
}
