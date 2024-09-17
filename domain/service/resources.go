package service

import (
	"GoRestify/domain/acc"
	"GoRestify/domain/base"
	"GoRestify/pkg/pkg_types"
)

// AllResources get list of all resources
func AllResources() []pkg_types.Resource {
	resource := []pkg_types.Resource{
		base.ActivityRead,
		base.SettingWrite,

		base.CityWrite, base.CityRead,
		base.RegionWrite, base.RegionRead,

		base.RoleWrite, base.RoleRead,
		base.UserWrite, base.UserRead,
		base.AccountWrite, base.AccountRead,

		acc.CurrencyWrite, acc.CurrencyRead,
		acc.TransactionWrite, acc.TransactionRead,
	}

	return resource
}
