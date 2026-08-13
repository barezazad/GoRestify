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

const transactionCurrencyJoin = "INNER JOIN acc_currencies ON base_transactions.currency_id = acc_currencies.id"

// FindByID finds the transaction via its id
func (r *TransactionRepo) FindByID(tx tx.Tx, id uint) (transaction acc_model.Transaction, err error) {
	err = tx.GetDB(r.Engine.DB, true).Table(acc_model.TransactionTable).
		Joins(transactionCurrencyJoin).
		Where("base_transactions.id = ?", id).
		First(&transaction).Error

	err = db_error.Parse(err, acc_term.Transactions, validator.Find)
	return
}

// FindByHash finds a transaction by idempotency hash
func (r *TransactionRepo) FindByHash(tx tx.Tx, hash string) (transaction acc_model.Transaction, err error) {
	err = tx.GetDB(r.Engine.DB, true).Table(acc_model.TransactionTable).
		Joins(transactionCurrencyJoin).
		Where("base_transactions.hash = ?", hash).
		First(&transaction).Error

	err = db_error.Parse(err, acc_term.Transactions, validator.Find)
	return
}

// GetAll returns all transactions without filter, order, or pagination
func (r *TransactionRepo) GetAll() (transactions []acc_model.Transaction, err error) {
	err = r.Engine.DB.Table(acc_model.TransactionTable).
		Joins(transactionCurrencyJoin).
		Find(&transactions).Error
	err = db_error.Parse(err, acc_term.Transactions, validator.List)
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
	var whereArgs []interface{}
	if whereStr, whereArgs, err = params.ParseWhere(r.Cols); err != nil {
		err = pkg_err.Take(err, "E1151756").Custom(pkg_err.ValidationFailedErr).Build()
		return
	}

	err = r.Engine.DB.Table(acc_model.TransactionTable).Select(colsStr).
		Joins(transactionCurrencyJoin).
		Where(whereStr, whereArgs...).
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
	var whereArgs []interface{}
	if whereStr, whereArgs, err = params.ParseWhere(r.Cols); err != nil {
		err = pkg_err.Take(err, "E1143725").Custom(pkg_err.ValidationFailedErr).Build()
		return
	}

	err = r.Engine.DB.Table(acc_model.TransactionTable).
		Joins(transactionCurrencyJoin).
		Where(whereStr, whereArgs...).
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

// Save TransactionRepo — prefer immutability at service layer; kept for repair tools.
func (r *TransactionRepo) Save(tx tx.Tx, transaction acc_model.Transaction) (u acc_model.Transaction, err error) {
	err = tx.GetDB(r.Engine.DB).Table(acc_model.TransactionTable).Save(&transaction).Find(&u).Error

	err = db_error.Parse(err, acc_term.Transactions, validator.Update)
	return
}

// Delete transaction — hard delete; prefer immutability at service layer.
func (r *TransactionRepo) Delete(tx tx.Tx, transaction acc_model.Transaction) (err error) {
	err = tx.GetDB(r.Engine.DB).Table(acc_model.TransactionTable).Unscoped().Delete(&transaction).Error

	err = db_error.Parse(err, acc_term.Transactions, validator.Delete)
	return
}
