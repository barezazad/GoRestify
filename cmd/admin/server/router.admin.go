package server

import (
	"GoRestify/internal/core"
	"GoRestify/internal/core/enum/domain_app"

	"GoRestify/pkg/middleware"

	"github.com/gin-gonic/gin"
)

// Route trigger router and api methods
func Route(rg gin.RouterGroup, engine *core.Engine) {

	// Base Domain
	// baseCityAPI := wire.InitBaseCityAPI(engine)
	// baseRegionAPI := wire.InitBaseRegionAPI(engine)
	// PkgAPI := wire.InitBasePkgAPI(engine)

	// set X-Domain domain name, to filter APIs base on domain
	rg.Use(middleware.SetDomainInHeader(domain_app.Admin))

	{
		// Pkg APIs
		// rg.GET("/settings", app.CheckLocalAuth(base.SettingWrite), PkgAPI.SettingList)
		// rg.PUT("/settings/:settingID", app.CheckLocalAuth(base.SettingWrite), PkgAPI.SettingUpdate)

		// rg.GET("/activities", app.CheckLocalAuth(base.ActivityRead), PkgAPI.ActivitiesList)

		// rg.PUT("/clear-cache/:key", app.CheckLocalAuth(base.SettingWrite), PkgAPI.RedisResetCacheByKey)
		// rg.PUT("/clear-cache/user/:userID", app.CheckLocalAuth(base.SettingWrite), PkgAPI.RedisClearCacheToUser)

		// // Base Domain
		// rg.GET("/cities", app.CheckLocalAuth(base.CityRead), baseCityAPI.List)
		// rg.GET("/all/cities", app.CheckLocalAuth(base.CityRead), baseCityAPI.GetAll)
		// rg.GET("/cities/:cityID", app.CheckLocalAuth(base.CityRead), baseCityAPI.FindByID)
		// rg.POST("/cities", app.CheckLocalAuth(base.CityWrite), baseCityAPI.Create)
		// rg.PUT("/cities/:cityID", app.CheckLocalAuth(base.CityWrite), baseCityAPI.Update)
		// rg.DELETE("/cities/:cityID", app.CheckLocalAuth(base.CityWrite), baseCityAPI.Delete)
		// rg.GET("/excel/cities", app.CheckLocalAuth(base.CityRead), baseCityAPI.Excel)

		// rg.GET("/regions", app.CheckLocalAuth(base.RegionRead), baseRegionAPI.List)
		// rg.GET("/all/regions", app.CheckLocalAuth(base.RegionRead), baseRegionAPI.GetAll)
		// rg.GET("/regions/:regionID", app.CheckLocalAuth(base.RegionRead), baseRegionAPI.FindByID)
		// rg.POST("/regions", app.CheckLocalAuth(base.RegionWrite), baseRegionAPI.Create)
		// rg.PUT("/regions/:regionID", app.CheckLocalAuth(base.RegionWrite), baseRegionAPI.Update)
		// rg.DELETE("/regions/:regionID", app.CheckLocalAuth(base.RegionWrite), baseRegionAPI.Delete)
	}
}
