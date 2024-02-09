package base_repo

import (
	"GoRestify/domain/base/base_model"
	"GoRestify/domain/base/base_term"
	"GoRestify/internal/core"
	"reflect"
	"strings"

	"GoRestify/pkg/db_error"
	"GoRestify/pkg/param"
	"GoRestify/pkg/pkg_err"
	"GoRestify/pkg/tx"

	"GoRestify/pkg/validator"

	"GoRestify/pkg/pkg_sql"
)

// UserRepo for injecting engine
type UserRepo struct {
	Engine *core.Engine
	Cols   []string
}

// ProvideUserRepo is used in wire and initiate the Cols
func ProvideUserRepo(engine *core.Engine) UserRepo {
	return UserRepo{
		Engine: engine,
		Cols:   pkg_sql.ColumnExtractor(reflect.TypeOf(base_model.User{}), base_model.UserTable),
	}
}

// FindByID finds the user via its id
func (r *UserRepo) FindByID(tx tx.Tx, id uint) (user base_model.User, err error) {

	var colsStr string
	if colsStr, err = validator.CheckColumns(r.Cols, "*"); err != nil {
		err = pkg_err.Take(err, "E1164621").Build()
		return
	}

	err = tx.GetDB(r.Engine.DB, true).Table(base_model.UserTable).Select(colsStr).
		Joins("INNER JOIN base_roles ON base_users.role_id = base_roles.id").
		Where("base_users.id = ?", id).
		First(&user).Error

	err = db_error.Parse(err, base_term.Users, validator.Find)
	return
}

// GetUserResources is used for finding all resources
func (r *UserRepo) GetUserResources(userID uint) (resourceList base_model.ResourceList, err error) {

	err = r.Engine.DB.Table(base_model.UserTable).
		Select("base_roles.resources AS resources").
		Joins("LEFT JOIN base_roles ON base_users.role_id = base_roles.id").
		Where("base_users.id = ?", userID).
		Scan(&resourceList).Error

	resourceList.ResourcesArray = strings.Split(resourceList.Resources, ",")

	err = db_error.Parse(err, base_term.Users, validator.Find)

	return
}

// List of users
func (r *UserRepo) List(params param.Param) (users []base_model.User, err error) {

	var colsStr string
	if colsStr, err = validator.CheckColumns(r.Cols, params.Select); err != nil {
		err = pkg_err.Take(err, "E1190210").Build()
		return
	}

	var whereStr string
	if whereStr, err = params.ParseWhere(r.Cols); err != nil {
		err = pkg_err.Take(err, "E1158754").Custom(pkg_err.ValidationFailedErr).Build()
		return
	}

	err = r.Engine.DB.Table(base_model.UserTable).Select(colsStr).
		Joins("INNER JOIN base_roles ON base_users.role_id = base_roles.id").
		Where(whereStr).
		Order(params.Order).
		Limit(params.Limit).
		Offset(params.Offset).
		Find(&users).Error

	err = db_error.Parse(err, base_term.Users, validator.List)
	return
}

// Count of users
func (r *UserRepo) Count(params param.Param) (count int64, err error) {
	var whereStr string
	if whereStr, err = params.ParseWhere(r.Cols); err != nil {
		err = pkg_err.Take(err, "E1195269").Custom(pkg_err.ValidationFailedErr).Build()
		return
	}

	err = r.Engine.DB.Table(base_model.UserTable).
		Joins("INNER JOIN base_roles ON base_users.role_id = base_roles.id").
		Where(whereStr).
		Count(&count).Error

	err = db_error.Parse(err, base_term.Users, validator.List)
	return
}

// Create is used for creating user in tx mode
func (r *UserRepo) Create(tx tx.Tx, user base_model.User) (u base_model.User, err error) {
	err = tx.GetDB(r.Engine.DB).Table(base_model.UserTable).Create(&user).Scan(&u).Error

	err = db_error.Parse(err, base_term.Users, validator.Create)
	return
}

// Save UserRepo
func (r *UserRepo) Save(tx tx.Tx, user base_model.User) (u base_model.User, err error) {
	err = tx.GetDB(r.Engine.DB).Table(base_model.UserTable).Save(&user).Find(&u).Error

	err = db_error.Parse(err, base_term.Users, validator.Update)
	return
}

// Delete user
func (r *UserRepo) Delete(tx tx.Tx, user base_model.User) (err error) {
	err = tx.GetDB(r.Engine.DB).Table(base_model.UserTable).Unscoped().Delete(&user).Error

	err = db_error.Parse(err, base_term.Users, validator.Delete)
	return
}
