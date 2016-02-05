package utils

// FirstError returns the first error in a collection of (potential)
// errors or nil.
func FirstError(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}
