package validator

import (
	"context"
	ut "github.com/go-playground/universal-translator"
	validator2 "github.com/go-playground/validator/v10"
	"github.com/sandwich-go/boost/xerror"
	"reflect"
)

type (
	// FieldError contains all functions to get error details
	FieldError = validator2.FieldError
	// FieldLevel contains all the information and helper functions
	// to validate a field
	FieldLevel = validator2.FieldLevel
	// StructLevel contains all the information and helper functions
	// to validate a struct
	StructLevel = validator2.StructLevel
	// TagNameFunc allows for adding of a custom tag name parser
	TagNameFunc = validator2.TagNameFunc
	// Func accepts a context.Context and FieldLevel interface for all
	// validation needs. The return value should be true when validation succeeds.
	Func = validator2.FuncCtx
	// StructLevelFunc accepts all values needed for struct level validation
	// but also allows passing of contextual validation information via context.Context.
	StructLevelFunc = validator2.StructLevelFuncCtx
	// CustomTypeFunc allows for overriding or adding custom field type handler functions
	// field = field value of the type to return a value to be validated
	// example Valuer from sql drive see https://golang.org/src/database/sql/driver/types.go?s=1210:1293#L29
	CustomTypeFunc = validator2.CustomTypeFunc
	// RegisterTranslationsFunc allows for registering of translations
	// for a 'ut.Translator' for use within the 'TranslationFunc'
	RegisterTranslationsFunc = validator2.RegisterTranslationsFunc
	// TranslationFunc is the function type used to register or override
	// custom translations
	TranslationFunc = validator2.TranslationFunc
	// FilterFunc is the type used to filter fields using
	// StructFiltered(...) function.
	// returning true results in the field being filtered/skiped from
	// validation
	FilterFunc = validator2.FilterFunc

	// ValidationErrors is an array of FieldError's
	// for use in custom error messages post validation.
	ValidationErrors = validator2.ValidationErrors
)

// Validator validator
type Validator interface {
	// SetTagName allows for changing of the default tag name of Validator
	SetTagName(name string)

	// ValidateMap validates a map using a map of validation rules and allows passing of contextual
	// validation validation information via context.Context.
	ValidateMap(ctx context.Context, data map[string]interface{}, rules map[string]interface{}) map[string]interface{}

	// RegisterTagNameFunc registers a function to get alternate names for StructFields.
	//
	// eg. to use the names which have been specified for JSON representations of structs, rather than normal Go field names:
	//
	//    validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
	//        name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
	//        // skip if tag key says it should be ignored
	//        if name == "-" {
	//            return ""
	//        }
	//        return name
	//    })
	RegisterTagNameFunc(fn TagNameFunc)

	// RegisterAlias registers a mapping of a single validation tag that
	// defines a common or complex set of validation(s) to simplify adding validation
	// to structs.
	//
	// NOTE: this function is not thread-safe it is intended that these all be registered prior to any validation
	RegisterAlias(alias, tags string)

	// RegisterValidation adds a validation with the given tag
	// allowing context.Context validation support.
	//
	// NOTES:
	// - if the key already exists, the previous validation function will be replaced.
	// - this method is not thread-safe it is intended that these all be registered prior to any validation
	RegisterValidation(tag string, fn Func, callValidationEvenIfNull ...bool) error

	// RegisterStructValidation registers a StructLevelFunc against a number of types and allows passing
	// of contextual validation information via context.Context.
	//
	// NOTE:
	// - this method is not thread-safe it is intended that these all be registered prior to any validation
	RegisterStructValidation(fn StructLevelFunc, types ...interface{})

	// RegisterStructValidationMapRules registers validate map rules.
	// Be aware that map validation rules supersede those defined on a/the struct if present.
	//
	// NOTE: this method is not thread-safe it is intended that these all be registered prior to any validation
	RegisterStructValidationMapRules(rules map[string]string, types ...interface{})

	// RegisterCustomTypeFunc registers a CustomTypeFunc against a number of types
	//
	// NOTE: this method is not thread-safe it is intended that these all be registered prior to any validation
	RegisterCustomTypeFunc(fn CustomTypeFunc, types ...interface{})

	// RegisterTranslation registers translations against the provided tag.
	RegisterTranslation(tag string, trans ut.Translator, registerFn RegisterTranslationsFunc, translationFn TranslationFunc) error

	// Struct validates a structs exposed fields, and automatically validates nested structs, unless otherwise specified
	// and also allows passing of context.Context for contextual validation information.
	//
	// It returns InvalidValidationError for bad values passed in and nil or ValidationErrors as error otherwise.
	// You will need to assert the error if it's not nil eg. err.(validator.ValidationErrors) to access the array of errors.
	Struct(ctx context.Context, s interface{}) error

	// StructFiltered validates a structs exposed fields, that pass the FilterFunc check and automatically validates
	// nested structs, unless otherwise specified and also allows passing of contextual validation information via
	// context.Context
	//
	// It returns InvalidValidationError for bad values passed in and nil or ValidationErrors as error otherwise.
	// You will need to assert the error if it's not nil eg. err.(validator.ValidationErrors) to access the array of errors.
	StructFiltered(ctx context.Context, s interface{}, fn FilterFunc) error

	// StructPartial validates the fields passed in only, ignoring all others and allows passing of contextual
	// validation validation information via context.Context
	// Fields may be provided in a namespaced fashion relative to the  struct provided
	// eg. NestedStruct.Field or NestedArrayField[0].Struct.Name
	//
	// It returns InvalidValidationError for bad values passed in and nil or ValidationErrors as error otherwise.
	// You will need to assert the error if it's not nil eg. err.(validator.ValidationErrors) to access the array of errors.
	StructPartial(ctx context.Context, s interface{}, fields ...string) error

	// StructExcept validates all fields except the ones passed in and allows passing of contextual
	// validation validation information via context.Context
	// Fields may be provided in a namespaced fashion relative to the  struct provided
	// i.e. NestedStruct.Field or NestedArrayField[0].Struct.Name
	//
	// It returns InvalidValidationError for bad values passed in and nil or ValidationErrors as error otherwise.
	// You will need to assert the error if it's not nil eg. err.(validator.ValidationErrors) to access the array of errors.
	StructExcept(ctx context.Context, s interface{}, fields ...string) error

	// Var validates a single variable using tag style validation and allows passing of contextual
	// validation validation information via context.Context.
	// eg.
	// var i int
	// validate.Var(i, "gt=1,lt=10")
	//
	// WARNING: a struct can be passed for validation eg. time.Time is a struct or
	// if you have a custom type and have registered a custom type handler, so must
	// allow it; however unforeseen validations will occur if trying to validate a
	// struct that is meant to be passed to 'validate.Struct'
	//
	// It returns InvalidValidationError for bad values passed in and nil or ValidationErrors as error otherwise.
	// You will need to assert the error if it's not nil eg. err.(validator.ValidationErrors) to access the array of errors.
	// validate Array, Slice and maps fields which may contain more than one error
	Var(ctx context.Context, val interface{}, tag string) error

	// VarWithValue validates a single variable, against another variable/field's value using tag style validation and
	// allows passing of contextual validation validation information via context.Context.
	// eg.
	// s1 := "abcd"
	// s2 := "abcd"
	// validate.VarWithValue(s1, s2, "eqcsfield") // returns true
	//
	// WARNING: a struct can be passed for validation eg. time.Time is a struct or
	// if you have a custom type and have registered a custom type handler, so must
	// allow it; however unforeseen validations will occur if trying to validate a
	// struct that is meant to be passed to 'validate.Struct'
	//
	// It returns InvalidValidationError for bad values passed in and nil or ValidationErrors as error otherwise.
	// You will need to assert the error if it's not nil eg. err.(validator.ValidationErrors) to access the array of errors.
	// validate Array, Slice and maps fields which may contain more than one error
	VarWithValue(ctx context.Context, field interface{}, other interface{}, tag string) error
}

var (
	Default  = New()
	Validate = Default
)

type validate struct {
	v *validator2.Validate
}

// New returns a new Validator with sane defaults.
// Validator is designed to be thread-safe and used as a singleton instance.
// It caches information about your struct and validations,
// in essence only parsing your validation tags once per struct type.
// Using multiple instances neglects the benefit of caching.
func New() Validator {
	return &validate{v: validator2.New()}
}

func (v *validate) SetTagName(name string)             { v.v.SetTagName(name) }
func (v *validate) RegisterTagNameFunc(fn TagNameFunc) { v.v.RegisterTagNameFunc(fn) }
func (v *validate) RegisterAlias(alias, tags string)   { v.v.RegisterAlias(alias, tags) }
func (v *validate) RegisterValidation(tag string, fn Func, callValidationEvenIfNull ...bool) error {
	return v.v.RegisterValidationCtx(tag, fn, callValidationEvenIfNull...)
}
func (v *validate) RegisterStructValidation(fn StructLevelFunc, types ...interface{}) {
	v.v.RegisterStructValidationCtx(fn, types...)
}
func (v *validate) RegisterStructValidationMapRules(rules map[string]string, types ...interface{}) {
	v.v.RegisterStructValidationMapRules(rules, types...)
}
func (v *validate) RegisterCustomTypeFunc(fn CustomTypeFunc, types ...interface{}) {
	v.v.RegisterCustomTypeFunc(fn, types...)
}
func (v *validate) RegisterTranslation(tag string, trans ut.Translator, registerFn RegisterTranslationsFunc, translationFn TranslationFunc) error {
	return v.v.RegisterTranslation(tag, trans, registerFn, translationFn)
}
func (v *validate) ValidateMap(ctx context.Context, data map[string]interface{}, rules map[string]interface{}) map[string]interface{} {
	return v.v.ValidateMapCtx(ctx, data, rules)
}
func (v *validate) Var(ctx context.Context, val interface{}, tag string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = xerror.NewText("validate panic with val: %v, tag: %s, reason: %v", val, tag, r)
		}
	}()
	err = v.v.VarCtx(ctx, val, tag)
	return
}
func (v *validate) Struct(ctx context.Context, s interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = xerror.NewText("validate panic with struct: %s, reason: %v", reflect.TypeOf(s).String(), r)
		}
	}()
	err = v.v.StructCtx(ctx, s)
	return
}
func (v *validate) StructFiltered(ctx context.Context, s interface{}, fn FilterFunc) error {
	return v.v.StructFilteredCtx(ctx, s, fn)
}
func (v *validate) StructPartial(ctx context.Context, s interface{}, fields ...string) error {
	return v.v.StructPartialCtx(ctx, s, fields...)
}
func (v *validate) StructExcept(ctx context.Context, s interface{}, fields ...string) error {
	return v.v.StructExceptCtx(ctx, s, fields...)
}
func (v *validate) VarWithValue(ctx context.Context, field interface{}, other interface{}, tag string) error {
	return v.v.VarWithValueCtx(ctx, field, other, tag)
}
