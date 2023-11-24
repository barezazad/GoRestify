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

// RegionRepo for injecting engine
type RegionRepo struct {
	Engine *core.Engine
	Cols   []string
}

// ProvideRegionRepo is used in wire and initiate the Cols
func ProvideRegionRepo(engine *core.Engine) RegionRepo {
	return RegionRepo{
		Engine: engine,
		Cols:   pkg_sql.ColumnExtractor(reflect.TypeOf(base_model.Region{}), base_model.RegionTable),
	}
}

// FindByID finds the region via its id
func (r *RegionRepo) FindByID(tx tx.Tx, id uint) (region base_model.Region, err error) {
	err = tx.GetDB(r.Engine.DB, true).Table(base_model.RegionTable).
		Where("id = ?", id).
		First(&region).Error

	err = db_error.Parse(err, base_term.Regions, validator.Find)
	return
}

// List of regions
func (r *RegionRepo) List(params param.Param) (regions []base_model.Region, err error) {

	var colsStr string
	if colsStr, err = validator.CheckColumns(r.Cols, params.Select); err != nil {
		err = pkg_err.Take(err, "E1177123").Build()
		return
	}

	var whereStr string
	if whereStr, err = params.ParseWhere(r.Cols); err != nil {
		err = pkg_err.Take(err, "E1174655").Custom(pkg_err.ValidationFailedErr).Build()
		return
	}

	err = r.Engine.DB.Table(base_model.RegionTable).Select(colsStr).
		Where(whereStr).
		Order(params.Order).
		Limit(params.Limit).
		Offset(params.Offset).
		Find(&regions).Error

	err = db_error.Parse(err, base_term.Regions, validator.List)
	return
}

// Count of regions
func (r *RegionRepo) Count(params param.Param) (count int64, err error) {
	var whereStr string
	if whereStr, err = params.ParseWhere(r.Cols); err != nil {
		err = pkg_err.Take(err, "E1124248").Custom(pkg_err.ValidationFailedErr).Build()
		return
	}

	err = r.Engine.DB.Table(base_model.RegionTable).
		Where(whereStr).
		Count(&count).Error

	err = db_error.Parse(err, base_term.Regions, validator.List)
	return
}

// Create is used for creating region in tx mode
func (r *RegionRepo) Create(tx tx.Tx, region base_model.Region) (u base_model.Region, err error) {
	err = tx.GetDB(r.Engine.DB).Table(base_model.RegionTable).Create(&region).Scan(&u).Error

	err = db_error.Parse(err, base_term.Regions, validator.Create)
	return
}

// Save RegionRepo
func (r *RegionRepo) Save(tx tx.Tx, region base_model.Region) (u base_model.Region, err error) {
	err = tx.GetDB(r.Engine.DB).Table(base_model.RegionTable).Save(&region).Find(&u).Error

	err = db_error.Parse(err, base_term.Regions, validator.Update)
	return
}

// Delete region
func (r *RegionRepo) Delete(tx tx.Tx, region base_model.Region) (err error) {
	err = tx.GetDB(r.Engine.DB).Table(base_model.RegionTable).Unscoped().Delete(&region).Error

	err = db_error.Parse(err, base_term.Regions, validator.Delete)
	return
}
