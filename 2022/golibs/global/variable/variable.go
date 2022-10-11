package variable

import (
	"flag"
	"golibs/core"
	"golibs/global/consts"
	"os"
	"strings"
)

var (
	// BasePath 项目根目录
	BasePath = ""
	// ConfPath 配置目录
	ConfPath = ""
)

var workdir = flag.String("w", "", "指定工作目录 -w /workdir")
var confdir = flag.String("c", "", "指定配置目录 -c /confdir")

func init() {
	if len(os.Args) > 1 && strings.HasPrefix(os.Args[1], "-test") {
		BasePath = core.GetDevBasePath()
		if BasePath == "" {
			panic("Environment variable [" + consts.DevBasePathKey + "] not set. Specify the project root directory.")
		}
	} else {
		flag.Parse()
		if *workdir != "" {
			BasePath = *workdir
		} else if pwd, err := os.Getwd(); err == nil {
			BasePath = pwd
		}
	}
	ConfPath = BasePath
	if *confdir != "" {
		ConfPath = *confdir
	}
}
