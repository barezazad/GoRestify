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

func InitBasePkgAPI(e *core.Engine) base_api.PkgAPI {
	wire.Build(base_api.ProvidePkgAPI)
	return base_api.PkgAPI{}
}
