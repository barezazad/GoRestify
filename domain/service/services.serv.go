package service

import (
	"GoRestify/domain/acc/acc_repo"
	"GoRestify/domain/base/base_repo"
	"GoRestify/internal/core"
)

// please insert based on alphabetic sort
var (

	// base domain
	BaseCityService    BaseCityServ
	BaseRegionService  BaseRegionServ
	BaseAccountService BaseAccountServ
	BaseRoleService    BaseRoleServ
	BaseUserService    BaseUserServ
	BaseAuthService    BaseAuthServ

	// accounting domain
	AccCurrencyService      AccCurrencyServ
	AccSlotService          AccSlotServ
	AccTransactionService   AccTransactionServ
	AccAccountCreditService AccAccountCreditServ
)

// InitAllServices initiate all service
func InitAllServices(engine *core.Engine) {

	// base domain
	BaseCityService = ProvideBaseCityService(base_repo.ProvideCityRepo(engine))
	BaseRegionService = ProvideBaseRegionService(base_repo.ProvideRegionRepo(engine))
	BaseAccountService = ProvideBaseAccountService(base_repo.ProvideAccountRepo(engine))
	BaseRoleService = ProvideBaseRoleService(base_repo.ProvideRoleRepo(engine))
	BaseUserService = ProvideBaseUserService(base_repo.ProvideUserRepo(engine))
	BaseAuthService = ProvideBaseAuthService(engine)

	// accounting domain
	AccCurrencyService = ProvideAccCurrencyService(acc_repo.ProvideCurrencyRepo(engine))
	AccSlotService = ProvideAccSlotService(acc_repo.ProvideSlotRepo(engine))
	AccTransactionService = ProvideAccTransactionService(acc_repo.ProvideTransactionRepo(engine))
	AccAccountCreditService = ProvideAccAccountCreditService(acc_repo.ProvideAccountCreditRepo(engine))
}
