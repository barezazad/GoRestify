package setting

import (
	"fmt"
	"strings"
	"time"

	"GoRestify/pkg/models"
	"GoRestify/pkg/param"
	"GoRestify/pkg/pkg_config"
	"GoRestify/pkg/pkg_consts"
	"GoRestify/pkg/pkg_log"
	"GoRestify/pkg/pkg_terms"
	"GoRestify/pkg/pkg_types"
)

// Envs . setting envs map
var Envs map[pkg_types.Setting]pkg_types.SettingMap

// LoadSettingTick to tick like cronjob
func LoadSettingTick(domain pkg_types.Enum) {
	for {
		time.Sleep(pkg_consts.CheckToUpdateSetting * time.Second)
		go LoadSetting(domain)
	}
}

// LoadSetting read settings from database and assign them
func LoadSetting(domain pkg_types.Enum) {

	result := pkg_config.Config.Redis.Get(pkg_terms.Settings)
	if strings.Contains(result, string(domain)) && Envs != nil {
		return
	}

	params := param.New()
	params.Limit = 10000

	var settings []models.Setting
	var err error
	if settings, err = List(params); err != nil {
		pkg_log.Fatal(err, "load settings failed")
	}

	var settingEnvs = make(map[pkg_types.Setting]pkg_types.SettingMap)

	for _, v := range settings {
		settingVal := pkg_types.SettingMap{
			Value: v.Value,
		}
		settingEnvs[v.Property] = settingVal
	}

	// if result don't contain domain, append domain to value
	redisValue := result
	if !strings.Contains(result, string(domain)) {
		redisValue = fmt.Sprintf("%v-%v", result, domain)
	}

	if err = pkg_config.Config.Redis.Set(pkg_terms.Settings, redisValue); err != nil {
		pkg_log.Fatal(err, "set setting in cache failed")
	}

	Envs = settingEnvs
}
