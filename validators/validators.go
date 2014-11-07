package validators

import (
	"fmt"
	"reflect"
	"regexp"
	"unicode/utf8"

	"github.com/PlanitarInc/validate"
)

const (
	/* XXX This should be sufficient for 99% of the cases */
	idPattern         = `([0-9a-zA-Z][-0-9a-zA-Z.+_]*)?[0-9a-zA-Z]`
	domainnamePattern = idPattern + `\.[a-zA-Z]{2,10}`
	emailPattern      = idPattern + `@` + domainnamePattern
)

var (
	V = validate.V{
		"nonnegative":    nonnegativeValidator,
		"nonempty":       nonemptyValidator,
		"notnull":        notnullValidator,
		"strlimit-2-2":   StrLimit(2, 2),
		"strlimit-1-20":  StrLimit(1, 20),
		"strlimit-1-128": StrLimit(1, 128),
		"email":          REMatch(emailPattern, "invalid email"),
	}
)

func nonnegativeValidator(src interface{}) interface{} {
	negative := false

	switch src.(type) {
	default:
		return "Should be an integer"

	case int8:
		n := src.(int8)
		negative = n < 0
	case int16:
		n := src.(int16)
		negative = n < 0
	case int32:
		n := src.(int32)
		negative = n < 0
	case int64:
		n := src.(int64)
		negative = n < 0
	case int:
		n := src.(int)
		negative = n < 0
	}

	if negative {
		return "Should be nonnegative"
	}

	return nil
}

func nonemptyValidator(src interface{}) interface{} {
	str, ok := src.(string)
	if !ok {
		return "Should be a string"
	}

	if len(str) == 0 {
		return "Should be nonempty"
	}

	return nil
}

func StrLimit(min, max uint) validate.ValidatorFn {
	typErr := "Should be a string or byte array"
	minErr := fmt.Sprintf("Minimal length is %d", min)
	maxErr := fmt.Sprintf("Maximal length is %d", max)

	return func(src interface{}) interface{} {
		length := uint(0)
		switch src.(type) {
		default:
			return typErr
		case string:
			length = uint(utf8.RuneCountInString(src.(string)))
		case []byte:
			length = uint(utf8.RuneCount(src.([]byte)))
		}
		if length < min {
			return minErr
		}
		if length > max {
			return maxErr
		}
		return nil
	}
}

func notnullValidator(src interface{}) interface{} {
	val := reflect.ValueOf(src)

	switch val.Kind() {
	default:
		if src == nil {
			return "Expected non null pointer"
		}
		return nil

	case reflect.Interface:
		fallthrough
	case reflect.Map:
		fallthrough
	case reflect.Ptr:
		fallthrough
	case reflect.Slice:
		if val.IsNil() {
			return "Expected non null pointer"
		}

		return nil
	}
}

func REMatch(pattern string, mismatchError ...interface{}) validate.ValidatorFn {
	re := regexp.MustCompile(pattern)
	var mismatchErr interface{}
	if len(mismatchError) == 0 {
		mismatchErr = "Value should match the pattern: " + pattern
	} else {
		mismatchErr = mismatchError[0]
	}

	return func(src interface{}) interface{} {
		var match bool
		switch src.(type) {
		default:
			return "Unsupported type"

		case []byte:
			match = re.Match(src.([]byte))
		case string:
			match = re.MatchString(src.(string))
		}
		if !match {
			return mismatchErr
		}
		return nil
	}
}
