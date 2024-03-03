package zip

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// 压缩文件夹内的所有文件到目标，不包括最外层文件夹
func ZipFolderInner(folderPath, dst string) error {
	src, _ := filepath.Abs(folderPath)
	dst, _ = filepath.Abs(dst)

	fw, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer fw.Close()

	zw := zip.NewWriter(fw)
	defer func() {
		if err := zw.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	err = filepath.Walk(src, func(path string, fi os.FileInfo, errBack error) (err error) {
		if path == src {
			return nil
		}

		if errBack != nil {
			return errBack
		}

		// 通过文件信息，创建 zip 的文件信息
		fh, err := zip.FileInfoHeader(fi)
		if err != nil {
			return err
		}

		// 替换文件信息中的文件名
		fh.Name = strings.TrimPrefix(path, src)
		if strings.HasPrefix(fh.Name, "\\") {
			fh.Name = strings.TrimPrefix(fh.Name, "\\")
		} else {
			fh.Name = strings.TrimPrefix(fh.Name, "/")
		}

		//统一整成linux的分隔符形式
		fh.Name = strings.ReplaceAll(fh.Name, "\\", "/")

		// 这步开始没有加，会发现解压的时候说它不是个目录
		if fi.IsDir() {
			fh.Name += "/"
		}
		fh.Method = zip.Deflate

		// 写入文件信息，并返回一个 Write 结构
		w, err := zw.CreateHeader(fh)
		if err != nil {
			return err
		}

		// 检测，如果不是标准文件就只写入头信息，不写入文件数据到 w, 如目录，也没有数据需要写
		if !fh.Mode().IsRegular() {
			return nil
		}

		// 打开要压缩的文件
		fr, err := os.Open(path)
		if err != nil {
			return
		}
		defer fr.Close()

		// 将打开的文件 Copy 到 w
		_, err = io.Copy(w, fr)
		return err
	})
	if err != nil {
		return err
	}

	return nil
}

// 解压文件到指定目录
func Unzip(src, dst string) error {
	// 1. Open the zip file
	reader, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer reader.Close()

	// 2. Get the absolute destination path
	dst, err = filepath.Abs(dst)
	if err != nil {
		return err
	}

	// 加压方法
	decompresser := func(f *zip.File, destination string) error {
		// 4. Check if file paths are not vulnerable to Zip Slip
		filePath := filepath.Join(destination, f.Name)
		if !strings.HasPrefix(filePath, filepath.Clean(destination)+string(os.PathSeparator)) {
			return fmt.Errorf("invalid file path: %s", filePath)
		}

		// 5. Create directory tree
		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
				return err
			}
			return nil
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return err
		}

		// 6. Create a destination file for unzipped content
		destinationFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		defer destinationFile.Close()

		// 7. Unzip the content of a file and copy it to the destination file
		zippedFile, err := f.Open()
		if err != nil {
			return err
		}
		defer zippedFile.Close()

		if _, err := io.Copy(destinationFile, zippedFile); err != nil {
			return err
		}
		return nil
	}

	// 3. Iterate over zip files inside the archive and unzip each of them
	for _, zf := range reader.File {
		err := decompresser(zf, dst)
		if err != nil {
			return err
		}
	}

	return nil
}
