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

// FindByID finds the slot via its id
func (r *SlotRepo) FindByID(tx tx.Tx, id uint) (slot acc_model.Slot, err error) {
	err = tx.GetDB(r.Engine.DB, true).Table(acc_model.SlotTable).
		Where("id = ?", id).
		First(&slot).Error

	err = db_error.Parse(err, acc_term.Slots, validator.Find)
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
	if whereStr, err = params.ParseWhere(r.Cols); err != nil {
		err = pkg_err.Take(err, "E1132412").Custom(pkg_err.ValidationFailedErr).Build()
		return
	}

	err = r.Engine.DB.Table(acc_model.SlotTable).Select(colsStr).
		Where(whereStr).
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
	if whereStr, err = params.ParseWhere(r.Cols); err != nil {
		err = pkg_err.Take(err, "E1171787").Custom(pkg_err.ValidationFailedErr).Build()
		return
	}

	err = r.Engine.DB.Table(acc_model.SlotTable).
		Where(whereStr).
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

// Save SlotRepo
func (r *SlotRepo) Save(tx tx.Tx, slot acc_model.Slot) (u acc_model.Slot, err error) {
	err = tx.GetDB(r.Engine.DB).Table(acc_model.SlotTable).Save(&slot).Find(&u).Error

	err = db_error.Parse(err, acc_term.Slots, validator.Update)
	return
}

// Delete slot
func (r *SlotRepo) Delete(tx tx.Tx, slot acc_model.Slot) (err error) {
	err = tx.GetDB(r.Engine.DB).Table(acc_model.SlotTable).Unscoped().Delete(&slot).Error

	err = db_error.Parse(err, acc_term.Slots, validator.Delete)
	return
}
