package driver

// Result is the result of an operation that does not return a value.
type Result struct {
	Err error
}

// Get returns the result value and error.
// It panics if the operation has not yet been executed.
func (r *Result) Get() error {
	if r == nil {
		panic("operation has not been executed")
	}

	return r.Err
}
