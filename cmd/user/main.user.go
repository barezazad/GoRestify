package main

import (
	"log"

	"GoRestify/cmd/user/server"
	"GoRestify/domain/service"
	"GoRestify/internal/core"
	"GoRestify/internal/core/action"
	"GoRestify/internal/core/enum/domain_app"
	"GoRestify/pkg"

	"GoRestify/pkg/pkg_config"
	"GoRestify/pkg/pkg_redis"
	"GoRestify/pkg/pkg_sql"
	"GoRestify/pkg/utils"
)

func main() {

	engine := core.LoadEnvs(core.ListUserEnv)

	utils.GenerateErrCode()

	confCoreEngine := pkg_config.Cnf{
		IsDebug: engine.Envs[core.GinMode],

		DbDSN:         engine.Envs[core.DatabaseDataDSN],
		ActivityDbDSN: engine.Envs[core.DatabaseActivityDSN],
		RedisAddress:  engine.Envs[core.RedisCacheAPI],

		MustBeInTypes: action.MustBeInTypes,

		SettingActive:    true,
		SettingDomainApp: domain_app.Admin,

		JWTPrivateKeyPath: engine.Envs[core.JWTPrivateKeyPath],
		JWTPublicKeyPath:  engine.Envs[core.JWTPublicKeyPath],
		JWTAccessMinutes:  engine.Envs.ToInt(core.JWTAccessMinutes),
		JWTRefreshHours:   engine.Envs.ToInt(core.JWTRefreshHours),
	}

	pkg.Init(confCoreEngine)

	// connect the database
	engine.DB = pkg_sql.MySQLConnectDB(engine.Envs[core.DatabaseDataDSN])

	// establish redis connections
	var err error
	if engine.RedisCacheAPI, err = pkg_redis.ConnectRedis(engine.Envs[core.RedisCacheAPI], pkg_config.Config.IsDebug); err != nil {
		log.Fatal("Redis server connection failed")
		return
	}

	service.InitAllServices(engine)

	// start the API
	server.Start(engine)

}
