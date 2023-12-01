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

// RoleRepo for injecting engine
type RoleRepo struct {
	Engine *core.Engine
	Cols   []string
}

// ProvideRoleRepo is used in wire and initiate the Cols
func ProvideRoleRepo(engine *core.Engine) RoleRepo {
	return RoleRepo{
		Engine: engine,
		Cols:   pkg_sql.ColumnExtractor(reflect.TypeOf(base_model.Role{}), base_model.RoleTable),
	}
}

// FindByID finds the role via its id
func (r *RoleRepo) FindByID(tx tx.Tx, id uint) (role base_model.Role, err error) {
	err = tx.GetDB(r.Engine.DB, true).Table(base_model.RoleTable).
		Where("id = ?", id).
		First(&role).Error

	err = db_error.Parse(err, base_term.Roles, validator.Find)
	return
}

// List of roles
func (r *RoleRepo) List(params param.Param) (roles []base_model.Role, err error) {

	var colsStr string
	if colsStr, err = validator.CheckColumns(r.Cols, params.Select); err != nil {
		err = pkg_err.Take(err, "E1139556").Build()
		return
	}

	var whereStr string
	if whereStr, err = params.ParseWhere(r.Cols); err != nil {
		err = pkg_err.Take(err, "E1173432").Custom(pkg_err.ValidationFailedErr).Build()
		return
	}

	err = r.Engine.DB.Table(base_model.RoleTable).Select(colsStr).
		Where(whereStr).
		Order(params.Order).
		Limit(params.Limit).
		Offset(params.Offset).
		Find(&roles).Error

	err = db_error.Parse(err, base_term.Roles, validator.List)
	return
}

// Count of roles
func (r *RoleRepo) Count(params param.Param) (count int64, err error) {
	var whereStr string
	if whereStr, err = params.ParseWhere(r.Cols); err != nil {
		err = pkg_err.Take(err, "E1129491").Custom(pkg_err.ValidationFailedErr).Build()
		return
	}

	err = r.Engine.DB.Table(base_model.RoleTable).
		Where(whereStr).
		Count(&count).Error

	err = db_error.Parse(err, base_term.Roles, validator.List)
	return
}

// Create is used for creating role in tx mode
func (r *RoleRepo) Create(tx tx.Tx, role base_model.Role) (u base_model.Role, err error) {
	err = tx.GetDB(r.Engine.DB).Table(base_model.RoleTable).Create(&role).Scan(&u).Error

	err = db_error.Parse(err, base_term.Roles, validator.Create)
	return
}

// Save RoleRepo
func (r *RoleRepo) Save(tx tx.Tx, role base_model.Role) (u base_model.Role, err error) {
	err = tx.GetDB(r.Engine.DB).Table(base_model.RoleTable).Save(&role).Find(&u).Error

	err = db_error.Parse(err, base_term.Roles, validator.Update)
	return
}

// Delete role
func (r *RoleRepo) Delete(tx tx.Tx, role base_model.Role) (err error) {
	err = tx.GetDB(r.Engine.DB).Table(base_model.RoleTable).Unscoped().Delete(&role).Error

	err = db_error.Parse(err, base_term.Roles, validator.Delete)
	return
}
