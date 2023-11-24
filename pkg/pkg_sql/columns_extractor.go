package pkg_sql

import (
	"fmt"
	"reflect"
	"regexp"
)

type extractor struct {
	arr   []string
	re    *regexp.Regexp
	table string
}

// ColumnExtractor extract the name of table and field from json and table tag
func ColumnExtractor(t reflect.Type, table string) []string {
	ext := extractor{
		arr:   []string{},
		re:    regexp.MustCompile(`\w+`),
		table: table,
	}

	ext.getTag(t)

	return ext.arr
}

func (p *extractor) getTag(t reflect.Type) {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		externalGorm := field.Tag.Get("gorm")
		if externalGorm == "-" {
			continue
		}

		externalTable := field.Tag.Get("table")

		if externalTable == "" && externalGorm == "-" {
			continue
		}

		if field.Type.Kind() == reflect.Struct && field.Type.Name() == "Model" {
			p.arr = append(p.arr, p.table+"."+"id")
			p.arr = append(p.arr, p.table+"."+"created_at")
			p.arr = append(p.arr, p.table+"."+"updated_at")
			p.arr = append(p.arr, p.table+"."+"deleted_at")
		} else {
			// below code is for recursive
			if field.Type.Kind() == reflect.Struct && externalTable == "" {
				p.getTag(field.Type)
				continue
			}
		}

		column := field.Tag.Get("json")
		if column == "" {
			continue
		}
		column = p.re.FindString(column)

		switch {
		case externalTable == "-":
			continue
		case externalTable != "":
			column = externalTable
		default:
			column = fmt.Sprintf("%v.%v", p.table, column)
		}

		p.arr = append(p.arr, column)
	}

}
