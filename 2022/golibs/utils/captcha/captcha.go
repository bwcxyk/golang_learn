package captcha

import "github.com/mojocn/base64Captcha"

/*
type configjsonbody struct {
	id            string
	captchatype   string
	verifyvalue   string
	driveraudio   *base64captcha.driveraudio
	driverstring  *base64captcha.driverstring
	driverchinese *base64captcha.driverchinese
	drivermath    *base64captcha.drivermath
	driverdigit   *base64captcha.driverdigit
}
*/

//var driver base64Captcha.Driver
var store = rdsMemStore
var driver = base64Captcha.DefaultDriverDigit

var err error

func Generate() (id, b64s string) {
	c := base64Captcha.NewCaptcha(base64Captcha.DefaultDriverDigit, store)
	id, b64s, err = c.Generate()
	if err != nil {
		panic(err)
	}
	return
}

func Verify(id, answer string) bool {
	c := base64Captcha.NewCaptcha(base64Captcha.DefaultDriverDigit, store)
	return c.Verify(id, answer, true)
}
