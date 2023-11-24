package core

import (
	"GoRestify/pkg/pkg_types"
	"GoRestify/pkg/utils"
)

// list of core environment keys
const (
	Port    pkg_types.Envkey = "CORE_PORT"
	GinMode pkg_types.Envkey = "GIN_MODE"

	AutoMigrate pkg_types.Envkey = "CORE_AUTO_MIGRATE"

	DatabaseDataDSN     pkg_types.Envkey = "CORE_DATABASE_DATA_DSN"
	DatabaseActivityDSN pkg_types.Envkey = "CORE_DATABASE_ACTIVITY_DSN"

	RedisCacheAPI pkg_types.Envkey = "CORE_REDIS_CACHE_API"

	JWTSecretKey pkg_types.Envkey = "CORE_JWT_SECRET_KEY"
)

// ListAdminEnv list of env for admin
var ListAdminEnv = []pkg_types.Envkey{
	Port,
	GinMode,

	AutoMigrate,

	DatabaseDataDSN,

	DatabaseActivityDSN,

	RedisCacheAPI,
}

// ListUserEnv list of env for user
var ListUserEnv = []pkg_types.Envkey{
	Port,
	GinMode,

	DatabaseDataDSN,

	RedisCacheAPI,
}

// LoadEnvs load environment from env file
func LoadEnvs(envList []pkg_types.Envkey) *Engine {
	var engine Engine
	engine.Envs = utils.SetENVs(envList)
	return &engine
}
