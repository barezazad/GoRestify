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

// Create transaction
func (a *TransactionAPI) Create(c *gin.Context) {
	resp, params := response.NewParam(c, acc_model.TransactionTable)
	var transaction, createdTransaction acc_model.Transaction
	var err error

	if err = resp.Bind(&transaction, "E1115292", acc_term.Transaction); err != nil {
		return
	}

	// start tx
	params.Tx.DB = a.Engine.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			pkg_log.LogError(fmt.Errorf("panic happened in tx mode for %v",
				"transaction"), "rollback recover create transaction")
			err = pkg_err.New(pkg_err.SomethingWentWrong, "E1154151").
				Message(pkg_err.SomethingWentWrong).Build()
			// rollback tx
			params.Tx.DB.Rollback()
			return
		}
	}()

	if createdTransaction, err = a.Service.Create(params.Tx, transaction); err != nil {
		resp.Error(err).JSON()
		// rollback tx
		params.Tx.DB.Rollback()
		return
	}

	// commit tx
	params.Tx.DB.Commit()

	resp.Record(acc.CreateTransaction, transaction, createdTransaction)
	resp.Status(http.StatusOK).
		Message(pkg_terms.VCreatedSuccessfully, acc_term.Transaction).
		JSON(createdTransaction)
}

// Update transaction
func (a *TransactionAPI) Update(c *gin.Context) {
	resp, params := response.NewParam(c, acc_model.TransactionTable)
	var err error
	var transaction, transactionBefore, transactionUpdated acc_model.Transaction

	if err = resp.Bind(&transaction, "E1198257", acc_term.Transaction); err != nil {
		return
	}

	if transaction.ID, err = resp.GetID(c.Param("transactionID"), "E1188360", acc_term.Transaction); err != nil {
		return
	}

	// start tx
	params.Tx.DB = a.Engine.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			pkg_log.LogError(fmt.Errorf("panic happened in tx mode for %v",
				"transaction"), "rollback recover create transaction")
			err = pkg_err.New(pkg_err.SomethingWentWrong, "E1197208").
				Message(pkg_err.SomethingWentWrong).Build()
			// rollback tx
			params.Tx.DB.Rollback()
			return
		}
	}()

	if transactionUpdated, transactionBefore, err = a.Service.Save(params.Tx, transaction); err != nil {
		resp.Error(err).JSON()
		// rollback tx
		params.Tx.DB.Rollback()
		return
	}

	// commit tx
	params.Tx.DB.Commit()

	resp.Record(acc.UpdateTransaction, transactionBefore, transaction)
	resp.Status(http.StatusOK).
		Message(pkg_terms.VUpdatedSuccessfully, acc_term.Transaction).
		JSON(transactionUpdated)
}

// Delete transaction
func (a *TransactionAPI) Delete(c *gin.Context) {
	resp, params := response.NewParam(c, acc_model.TransactionTable)
	var err error
	var transaction acc_model.Transaction
	var id uint

	if id, err = resp.GetID(c.Param("transactionID"), "E1112945", acc_term.Transaction); err != nil {
		return
	}

	// start tx
	params.Tx.DB = a.Engine.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			pkg_log.LogError(fmt.Errorf("panic happened in tx mode for %v",
				"transaction"), "rollback recover create transaction")
			err = pkg_err.New(pkg_err.SomethingWentWrong, "E1115395").
				Message(pkg_err.SomethingWentWrong).Build()
			// rollback tx
			params.Tx.DB.Rollback()
			return
		}
	}()

	if transaction, err = a.Service.Delete(params.Tx, id); err != nil {
		resp.Error(err).JSON()
		// rollback tx
		params.Tx.DB.Rollback()
		return
	}

	// commit tx
	params.Tx.DB.Commit()

	resp.Record(acc.DeleteTransaction, transaction)
	resp.Status(http.StatusOK).
		Message(pkg_terms.VDeletedSuccessfully, acc_term.Transaction).
		JSON()
}
