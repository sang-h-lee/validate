package validators

import (
	"fmt"
	"reflect"
	"regexp"
	"unicode/utf8"

	"github.com/PlanitarInc/validate"
)

const (
	/* Note: synchronize it with PlanitarInc/walk-inside-app:app/lib/utils.js, or
	 * https://github.com/PlanitarInc/walk-inside-app/blob/master/app/lib/utils.js
	 */
	/* XXX This should be sufficient for 99% of the cases */
	idPattern         = "([0-9a-zA-Z][-0-9a-zA-Z.+_']*)?[0-9a-zA-Z_]"
	domainnamePattern = idPattern + "\\.[a-zA-Z]{2,10}"
	emailPattern      = "^" + idPattern + "@" + domainnamePattern + "$"
)

var (
	V = validate.V{
		"nonnegative":     nonnegativeValidator,
		"nonempty":        nonemptyValidator,
		"notnull":         notnullValidator,
		"strlimit-2-2":    StrLimit(2, 2),
		"strlimit-1-20":   StrLimit(1, 20),
		"strlimit-1-128":  StrLimit(1, 128),
		"strlimit-1-256":  StrLimit(1, 256),
		"strlimit-1-512":  StrLimit(1, 512),
		"strlimit-1-1024": StrLimit(1, 1024),
		"strlimit-0-20":   StrLimit(0, 20),
		"strlimit-0-256":  StrLimit(0, 256),
		"strlimit-0-512":  StrLimit(0, 512),
		"strlimit-0-1024": StrLimit(0, 1024),
		"strlimit-0-2048": StrLimit(0, 2048),
		"email":           REMatch(emailPattern, "invalid email"),
		"password":        PasswordValidator,
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
	minErr := fmt.Sprintf("Minimum length is %d", min)
	maxErr := fmt.Sprintf("Maximum length is %d", max)
	validate := func(length uint) interface{} {
		if length < min {
			return minErr
		}
		if length > max {
			return maxErr
		}
		return nil
	}

	return func(src interface{}) interface{} {
		length := uint(0)
		switch src.(type) {
		default:
			return typErr
		case []string:
			arr := src.([]string)
			errs := map[int]interface{}{}
			for i := range arr {
				if e := validate(uint(utf8.RuneCountInString(arr[i]))); e != nil {
					errs[i] = e
				}
			}
			if len(errs) > 0 {
				return errs
			}
			return nil

		case string:
			length = uint(utf8.RuneCountInString(src.(string)))
		case []byte:
			length = uint(utf8.RuneCount(src.([]byte)))
		}
		return validate(length)
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
		case []string:
			arr := src.([]string)
			errArr := make([]interface{}, len(arr))
			match = true
			for i := range arr {
				if re.MatchString(arr[i]) {
					continue
				}
				errArr[i] = mismatchErr
				match = false
			}
			if !match {
				return errArr
			}
			return nil
		}
		if !match {
			return mismatchErr
		}
		return nil
	}
}

func PasswordValidator(src interface{}) interface{} {
	str, ok := src.(string)
	if !ok {
		return "invalid password"
	}

	if len(str) < 8 || len(str) > 128 {
		return "invalid password"
	}
	if m, e := regexp.MatchString("[a-z]", str); !m || e != nil {
		return "invalid password"
	}
	if m, e := regexp.MatchString("[A-Z]", str); !m || e != nil {
		return "invalid password"
	}
	if m, e := regexp.MatchString("[0-9]", str); !m || e != nil {
		return "invalid password"
	}
	return nil
}
