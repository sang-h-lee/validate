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
	maxErr := func(n int) string { return fmt.Sprintf("Maximum length is %d", n) }
	minErr := func(n int) string { return fmt.Sprintf("Minimum length is %d", n) }

	Ω(StrLimit(0, 0)("")).ShouldNot(HaveOccurred())
	Ω(StrLimit(2, 2)("aa")).ShouldNot(HaveOccurred())
	Ω(StrLimit(0, 5)("")).ShouldNot(HaveOccurred())
	Ω(StrLimit(3, 5)("12345")).ShouldNot(HaveOccurred())
	Ω(StrLimit(2, 4)("123")).ShouldNot(HaveOccurred())

	Ω(StrLimit(1, 2)("")).Should(Equal(minErr(1)))
	Ω(StrLimit(10, 20)("abcd ef")).Should(Equal(minErr(10)))
	Ω(StrLimit(0, 5)("123456")).Should(Equal(maxErr(5)))
	Ω(StrLimit(0, 3)("1234567")).Should(Equal(maxErr(3)))

	Ω(StrLimit(0, 1)(nil)).Should(Equal(nonstringErr))
	Ω(StrLimit(0, 2)(1)).Should(Equal(nonstringErr))
	Ω(StrLimit(1, 40)(12.1)).Should(Equal(nonstringErr))

	arr := []string{}
	Ω(StrLimit(1, 1)(arr)).Should(BeNil())
	arr = []string{""}
	errs := map[int]string{0: minErr(1)}
	Ω(StrLimit(1, 1)(arr)).ShouldNot(Equal(errs))
	arr = []string{"a"}
	Ω(StrLimit(1, 1)(arr)).Should(BeNil())
	arr = []string{"", "asd", "bsd", "qs", ""}
	e := StrLimit(1, 2)(arr)
	Ω(e).Should(HaveKeyWithValue(0, minErr(1)))
	Ω(e).Should(HaveKeyWithValue(1, maxErr(2)))
	Ω(e).Should(HaveKeyWithValue(2, maxErr(2)))
	Ω(e).Should(HaveKeyWithValue(4, minErr(1)))
	Ω(e).Should(HaveLen(4))
	arr = []string{"aa", "ab", "ac", "a"}
	Ω(StrLimit(1, 2)(arr)).Should(BeNil())
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

func TestRegexpValidator(t *testing.T) {
	RegisterTestingT(t)

	errMsg := func(p string) string {
		return "Value should match the pattern: " + p
	}

	Ω(REMatch("")("")).Should(BeNil())
	Ω(REMatch("abc")("a")).Should(Equal(errMsg("abc")))
	Ω(REMatch("^abc")("aabc")).Should(Equal(errMsg("^abc")))

	Ω(REMatch("")([]byte{})).Should(BeNil())
	Ω(REMatch("w")([]byte("qwe"))).Should(BeNil())
	Ω(REMatch("a?b?c")([]byte("bbb"))).Should(Equal(errMsg("a?b?c")))

	Ω(REMatch("a")([]string{})).Should(BeNil())
	Ω(REMatch("a")([]string{"a", "ba", "bab"})).Should(BeNil())
	Ω(REMatch("a")([]string{"a", "bb", "bab"})).Should(Equal([]interface{}{
		nil, errMsg("a"), nil,
	}))
	Ω(REMatch("c")([]string{"a", "bb", "bab"})).Should(Equal([]interface{}{
		errMsg("c"), errMsg("c"), errMsg("c"),
	}))

	v := REMatch("^ab+a$", "fail")
	Ω(v("aba")).Should(BeNil())
	Ω(v([]byte("abbbba"))).Should(BeNil())
	Ω(v("ababa")).Should(Equal("fail"))
	Ω(v([]byte("aa"))).Should(Equal("fail"))

	Ω(notnullValidator(1)).ShouldNot(Equal("Unsupported type"))
}

func TestEmailValidator(t *testing.T) {
	RegisterTestingT(t)

	v, ok := V["email"]
	Ω(ok).Should(BeTrue())

	Ω(v("d@p.aa")).Should(BeNil())
	Ω(v("dmitri@planitar.com")).Should(BeNil())
	Ω(v("D.m.I.t.R.i@p.L.a.N.i.T.a.R.cOm")).Should(BeNil())
	Ω(v("angel's@yahoo.com")).Should(BeNil())
	Ω(v("alyson.o'laughlin@wsdevelopment.com")).Should(BeNil())
	Ω(v("happily_@planitar.com")).Should(BeNil())
	Ω(v("happily-@addr.com")).Should(BeNil())
	Ω(v("-happily@addr.com")).Should(BeNil())

	Ω(v("-bad.@addr.com")).Should(Equal("invalid email"))
	Ω(v(".bad.@addr.com")).Should(Equal("invalid email"))
	Ω(v("bad.@addr.com")).Should(Equal("invalid email"))
	Ω(v("bad..mail@addr.com")).Should(Equal("invalid email"))
	Ω(v("@bad.com")).Should(Equal("invalid email"))
	Ω(v("a@bad.")).Should(Equal("invalid email"))
	Ω(v("a@.bad.com")).Should(Equal("invalid email"))
	Ω(v("a@a")).Should(Equal("invalid email"))
	Ω(v("@")).Should(Equal("invalid email"))
	Ω(v("")).Should(Equal("invalid email"))
	Ω(v("sdasd.asdas.com")).Should(Equal("invalid email"))
}

func TestPasswordValidator(t *testing.T) {
	RegisterTestingT(t)

	Ω(PasswordValidator("")).Should(Equal("invalid password"))
	Ω(PasswordValidator("bcDEF67")).Should(Equal("invalid password"))
	Ω(PasswordValidator("aaaaaaaaa")).Should(Equal("invalid password"))
	Ω(PasswordValidator("AAAAAAAAA")).Should(Equal("invalid password"))
	Ω(PasswordValidator("aAaAaAaAa")).Should(Equal("invalid password"))
	Ω(PasswordValidator("1010101010")).Should(Equal("invalid password"))
	Ω(PasswordValidator("aaaa101010")).Should(Equal("invalid password"))

	Ω(PasswordValidator("bcDEF67_")).Should(BeNil())
	Ω(PasswordValidator("Aaaa101010")).Should(BeNil())
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
	Ω(v("1")).Should(Equal("Minimum length is 2"))
	Ω(v("123")).Should(Equal("Maximum length is 2"))

	v, ok = V["strlimit-1-20"]
	Ω(ok).Should(BeTrue())
	Ω(v("")).Should(Equal("Minimum length is 1"))
	Ω(v(strings.Repeat("1", 21))).Should(Equal("Maximum length is 20"))

	v, ok = V["strlimit-1-128"]
	Ω(ok).Should(BeTrue())
	Ω(v("")).Should(Equal("Minimum length is 1"))
	Ω(v(strings.Repeat("1", 129))).Should(Equal("Maximum length is 128"))

	/* Presence of email validator was tested in TestEmailValidator() */

	v, ok = V["password"]
	Ω(ok).Should(BeTrue())
	Ω(v("")).Should(Equal("invalid password"))
}
