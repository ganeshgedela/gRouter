package nats

// ValidateFunc is a function that validates message data.
type ValidateFunc func(data []byte) error

// MapValidator implements the Validator interface using a map of validation functions.
type MapValidator struct {
	validators map[string]ValidateFunc
}

// NewMapValidator creates a new MapValidator.
func NewMapValidator() *MapValidator {
	return &MapValidator{
		validators: make(map[string]ValidateFunc),
	}
}

// Register adds a validation function for a specific message type.
func (v *MapValidator) Register(msgType string, fn ValidateFunc) {
	v.validators[msgType] = fn
}

// Validate checks if the data matches the schema for the given message type.
func (v *MapValidator) Validate(msgType string, data []byte) error {
	fn, ok := v.validators[msgType]
	if !ok {
		// If no validator is registered for this type, we assume it's valid.
		// Alternatively, we could return an error if strict validation is required.
		return nil
	}
	return fn(data)
}

// Ensure MapValidator implements Validator interface.
var _ Validator = (*MapValidator)(nil)
