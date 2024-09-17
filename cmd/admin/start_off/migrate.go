package start_off

import (
	"GoRestify/domain/acc/acc_model"
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
	onMigrate(engine.DB.Table(base_model.AccountTable).AutoMigrate(&base_model.Account{}))

	onMigrate(engine.DB.Table(base_model.RoleTable).AutoMigrate(&base_model.Role{}))

	onMigrate(engine.DB.Table(base_model.UserTable).AutoMigrate(&base_model.User{}))
	engine.DB.Exec(pkg_sql.ForeignKey(base_model.UserTable, base_model.RoleTable, "role_id"))

	onMigrate(engine.DB.Table(base_model.RegionTable).AutoMigrate(&base_model.Region{}))

	onMigrate(engine.DB.Table(base_model.CityTable).AutoMigrate(&base_model.City{}))
	engine.DB.Exec(pkg_sql.ForeignKey(base_model.CityTable, base_model.RegionTable, "region_id"))

	// acc domain
	onMigrate(engine.DB.Table(acc_model.CurrencyTable).AutoMigrate(&acc_model.Currency{}))

	onMigrate(engine.DB.Table(acc_model.AccountCreditTable).AutoMigrate(&acc_model.AccountCredit{}))
	engine.DB.Exec(pkg_sql.ForeignKey(acc_model.AccountCreditTable, base_model.AccountTable, "account_id"))
	engine.DB.Exec(pkg_sql.ForeignKey(acc_model.AccountCreditTable, acc_model.CurrencyTable, "currency_id"))

	onMigrate(engine.DB.Table(acc_model.TransactionTable).AutoMigrate(&acc_model.Transaction{}))
	engine.DB.Exec(pkg_sql.ForeignKey(acc_model.TransactionTable, base_model.AccountTable, "sender_id"))
	engine.DB.Exec(pkg_sql.ForeignKey(acc_model.TransactionTable, base_model.AccountTable, "receiver_id"))
	engine.DB.Exec(pkg_sql.ForeignKey(acc_model.TransactionTable, acc_model.CurrencyTable, "currency_id"))

	onMigrate(engine.DB.Table(acc_model.SlotTable).AutoMigrate(&acc_model.Slot{}))
	engine.DB.Exec(pkg_sql.ForeignKey(acc_model.SlotTable, base_model.AccountTable, "account_id"))
	engine.DB.Exec(pkg_sql.ForeignKey(acc_model.SlotTable, acc_model.CurrencyTable, "currency_id"))
	engine.DB.Exec(pkg_sql.ForeignKey(acc_model.SlotTable, acc_model.TransactionTable, "transaction_id"))
}
