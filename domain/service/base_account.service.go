package service

import (
	"GoRestify/domain/base/base_model"
	"GoRestify/domain/base/base_repo"
	"GoRestify/domain/base/base_term"
	"GoRestify/domain/base/enum/account_type"
	"GoRestify/internal/core"
	"fmt"

	"GoRestify/pkg/param"
	"GoRestify/pkg/pkg_err"
	"GoRestify/pkg/pkg_log"
	"GoRestify/pkg/pkg_password"
	"GoRestify/pkg/tx"

	"GoRestify/pkg/validator"
)

// BaseAccountServ for injecting  base_repo
type BaseAccountServ struct {
	Repo   base_repo.AccountRepo
	Engine *core.Engine
}

// ProvideBaseAccountService for account is used in wire
func ProvideBaseAccountService(accountRepo base_repo.AccountRepo) BaseAccountServ {
	return BaseAccountServ{
		Repo:   accountRepo,
		Engine: accountRepo.Engine,
	}
}

// FindByID for getting account by its id
func (s *BaseAccountServ) FindByID(tx tx.Tx, id uint) (account base_model.Account, err error) {

	key := fmt.Sprintf("%v-%v", base_term.Account, id)
	if ok := s.Engine.RedisCacheAPI.GetCache(tx, key, &account); ok {
		return
	}

	if account, err = s.Repo.FindByID(tx, id); err != nil {
		pkg_err.Log(err, "E1133637", "can't fetch the account", id)
		return
	}

	err = s.Engine.RedisCacheAPI.Set(key, account)

	return
}

// FindByUsername for getting account by its username
func (s *BaseAccountServ) FindByUsername(tx tx.Tx, username string) (account base_model.Account, err error) {

	if account, err = s.Repo.FindByUsername(tx, username); err != nil {
		pkg_err.Log(err, "E1167315", "can't fetch the user", username)
		return
	}

	return
}

// GetAll of accounts, it supports pagination and search and return count
func (s *BaseAccountServ) GetAll(params param.Param) (accounts []base_model.Account, err error) {

	if ok := s.Engine.RedisCacheAPI.GetCache(params.Tx, base_term.Accounts, &accounts); ok {
		return
	}

	params.Limit = 100000
	if accounts, err = s.Repo.List(params); err != nil {
		pkg_log.CheckError(err, "error in accounts list")
		return
	}

	for i := range accounts {
		accounts[i].Password = ""
	}

	err = s.Engine.RedisCacheAPI.Set(base_term.Accounts, accounts)

	return
}

// List of accounts, it supports pagination and search and return count
func (s *BaseAccountServ) List(params param.Param) (accounts []base_model.Account,
	count int64, err error) {

	if accounts, err = s.Repo.List(params); err != nil {
		pkg_log.CheckError(err, "error in accounts list")
		return
	}

	for i := range accounts {
		accounts[i].Password = ""
	}

	if count, err = s.Repo.Count(params); err != nil {
		pkg_log.CheckError(err, "error in accounts count")
	}

	return
}

// Create a account
func (s *BaseAccountServ) Create(tx tx.Tx, account base_model.Account) (createdAccount base_model.Account, err error) {

	if err = validator.ValidateModel(account, base_term.Account, validator.Create); err != nil {
		err = pkg_err.TickValidate(err, "E1192923", pkg_err.ValidationFailed, account)
		return
	}

	if account.Password, err = pkg_password.Hash(account.Password, s.Engine.Envs[core.PasswordSalt]); err != nil {
		err = pkg_err.Log(err, "E1181138", "error in hashing password", account)
		return
	}

	if createdAccount, err = s.Repo.Create(tx, account); err != nil {
		pkg_err.Log(err, "E1129381", "error in creating account", account)
		return
	}
	createdAccount.Password = ""

	switch account.Type {
	case account_type.User:
		account.User = base_model.User{
			ID:     createdAccount.ID,
			RoleID: account.RoleID,
		}
		if _, err = BaseUserService.Create(tx, account.User); err != nil {
			pkg_err.Log(err, "E1140039", "error in creating user", account)
			return
		}
	}

	s.Engine.RedisCacheAPI.Delete(base_term.Accounts)

	return
}

// Save a account, if it is exists update it, if not create it
func (s *BaseAccountServ) Save(tx tx.Tx, account base_model.Account) (updatedAccount, accountBefore base_model.Account, err error) {

	if err = validator.ValidateModel(account, base_term.Account, validator.Update); err != nil {
		err = pkg_err.TickValidate(err, "E1158219", pkg_err.ValidationFailed, account)
		return
	}

	if accountBefore, err = s.FindByID(tx, account.ID); err != nil {
		pkg_err.Log(err, "E1124905", "can't fetch account by id for saving it", account.ID)
		return
	}

	if account.Password != "" {
		if account.Password, err = pkg_password.Hash(account.Password, s.Engine.Envs[core.PasswordSalt]); err != nil {
			err = pkg_err.Log(err, "E1115211", "error in hashing password", account)
			return
		}
	} else {
		account.Password = accountBefore.Password
	}

	if updatedAccount, err = s.Repo.Save(tx, account); err != nil {
		pkg_err.Log(err, "E1116402", "account not saved")
		return
	}
	updatedAccount.Password = ""

	switch account.Type {
	case account_type.User:
		account.User = base_model.User{
			ID:     updatedAccount.ID,
			RoleID: account.RoleID,
		}
		if _, _, err = BaseUserService.Save(tx, account.User); err != nil {
			pkg_err.Log(err, "E1125792", "error in creating user", account)
			return
		}
	}

	key := fmt.Sprintf("%v-%v", base_term.Account, updatedAccount.ID)
	if err = s.Engine.RedisCacheAPI.Delete(key); err != nil {
		return
	}

	s.Engine.RedisCacheAPI.Delete(base_term.Accounts)

	return
}

// Delete account, it is soft delete
func (s *BaseAccountServ) Delete(tx tx.Tx, id uint) (account base_model.Account, err error) {

	if account, err = s.FindByID(tx, id); err != nil {
		pkg_err.Log(err, "E1160391", "account not found for deleting")
		return
	}

	if err = s.Repo.Delete(tx, account); err != nil {
		pkg_err.Log(err, "E1172081", "account not deleted")
		return
	}

	key := fmt.Sprintf("%v-%v", base_term.Account, account.ID)
	s.Engine.RedisCacheAPI.Delete(key)
	s.Engine.RedisCacheAPI.Delete(base_term.Accounts)

	return
}
