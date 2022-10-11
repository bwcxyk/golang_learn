package core

import (
	"golibs/global/consts"
	"os"
)

func GetDevBasePath() string {
	env, _ := os.LookupEnv(consts.DevBasePathKey)
	return env
}
