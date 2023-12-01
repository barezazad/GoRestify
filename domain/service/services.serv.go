package service

import (
	"GoRestify/domain/base/base_repo"
	"GoRestify/internal/core"
)

// please insert based on alphabetic sort
var (

	// base domain
	BaseCityService   BaseCityServ
	BaseRegionService BaseRegionServ
	BaseRoleService   BaseRoleServ
	BaseUserService   BaseUserServ
	BaseAuthService   BaseAuthServ
)

// InitAllServices initiate all service
func InitAllServices(engine *core.Engine) {
	BaseCityService = ProvideBaseCityService(base_repo.ProvideCityRepo(engine))
	BaseRegionService = ProvideBaseRegionService(base_repo.ProvideRegionRepo(engine))
	BaseRoleService = ProvideBaseRoleService(base_repo.ProvideRoleRepo(engine))
	BaseUserService = ProvideBaseUserService(base_repo.ProvideUserRepo(engine))
	BaseAuthService = ProvideBaseAuthService(engine)
}
