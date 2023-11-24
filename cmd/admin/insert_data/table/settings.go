package table

import (
	"GoRestify/internal/core/setting_key"

	"GoRestify/pkg/dictionary"
	"GoRestify/pkg/models"
	"GoRestify/pkg/pkg_log"
	"GoRestify/pkg/setting"
)

// InsertSettings for add required settings
func InsertSettings() {
	settings := []models.Setting{
		{
			ID:          1,
			Property:    setting_key.DefaultLang,
			Value:       string(dictionary.En),
			Description: "default language example for setting env",
		},
	}

	for _, v := range settings {

		if _, err := setting.FindByID(v.ID); err == nil {
			if _, err = setting.Save(v); err != nil {
				pkg_log.Fatal("error in saving settings", err)
			}
		} else {
			if _, err = setting.Create(v); err != nil {
				pkg_log.Fatal("error in creating settings", err)
			}
		}
	}

}
