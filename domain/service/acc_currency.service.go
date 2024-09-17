package service

import (
	"GoRestify/domain/acc/acc_model"
	"GoRestify/domain/acc/acc_repo"
	"GoRestify/domain/acc/acc_term"
	"GoRestify/internal/core"
	"fmt"

	"GoRestify/pkg/param"
	"GoRestify/pkg/pkg_err"
	"GoRestify/pkg/pkg_log"
	"GoRestify/pkg/tx"

	"GoRestify/pkg/validator"
)

// AccCurrencyServ for injecting  acc_repo
type AccCurrencyServ struct {
	Repo   acc_repo.CurrencyRepo
	Engine *core.Engine
}

// ProvideAccCurrencyService for currency is used in wire
func ProvideAccCurrencyService(currencyRepo acc_repo.CurrencyRepo) AccCurrencyServ {
	return AccCurrencyServ{
		Repo:   currencyRepo,
		Engine: currencyRepo.Engine,
	}
}

// FindByID for getting currency by its id
func (s *AccCurrencyServ) FindByID(tx tx.Tx, id uint) (currency acc_model.Currency, err error) {

	key := fmt.Sprintf("%v-%v", acc_term.Currency, id)
	if ok := s.Engine.RedisCacheAPI.GetCache(tx, key, &currency); ok {
		return
	}

	if currency, err = s.Repo.FindByID(tx, id); err != nil {
		pkg_err.Log(err, "E1127722", "can't fetch the currency", id)
		return
	}

	err = s.Engine.RedisCacheAPI.Set(key, currency)

	return
}

// GetAll of currencies, it supports pagination and search and return count
func (s *AccCurrencyServ) GetAll(params param.Param) (currencies []acc_model.Currency, err error) {

	if ok := s.Engine.RedisCacheAPI.GetCache(params.Tx, acc_term.Currencies, &currencies); ok {
		return
	}

	params.Limit = 100000
	if currencies, err = s.Repo.List(params); err != nil {
		pkg_log.CheckError(err, "error in currencies list")
		return
	}

	err = s.Engine.RedisCacheAPI.Set(acc_term.Currencies, currencies)

	return
}

// List of currencies, it supports pagination and search and return count
func (s *AccCurrencyServ) List(params param.Param) (currencies []acc_model.Currency,
	count int64, err error) {

	if currencies, err = s.Repo.List(params); err != nil {
		pkg_log.CheckError(err, "error in currencies list")
		return
	}

	if count, err = s.Repo.Count(params); err != nil {
		pkg_log.CheckError(err, "error in currencies count")
	}

	return
}

// Create a currency
func (s *AccCurrencyServ) Create(tx tx.Tx, currency acc_model.Currency) (createdCurrency acc_model.Currency, err error) {

	if err = validator.ValidateModel(currency, acc_term.Currency, validator.Create); err != nil {
		err = pkg_err.TickValidate(err, "E1156321", pkg_err.ValidationFailed, currency)
		return
	}

	if createdCurrency, err = s.Repo.Create(tx, currency); err != nil {
		pkg_err.Log(err, "E1196495", "error in creating currency", currency)
		return
	}

	s.Engine.RedisCacheAPI.Delete(acc_term.Currencies)

	return
}

// Save a currency, if it is exists update it, if not create it
func (s *AccCurrencyServ) Save(tx tx.Tx, currency acc_model.Currency) (updatedCurrency, currencyBefore acc_model.Currency, err error) {

	if err = validator.ValidateModel(currency, acc_term.Currency, validator.Update); err != nil {
		err = pkg_err.TickValidate(err, "E1128005", pkg_err.ValidationFailed, currency)
		return
	}

	if currencyBefore, err = s.FindByID(tx, currency.ID); err != nil {
		pkg_err.Log(err, "E1134515", "can't fetch currency by id for saving it", currency.ID)
		return
	}

	if updatedCurrency, err = s.Repo.Save(tx, currency); err != nil {
		pkg_err.Log(err, "E1147722", "currency not saved")
		return
	}

	key := fmt.Sprintf("%v-%v", acc_term.Currency, updatedCurrency.ID)
	if err = s.Engine.RedisCacheAPI.Delete(key); err != nil {
		return
	}

	s.Engine.RedisCacheAPI.Delete(acc_term.Currencies)

	return
}

// Delete currency, it is soft delete
func (s *AccCurrencyServ) Delete(tx tx.Tx, id uint) (currency acc_model.Currency, err error) {

	if currency, err = s.FindByID(tx, id); err != nil {
		pkg_err.Log(err, "E1150148", "currency not found for deleting")
		return
	}

	if err = s.Repo.Delete(tx, currency); err != nil {
		pkg_err.Log(err, "E1198830", "currency not deleted")
		return
	}

	key := fmt.Sprintf("%v-%v", acc_term.Currency, currency.ID)
	s.Engine.RedisCacheAPI.Delete(key)
	s.Engine.RedisCacheAPI.Delete(acc_term.Currencies)

	return
}
