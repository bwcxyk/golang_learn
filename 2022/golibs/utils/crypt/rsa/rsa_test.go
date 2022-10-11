package rsa

import (
	"encoding/base64"
	"os"
	"testing"
)

func TestGenerateKey(t *testing.T) {
	t.Run("生成秘钥", func(t *testing.T) {
		GeneralRsaKey()
	})

	t.Run("验证加密与解密", func(t *testing.T) {
		txt := []byte("hello")
		pub, err := os.ReadFile(publicFileName)
		if err != nil {
			panic(err)
		}
		etxt, err := Encrypt(txt, pub)
		if err != nil {
			panic(err)
		}
		prv, err := os.ReadFile(privateFileName)
		if err != nil {
			panic(err)
		}
		dtxt, err := Decrypt(etxt, prv)
		if err != nil {
			panic(err)
		}
		if string(dtxt) != string(txt) {
			t.Error("解密失败")
		}
	})

	t.Run("解密验证", func(t *testing.T) {
		private := []byte(`-----BEGIN RSA PRIVATE KEY-----
MIICWwIBAAKBgQC+l2SrPNcC/UwgTrcfpe82TDL3Gp3fsWPe5LeGfCxLL6Bp8etl
CGHQB3sxMSV+6G+YEaVd6RQY3cOR10Ejy0MaG4pe9iI16jLEya8gAGXrJAksSycc
AMkWB1G3o4YH8lTcHTnn6ggshlqSquZEiOAT/JJowrZF1k8+Ilu6RjhgFQIDAQAB
AoGAGlhM9wqS4fvnYPPghkRmm4fO569LMVeQ6YeOAs35RW9Q6jQhiLau5pWSJsuX
VkKE4m1WAXQtbf6BgRBTFcKMBLtaVAvwXxJp+XhU0th54Ogr7IC8DvoGV4X8+ZFC
AcxKpUCm/mPSrAXxmwkQa0wjdHFbC5RFvLBSqXw8m+n8XUECQQD4flTJ/k83vLiL
E/763e+C9MhljopYnvcnS/ahE7tcpMlSx6RUogc7vEynHS7u5uTWu1z2Isd5QRu5
ds1q1wXxAkEAxFlKAJv5YevOoYE9ExM7v3aLKJUoY5lG2j5XZfzEgh6xWD52tatU
GXxraTkfMZ6DI7TwaDDJDQPo9jABcOuIZQJAOgQcvbBPRH8eQvJfNKj+U3/dDcDy
0fADWjxlD4Rw2RdsHZSY7F2L/MlgyB+LJqHnya6i7KCAG/I0i9+N5CThsQJABmoY
Vc7CpeVLOdF8Ckx9jjK0Nx4wVJikTOrudgF89fdEuepIfITdWXvMEoLwNvHDvP3j
QLQfTVKMYMaOtX6sVQJAcZqH1xVVkJwPI/eU5Nm4I5OstDkZ4hv48hUwO6ZrFnue
11lv0VjsOnH6189jBWTxtI+SFBmP4GYNNUh9mr+nXA==
-----END RSA PRIVATE KEY-----`)
		//plain := "ASCGwszKOFh7XxhHJI0SL0cd3iTDp2TPIJX9WwouzmwMirjneZAaYwUw7iT5whoIy9kOhreDvTG/HWPDArZ8weFhbLye0jpmGnlAtHHUqtFCGqZzJLj5pPpfnjiN1Llowvmajwk5wsALAGcVM8V7wA4+FKJPUo59acyynuLEVbE="
		plain := "PRnChFItl7MJD80+eMoHaOMR67jqcJ04iCF0mwkIYbGIiFEY6vqwG6jFfzU994d9QhKTvEdLFCtvDu4r59iEhhWF8/X1LqvD091be7OSKVxarUcVSKwAVpgnPYlRjC1gjmyiBVDTWBsWGDaGKmftg7aYz9V134lR36C6W1Fu/gk="
		dplain, err := base64.StdEncoding.DecodeString(plain)
		txt, err := Decrypt(dplain, private)
		if err != nil {
			panic(err)
		}
		t.Log(string(txt))
	})
}
