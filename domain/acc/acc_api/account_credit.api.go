package acc_api

import (
	"net/http"

	"GoRestify/domain/acc"
	"GoRestify/domain/acc/acc_model"
	"GoRestify/domain/acc/acc_term"
	"GoRestify/domain/service"
	"GoRestify/internal/core"

	"GoRestify/pkg/pkg_err"
	"GoRestify/pkg/pkg_terms"
	"GoRestify/pkg/response"

	"github.com/gin-gonic/gin"
)

// AccountCreditAPI for injecting accountCredit service (read-only HTTP surface; writes via DoTransaction)
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
