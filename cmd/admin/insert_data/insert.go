package insert_data

import (
	"GoRestify/cmd/admin/insert_data/table"
	"GoRestify/internal/core"
	"GoRestify/internal/core/enum/domain_app"

	"GoRestify/pkg/pkg_log"
	"GoRestify/pkg/setting"
)

// Insert is used for add static rows to database
func Insert(engine *core.Engine) {

	if err := engine.RedisCacheAPI.FlushDB(); err != nil {
		pkg_log.Fatal(err, "couldn't flush redis cache")
	}

	table.InsertSettings()
	setting.LoadSetting(domain_app.Admin)

	table.InsertRoles(engine)
	table.InsertAccounts(engine)
	table.InsertRegions(engine)
	table.InsertCities(engine)

}
