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

// CurrencyAPI for injecting currency service
type CurrencyAPI struct {
	Service service.AccCurrencyServ
	Engine  *core.Engine
}

// ProvideCurrencyAPI .
func ProvideCurrencyAPI(c service.AccCurrencyServ) CurrencyAPI {
	return CurrencyAPI{Service: c, Engine: c.Engine}
}

// FindByID is used for fetch an currency by its id
func (a *CurrencyAPI) FindByID(c *gin.Context) {
	resp, params := response.NewParam(c, acc_model.CurrencyTable)
	var err error
	var currency acc_model.Currency
	var id uint

	if id, err = resp.GetID(c.Param("currencyID"), "E1173530", acc_term.Currency); err != nil {
		return
	}

	if currency, err = a.Service.FindByID(params.Tx, id); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(acc.ViewCurrency)
	resp.Status(http.StatusOK).
		Message(pkg_terms.VInfo, acc_term.Currency).
		JSON(currency)
}

// GetAll list of currencies
func (a *CurrencyAPI) GetAll(c *gin.Context) {
	resp, params := response.NewParam(c, acc_model.CurrencyTable)
	var currencies []acc_model.Currency
	var err error

	if currencies, err = a.Service.GetAll(params); err != nil {
		err = pkg_err.Take(err, "E1151472").Message(pkg_err.SomethingWentWrong).Build()
		resp.Error(err).JSON()
		return
	}

	resp.Record(acc.ListCurrency)
	resp.Status(http.StatusOK).
		Message(pkg_terms.ListOfV, acc_term.Currencies).
		JSON(currencies)
}

// List of currencies
func (a *CurrencyAPI) List(c *gin.Context) {
	resp, params := response.NewParam(c, acc_model.CurrencyTable)

	data := make(map[string]interface{})
	var err error

	if data["list"], data["count"], err = a.Service.List(params); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(acc.ListCurrency)
	resp.Status(http.StatusOK).
		Message(pkg_terms.ListOfV, acc_term.Currencies).
		JSON(data)
}

// Create currency
func (a *CurrencyAPI) Create(c *gin.Context) {
	resp, params := response.NewParam(c, acc_model.CurrencyTable)
	var currency, createdCurrency acc_model.Currency
	var err error

	if err = resp.Bind(&currency, "E1196146", acc_term.Currency); err != nil {
		return
	}

	// start tx
	params.Tx.DB = a.Engine.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			pkg_log.LogError(fmt.Errorf("panic happened in tx mode for %v",
				"currency"), "rollback recover create currency")
			err = pkg_err.New(pkg_err.SomethingWentWrong, "E1174851").
				Message(pkg_err.SomethingWentWrong).Build()
			// rollback tx
			params.Tx.DB.Rollback()
			return
		}
	}()

	if createdCurrency, err = a.Service.Create(params.Tx, currency); err != nil {
		resp.Error(err).JSON()
		// rollback tx
		params.Tx.DB.Rollback()
		return
	}

	// commit tx
	params.Tx.DB.Commit()

	resp.Record(acc.CreateCurrency, currency, createdCurrency)
	resp.Status(http.StatusOK).
		Message(pkg_terms.VCreatedSuccessfully, acc_term.Currency).
		JSON(createdCurrency)
}

// Update currency
func (a *CurrencyAPI) Update(c *gin.Context) {
	resp, params := response.NewParam(c, acc_model.CurrencyTable)
	var err error
	var currency, currencyBefore, currencyUpdated acc_model.Currency

	if err = resp.Bind(&currency, "E1121644", acc_term.Currency); err != nil {
		return
	}

	if currency.ID, err = resp.GetID(c.Param("currencyID"), "E1159908", acc_term.Currency); err != nil {
		return
	}

	// start tx
	params.Tx.DB = a.Engine.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			pkg_log.LogError(fmt.Errorf("panic happened in tx mode for %v",
				"currency"), "rollback recover create currency")
			err = pkg_err.New(pkg_err.SomethingWentWrong, "E1112114").
				Message(pkg_err.SomethingWentWrong).Build()
			// rollback tx
			params.Tx.DB.Rollback()
			return
		}
	}()

	if currencyUpdated, currencyBefore, err = a.Service.Save(params.Tx, currency); err != nil {
		resp.Error(err).JSON()
		// rollback tx
		params.Tx.DB.Rollback()
		return
	}

	// commit tx
	params.Tx.DB.Commit()

	resp.Record(acc.UpdateCurrency, currencyBefore, currency)
	resp.Status(http.StatusOK).
		Message(pkg_terms.VUpdatedSuccessfully, acc_term.Currency).
		JSON(currencyUpdated)
}

// Delete currency
func (a *CurrencyAPI) Delete(c *gin.Context) {
	resp, params := response.NewParam(c, acc_model.CurrencyTable)
	var err error
	var currency acc_model.Currency
	var id uint

	if id, err = resp.GetID(c.Param("currencyID"), "E1136191", acc_term.Currency); err != nil {
		return
	}

	// start tx
	params.Tx.DB = a.Engine.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			pkg_log.LogError(fmt.Errorf("panic happened in tx mode for %v",
				"currency"), "rollback recover create currency")
			err = pkg_err.New(pkg_err.SomethingWentWrong, "E1125252").
				Message(pkg_err.SomethingWentWrong).Build()
			// rollback tx
			params.Tx.DB.Rollback()
			return
		}
	}()

	if currency, err = a.Service.Delete(params.Tx, id); err != nil {
		resp.Error(err).JSON()
		// rollback tx
		params.Tx.DB.Rollback()
		return
	}

	// commit tx
	params.Tx.DB.Commit()

	resp.Record(acc.DeleteCurrency, currency)
	resp.Status(http.StatusOK).
		Message(pkg_terms.VDeletedSuccessfully, acc_term.Currency).
		JSON()
}
