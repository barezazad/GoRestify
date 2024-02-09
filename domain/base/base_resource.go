package base

import "GoRestify/pkg/pkg_types"

// list of resources for base domain
const (
	Domain string = "base"

	ActivityRead pkg_types.Resource = "activity:read"
	SettingWrite pkg_types.Resource = "setting:write"

	CityWrite pkg_types.Resource = "city:write"
	CityRead  pkg_types.Resource = "city:read"

	RegionWrite pkg_types.Resource = "region:write"
	RegionRead  pkg_types.Resource = "region:read"

	AccountWrite pkg_types.Resource = "account:write"
	AccountRead  pkg_types.Resource = "account:read"

	RoleWrite pkg_types.Resource = "role:write"
	RoleRead  pkg_types.Resource = "role:read"

	UserWrite pkg_types.Resource = "user:write"
	UserRead  pkg_types.Resource = "user:read"
)
