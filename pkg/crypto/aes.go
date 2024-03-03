package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
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
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(content, padtext...)
}

// 去填充，解密时使用
func pkcs7unpadding(content []byte) ([]byte, error) {
	//获取数据长度
	length := len(content)
	if length == 0 {
		return nil, errors.New("pkcs7unpadding error")
	}

	//获取填充字符串长度
	unpadding := int(content[length-1])
	//截取切片，删除填充字节，并且返回明文
	return content[:(length - unpadding)], nil
}

func pkcs7unpadding2(content []byte) ([]byte, error) {
	//获取数据长度
	length := len(content)
	if length == 0 {
		return nil, errors.New("pkcs7unpadding error: content is empty")
	}

	//获取填充字节
	padding := int(content[length-1])
	if padding == 0 || padding > length {
		return nil, errors.New("pkcs7unpadding error: invalid padding")
	}

	//检查填充字节是否合法
	for i := length - padding; i < length-1; i++ {
		if content[i] != byte(padding) {
			return nil, errors.New("pkcs7unpadding error: invalid padding bytes")
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

// FIXME: 加解密后文件末尾出现 \0x00 等字符
func (c *cbc) EncryptFile(src, dst string) error {
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

	blockSize := block.BlockSize()
	for {
		buf := make([]byte, blockSize)
		n, err := io.ReadFull(in, buf)
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
			buf = pkcs7padding(buf, blockSize)
		}

		bufRes := []byte{}
		mode.CryptBlocks(bufRes, buf)
		_, err = out.Write(bufRes)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *cbc) EncryptFileFull(src, dst string) error {
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

// func (c *cbc) DecryptFile(src, dst string) error {
// 	key := c.key
// 	iv := c.iv

// 	caeFile, err := os.Open(src)
// 	if err != nil {
// 		return err
// 	}
// 	defer caeFile.Close()

// 	block, err := aes.NewCipher([]byte(key))
// 	if err != nil {
// 		return err
// 	}
// 	mode := cipher.NewCBCDecrypter(block, []byte(iv))

// 	out, err := os.Create(dst)
// 	if err != nil {
// 		return err
// 	}
// 	defer out.Close()
// 	var fileSize int64 = 0

// 	blockSize := block.BlockSize()
// 	for {
// 		buf := make([]byte, blockSize)
// 		n, err := io.ReadFull(caeFile, buf)
// 		if err == io.EOF {
// 			break
// 		} else if err != nil && err != io.ErrUnexpectedEOF {
// 			return err
// 		}

// 		if n == 0 {
// 			break
// 		}

// 		mode.CryptBlocks(buf, buf)
// 		_, err = out.Write(buf)
// 		if err != nil {
// 			return err
// 		}
// 		fileSize += int64(n)
// 	}

// 	lastBuf := make([]byte, blockSize)
// 	offset := fileSize - int64(blockSize)
// 	_, err = out.Seek(offset, 0)
// 	if err != nil {
// 		return err
// 	}
// 	_, err = out.Read(lastBuf)
// 	if err != nil {
// 		return err
// 	}

// 	lastBuf, err = pkcs7unpadding(lastBuf)
// 	if err != nil {
// 		return err
// 	}

// 	err = out.Truncate(offset)
// 	if err != nil {
// 		return err
// 	}
// 	_, err = out.Write(lastBuf)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// FIXME: 加解密后文件末尾出现 \0x00 等字符
func (c *cbc) DecryptFile(src, dst string) error {
	key := c.key
	iv := c.iv

	caeFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer caeFile.Close()

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

	blockSize := block.BlockSize()
	var lastBlock []byte
	for {
		buf := make([]byte, blockSize)
		n, err := io.ReadFull(caeFile, buf)
		if err == io.EOF {
			break
		} else if err != nil && err != io.ErrUnexpectedEOF {
			return err
		}

		if n == 0 {
			break
		}

		mode.CryptBlocks(buf, buf)
		if n < blockSize {
			lastBlock = buf // 保存最后一个块
			continue
		}
		_, err = out.Write(buf)
		if err != nil {
			return err
		}
		fileSize += int64(n)
	}

	// 去除填充
	if len(lastBlock) > 0 {
		lastBlock, err = pkcs7unpadding2(lastBlock)
		if err != nil {
			return err
		}
		_, err := out.Write(lastBlock)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *cbc) DecryptFileFull(src, dst string) error {
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
