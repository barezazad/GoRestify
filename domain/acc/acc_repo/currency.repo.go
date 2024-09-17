package acc_repo

import (
	"GoRestify/domain/acc/acc_model"
	"GoRestify/domain/acc/acc_term"
	"GoRestify/internal/core"
	"reflect"

	"GoRestify/pkg/db_error"
	"GoRestify/pkg/param"
	"GoRestify/pkg/pkg_err"
	"GoRestify/pkg/tx"

	"GoRestify/pkg/validator"

	"GoRestify/pkg/pkg_sql"
)

// CurrencyRepo for injecting engine
type CurrencyRepo struct {
	Engine *core.Engine
	Cols   []string
}

// ProvideCurrencyRepo is used in wire and initiate the Cols
func ProvideCurrencyRepo(engine *core.Engine) CurrencyRepo {
	return CurrencyRepo{
		Engine: engine,
		Cols:   pkg_sql.ColumnExtractor(reflect.TypeOf(acc_model.Currency{}), acc_model.CurrencyTable),
	}
}

// FindByID finds the currency via its id
func (r *CurrencyRepo) FindByID(tx tx.Tx, id uint) (currency acc_model.Currency, err error) {
	err = tx.GetDB(r.Engine.DB, true).Table(acc_model.CurrencyTable).
		Where("id = ?", id).
		First(&currency).Error

	err = db_error.Parse(err, acc_term.Currencies, validator.Find)
	return
}

// List of currencies
func (r *CurrencyRepo) List(params param.Param) (currencies []acc_model.Currency, err error) {

	var colsStr string
	if colsStr, err = validator.CheckColumns(r.Cols, params.Select); err != nil {
		err = pkg_err.Take(err, "E1171238").Build()
		return
	}

	var whereStr string
	if whereStr, err = params.ParseWhere(r.Cols); err != nil {
		err = pkg_err.Take(err, "E1148650").Custom(pkg_err.ValidationFailedErr).Build()
		return
	}

	err = r.Engine.DB.Table(acc_model.CurrencyTable).Select(colsStr).
		Where(whereStr).
		Order(params.Order).
		Limit(params.Limit).
		Offset(params.Offset).
		Find(&currencies).Error

	err = db_error.Parse(err, acc_term.Currencies, validator.List)
	return
}

// Count of currencies
func (r *CurrencyRepo) Count(params param.Param) (count int64, err error) {
	var whereStr string
	if whereStr, err = params.ParseWhere(r.Cols); err != nil {
		err = pkg_err.Take(err, "E1153286").Custom(pkg_err.ValidationFailedErr).Build()
		return
	}

	err = r.Engine.DB.Table(acc_model.CurrencyTable).
		Where(whereStr).
		Count(&count).Error

	err = db_error.Parse(err, acc_term.Currencies, validator.List)
	return
}

// Create is used for creating currency in tx mode
func (r *CurrencyRepo) Create(tx tx.Tx, currency acc_model.Currency) (u acc_model.Currency, err error) {
	err = tx.GetDB(r.Engine.DB).Table(acc_model.CurrencyTable).Create(&currency).Scan(&u).Error

	err = db_error.Parse(err, acc_term.Currencies, validator.Create)
	return
}

// Save CurrencyRepo
func (r *CurrencyRepo) Save(tx tx.Tx, currency acc_model.Currency) (u acc_model.Currency, err error) {
	err = tx.GetDB(r.Engine.DB).Table(acc_model.CurrencyTable).Save(&currency).Find(&u).Error

	err = db_error.Parse(err, acc_term.Currencies, validator.Update)
	return
}

// Delete currency
func (r *CurrencyRepo) Delete(tx tx.Tx, currency acc_model.Currency) (err error) {
	err = tx.GetDB(r.Engine.DB).Table(acc_model.CurrencyTable).Unscoped().Delete(&currency).Error

	err = db_error.Parse(err, acc_term.Currencies, validator.Delete)
	return
}
