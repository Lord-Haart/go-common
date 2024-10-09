// 使用PBEWITHHMACSHA512ANDAES_256算法进行加密和解密。

package encryption

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"hash"
)

// getDerivedKey 生成派生Key，用于执行后续解密加密。
func getDerivedKey(prf hash.Hash, salt []byte, count, keySize int) (dk []byte) {
	hSize := prf.BlockSize()

	intL := (keySize + hSize - 1) / hSize // ceiling
	intR := keySize - (intL-1)*hSize      // residue
	ti := make([]byte, hSize)

	for i := 1; i <= intL; i++ {
		prf.Write(salt)

		ibytes := make([]byte, 4)
		ibytes[3] = byte(i)
		ibytes[2] = byte((i >> 8) & 0xff)
		ibytes[1] = byte((i >> 16) & 0xff)
		ibytes[0] = byte((i >> 24) & 0xff)
		prf.Write(ibytes)

		ui := prf.Sum(nil)
		prf.Reset()
		copy(ti, ui)

		for j := 2; j <= count; j++ {
			prf.Write(ui)
			ui = prf.Sum(nil)
			prf.Reset()
			for k := 0; k < len(ui); k++ {
				ti[k] ^= ui[k]
			}
		}

		dk = make([]byte, (i-1)*hSize)
		if i == intL {
			dk = append(dk, ti[:intR]...)
		} else {
			dk = append(dk, ti[:hSize]...)
		}
	}

	return
}

func pad(src []byte, blockSize int) []byte {
	padNum := blockSize - len(src)%blockSize
	padding := bytes.Repeat([]byte{byte(padNum)}, padNum)
	return append(src, padding...)
}

func unpad(src []byte, blockSize int) []byte {
	length := len(src)
	if length == 0 {
		return src
	}
	padding := int(src[length-1])
	if padding > length || padding > blockSize {
		return src
	}
	for i := 0; i < padding; i++ {
		if int(src[length-1-i]) != padding {
			return src
		}
	}
	return src[:length-padding]
}

func aes256Encrypt(origData, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	origData = pad(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, iv)
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

func aes256Decrypt(crypted, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockMode := cipher.NewCBCDecrypter(block, iv)
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = unpad(origData, block.BlockSize())
	return origData, nil
}

// decryptWithPBEHmacSHA512AndAES_256 使用指定的密钥进行解密。密文应当是经过base64编码的。
func decryptWithPBEHmacSHA512AndAES_256(msg string, password []byte) (string, error) {
	if msgBytes, err := base64.StdEncoding.DecodeString(msg); err != nil {
		return "", err
	} else {
		keySize := 32  // key size of AES256
		saltSize := 16 // salt size of AES256
		if len(msgBytes) <= saltSize*2 {
			return "", fmt.Errorf("invalid encrypted data(length: %d)", len(msgBytes))
		}
		salt := msgBytes[:saltSize]
		iv := msgBytes[saltSize : saltSize+saltSize]
		encText := msgBytes[saltSize+saltSize:]

		dk := getDerivedKey(hmac.New(sha512.New, password), salt, 1000, keySize)

		if text, err := aes256Decrypt(encText, dk, iv); err != nil {
			return "", err
		} else {
			return string(text), nil
		}
	}
}

// encryptWithPBEHmacSHA512AndAES_256 使用指定的密钥进行加密。加密结果是经过base64编码的。
func encryptWithPBEHmacSHA512AndAES_256(msg string, password []byte) (string, error) {
	keySize := 32  // key size of AES256
	saltSize := 16 // salt size of AES256
	salt := make([]byte, saltSize)
	rand.Read(salt)
	// TODO: 使用随机IV
	iv := []byte{43, 84, 101, 91, 46, 4, 56, 110, 101, 19, 127, 42, 92, 32, 99, 57}

	dk := getDerivedKey(hmac.New(sha512.New, password), salt, 1000, keySize)

	if text, err := aes256Encrypt([]byte(msg), dk, iv); err != nil {
		return "", err
	} else {
		r := append(salt, iv...)
		r = append(r, text...)
		return base64.StdEncoding.EncodeToString(r), nil
	}
}
