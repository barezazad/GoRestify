package core

import (
	"GoRestify/pkg/pkg_redis"
	"GoRestify/pkg/pkg_types"

	"gorm.io/gorm"
)

// Engine to keep all database connections and
// logs configuration and environments and etc
type Engine struct {
	DB            *gorm.DB
	Envs          pkg_types.Envs
	RedisCacheAPI pkg_redis.RedisCon
}
