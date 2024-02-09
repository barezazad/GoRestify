package table

import (
	"GoRestify/domain/base/base_model"
	"GoRestify/domain/base/enum/account_status"
	"GoRestify/domain/base/enum/account_type"
	"GoRestify/domain/service"
	"GoRestify/internal/core"

	"GoRestify/pkg/pkg_log"
	"GoRestify/pkg/tx"
)

// InsertAccounts for add required accounts
func InsertAccounts(engine *core.Engine) {
	accounts := []base_model.Account{
		{
			ID:       1,
			FullName: "admin",
			Username: "admin",
			Password: "admin123Aa",
			Email:    "barezazad100@gmail.com",
			Phone:    "9647705549911",
			Type:     account_type.User,
			Status:   account_status.Active,
			RoleID:   1,
		},
	}

	for _, v := range accounts {
		if _, err := service.BaseAccountService.FindByID(tx.Tx{}, v.ID); err == nil {
			if _, _, err = service.BaseAccountService.Save(tx.Tx{}, v); err != nil {
				pkg_log.Fatal("error in saving accounts", err)
			}
		} else {
			if _, err = service.BaseAccountService.Create(tx.Tx{}, v); err != nil {
				pkg_log.Fatal("error in creating accounts", err)
			}
		}
	}

}
