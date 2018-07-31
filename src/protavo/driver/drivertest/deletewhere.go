package drivertest

import (
	"context"
	"errors"

	"github.com/jmalloc/protavo/src/protavo"
	"github.com/jmalloc/protavo/src/protavo/document"
	g "github.com/onsi/ginkgo"
	m "github.com/onsi/gomega"
)

// describeDeleteWhere defines the standard test suite for the
// protavo.DeleteWhere() operation.
func describeDeleteWhere(
	before func() (*protavo.DB, error),
	after func(),
) {
	ctx := context.Background()
	var doc1, doc2, doc3 *document.Document

	g.Describe("DeleteWhere", func() {
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
				Keys: document.Keys{
					"uniq": document.UniqueKey,
					"foo":  document.SharedKey,
					"bar":  document.SharedKey,
				},
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
				op := protavo.DeleteWhere(
					nil, // each-func is optional
				)

				err := db.Write(ctx, op)
				m.Expect(err).ShouldNot(m.HaveOccurred())
			})
		})

		g.When("there are documents in the database", func() {
			g.BeforeEach(func() {
				err := db.Save(ctx, doc1, doc2, doc3)
				m.Expect(err).ShouldNot(m.HaveOccurred())
			})

			g.It("does not return an error", func() {
				op := protavo.DeleteWhere(nil)

				err := db.Write(ctx, op)
				m.Expect(err).ShouldNot(m.HaveOccurred())
			})

			g.It("deletes the matching documents", func() {
				op := protavo.DeleteWhere(
					nil,
					protavo.HasKeys("foo", "bar"),
				)

				err := db.Write(ctx, op)
				m.Expect(err).ShouldNot(m.HaveOccurred())

				_, ok, err := db.Load(ctx, "doc-1")
				m.Expect(err).ShouldNot(m.HaveOccurred())
				m.Expect(ok).To(m.BeTrue())

				_, ok, err = db.Load(ctx, "doc-2")
				m.Expect(err).ShouldNot(m.HaveOccurred())
				m.Expect(ok).To(m.BeFalse()) // deleted document

				_, ok, err = db.Load(ctx, "doc-3")
				m.Expect(err).ShouldNot(m.HaveOccurred())
				m.Expect(ok).To(m.BeTrue())
			})

			g.It("does not delete any documents if called without conditions", func() {
				op := protavo.DeleteWhere(nil)

				err := db.Write(ctx, op)
				m.Expect(err).ShouldNot(m.HaveOccurred())

				docs, err := db.LoadMany(ctx, "doc-1", "doc-2", "doc-3")
				m.Expect(err).ShouldNot(m.HaveOccurred())
				m.Expect(docs).To(m.HaveLen(3))
			})

			g.When("an each-func is provided", func() {
				g.It("deletes the matching documents", func() {
					op := protavo.DeleteWhere(
						nil,
						protavo.HasKeys("foo", "bar"),
					)

					err := db.Write(ctx, op)
					m.Expect(err).ShouldNot(m.HaveOccurred())

					docs, err := db.LoadMany(ctx, "doc-1", "doc-2", "doc-3")
					m.Expect(err).ShouldNot(m.HaveOccurred())
					m.Expect(docs).To(m.HaveLen(2))
				})

				g.It("does not delete any documents if called without conditions", func() {
					op := protavo.DeleteWhere(
						func(id string) error {
							return nil
						},
					)

					err := db.Write(ctx, op)
					m.Expect(err).ShouldNot(m.HaveOccurred())

					docs, err := db.LoadMany(ctx, "doc-1", "doc-2", "doc-3")
					m.Expect(err).ShouldNot(m.HaveOccurred())
					m.Expect(docs).To(m.HaveLen(3))
				})

				g.It("calls the each-func for each of the matching documents", func() {
					ids := map[string]struct{}{}

					op := protavo.DeleteWhere(
						func(id string) error {
							ids[id] = struct{}{}
							return nil
						},
						protavo.HasKeys("foo", "bar"),
					)

					err := db.Write(ctx, op)
					m.Expect(err).ShouldNot(m.HaveOccurred())
					m.Expect(ids).To(m.HaveLen(1))
					m.Expect(ids).To(m.HaveKey("doc-2"))
				})

				g.It("stops iterating if the each-func returns a non-nil error", func() {
					count := 0
					expected := errors.New("<error>")

					op := protavo.DeleteWhere(
						func(string) error {
							count++
							return expected
						},
						protavo.HasKeys("foo"),
					)

					err := db.Write(ctx, op)
					m.Expect(err).To(m.MatchError(expected))
					m.Expect(count).To(m.Equal(1))
				})
			})

			g.It("allows other documents to use the deleted document's unique key", func() {
				op := protavo.DeleteWhere(
					nil,
					protavo.HasKeys("foo", "bar"),
				)

				err := db.Write(ctx, op)
				m.Expect(err).ShouldNot(m.HaveOccurred())

				doc1.Keys = document.UniqueKeys("uniq")
				err = db.Save(ctx, doc1)
				m.Expect(err).ShouldNot(m.HaveOccurred())
			})
		})
	})
}
