package dictionary

import (
	"fmt"
	"strings"
)

// SafeTranslate it return word with translate
func SafeTranslate(lang Lang, str string, params ...interface{}) string {

	term, ok := thisTerms[strings.ToLower(str)]
	pattern := str

	if ok {
		switch lang {
		case En:
			pattern = term.En
		case Ku:
			pattern = term.Ku
		case Ar:
			pattern = term.Ar
		}
	}

	translatedParams := make([]interface{}, len(params))

	for i, v := range params {
		switch v := v.(type) {
		case string:
			translatedParams[i] = Translate(lang, v)
		default:
			translatedParams[i] = fmt.Sprint(v)
		}
	}

	if len(translatedParams) > 0 {
		pattern = fmt.Sprintf(pattern, translatedParams...)
	}

	return pattern

}

// Translate the requested term
func Translate(lang Lang, str string, params ...interface{}) string {

	str = strings.ReplaceAll(str, "!!!", "")
	return SafeTranslate(lang, str, params...)
}
