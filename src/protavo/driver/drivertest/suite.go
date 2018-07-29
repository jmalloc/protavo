package drivertest

import (
	"github.com/jmalloc/protavo/src/protavo/driver"
)

// Describe defines the standard test suite for an implementation of
// protavo.Driver.
func Describe(
	name string,
	before func() (driver.Driver, error),
	after func(),
) {
	describeFetchAll(before, after)
}
