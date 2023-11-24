package validator

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"

	"GoRestify/pkg/decimal"
	"GoRestify/pkg/pkg_config"
	"GoRestify/pkg/pkg_err"
	"GoRestify/pkg/utils"
)

// Action core action
type Action string

// main actions
const (
	Find   Action = "find"
	List   Action = "list"
	Create Action = "create"
	Update Action = "update"
	Save   Action = "save"
	Delete Action = "delete"
)

// ValidateModel extract the tag of bind and field from json and bind tag
func ValidateModel(model interface{}, entity string, action Action) (err error) {

	// reflect interface to get value and struct tags
	reflexType := reflect.TypeOf(model)
	reflectValue := reflect.ValueOf(model)

	for i := 0; i < reflexType.NumField(); i++ {
		// get field
		field := reflexType.Field(i)
		// get value
		value := reflectValue.Field(i).Interface()

		// to get tag bind value in field
		bindTag := field.Tag.Get("bind")

		if bindTag == "-" || bindTag == "" {
			continue
		}

		// if reflect is struct, and it will validate those fields in nested struct
		if reflect.ValueOf(value).Kind() == reflect.Struct {
			// reflect interface to get value and struct tags

			reflexTypeN := reflect.TypeOf(value)
			reflectValueN := reflect.ValueOf(value)

			if reflexTypeN.String() == "time.Time" {
				continue
			}

			for j := 0; j < reflexTypeN.NumField(); j++ {
				// get field and value
				fieldN := reflexTypeN.Field(j)
				valueN := reflectValueN.Field(j).Interface()

				// validate method to check value and compare with tags,then get proper error
				err = validationCase(err, fieldN, valueN, action)
			}
		} else {
			// validate method to check value and compare with tags,then get proper error
			err = validationCase(err, field, value, action)
		}
	}

	if err != nil {
		err = pkg_err.Take(err, "E1088589").
			Message(pkg_err.ValidationForXFailedInX, entity, action).
			Custom(pkg_err.ValidationFailedErr).Build()
	}
	return
}

// validationCase compares value with bind tag conditions, and it gets proper error message
func validationCase(err error, field reflect.StructField, value interface{}, action Action) (errors error) {

	if reflect.ValueOf(value).Kind() == reflect.Ptr {
		value = reflectPointerToValue(value)
	}

	// get bind tag to binding field
	bindTag := field.Tag.Get("bind")

	// get field name
	jsonTag := field.Tag.Get("json")
	regex := regexp.MustCompile(`\w+`)
	fieldName := regex.FindString(jsonTag)
	// to remove _id in field name like city_id is required will be city is required
	fieldName = strings.ReplaceAll(fieldName, "_id", "")

	// split tags by comma and find all
	tagByAction := bindTagByAction(bindTag, action)

	allTags := strings.Split(tagByAction, ",")

	if strings.Contains(tagByAction, "if_exist") {
		strValue := fmt.Sprintf("%v", value)
		if strValue == "" || strValue == "0" {
			return err
		}
	}

	for _, v := range allTags {

		// split taq by equal to find key and value
		splitTag := strings.Split(v, "=")
		var tagKey string
		var tagValue interface{}

		if len(splitTag) > 0 {
			tagKey = splitTag[0]
		}
		if len(splitTag) > 1 {
			tagValue = splitTag[1]
		}

		// custom validation cases
		switch tagKey {

		case "required":
			strValue := fmt.Sprintf("%v", value)
			if strValue == "" || strValue == "0" {
				err = pkg_err.AddInvalidParam(err, fieldName,
					pkg_err.XIsRequired, fieldName)
			}

		case "min":
			intTagValue, _ := utils.StrToInt(tagValue.(string))
			if len(value.(string)) < intTagValue {
				err = pkg_err.AddInvalidParam(err, fieldName,
					pkg_err.MinimumAcceptedCharacterForXisX, fieldName, tagValue)
			}

		case "max":
			intTagValue, _ := utils.StrToInt(tagValue.(string))
			if len(value.(string)) > intTagValue {
				err = pkg_err.AddInvalidParam(err, fieldName,
					pkg_err.MaximumAcceptedCharacterForXisX, fieldName, tagValue)
			}

		case "lte":
			switch v := value.(type) {
			case float64:
				floatTagValue, _ := utils.StrToFloat(tagValue.(string))
				if v > floatTagValue {
					err = pkg_err.AddInvalidParam(err, fieldName,
						pkg_err.MinimumAcceptedValueForXIsX, fieldName, tagValue)
				}

			case uint:
				uintTagValue, _ := utils.StrToUint(tagValue.(string))

				if v > uintTagValue {
					err = pkg_err.AddInvalidParam(err, fieldName,
						pkg_err.MinimumAcceptedValueForXIsX, fieldName, tagValue)
				}

			case decimal.Decimal:
				decimalTagValue, _ := decimal.NewFromString(tagValue.(string))
				if !v.LessThanOrEqual(decimalTagValue) {
					err = pkg_err.AddInvalidParam(err, fieldName,
						pkg_err.MinimumAcceptedValueForXIsX, fieldName, tagValue)
				}
			}

		case "gte":
			switch v := value.(type) {
			case float64:
				floatTagValue, _ := utils.StrToFloat(tagValue.(string))
				if v < floatTagValue {
					err = pkg_err.AddInvalidParam(err, fieldName,
						pkg_err.MaximumAcceptedValueForXIsX, fieldName, tagValue)
				}
			case int64:
				int64TagValue, _ := utils.StrToInt64(tagValue.(string))
				if v < int64TagValue {
					err = pkg_err.AddInvalidParam(err, fieldName,
						pkg_err.MaximumAcceptedValueForXIsX, fieldName, tagValue)
				}
			case uint:
				uintTagValue, _ := utils.StrToUint(tagValue.(string))
				if v < uintTagValue {
					err = pkg_err.AddInvalidParam(err, fieldName,
						pkg_err.MaximumAcceptedValueForXIsX, fieldName, tagValue)
				}
			case decimal.Decimal:
				decimalTagValue, _ := decimal.NewFromString(tagValue.(string))
				if !v.GreaterThanOrEqual(decimalTagValue) {
					err = pkg_err.AddInvalidParam(err, fieldName,
						pkg_err.MaximumAcceptedValueForXIsX, fieldName, tagValue)
				}
			}

		case "one_of":
			types := pkg_config.Config.EnumLists[tagValue.(string)]
			if ok, _ := utils.ArrayIncludes(types, value); !ok {
				err = pkg_err.AddInvalidParam(err, fieldName,
					pkg_err.AcceptedValueForXAreX, fieldName, types)
			}

		case "contain":
			if ok := strings.Contains(value.(string), tagValue.(string)); !ok {
				err = pkg_err.AddInvalidParam(err, fieldName,
					pkg_err.InvalidValueDoNotIncludeX, tagValue)
			}

		case "password":
			secure := true
			regs := []string{".{8,}", "[a-z]", "[A-Z]", "[0-9]"}
			for _, rg := range regs {
				t, _ := regexp.MatchString(rg, value.(string))
				if !t {
					secure = false
					break
				}
			}
			if !secure {
				err = pkg_err.AddInvalidParam(err, fieldName,
					pkg_err.XisNotValid, fieldName)
			}

		case "username":
			re := regexp.MustCompile("^[a-zA-Z0-9\\._-]+$")
			if !re.MatchString(value.(string)) {
				err = pkg_err.AddInvalidParam(err, fieldName, pkg_err.XisNotValid, fieldName)
			}

		case "phone":

			phoneNumber := value.(string)

			re := regexp.MustCompile("[^0-9]+")
			phoneNumber = fmt.Sprintf("%v", re.ReplaceAllString(phoneNumber, ""))

			phoneRegex := "^(964[0-9]{10})$"
			match, _ := regexp.MatchString(phoneRegex, phoneNumber)

			if !match {
				err = pkg_err.AddInvalidParam(err, fieldName, pkg_err.XisNotValid, fieldName)
			}

		case "email":
			re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
			if !re.MatchString(value.(string)) {
				err = pkg_err.AddInvalidParam(err, fieldName,
					pkg_err.XisNotValid, fieldName)
			}

		case "birthday":

			switch v := value.(type) {
			case time.Time:

				minAge := time.Now().AddDate(-90, 0, 0)
				maxAge := time.Now().AddDate(-10, 0, 0)

				if v.Before(minAge) || v.After(maxAge) {
					err = pkg_err.AddInvalidParam(err, fieldName,
						pkg_err.XisNotValid, fieldName)
				}
			}

		case "pin":
			re := regexp.MustCompile("^[0-9]+$")
			if !re.MatchString(value.(string)) {
				err = pkg_err.AddInvalidParam(err, fieldName,
					pkg_err.XisNotValid, fieldName)
			}
		}

	}

	return err
}

// bindTagByAction this function separate validation per type of action
func bindTagByAction(tag string, action Action) (result string) {

	perAction := utils.RegexFindBetweenTwoPattern(fmt.Sprintf("%v:", action), `\|`, tag)
	if perAction != "" {
		result = perAction
	}

	all := utils.RegexFindBetweenTwoPattern(`all:`, `\|`, tag)
	if all != "" {
		result += "," + all
	}

	if result == "" && !strings.Contains(tag, "|") && !strings.Contains(tag, ":") {
		result = tag
	}

	return
}

func reflectPointerToValue(v interface{}) interface{} {

	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
		rv = rv.Elem()
	}

	if rv.IsValid() {
		return rv.Interface()
	}
	return ""
}
