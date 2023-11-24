package main

import (
	"log"

	"GoRestify/cmd/admin/insert_data"
	"GoRestify/cmd/admin/server"
	"GoRestify/cmd/admin/start_off"
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
	engine := core.LoadEnvs(core.ListAdminEnv)

	utils.GenerateErrCode()

	confCoreEngine := pkg_config.Cnf{
		IsDebug: engine.Envs[core.GinMode],

		DbDSN:         engine.Envs[core.DatabaseDataDSN],
		ActivityDbDSN: engine.Envs[core.DatabaseActivityDSN],
		RedisAddress:  engine.Envs[core.RedisCacheAPI],

		MustBeInTypes: action.MustBeInTypes,

		SettingActive:    true,
		SettingDomainApp: domain_app.Admin,

		JWTSecretKey: engine.Envs[core.JWTSecretKey],
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

	if engine.Envs.ToBool(core.AutoMigrate) {
		// migrate the database
		start_off.Migrate(engine)
		// insert basic data
		insert_data.Insert(engine)
	}

	// start the API
	server.Start(engine)
}
