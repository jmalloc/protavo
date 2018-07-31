package drivertest

import (
	"context"

	"github.com/jmalloc/protavo/src/protavo"
	"github.com/jmalloc/protavo/src/protavo/document"
	"github.com/jmalloc/protavo/src/protavo/filter"
	g "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/extensions/table"
	m "github.com/onsi/gomega"
)

// describeFilters defines the standard test suite for Protavo filters.
func describeFilters(
	before func() (*protavo.DB, error),
	after func(),
) {
	ctx := context.Background()

	g.Describe("Filters", func() {
		var db *protavo.DB

		g.BeforeEach(func() {
			var err error
			db, err = before()
			m.Expect(err).ShouldNot(m.HaveOccurred())

			err = db.Save(
				ctx,
				&document.Document{
					ID: "doc-1",
					Keys: document.Keys{
						"uniq-1": document.UniqueKey,
						"shar-a": document.SharedKey,
					},
					Content: document.StringContent(""),
				},
				&document.Document{
					ID: "doc-2",
					Keys: document.Keys{
						"uniq-2": document.UniqueKey,
						"shar-a": document.SharedKey,
						"shar-b": document.SharedKey,
					},
					Content: document.StringContent(""),
				},
				&document.Document{
					ID: "doc-3",
					Keys: document.Keys{
						"uniq-3": document.UniqueKey,
						"shar-a": document.SharedKey,
						"shar-b": document.SharedKey,
						"shar-c": document.SharedKey,
					},
					Content: document.StringContent(""),
				},
				&document.Document{
					ID: "doc-4",
					Keys: document.Keys{
						"uniq-4": document.UniqueKey,
						"shar-a": document.SharedKey,
						"shar-b": document.SharedKey,
						"shar-c": document.SharedKey,
						"shar-d": document.SharedKey,
					},
					Content: document.StringContent(""),
				},
				&document.Document{
					ID: "doc-5",
					Keys: document.Keys{
						"uniq-5": document.UniqueKey,
						"shar-a": document.SharedKey,
						"shar-b": document.SharedKey,
						"shar-c": document.SharedKey,
						"shar-d": document.SharedKey,
						"shar-e": document.SharedKey,
					},
					Content: document.StringContent(""),
				},
			)

			m.Expect(err).ShouldNot(m.HaveOccurred())
		})

		g.AfterEach(func() {
			_ = db.Close()

			if after != nil {
				after()
			}
		})

		entries := []table.TableEntry{
			table.Entry(
				"IsOneOf",
				[]filter.Condition{
					protavo.IsOneOf("doc-1", "doc-3", "non-existent"),
				},
				[]string{"doc-1", "doc-3"},
			),
			table.Entry(
				"HasUniqueKeyIn",
				[]filter.Condition{
					protavo.HasUniqueKeyIn("uniq-1", "uniq-3", "non-existent"),
				},
				[]string{"doc-1", "doc-3"},
			),
			table.Entry(
				"HasKeys",
				[]filter.Condition{
					protavo.HasKeys("shar-a", "shar-b", "shar-c"),
				},
				[]string{"doc-3", "doc-4", "doc-5"},
			),
			table.Entry(
				"Everything",
				[]filter.Condition{
					protavo.IsOneOf("doc-3", "doc-4", "doc-5"),           // start with docs 3, 4, and 5
					protavo.HasUniqueKeyIn("uniq-2", "uniq-3", "uniq-4"), // filter down to 3 and 4
					protavo.HasKeys("shar-b", "shar-c", "shar-d"),        // filter down to 4
				},
				[]string{"doc-4"},
			),

			// The following entries test combinations of different constraints.
			//
			// Note that one of the conditions is always "more constrainted" than the
			// others. This an attempt to his the different "query strategies" used by
			// the BoltDB implementation, and any other implementations that may change
			// strategies dynamically. The order of the conditions in the slice does not
			// matter.

			// IsOneOf first ...
			table.Entry(
				"IsOneOf, then HasUniqueKeyIn",
				[]filter.Condition{
					protavo.IsOneOf("doc-1", "doc-4"),
					protavo.HasUniqueKeyIn("uniq-1", "uniq-2", "uniq-3"),
				},
				[]string{"doc-1"},
			),
			table.Entry(
				"IsOneOf, then HasKeys",
				[]filter.Condition{
					protavo.IsOneOf("doc-1", "doc-4"),
					protavo.HasKeys("shar-a", "shar-b", "shar-c"),
				},
				[]string{"doc-4"},
			),

			// HasUniqueKeyIn first ...
			table.Entry(
				"HasUniqueKeyIn, then IsOneOf",
				[]filter.Condition{
					protavo.HasUniqueKeyIn("uniq-1", "uniq-4"),
					protavo.IsOneOf("doc-1", "doc-2", "doc-3"),
				},
				[]string{"doc-1"},
			),
			table.Entry(
				"HasUniqueKeyIn, then HasKeys",
				[]filter.Condition{
					protavo.HasUniqueKeyIn("uniq-1", "uniq-3"),
					protavo.HasKeys("shar-a", "shar-b", "shar-c"),
				},
				[]string{"doc-3"},
			),

			// HasKeys first ...
			table.Entry(
				"HasKeys, then IsOneOf",
				[]filter.Condition{
					protavo.HasKeys("shar-a", "shar-b"),
					protavo.IsOneOf("doc-1", "doc-2", "doc-3"),
				},
				[]string{"doc-2", "doc-3"},
			),
			table.Entry(
				"HasKeys, then HasUniqueKeyIn",
				[]filter.Condition{
					protavo.HasKeys("shar-a", "shar-b"),
					protavo.HasUniqueKeyIn("uniq-1", "uniq-2", "uniq-3"),
				},
				[]string{"doc-2", "doc-3"},
			),
		}

		table.DescribeTable(
			"FetchWhere",
			func(
				c []filter.Condition,
				expectedIDs []string,
			) {
				docs, err := db.LoadManyWhere(ctx, c...)
				m.Expect(err).ShouldNot(m.HaveOccurred())

				ids := make([]string, len(docs))
				for i, doc := range docs {
					ids[i] = doc.ID
				}

				m.Expect(ids).To(m.ConsistOf(expectedIDs))
			},
			entries...,
		)

		table.DescribeTable(
			"DeleteWhere",
			func(
				c []filter.Condition,
				expectedIDs []string,
			) {
				ids, err := db.DeleteWhere(ctx, c...)
				m.Expect(err).ShouldNot(m.HaveOccurred())
				m.Expect(ids).To(m.ConsistOf(expectedIDs))
			},
			entries...,
		)
	})
}
