package bootstrap

import (
	"bufio"
	"io"
	"log"
	"net"
	"os"
	"regexp"
	"time"

	"golibs/cerrors"
	"golibs/config"
	"golibs/global/consts"
	"golibs/global/variable"
	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/core/router"
	"github.com/kataras/iris/v12/hero"
	"github.com/kataras/iris/v12/middleware/accesslog"
	"github.com/kataras/iris/v12/middleware/pprof"
	recover2 "github.com/kataras/iris/v12/middleware/recover"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
)

var c = config.Config

type Configurator func(*Bootstrapper)

type Bootstrapper struct {
	*iris.Application
	AppName      string
	AppOwner     string
	AppSpawnDate time.Time
	//Sessions     *sessions.Sessions
}

// New returns a new Bootstrapper.
func New(appName, appOwner string, cfgs ...Configurator) *Bootstrapper {
	b := &Bootstrapper{
		AppName:      appName,
		AppOwner:     appOwner,
		AppSpawnDate: time.Now(),
		Application:  iris.New(),
	}

	for _, cfg := range cfgs {
		cfg(b)
	}

	return b
}

// SetupPprof 设置性能监控
func (b *Bootstrapper) SetupPprof() {
	if c.GetActive() == consts.EnvProduct {
		return
	}
	p := pprof.New()
	b.Any("/debug/pprof", p)
	b.Any("/debug/pprof/{action:path}", p)
}

// SetupLogging 设置请求日志
func (b *Bootstrapper) SetupLogging() {
	pathToAccessLog := c.GetLogFile()
	var w *rotatelogs.RotateLogs
	var err error

	if pathToAccessLog != "" {
		log.Print(pathToAccessLog)
		w, err = rotatelogs.New(
			pathToAccessLog,
			rotatelogs.WithMaxAge(24*time.Hour),
			rotatelogs.WithRotationTime(time.Hour))
		if err != nil {
			panic(err)
		}
	}

	// Initialize a new access log middleware.
	var ioWriter io.Writer
	// 生产只输出到文件
	if pathToAccessLog == "" {
		ioWriter = io.MultiWriter(os.Stdout)
	} else if c.GetActive() == consts.EnvProduct {
		ioWriter = bufio.NewWriter(w)
	} else {
		// 非生产也输出到控制台
		ioWriter = io.MultiWriter(w, os.Stdout)
	}
	ac := accesslog.New(ioWriter)
	defer ac.Close()
	ac.ResponseBody = true
	// 若为true会导致iris mvc自动转int类型为数字科学计数法
	ac.BodyMinify = false

	ac.AddFields(func(ctx iris.Context, fields *accesslog.Fields) {
		uuid := ctx.Values().Get("xRequestUuid")
		fields.Set("UUID", uuid)
	})
	//ac.Delim = ' ' // change the separator from '|' to space.
	//ac.SetFormatter(&accesslog.Template{
	//	Text: "{{.Now.Format .TimeFormat}}|{{.Latency}}|{{.Code}}|{{.Method}}|{{.Path}}|{{.IP}}|{{.RequestValuesLine}}|{{.BytesReceivedLine}}|{{.BytesSentLine}}|{{.Request}}|{{.Response}}|\n",
	//	// Default ^
	//})
	b.UseRouter(ac.Handler)
	b.Logger().SetOutput(ac.Writer)
	log.SetOutput(ac.Writer)
}

// SetupValidator 参数验证 go get github.com/go-playground/validator/v10
func (b *Bootstrapper) SetupValidator() {
	validate := validator.New()
	// 手机号验证
	validate.RegisterValidation("phone", func(fl validator.FieldLevel) bool {
		is, _ := regexp.Match("^1[0-9]{10}$", []byte(fl.Field().String()))
		return is
	})
	b.Validator = validate
}

// SetupViews loads the templates.
func (b *Bootstrapper) SetupViews(viewsDir string) {
	b.RegisterView(iris.HTML(viewsDir, ".html").
		Layout("shared/layout.html").
		Reload(c.GetActive() != consts.EnvProduct))
}

//func (b *Bootstrapper) SetupSessions(expires time.Duration, cookieHashKey, cookieBlockKey []byte) {
//	b.Sessions = sessions.New(sessions.Config{
//		Cookie:   "SECRET_SESS_COOKIE_" + strings.ToUpper(b.AppName),
//		Expires:  expires,
//		Encoding: securecookie.New(cookieHashKey, cookieBlockKey),
//	})
//	b.Sessions.UseDatabase(redis.New(redis.Config{
//		Addr:     c.GetRedis().GetHost() + ":" + strconv.Itoa(c.GetRedis().GetPort()),
//		Database: strconv.Itoa(c.GetRedis().GetDb()),
//	}))
//}

// SetupErrorHandlers `(context.StatusCodeNotSuccessful`,  which defaults to >=400 (but you can change it).
func (b *Bootstrapper) SetupErrorHandlers() {
	// 处理控制器返回的 error
	b.APIBuilder.ConfigureContainer().Container.GetErrorHandler = func(*context.Context) hero.ErrorHandler {
		return hero.ErrorHandlerFunc(func(ctx *context.Context, err error) {
			if err != hero.ErrStopExecution {
				if status := ctx.GetStatusCode(); status == 0 || !context.StatusCodeNotSuccessful(status) {
					ctx.StatusCode(hero.DefaultErrStatusCode)
				}
				if isOutJson(ctx) {
					ctx.ContentType(context.ContentJSONHeaderValue)
				}
				_, _ = ctx.WriteString(err.Error())
			}

			ctx.StopExecution()
		})
	}

	// 处理控制器里的 panic等错误
	b.OnAnyErrorCode(func(ctx iris.Context) {
		res := cerrors.Errors{}
		res.SetCode(ctx.GetStatusCode())
		res.SetMessage(iris.StatusText(ctx.GetStatusCode()))

		if isOutJson(ctx) || !c.GetIsSetupView() {
			ctx.JSON(res)
			return
		}

		err := iris.Map{
			"code":    res.GetCode(),
			"message": res.GetMessage(),
		}
		ctx.ViewData("Err", err)
		ctx.ViewData("Title", "Error")
		ctx.View("shared/error.html")
	})
}

func isOutJson(ctx iris.Context) bool {
	return ctx.IsAjax() ||
		ctx.URLParamExists("json") ||
		ctx.GetHeader("Accept") == context.ContentJSONHeaderValue
}

var (
	// StaticAssets is the root directory for public assets like images, css, js.
	StaticAssets = variable.BasePath + "/web/public/"
	// Favicon is the relative 9to the "StaticAssets") favicon path for our app.
	//Favicon = "favicon.ico"
)

// Configure accepts configurations and runs them inside the Bootstraper's context.
func (b *Bootstrapper) Configure(cs ...Configurator) {
	for _, c := range cs {
		c(b)
	}
}

func (b *Bootstrapper) Bootstrap() *Bootstrapper {

	// 设置性能监控
	b.SetupPprof()

	// 设置日志
	b.SetupLogging()

	// 设置参数验证器
	b.SetupValidator()

	if c.GetIsSetupView() {
		// 设置golang views目录
		b.SetupViews(variable.BasePath + "/web/views")
	}

	// 设置session
	//hashKey := securecookie.GenerateRandomKey(64)
	//blockKey := securecookie.GenerateRandomKey(32)
	//b.SetupSessions(10*time.Minute, hashKey, blockKey)

	// 设置错误捕获
	b.SetupErrorHandlers()

	if c.GetIsSetupView() {
		// static files
		//b.Favicon(StaticAssets + Favicon)
		router.DefaultDirOptions.IndexName = "/index.html"
		b.HandleDir("public", iris.Dir(StaticAssets), router.DefaultDirOptions)
	}
	b.Use(recover2.New())
	//b.Use(b.Sessions.Handler())

	return b
}

func (b *Bootstrapper) Listen(addr string, cfgs ...iris.Configurator) {
	b.Run(iris.Addr(addr), cfgs...)
}

func (b *Bootstrapper) ListenSock(l net.Listener, cfgs ...iris.Configurator) {
	b.Run(iris.Listener(l), cfgs...)
}
