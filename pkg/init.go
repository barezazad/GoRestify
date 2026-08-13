package pkg

import (
	"GoRestify/pkg/pkg_config"
	"GoRestify/pkg/pkg_jwt"
	"GoRestify/pkg/response"
	"GoRestify/pkg/setting"
	"log"
)

// Init this function init core pkg_config
func Init(cnf pkg_config.Cnf) {

	if _, err := pkg_config.InitConfig(cnf); err != nil {
		log.Fatalf("Error in Core Config: %v", err)
	}

	if err := pkg_jwt.Init(cnf.JWTPrivateKeyPath, cnf.JWTPublicKeyPath, cnf.JWTAccessMinutes, cnf.JWTRefreshHours); err != nil {
		log.Fatalf("Error in JWT Config: %v", err)
	}

	// load setting env for config, in case we activate setting
	if cnf.SettingActive {
		setting.LoadSetting(cnf.SettingDomainApp)
		go setting.LoadSettingTick(cnf.SettingDomainApp)
	}

	// watch activities
	go response.ActivityWatcher()

	return
}
