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

// SlotRepo for injecting engine
type SlotRepo struct {
	Engine *core.Engine
	Cols   []string
}

// ProvideSlotRepo is used in wire and initiate the Cols
func ProvideSlotRepo(engine *core.Engine) SlotRepo {
	return SlotRepo{
		Engine: engine,
		Cols:   pkg_sql.ColumnExtractor(reflect.TypeOf(acc_model.Slot{}), acc_model.SlotTable),
	}
}

const slotJoins = `
INNER JOIN base_transactions ON acc_slots.transaction_id = base_transactions.id
INNER JOIN acc_currencies ON acc_slots.currency_id = acc_currencies.id`

// FindByID finds the slot via its id
func (r *SlotRepo) FindByID(tx tx.Tx, id uint) (slot acc_model.Slot, err error) {
	err = tx.GetDB(r.Engine.DB, true).Table(acc_model.SlotTable).
		Joins(slotJoins).
		Where("acc_slots.id = ?", id).
		First(&slot).Error

	err = db_error.Parse(err, acc_term.Slots, validator.Find)
	return
}

// ListByTransactionID returns slots for a transaction header
func (r *SlotRepo) ListByTransactionID(tx tx.Tx, transactionID uint) (slots []acc_model.Slot, err error) {
	err = tx.GetDB(r.Engine.DB).Table(acc_model.SlotTable).
		Joins(slotJoins).
		Where("acc_slots.transaction_id = ?", transactionID).
		Order("acc_slots.id asc").
		Find(&slots).Error

	err = db_error.Parse(err, acc_term.Slots, validator.List)
	return
}

// GetAll returns all slots without filter, order, or pagination
func (r *SlotRepo) GetAll() (slots []acc_model.Slot, err error) {
	err = r.Engine.DB.Table(acc_model.SlotTable).
		Joins(slotJoins).
		Find(&slots).Error
	err = db_error.Parse(err, acc_term.Slots, validator.List)
	return
}

// List of slots
func (r *SlotRepo) List(params param.Param) (slots []acc_model.Slot, err error) {

	var colsStr string
	if colsStr, err = validator.CheckColumns(r.Cols, params.Select); err != nil {
		err = pkg_err.Take(err, "E1134139").Build()
		return
	}

	var whereStr string
	var whereArgs []interface{}
	if whereStr, whereArgs, err = params.ParseWhere(r.Cols); err != nil {
		err = pkg_err.Take(err, "E1132412").Custom(pkg_err.ValidationFailedErr).Build()
		return
	}

	err = r.Engine.DB.Table(acc_model.SlotTable).Select(colsStr).
		Joins(slotJoins).
		Where(whereStr, whereArgs...).
		Order(params.Order).
		Limit(params.Limit).
		Offset(params.Offset).
		Find(&slots).Error

	err = db_error.Parse(err, acc_term.Slots, validator.List)
	return
}

// Count of slots
func (r *SlotRepo) Count(params param.Param) (count int64, err error) {
	var whereStr string
	var whereArgs []interface{}
	if whereStr, whereArgs, err = params.ParseWhere(r.Cols); err != nil {
		err = pkg_err.Take(err, "E1171787").Custom(pkg_err.ValidationFailedErr).Build()
		return
	}

	err = r.Engine.DB.Table(acc_model.SlotTable).
		Joins(slotJoins).
		Where(whereStr, whereArgs...).
		Count(&count).Error

	err = db_error.Parse(err, acc_term.Slots, validator.List)
	return
}

// Create is used for creating slot in tx mode
func (r *SlotRepo) Create(tx tx.Tx, slot acc_model.Slot) (u acc_model.Slot, err error) {
	err = tx.GetDB(r.Engine.DB).Table(acc_model.SlotTable).Create(&slot).Scan(&u).Error

	err = db_error.Parse(err, acc_term.Slots, validator.Create)
	return
}

// Save SlotRepo — prefer immutability at service layer.
func (r *SlotRepo) Save(tx tx.Tx, slot acc_model.Slot) (u acc_model.Slot, err error) {
	err = tx.GetDB(r.Engine.DB).Table(acc_model.SlotTable).Save(&slot).Find(&u).Error

	err = db_error.Parse(err, acc_term.Slots, validator.Update)
	return
}

// Delete slot — hard delete; prefer immutability at service layer.
func (r *SlotRepo) Delete(tx tx.Tx, slot acc_model.Slot) (err error) {
	err = tx.GetDB(r.Engine.DB).Table(acc_model.SlotTable).Unscoped().Delete(&slot).Error

	err = db_error.Parse(err, acc_term.Slots, validator.Delete)
	return
}
