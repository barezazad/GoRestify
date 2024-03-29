package excel

import (
	"GoRestify/pkg/pkg_consts"
	"bytes"
	"fmt"
	"reflect"
	"time"

	"github.com/xuri/excelize/v2"
)

// Builder is used for initiate the builder design pattern
type Builder struct {
	SheetCount  int
	ActiveSheet string
	File        *excelize.File

	err    error
	part   string
	Sheets map[string]*SheetInfo
}

// SheetInfo demonstrate sheet property
type SheetInfo struct {
	col    int
	Row    int
	header []string
}

// New initiate the functionality
func New(part string) *Builder {
	// return new(Builder)
	return &Builder{
		File:   excelize.NewFile(),
		part:   part,
		Sheets: make(map[string]*SheetInfo, 100),
	}
}

// AddSheet create new sheet, at first step rename Sheet1
func (b *Builder) AddSheet(name string) *Builder {
	if b.SheetCount == 0 {
		b.File.SetSheetName("Sheet1", name)
	} else {
		b.File.NewSheet(name)
	}
	b.SheetCount++
	return b
}

// Generate is used for save the output to the buffer
func (b *Builder) Generate() (*bytes.Buffer, string, error) {
	buff, err := b.File.WriteToBuffer()

	fileName := time.Now().UTC().Format("data-20060102150405.xlsx")

	return buff, fileName, err
}

// Active is used for choose the sheet for affect style and data
func (b *Builder) Active(sheet string) *Builder {
	b.ActiveSheet = sheet
	return b
}

// SetPageLayout is used for set orientation and size
func (b *Builder) SetPageLayout(orientation string, size string) *Builder {
	if b.err != nil {
		return b
	}

	var sizeInt int
	switch size {
	case "A4":
		sizeInt = 9

	}

	if orientation == "portrait" {

		b.err = b.File.SetPageLayout(
			b.ActiveSheet,
			excelize.PageLayoutOrientation(excelize.OrientationPortrait),
			excelize.PageLayoutPaperSize(sizeInt),
		)
	} else {
		b.err = b.File.SetPageLayout(
			b.ActiveSheet,
			excelize.PageLayoutOrientation(excelize.OrientationLandscape),
			excelize.PageLayoutPaperSize(sizeInt),
		)

	}

	return b
}

// SetPageMargins is used for set size of margin
func (b *Builder) SetPageMargins(size float64) *Builder {
	if b.err != nil {
		return b
	}

	b.err = b.File.SetPageMargins(b.ActiveSheet,
		excelize.PageMarginBottom(size),
		excelize.PageMarginFooter(size),
		excelize.PageMarginHeader(size),
		excelize.PageMarginLeft(size),
		excelize.PageMarginRight(size),
		excelize.PageMarginTop(size),
	)

	return b
}

// SetHeaderFooter put proper header and footer content
func (b *Builder) SetHeaderFooter() *Builder {
	if b.err != nil {
		return b
	}

	b.err = b.File.SetHeaderFooter(b.ActiveSheet,
		&excelize.FormatHeaderFooter{
			DifferentFirst:   true,
			DifferentOddEven: true,
			OddHeader:        "&R&P",
			OddFooter:        "&C&F",
			EvenHeader:       "&L&P",
			EvenFooter:       "&L&D&R&T",
			FirstHeader:      `&CCenter &"-,Bold"Bold&"-,Regular"HeaderU+000A&D`,
		})

	return b

}

// SetDocProps put document's property
func (b *Builder) SetDocProps() *Builder {
	if b.err != nil {
		return b
	}

	b.err = b.File.SetDocProps(&excelize.DocProperties{
		Category:      b.part,
		ContentStatus: "Done",
		Created:       time.Now().UTC().Format(pkg_consts.DateTimeLayout),
		Creator:       "Excel",
		Description:   "Excel",
		Identifier:    "xlsx",
		Keywords:      "Spreadsheet",
		Revision:      "0",
		Subject:       b.part,
		Title:         b.part,
		Language:      "en-US",
	})

	return b

}

// SetColWidth is used for changing the width of columns
func (b *Builder) SetColWidth(start, end string, size float64) *Builder {
	b.File.SetColWidth(b.ActiveSheet, start, end, size)
	return b
}

// WriteHeader is used for writing to the first line
func (b *Builder) WriteHeader(cols ...string) *Builder {

	var inter []interface{}

	for _, v := range cols {
		inter = append(inter, v)
	}

	b.Sheets[b.ActiveSheet] = &SheetInfo{col: len(cols)}

	b.File.SetSheetRow(b.ActiveSheet, "A1", &inter)
	return b
}

// SetSheetFields fill the Excel with rows
func (b *Builder) SetSheetFields(cols ...string) *Builder {
	b.Sheets[b.ActiveSheet].header = cols
	return b
}

// WriteData insert table to the sheet
func (b *Builder) WriteData(table interface{}) *Builder {

	rows := reflect.ValueOf(table)
	if rows.Kind() == reflect.Slice {
		for i := 0; i < rows.Len(); i++ {
			item := rows.Index(i)
			var inter []interface{}
			for _, v := range b.Sheets[b.ActiveSheet].header {
				f := reflect.Indirect(item).FieldByName(v)
				inter = append(inter, f)
			}
			b.File.SetSheetRow(b.ActiveSheet, fmt.Sprint("A", i+2), &inter)
		}

		b.Sheets[b.ActiveSheet].Row = rows.Len() + 1
	}

	return b

}

// AddTable style the sheet
func (b *Builder) AddTable() *Builder {
	colCount, _ := excelize.ColumnNumberToName(b.Sheets[b.ActiveSheet].col)
	b.File.AddTable(b.ActiveSheet, "A1",
		fmt.Sprint(colCount, b.Sheets[b.ActiveSheet].Row),
		`{"table_name":"table","table_style":"TableStyleMedium2", "show_first_column":true,"show_last_column":true,"show_row_stripes":false,"show_column_stripes":true}`)

	return b
}
