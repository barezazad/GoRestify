package base_repo

import (
	"GoRestify/domain/base/base_model"
	"GoRestify/domain/base/base_term"
	"GoRestify/internal/core"
	"reflect"

	"GoRestify/pkg/db_error"
	"GoRestify/pkg/param"
	"GoRestify/pkg/pkg_err"
	"GoRestify/pkg/pkg_sql"

	"GoRestify/pkg/tx"
	"GoRestify/pkg/validator"
)

// CityRepo for injecting engine
type CityRepo struct {
	Engine *core.Engine
	Cols   []string
}

// ProvideCityRepo is used in wire and initiate the Cols
func ProvideCityRepo(engine *core.Engine) CityRepo {
	return CityRepo{
		Engine: engine,
		Cols:   pkg_sql.ColumnExtractor(reflect.TypeOf(base_model.City{}), base_model.CityTable),
	}
}

// FindByID finds the city via its id
func (r *CityRepo) FindByID(tx tx.Tx, id uint) (city base_model.City, err error) {

	err = tx.GetDB(r.Engine.DB, true).Table(base_model.CityTable).Where("id = ?", id).First(&city).Error

	err = db_error.Parse(err, base_term.Cities, validator.Find)
	return
}

// List returns an array of cities
func (r *CityRepo) List(params param.Param) (cities []base_model.City, err error) {
	var colsStr string
	if colsStr, err = validator.CheckColumns(r.Cols, params.Select); err != nil {
		err = pkg_err.Take(err, "E1165529").Build()
		return
	}

	var whereStr string
	if whereStr, err = params.ParseWhere(r.Cols); err != nil {
		err = pkg_err.Take(err, "E1181807").Custom(pkg_err.ValidationFailedErr).Build()
		return
	}

	err = r.Engine.DB.Table(base_model.CityTable).Select(colsStr).
		Where(whereStr).
		Order(params.Order).
		Limit(params.Limit).
		Offset(params.Offset).
		Find(&cities).Error

	err = db_error.Parse(err, base_term.Cities, validator.List)
	return
}

// Count of cities, mainly calls with List
func (r *CityRepo) Count(params param.Param) (count int64, err error) {
	var whereStr string
	if whereStr, err = params.ParseWhere(r.Cols); err != nil {
		err = pkg_err.Take(err, "E1130747").Custom(pkg_err.ValidationFailedErr).Build()
		return
	}

	err = r.Engine.DB.Table(base_model.CityTable).
		Where(whereStr).
		Count(&count).Error

	err = db_error.Parse(err, base_term.Cities, validator.List)
	return
}

// Save the city, in case it is not exist create it
func (r *CityRepo) Save(tx tx.Tx, city base_model.City) (u base_model.City, err error) {
	err = tx.GetDB(r.Engine.DB).Table(base_model.CityTable).Save(&city).Find(&u).Error

	err = db_error.Parse(err, base_term.Cities, validator.Update)
	return
}

// Create a city
func (r *CityRepo) Create(tx tx.Tx, city base_model.City) (u base_model.City, err error) {
	err = tx.GetDB(r.Engine.DB).Table(base_model.CityTable).Create(&city).Scan(&u).Error

	tx.GetDB(r.Engine.DB)
	err = db_error.Parse(err, base_term.Cities, validator.Create)
	return
}

// Delete the city
func (r *CityRepo) Delete(tx tx.Tx, city base_model.City) (err error) {
	err = tx.GetDB(r.Engine.DB).Unscoped().Table(base_model.CityTable).Delete(&city).Error

	err = db_error.Parse(err, base_term.Cities, validator.Delete)
	return
}
