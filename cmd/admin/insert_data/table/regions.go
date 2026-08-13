package table

import (
	"GoRestify/domain/base/base_model"
	"GoRestify/domain/service"
	"GoRestify/internal/core"

	"GoRestify/pkg/pkg_log"
	"GoRestify/pkg/tx"
)

// InsertRegions for add required regions
func InsertRegions(engine *core.Engine) {
	regions := []base_model.Region{
		{ID: 1, Name: "Kurdistan"},
	}

	for _, v := range regions {
		if _, err := service.BaseRegionService.FindByID(tx.Tx{}, v.ID); err == nil {
			if _, _, err = service.BaseRegionService.Save(tx.Tx{}, v); err != nil {
				pkg_log.Fatal("error in saving regions", err)
			}
		} else {
			if _, err = service.BaseRegionService.Create(tx.Tx{}, v); err != nil {
				pkg_log.Fatal("error in creating regions", err)
			}
		}
	}

}
