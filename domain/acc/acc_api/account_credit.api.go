package acc_api

import (
	"GoRestify/domain/acc"
	"GoRestify/domain/acc/acc_model"
	"GoRestify/domain/acc/acc_term"
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

// AccountCreditAPI for injecting accountCredit service
type AccountCreditAPI struct {
	Service service.AccAccountCreditServ
	Engine  *core.Engine
}

// ProvideAccountCreditAPI .
func ProvideAccountCreditAPI(c service.AccAccountCreditServ) AccountCreditAPI {
	return AccountCreditAPI{Service: c, Engine: c.Engine}
}

// FindByID is used for fetch an accountCredit by its id
func (a *AccountCreditAPI) FindByID(c *gin.Context) {
	resp, params := response.NewParam(c, acc_model.AccountCreditTable)
	var err error
	var accountCredit acc_model.AccountCredit
	var id uint

	if id, err = resp.GetID(c.Param("accountCreditID"), "E1163666", acc_term.AccountCredit); err != nil {
		return
	}

	if accountCredit, err = a.Service.FindByID(params.Tx, id); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(acc.ViewAccountCredit)
	resp.Status(http.StatusOK).
		Message(pkg_terms.VInfo, acc_term.AccountCredit).
		JSON(accountCredit)
}

// GetAll list of accountCredits
func (a *AccountCreditAPI) GetAll(c *gin.Context) {
	resp, params := response.NewParam(c, acc_model.AccountCreditTable)
	var accountCredits []acc_model.AccountCredit
	var err error

	if accountCredits, err = a.Service.GetAll(params); err != nil {
		err = pkg_err.Take(err, "E1180593").Message(pkg_err.SomethingWentWrong).Build()
		resp.Error(err).JSON()
		return
	}

	resp.Record(acc.ListAccountCredit)
	resp.Status(http.StatusOK).
		Message(pkg_terms.ListOfV, acc_term.AccountCredits).
		JSON(accountCredits)
}

// List of accountCredits
func (a *AccountCreditAPI) List(c *gin.Context) {
	resp, params := response.NewParam(c, acc_model.AccountCreditTable)

	data := make(map[string]interface{})
	var err error

	if data["list"], data["count"], err = a.Service.List(params); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(acc.ListAccountCredit)
	resp.Status(http.StatusOK).
		Message(pkg_terms.ListOfV, acc_term.AccountCredits).
		JSON(data)
}

// Create accountCredit
func (a *AccountCreditAPI) Create(c *gin.Context) {
	resp, params := response.NewParam(c, acc_model.AccountCreditTable)
	var accountCredit, createdAccountCredit acc_model.AccountCredit
	var err error

	if err = resp.Bind(&accountCredit, "E1150888", acc_term.AccountCredit); err != nil {
		return
	}

	// start tx
	params.Tx.DB = a.Engine.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			pkg_log.LogError(fmt.Errorf("panic happened in tx mode for %v",
				"accountCredit"), "rollback recover create accountCredit")
			err = pkg_err.New(pkg_err.SomethingWentWrong, "E1145492").
				Message(pkg_err.SomethingWentWrong).Build()
			// rollback tx
			params.Tx.DB.Rollback()
			return
		}
	}()

	if createdAccountCredit, err = a.Service.Create(params.Tx, accountCredit); err != nil {
		resp.Error(err).JSON()
		// rollback tx
		params.Tx.DB.Rollback()
		return
	}

	// commit tx
	params.Tx.DB.Commit()

	resp.Record(acc.CreateAccountCredit, accountCredit, createdAccountCredit)
	resp.Status(http.StatusOK).
		Message(pkg_terms.VCreatedSuccessfully, acc_term.AccountCredit).
		JSON(createdAccountCredit)
}

// Update accountCredit
func (a *AccountCreditAPI) Update(c *gin.Context) {
	resp, params := response.NewParam(c, acc_model.AccountCreditTable)
	var err error
	var accountCredit, accountCreditBefore, accountCreditUpdated acc_model.AccountCredit

	if err = resp.Bind(&accountCredit, "E1135664", acc_term.AccountCredit); err != nil {
		return
	}

	if accountCredit.ID, err = resp.GetID(c.Param("accountCreditID"), "E1162016", acc_term.AccountCredit); err != nil {
		return
	}

	// start tx
	params.Tx.DB = a.Engine.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			pkg_log.LogError(fmt.Errorf("panic happened in tx mode for %v",
				"accountCredit"), "rollback recover create accountCredit")
			err = pkg_err.New(pkg_err.SomethingWentWrong, "E1138430").
				Message(pkg_err.SomethingWentWrong).Build()
			// rollback tx
			params.Tx.DB.Rollback()
			return
		}
	}()

	if accountCreditUpdated, accountCreditBefore, err = a.Service.Save(params.Tx, accountCredit); err != nil {
		resp.Error(err).JSON()
		// rollback tx
		params.Tx.DB.Rollback()
		return
	}

	// commit tx
	params.Tx.DB.Commit()

	resp.Record(acc.UpdateAccountCredit, accountCreditBefore, accountCredit)
	resp.Status(http.StatusOK).
		Message(pkg_terms.VUpdatedSuccessfully, acc_term.AccountCredit).
		JSON(accountCreditUpdated)
}

// Delete accountCredit
func (a *AccountCreditAPI) Delete(c *gin.Context) {
	resp, params := response.NewParam(c, acc_model.AccountCreditTable)
	var err error
	var accountCredit acc_model.AccountCredit
	var id uint

	if id, err = resp.GetID(c.Param("accountCreditID"), "E1144110", acc_term.AccountCredit); err != nil {
		return
	}

	// start tx
	params.Tx.DB = a.Engine.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			pkg_log.LogError(fmt.Errorf("panic happened in tx mode for %v",
				"accountCredit"), "rollback recover create accountCredit")
			err = pkg_err.New(pkg_err.SomethingWentWrong, "E1199821").
				Message(pkg_err.SomethingWentWrong).Build()
			// rollback tx
			params.Tx.DB.Rollback()
			return
		}
	}()

	if accountCredit, err = a.Service.Delete(params.Tx, id); err != nil {
		resp.Error(err).JSON()
		// rollback tx
		params.Tx.DB.Rollback()
		return
	}

	// commit tx
	params.Tx.DB.Commit()

	resp.Record(acc.DeleteAccountCredit, accountCredit)
	resp.Status(http.StatusOK).
		Message(pkg_terms.VDeletedSuccessfully, acc_term.AccountCredit).
		JSON()
}
