package base_api

import (
	"GoRestify/domain/base"
	"GoRestify/domain/base/base_model"
	"GoRestify/domain/base/base_term"
	"GoRestify/domain/base/enum/account_type"
	"GoRestify/domain/service"
	"GoRestify/internal/core"
	"fmt"
	"net/http"

	"GoRestify/pkg/pkg_err"
	"GoRestify/pkg/pkg_log"
	"GoRestify/pkg/pkg_terms"
	"GoRestify/pkg/pkg_types"
	"GoRestify/pkg/response"

	"github.com/gin-gonic/gin"
)

// AccountAPI for injecting account service
type AccountAPI struct {
	Service service.BaseAccountServ
	Engine  *core.Engine
}

// ProvideAccountAPI .
func ProvideAccountAPI(c service.BaseAccountServ) AccountAPI {
	return AccountAPI{Service: c, Engine: c.Engine}
}

func (a *AccountAPI) entityTerm(accountType pkg_types.Enum) (singular, plural string) {
	if accountType == account_type.Customer {
		return base_term.Customer, base_term.Customers
	}
	return base_term.User, base_term.Users
}

// ListUsers lists accounts with type=user
func (a *AccountAPI) ListUsers(c *gin.Context) {
	a.listByType(c, account_type.User)
}

// ListCustomers lists accounts with type=customer
func (a *AccountAPI) ListCustomers(c *gin.Context) {
	a.listByType(c, account_type.Customer)
}

func (a *AccountAPI) listByType(c *gin.Context, accountType pkg_types.Enum) {
	resp, params := response.NewParam(c, base_model.AccountTable)
	_, plural := a.entityTerm(accountType)

	data := make(map[string]interface{})
	var err error

	if data["list"], data["count"], err = a.Service.ListByType(params, accountType); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(base.ListAccount)
	resp.Status(http.StatusOK).
		Message(pkg_terms.ListOfV, plural).
		JSON(data)
}

// GetAllUsers returns all accounts with type=user
func (a *AccountAPI) GetAllUsers(c *gin.Context) {
	a.getAllByType(c, account_type.User)
}

// GetAllCustomers returns all accounts with type=customer
func (a *AccountAPI) GetAllCustomers(c *gin.Context) {
	a.getAllByType(c, account_type.Customer)
}

func (a *AccountAPI) getAllByType(c *gin.Context, accountType pkg_types.Enum) {
	resp, params := response.NewParam(c, base_model.AccountTable)
	_, plural := a.entityTerm(accountType)
	var accounts []base_model.Account
	var err error

	if accounts, err = a.Service.GetAllByType(params, accountType); err != nil {
		err = pkg_err.Take(err, "E1112975").Message(pkg_err.SomethingWentWrong).Build()
		resp.Error(err).JSON()
		return
	}

	resp.Record(base.ListAccount)
	resp.Status(http.StatusOK).
		Message(pkg_terms.ListOfV, plural).
		JSON(accounts)
}

// FindUserByID fetches a user account by id
func (a *AccountAPI) FindUserByID(c *gin.Context) {
	a.findByType(c, account_type.User, "userID")
}

// FindCustomerByID fetches a customer account by id
func (a *AccountAPI) FindCustomerByID(c *gin.Context) {
	a.findByType(c, account_type.Customer, "customerID")
}

func (a *AccountAPI) findByType(c *gin.Context, accountType pkg_types.Enum, paramName string) {
	resp, params := response.NewParam(c, base_model.AccountTable)
	singular, _ := a.entityTerm(accountType)
	var err error
	var account base_model.Account
	var id uint

	if id, err = resp.GetID(c.Param(paramName), "E1130152", singular); err != nil {
		return
	}

	if account, err = a.Service.FindByIDAndType(params.Tx, id, accountType); err != nil {
		resp.Error(err).JSON()
		return
	}

	account.Password = ""

	resp.Record(base.ViewAccount)
	resp.Status(http.StatusOK).
		Message(pkg_terms.VInfo, singular).
		JSON(account)
}

// CreateUser creates an account with type=user (requires role_id)
func (a *AccountAPI) CreateUser(c *gin.Context) {
	a.createByType(c, account_type.User)
}

// CreateCustomer creates an account with type=customer
func (a *AccountAPI) CreateCustomer(c *gin.Context) {
	a.createByType(c, account_type.Customer)
}

func (a *AccountAPI) createByType(c *gin.Context, accountType pkg_types.Enum) {
	resp, params := response.NewParam(c, base_model.AccountTable)
	singular, _ := a.entityTerm(accountType)
	var account, createdAccount base_model.Account
	var err error

	if err = resp.Bind(&account, "E1139696", singular); err != nil {
		return
	}

	account.Type = accountType
	if accountType == account_type.Customer {
		account.RoleID = 0
	}

	params.Tx.DB = a.Engine.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			pkg_log.LogError(fmt.Errorf("panic happened in tx mode for %v",
				singular), "rollback recover create account")
			err = pkg_err.New(pkg_err.SomethingWentWrong, "E1128312").
				Message(pkg_err.SomethingWentWrong).Build()
			params.Tx.DB.Rollback()
			return
		}
	}()

	if createdAccount, err = a.Service.Create(params.Tx, account); err != nil {
		resp.Error(err).JSON()
		params.Tx.DB.Rollback()
		return
	}

	params.Tx.DB.Commit()

	resp.Record(base.CreateAccount, account, createdAccount)
	resp.Status(http.StatusOK).
		Message(pkg_terms.VCreatedSuccessfully, singular).
		JSON(createdAccount)
}

// UpdateUser updates an account with type=user
func (a *AccountAPI) UpdateUser(c *gin.Context) {
	a.updateByType(c, account_type.User, "userID")
}

// UpdateCustomer updates an account with type=customer
func (a *AccountAPI) UpdateCustomer(c *gin.Context) {
	a.updateByType(c, account_type.Customer, "customerID")
}

func (a *AccountAPI) updateByType(c *gin.Context, accountType pkg_types.Enum, paramName string) {
	resp, params := response.NewParam(c, base_model.AccountTable)
	singular, _ := a.entityTerm(accountType)
	var err error
	var account, accountBefore, accountUpdated base_model.Account

	if err = resp.Bind(&account, "E1126876", singular); err != nil {
		return
	}

	if account.ID, err = resp.GetID(c.Param(paramName), "E1182300", singular); err != nil {
		return
	}

	account.Type = accountType
	if accountType == account_type.Customer {
		account.RoleID = 0
	}

	params.Tx.DB = a.Engine.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			pkg_log.LogError(fmt.Errorf("panic happened in tx mode for %v",
				singular), "rollback recover update account")
			err = pkg_err.New(pkg_err.SomethingWentWrong, "E1163983").
				Message(pkg_err.SomethingWentWrong).Build()
			params.Tx.DB.Rollback()
			return
		}
	}()

	if accountUpdated, accountBefore, err = a.Service.Save(params.Tx, account); err != nil {
		resp.Error(err).JSON()
		params.Tx.DB.Rollback()
		return
	}

	params.Tx.DB.Commit()

	resp.Record(base.UpdateAccount, accountBefore, account)
	resp.Status(http.StatusOK).
		Message(pkg_terms.VUpdatedSuccessfully, singular).
		JSON(accountUpdated)
}

// DeleteUser deletes an account with type=user
func (a *AccountAPI) DeleteUser(c *gin.Context) {
	a.deleteByType(c, account_type.User, "userID")
}

// DeleteCustomer deletes an account with type=customer
func (a *AccountAPI) DeleteCustomer(c *gin.Context) {
	a.deleteByType(c, account_type.Customer, "customerID")
}

func (a *AccountAPI) deleteByType(c *gin.Context, accountType pkg_types.Enum, paramName string) {
	resp, params := response.NewParam(c, base_model.AccountTable)
	singular, _ := a.entityTerm(accountType)
	var err error
	var account base_model.Account
	var id uint

	if id, err = resp.GetID(c.Param(paramName), "E1148834", singular); err != nil {
		return
	}

	params.Tx.DB = a.Engine.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			pkg_log.LogError(fmt.Errorf("panic happened in tx mode for %v",
				singular), "rollback recover delete account")
			err = pkg_err.New(pkg_err.SomethingWentWrong, "E1119314").
				Message(pkg_err.SomethingWentWrong).Build()
			params.Tx.DB.Rollback()
			return
		}
	}()

	if account, err = a.Service.DeleteByType(params.Tx, id, accountType); err != nil {
		resp.Error(err).JSON()
		params.Tx.DB.Rollback()
		return
	}

	params.Tx.DB.Commit()

	resp.Record(base.DeleteAccount, account)
	resp.Status(http.StatusOK).
		Message(pkg_terms.VDeletedSuccessfully, singular).
		JSON()
}
