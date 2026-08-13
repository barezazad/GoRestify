package acc_repo

import (
	"GoRestify/domain/acc/acc_model"
	"GoRestify/domain/acc/acc_term"
	"GoRestify/internal/core"
	"reflect"

	"GoRestify/pkg/db_error"
	"GoRestify/pkg/param"
	"GoRestify/pkg/pkg_err"
	"GoRestify/pkg/pkg_sql"
	"GoRestify/pkg/tx"
	"GoRestify/pkg/validator"
)

// AccountCreditRepo for injecting engine
type AccountCreditRepo struct {
	Engine *core.Engine
	Cols   []string
}

// ProvideAccountCreditRepo is used in wire and initiate the Cols
func ProvideAccountCreditRepo(engine *core.Engine) AccountCreditRepo {
	return AccountCreditRepo{
		Engine: engine,
		Cols:   pkg_sql.ColumnExtractor(reflect.TypeOf(acc_model.AccountCredit{}), acc_model.AccountCreditTable),
	}
}

const accountCreditCurrencyJoin = "INNER JOIN acc_currencies ON acc_account_credits.currency_id = acc_currencies.id"

// FindByID finds the accountCredit via its id
func (r *AccountCreditRepo) FindByID(tx tx.Tx, id uint) (accountCredit acc_model.AccountCredit, err error) {
	err = tx.GetDB(r.Engine.DB, true).Table(acc_model.AccountCreditTable).
		Joins(accountCreditCurrencyJoin).
		Where("acc_account_credits.id = ?", id).
		First(&accountCredit).Error

	err = db_error.Parse(err, acc_term.AccountCredits, validator.Find)
	return
}

// GetAll returns all accountcredits without filter, order, or pagination
func (r *AccountCreditRepo) GetAll() (accountCredits []acc_model.AccountCredit, err error) {
	err = r.Engine.DB.Table(acc_model.AccountCreditTable).
		Joins(accountCreditCurrencyJoin).
		Find(&accountCredits).Error
	err = db_error.Parse(err, acc_term.AccountCredits, validator.List)
	return
}

// FindByAccountIDAndCurrency finds the accountCredit via account + currency
func (r *AccountCreditRepo) FindByAccountIDAndCurrency(tx tx.Tx, accountID, currencyID uint) (accountCredit acc_model.AccountCredit, err error) {

	err = tx.GetDB(r.Engine.DB, true).Table(acc_model.AccountCreditTable).
		Joins(accountCreditCurrencyJoin).
		Where("acc_account_credits.account_id = ? and acc_account_credits.currency_id = ?", accountID, currencyID).
		First(&accountCredit).Error

	err = db_error.Parse(err, acc_term.AccountCredits, validator.Find)
	return
}

// List of accountCredits
func (r *AccountCreditRepo) List(params param.Param) (accountCredits []acc_model.AccountCredit, err error) {

	var colsStr string
	if colsStr, err = validator.CheckColumns(r.Cols, params.Select); err != nil {
		err = pkg_err.Take(err, "E1177123").Build()
		return
	}

	var whereStr string
	var whereArgs []interface{}
	if whereStr, whereArgs, err = params.ParseWhere(r.Cols); err != nil {
		err = pkg_err.Take(err, "E1174655").Custom(pkg_err.ValidationFailedErr).Build()
		return
	}

	err = r.Engine.DB.Table(acc_model.AccountCreditTable).Select(colsStr).
		Joins(accountCreditCurrencyJoin).
		Where(whereStr, whereArgs...).
		Order(params.Order).
		Limit(params.Limit).
		Offset(params.Offset).
		Find(&accountCredits).Error

	err = db_error.Parse(err, acc_term.AccountCredits, validator.List)
	return
}

// Count of accountCredits
func (r *AccountCreditRepo) Count(params param.Param) (count int64, err error) {
	var whereStr string
	var whereArgs []interface{}
	if whereStr, whereArgs, err = params.ParseWhere(r.Cols); err != nil {
		err = pkg_err.Take(err, "E1124248").Custom(pkg_err.ValidationFailedErr).Build()
		return
	}

	err = r.Engine.DB.Table(acc_model.AccountCreditTable).
		Joins(accountCreditCurrencyJoin).
		Where(whereStr, whereArgs...).
		Count(&count).Error

	err = db_error.Parse(err, acc_term.AccountCredits, validator.List)
	return
}

// CountByCurrency counts credit rows for a currency (delete guard)
func (r *AccountCreditRepo) CountByCurrency(tx tx.Tx, currencyID uint) (count int64, err error) {
	err = tx.GetDB(r.Engine.DB).Table(acc_model.AccountCreditTable).
		Where("currency_id = ?", currencyID).
		Count(&count).Error
	err = db_error.Parse(err, acc_term.AccountCredits, validator.List)
	return
}

// Create is used for creating accountCredit in tx mode
func (r *AccountCreditRepo) Create(tx tx.Tx, accountCredit acc_model.AccountCredit) (u acc_model.AccountCredit, err error) {
	err = tx.GetDB(r.Engine.DB).Table(acc_model.AccountCreditTable).Create(&accountCredit).Scan(&u).Error

	err = db_error.Parse(err, acc_term.AccountCredits, validator.Create)
	return
}

// Save AccountCreditRepo
func (r *AccountCreditRepo) Save(tx tx.Tx, accountCredit acc_model.AccountCredit) (u acc_model.AccountCredit, err error) {
	err = tx.GetDB(r.Engine.DB).Table(acc_model.AccountCreditTable).Save(&accountCredit).Find(&u).Error

	err = db_error.Parse(err, acc_term.AccountCredits, validator.Update)
	return
}

// Delete accountCredit — hard delete; prefer immutability for balances except via DoTransaction.
func (r *AccountCreditRepo) Delete(tx tx.Tx, accountCredit acc_model.AccountCredit) (err error) {
	err = tx.GetDB(r.Engine.DB).Table(acc_model.AccountCreditTable).Unscoped().Delete(&accountCredit).Error

	err = db_error.Parse(err, acc_term.AccountCredits, validator.Delete)
	return
}
