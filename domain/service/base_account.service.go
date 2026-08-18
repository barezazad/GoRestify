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
	"GoRestify/pkg/pkg_types"
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

func accountsCacheKey(accountType pkg_types.Enum) string {
	return fmt.Sprintf("%v-%v", base_term.Accounts, accountType)
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

// FindByIDAndType for getting account by its id and type
func (s *BaseAccountServ) FindByIDAndType(tx tx.Tx, id uint, accountType pkg_types.Enum) (account base_model.Account, err error) {

	key := fmt.Sprintf("%v-%v", base_term.Account, id)
	if ok := s.Engine.RedisCacheAPI.GetCache(tx, key, &account); ok {
		if account.Type != accountType {
			err = pkg_err.New(pkg_err.RecordNotFound, "E1178401").
				Custom(pkg_err.NotFoundErr).Message(pkg_err.RecordNotFound).Build()
			account = base_model.Account{}
			return
		}
		return
	}

	if account, err = s.Repo.FindByIDAndType(tx, id, string(accountType)); err != nil {
		pkg_err.Log(err, "E1178402", "can't fetch the account by type", id, accountType)
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

// GetAll of accounts without filter, order, or pagination
func (s *BaseAccountServ) GetAll(params param.Param) (accounts []base_model.Account, err error) {

	if ok := s.Engine.RedisCacheAPI.GetCache(params.Tx, base_term.Accounts, &accounts); ok {
		return
	}

	if accounts, err = s.Repo.GetAll(); err != nil {
		pkg_log.CheckError(err, "error in accounts list")
		return
	}

	for i := range accounts {
		accounts[i].Password = ""
	}

	err = s.Engine.RedisCacheAPI.Set(base_term.Accounts, accounts)

	return
}

// GetAllByType of accounts for a given type without filter, order, or pagination
func (s *BaseAccountServ) GetAllByType(params param.Param, accountType pkg_types.Enum) (accounts []base_model.Account, err error) {

	cacheKey := accountsCacheKey(accountType)
	if ok := s.Engine.RedisCacheAPI.GetCache(params.Tx, cacheKey, &accounts); ok {
		return
	}

	if accounts, err = s.Repo.GetAllByType(string(accountType)); err != nil {
		pkg_log.CheckError(err, "error in accounts list by type")
		return
	}

	for i := range accounts {
		accounts[i].Password = ""
	}

	err = s.Engine.RedisCacheAPI.Set(cacheKey, accounts)

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

// ListByType lists accounts filtered by type, with pagination and search
func (s *BaseAccountServ) ListByType(params param.Param, accountType pkg_types.Enum) (accounts []base_model.Account,
	count int64, err error) {

	if err = params.AddPreCondition("base_accounts.type = ?", string(accountType)); err != nil {
		pkg_err.Log(err, "E1178403", "invalid account type precondition", accountType)
		return
	}

	return s.List(params)
}

// Create a account
func (s *BaseAccountServ) Create(tx tx.Tx, account base_model.Account) (createdAccount base_model.Account, err error) {

	if err = validator.ValidateModel(account, base_term.Account, validator.Create); err != nil {
		err = pkg_err.TickValidate(err, "E1192923", pkg_err.ValidationFailed, account)
		return
	}

	if account.Type == account_type.User && account.RoleID == 0 {
		err = pkg_err.New("role_id is required for user accounts", "E1178404").
			Custom(pkg_err.ValidationFailedErr).Message("role_id is required for user accounts").Build()
		return
	}

	if account.Type == account_type.Customer {
		account.RoleID = 0
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
	s.Engine.RedisCacheAPI.Delete(accountsCacheKey(account.Type))

	return
}

// Save a account, if it is exists update it, if not create it
func (s *BaseAccountServ) Save(tx tx.Tx, account base_model.Account) (updatedAccount, accountBefore base_model.Account, err error) {

	if err = validator.ValidateModel(account, base_term.Account, validator.Update); err != nil {
		err = pkg_err.TickValidate(err, "E1158219", pkg_err.ValidationFailed, account)
		return
	}

	if account.Type != "" {
		if accountBefore, err = s.FindByIDAndType(tx, account.ID, account.Type); err != nil {
			pkg_err.Log(err, "E1124905", "can't fetch account by id and type for saving it", account.ID, account.Type)
			return
		}
	} else if accountBefore, err = s.FindByID(tx, account.ID); err != nil {
		pkg_err.Log(err, "E1124906", "can't fetch account by id for saving it", account.ID)
		return
	}

	account.Type = accountBefore.Type

	if account.Type == account_type.User && account.RoleID == 0 {
		err = pkg_err.New("role_id is required for user accounts", "E1178406").
			Custom(pkg_err.ValidationFailedErr).Message("role_id is required for user accounts").Build()
		return
	}

	if account.Type == account_type.Customer {
		account.RoleID = 0
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
	s.Engine.RedisCacheAPI.Delete(accountsCacheKey(account.Type))

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
	s.Engine.RedisCacheAPI.Delete(accountsCacheKey(account.Type))

	return
}

// DeleteByType deletes an account only when it matches the given type
func (s *BaseAccountServ) DeleteByType(tx tx.Tx, id uint, accountType pkg_types.Enum) (account base_model.Account, err error) {

	if account, err = s.FindByIDAndType(tx, id, accountType); err != nil {
		pkg_err.Log(err, "E1178407", "account not found for deleting by type")
		return
	}

	if err = s.Repo.Delete(tx, account); err != nil {
		pkg_err.Log(err, "E1178408", "account not deleted")
		return
	}

	key := fmt.Sprintf("%v-%v", base_term.Account, account.ID)
	s.Engine.RedisCacheAPI.Delete(key)
	s.Engine.RedisCacheAPI.Delete(base_term.Accounts)
	s.Engine.RedisCacheAPI.Delete(accountsCacheKey(account.Type))

	return
}
