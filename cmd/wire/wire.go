//go:build wireinject
// +build wireinject

package wire

import (
	"GoRestify/domain/acc/acc_api"
	"GoRestify/domain/acc/acc_repo"
	"GoRestify/domain/base/base_api"
	"GoRestify/domain/base/base_repo"
	"GoRestify/domain/service"
	"GoRestify/internal/core"

	"github.com/google/wire"
)

// Base Domain
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

func InitBaseAccountAPI(e *core.Engine) base_api.AccountAPI {
	wire.Build(base_repo.ProvideAccountRepo, service.ProvideBaseAccountService,
		base_api.ProvideAccountAPI)
	return base_api.AccountAPI{}
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

// Acc Domain
func InitAccTransactionAPI(e *core.Engine) acc_api.TransactionAPI {
	wire.Build(acc_repo.ProvideTransactionRepo, service.ProvideAccTransactionService,
		acc_api.ProvideTransactionAPI)
	return acc_api.TransactionAPI{}
}

func InitAccSlotAPI(e *core.Engine) acc_api.SlotAPI {
	wire.Build(acc_repo.ProvideSlotRepo, service.ProvideAccSlotService,
		acc_api.ProvideSlotAPI)
	return acc_api.SlotAPI{}
}

func InitAccCurrencyAPI(e *core.Engine) acc_api.CurrencyAPI {
	wire.Build(acc_repo.ProvideCurrencyRepo, service.ProvideAccCurrencyService,
		acc_api.ProvideCurrencyAPI)
	return acc_api.CurrencyAPI{}
}

func InitAccAccountCreditAPI(e *core.Engine) acc_api.AccountCreditAPI {
	wire.Build(acc_repo.ProvideAccountCreditRepo, service.ProvideAccAccountCreditService,
		acc_api.ProvideAccountCreditAPI)
	return acc_api.AccountCreditAPI{}
}

func InitBaseDocumentAPI(e *core.Engine) base_api.DocumentAPI {
	wire.Build(base_repo.ProvideDocumentRepo, service.ProvideBaseDocumentService,
		base_api.ProvideDocumentAPI)
	return base_api.DocumentAPI{}
}
