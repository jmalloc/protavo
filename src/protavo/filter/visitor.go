package filter

// Visitor an interface for performing operations with filter conditions.
type Visitor interface {
	IsOneOf(*IsOneOf) (bool, error)
	HasUniqueKeyIn(*HasUniqueKeyIn) (bool, error)
	HasKeys(*HasKeys) (bool, error)
}
