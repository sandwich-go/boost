package validator

import (
	"bytes"
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/fr"
	"github.com/go-playground/locales/nl"
	ut "github.com/go-playground/universal-translator"
	. "github.com/smartystreets/goconvey/convey"
	"reflect"
	"strings"
	"testing"
	"time"
)

type StructLevelInvalidErr struct {
	Value string
}

type SubTest struct {
	Test string `validate:"required"`
}

type TestStruct struct {
	String string `validate:"required" json:"StringVal"`
}

type TestPartial struct {
	NoTag     string
	BlankTag  string     `validate:""`
	Required  string     `validate:"required"`
	SubSlice  []*SubTest `validate:"required,dive"`
	Sub       *SubTest
	SubIgnore *SubTest `validate:"-"`
	Anonymous struct {
		A         string     `validate:"required"`
		ASubSlice []*SubTest `validate:"required,dive"`

		SubAnonStruct []struct {
			Test      string `validate:"required"`
			OtherTest string `validate:"required"`
		} `validate:"required,dive"`
	}
}

func StructLevelInvalidError(_ context.Context, sl StructLevel) {
	top := sl.Top().Interface().(StructLevelInvalidErr)
	s := sl.Current().Interface().(StructLevelInvalidErr)

	if top.Value == s.Value {
		sl.ReportError(nil, "Value", "Value", "required", "")
	}
}

type valuer struct {
	Name string
}

func (v valuer) Value() (driver.Value, error) {
	if v.Name == "errorme" {
		panic("SQL Driver Valuer error: some kind of error")
		// return nil, errors.New("some kind of error")
	}

	if len(v.Name) == 0 {
		return nil, nil
	}

	return v.Name, nil
}

func ValidateValuerType(field reflect.Value) interface{} {
	if valuer, ok := field.Interface().(driver.Valuer); ok {

		val, err := valuer.Value()
		if err != nil {
			// handle the error how you want
			return nil
		}

		return val
	}

	return nil
}

type MadeUpCustomType struct {
	FirstName string
	LastName  string
}

func ValidateCustomType(field reflect.Value) interface{} {
	if cust, ok := field.Interface().(MadeUpCustomType); ok {

		if len(cust.FirstName) == 0 || len(cust.LastName) == 0 {
			return ""
		}

		return cust.FirstName + " " + cust.LastName
	}

	return ""
}

type CustomMadeUpStruct struct {
	MadeUp        MadeUpCustomType `validate:"required"`
	OverriddenInt int              `validate:"gt=1"`
}

func OverrideIntTypeForSomeReason(field reflect.Value) interface{} {
	if i, ok := field.Interface().(int); ok {
		if i == 1 {
			return "1"
		}

		if i == 2 {
			return "12"
		}
	}

	return ""
}

func TestValidator(t *testing.T) {
	Convey(`test validator.SetTagName`, t, func() {
		Default.SetTagName("val")
		type Test struct {
			Name string `val:"len=4"`
		}
		s := &Test{
			Name: "TEST",
		}
		So(Default.Struct(context.Background(), s), ShouldBeNil)
		s.Name = ""
		err := Default.Struct(context.Background(), s)
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, "Key: 'Test.Name' Error:Field validation for 'Name' failed on the 'len' tag")
	})

	Convey(`test validator.ValidateMap`, t, func() {
		type args struct {
			data  map[string]interface{}
			rules map[string]interface{}
		}
		tests := []struct {
			name string
			args args
			want int
		}{
			{
				name: "test nested map in slice",
				args: args{
					data: map[string]interface{}{
						"Test_A": map[string]interface{}{
							"Test_B": "Test_B",
							"Test_C": []map[string]interface{}{
								{
									"Test_D": "Test_D",
								},
							},
							"Test_E": map[string]interface{}{
								"Test_F": "Test_F",
							},
						},
					},
					rules: map[string]interface{}{
						"Test_A": map[string]interface{}{
							"Test_B": "min=2",
							"Test_C": map[string]interface{}{
								"Test_D": "min=2",
							},
							"Test_E": map[string]interface{}{
								"Test_F": "min=2",
							},
						},
					},
				},
				want: 0,
			},

			{
				name: "test nested map error",
				args: args{
					data: map[string]interface{}{
						"Test_A": map[string]interface{}{
							"Test_B": "Test_B",
							"Test_C": []interface{}{"Test_D"},
							"Test_E": map[string]interface{}{
								"Test_F": "Test_F",
							},
							"Test_G": "Test_G",
							"Test_I": []map[string]interface{}{
								{
									"Test_J": "Test_J",
								},
							},
						},
					},
					rules: map[string]interface{}{
						"Test_A": map[string]interface{}{
							"Test_B": "min=2",
							"Test_C": map[string]interface{}{
								"Test_D": "min=2",
							},
							"Test_E": map[string]interface{}{
								"Test_F": "min=100",
							},
							"Test_G": map[string]interface{}{
								"Test_H": "min=2",
							},
							"Test_I": map[string]interface{}{
								"Test_J": "min=100",
							},
						},
					},
				},
				want: 1,
			},
		}
		for _, tt := range tests {
			got := Validate.ValidateMap(context.Background(), tt.args.data, tt.args.rules)
			So(len(got), ShouldEqual, tt.want)
		}
	})

	Convey(`test validator.RegisterTagNameFunc`, t, func() {
		v := New()
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

			if name == "-" {
				return ""
			}

			return name
		})

		type Test struct {
			A interface{}
		}

		tst := &Test{
			A: struct {
				A string `validate:"required"`
			}{
				A: "",
			},
		}

		err := v.Struct(context.Background(), tst)
		So(err, ShouldNotBeNil)

		errs := err.(ValidationErrors)

		So(len(errs), ShouldEqual, 1)
		So(errs.Error(), ShouldEqual, "Key: 'Test.A.A' Error:Field validation for 'A' failed on the 'required' tag")

		tst = &Test{
			A: struct {
				A string `validate:"omitempty,required"`
			}{
				A: "",
			},
		}

		err = v.Struct(context.Background(), tst)
		So(err, ShouldBeNil)
	})

	Convey(`test validator.RegisterAlias`, t, func() {
		v := New()
		v.RegisterAlias("iscoloralias", "hexcolor|rgb|rgba|hsl|hsla")

		s := "rgb(255,255,255)"
		errs := v.Var(context.Background(), s, "iscoloralias")
		So(errs, ShouldBeNil)

		s = ""
		errs = v.Var(context.Background(), s, "omitempty,iscoloralias")
		So(errs, ShouldBeNil)

		s = "rgb(255,255,0)"
		errs = v.Var(context.Background(), s, "iscoloralias,len=5")
		So(errs, ShouldNotBeNil)
		So(errs.Error(), ShouldEqual, "Key: '' Error:Field validation for '' failed on the 'len' tag")

		type Test struct {
			Color string `validate:"iscoloralias"`
		}

		tst := &Test{
			Color: "#000",
		}

		errs = v.Struct(context.Background(), tst)
		So(errs, ShouldBeNil)

		tst.Color = "cfvre"
		errs = v.Struct(context.Background(), tst)
		So(errs, ShouldNotBeNil)
		So(errs.Error(), ShouldEqual, "Key: 'Test.Color' Error:Field validation for 'Color' failed on the 'iscoloralias' tag")

		v.RegisterAlias("req", "required,dive,iscoloralias")
		arr := []string{"val1", "#fff", "#000"}

		errs = v.Var(context.Background(), arr, "req")
		So(errs, ShouldNotBeNil)
		So(errs.Error(), ShouldEqual, "Key: '[0]' Error:Field validation for '[0]' failed on the 'iscoloralias' tag")

		So(func() { v.RegisterAlias("exists!", "gt=5,lt=10") }, ShouldPanic)
	})

	Convey(`test validator.RegisterValidation`, t, func() {
		var tag string

		type Test struct {
			String string `validate:"mytag"`
		}

		val := New()
		_ = val.RegisterValidation("mytag", func(ctx context.Context, fl FieldLevel) bool {
			tag = fl.GetTag()
			return true
		})

		var test Test
		errs := val.Struct(context.Background(), test)
		So(errs, ShouldBeNil)
		So(tag, ShouldEqual, "mytag")
	})

	Convey(`test validator.RegisterStructValidation`, t, func() {
		v := New()
		v.RegisterStructValidation(StructLevelInvalidError, StructLevelInvalidErr{})

		var test StructLevelInvalidErr

		err := v.Struct(context.Background(), test)
		So(err, ShouldNotBeNil)

		errs, ok := err.(ValidationErrors)
		So(ok, ShouldBeTrue)

		fe := errs[0]
		So(fe.Field(), ShouldEqual, "Value")
		So(fe.StructField(), ShouldEqual, "Value")
		So(fe.Namespace(), ShouldEqual, "StructLevelInvalidErr.Value")
		So(fe.StructNamespace(), ShouldEqual, "StructLevelInvalidErr.Value")
		So(fe.Tag(), ShouldEqual, "required")
		So(fe.ActualTag(), ShouldEqual, "required")
		So(fe.Kind(), ShouldEqual, reflect.Invalid)
		So(fe.Type(), ShouldEqual, reflect.TypeOf(nil))
	})

	Convey(`test validator.RegisterStructValidationMapRules`, t, func() {
		type Data struct {
			Name string
			Age  uint32
		}

		data := Data{
			Name: "leo",
			Age:  1000,
		}

		rules := map[string]string{
			"Name": "min=4,max=6",
			"Age":  "min=4,max=6",
		}

		v := New()
		v.RegisterStructValidationMapRules(rules, Data{})

		err := v.Struct(context.Background(), data)
		errs, ok := err.(ValidationErrors)
		So(ok, ShouldBeTrue)
		So(len(errs), ShouldEqual, 2)
		So(errs[0].Error(), ShouldEqual, "Key: 'Data.Name' Error:Field validation for 'Name' failed on the 'min' tag")
		So(errs[1].Error(), ShouldEqual, "Key: 'Data.Age' Error:Field validation for 'Age' failed on the 'max' tag")
	})

	Convey(`test validator.RegisterCustomTypeFunc`, t, func() {
		v := New()
		v.RegisterCustomTypeFunc(ValidateValuerType, valuer{}, (*driver.Valuer)(nil), sql.NullString{}, sql.NullInt64{}, sql.NullBool{}, sql.NullFloat64{})
		v.RegisterCustomTypeFunc(ValidateCustomType, MadeUpCustomType{})
		v.RegisterCustomTypeFunc(OverrideIntTypeForSomeReason, 1)

		val := valuer{
			Name: "",
		}

		errs := v.Var(context.Background(), val, "required")
		So(errs, ShouldNotBeNil)
		So(errs.Error(), ShouldEqual, "Key: '' Error:Field validation for '' failed on the 'required' tag")

		val.Name = "Valid Name"
		errs = v.Var(context.Background(), val, "required")
		So(errs, ShouldBeNil)

		val.Name = "errorme"
		errs = v.Var(context.Background(), val, "required")
		So(errs, ShouldNotBeNil)
		So(errs.Error(), ShouldEqual, "validate panic with val: {errorme}, tag: required, reason: SQL Driver Valuer error: some kind of error")

		myVal := valuer{
			Name: "",
		}

		errs = v.Var(context.Background(), myVal, "required")
		So(errs, ShouldNotBeNil)
		So(errs.Error(), ShouldEqual, "Key: '' Error:Field validation for '' failed on the 'required' tag")

		cust := MadeUpCustomType{
			FirstName: "Joey",
			LastName:  "Bloggs",
		}

		c := CustomMadeUpStruct{MadeUp: cust, OverriddenInt: 2}

		errs = v.Struct(context.Background(), c)
		So(errs, ShouldBeNil)

		c.MadeUp.FirstName = ""
		c.OverriddenInt = 1

		err := v.Struct(context.Background(), c)
		So(err, ShouldNotBeNil)
		errs0, ok := err.(ValidationErrors)
		So(ok, ShouldBeTrue)
		So(len(errs0), ShouldEqual, 2)
		So(errs0[0].Error(), ShouldEqual, "Key: 'CustomMadeUpStruct.MadeUp' Error:Field validation for 'MadeUp' failed on the 'required' tag")
		So(errs0[1].Error(), ShouldEqual, "Key: 'CustomMadeUpStruct.OverriddenInt' Error:Field validation for 'OverriddenInt' failed on the 'gt' tag")
	})

	Convey(`test validator.RegisterTranslation`, t, func() {
		en0 := en.New()
		uni := ut.New(en0, en0, fr.New())

		trans, _ := uni.GetTranslator("en")
		fr0, _ := uni.GetTranslator("fr")

		v := New()
		err := v.RegisterTranslation("required", trans,
			func(ut ut.Translator) (err error) {
				// using this stype because multiple translation may have to be added for the full translation
				if err = ut.Add("required", "{0} is a required field", false); err != nil {
					return
				}

				return
			}, func(ut ut.Translator, fe FieldError) string {
				t0, err := ut.T(fe.Tag(), fe.Field())
				if err != nil {
					fmt.Printf("warning: error translating FieldError: %#v", fe)
					return fe.Error()
				}

				return t0
			})
		So(err, ShouldBeNil)

		err = v.RegisterTranslation("required", fr0,
			func(ut ut.Translator) (err error) {
				// using this stype because multiple translation may have to be added for the full translation
				if err = ut.Add("required", "{0} est un champ obligatoire", false); err != nil {
					return
				}

				return
			}, func(ut ut.Translator, fe FieldError) string {
				t0, transErr := ut.T(fe.Tag(), fe.Field())
				if transErr != nil {
					fmt.Printf("warning: error translating FieldError: %#v", fe)
					return fe.Error()
				}

				return t0
			})

		So(err, ShouldBeNil)

		type Test struct {
			Value interface{} `validate:"required"`
		}

		var test Test

		err = v.Struct(context.Background(), test)
		So(err, ShouldNotBeNil)

		errs := err.(ValidationErrors)
		So(len(errs), ShouldEqual, 1)

		fe := errs[0]
		So(fe.Tag(), ShouldEqual, "required")
		So(fe.Namespace(), ShouldEqual, "Test.Value")
		So(fe.Translate(trans), ShouldEqual, fmt.Sprintf("%s is a required field", fe.Field()))
		So(fe.Translate(fr0), ShouldEqual, fmt.Sprintf("%s est un champ obligatoire", fe.Field()))

		nl0 := nl.New()
		uni2 := ut.New(nl0, nl0)
		trans2, _ := uni2.GetTranslator("nl")
		So(fe.Translate(trans2), ShouldEqual, "Key: 'Test.Value' Error:Field validation for 'Value' failed on the 'required' tag")

		terrs := errs.Translate(trans)
		So(len(terrs), ShouldEqual, 1)

		v0, ok := terrs["Test.Value"]
		So(ok, ShouldBeTrue)
		So(v0, ShouldEqual, fmt.Sprintf("%s is a required field", fe.Field()))

		terrs = errs.Translate(fr0)
		So(len(terrs), ShouldEqual, 1)

		v0, ok = terrs["Test.Value"]
		So(ok, ShouldBeTrue)
		So(v0, ShouldEqual, fmt.Sprintf("%s est un champ obligatoire", fe.Field()))

		type Test2 struct {
			Value string `validate:"gt=1"`
		}

		var t2 Test2

		err = v.Struct(context.Background(), t2)
		So(err, ShouldNotBeNil)

		errs = err.(ValidationErrors)
		So(len(errs), ShouldEqual, 1)

		fe = errs[0]
		So(fe.Tag(), ShouldEqual, "gt")
		So(fe.Namespace(), ShouldEqual, "Test2.Value")
		So(fe.Translate(trans), ShouldEqual, "Key: 'Test2.Value' Error:Field validation for 'Value' failed on the 'gt' tag")
	})

	Convey(`test validator.Struct`, t, func() {
		var ctxVal string

		fnCtx := func(ctx context.Context, fl FieldLevel) bool {
			ctxVal = ctx.Value(&ctxVal).(string)
			return true
		}

		var ctxSlVal string
		slFn := func(ctx context.Context, sl StructLevel) {
			ctxSlVal = ctx.Value(&ctxSlVal).(string)
		}

		type Test struct {
			Field string `validate:"val"`
		}

		var tst Test

		v := New()
		err := v.RegisterValidation("val", fnCtx)
		So(err, ShouldBeNil)

		v.RegisterStructValidation(slFn, Test{})

		ctx := context.WithValue(context.Background(), &ctxVal, "testval")
		ctx = context.WithValue(ctx, &ctxSlVal, "slVal")
		errs := v.Struct(ctx, tst)
		So(errs, ShouldBeNil)
		So(ctxVal, ShouldEqual, "testval")
		So(ctxSlVal, ShouldEqual, "slVal")
	})

	Convey(`test validator.StructFiltered`, t, func() {
		p1 := func(ns []byte) bool {
			if bytes.HasSuffix(ns, []byte("NoTag")) || bytes.HasSuffix(ns, []byte("Required")) {
				return false
			}

			return true
		}

		p2 := func(ns []byte) bool {
			if bytes.HasSuffix(ns, []byte("SubSlice[0].Test")) ||
				bytes.HasSuffix(ns, []byte("SubSlice[0]")) ||
				bytes.HasSuffix(ns, []byte("SubSlice")) ||
				bytes.HasSuffix(ns, []byte("Sub")) ||
				bytes.HasSuffix(ns, []byte("SubIgnore")) ||
				bytes.HasSuffix(ns, []byte("Anonymous")) ||
				bytes.HasSuffix(ns, []byte("Anonymous.A")) {
				return false
			}

			return true
		}

		p3 := func(ns []byte) bool {
			return !bytes.HasSuffix(ns, []byte("SubTest.Test"))
		}

		// p4 := []string{
		// 	"A",
		// }

		tPartial := &TestPartial{
			NoTag:    "NoTag",
			Required: "Required",

			SubSlice: []*SubTest{
				{

					Test: "Required",
				},
				{

					Test: "Required",
				},
			},

			Sub: &SubTest{
				Test: "1",
			},
			SubIgnore: &SubTest{
				Test: "",
			},
			Anonymous: struct {
				A             string     `validate:"required"`
				ASubSlice     []*SubTest `validate:"required,dive"`
				SubAnonStruct []struct {
					Test      string `validate:"required"`
					OtherTest string `validate:"required"`
				} `validate:"required,dive"`
			}{
				A: "1",
				ASubSlice: []*SubTest{
					{
						Test: "Required",
					},
					{
						Test: "Required",
					},
				},

				SubAnonStruct: []struct {
					Test      string `validate:"required"`
					OtherTest string `validate:"required"`
				}{
					{"Required", "RequiredOther"},
					{"Required", "RequiredOther"},
				},
			},
		}

		v := New()

		// the following should all return no errors as everything is valid in
		// the default state
		errs := v.StructFiltered(context.Background(), tPartial, p1)
		So(errs, ShouldBeNil)

		errs = v.StructFiltered(context.Background(), tPartial, p2)
		So(errs, ShouldBeNil)

		// this isn't really a robust test, but is ment to illustrate the ANON CASE below
		errs = v.StructFiltered(context.Background(), tPartial.SubSlice[0], p3)
		So(errs, ShouldBeNil)

		// mod tParial for required feild and re-test making sure invalid fields are NOT required:
		tPartial.Required = ""

		// inversion and retesting Partial to generate failures:
		errs = v.StructFiltered(context.Background(), tPartial, p1)
		So(errs, ShouldNotBeNil)
		So(errs.Error(), ShouldEqual, "Key: 'TestPartial.Required' Error:Field validation for 'Required' failed on the 'required' tag")

		// reset Required field, and set nested struct
		tPartial.Required = "Required"
		tPartial.Anonymous.A = ""

		// will pass as unset feilds is not going to be tested
		errs = v.StructFiltered(context.Background(), tPartial, p1)
		So(errs, ShouldBeNil)

		// will fail as unset feild is tested
		errs = v.StructFiltered(context.Background(), tPartial, p2)
		So(errs, ShouldNotBeNil)
		So(errs.Error(), ShouldEqual, "Key: 'TestPartial.Anonymous.A' Error:Field validation for 'A' failed on the 'required' tag")

		// reset nested struct and unset struct in slice
		tPartial.Anonymous.A = "Required"
		tPartial.SubSlice[0].Test = ""

		// these will pass as unset item is NOT tested
		errs = v.StructFiltered(context.Background(), tPartial, p1)
		So(errs, ShouldBeNil)

		errs = v.StructFiltered(context.Background(), tPartial, p2)
		So(errs, ShouldNotBeNil)
		So(errs.Error(), ShouldEqual, "Key: 'TestPartial.SubSlice[0].Test' Error:Field validation for 'Test' failed on the 'required' tag")
		So(len(errs.(ValidationErrors)), ShouldEqual, 1)

		// Unset second slice member concurrently to test dive behavior:
		tPartial.SubSlice[1].Test = ""

		errs = v.StructFiltered(context.Background(), tPartial, p1)
		So(errs, ShouldBeNil)

		errs = v.StructFiltered(context.Background(), tPartial, p2)
		So(errs, ShouldNotBeNil)
		So(errs.Error(), ShouldEqual, "Key: 'TestPartial.SubSlice[0].Test' Error:Field validation for 'Test' failed on the 'required' tag")
		So(len(errs.(ValidationErrors)), ShouldEqual, 1)

		// reset struct in slice, and unset struct in slice in unset posistion
		tPartial.SubSlice[0].Test = "Required"

		// these will pass as the unset item is NOT tested
		errs = v.StructFiltered(context.Background(), tPartial, p1)
		So(errs, ShouldBeNil)

		errs = v.StructFiltered(context.Background(), tPartial, p2)
		So(errs, ShouldBeNil)

		tPartial.SubSlice[1].Test = "Required"
		tPartial.Anonymous.SubAnonStruct[0].Test = ""

		// these will pass as the unset item is NOT tested
		errs = v.StructFiltered(context.Background(), tPartial, p1)
		So(errs, ShouldBeNil)

		errs = v.StructFiltered(context.Background(), tPartial, p2)
		So(errs, ShouldBeNil)

		dt := time.Now()
		err := v.StructFiltered(context.Background(), &dt, func(ns []byte) bool { return true })
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, "validator: (nil *time.Time)")
	})

	Convey(`test validator.StructPartial/validator.StructExcept`, t, func() {
		p1 := []string{
			"NoTag",
			"Required",
		}

		p2 := []string{
			"SubSlice[0].Test",
			"Sub",
			"SubIgnore",
			"Anonymous.A",
		}

		p3 := []string{
			"SubTest.Test",
		}

		p4 := []string{
			"A",
		}

		tPartial := &TestPartial{
			NoTag:    "NoTag",
			Required: "Required",

			SubSlice: []*SubTest{
				{

					Test: "Required",
				},
				{

					Test: "Required",
				},
			},

			Sub: &SubTest{
				Test: "1",
			},
			SubIgnore: &SubTest{
				Test: "",
			},
			Anonymous: struct {
				A             string     `validate:"required"`
				ASubSlice     []*SubTest `validate:"required,dive"`
				SubAnonStruct []struct {
					Test      string `validate:"required"`
					OtherTest string `validate:"required"`
				} `validate:"required,dive"`
			}{
				A: "1",
				ASubSlice: []*SubTest{
					{
						Test: "Required",
					},
					{
						Test: "Required",
					},
				},

				SubAnonStruct: []struct {
					Test      string `validate:"required"`
					OtherTest string `validate:"required"`
				}{
					{"Required", "RequiredOther"},
					{"Required", "RequiredOther"},
				},
			},
		}

		v := New()

		// the following should all return no errors as everything is valid in
		// the default state
		errs := v.StructPartial(context.Background(), tPartial, p1...)
		So(errs, ShouldBeNil)

		errs = v.StructPartial(context.Background(), tPartial, p2...)
		So(errs, ShouldBeNil)

		// this isn't really a robust test, but is ment to illustrate the ANON CASE below
		errs = v.StructPartial(context.Background(), tPartial.SubSlice[0], p3...)
		So(errs, ShouldBeNil)

		errs = v.StructExcept(context.Background(), tPartial, p1...)
		So(errs, ShouldBeNil)

		errs = v.StructExcept(context.Background(), tPartial, p2...)
		So(errs, ShouldBeNil)

		// mod tParial for required feild and re-test making sure invalid fields are NOT required:
		tPartial.Required = ""

		errs = v.StructExcept(context.Background(), tPartial, p1...)
		So(errs, ShouldBeNil)

		errs = v.StructPartial(context.Background(), tPartial, p2...)
		So(errs, ShouldBeNil)

		// inversion and retesting Partial to generate failures:
		errs = v.StructPartial(context.Background(), tPartial, p1...)
		So(errs, ShouldNotBeNil)
		So(errs.Error(), ShouldEqual, "Key: 'TestPartial.Required' Error:Field validation for 'Required' failed on the 'required' tag")

		errs = v.StructExcept(context.Background(), tPartial, p2...)
		So(errs, ShouldNotBeNil)
		So(errs.Error(), ShouldEqual, "Key: 'TestPartial.Required' Error:Field validation for 'Required' failed on the 'required' tag")

		// reset Required field, and set nested struct
		tPartial.Required = "Required"
		tPartial.Anonymous.A = ""

		// will pass as unset feilds is not going to be tested
		errs = v.StructPartial(context.Background(), tPartial, p1...)
		So(errs, ShouldBeNil)

		errs = v.StructExcept(context.Background(), tPartial, p2...)
		So(errs, ShouldBeNil)

		// ANON CASE the response here is strange, it clearly does what it is being told to
		errs = v.StructExcept(context.Background(), tPartial.Anonymous, p4...)
		So(errs, ShouldBeNil)

		// will fail as unset feild is tested
		errs = v.StructPartial(context.Background(), tPartial, p2...)
		So(errs, ShouldNotBeNil)
		So(errs.Error(), ShouldEqual, "Key: 'TestPartial.Anonymous.A' Error:Field validation for 'A' failed on the 'required' tag")

		errs = v.StructExcept(context.Background(), tPartial, p1...)
		So(errs, ShouldNotBeNil)
		So(errs.Error(), ShouldEqual, "Key: 'TestPartial.Anonymous.A' Error:Field validation for 'A' failed on the 'required' tag")

		// reset nested struct and unset struct in slice
		tPartial.Anonymous.A = "Required"
		tPartial.SubSlice[0].Test = ""

		// these will pass as unset item is NOT tested
		errs = v.StructPartial(context.Background(), tPartial, p1...)
		So(errs, ShouldBeNil)

		errs = v.StructExcept(context.Background(), tPartial, p2...)
		So(errs, ShouldBeNil)

		// these will fail as unset item IS tested
		errs = v.StructExcept(context.Background(), tPartial, p1...)
		So(errs, ShouldNotBeNil)
		So(errs.Error(), ShouldEqual, "Key: 'TestPartial.SubSlice[0].Test' Error:Field validation for 'Test' failed on the 'required' tag")
		So(len(errs.(ValidationErrors)), ShouldEqual, 1)

		errs = v.StructPartial(context.Background(), tPartial, p2...)
		So(errs, ShouldNotBeNil)
		So(errs.Error(), ShouldEqual, "Key: 'TestPartial.SubSlice[0].Test' Error:Field validation for 'Test' failed on the 'required' tag")
		So(len(errs.(ValidationErrors)), ShouldEqual, 1)

		// Unset second slice member concurrently to test dive behavior:
		tPartial.SubSlice[1].Test = ""

		errs = v.StructPartial(context.Background(), tPartial, p1...)
		So(errs, ShouldBeNil)

		// NOTE: When specifying nested items, it is still the users responsibility
		// to specify the dive tag, the library does not override this.
		errs = v.StructExcept(context.Background(), tPartial, p2...)
		So(errs, ShouldNotBeNil)
		So(errs.Error(), ShouldEqual, "Key: 'TestPartial.SubSlice[1].Test' Error:Field validation for 'Test' failed on the 'required' tag")

		errs = v.StructExcept(context.Background(), tPartial, p1...)
		So(errs, ShouldNotBeNil)
		So(len(errs.(ValidationErrors)), ShouldEqual, 2)
		So(errs.(ValidationErrors)[0].Error(), ShouldEqual, "Key: 'TestPartial.SubSlice[0].Test' Error:Field validation for 'Test' failed on the 'required' tag")
		So(errs.(ValidationErrors)[1].Error(), ShouldEqual, "Key: 'TestPartial.SubSlice[1].Test' Error:Field validation for 'Test' failed on the 'required' tag")

		errs = v.StructPartial(context.Background(), tPartial, p2...)
		So(errs, ShouldNotBeNil)
		So(len(errs.(ValidationErrors)), ShouldEqual, 1)
		So(errs.Error(), ShouldEqual, "Key: 'TestPartial.SubSlice[0].Test' Error:Field validation for 'Test' failed on the 'required' tag")

		// reset struct in slice, and unset struct in slice in unset posistion
		tPartial.SubSlice[0].Test = "Required"

		// these will pass as the unset item is NOT tested
		errs = v.StructPartial(context.Background(), tPartial, p1...)
		So(errs, ShouldBeNil)

		errs = v.StructPartial(context.Background(), tPartial, p2...)
		So(errs, ShouldBeNil)

		// testing for missing item by exception, yes it dives and fails
		errs = v.StructExcept(context.Background(), tPartial, p1...)
		So(errs, ShouldNotBeNil)
		So(len(errs.(ValidationErrors)), ShouldEqual, 1)
		So(errs.Error(), ShouldEqual, "Key: 'TestPartial.SubSlice[1].Test' Error:Field validation for 'Test' failed on the 'required' tag")

		errs = v.StructExcept(context.Background(), tPartial, p2...)
		So(errs, ShouldNotBeNil)
		So(errs.Error(), ShouldEqual, "Key: 'TestPartial.SubSlice[1].Test' Error:Field validation for 'Test' failed on the 'required' tag")

		tPartial.SubSlice[1].Test = "Required"

		tPartial.Anonymous.SubAnonStruct[0].Test = ""
		// these will pass as the unset item is NOT tested
		errs = v.StructPartial(context.Background(), tPartial, p1...)
		So(errs, ShouldBeNil)

		errs = v.StructPartial(context.Background(), tPartial, p2...)
		So(errs, ShouldBeNil)

		errs = v.StructExcept(context.Background(), tPartial, p1...)
		So(errs, ShouldNotBeNil)
		So(errs.Error(), ShouldEqual, "Key: 'TestPartial.Anonymous.SubAnonStruct[0].Test' Error:Field validation for 'Test' failed on the 'required' tag")

		errs = v.StructExcept(context.Background(), tPartial, p2...)
		So(errs, ShouldNotBeNil)
		So(errs.Error(), ShouldEqual, "Key: 'TestPartial.Anonymous.SubAnonStruct[0].Test' Error:Field validation for 'Test' failed on the 'required' tag")

		// Test for unnamed struct
		testStruct := &TestStruct{
			String: "test",
		}
		unnamedStruct := struct {
			String string `validate:"required" json:"StringVal"`
		}{String: "test"}
		composedUnnamedStruct := struct{ *TestStruct }{&TestStruct{String: "test"}}

		errs = v.StructPartial(context.Background(), testStruct, "String")
		So(errs, ShouldBeNil)

		errs = v.StructPartial(context.Background(), unnamedStruct, "String")
		So(errs, ShouldBeNil)

		errs = v.StructPartial(context.Background(), composedUnnamedStruct, "TestStruct.String")
		So(errs, ShouldBeNil)

		testStruct.String = ""
		errs = v.StructPartial(context.Background(), testStruct, "String")
		So(errs, ShouldNotBeNil)
		So(errs.Error(), ShouldEqual, "Key: 'TestStruct.String' Error:Field validation for 'String' failed on the 'required' tag")

		unnamedStruct.String = ""
		errs = v.StructPartial(context.Background(), unnamedStruct, "String")
		So(errs, ShouldNotBeNil)
		So(errs.Error(), ShouldEqual, "Key: 'String' Error:Field validation for 'String' failed on the 'required' tag")

		composedUnnamedStruct.String = ""
		errs = v.StructPartial(context.Background(), composedUnnamedStruct, "TestStruct.String")
		So(errs, ShouldNotBeNil)
		So(errs.Error(), ShouldEqual, "Key: 'TestStruct.String' Error:Field validation for 'String' failed on the 'required' tag")
	})

	Convey(`test validator.VarWithValue`, t, func() {
		var errs error
		v := New()

		type TimeTest struct {
			Start *time.Time `validate:"required,gt"`
			End   *time.Time `validate:"required,gt,gtfield=Start"`
		}

		now := time.Now()
		start := now.Add(time.Hour * 24)
		end := start.Add(time.Hour * 24)

		timeTest := &TimeTest{
			Start: &start,
			End:   &end,
		}

		errs = v.Struct(context.Background(), timeTest)
		So(errs, ShouldBeNil)

		timeTest = &TimeTest{
			Start: &end,
			End:   &start,
		}

		errs = v.Struct(context.Background(), timeTest)
		So(errs, ShouldNotBeNil)
		So(errs.Error(), ShouldEqual, "Key: 'TimeTest.End' Error:Field validation for 'End' failed on the 'gtfield' tag")

		errs = v.VarWithValue(context.Background(), &end, &start, "gtfield")
		So(errs, ShouldBeNil)

		errs = v.VarWithValue(context.Background(), &start, &end, "gtfield")
		So(errs, ShouldNotBeNil)
		So(errs.Error(), ShouldEqual, "Key: '' Error:Field validation for '' failed on the 'gtfield' tag")

		errs = v.VarWithValue(context.Background(), &end, &start, "gtfield")
		So(errs, ShouldBeNil)

		errs = v.VarWithValue(context.Background(), &timeTest, &end, "gtfield")
		So(errs, ShouldNotBeNil)
		So(errs.Error(), ShouldEqual, "Key: 'TimeTest.End' Error:Field validation for 'End' failed on the 'gtfield' tag")

		errs = v.VarWithValue(context.Background(), "test bigger", "test", "gtfield")
		So(errs, ShouldBeNil)

		// Tests for time.Duration type.

		// -- Validations for variables of time.Duration type.

		errs = v.VarWithValue(context.Background(), time.Hour, time.Hour-time.Minute, "gtfield")
		So(errs, ShouldBeNil)

		errs = v.VarWithValue(context.Background(), time.Hour, time.Hour, "gtfield")
		So(errs, ShouldNotBeNil)
		So(errs.Error(), ShouldEqual, "Key: '' Error:Field validation for '' failed on the 'gtfield' tag")

		errs = v.VarWithValue(context.Background(), time.Hour, time.Hour+time.Minute, "gtfield")
		So(errs, ShouldNotBeNil)
		So(errs.Error(), ShouldEqual, "Key: '' Error:Field validation for '' failed on the 'gtfield' tag")

		errs = v.VarWithValue(context.Background(), time.Duration(0), time.Hour, "omitempty,gtfield")
		So(errs, ShouldBeNil)

		// -- Validations for a struct with time.Duration type fields.

		type TimeDurationTest struct {
			First  time.Duration `validate:"gtfield=Second"`
			Second time.Duration
		}
		var timeDurationTest *TimeDurationTest

		timeDurationTest = &TimeDurationTest{time.Hour, time.Hour - time.Minute}
		errs = v.Struct(context.Background(), timeDurationTest)
		So(errs, ShouldBeNil)

		timeDurationTest = &TimeDurationTest{time.Hour, time.Hour}
		errs = v.Struct(context.Background(), timeDurationTest)
		So(errs, ShouldNotBeNil)
		So(errs.Error(), ShouldEqual, "Key: 'TimeDurationTest.First' Error:Field validation for 'First' failed on the 'gtfield' tag")

		timeDurationTest = &TimeDurationTest{time.Hour, time.Hour + time.Minute}
		errs = v.Struct(context.Background(), timeDurationTest)
		So(errs, ShouldNotBeNil)
		So(errs.Error(), ShouldEqual, "Key: 'TimeDurationTest.First' Error:Field validation for 'First' failed on the 'gtfield' tag")

		type TimeDurationOmitemptyTest struct {
			First  time.Duration `validate:"omitempty,gtfield=Second"`
			Second time.Duration
		}

		timeDurationOmitemptyTest := &TimeDurationOmitemptyTest{time.Duration(0), time.Hour}
		errs = v.Struct(context.Background(), timeDurationOmitemptyTest)
		So(errs, ShouldBeNil)

		// Tests for Ints types.

		type IntTest struct {
			Val1 int `validate:"required"`
			Val2 int `validate:"required,gtfield=Val1"`
		}

		intTest := &IntTest{
			Val1: 1,
			Val2: 5,
		}

		errs = v.Struct(context.Background(), intTest)
		So(errs, ShouldBeNil)

		intTest = &IntTest{
			Val1: 5,
			Val2: 1,
		}

		errs = v.Struct(context.Background(), intTest)
		So(errs, ShouldNotBeNil)
		So(errs.Error(), ShouldEqual, "Key: 'IntTest.Val2' Error:Field validation for 'Val2' failed on the 'gtfield' tag")

		errs = v.VarWithValue(context.Background(), int(5), int(1), "gtfield")
		So(errs, ShouldBeNil)

		errs = v.VarWithValue(context.Background(), int(1), int(5), "gtfield")
		So(errs, ShouldNotBeNil)
		So(errs.Error(), ShouldEqual, "Key: '' Error:Field validation for '' failed on the 'gtfield' tag")

		type UIntTest struct {
			Val1 uint `validate:"required"`
			Val2 uint `validate:"required,gtfield=Val1"`
		}

		uIntTest := &UIntTest{
			Val1: 1,
			Val2: 5,
		}

		errs = v.Struct(context.Background(), uIntTest)
		So(errs, ShouldBeNil)

		uIntTest = &UIntTest{
			Val1: 5,
			Val2: 1,
		}

		errs = v.Struct(context.Background(), uIntTest)
		So(errs, ShouldNotBeNil)
		So(errs.Error(), ShouldEqual, "Key: 'UIntTest.Val2' Error:Field validation for 'Val2' failed on the 'gtfield' tag")

		errs = v.VarWithValue(context.Background(), uint(5), uint(1), "gtfield")
		So(errs, ShouldBeNil)

		errs = v.VarWithValue(context.Background(), uint(1), uint(5), "gtfield")
		So(errs, ShouldNotBeNil)
		So(errs.Error(), ShouldEqual, "Key: '' Error:Field validation for '' failed on the 'gtfield' tag")

		type FloatTest struct {
			Val1 float64 `validate:"required"`
			Val2 float64 `validate:"required,gtfield=Val1"`
		}

		floatTest := &FloatTest{
			Val1: 1,
			Val2: 5,
		}

		errs = v.Struct(context.Background(), floatTest)
		So(errs, ShouldBeNil)

		floatTest = &FloatTest{
			Val1: 5,
			Val2: 1,
		}

		errs = v.Struct(context.Background(), floatTest)
		So(errs, ShouldNotBeNil)
		So(errs.Error(), ShouldEqual, "Key: 'FloatTest.Val2' Error:Field validation for 'Val2' failed on the 'gtfield' tag")

		errs = v.VarWithValue(context.Background(), float32(5), float32(1), "gtfield")
		So(errs, ShouldBeNil)

		errs = v.VarWithValue(context.Background(), float32(1), float32(5), "gtfield")
		So(errs, ShouldNotBeNil)
		So(errs.Error(), ShouldEqual, "Key: '' Error:Field validation for '' failed on the 'gtfield' tag")

		errs = v.VarWithValue(context.Background(), nil, 1, "gtfield")
		So(errs, ShouldNotBeNil)
		So(errs.Error(), ShouldEqual, "Key: '' Error:Field validation for '' failed on the 'gtfield' tag")

		errs = v.VarWithValue(context.Background(), 5, "T", "gtfield")
		So(errs, ShouldNotBeNil)
		So(errs.Error(), ShouldEqual, "Key: '' Error:Field validation for '' failed on the 'gtfield' tag")

		errs = v.VarWithValue(context.Background(), 5, start, "gtfield")
		So(errs, ShouldNotBeNil)
		So(errs.Error(), ShouldEqual, "Key: '' Error:Field validation for '' failed on the 'gtfield' tag")

		type TimeTest2 struct {
			Start *time.Time `validate:"required"`
			End   *time.Time `validate:"required,gtfield=NonExistantField"`
		}

		timeTest2 := &TimeTest2{
			Start: &start,
			End:   &end,
		}

		errs = v.Struct(context.Background(), timeTest2)
		So(errs, ShouldNotBeNil)
		So(errs.Error(), ShouldEqual, "Key: 'TimeTest2.End' Error:Field validation for 'End' failed on the 'gtfield' tag")

		type Other struct {
			Value string
		}

		type Test struct {
			Value Other
			Time  time.Time `validate:"gtfield=Value"`
		}

		tst := Test{
			Value: Other{Value: "StringVal"},
			Time:  end,
		}

		errs = v.Struct(context.Background(), tst)
		So(errs, ShouldNotBeNil)
		So(errs.Error(), ShouldEqual, "Key: 'Test.Time' Error:Field validation for 'Time' failed on the 'gtfield' tag")
	})
}
