/*
   Copyright (c) [2021] IT.SOS
   golibs is licensed under Mulan PSL v2.
   You can use this software according to the terms and conditions of the Mulan PSL v2.
   You may obtain a copy of Mulan PSL v2 at:
            http://license.coscl.org.cn/MulanPSL2
   THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
   See the Mulan PSL v2 for more details.
*/

package services

import (
	"fmt"
	"golibs/cerrors"
	"golibs/config"
	"golibs/utils/crypt/aes"
	"log"
	"strconv"
	"strings"
	"time"
)

func GetLoginId(token string) (loginId string, err error) {
	if token == "" {
		err = cerrors.Error("unauthorized_access")
		return
	}
	deToken, err := aes.JavaDecryptCBC(token, config.Config.GetCryptAesToken())
	if err != nil {
		log.Panicf("token aes decode fail: %v, data: %v", err, token)
	}
	tokenSplit := strings.Split(deToken, "|")
	if len(tokenSplit) != 2 {
		log.Panicf("token is error: %v", deToken)
	}
	unix, err := strconv.ParseInt(tokenSplit[1], 10, 64)
	if err != nil {
		log.Panicf("token unix time parseInt err: %v, data: %v", err, tokenSplit[1])
	}
	if time.Now().After(time.Unix(unix, 0)) {
		err = cerrors.Error("unauthorized_access")
	} else {
		loginId = tokenSplit[0]
	}
	return
}

func GetToken(loginId string, lifetime time.Duration) string {
	ttl := time.Now()
	if lifetime > 0 {
		ttl = ttl.Add(lifetime)
	} else {
		// 默认一周
		ttl = ttl.Add(time.Hour * 24 * 7)
	}
	tk := fmt.Sprintf("%s|%d", loginId, ttl.Unix())
	token, err := aes.JavaEncryptCBC(tk, config.Config.GetCryptAesToken())
	if err != nil {
		log.Panicf("token aes encode fail: %v, data: %v", err, tk)
	}
	return token
}
