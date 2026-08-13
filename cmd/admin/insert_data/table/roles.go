package table

import (
	"GoRestify/domain/base/base_model"
	"GoRestify/domain/service"
	"GoRestify/internal/core"

	"GoRestify/pkg/pkg_log"
	"GoRestify/pkg/pkg_types"
	"GoRestify/pkg/tx"
)

// InsertRoles for add required roles
func InsertRoles(engine *core.Engine) {
	roles := []base_model.Role{
		{ID: 1, Name: "Admin", Resources: pkg_types.ResourceJoin(service.AllResources())},
	}

	for _, v := range roles {
		if _, err := service.BaseRoleService.FindByID(tx.Tx{}, v.ID); err == nil {
			if _, _, err = service.BaseRoleService.Save(tx.Tx{}, v); err != nil {
				pkg_log.Fatal("error in saving roles", err)
			}
		} else {
			if _, err = service.BaseRoleService.Create(tx.Tx{}, v); err != nil {
				pkg_log.Fatal("error in creating roles", err)
			}
		}
	}

}
