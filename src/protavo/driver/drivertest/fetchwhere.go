package drivertest

import (
	"context"
	"errors"

	"github.com/jmalloc/protavo/src/protavo"
	"github.com/jmalloc/protavo/src/protavo/document"
	g "github.com/onsi/ginkgo"
	m "github.com/onsi/gomega"
)

// describeFetchWhere defines the standard test suite for the
// protavo.FetchWhere() operation.
func describeFetchWhere(
	before func() (*protavo.DB, error),
	after func(),
) {
	ctx := context.Background()
	var doc1, doc2, doc3 *document.Document

	g.Describe("FetchWhere", func() {
		var db *protavo.DB

		g.BeforeEach(func() {
			var err error
			db, err = before()
			m.Expect(err).ShouldNot(m.HaveOccurred())

			doc1 = &document.Document{
				ID:      "doc-1",
				Content: document.StringContent("content-1"),
				Keys:    document.SharedKeys("foo"),
			}

			doc2 = &document.Document{
				ID:      "doc-2",
				Content: document.StringContent("content-2"),
				Keys:    document.SharedKeys("foo", "bar"),
			}

			doc3 = &document.Document{
				ID:      "doc-3",
				Content: document.StringContent("content-3"),
				Keys:    document.SharedKeys("bar"),
			}
		})

		g.AfterEach(func() {
			_ = db.Close()

			if after != nil {
				after()
			}
		})

		g.When("there are no documents in the database", func() {
			g.It("does not return an error", func() {
				op := protavo.FetchWhere(
					nil, // should never be invoked
				)

				err := db.Read(ctx, op)
				m.Expect(err).ShouldNot(m.HaveOccurred())
			})
		})

		g.When("there are documents in the database", func() {
			g.BeforeEach(func() {
				err := db.Save(ctx, doc1, doc2, doc3)
				m.Expect(err).ShouldNot(m.HaveOccurred())
			})

			g.It("does not return an error", func() {
				op := protavo.FetchWhere(
					func(*document.Document) (bool, error) {
						return false, nil
					},
				)

				err := db.Read(ctx, op)
				m.Expect(err).ShouldNot(m.HaveOccurred())
			})

			g.It("does not match any documents if called without conditions", func() {
				op := protavo.FetchWhere(
					func(*document.Document) (bool, error) {
						return false, errors.New("each-func was called unexpectedly")
					},
				)

				err := db.Read(ctx, op)
				m.Expect(err).ShouldNot(m.HaveOccurred())
			})

			g.It("calls the each-func for each of the documents in the store", func() {
				docs := map[string]*document.Document{}

				op := protavo.FetchWhere(
					func(doc *document.Document) (bool, error) {
						docs[doc.ID] = doc
						return true, nil
					},
					protavo.HasKeys("foo", "bar"),
				)

				err := db.Read(ctx, op)
				m.Expect(err).ShouldNot(m.HaveOccurred())
				m.Expect(docs).To(m.HaveLen(1))

				m.Expect(docs).To(m.HaveKey("doc-2"))
				m.Expect(
					docs["doc-2"].Equal(doc2),
				).To(m.BeTrue())
			})

			g.It("stops iterating if the each-func returns false", func() {
				count := 0

				op := protavo.FetchWhere(
					func(doc *document.Document) (bool, error) {
						count++
						return false, nil
					},
					protavo.HasKeys("foo"),
				)

				err := db.Read(ctx, op)
				m.Expect(err).ShouldNot(m.HaveOccurred())
				m.Expect(count).To(m.Equal(1))
			})

			g.It("stops iterating if the each-func returns a non-nil error", func() {
				count := 0
				expected := errors.New("<error>")

				op := protavo.FetchWhere(
					func(doc *document.Document) (bool, error) {
						count++
						return false, expected
					},
					protavo.HasKeys("foo"),
				)

				err := db.Read(ctx, op)
				m.Expect(err).To(m.MatchError(expected))
				m.Expect(count).To(m.Equal(1))
			})
		})
	})
}
