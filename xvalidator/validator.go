package xvalidator

import (
	protovalidate "github.com/bufbuild/protovalidate-go"
	"google.golang.org/protobuf/proto"
)

type (
	// A Option modifies the default configuration of a Validator. See the
	// individual options for their defaults and affects on the fallibility of
	// configuring a Validator.
	Option = protovalidate.ValidatorOption

	// StandardConstraintResolver is responsible for resolving the standard
	// constraints from the provided protoreflect.Descriptor. The default resolver
	// can be intercepted and modified using WithStandardConstraintInterceptor.
	StandardConstraintResolver = protovalidate.StandardConstraintResolver

	// StandardConstraintInterceptor can be provided to
	// WithStandardConstraintInterceptor to allow modifying a
	// StandardConstraintResolver.
	StandardConstraintInterceptor = protovalidate.StandardConstraintInterceptor
)

var (
	// WithUTC specifies whether timestamp operations should use UTC or the OS's
	// local timezone for timestamp related values. By default, the local timezone
	// is used.
	WithUTC = protovalidate.WithUTC

	// WithFailFast specifies whether validation should fail on the first constraint
	// violation encountered or if all violations should be accumulated. By default,
	// all violations are accumulated.
	WithFailFast = protovalidate.WithFailFast

	// WithMessages allows warming up the Validator with messages that are
	// expected to be validated. Messages included transitively (i.e., fields with
	// message values) are automatically handled.
	WithMessages = protovalidate.WithMessages

	// WithDescriptors allows warming up the Validator with message
	// descriptors that are expected to be validated. Messages included transitively
	// (i.e., fields with message values) are automatically handled.
	WithDescriptors = protovalidate.WithDescriptors

	// WithDisableLazy prevents the Validator from lazily building validation logic
	// for a message it has not encountered before. Disabling lazy logic
	// additionally eliminates any internal locking as the validator becomes
	// read-only.
	//
	// Note: All expected messages must be provided by WithMessages or
	// WithDescriptors during initialization.
	WithDisableLazy = protovalidate.WithDisableLazy

	// WithStandardConstraintInterceptor allows intercepting the
	// StandardConstraintResolver used by the Validator to modify or replace it.
	WithStandardConstraintInterceptor = protovalidate.WithStandardConstraintInterceptor
)

type Validator interface {
	Validate(msg proto.Message) error
}

var (
	Default, _ = New()
	Validate   = Default
)

// New creates a Validator with the given options. An error may occur in setting
// up the CEL execution environment if the configuration is invalid. See the
// individual ValidatorOption for how they impact the fallibility of New.
func New(opts ...Option) (Validator, error) { return protovalidate.New(opts...) }
