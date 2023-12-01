package server

import (
	"GoRestify/cmd/wire"
	"GoRestify/domain/base"
	"GoRestify/internal/core"
	"GoRestify/internal/core/check_access"
	"GoRestify/internal/core/enum/domain_app"

	"GoRestify/pkg/middleware"

	"github.com/gin-gonic/gin"
)

// Route trigger router and api methods
func Route(rg gin.RouterGroup, engine *core.Engine) {

	baseAuthAPI := wire.InitBaseAuthAPI(engine)
	baseRoleAPI := wire.InitBaseRoleAPI(engine)
	baseUserAPI := wire.InitBaseUserAPI(engine)
	baseRegionAPI := wire.InitBaseRegionAPI(engine)
	baseCityAPI := wire.InitBaseCityAPI(engine)
	PkgAPI := wire.InitBasePkgAPI(engine)

	// set X-Domain domain name, to filter APIs base on domain
	rg.Use(middleware.SetDomainInHeader(domain_app.Admin))

	rg.POST("/login", baseAuthAPI.Login)
	{
		rg.Use(middleware.JwtAuthGuard())

		// Pkg APIs
		rg.GET("/settings", check_access.IsAllow(base.SettingWrite), PkgAPI.SettingList)
		rg.PUT("/settings/:settingID", check_access.IsAllow(base.SettingWrite), PkgAPI.SettingUpdate)

		rg.GET("/activities", check_access.IsAllow(base.ActivityRead), PkgAPI.ActivitiesList)

		rg.PUT("/clear-cache/:key", check_access.IsAllow(base.SettingWrite), PkgAPI.RedisResetCacheByKey)
		rg.PUT("/clear-cache/user/:userID", check_access.IsAllow(base.SettingWrite), PkgAPI.RedisClearCacheToUser)

		// Base Domain
		rg.GET("/roles", check_access.IsAllow(base.RoleRead), baseRoleAPI.List)
		rg.GET("/resources", check_access.IsAllow(base.RoleRead), baseRoleAPI.GetResources)
		rg.GET("/all/roles", check_access.IsAllow(base.RoleRead), baseRoleAPI.GetAll)
		rg.GET("/roles/:roleID", check_access.IsAllow(base.RoleRead), baseRoleAPI.FindByID)
		rg.POST("/roles", check_access.IsAllow(base.RoleWrite), baseRoleAPI.Create)
		rg.PUT("/roles/:roleID", check_access.IsAllow(base.RoleWrite), baseRoleAPI.Update)
		rg.DELETE("/roles/:roleID", check_access.IsAllow(base.RoleWrite), baseRoleAPI.Delete)

		rg.GET("/users", check_access.IsAllow(base.UserRead), baseUserAPI.List)
		rg.GET("/all/users", check_access.IsAllow(base.UserRead), baseUserAPI.GetAll)
		rg.GET("/users/:userID", check_access.IsAllow(base.UserRead), baseUserAPI.FindByID)
		rg.POST("/users", check_access.IsAllow(base.UserWrite), baseUserAPI.Create)
		rg.PUT("/users/:userID", check_access.IsAllow(base.UserWrite), baseUserAPI.Update)
		rg.DELETE("/users/:userID", check_access.IsAllow(base.UserWrite), baseUserAPI.Delete)

		rg.GET("/cities", check_access.IsAllow(base.CityRead), baseCityAPI.List)
		rg.GET("/all/cities", check_access.IsAllow(base.CityRead), baseCityAPI.GetAll)
		rg.GET("/cities/:cityID", check_access.IsAllow(base.CityRead), baseCityAPI.FindByID)
		rg.POST("/cities", check_access.IsAllow(base.CityWrite), baseCityAPI.Create)
		rg.PUT("/cities/:cityID", check_access.IsAllow(base.CityWrite), baseCityAPI.Update)
		rg.DELETE("/cities/:cityID", check_access.IsAllow(base.CityWrite), baseCityAPI.Delete)
		rg.GET("/excel/cities", check_access.IsAllow(base.CityRead), baseCityAPI.Excel)

		rg.GET("/regions", check_access.IsAllow(base.RegionRead), baseRegionAPI.List)
		rg.GET("/all/regions", check_access.IsAllow(base.RegionRead), baseRegionAPI.GetAll)
		rg.GET("/regions/:regionID", check_access.IsAllow(base.RegionRead), baseRegionAPI.FindByID)
		rg.POST("/regions", check_access.IsAllow(base.RegionWrite), baseRegionAPI.Create)
		rg.PUT("/regions/:regionID", check_access.IsAllow(base.RegionWrite), baseRegionAPI.Update)
		rg.DELETE("/regions/:regionID", check_access.IsAllow(base.RegionWrite), baseRegionAPI.Delete)

	}
}
