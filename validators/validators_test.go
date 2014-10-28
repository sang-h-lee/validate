package validators

import (
	"fmt"
	"strings"
	"testing"

	"github.com/PlanitarInc/validate"
	. "github.com/onsi/gomega"
)

func TestNonnegativeValidator(t *testing.T) {
	RegisterTestingT(t)

	nonnegativeErr := "Should be nonnegative"
	nonintegerErr := "Should be an integer"

	Ω(nonnegativeValidator(0)).ShouldNot(HaveOccurred())
	Ω(nonnegativeValidator(123)).ShouldNot(HaveOccurred())
	Ω(nonnegativeValidator(-1)).Should(Equal(nonnegativeErr))

	Ω(nonnegativeValidator(int8(0))).ShouldNot(HaveOccurred())
	Ω(nonnegativeValidator(int8(32))).ShouldNot(HaveOccurred())
	Ω(nonnegativeValidator(int8(-11))).Should(Equal(nonnegativeErr))

	Ω(nonnegativeValidator(int16(0))).ShouldNot(HaveOccurred())
	Ω(nonnegativeValidator(int16(123))).ShouldNot(HaveOccurred())
	Ω(nonnegativeValidator(int16(-111))).Should(Equal(nonnegativeErr))

	Ω(nonnegativeValidator(int32(0))).ShouldNot(HaveOccurred())
	Ω(nonnegativeValidator(int32(1))).ShouldNot(HaveOccurred())
	Ω(nonnegativeValidator(int32(-1))).Should(Equal(nonnegativeErr))

	Ω(nonnegativeValidator(int64(0))).ShouldNot(HaveOccurred())
	Ω(nonnegativeValidator(int64(131))).ShouldNot(HaveOccurred())
	Ω(nonnegativeValidator(int64(-97))).Should(Equal(nonnegativeErr))

	Ω(nonnegativeValidator(1.1)).Should(Equal(nonintegerErr))
	Ω(nonnegativeValidator("1")).Should(Equal(nonintegerErr))
	Ω(nonnegativeValidator(nil)).Should(Equal(nonintegerErr))
}

func TestNonemptyValidator(t *testing.T) {
	RegisterTestingT(t)

	nonstringErr := "Should be a string"
	nonemptyErr := "Should be nonempty"

	Ω(nonemptyValidator("")).Should(Equal(nonemptyErr))
	Ω(nonemptyValidator(".")).ShouldNot(HaveOccurred())
	Ω(nonemptyValidator("asb")).ShouldNot(HaveOccurred())

	Ω(nonemptyValidator(nil)).Should(Equal(nonstringErr))
	Ω(nonemptyValidator(1)).Should(Equal(nonstringErr))
	Ω(nonemptyValidator(1.1)).Should(Equal(nonstringErr))
}

func TestStrLimitValidator(t *testing.T) {
	RegisterTestingT(t)

	nonstringErr := "Should be a string or byte array"
	maxErr := func(n int) string { return fmt.Sprintf("Maximal length is %d", n) }
	minErr := func(n int) string { return fmt.Sprintf("Minimal length is %d", n) }

	Ω(strLimit(0, 0)("")).ShouldNot(HaveOccurred())
	Ω(strLimit(2, 2)("aa")).ShouldNot(HaveOccurred())
	Ω(strLimit(0, 5)("")).ShouldNot(HaveOccurred())
	Ω(strLimit(3, 5)("12345")).ShouldNot(HaveOccurred())
	Ω(strLimit(2, 4)("123")).ShouldNot(HaveOccurred())

	Ω(strLimit(1, 2)("")).Should(Equal(minErr(1)))
	Ω(strLimit(10, 20)("abcd ef")).Should(Equal(minErr(10)))
	Ω(strLimit(0, 5)("123456")).Should(Equal(maxErr(5)))
	Ω(strLimit(0, 3)("1234567")).Should(Equal(maxErr(3)))

	Ω(strLimit(0, 1)(nil)).Should(Equal(nonstringErr))
	Ω(strLimit(0, 2)(1)).Should(Equal(nonstringErr))
	Ω(strLimit(1, 40)(12.1)).Should(Equal(nonstringErr))
}

func TestNotNull(t *testing.T) {
	RegisterTestingT(t)

	Ω(notnullValidator(nil)).Should(Equal("Expected non null pointer"))
	{
		var src interface{}
		Ω(notnullValidator(src)).Should(Equal("Expected non null pointer"))
	}
	{
		var src map[string]interface{}
		Ω(notnullValidator(src)).Should(Equal("Expected non null pointer"))
	}
	{
		var src []int
		Ω(notnullValidator(src)).Should(Equal("Expected non null pointer"))
	}
	{
		var src *struct{ X int }
		Ω(notnullValidator(src)).Should(Equal("Expected non null pointer"))
	}

	{
		var src interface{}
		Ω(notnullValidator(&src)).ShouldNot(HaveOccurred())
	}
	{
		src := &struct{ X int }{}
		Ω(notnullValidator(src)).ShouldNot(HaveOccurred())
	}
	{
		var src int
		Ω(notnullValidator(&src)).ShouldNot(HaveOccurred())
	}

	Ω(notnullValidator(1)).ShouldNot(HaveOccurred())
	Ω(notnullValidator(1.13)).ShouldNot(HaveOccurred())
	Ω(notnullValidator("")).ShouldNot(HaveOccurred())
	Ω(notnullValidator(struct{ X int }{X: 1})).ShouldNot(HaveOccurred())
	Ω(notnullValidator(struct{ X int }{})).ShouldNot(HaveOccurred())
}

func TestValidatorArray(t *testing.T) {
	RegisterTestingT(t)

	var ok bool
	var v validate.ValidatorFn

	_, ok = V["nonempty"]
	Ω(ok).Should(BeTrue())

	_, ok = V["nonnegative"]
	Ω(ok).Should(BeTrue())

	v, ok = V["notnull"]
	Ω(ok).Should(BeTrue())
	Ω(v(nil)).Should(Equal("Expected non null pointer"))
	Ω(v(&struct{}{})).ShouldNot(HaveOccurred())

	v, ok = V["strlimit-2-2"]
	Ω(ok).Should(BeTrue())
	Ω(v("1")).Should(Equal("Minimal length is 2"))
	Ω(v("123")).Should(Equal("Maximal length is 2"))

	v, ok = V["strlimit-1-20"]
	Ω(ok).Should(BeTrue())
	Ω(v("")).Should(Equal("Minimal length is 1"))
	Ω(v(strings.Repeat("1", 21))).Should(Equal("Maximal length is 20"))

	v, ok = V["strlimit-1-128"]
	Ω(ok).Should(BeTrue())
	Ω(v("")).Should(Equal("Minimal length is 1"))
	Ω(v(strings.Repeat("1", 129))).Should(Equal("Maximal length is 128"))
}
