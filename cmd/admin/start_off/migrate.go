package start_off

import (
	"GoRestify/domain/base/base_model"
	"GoRestify/internal/core"
	"log"

	"GoRestify/pkg/pkg_sql"
)

func onMigrate(err error) {
	if err != nil {
		log.Fatalf("an error occurred while migrate: %v\n", err)
	}
}

// Migrate the database for creating tables
func Migrate(engine *core.Engine) {

	// Base Domain
	onMigrate(engine.DB.Table(base_model.RegionTable).AutoMigrate(&base_model.Region{}))

	onMigrate(engine.DB.Table(base_model.CityTable).AutoMigrate(&base_model.City{}))
	engine.DB.Exec(pkg_sql.ForeignKey(base_model.CityTable, base_model.RegionTable, "region_id"))
}
