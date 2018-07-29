package protavobolt_test

import (
	"io/ioutil"
	"os"

	bolt "github.com/coreos/bbolt"
	"github.com/jmalloc/protavo/src/protavo/driver"
	"github.com/jmalloc/protavo/src/protavo/driver/drivertest"
	. "github.com/jmalloc/protavo/src/protavobolt"
)

func init() {
	var path string

	drivertest.Describe(
		"ExclusiveDriver",
		func() (driver.Driver, error) {
			fp, err := ioutil.TempFile("", "protavo-")
			if err != nil {
				return nil, err
			}

			path = fp.Name()

			err = fp.Close()
			if err != nil {
				return nil, err
			}

			err = os.Remove(path)
			if err != nil {
				return nil, err
			}

			db, err := bolt.Open(
				path,
				0600,
				nil,
			)
			if err != nil {
				return nil, err
			}

			return &ExclusiveDriver{DB: db}, nil
		},
		func() {
			if err := os.Remove(path); err != nil {
				panic(err)
			}
		},
	)
}
