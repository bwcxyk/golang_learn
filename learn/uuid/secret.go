/*
@Author : YaoKun
@Time : 2022/12/29 16:46
*/

package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/google/uuid"
	"sort"
	"strconv"
	"strings"
	//_uuid "myproject/uuid"
)

const (
	ServerName = "MyProject_abc123"
)

var (
	chars = []string{"a", "b", "c", "d", "e", "f",
		"g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s",
		"t", "u", "v", "w", "x", "y", "z", "0", "1", "2", "3", "4", "5",
		"6", "7", "8", "9", "A", "B", "C", "D", "E", "F", "G", "H", "I",
		"J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V",
		"W", "X", "Y", "Z"}
)

func getAppKey() string {
	// 1.
	uuidByte, err := uuid.NewUUID()
	if err != nil {
		panic("new uuid error")
	}
	MyUuid := strings.ReplaceAll(uuidByte.String(), "-", "")

	// 2.
	//uuid := _uuid.Rand().Hex()
	//uuid = strings.ReplaceAll(uuid, "-", "")

	appKey := ""
	for i := 0; i < 8; i++ {
		str := MyUuid[i*4 : i*4+4]
		x, _ := strconv.ParseInt(str, 16, 64)
		appKey += chars[x%0x3e]
	}
	return appKey
}

func getAppSecret(appKey string) string {
	sli := []string{appKey, ServerName}
	sort.Strings(sli)
	str := strings.Join(sli, "")

	h := sha1.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func main() {
	appKey := getAppKey()
	appSecret := getAppSecret(appKey)
	fmt.Println("appKey", appKey)    // EO8bvxGl
	fmt.Println("appKey", appSecret) // 215bc9d7e1ab16c663289d40cea8164b9b00e307
}
