package zip

import (
	// "archive/zip"
	"fmt"
	"io"
	"os"

	"path/filepath"
	"strings"

	"github.com/klauspost/compress/zip"
)

// 压缩文件夹内容，速度更快
func CompressFolderContent(src string, dst string) error {
	src, _ = filepath.Abs(src)
	dst, _ = filepath.Abs(dst)

	fw, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer fw.Close()

	zw := zip.NewWriter(fw)
	defer zw.Close()

	err = filepath.Walk(src, func(path string, fi os.FileInfo, errBack error) (err error) {
		if path == src {
			return nil
		}

		if errBack != nil {
			return errBack
		}

		return makeCompress(zw, path, src, fi)

		// // 通过文件信息，创建 zip 的文件信息
		// fh, err := zip.FileInfoHeader(fi)
		// if err != nil {
		// 	return err
		// }

		// // 替换文件信息中的文件名
		// fh.Name = strings.TrimPrefix(path, src)
		// if strings.HasPrefix(fh.Name, "\\") {
		// 	fh.Name = strings.TrimPrefix(fh.Name, "\\")
		// } else {
		// 	fh.Name = strings.TrimPrefix(fh.Name, "/")
		// }

		// //统一整成linux的分隔符形式
		// fh.Name = strings.ReplaceAll(fh.Name, "\\", "/")

		// // 这步开始没有加，会发现解压的时候说它不是个目录
		// if fi.IsDir() {
		// 	fh.Name += "/"
		// }
		// fh.Method = zip.Deflate

		// // 写入文件信息，并返回一个 Write 结构
		// w, err := zw.CreateHeader(fh)
		// if err != nil {
		// 	return err
		// }

		// // 检测，如果不是标准文件就只写入头信息，不写入文件数据到 w, 如目录，也没有数据需要写
		// if !fh.Mode().IsRegular() {
		// 	return nil
		// }

		// // 打开要压缩的文件
		// fr, err := os.Open(path)
		// if err != nil {
		// 	return
		// }
		// defer fr.Close()

		// // 将打开的文件 Copy 到 w
		// _, err = io.Copy(w, fr)
		// return err
	})
	if err != nil {
		return err
	}

	return nil
}

// 解压缩，速度较快
func UncompressArchive(src, dst string) error {
	// 1. Open the zip file
	reader, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer reader.Close()

	dst, err = filepath.Abs(dst)
	if err != nil {
		return err
	}

	for _, file := range reader.File {
		if err := makeUncompress(file, dst); err != nil {
			return err
		}
	}

	return nil
}

func makeCompress(zw *zip.Writer, fpath string, root string, fi os.FileInfo) error {
	// 通过文件信息，创建 zip 的文件信息
	dstHeader, err := zip.FileInfoHeader(fi)
	if err != nil {
		return err
	}

	// 替换文件信息中的文件名
	dstHeader.Name = strings.TrimPrefix(fpath, root)
	if strings.HasPrefix(dstHeader.Name, "\\") {
		dstHeader.Name = strings.TrimPrefix(dstHeader.Name, "\\")
	} else {
		dstHeader.Name = strings.TrimPrefix(dstHeader.Name, "/")
	}

	//统一整成linux的分隔符形式
	dstHeader.Name = strings.ReplaceAll(dstHeader.Name, "\\", "/")

	// 这步开始没有加，会发现解压的时候说它不是个目录
	if fi.IsDir() {
		dstHeader.Name += "/"
	}
	dstHeader.Method = zip.Deflate

	// 写入文件信息，并返回一个 Write 结构
	dstWriter, err := zw.CreateHeader(dstHeader)
	if err != nil {
		return err
	}

	if !dstHeader.Mode().IsRegular() {
		return nil
	}

	srcReader, err := os.Open(fpath)
	if err != nil {
		return err
	}
	defer srcReader.Close()

	if _, err := io.Copy(dstWriter, srcReader); err != nil {
		return err
	}

	return nil
}

func makeUncompress(f *zip.File, dst string) error {
	filePath := filepath.Join(dst, f.Name)
	if !strings.HasPrefix(filePath, filepath.Clean(dst)+string(os.PathSeparator)) {
		return fmt.Errorf("invalid file path(%s)", filePath)
	}

	if f.FileInfo().IsDir() {
		if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
			return err
		}
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return err
	}

	dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		return err
	}
	defer dstFile.Close()

	srcFile, err := f.Open()
	if err != nil {
		return err
	}
	defer srcFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	return nil
}
