package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"fmt"
	"io"
	"os"
)

// func padding(content []byte, blockSize int) []byte {
// 	//计算要填充的长度
// 	n := blockSize - len(content)%blockSize
// 	//对原来的明文填充n个n
// 	temp := bytes.Repeat([]byte{byte(n)}, n)
// 	content = append(content, temp...)
// 	return content
// }

// 填充，加密时使用
func pkcs7padding(content []byte, blockSize int) []byte {
	padding := blockSize - len(content)%blockSize
	fmt.Println("padding", padding)
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(content, padtext...)
}

// // 去填充，解密时使用
// func pkcs7unpadding(content []byte) ([]byte, error) {
// 	//获取数据长度
// 	length := len(content)
// 	if length == 0 {
// 		return nil, errors.New("pkcs7unpadding error")
// 	}

// 	//获取填充字符串长度
// 	unpadding := int(content[length-1])
// 	//截取切片，删除填充字节，并且返回明文
// 	return content[:(length - unpadding)], nil
// }

func pkcs7unpadding(content []byte) ([]byte, error) {
	//获取数据长度
	length := len(content)
	if length == 0 {
		return nil, errors.New("pkcs7unpadding error: content is empty")
	}

	//获取填充字节
	padding := int(content[length-1])
	if padding == 0 || padding > length {
		return content, errors.New("pkcs7unpadding error: invalid padding")
	}

	//检查填充字节是否合法，不合法则说明无需 padding
	for i := length - padding; i < length-1; i++ {
		if content[i] != byte(padding) {
			return content, errors.New("pkcs7unpadding error: invalid padding bytes")
		}
	}

	//截取切片，删除填充字节，并且返回明文
	return content[:length-padding], nil
}

type cbc struct {
	key string
	iv  string
}

func NewCbc(key string, iv string) *cbc {
	return &cbc{key, iv}
}

// 流式加密文件
func (c *cbc) EncryptFileStream(src, dst string) error {
	key := c.key
	iv := c.iv

	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return err
	}
	mode := cipher.NewCBCEncrypter(block, []byte(iv))

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	blockSize := 100 * block.BlockSize()
	for {
		buf := make([]byte, blockSize)
		n, err := in.Read(buf)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		if n == 0 {
			break
		}

		// padding
		if n < blockSize {
			buf = pkcs7padding(buf[:n], block.BlockSize())
		}

		mode.CryptBlocks(buf, buf)
		_, err = out.Write(buf[:])
		if err != nil {
			return err
		}
	}
	return nil
}

// 整体加密（文件内容加载到内存后加密）
func (c *cbc) EncryptFile(src, dst string) error {
	in, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	out, err := c.Encrypt(in)
	if err != nil {
		return err
	}

	return os.WriteFile(dst, out, os.ModePerm)
}

// 流式解密文件
func (c *cbc) DecryptFileStream(src, dst string) error {
	key := c.key
	iv := c.iv

	input, err := os.Open(src)
	if err != nil {
		return err
	}
	defer input.Close()

	inputStat, err := input.Stat()
	if err != nil {
		return err
	}
	inputFileSize := inputStat.Size()

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return err
	}
	mode := cipher.NewCBCDecrypter(block, []byte(iv))

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	var fileSize int64 = 0

	blockSize := 100 * block.BlockSize()
	var lastBlock []byte
	for {
		buf := make([]byte, blockSize)
		n, err := input.Read(buf)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		if n == 0 {
			break
		}

		fileSize += int64(n)
		if fileSize < inputFileSize {
			mode.CryptBlocks(buf, buf)
			_, err = out.Write(buf)
			if err != nil {
				return err
			}
		} else {
			lastBlock = buf[:n]
			break
		}
	}

	if len(lastBlock) > 0 {
		mode.CryptBlocks(lastBlock, lastBlock)
		// 去除填充
		lastBlock, _ = pkcs7unpadding(lastBlock)
		_, err = out.Write(lastBlock[:])
		if err != nil {
			return err
		}
	}

	return nil
}

// 整体解密（文件内容加载到内存后解密）
func (c *cbc) DecryptFile(src, dst string) error {
	in, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	out, err := c.Decrypt(in)
	if err != nil {
		return err
	}

	return os.WriteFile(dst, out, os.ModePerm)
}

func (c *cbc) Decrypt(content []byte) ([]byte, error) {
	key := c.key
	iv := c.iv

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}

	blockMode := cipher.NewCBCDecrypter(block, []byte(iv))
	result := make([]byte, len(content))
	blockMode.CryptBlocks(result, content)
	result, err = pkcs7unpadding(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *cbc) Encrypt(content []byte) ([]byte, error) {
	key := c.key
	iv := c.iv

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}
	content = pkcs7padding(content, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, []byte(iv))
	result := make([]byte, len(content))
	blockMode.CryptBlocks(result, content)
	return result, nil
}
