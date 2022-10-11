/*
   Copyright (c) [2021] IT.SOS
   golibs is licensed under Mulan PSL v2.
   You can use this software according to the terms and conditions of the Mulan PSL v2.
   You may obtain a copy of Mulan PSL v2 at:
            http://license.coscl.org.cn/MulanPSL2
   THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
   See the Mulan PSL v2 for more details.
*/

package aes

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"strings"
)

// pkcs5Padding 填充明文
func pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// pkcs5UnPadding 去除填充数据
func pkcs5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// 转大写形式16进制并base64编码
func encode(src []byte) []byte {
	return []byte(base64.StdEncoding.EncodeToString(
		[]byte(strings.ToUpper(
			hex.EncodeToString(src)))))
}

// 解base64/转小写/转byte
func decode(encrypted []byte) (decrypted []byte, err error) {
	// 解 base64
	encrypted, err = base64.StdEncoding.DecodeString(string(encrypted))
	if err != nil {
		return
	}
	// to byte
	decrypted, err = hex.DecodeString(
		strings.ToLower(string(encrypted)))
	if err != nil {
		return
	}
	return
}
