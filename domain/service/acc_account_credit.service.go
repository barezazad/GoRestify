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

// AccAccountCreditServ for injecting  acc_repo
type AccAccountCreditServ struct {
	Repo   acc_repo.AccountCreditRepo
	Engine *core.Engine
}

// ProvideAccAccountCreditService for accountCredit is used in wire
func ProvideAccAccountCreditService(accountCreditRepo acc_repo.AccountCreditRepo) AccAccountCreditServ {
	return AccAccountCreditServ{
		Repo:   accountCreditRepo,
		Engine: accountCreditRepo.Engine,
	}
}

// FindByID for getting accountCredit by its id
func (s *AccAccountCreditServ) FindByID(tx tx.Tx, id uint) (accountCredit acc_model.AccountCredit, err error) {

	key := fmt.Sprintf("%v-%v", acc_term.AccountCredit, id)
	if ok := s.Engine.RedisCacheAPI.GetCache(tx, key, &accountCredit); ok {
		return
	}

	if accountCredit, err = s.Repo.FindByID(tx, id); err != nil {
		pkg_err.Log(err, "E1673780", "can't fetch the accountCredit", id)
		return
	}

	err = s.Engine.RedisCacheAPI.Set(key, accountCredit)

	return
}

// FindByAccountIDAndCurrency for getting accountCredit by its id
func (s *AccAccountCreditServ) FindByAccountIDAndCurrency(tx tx.Tx, accountID, currencyID uint) (accountCredit acc_model.AccountCredit, err error) {

	if accountCredit, err = s.Repo.FindByAccountIDAndCurrency(tx, accountID, currencyID); err != nil {
		err = pkg_err.Log(err, "E1196479", "can't fetch the accountCredit", accountID)
		return
	}

	return
}

// GetAll of accountCredits, it supports pagination and search and return count
func (s *AccAccountCreditServ) GetAll(params param.Param) (accountCredits []acc_model.AccountCredit, err error) {

	if ok := s.Engine.RedisCacheAPI.GetCache(params.Tx, acc_term.AccountCredits, &accountCredits); ok {
		return
	}

	params.Pagination.Limit = 100000
	if accountCredits, err = s.Repo.List(params); err != nil {
		pkg_log.CheckError(err, "error in accountCredits list")
		return
	}

	err = s.Engine.RedisCacheAPI.Set(acc_term.AccountCredits, accountCredits)

	return
}

// List of accountCredits, it supports pagination and search and return count
func (s *AccAccountCreditServ) List(params param.Param) (accountCredits []acc_model.AccountCredit,
	count int64, err error) {

	if accountCredits, err = s.Repo.List(params); err != nil {
		pkg_log.CheckError(err, "error in accountCredits list")
		return
	}

	if count, err = s.Repo.Count(params); err != nil {
		pkg_log.CheckError(err, "error in accountCredits count")
	}

	return
}

// Create a accountCredit
func (s *AccAccountCreditServ) Create(tx tx.Tx, accountCredit acc_model.AccountCredit) (createdAccountCredit acc_model.AccountCredit, err error) {

	if err = validator.ValidateModel(accountCredit, acc_term.AccountCredit, validator.Create); err != nil {
		err = pkg_err.TickValidate(err, "E1680067", pkg_err.ValidationFailed, accountCredit)
		return
	}

	if createdAccountCredit, err = s.Repo.Create(tx, accountCredit); err != nil {
		pkg_err.Log(err, "E1626674", "error in creating accountCredit", accountCredit)
		return
	}

	s.Engine.RedisCacheAPI.Delete(acc_term.AccountCredits)

	return
}

// Save a accountCredit, if it is exists update it, if not create it
func (s *AccAccountCreditServ) Save(tx tx.Tx, accountCredit acc_model.AccountCredit) (updatedAccountCredit, accountCreditBefore acc_model.AccountCredit, err error) {

	if err = validator.ValidateModel(accountCredit, acc_term.AccountCredit, validator.Update); err != nil {
		err = pkg_err.TickValidate(err, "E1679868", pkg_err.ValidationFailed, accountCredit)
		return
	}

	if accountCreditBefore, err = s.FindByID(tx, accountCredit.ID); err != nil {
		pkg_err.Log(err, "E1625869", "can't fetch accountCredit by id for saving it", accountCredit.ID)
		return
	}

	if updatedAccountCredit, err = s.Repo.Save(tx, accountCredit); err != nil {
		pkg_err.Log(err, "E1139340", "accountCredit not saved")
		return
	}

	key := fmt.Sprintf("%v-%v", acc_term.AccountCredit, updatedAccountCredit.ID)
	if err = s.Engine.RedisCacheAPI.Delete(key); err != nil {
		return
	}

	s.Engine.RedisCacheAPI.Delete(acc_term.AccountCredits)

	return
}

// Delete accountCredit, it is soft delete
func (s *AccAccountCreditServ) Delete(tx tx.Tx, id uint) (accountCredit acc_model.AccountCredit, err error) {

	if accountCredit, err = s.FindByID(tx, id); err != nil {
		pkg_err.Log(err, "E1653653", "accountCredit not found for deleting")
		return
	}

	if err = s.Repo.Delete(tx, accountCredit); err != nil {
		pkg_err.Log(err, "E1681259", "accountCredit not deleted")
		return
	}

	key := fmt.Sprintf("%v-%v", acc_term.AccountCredit, accountCredit.ID)
	s.Engine.RedisCacheAPI.Delete(key)
	s.Engine.RedisCacheAPI.Delete(acc_term.AccountCredits)

	return
}
