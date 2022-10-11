package consts

const (
	// DevBasePathKey 开发环境项目根目录，ide或系统上设置环境变量并制定目录 如：export DEVBASEPATH=/project/golibs
	DevBasePathKey string = "DEVBASEPATH"

	// EnvProduct 生产环境
	EnvProduct string = "prod"

	// EnvPre 灰度环境
	EnvPre string = "pre"

	// EnvUat 用户验证环境
	EnvUat string = "uat"

	// EnvTest 测试环境
	EnvTest string = "test"

	// EnvDev 开发环境
	EnvDev string = "dev"
)
