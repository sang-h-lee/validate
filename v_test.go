package validate

import (
	"fmt"
	"reflect"
	"testing"
)

func ExampleV_Validate() {
	type X struct {
		A string `validate:"long"`
		B string `validate:"short"`
		C string `validate:"long,short"`
		D string
	}

	vd := make(V)
	vd["long"] = func(i interface{}) interface{} {
		s := i.(string)
		if len(s) < 5 {
			return fmt.Errorf("%q is too short", s)
		}
		return nil
	}
	vd["short"] = func(i interface{}) interface{} {
		s := i.(string)
		if len(s) >= 5 {
			return fmt.Errorf("%q is too long", s)
		}
		return nil
	}

	fmt.Println(vd.Validate(X{
		A: "hello there",
		B: "hi",
		C: "help me",
		D: "I am not validated",
	}))

	// Output: map[C:"help me" is too long]
}

func TestV_Validate_allgood(t *testing.T) {
	type X struct {
		A int `validate:"odd"`
	}

	vd := make(V)
	vd["odd"] = func(i interface{}) interface{} {
		n := i.(int)
		if n&1 == 0 {
			return fmt.Errorf("%d is not odd", n)
		}
		return nil
	}

	errs := vd.Validate(X{
		A: 1,
	})

	if errs != nil {
		t.Fatalf("unexpected errors for a valid struct: %v", errs)
	}
}

func TestV_Validate_undef(t *testing.T) {
	type X struct {
		A string `validate:"oops"`
	}

	vd := make(V)

	errs := vd.Validate(X{
		A: "oh my",
	})

	if len(errs) == 0 {
		t.Fatal("no errors returned for an undefined validator")
	}
	if len(errs) != 1 {
		t.Fatalf("too many errors returns for an undefined validator: %v", errs)
	}
	if errs["A"] == nil {
		t.Fatalf("expected error for field A: %v", errs)
	}
	if errs["A"].(error).Error() != `undefined validator: "oops"` {
		t.Fatal("wrong message for an undefined validator:", errs["A"])
	}
}

func TestV_Validate_multi_first_fails(t *testing.T) {
	type X struct {
		A int `validate:"nonzero,odd"`
	}

	vd := make(V)
	vd["nonzero"] = func(i interface{}) interface{} {
		n := i.(int)
		if n == 0 {
			return fmt.Errorf("should be nonzero")
		}
		return nil
	}
	vd["odd"] = func(i interface{}) interface{} {
		n := i.(int)
		if n&1 == 0 {
			return fmt.Errorf("%d is not odd", n)
		}
		return nil
	}
	errs := vd.Validate(X{
		A: 0,
	})
	if len(errs) != 1 {
		t.Fatalf("wrong number of errors for two failures: %v", errs)
	}
	if errs["A"] == nil {
		t.Fatalf("expected error for field A: %v", errs)
	}
	if errs["A"].(error).Error() != "should be nonzero" {
		t.Fatal("first error should be nonzero:", errs["A"])
	}
}

func TestV_Validate_multi_first_passes(t *testing.T) {
	type X struct {
		A int `validate:"nonzero,odd"`
	}

	vd := make(V)
	vd["nonzero"] = func(i interface{}) interface{} {
		n := i.(int)
		if n == 0 {
			return fmt.Errorf("should be nonzero")
		}
		return nil
	}
	vd["odd"] = func(i interface{}) interface{} {
		n := i.(int)
		if n&1 == 0 {
			return fmt.Errorf("%d is not odd", n)
		}
		return nil
	}
	errs := vd.Validate(X{
		A: 2,
	})
	if len(errs) != 1 {
		t.Fatalf("wrong number of errors for two failures: %v", errs)
	}
	if errs["A"] == nil {
		t.Fatalf("expected error for field A: %v", errs)
	}
	if errs["A"].(error).Error() != "2 is not odd" {
		t.Fatal("second error should be odd:", errs["A"])
	}
}

func ExampleV_Validate_struct() {
	type X struct {
		A int `validate:"nonzero"`
	}

	type Y struct {
		X `validate:"struct,odd"`
	}

	vd := make(V)
	vd["nonzero"] = func(i interface{}) interface{} {
		n := i.(int)
		if n == 0 {
			return fmt.Errorf("should be nonzero")
		}
		return nil
	}
	vd["odd"] = func(i interface{}) interface{} {
		x := i.(X)
		if x.A&1 == 0 {
			return fmt.Errorf("%d is not odd", x.A)
		}
		return nil
	}

	errs := vd.Validate(Y{X{
		A: 0,
	}})

	for k, err := range errs {
		fmt.Printf("%s=>%v\n", k, err)
	}

	// Output: X=>map[A:should be nonzero]
}

func ExampleV_Validate_struct2() {
	type X struct {
		A int `validate:"nonzero"`
	}

	type Y struct {
		X `validate:"struct,odd"`
	}

	vd := make(V)
	vd["nonzero"] = func(i interface{}) interface{} {
		n := i.(int)
		if n == 0 {
			return fmt.Errorf("should be nonzero")
		}
		return nil
	}
	vd["odd"] = func(i interface{}) interface{} {
		x := i.(X)
		if x.A&1 == 0 {
			return fmt.Errorf("%d is not odd", x.A)
		}
		return nil
	}

	errs := vd.Validate(Y{X{
		A: 2,
	}})

	for k, err := range errs {
		fmt.Printf("%s=>%v\n", k, err)
	}

	// Output: X=>2 is not odd
}

func TestV_Validate_uninterfaceable(t *testing.T) {
	type X struct {
		a int `validate:"nonzero"`
	}

	vd := make(V)
	vd["nonzero"] = func(i interface{}) interface{} {
		n := i.(int)
		if n == 0 {
			return fmt.Errorf("should be nonzero")
		}
		return nil
	}

	errs := vd.Validate(X{
		a: 0,
	})

	if len(errs) != 0 {
		t.Fatal("wrong number of errors for two failures:", errs)
	}
}

func TestV_Validate_nonstruct(t *testing.T) {
	vd := make(V)
	vd["wrong"] = func(i interface{}) interface{} {
		return fmt.Errorf("WRONG: %v", i)
	}

	errs := vd.Validate(7)
	if errs != nil {
		t.Fatalf("non-structs should always pass validation: %v", errs)
	}
}

func TestV_Validate_json_name(t *testing.T) {
	type X struct {
		A int `json:"z" validate:"nonzero"`
	}

	vd := make(V)
	vd["nonzero"] = func(i interface{}) interface{} {
		n := i.(int)
		if n == 0 {
			return fmt.Errorf("should be nonzero")
		}
		return nil
	}

	errs := vd.Validate(X{
		A: 0,
	})

	if len(errs) != 1 {
		t.Fatal("wrong number of errors for two failures:", errs)
	}
	if errs["z"] == nil {
		t.Fatal("an error for z field should be present:", errs)
	}
	if errs["z"].(error).Error() != "should be nonzero" {
		t.Fatalf("the error should be nonzero: %s", errs["z"])
	}
}

func TestV_Validate_json_name_nested(t *testing.T) {
	type Z struct {
		B int `json:"z" validate:"nonzero"`
	}
	type X struct {
		A Z `json:"xxx" validate:"struct"`
	}

	vd := make(V)
	vd["nonzero"] = func(i interface{}) interface{} {
		n := i.(int)
		if n == 0 {
			return fmt.Errorf("should be nonzero")
		}
		return nil
	}

	errs := vd.Validate(X{A: Z{
		B: 0,
	}})

	if len(errs) != 1 {
		t.Fatal("wrong number of errors for two failures:", errs)
	}
	if errs["xxx"] == nil {
		t.Fatal("an error for xxx field should be present:", errs)
	}
	nested_errs, ok := errs["xxx"].(map[string]interface{})
	if !ok {
		t.Fatal("an error for xxx should be a map:", nested_errs)
	}
	if nested_errs["z"] == nil {
		t.Fatal("an error for z field should be present:", nested_errs)
	}
	if nested_errs["z"].(error).Error() != "should be nonzero" {
		t.Fatalf("the error should be nonzero: %s", nested_errs["z"])
	}
}

func TestV_Validate_field_order(t *testing.T) {
	type X struct {
		Z string `validate:"longer"`
		A string `validate:"longer"`
		B string `validate:"longer"`
		C string
		D string `validate:"longer"`
		E string
	}

	maxLen := 4
	vd := make(V)
	vd["long"] = func(i interface{}) interface{} {
		s := i.(string)
		if len(s) < 5 {
			return fmt.Errorf("%q is too short", s)
		}
		maxLen = len(s)
		return nil
	}
	vd["longer"] = func(i interface{}) interface{} {
		s := i.(string)
		if len(s) <= maxLen {
			return fmt.Errorf("%q is too short, should be longer than %d",
				s, maxLen)
		}
		maxLen = len(s)
		return nil
	}

	maxLen = 4
	errs := vd.Validate(X{
		Z: "12345",
		A: "hello there",
		B: "hi, hi, hi!!!",
		C: "help me",
		D: "I am not validated",
	})
	if len(errs) != 0 {
		t.Fatal("wrong number of errors:", errs)
	}

	maxLen = 3
	errs = vd.Validate(X{
		Z: "123",
		A: "hello there",
		B: "hi, hi, hi!!!",
		C: "help me",
		D: "I am not validated",
	})
	if len(errs) != 1 {
		t.Fatal("wrong number of errors:", errs)
	}
	if errs["Z"].(error).Error() != `"123" is too short, should be longer than 3` {
		t.Fatal("error for Z field is wrong:", errs)
	}

	maxLen = 3
	errs = vd.Validate(X{
		Z: "123",
		A: "h",
		B: "2",
		C: "help me",
		D: "I",
	})
	if len(errs) != 4 {
		t.Fatal("wrong number of errors:", errs)
	}
	if errs["Z"].(error).Error() != `"123" is too short, should be longer than 3` {
		t.Fatal("error for Z field is wrong:", errs)
	}
	if errs["A"].(error).Error() != `"h" is too short, should be longer than 3` {
		t.Fatal("error for Z field is wrong:", errs)
	}
	if errs["B"].(error).Error() != `"2" is too short, should be longer than 3` {
		t.Fatal("error for Z field is wrong:", errs)
	}
	if errs["D"].(error).Error() != `"I" is too short, should be longer than 3` {
		t.Fatal("error for Z field is wrong:", errs)
	}

	// Output: map[C:"help me" is too long]
}

type ValidatorExample struct {
	Error interface{}
}

func (v ValidatorExample) Validate() interface{} {
	return v.Error
}

func TestV_Validate_Validator(t *testing.T) {
	type X struct {
		V ValidatorExample
	}

	vd := make(V)

	xOk := X{V: ValidatorExample{}}
	if errs := vd.Validate(&xOk); len(errs) != 0 {
		t.Fatal("wrong error: expeted nil; got:", errs)
	}

	err1 := map[string]string{
		"one": "qwe",
		"two": "asd",
	}
	x1 := X{V: ValidatorExample{Error: err1}}
	if errs := vd.Validate(&x1); !reflect.DeepEqual(errs["V"], err1) {
		t.Fatal("wrong error: expeted ", err1, "; got:", errs)
	}

	err2 := map[string]string{
		"two": "asd",
	}
	x2 := X{V: ValidatorExample{Error: err2}}
	if errs := vd.Validate(&x2); !reflect.DeepEqual(errs["V"], err2) {
		t.Fatal("wrong error: expeted ", err2, "; got:", errs)
	}
}

type ArrMapper struct {
	Arr []string
}

func (a ArrMapper) MapValue() interface{} {
	return a.Arr
}

func TestV_Validate_ValueMapper(t *testing.T) {
	type X struct {
		A ArrMapper `validate:"long"`
	}

	vd := make(V)
	vd["long"] = func(i interface{}) interface{} {
		s := i.([]string)
		res := map[int]error{}
		for i := range s {
			if len(s[i]) < 5 {
				res[i] = fmt.Errorf("%q is too short", s[i])
			}
		}
		if len(res) != 0 {
			return res
		}
		return nil
	}

	xEmpty := X{A: ArrMapper{Arr: []string{}}}
	if errs := vd.Validate(&xEmpty); len(errs) != 0 {
		t.Fatal("wrong error: expeted nil; got:", errs)
	}

	xOk := X{A: ArrMapper{Arr: []string{"qweas", "12345"}}}
	if errs := vd.Validate(&xOk); len(errs) != 0 {
		t.Fatal("wrong error: expeted nil; got:", errs)
	}

	x := X{A: ArrMapper{Arr: []string{"q", "qweqwe", "asd"}}}
	errs := vd.Validate(&x)
	if len(errs) != 1 {
		t.Fatal(`wrong number of errors: expected 1; got:`, errs)
	}
	if errs["A"] == nil {
		t.Fatal(`wrong errors: expected error for "A"; got:`, errs)
	}
	aErrs := errs["A"].(map[int]error)
	if aErrs[0].(error).Error() != `"q" is too short` {
		t.Fatal(`wrong :`, aErrs)
	}
	if aErrs[2].(error).Error() != `"asd" is too short` {
		t.Fatal(`wrong :`, aErrs)
	}
	if len(aErrs) != 2 {
		t.Fatal(`wrong number of errors: expected 2; got:`, errs)
	}
}
