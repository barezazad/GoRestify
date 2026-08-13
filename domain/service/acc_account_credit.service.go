package service

import (
	"fmt"

	"GoRestify/domain/acc/acc_model"
	"GoRestify/domain/acc/acc_repo"
	"GoRestify/domain/acc/acc_term"
	"GoRestify/internal/core"

	"GoRestify/pkg/decimal"
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

// FindByAccountIDAndCurrency for getting accountCredit by account + currency
func (s *AccAccountCreditServ) FindByAccountIDAndCurrency(tx tx.Tx, accountID, currencyID uint) (accountCredit acc_model.AccountCredit, err error) {

	if accountCredit, err = s.Repo.FindByAccountIDAndCurrency(tx, accountID, currencyID); err != nil {
		err = pkg_err.Log(err, "E1196479", "can't fetch the accountCredit", accountID)
		return
	}

	return
}

// Ensure returns the credit row for (account, currency), creating a zero balance if missing.
// When tx.IsLock is set, the returned row is locked for update.
func (s *AccAccountCreditServ) Ensure(txSession tx.Tx, accountID, currencyID uint) (accountCredit acc_model.AccountCredit, err error) {
	accountCredit, err = s.Repo.FindByAccountIDAndCurrency(txSession, accountID, currencyID)
	if err == nil {
		return
	}

	zero := decimal.NewFromInt(0).Build()
	created, createErr := s.Repo.Create(txSession, acc_model.AccountCredit{
		AccountID:  accountID,
		CurrencyID: currencyID,
		Balance:    zero,
	})
	if createErr != nil {
		// race: another TX created it — reload with lock
		accountCredit, err = s.Repo.FindByAccountIDAndCurrency(txSession, accountID, currencyID)
		if err != nil {
			err = pkg_err.Take(createErr, "E1201020", "ensure account credit failed").Build()
		}
		return
	}

	_ = s.Engine.RedisCacheAPI.Delete(acc_term.AccountCredits)
	accountCredit = created

	// re-fetch under lock if requested
	if txSession.IsLock {
		accountCredit, err = s.Repo.FindByAccountIDAndCurrency(txSession, accountID, currencyID)
	}
	return
}

// GetAll of accountCredits
func (s *AccAccountCreditServ) GetAll(params param.Param) (accountCredits []acc_model.AccountCredit, err error) {

	if ok := s.Engine.RedisCacheAPI.GetCache(params.Tx, acc_term.AccountCredits, &accountCredits); ok {
		return
	}

	if accountCredits, err = s.Repo.GetAll(); err != nil {
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

// Create a accountCredit (bootstrap zero balance). Prefer Ensure from DoTransaction.
func (s *AccAccountCreditServ) Create(tx tx.Tx, accountCredit acc_model.AccountCredit) (createdAccountCredit acc_model.AccountCredit, err error) {

	if err = validator.ValidateModel(accountCredit, acc_term.AccountCredit, validator.Create); err != nil {
		err = pkg_err.TickValidate(err, "E1680067", pkg_err.ValidationFailed, accountCredit)
		return
	}

	if accountCredit.Balance == nil {
		accountCredit.Balance = decimal.NewFromInt(0).Build()
	}

	if createdAccountCredit, err = s.Repo.Create(tx, accountCredit); err != nil {
		pkg_err.Log(err, "E1626674", "error in creating accountCredit", accountCredit)
		return
	}

	s.Engine.RedisCacheAPI.Delete(acc_term.AccountCredits)

	return
}

// Save updates balance — used by DoTransaction. Direct HTTP update should not be exposed.
func (s *AccAccountCreditServ) Save(txSession tx.Tx, accountCredit acc_model.AccountCredit) (updatedAccountCredit, accountCreditBefore acc_model.AccountCredit, err error) {

	if accountCreditBefore, err = s.Repo.FindByID(txSession, accountCredit.ID); err != nil {
		pkg_err.Log(err, "E1625869", "can't fetch accountCredit by id for saving it", accountCredit.ID)
		return
	}

	if updatedAccountCredit, err = s.Repo.Save(txSession, accountCredit); err != nil {
		pkg_err.Log(err, "E1139340", "accountCredit not saved")
		return
	}

	key := fmt.Sprintf("%v-%v", acc_term.AccountCredit, updatedAccountCredit.ID)
	_ = s.Engine.RedisCacheAPI.Delete(key)
	_ = s.Engine.RedisCacheAPI.Delete(acc_term.AccountCredits)

	return
}

// Delete is disabled — balances are updated only via DoTransaction.
func (s *AccAccountCreditServ) Delete(tx tx.Tx, id uint) (accountCredit acc_model.AccountCredit, err error) {
	err = pkg_err.New(acc_term.LedgerRecordsAreImmutable, "E1201021").
		Message(acc_term.LedgerRecordsAreImmutable).Custom(pkg_err.ValidationFailedErr).Build()
	return
}
