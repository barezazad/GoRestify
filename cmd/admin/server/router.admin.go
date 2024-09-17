package server

import (
	"GoRestify/cmd/wire"
	"GoRestify/domain/acc"
	"GoRestify/domain/base"
	"GoRestify/internal/core"
	"GoRestify/internal/core/check_access"
	"GoRestify/internal/core/enum/domain_app"

	"GoRestify/pkg/middleware"

	"github.com/gin-gonic/gin"
)

// Route trigger router and api methods
func Route(rg gin.RouterGroup, engine *core.Engine) {

	// base domain
	baseAuthAPI := wire.InitBaseAuthAPI(engine)
	baseRoleAPI := wire.InitBaseRoleAPI(engine)
	baseAccountAPI := wire.InitBaseAccountAPI(engine)
	baseRegionAPI := wire.InitBaseRegionAPI(engine)
	baseCityAPI := wire.InitBaseCityAPI(engine)
	PkgAPI := wire.InitBasePkgAPI(engine)

	// acc domain
	accTransactionAPI := wire.InitAccTransactionAPI(engine)
	accCurrencyAPI := wire.InitAccCurrencyAPI(engine)

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

		rg.GET("/roles", check_access.IsAllow(base.RoleRead), baseRoleAPI.List)
		rg.GET("/resources", check_access.IsAllow(base.RoleRead), baseRoleAPI.GetResources)
		rg.GET("/all/roles", check_access.IsAllow(base.RoleRead), baseRoleAPI.GetAll)
		rg.GET("/roles/:roleID", check_access.IsAllow(base.RoleRead), baseRoleAPI.FindByID)
		rg.POST("/roles", check_access.IsAllow(base.RoleWrite), baseRoleAPI.Create)
		rg.PUT("/roles/:roleID", check_access.IsAllow(base.RoleWrite), baseRoleAPI.Update)
		rg.DELETE("/roles/:roleID", check_access.IsAllow(base.RoleWrite), baseRoleAPI.Delete)

		rg.GET("/accounts", check_access.IsAllow(base.AccountRead), baseAccountAPI.List)
		rg.GET("/all/accounts", check_access.IsAllow(base.AccountRead), baseAccountAPI.GetAll)
		rg.GET("/accounts/:accountID", check_access.IsAllow(base.AccountRead), baseAccountAPI.FindByID)
		rg.POST("/accounts", check_access.IsAllow(base.AccountWrite), baseAccountAPI.Create)
		rg.PUT("/accounts/:accountID", check_access.IsAllow(base.AccountWrite), baseAccountAPI.Update)
		rg.DELETE("/accounts/:accountID", check_access.IsAllow(base.AccountWrite), baseAccountAPI.Delete)

		// acc domain
		rg.GET("/currencies", check_access.IsAllow(acc.CurrencyRead), accCurrencyAPI.List)
		rg.GET("/all/currencies", check_access.IsAllow(acc.CurrencyRead), accCurrencyAPI.GetAll)
		rg.GET("/currencies/:currencyID", check_access.IsAllow(acc.CurrencyRead), accCurrencyAPI.FindByID)
		rg.POST("/currencies", check_access.IsAllow(acc.CurrencyWrite), accCurrencyAPI.Create)
		rg.PUT("/currencies/:currencyID", check_access.IsAllow(acc.CurrencyWrite), accCurrencyAPI.Update)
		rg.DELETE("/currencies/:currencyID", check_access.IsAllow(acc.CurrencyWrite), accCurrencyAPI.Delete)

		rg.GET("/transactions", check_access.IsAllow(acc.TransactionRead), accTransactionAPI.List)
		rg.GET("/all/transactions", check_access.IsAllow(acc.TransactionRead), accTransactionAPI.GetAll)
		rg.GET("/transactions/:transactionID", check_access.IsAllow(acc.TransactionRead), accTransactionAPI.FindByID)
	}
}
