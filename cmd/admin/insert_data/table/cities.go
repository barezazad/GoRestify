package table

import (
	"GoRestify/domain/base/base_model"
	"GoRestify/domain/service"
	"GoRestify/internal/core"

	"GoRestify/pkg/pkg_log"
	"GoRestify/pkg/tx"
)

// InsertCities for add required cities
func InsertCities(engine *core.Engine) {

	cities := []base_model.City{
		{ID: 1, RegionID: 1, Name: "Sulaymaniyah"},
	}

	for _, v := range cities {
		if _, err := service.BaseCityService.FindByID(tx.Tx{}, v.ID); err == nil {
			if _, _, err = service.BaseCityService.Save(tx.Tx{}, v); err != nil {
				pkg_log.Fatal("error in saving cities", err)
			}
		} else {
			if _, err = service.BaseCityService.Create(tx.Tx{}, v); err != nil {
				pkg_log.Fatal("error in creating cities", err)
			}
		}
	}

}
