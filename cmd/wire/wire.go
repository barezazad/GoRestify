//go:build wireinject
// +build wireinject

package wire

import (
	"GoRestify/domain/base/base_api"
	"GoRestify/domain/base/base_repo"
	"GoRestify/domain/service"
	"GoRestify/internal/core"

	"github.com/google/wire"
)

func InitBaseCityAPI(e *core.Engine) base_api.CityAPI {
	wire.Build(base_repo.ProvideCityRepo, service.ProvideBaseCityService,
		base_api.ProvideCityAPI)
	return base_api.CityAPI{}
}

func InitBaseRegionAPI(e *core.Engine) base_api.RegionAPI {
	wire.Build(base_repo.ProvideRegionRepo, service.ProvideBaseRegionService,
		base_api.ProvideRegionAPI)
	return base_api.RegionAPI{}
}

func InitBaseRoleAPI(e *core.Engine) base_api.RoleAPI {
	wire.Build(base_repo.ProvideRoleRepo, service.ProvideBaseRoleService,
		base_api.ProvideRoleAPI)
	return base_api.RoleAPI{}
}

func InitBaseUserAPI(e *core.Engine) base_api.UserAPI {
	wire.Build(base_repo.ProvideUserRepo, service.ProvideBaseUserService,
		base_api.ProvideUserAPI)
	return base_api.UserAPI{}
}

func InitBaseAuthAPI(e *core.Engine) base_api.AuthAPI {
	wire.Build(service.ProvideBaseAuthService,
		base_api.ProvideAuthAPI)
	return base_api.AuthAPI{}
}

func InitBasePkgAPI(e *core.Engine) base_api.PkgAPI {
	wire.Build(base_api.ProvidePkgAPI)
	return base_api.PkgAPI{}
}
