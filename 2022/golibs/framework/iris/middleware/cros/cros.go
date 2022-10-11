package cros

import (
	"golibs/config"
	"golibs/framework/iris/bootstrap"
	"golibs/utils/array"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/accesslog"
	"strings"
)

// Cros 跨域处理
func Cros(b *bootstrap.Bootstrapper) {
	b.UseGlobal(
		func(ctx iris.Context) {
			// 当出现panic时是会以状态码500重复调用导致cros异常
			referer := getReferer(ctx)
			if ctx.GetStatusCode() < iris.StatusInternalServerError && referer != "" {
				// 设置允许跨域访问
				ctx.Header("Access-Control-Allow-Origin", referer)
				ctx.Header("Access-Control-Allow-Credentials", "true")
				ctx.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
				ctx.Header("Access-Control-Allow-Headers", getAllowHeaders())
				ctx.Header("Access-Control-Expose-Headers", "*")
			}

			// 预检查 options 直接放行
			if ctx.Method() == "OPTIONS" {
				accesslog.Skip(ctx)
				ctx.StopWithStatus(iris.StatusOK)
				return
			}

			ctx.Next()
		})
}

func getReferer(ctx iris.Context) (referer string) {
	cros := config.Config.GetCrosAllowOrigin()
	if is, _ := array.InArray("*", cros); is {
		return "*"
	}
	if referer = ctx.GetHeader("Referer"); referer == "" {
		if referer = ctx.GetHeader("Origin"); referer == "null" {
			return ""
		}
	}
	start := strings.Index(referer, ":") + 3
	end := strings.Index(referer[start:], "/")
	if end > -1 {
		referer = referer[:start+end]
	}
	if is, _ := array.InArray(referer, cros); is {
		return referer
	}
	return ""
}

func getAllowHeaders() string {
	allowHeader := "Authorization, Content-Type, Depth,User-Agent, " +
		"X-File-Size, X-Requested-With, X-Requested-By, If-Modified-Since, " +
		"X-File-Name, X-File-Type, Cache-Control, Origin, token, " +
		"x-requested-with, ts, nonce, sign"
	if s := config.Config.GetCrosAllowHeaders(); s != "" {
		allowHeader += ", " + s
	}
	return allowHeader
}
