package base_api

import (
	"GoRestify/domain/base"
	"GoRestify/domain/base/base_model"
	"GoRestify/domain/base/base_term"
	"GoRestify/domain/service"
	"GoRestify/internal/core"
	"fmt"
	"net/http"

	"GoRestify/pkg/pkg_err"
	"GoRestify/pkg/pkg_log"
	"GoRestify/pkg/pkg_terms"
	"GoRestify/pkg/response"

	"github.com/gin-gonic/gin"
)

// AccountAPI for injecting account service
type AccountAPI struct {
	Service service.BaseAccountServ
	Engine  *core.Engine
}

// ProvideAccountAPI .
func ProvideAccountAPI(c service.BaseAccountServ) AccountAPI {
	return AccountAPI{Service: c, Engine: c.Engine}
}

// FindByID is used for fetch an account by its id
func (a *AccountAPI) FindByID(c *gin.Context) {
	resp, params := response.NewParam(c, base_model.AccountTable)
	var err error
	var account base_model.Account
	var id uint

	if id, err = resp.GetID(c.Param("accountID"), "E1163666", base_term.Account); err != nil {
		return
	}

	if account, err = a.Service.FindByID(params.Tx, id); err != nil {
		resp.Error(err).JSON()
		return
	}

	account.Password = ""

	resp.Record(base.ViewAccount)
	resp.Status(http.StatusOK).
		Message(pkg_terms.VInfo, base_term.Account).
		JSON(account)
}

// GetAll list of accounts
func (a *AccountAPI) GetAll(c *gin.Context) {
	resp, params := response.NewParam(c, base_model.AccountTable)
	var accounts []base_model.Account
	var err error

	if accounts, err = a.Service.GetAll(params); err != nil {
		err = pkg_err.Take(err, "E1180593").Message(pkg_err.SomethingWentWrong).Build()
		resp.Error(err).JSON()
		return
	}

	resp.Record(base.ListAccount)
	resp.Status(http.StatusOK).
		Message(pkg_terms.ListOfV, base_term.Accounts).
		JSON(accounts)
}

// List of accounts
func (a *AccountAPI) List(c *gin.Context) {
	resp, params := response.NewParam(c, base_model.AccountTable)

	data := make(map[string]interface{})
	var err error

	if data["list"], data["count"], err = a.Service.List(params); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(base.ListAccount)
	resp.Status(http.StatusOK).
		Message(pkg_terms.ListOfV, base_term.Accounts).
		JSON(data)
}

// Create account
func (a *AccountAPI) Create(c *gin.Context) {
	resp, params := response.NewParam(c, base_model.AccountTable)
	var account, createdAccount base_model.Account
	var err error

	if err = resp.Bind(&account, "E1150888", base_term.Account); err != nil {
		return
	}

	// start tx
	params.Tx.DB = a.Engine.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			pkg_log.LogError(fmt.Errorf("panic happened in tx mode for %v",
				"account"), "rollback recover create account")
			err = pkg_err.New(pkg_err.SomethingWentWrong, "E1145492").
				Message(pkg_err.SomethingWentWrong).Build()
			// rollback tx
			params.Tx.DB.Rollback()
			return
		}
	}()

	if createdAccount, err = a.Service.Create(params.Tx, account); err != nil {
		resp.Error(err).JSON()
		// rollback tx
		params.Tx.DB.Rollback()
		return
	}

	// commit tx
	params.Tx.DB.Commit()

	resp.Record(base.CreateAccount, account, createdAccount)
	resp.Status(http.StatusOK).
		Message(pkg_terms.VCreatedSuccessfully, base_term.Account).
		JSON(createdAccount)
}

// Update account
func (a *AccountAPI) Update(c *gin.Context) {
	resp, params := response.NewParam(c, base_model.AccountTable)
	var err error
	var account, accountBefore, accountUpdated base_model.Account

	if err = resp.Bind(&account, "E1135664", base_term.Account); err != nil {
		return
	}

	if account.ID, err = resp.GetID(c.Param("accountID"), "E1162016", base_term.Account); err != nil {
		return
	}

	// start tx
	params.Tx.DB = a.Engine.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			pkg_log.LogError(fmt.Errorf("panic happened in tx mode for %v",
				"account"), "rollback recover create account")
			err = pkg_err.New(pkg_err.SomethingWentWrong, "E1138430").
				Message(pkg_err.SomethingWentWrong).Build()
			// rollback tx
			params.Tx.DB.Rollback()
			return
		}
	}()

	if accountUpdated, accountBefore, err = a.Service.Save(params.Tx, account); err != nil {
		resp.Error(err).JSON()
		// rollback tx
		params.Tx.DB.Rollback()
		return
	}

	// commit tx
	params.Tx.DB.Commit()

	resp.Record(base.UpdateAccount, accountBefore, account)
	resp.Status(http.StatusOK).
		Message(pkg_terms.VUpdatedSuccessfully, base_term.Account).
		JSON(accountUpdated)
}

// Delete account
func (a *AccountAPI) Delete(c *gin.Context) {
	resp, params := response.NewParam(c, base_model.AccountTable)
	var err error
	var account base_model.Account
	var id uint

	if id, err = resp.GetID(c.Param("accountID"), "E1144110", base_term.Account); err != nil {
		return
	}

	// start tx
	params.Tx.DB = a.Engine.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			pkg_log.LogError(fmt.Errorf("panic happened in tx mode for %v",
				"account"), "rollback recover create account")
			err = pkg_err.New(pkg_err.SomethingWentWrong, "E1199821").
				Message(pkg_err.SomethingWentWrong).Build()
			// rollback tx
			params.Tx.DB.Rollback()
			return
		}
	}()

	if account, err = a.Service.Delete(params.Tx, id); err != nil {
		resp.Error(err).JSON()
		// rollback tx
		params.Tx.DB.Rollback()
		return
	}

	// commit tx
	params.Tx.DB.Commit()

	resp.Record(base.DeleteAccount, account)
	resp.Status(http.StatusOK).
		Message(pkg_terms.VDeletedSuccessfully, base_term.Account).
		JSON()
}
