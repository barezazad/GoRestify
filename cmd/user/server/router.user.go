package server

import (
	"GoRestify/cmd/wire"
	"GoRestify/internal/core"
	"GoRestify/internal/core/enum/domain_app"

	"GoRestify/pkg/middleware"

	"github.com/gin-gonic/gin"
)

// Route trigger router and api methods
func Route(rg gin.RouterGroup, engine *core.Engine) {

	// Base Domain
	baseCityAPI := wire.InitBaseCityAPI(engine)

	// set X-Domain domain name, to filter APIs base on domain
	rg.Use(middleware.SetDomainInHeader(domain_app.User))

	rg.Use(middleware.JwtAuthGuard())
	{
		rg.GET("/cities", baseCityAPI.GetAll)
	}
}
