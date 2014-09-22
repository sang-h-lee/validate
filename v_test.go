package validate

import (
	"fmt"
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
	vd["long"] = func(i interface{}) error {
		s := i.(string)
		if len(s) < 5 {
			return fmt.Errorf("%q is too short", s)
		}
		return nil
	}
	vd["short"] = func(i interface{}) error {
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
	vd["odd"] = func(i interface{}) error {
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
	vd["nonzero"] = func(i interface{}) error {
		n := i.(int)
		if n == 0 {
			return fmt.Errorf("should be nonzero")
		}
		return nil
	}
	vd["odd"] = func(i interface{}) error {
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
	vd["nonzero"] = func(i interface{}) error {
		n := i.(int)
		if n == 0 {
			return fmt.Errorf("should be nonzero")
		}
		return nil
	}
	vd["odd"] = func(i interface{}) error {
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
	vd["nonzero"] = func(i interface{}) error {
		n := i.(int)
		if n == 0 {
			return fmt.Errorf("should be nonzero")
		}
		return nil
	}
	vd["odd"] = func(i interface{}) error {
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
	vd["nonzero"] = func(i interface{}) error {
		n := i.(int)
		if n == 0 {
			return fmt.Errorf("should be nonzero")
		}
		return nil
	}
	vd["odd"] = func(i interface{}) error {
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
	vd["nonzero"] = func(i interface{}) error {
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
	vd["wrong"] = func(i interface{}) error {
		return fmt.Errorf("WRONG: %v", i)
	}

	errs := vd.Validate(7)
	if errs != nil {
		t.Fatalf("non-structs should always pass validation: %v", errs)
	}
}
