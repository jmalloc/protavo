package drivertest

import "github.com/jmalloc/protavo/src/protavo"

// Describe defines the standard test suite for an implementation of
// protavo.Driver.
func Describe(
	name string,
	before func() (*protavo.DB, error),
	after func(),
) {
	describeFetchAll(before, after)
}
