package service

import (
	"GoRestify/domain/base"
	"GoRestify/pkg/pkg_types"
)

// AllResources get list of all resources
func AllResources() []pkg_types.Resource {
	resource := []pkg_types.Resource{
		base.ActivityRead,
		base.SettingWrite,
		base.CityWrite,
		base.CityRead,
		base.RegionWrite,
		base.RegionRead,
		base.RoleWrite,
		base.RoleRead,
		base.UserWrite,
		base.UserRead,
	}

	return resource
}
