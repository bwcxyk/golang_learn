/*
   Copyright (c) [2021] IT.SOS
   golibs is licensed under Mulan PSL v2.
   You can use this software according to the terms and conditions of the Mulan PSL v2.
   You may obtain a copy of Mulan PSL v2 at:
            http://license.coscl.org.cn/MulanPSL2
   THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
   See the Mulan PSL v2 for more details.
*/

package crypt

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
)

func Md5(s string) string {
	return hex.EncodeToString(Md5Byte([]byte(s)))
}

func Md5Byte(s []byte) []byte {
	h := md5.New()
	h.Write(s)
	return h.Sum(nil)
}

func Sha1(s string) string {
	return hex.EncodeToString(Sha1Byte([]byte(s)))
}

func Sha1Byte(data []byte) []byte {
	h := sha1.New()
	h.Write(data)
	return h.Sum(nil)
}
