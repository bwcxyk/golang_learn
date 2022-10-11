package identity

import (
	"fmt"
	"golibs/framework/iris/bootstrap"
	"github.com/google/uuid"
	"github.com/kataras/golog"
	"github.com/kataras/iris/v12"
	"log"
	"strconv"
	"time"
)

func Identity(b *bootstrap.Bootstrapper) {
	b.UseGlobal(func(ctx iris.Context) {
		if ctx.GetStatusCode() < iris.StatusInternalServerError && ctx.Method() != "OPTIONS" {
			// 重置 session 未操作计时器，长时不操作 session 自动销毁
			//b.Sessions.ShiftExpiration(ctx)

			// response headers
			ctx.Header("App-Name", b.AppName)
			ctx.Header("App-Owner", b.AppOwner)
			ctx.Header("App-Since", time.Since(b.AppSpawnDate).String())

			ctx.Header("x-time", strconv.Itoa(int(time.Now().Unix())))
			uuid, _ := uuid.NewUUID()
			ctx.Header("x-request-uuid", uuid.String())
			ctx.Values().Set("xRequestUuid", uuid.String())
			pUuid := fmt.Sprintf("UUID=%v ", uuid.String())
			log.SetPrefix(pUuid)
			golog.SetPrefix(pUuid)

			// view data if ctx.View or c.Tmpl = "$page.html" will be called next.
			ctx.ViewData("AppName", b.AppName)
			ctx.ViewData("AppOwner", b.AppOwner)

			ctx.Next()
		}
	})
}
