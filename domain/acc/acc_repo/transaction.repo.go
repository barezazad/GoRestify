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

// TransactionRepo for injecting engine
type TransactionRepo struct {
	Engine *core.Engine
	Cols   []string
}

// ProvideTransactionRepo is used in wire and initiate the Cols
func ProvideTransactionRepo(engine *core.Engine) TransactionRepo {
	return TransactionRepo{
		Engine: engine,
		Cols:   pkg_sql.ColumnExtractor(reflect.TypeOf(acc_model.Transaction{}), acc_model.TransactionTable),
	}
}

// FindByID finds the transaction via its id
func (r *TransactionRepo) FindByID(tx tx.Tx, id uint) (transaction acc_model.Transaction, err error) {
	err = tx.GetDB(r.Engine.DB, true).Table(acc_model.TransactionTable).
		Joins("INNER JOIN acc_currencies ON base_transactions.currency_id = acc_currencies.id").
		Where("id = ?", id).
		First(&transaction).Error

	err = db_error.Parse(err, acc_term.Transactions, validator.Find)
	return
}

// List of transactions
func (r *TransactionRepo) List(params param.Param) (transactions []acc_model.Transaction, err error) {

	var colsStr string
	if colsStr, err = validator.CheckColumns(r.Cols, params.Select); err != nil {
		err = pkg_err.Take(err, "E1169452").Build()
		return
	}

	var whereStr string
	if whereStr, err = params.ParseWhere(r.Cols); err != nil {
		err = pkg_err.Take(err, "E1151756").Custom(pkg_err.ValidationFailedErr).Build()
		return
	}

	err = r.Engine.DB.Table(acc_model.TransactionTable).Select(colsStr).
		Joins("INNER JOIN acc_currencies ON base_transactions.currency_id = acc_currencies.id").
		Where(whereStr).
		Order(params.Order).
		Limit(params.Limit).
		Offset(params.Offset).
		Find(&transactions).Error

	err = db_error.Parse(err, acc_term.Transactions, validator.List)
	return
}

// Count of transactions
func (r *TransactionRepo) Count(params param.Param) (count int64, err error) {
	var whereStr string
	if whereStr, err = params.ParseWhere(r.Cols); err != nil {
		err = pkg_err.Take(err, "E1143725").Custom(pkg_err.ValidationFailedErr).Build()
		return
	}

	err = r.Engine.DB.Table(acc_model.TransactionTable).
		Joins("INNER JOIN acc_currencies ON base_transactions.currency_id = acc_currencies.id").
		Where(whereStr).
		Count(&count).Error

	err = db_error.Parse(err, acc_term.Transactions, validator.List)
	return
}

// Create is used for creating transaction in tx mode
func (r *TransactionRepo) Create(tx tx.Tx, transaction acc_model.Transaction) (u acc_model.Transaction, err error) {
	err = tx.GetDB(r.Engine.DB).Table(acc_model.TransactionTable).Create(&transaction).Scan(&u).Error

	err = db_error.Parse(err, acc_term.Transactions, validator.Create)
	return
}

// Save TransactionRepo
func (r *TransactionRepo) Save(tx tx.Tx, transaction acc_model.Transaction) (u acc_model.Transaction, err error) {
	err = tx.GetDB(r.Engine.DB).Table(acc_model.TransactionTable).Save(&transaction).Find(&u).Error

	err = db_error.Parse(err, acc_term.Transactions, validator.Update)
	return
}

// Delete transaction
func (r *TransactionRepo) Delete(tx tx.Tx, transaction acc_model.Transaction) (err error) {
	err = tx.GetDB(r.Engine.DB).Table(acc_model.TransactionTable).Unscoped().Delete(&transaction).Error

	err = db_error.Parse(err, acc_term.Transactions, validator.Delete)
	return
}
