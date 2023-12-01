package table

import (
	"GoRestify/domain/base/base_model"
	"GoRestify/domain/base/enum/user_status"
	"GoRestify/domain/service"
	"GoRestify/internal/core"

	"GoRestify/pkg/pkg_log"
	"GoRestify/pkg/tx"
)

// InsertUsers for add required users
func InsertUsers(engine *core.Engine) {
	users := []base_model.User{
		{
			ID:       1,
			RoleID:   1,
			FullName: "admin",
			Username: "admin",
			Password: "admin123Aa",
			Email:    "barezazad100@gmail.com",
			Phone:    "9647705549911",
			Status:   user_status.Active,
		},
	}

	for _, v := range users {
		if _, err := service.BaseUserService.FindByID(tx.Tx{}, v.ID); err == nil {
			if _, _, err = service.BaseUserService.Save(tx.Tx{}, v); err != nil {
				pkg_log.Fatal("error in saving users", err)
			}
		} else {
			if _, err = service.BaseUserService.Create(tx.Tx{}, v); err != nil {
				pkg_log.Fatal("error in creating users", err)
			}
		}
	}

}
