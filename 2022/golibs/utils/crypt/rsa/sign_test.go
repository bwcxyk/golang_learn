package rsa

import (
	"encoding/hex"
	"fmt"
	"os"
	"testing"
)

func TestSign(t *testing.T) {
	t.Run("数字签名", func(t *testing.T) {
		publicKey, err := os.ReadFile(publicFileName)
		if err != nil {
			panic(err)
		}
		privateKey, err := os.ReadFile(privateFileName)
		if err != nil {
			panic(err)
		}
		msg := []byte("RSA数字签名测试")
		signmsg, err := Sign(msg, privateKey)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("RSA数字签名的消息为：", hex.EncodeToString(signmsg))

		// 验证数字签名正不正确
		result := VerifySign(msg, signmsg, publicKey)
		if result { // 如果result返回的是 true 那么就是本人签名，否则不是，只有私钥加密，相对的公钥验证才可以认为是本人
			fmt.Println("RSA数字签名正确，是本人")
		} else {
			fmt.Println("RSA数字签名错误，不是本人")
		}
	})
}
