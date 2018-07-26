package protobolt

import "fmt"

// OptimisticLockError is an error that occurs when attempt to modify a document
// fails because the incorrect document version was provided with the request.
type OptimisticLockError struct {
	DocumentID    string
	GivenVersion  uint64
	ActualVersion uint64
	Action        string
}

func (e *OptimisticLockError) Error() string {
	return fmt.Sprintf(
		"optimistic lock failure attempting to %s '%s', %d != %d",
		e.Action,
		e.DocumentID,
		e.GivenVersion,
		e.ActualVersion,
	)
}

// IsOptimisticLockError returns true if err represents an optimistic lock
// failure.
func IsOptimisticLockError(err error) bool {
	_, ok := err.(*OptimisticLockError)
	return ok
}

// DuplicateKeyError is an error that occurs when attempt is made to save a
// document with a unique key that is already used by a different document.
type DuplicateKeyError struct {
	DocumentID            string
	ConflictingDocumentID string
	UniqueKey             string
}

func (e *DuplicateKeyError) Error() string {
	return fmt.Sprintf(
		"cannot save '%s', unique key '%s' conflicts with '%s'",
		e.DocumentID,
		e.ConflictingDocumentID,
		e.UniqueKey,
	)
}
