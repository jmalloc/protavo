package document

// KeyType is an enumeration of the different types of document keys.
type KeyType int32

const (
	// UniqueKey is the KeyType for keys that are always exclusive to a single document.
	// Uniqye keys are useful for addressing specific documents by some role they
	// fill or property they hold.
	UniqueKey KeyType = 1

	// SharedKey is the KeyType for keys that may be shared by multiple documents.
	// Shared keys are useful for quickly locating sets of related documents.
	SharedKey KeyType = 2
)

// Keys is a map of key name to type.
type Keys map[string]KeyType

// UniqueKeys is a convenience function for creating a key-map consisting of a
// set of unique keys.
func UniqueKeys(keys ...string) Keys {
	m := make(Keys, len(keys))

	for _, k := range keys {
		m[k] = UniqueKey
	}

	return m
}

// SharedKeys is a convenience function for creating a key-map consisting of a
// set of shared keys.
func SharedKeys(keys ...string) Keys {
	m := make(Keys, len(keys))

	for _, k := range keys {
		m[k] = SharedKey
	}

	return m
}
