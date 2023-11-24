package pkg_sql

import "fmt"

// ForeignKey create foreignKey query for table
func ForeignKey(table, refTable, field string, refField ...string) (query string) {

	var referenceField string
	if len(refField) > 0 {
		referenceField = refField[0]
	} else {
		referenceField = "id"
	}

	query = fmt.Sprintf("ALTER TABLE %v ADD CONSTRAINT `fk_%[1]v_%v_%v` ", table, refTable, field)
	query += fmt.Sprintf("FOREIGN KEY (%v) REFERENCES %v(%v) ", field, refTable, referenceField)
	query += fmt.Sprintf("ON DELETE RESTRICT ON UPDATE RESTRICT;")

	return
}
