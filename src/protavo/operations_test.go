package protavo_test

import (
	"context"
	"fmt"

	. "github.com/jmalloc/protavo/src/protavo"
	"github.com/jmalloc/protavo/src/protavo/document"
	"github.com/jmalloc/protavo/src/protavobolt"
)

func ExampleFetchAll() {
	// First, initialize a database. We'll use the BoltDB driver for examples.
	db, err := protavobolt.OpenTemp(0600, nil)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Next, save some documents so we have something to fetch.
	if err := db.Save(
		context.Background(),
		&document.Document{
			ID:      "person:1",
			Content: document.StringContent("Alice"),
		},
		&document.Document{
			ID:      "person:2",
			Content: document.StringContent("Bob"),
		},
	); err != nil {
		panic(err)
	}

	// Next, we use FetchAll to (rather inefficiently) count the total number of
	// people in our database.
	count := 0

	if err := db.Read(
		context.Background(),
		FetchAll(
			func(doc *document.Document) (bool, error) {
				count++
				return true, nil
			},
		),
	); err != nil {
		panic(err)
	}

	fmt.Printf("found %d documents\n", count)

	// Output: found 2 documents
}

func ExampleFetchWhere() {
	// First, initialize a database. We'll use the BoltDB driver for examples.
	db, err := protavobolt.OpenTemp(0600, nil)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Next, save some documents so we have something to fetch.
	if err := db.Save(
		context.Background(),
		&document.Document{
			ID: "person:1",
			Keys: map[string]document.KeyType{
				"hobby:cycling": document.SharedKey,
			},
			Content: document.StringContent("Alice"),
		},
		&document.Document{
			ID: "person:2",
			Keys: map[string]document.KeyType{
				"hobby:cycling": document.SharedKey,
				"hobby:origami": document.SharedKey,
			},
			Content: document.StringContent("Bob"),
		},
		&document.Document{
			ID: "person:3",
			Keys: map[string]document.KeyType{
				"hobby:origami": document.SharedKey,
			},
			Content: document.StringContent("Carlos"),
		},
	); err != nil {
		panic(err)
	}

	// Next, we use FetchWhere to (rather inefficiently) count the number of
	// people that enjoy Origami.
	count := 0

	if err := db.Read(
		context.Background(),
		FetchWhere(
			func(doc *document.Document) (bool, error) {
				count++
				return true, nil
			},
			HasKeys("hobby:origami"),
		),
	); err != nil {
		panic(err)
	}

	fmt.Printf("found %d documents\n", count)

	// Output: found 2 documents
}
