package auth

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"golibs/cerrors"
	"golibs/config"
	"golibs/framework/iris/bootstrap"
	"golibs/framework/iris/caches"
	"golibs/utils/crypt"
	"golibs/utils/crypt/rsa"
	"github.com/kataras/golog"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

// CheckSign 验证签名
func CheckSign(b *bootstrap.Bootstrapper) {
	b.UseGlobal(func(ctx iris.Context) {
		if ctx.GetStatusCode() < iris.StatusInternalServerError && ctx.Method() != "OPTIONS" {
			if SignExclude(ctx) {
				ctx.Next()
				return
			}
			var mJson []byte
			contentType := ctx.GetContentTypeRequested()
			if contentType == context.ContentFormMultipartHeaderValue ||
				contentType == context.ContentFormHeaderValue {
				res := ctx.FormValues()
				if len(res) > 0 {
					mJson = toJson(res)
				}
			} else {
				body, _ := ctx.GetBody()
				uParams := ctx.URLParams()
				params := map[string]interface{}{}
				if len(body) > 0 {
					err := json.Unmarshal(body, &params)
					if err != nil {
						golog.Error(err)
						return
					}
				}
				for k, v := range uParams {
					params[k] = v
				}
				if len(params) > 0 {
					mJson = toJson(params)
				}
			}
			data := filter(mJson)

			token := ctx.GetHeader("token")
			sign := ctx.GetHeader("sign")
			ts := ctx.GetHeader("ts")
			nonce := ctx.GetHeader("nonce")

			// 300s 五分钟内有效
			tss, _ := strconv.Atoi(ts)
			if time.Now().Unix() > (int64(tss) + 300) {
				forbidden(ctx)
				golog.Errorf("sign fail: time expired.[token:%s]", token)
				return
			}

			dPlain, _ := base64.StdEncoding.DecodeString(nonce)
			priv := []byte(config.Config.GetCryptRsaPriv())
			key, err := rsa.Decrypt(dPlain, priv)
			if err != nil {
				forbidden(ctx)
				golog.Errorf("sign fail: nonce decode fail.[token:%s]", token)
				return
			}

			if crypt.Md5(token+ts+string(key)+data+nonce) != sign {
				forbidden(ctx)
				golog.Errorf("sign fail: not equals.[token:%s]", token)
				return
			}

			if !caches.SignSet(sign) {
				forbidden(ctx)
				golog.Errorf("sign is exists.[token:%s]", token)
				return
			}

			ctx.Next()
		}
	})
}

func SignExclude(ctx iris.Context) bool {
	ruleSign := config.Config.GetSignatureExclude()
	if len(ruleSign) == 0 {
		return false
	}
	uri := ctx.Request().RequestURI
	for _, r := range ruleSign {
		if is, _ := regexp.MatchString(r, uri); is {
			return true
		}
	}
	return false
}

func toJson(o interface{}) []byte {
	bf := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(bf)
	jsonEncoder.SetEscapeHTML(false)
	jsonEncoder.Encode(o)
	return bytes.Trim(bf.Bytes(), "\n")
}

func forbidden(ctx iris.Context) {
	ctx.StopWithError(iris.StatusForbidden, cerrors.Error("forbidden_access"))
}

func filter(str []byte) string {
	mString := string(str)
	mString = strings.ReplaceAll(mString, "{", "")
	mString = strings.ReplaceAll(mString, "}", "")
	mString = strings.ReplaceAll(mString, "[", "")
	mString = strings.ReplaceAll(mString, "]", "")
	mString = strings.ReplaceAll(mString, ",", "")
	mString = strings.ReplaceAll(mString, "\"", "")
	mString = strings.ReplaceAll(mString, ":", "")
	sString := strings.SplitN(mString, "", len(mString))
	sort.Strings(sString)
	return strings.Join(sString, "")
}
