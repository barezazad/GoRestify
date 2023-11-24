package db_error

import (
	"errors"
	"strings"

	"GoRestify/pkg/pkg_err"
	"GoRestify/pkg/utils"
	"GoRestify/pkg/validator"

	"gorm.io/gorm"
)

// Parse is an internal method for generate proper database error
func Parse(err error, entity string, action validator.Action) error {

	switch {

	// in case there is no error
	case err == nil:
		err = nil

		// record not found
	case errors.Is(err, gorm.ErrRecordNotFound):

		err = pkg_err.Take(err, "E1057706").
			Message(pkg_err.RecordNotFoundInX, entity).Custom(pkg_err.BadRequestErr).Build()

		// duplicate error
	case strings.Contains(strings.ToUpper(err.Error()), "DUPLICATE"):

		field := utils.RegexFindBetweenTwoPattern(`for key`, ``, err.Error())
		// value := utils.RegexFindBetweenTwoPattern(`Duplicate entry`, `for key`, err.Error())

		if strings.Contains(field, "'") {
			field = strings.ReplaceAll(field, "'", "")
		}
		if strings.Contains(field, ".") {
			field = strings.Split(field, ".")[1]
		}

		err = pkg_err.AddInvalidParam(err, field, pkg_err.XIsAlreadyExist, field)
		err = pkg_err.Take(err, "E1084012").
			Message(pkg_err.XIsUnique, field).Custom(pkg_err.DuplicateErr).Build()

		// unknown column error
	case strings.Contains(strings.ToUpper(err.Error()), "UNKNOWN COLUMN"):

		err = pkg_err.Take(err, "E1090487").
			Message(pkg_err.ValidationForXFailedInX, entity, action).
			Custom(pkg_err.ValidationFailedErr).Build()

		// foreign key error
	case strings.Contains(strings.ToUpper(err.Error()), "FOREIGN"):
		field := utils.RegexFindBetweenTwoPattern(`FOREIGN KEY \(`, `\) REFERENCES`, err.Error())

		switch action {
		case validator.Create:
			err = pkg_err.Take(err, "E1035604").
				Message(pkg_err.CheckThisXField, field).
				Custom(pkg_err.ForeignErr).Build()
		default:
			err = pkg_err.Take(err, "E1063447").
				Message(pkg_err.ItHasRelationWithSomeElementSoItIsNotX, action).
				Custom(pkg_err.ForeignErr).Build()
		}

		// default
	default:
		err = pkg_err.Take(err, "E1053857").
			Message(pkg_err.BadRequestForXInX, entity, action).
			Custom(pkg_err.BadRequestErr).Build()
	}

	return err
}
