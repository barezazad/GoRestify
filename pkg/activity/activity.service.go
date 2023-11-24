package activity

import (
	"GoRestify/pkg/db_error"
	"GoRestify/pkg/models"
	"GoRestify/pkg/param"
	"GoRestify/pkg/pkg_config"
	"GoRestify/pkg/pkg_err"
	"GoRestify/pkg/validator"
)

// CreateBatch ActivityRepo
func CreateBatch(activities []models.Activity) (u models.Activity, err error) {

	err = pkg_config.Config.ActivityDB.
		Table(models.ActivityTable).
		Create(&activities).Error

	return
}

// List of activities
func List(params param.Param) (activities []models.Activity, err error) {

	var colsStr string
	if colsStr, err = validator.CheckColumns(models.ColumnsActivity, params.Select); err != nil {
		err = pkg_err.Take(err, "E1151963").Build()
		return
	}

	var whereStr string
	if whereStr, err = params.ParseWhere(models.ColumnsActivity); err != nil {
		err = pkg_err.Take(err, "E1122102").Custom(pkg_err.ValidationFailedErr).Build()
		return
	}

	err = pkg_config.Config.ActivityDB.Table(models.ActivityTable).Select(colsStr).
		Where(whereStr).
		Order(params.Order).
		Limit(params.Limit).
		Offset(params.Offset).
		Find(&activities).Error

	err = db_error.Parse(err, "activities", validator.List)
	return
}

// Count of activities
func Count(params param.Param) (count int64, err error) {
	var whereStr string
	if whereStr, err = params.ParseWhere(models.ColumnsActivity); err != nil {
		err = pkg_err.Take(err, "E1127507").Custom(pkg_err.ValidationFailedErr).Build()
		return
	}

	err = pkg_config.Config.ActivityDB.Table(models.ActivityTable).
		Where(whereStr).
		Count(&count).Error

	err = db_error.Parse(err, "activities", validator.List)
	return
}

// Delete activity
func Delete(id uint) (err error) {
	var activity models.Activity
	err = pkg_config.Config.ActivityDB.Table(models.ActivityTable).Unscoped().Where("id = ?", id).Delete(&activity).Error

	err = db_error.Parse(err, "activities", validator.Delete)
	return
}
