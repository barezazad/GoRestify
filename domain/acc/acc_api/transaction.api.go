package acc_api

import (
	"fmt"
	"net/http"

	"GoRestify/domain/acc"
	"GoRestify/domain/acc/acc_model"
	"GoRestify/domain/acc/acc_term"
	"GoRestify/domain/service"
	"GoRestify/internal/core"

	"GoRestify/pkg/pkg_err"
	"GoRestify/pkg/pkg_log"
	"GoRestify/pkg/pkg_terms"
	"GoRestify/pkg/response"

	"github.com/gin-gonic/gin"
)

// TransactionAPI for injecting transaction service
type TransactionAPI struct {
	Service service.AccTransactionServ
	Engine  *core.Engine
}

// ProvideTransactionAPI .
func ProvideTransactionAPI(c service.AccTransactionServ) TransactionAPI {
	return TransactionAPI{Service: c, Engine: c.Engine}
}

// FindByID is used for fetch an transaction by its id
func (a *TransactionAPI) FindByID(c *gin.Context) {
	resp, params := response.NewParam(c, acc_model.TransactionTable)
	var err error
	var transaction acc_model.Transaction
	var id uint

	if id, err = resp.GetID(c.Param("transactionID"), "E1162700", acc_term.Transaction); err != nil {
		return
	}

	if transaction, err = a.Service.FindByID(params.Tx, id); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(acc.ViewTransaction)
	resp.Status(http.StatusOK).
		Message(pkg_terms.VInfo, acc_term.Transaction).
		JSON(transaction)
}

// GetAll list of transactions
func (a *TransactionAPI) GetAll(c *gin.Context) {
	resp, params := response.NewParam(c, acc_model.TransactionTable)
	var transactions []acc_model.Transaction
	var err error

	if transactions, err = a.Service.GetAll(params); err != nil {
		err = pkg_err.Take(err, "E1160850").Message(pkg_err.SomethingWentWrong).Build()
		resp.Error(err).JSON()
		return
	}

	resp.Record(acc.ListTransaction)
	resp.Status(http.StatusOK).
		Message(pkg_terms.ListOfV, acc_term.Transactions).
		JSON(transactions)
}

// List of transactions
func (a *TransactionAPI) List(c *gin.Context) {
	resp, params := response.NewParam(c, acc_model.TransactionTable)

	data := make(map[string]interface{})
	var err error

	if data["list"], data["count"], err = a.Service.List(params); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(acc.ListTransaction)
	resp.Status(http.StatusOK).
		Message(pkg_terms.ListOfV, acc_term.Transactions).
		JSON(data)
}

// Create posts a transaction via DoTransaction (balanced slots required).
func (a *TransactionAPI) Create(c *gin.Context) {
	resp, params := response.NewParam(c, acc_model.TransactionTable)
	var transaction, createdTransaction acc_model.Transaction
	var err error

	if err = resp.Bind(&transaction, "E1115292", acc_term.Transaction); err != nil {
		return
	}

	params.Tx.DB = a.Engine.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			pkg_log.LogError(fmt.Errorf("panic happened in tx mode for %v",
				"transaction"), "rollback recover create transaction")
			err = pkg_err.New(pkg_err.SomethingWentWrong, "E1154151").
				Message(pkg_err.SomethingWentWrong).Build()
			params.Tx.DB.Rollback()
			return
		}
	}()

	if createdTransaction, err = a.Service.DoTransaction(params.Tx, transaction); err != nil {
		resp.Error(err).JSON()
		params.Tx.DB.Rollback()
		return
	}

	params.Tx.DB.Commit()

	resp.Record(acc.CreateTransaction, transaction, createdTransaction)
	resp.Status(http.StatusOK).
		Message(pkg_terms.VCreatedSuccessfully, acc_term.Transaction).
		JSON(createdTransaction)
}
