package protobolt_test

import (
	"context"
	"io/ioutil"
	"os"
	"time"

	. "github.com/jmalloc/protobolt/src/protobolt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/jmalloc/protobolt/src/protobolt/internal/testtypes"
)

func init() {
	describeDBTest("DB (shared)", true)
	describeDBTest("DB (exclusive)", false)
}

func describeDBTest(
	n string,
	shared bool,
) {
	Describe(n, func() {
		ctx := context.Background()
		var (
			db   *DB
			path string
		)

		BeforeEach(func() {
			fp, err := ioutil.TempFile("", "protobolt-")
			Expect(err).ShouldNot(HaveOccurred())

			err = fp.Close()
			Expect(err).ShouldNot(HaveOccurred())

			path = fp.Name()
			err = os.Remove(path)
			Expect(err).ShouldNot(HaveOccurred())

			db, err = Open(path, shared, 0600, nil)
			Expect(err).ShouldNot(HaveOccurred())
		})

		AfterEach(func() {
			_ = db.Close()
			_ = os.Remove(path)
		})

		Describe("Load", func() {
			When("the document exists", func() {
				var doc *Document

				BeforeEach(func() {
					doc = NewDocument(
						"doc-id",
						&testtypes.TestContent{
							Data: "<content-1>",
						},
					)

					var err error
					doc, err = db.Save(ctx, doc)
					Expect(err).ShouldNot(HaveOccurred())
				})

				It("returns true", func() {
					_, ok, err := db.Load(ctx, "doc-id")
					Expect(err).ShouldNot(HaveOccurred())
					Expect(ok).To(BeTrue())
				})

				It("returns the document", func() {
					d, _, err := db.Load(ctx, "doc-id")
					Expect(err).ShouldNot(HaveOccurred())
					Expect(d.Equal(doc)).To(BeTrue())
				})
			})

			When("the document does not exist", func() {
				It("returns false", func() {
					_, ok, err := db.Load(ctx, "doc-id")
					Expect(err).ShouldNot(HaveOccurred())
					Expect(ok).To(BeFalse())
				})
			})
		})

		Describe("Save", func() {
			When("a new document is created", func() {
				var (
					savedDoc    *Document
					returnedDoc *Document
				)

				BeforeEach(func() {
					savedDoc = NewDocument("doc-id", &testtypes.TestContent{
						Data: "<content-1>",
					})

					var err error
					returnedDoc, err = db.Save(ctx, savedDoc)
					Expect(err).ShouldNot(HaveOccurred())
				})

				It("returns a copy of the document at version 1", func() {
					Expect(returnedDoc.Version()).To(
						BeNumerically("==", 1),
					)
				})

				It("sets the created and updated timestamps", func() {
					ct, ok := returnedDoc.CreatedAt()
					Expect(ok).To(BeTrue())
					Expect(ct).To(
						BeTemporally("~", time.Now(), 50*time.Millisecond),
					)

					ut, ok := returnedDoc.UpdatedAt()
					Expect(ok).To(BeTrue())
					Expect(ut.Equal(ct)).To(BeTrue())
				})

				It("persists the document as returned", func() {
					d, _, err := db.Load(ctx, "doc-id")
					Expect(err).ShouldNot(HaveOccurred())
					Expect(d.Equal(returnedDoc)).To(BeTrue())
				})

				It("returns an optimistic lock error if the version is non-zero", func() {
					// first delete the document
					err := db.Delete(ctx, returnedDoc)
					Expect(err).ShouldNot(HaveOccurred())

					// then try to save it with the version of the document @ v1
					_, err = db.Save(ctx, returnedDoc)
					Expect(err).To(
						MatchError(&OptimisticLockError{
							DocumentID:    "doc-id",
							GivenVersion:  1,
							ActualVersion: 0,
							Action:        "save",
						}),
					)
				})
			})
		})

		When("an existing document is updated", func() {
			var (
				savedDoc    *Document
				returnedDoc *Document
			)

			BeforeEach(func() {
				savedDoc = NewDocument(
					"doc-id",
					&testtypes.TestContent{
						Data: "<content-1>",
					},
				)

				var err error
				savedDoc, err = db.Save(ctx, savedDoc)
				Expect(err).ShouldNot(HaveOccurred())

				savedDoc = savedDoc.WithContent(
					&testtypes.TestContent{
						Data: "<content-2>",
					},
				)

				returnedDoc, err = db.Save(ctx, savedDoc)
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("returns a copy of the document at version 2", func() {
				Expect(returnedDoc.Version()).To(
					BeNumerically("==", 2),
				)
			})

			It("does not modify the created timestamp", func() {
				after, _ := returnedDoc.CreatedAt()
				before, _ := savedDoc.CreatedAt()

				Expect(after).To(Equal(before))
			})

			It("sets the updated timestamp", func() {
				t, ok := returnedDoc.UpdatedAt()
				Expect(ok).To(BeTrue())
				Expect(t).To(
					BeTemporally("~", time.Now(), 50*time.Millisecond),
				)
			})

			It("persists the document as returned", func() {
				d, _, err := db.Load(ctx, "doc-id")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(d.Equal(returnedDoc)).To(BeTrue())
			})

			It("returns an optimistic lock error if the version is not equal", func() {
				_, err := db.Save(ctx, savedDoc)
				Expect(err).To(MatchError(&OptimisticLockError{
					DocumentID:    "doc-id",
					GivenVersion:  1,
					ActualVersion: 2,
					Action:        "save",
				}))
			})
		})

		When("the document has a unique key", func() {
			var (
				savedDoc *Document
			)

			BeforeEach(func() {
				savedDoc = NewDocument(
					"doc-1",
					&testtypes.TestContent{
						Data: "<content-1>",
					},
				).WithUniqueKeys("<uniq>")

				var err error
				savedDoc, err = db.Save(ctx, savedDoc)
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("can be updated while retaining the unique key", func() {
				doc := savedDoc.WithContent(
					&testtypes.TestContent{
						Data: "<content-3>",
					},
				)

				_, err := db.Save(ctx, doc)
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("returns a duplicate key error if another document is saved with the same unique key", func() {
				doc := NewDocument(
					"doc-2",
					&testtypes.TestContent{
						Data: "<content-1>",
					},
				).WithUniqueKeys("<uniq>")

				_, err := db.Save(ctx, doc)
				Expect(err).To(MatchError(&DuplicateKeyError{
					DocumentID:            "doc-2",
					ConflictingDocumentID: "doc-1",
					UniqueKey:             "<uniq>",
				}))
			})
		})

		Describe("Delete", func() {
			BeforeEach(func() {
				savedDoc := NewDocument(
					"doc-id",
					&testtypes.TestContent{
						Data: "<content-1>",
					},
				).WithUniqueKeys("<uniq>")

				var err error
				savedDoc, err = db.Save(ctx, savedDoc)
				Expect(err).ShouldNot(HaveOccurred())

				err = db.Delete(ctx, savedDoc)
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("removes the document", func() {
				_, ok, err := db.Load(ctx, "doc-id")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(ok).To(BeFalse())
			})

			It("allows other documents with the same unique key to be saved", func() {
				doc := NewDocument(
					"doc-id",
					&testtypes.TestContent{
						Data: "<content-1>",
					},
				).WithUniqueKeys("<uniq>")

				_, err := db.Save(ctx, doc)
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
	})
}
