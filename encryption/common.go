// PBE加密。

package encryption

import "strings"

var (
	password []byte
)

// 初始化`PBE算法`。
// p `PBE算法`使用的密钥，不能少于8个字符。
func InitPBE(p string) error {
	password = []byte(p)
	return nil
}

// 执行PBE加密。
// src 明文。
// 返回对应的密文。
func Encrypt(src string) (string, error) {
	return encryptWithPBEHmacSHA512AndAES_256(src, password)
}

// 执行PBE解密。
// src 密文。
// 返回对应的明文。
func Decrypt(src string) (string, error) {
	return decryptWithPBEHmacSHA512AndAES_256(src, password)
}

func TryDecrypt(src string) (string, error) {
	if strings.HasPrefix(src, "ENC(") && strings.HasSuffix(src, ")") {
		return Decrypt(src[4 : len(src)-1])
	} else {
		return src, nil
	}
}
