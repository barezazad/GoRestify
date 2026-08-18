package base_repo

import (
	"GoRestify/domain/base/base_model"
	"GoRestify/domain/base/base_term"
	"GoRestify/internal/core"
	"reflect"

	"GoRestify/pkg/db_error"
	"GoRestify/pkg/param"
	"GoRestify/pkg/pkg_err"
	"GoRestify/pkg/tx"

	"GoRestify/pkg/validator"

	"GoRestify/pkg/pkg_sql"
)

// AccountRepo for injecting engine
type AccountRepo struct {
	Engine *core.Engine
	Cols   []string
}

// ProvideAccountRepo is used in wire and initiate the Cols
func ProvideAccountRepo(engine *core.Engine) AccountRepo {
	return AccountRepo{
		Engine: engine,
		Cols:   pkg_sql.ColumnExtractor(reflect.TypeOf(base_model.Account{}), base_model.AccountTable),
	}
}

// FindByID finds the account via its id
func (r *AccountRepo) FindByID(tx tx.Tx, id uint) (account base_model.Account, err error) {
	err = tx.GetDB(r.Engine.DB, true).Table(base_model.AccountTable).
		Where("id = ?", id).
		First(&account).Error

	err = db_error.Parse(err, base_term.Accounts, validator.Find)
	return
}

// FindByIDAndType finds the account via its id and type
func (r *AccountRepo) FindByIDAndType(tx tx.Tx, id uint, accountType string) (account base_model.Account, err error) {
	err = tx.GetDB(r.Engine.DB, true).Table(base_model.AccountTable).
		Where("id = ? AND type = ?", id, accountType).
		First(&account).Error

	err = db_error.Parse(err, base_term.Accounts, validator.Find)
	return
}

// GetAll returns all accounts without filter, order, or pagination
func (r *AccountRepo) GetAll() (accounts []base_model.Account, err error) {
	err = r.Engine.DB.Table(base_model.AccountTable).Find(&accounts).Error
	err = db_error.Parse(err, base_term.Accounts, validator.List)
	return
}

// GetAllByType returns all accounts of a given type without filter, order, or pagination
func (r *AccountRepo) GetAllByType(accountType string) (accounts []base_model.Account, err error) {
	err = r.Engine.DB.Table(base_model.AccountTable).
		Where("type = ?", accountType).
		Find(&accounts).Error
	err = db_error.Parse(err, base_term.Accounts, validator.List)
	return
}

// FindByUsername finds the account via its username
func (r *AccountRepo) FindByUsername(tx tx.Tx, username string) (account base_model.Account, err error) {

	var colsStr string
	if colsStr, err = validator.CheckColumns(r.Cols, "*"); err != nil {
		err = pkg_err.Take(err, "E1126008").Build()
		return
	}

	err = tx.GetDB(r.Engine.DB, true).Table(base_model.AccountTable).Select(colsStr).
		Where("base_accounts.username = ?", username).First(&account).Error

	err = db_error.Parse(err, base_term.Accounts, validator.Find)
	return
}

// List of accounts
func (r *AccountRepo) List(params param.Param) (accounts []base_model.Account, err error) {

	var colsStr string
	if colsStr, err = validator.CheckColumns(r.Cols, params.Select); err != nil {
		err = pkg_err.Take(err, "E1194599").Build()
		return
	}

	var whereStr string
	var whereArgs []interface{}
	if whereStr, whereArgs, err = params.ParseWhere(r.Cols); err != nil {
		err = pkg_err.Take(err, "E1123304").Custom(pkg_err.ValidationFailedErr).Build()
		return
	}

	err = r.Engine.DB.Table(base_model.AccountTable).Select(colsStr).
		Where(whereStr, whereArgs...).
		Order(params.Order).
		Limit(params.Limit).
		Offset(params.Offset).
		Find(&accounts).Error

	err = db_error.Parse(err, base_term.Accounts, validator.List)
	return
}

// Count of accounts
func (r *AccountRepo) Count(params param.Param) (count int64, err error) {
	var whereStr string
	var whereArgs []interface{}
	if whereStr, whereArgs, err = params.ParseWhere(r.Cols); err != nil {
		err = pkg_err.Take(err, "E1184792").Custom(pkg_err.ValidationFailedErr).Build()
		return
	}

	err = r.Engine.DB.Table(base_model.AccountTable).
		Where(whereStr, whereArgs...).
		Count(&count).Error

	err = db_error.Parse(err, base_term.Accounts, validator.List)
	return
}

// Create is used for creating account in tx mode
func (r *AccountRepo) Create(tx tx.Tx, account base_model.Account) (u base_model.Account, err error) {
	err = tx.GetDB(r.Engine.DB).Table(base_model.AccountTable).Create(&account).Scan(&u).Error

	err = db_error.Parse(err, base_term.Accounts, validator.Create)
	return
}

// Save AccountRepo
func (r *AccountRepo) Save(tx tx.Tx, account base_model.Account) (u base_model.Account, err error) {
	err = tx.GetDB(r.Engine.DB).Table(base_model.AccountTable).Save(&account).Find(&u).Error

	err = db_error.Parse(err, base_term.Accounts, validator.Update)
	return
}

// Delete account
func (r *AccountRepo) Delete(tx tx.Tx, account base_model.Account) (err error) {
	err = tx.GetDB(r.Engine.DB).Table(base_model.AccountTable).Unscoped().Delete(&account).Error

	err = db_error.Parse(err, base_term.Accounts, validator.Delete)
	return
}
