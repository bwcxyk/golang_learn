package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"golibs/utils/crypt"
	"net/url"
)

// aesSha1prng 模拟 java SHA1PRNG 处理
func aesSha1prng(keyBytes []byte, encryptLength int) ([]byte, error) {
	hashs := crypt.Sha1Byte(crypt.Sha1Byte(keyBytes))
	maxLen := len(hashs)
	realLen := encryptLength / 8
	if realLen > maxLen {
		return nil, errors.New("invalid length")
	}
	return hashs[0:realLen], nil
}

func generateKey(key []byte) (genKey []byte) {
	genKey = make([]byte, 16)
	copy(genKey, key)
	for i := 16; i < len(key); {
		for j := 0; j < 16 && i < len(key); j, i = j+1, i+1 {
			genKey[j] ^= key[i]
		}
	}
	return genKey
}

// JavaEncryptECB 适用于 Java SHA1PRNG 的加密
func JavaEncryptECB(srcStr, keyStr string) (encryptedStr string, err error) {
	var src, encrypted, key []byte

	src = []byte(srcStr)
	key = []byte(keyStr)

	key, err = aesSha1prng(key, 128)
	if err != nil {
		return
	}
	key = generateKey(key)
	// 分组秘钥
	ciphers, err := aes.NewCipher(key)
	if err != nil {
		return
	}
	// urlencode
	src = []byte(url.QueryEscape(string(src)))
	length := (len(src) + aes.BlockSize) / aes.BlockSize
	plain := make([]byte, length*aes.BlockSize)
	copy(plain, src)
	pad := byte(len(plain) - len(src))
	for i := len(src); i < len(plain); i++ {
		plain[i] = pad
	}
	encrypted = make([]byte, len(plain))
	// 分组分块加密
	for bs, be := 0, ciphers.BlockSize(); bs <= len(src); bs, be = bs+ciphers.BlockSize(), be+ciphers.BlockSize() {
		ciphers.Encrypt(encrypted[bs:be], plain[bs:be])
	}
	// 转大写形式16进制并base64编码
	encryptedStr = string(encode(encrypted))
	return
}

// JavaDecryptECB 适用于 Java SHA1PRNG 的解密
func JavaDecryptECB(encryptedStr, keyStr string) (decryptedStr string, err error) {
	var encrypted, key, decrypted []byte

	encrypted = []byte(encryptedStr)
	key = []byte(keyStr)

	key, err = aesSha1prng(key, 128)
	if err != nil {
		return
	}
	key = generateKey(key)
	// 分组秘钥
	ciphers, err := aes.NewCipher(key)
	if err != nil {
		return
	}
	// 解码
	encrypted, err = decode(encrypted)
	if err != nil {
		return
	}
	decrypted = make([]byte, len(encrypted))
	for bs, be := 0, ciphers.BlockSize(); bs < len(encrypted); bs, be = bs+ciphers.BlockSize(), be+ciphers.BlockSize() {
		ciphers.Decrypt(decrypted[bs:be], encrypted[bs:be])
	}
	trim := 0
	if len(decrypted) > 0 {
		trim = len(decrypted) - int(decrypted[len(decrypted)-1])
	}
	decrypted = decrypted[:trim]
	// urldecode
	decryptedStr, err = url.QueryUnescape(string(decrypted))
	if err != nil {
		return
	}
	return
}

// JavaEncryptCBC 适用于 Java SHA1PRNG 的加密
func JavaEncryptCBC(srcStr, keyStr string) (encryptedStr string, err error) {
	var src, encrypted, key []byte

	src = []byte(srcStr)
	key = []byte(keyStr)

	// 适配 Java SHA1PRNG 模块，生成秘钥
	key, err = aesSha1prng(key, 128)
	if err != nil {
		return
	}
	key = generateKey(key)
	// 分组秘钥
	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// urlencode
	src = []byte(url.QueryEscape(string(src)))
	// 填充
	src = pkcs5Padding(src, blockSize)
	// 初始化向量 对应
	iv := make([]byte, blockSize)
	// 加密模式
	blockMode := cipher.NewCBCEncrypter(block, iv)
	// 创建数组
	encrypted = make([]byte, len(src))
	// 执行加密
	blockMode.CryptBlocks(encrypted, src)
	// 转大写形式16进制并base64编码
	encryptedStr = string(encode(encrypted))
	return
}

// JavaDecryptCBC 适用于 Java SHA1PRNG 的解密
func JavaDecryptCBC(encryptedStr, keyStr string) (decryptedStr string, err error) {
	var encrypted, key, decrypted []byte

	encrypted = []byte(encryptedStr)
	key = []byte(keyStr)

	// 适配 Java SHA1PRNG 模块，生成秘钥
	key, err = aesSha1prng(key, 128)
	if err != nil {
		return
	}
	key = generateKey(key)
	// 分组秘钥
	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 初始化向量 对应
	iv := make([]byte, blockSize)
	// 加密模式
	blockMode := cipher.NewCBCDecrypter(block, iv)
	// 解码
	encrypted, err = decode(encrypted)
	if err != nil {
		return
	}
	// 创建数组
	decrypted = make([]byte, len(encrypted))
	// 执行加密
	blockMode.CryptBlocks(decrypted, encrypted)
	// 去除补全码
	decrypted = pkcs5UnPadding(decrypted)
	// urldecode
	decryptedStr, err = url.QueryUnescape(string(decrypted))
	if err != nil {
		return
	}
	return
}
