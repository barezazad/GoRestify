package tx

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Tx .
type Tx struct {
	DB     *gorm.DB // to handle transaction for repo
	SecDB  *gorm.DB // to handle transaction for repo
	IsLock bool     // to lock a row in database
}

// GetDB it check for getting db connection in tx or in engine
func (tx *Tx) GetDB(db *gorm.DB, enableRowLocking ...bool) *gorm.DB {

	IsEnableRowLocking := false
	if len(enableRowLocking) > 0 {
		IsEnableRowLocking = enableRowLocking[0]
	}

	if tx.DB != nil {
		if tx.IsLock && IsEnableRowLocking {
			return tx.DB.Clauses(clause.Locking{Strength: "UPDATE"})
		}
		return tx.DB
	}
	return db
}

// GetSecDB it check for getting db connection in tx or in engine
func (tx *Tx) GetSecDB(db *gorm.DB, enableRowLocking ...bool) *gorm.DB {

	IsEnableRowLocking := false
	if len(enableRowLocking) > 0 {
		IsEnableRowLocking = enableRowLocking[0]
	}

	if tx.SecDB != nil {
		if tx.IsLock && IsEnableRowLocking {
			return tx.SecDB.Clauses(clause.Locking{Strength: "UPDATE"})
		}
		return tx.SecDB
	}
	return db
}
