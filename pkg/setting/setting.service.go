package setting

import (
	"GoRestify/pkg/db_error"
	"GoRestify/pkg/models"
	"GoRestify/pkg/param"
	"GoRestify/pkg/pkg_config"
	"GoRestify/pkg/pkg_err"
	"GoRestify/pkg/pkg_terms"
	"GoRestify/pkg/validator"
)

// FindByID finds the setting via its id
func FindByID(id uint) (setting models.Setting, err error) {
	err = pkg_config.Config.DB.Table(models.SettingTable).Where("id = ?", id).First(&setting).Error

	err = db_error.Parse(err, "settings", validator.Find)
	return
}

// List of settings
func List(params param.Param) (settings []models.Setting, err error) {

	var colsStr string
	if colsStr, err = validator.CheckColumns(models.ColumnsSetting, params.Select); err != nil {
		err = pkg_err.Take(err, "E1144106").Build()
		return
	}

	var whereStr string
	if whereStr, err = params.ParseWhere(models.ColumnsSetting); err != nil {
		err = pkg_err.Take(err, "E1146731").Custom(pkg_err.ValidationFailedErr).Build()
		return
	}

	err = pkg_config.Config.DB.Table(models.SettingTable).Select(colsStr).
		Where(whereStr).
		Order(params.Order).
		Limit(params.Limit).
		Offset(params.Offset).
		Find(&settings).Error

	err = db_error.Parse(err, "settings", validator.List)
	return
}

// Count of settings
func Count(params param.Param) (count int64, err error) {
	var whereStr string
	if whereStr, err = params.ParseWhere(models.ColumnsSetting); err != nil {
		err = pkg_err.Take(err, "E1131934").Custom(pkg_err.ValidationFailedErr).Build()
		return
	}

	err = pkg_config.Config.DB.Table(models.SettingTable).
		Where(whereStr).
		Count(&count).Error

	err = db_error.Parse(err, "settings", validator.List)
	return
}

// Create setting
func Create(setting models.Setting) (u models.Setting, err error) {
	err = pkg_config.Config.DB.Table(models.SettingTable).Create(&setting).Scan(&u).Error

	// reset redis
	pkg_config.Config.Redis.Set(pkg_terms.Settings, "")

	err = db_error.Parse(err, "settings", validator.Create)
	return
}

// Save setting
func Save(setting models.Setting) (u models.Setting, err error) {

	// get setting before
	beforeSetting, err := FindByID(setting.ID)
	if err != nil {
		return
	}
	setting.Property = beforeSetting.Property

	// update in mysql
	err = pkg_config.Config.DB.Table(models.SettingTable).Save(&setting).Find(&u).Error

	// reset redis
	pkg_config.Config.Redis.Set(pkg_terms.Settings, "")

	err = db_error.Parse(err, "settings", validator.Update)
	return
}

// Delete setting
func Delete(id uint) (err error) {

	var setting models.Setting
	err = pkg_config.Config.DB.Table(models.SettingTable).Unscoped().Where("id = ?", id).Delete(&setting).Error

	// reset redis
	pkg_config.Config.Redis.Set(pkg_terms.Settings, "")

	err = db_error.Parse(err, "settings", validator.Delete)
	return
}
