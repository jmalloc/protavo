package drivertest

import (
	"context"
	"time"

	"github.com/jmalloc/protavo/src/protavo"
	"github.com/jmalloc/protavo/src/protavo/document"
	g "github.com/onsi/ginkgo"
	m "github.com/onsi/gomega"
)

// describeForceSave defines the standard test suite for the protavo.Save()
// operation.
func describeForceSave(
	before func() (*protavo.DB, error),
	after func(),
) {
	ctx := context.Background()
	var doc1, doc2 *document.Document

	g.Describe("ForceSave", func() {
		var db *protavo.DB

		g.BeforeEach(func() {
			var err error
			db, err = before()
			m.Expect(err).ShouldNot(m.HaveOccurred())

			doc1 = &document.Document{
				ID: "doc-1",
				Keys: document.Keys{
					"<unique-key>": document.UniqueKey,
					"<shared-key>": document.SharedKey,
				},
				Headers: document.Headers{
					"<header-key>": "<header-value>",
				},
				Content: document.StringContent("<content-1>"),
			}

			doc2 = &document.Document{
				ID:      "doc-2",
				Content: document.StringContent("<content-2>"),
			}
		})

		g.AfterEach(func() {
			_ = db.Close()

			if after != nil {
				after()
			}
		})

		g.When("creating a new document", func() {
			g.It("persists the document faithfully", func() {
				op := protavo.ForceSave(doc1)

				err := db.Write(ctx, op)
				m.Expect(err).ShouldNot(m.HaveOccurred())

				doc, ok, err := db.Load(ctx, "doc-1")
				m.Expect(err).ShouldNot(m.HaveOccurred())
				m.Expect(ok).To(m.BeTrue())

				m.Expect(
					doc.Equal(doc1),
				).To(m.BeTrue())
			})

			g.It("sets the revision and timestamps on the saved document", func() {
				op := protavo.ForceSave(doc1)

				err := db.Write(ctx, op)
				m.Expect(err).ShouldNot(m.HaveOccurred())

				m.Expect(doc1.Revision).To(m.Equal(uint64(1)))
				m.Expect(doc1.CreatedAt).To(
					m.BeTemporally("~", time.Now(), 50*time.Millisecond),
				)
				m.Expect(doc1.UpdatedAt).To(
					m.BeTemporally("==", doc1.CreatedAt),
				)
			})

			g.It("persists the document, even if the revision is not 0", func() {
				doc1.Revision = 123
				op := protavo.ForceSave(doc1)

				err := db.Write(ctx, op)
				m.Expect(err).ShouldNot(m.HaveOccurred())

				_, ok, err := db.Load(ctx, "doc-1")
				m.Expect(err).ShouldNot(m.HaveOccurred())
				m.Expect(ok).To(m.BeTrue())
			})
		})

		g.When("updating a document", func() {
			g.BeforeEach(func() {
				err := db.Save(ctx, doc1)
				m.Expect(err).ShouldNot(m.HaveOccurred())

				// modify the document headers and content
				doc1.Headers["<header-key>"] = "<updated-header-value>"
				doc1.Content = document.StringContent("<updated-content>")
			})

			g.It("persists the document faithfully", func() {
				op := protavo.ForceSave(doc1)

				err := db.Write(ctx, op)
				m.Expect(err).ShouldNot(m.HaveOccurred())

				doc, ok, err := db.Load(ctx, "doc-1")
				m.Expect(err).ShouldNot(m.HaveOccurred())
				m.Expect(ok).To(m.BeTrue())

				m.Expect(
					doc.Equal(doc1),
				).To(m.BeTrue())
			})

			g.It("sets the revision and timestamps on the saved document", func() {
				createdAt := doc1.CreatedAt

				op := protavo.ForceSave(doc1)

				err := db.Write(ctx, op)
				m.Expect(err).ShouldNot(m.HaveOccurred())

				m.Expect(doc1.Revision).To(m.Equal(uint64(2)))
				m.Expect(doc1.CreatedAt).To(
					m.BeTemporally("==", createdAt),
				)
				m.Expect(doc1.UpdatedAt).To(
					m.BeTemporally("~", time.Now(), 50*time.Millisecond),
				)
			})

			g.It("does not return an error, even if the provided revision is not 1", func() {
				doc1.Revision = 123
				op := protavo.ForceSave(doc1)

				err := db.Write(ctx, op)
				m.Expect(err).ShouldNot(m.HaveOccurred())

				doc, ok, err := db.Load(ctx, "doc-1")
				m.Expect(err).ShouldNot(m.HaveOccurred())
				m.Expect(ok).To(m.BeTrue())
				m.Expect(doc.Revision).To(m.Equal(uint64(2)))
			})
		})

		g.It("aborts the save if a unique key conflicts with another document", func() {
			op := protavo.ForceSave(doc1)

			err := db.Write(ctx, op)
			m.Expect(err).ShouldNot(m.HaveOccurred())

			doc2.Keys = document.UniqueKeys("<unique-key>") // doc1 already has this key
			op = protavo.ForceSave(doc2)

			err = db.Write(ctx, op)
			m.Expect(err).To(m.Equal(
				&protavo.DuplicateKeyError{
					DocumentID:            "doc-2",
					ConflictingDocumentID: "doc-1",
					UniqueKey:             "<unique-key>",
				},
			))
		})
	})
}
