package domain

import "errors"

// Sentinel errors expressing failure categories. The HTTP layer maps each to a
// status code (see httpapi.writeError), keeping transport concerns out of the
// domain and service layers.
var (
	// ErrNotFound indicates a requested entity does not exist.
	ErrNotFound = errors.New("not found")
	// ErrValidation indicates the caller supplied invalid input.
	ErrValidation = errors.New("validation failed")
	// ErrConflict indicates the operation conflicts with current state, e.g. a
	// duplicate name or a cyclic container parent.
	ErrConflict = errors.New("conflict")
)

// ValidationError carries a human-readable message while still matching
// errors.Is(err, ErrValidation), so callers can branch on the category and
// surface the message.
type ValidationError struct{ Message string }

func (e *ValidationError) Error() string      { return e.Message }
func (e *ValidationError) Is(target error) bool { return target == ErrValidation }

// Invalid constructs a ValidationError with the given message.
func Invalid(message string) error { return &ValidationError{Message: message} }
