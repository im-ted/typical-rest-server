package repository_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-rest-server/internal/server/repository"
	"github.com/typical-go/typical-rest-server/pkg/dbkit"
)

type bookRepoFn func(sqlmock.Sqlmock)

func createBookRepo(fn bookRepoFn) (repository.BookRepo, *sql.DB) {
	db, mock, _ := sqlmock.New()
	if fn != nil {
		fn(mock)
	}
	return &repository.BookRepoImpl{DB: db}, db
}

func TestBookRepoImpl_Create(t *testing.T) {
	testcases := []struct {
		testName           string
		book               *repository.Book
		bookRepoFn         bookRepoFn
		expectedInsertedID int64
		expectedErr        string
	}{
		{
			book: &repository.Book{
				Title:  "some-title",
				Author: "some-author",
			},
			expectedErr: "some-error",
			bookRepoFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO books (title,author,created_at,updated_at) VALUES ($1,$2,$3,$4) RETURNING "id"`)).
					WithArgs("some-title", "some-author", sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(errors.New("some-error"))
			},
		},
		{
			book: &repository.Book{
				Title:  "some-title",
				Author: "some-author",
			},
			bookRepoFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO books (title,author,created_at,updated_at) VALUES ($1,$2,$3,$4) RETURNING "id"`)).
					WithArgs("some-title", "some-author", sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(999))
			},
			expectedInsertedID: 999,
		},
	}

	for _, tt := range testcases {
		t.Run(tt.testName, func(t *testing.T) {
			repo, db := createBookRepo(tt.bookRepoFn)
			defer db.Close()

			id, err := repo.Create(context.Background(), tt.book)
			if tt.expectedErr != "" {
				require.EqualError(t, err, tt.expectedErr)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.expectedInsertedID, id)
		})
	}
}

func TestBookRepoImpl_Update(t *testing.T) {
	testcases := []struct {
		testName            string
		book                *repository.Book
		bookRepoFn          bookRepoFn
		opt                 dbkit.UpdateOption
		expectedErr         string
		expectedAffectedRow int64
	}{
		{
			book: &repository.Book{
				Title:  "new-title",
				Author: "new-author",
			},
			opt: dbkit.Equal(repository.BookTable.ID, 888),
			bookRepoFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE books SET title = $1, author = $2, updated_at = $3 WHERE id = $4`)).
					WithArgs("new-title", "new-author", sqlmock.AnyArg(), 888).
					WillReturnError(fmt.Errorf("some-update-error"))
			},
			expectedErr:         "some-update-error",
			expectedAffectedRow: -1,
		},
		{
			book: &repository.Book{
				Title:  "new-title",
				Author: "new-author",
			},
			opt: dbkit.Equal(repository.BookTable.ID, 888),
			bookRepoFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE books SET title = $1, author = $2, updated_at = $3 WHERE id = $4`)).
					WithArgs("new-title", "new-author", sqlmock.AnyArg(), 888).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectedAffectedRow: 1,
		},
		{
			testName: "empty author",
			book: &repository.Book{
				Title: "new-title",
			},
			opt: dbkit.Equal(repository.BookTable.ID, 888),
			bookRepoFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE books SET title = $1, author = $2, updated_at = $3 WHERE id = $4`)).
					WithArgs("new-title", "", sqlmock.AnyArg(), 888).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectedAffectedRow: 1,
		},
		{
			testName: "empty title",
			book: &repository.Book{
				Author: "new-author",
			},
			opt: dbkit.Equal(repository.BookTable.ID, 888),
			bookRepoFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE books SET title = $1, author = $2, updated_at = $3 WHERE id = $4`)).
					WithArgs("", "new-author", sqlmock.AnyArg(), 888).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectedAffectedRow: 1,
		},
	}
	for _, tt := range testcases {
		t.Run(tt.testName, func(t *testing.T) {
			repo, db := createBookRepo(tt.bookRepoFn)
			defer db.Close()

			affectedRow, err := repo.Update(context.Background(), tt.book, tt.opt)
			if tt.expectedErr != "" {
				require.EqualError(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedAffectedRow, affectedRow)
			}
		})
	}
}

func TestBookRepoImpl_Patch(t *testing.T) {
	testcases := []struct {
		testName            string
		book                *repository.Book
		bookRepoFn          bookRepoFn
		opt                 dbkit.UpdateOption
		expectedErr         string
		expectedAffectedRow int64
	}{
		{
			book: &repository.Book{
				Title:  "new-title",
				Author: "new-author",
			},
			opt: dbkit.Equal(repository.BookTable.ID, 888),
			bookRepoFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE books SET title = $1, author = $2, updated_at = $3 WHERE id = $4`)).
					WithArgs("new-title", "new-author", sqlmock.AnyArg(), 888).
					WillReturnError(fmt.Errorf("some-update-error"))
			},
			expectedErr:         "some-update-error",
			expectedAffectedRow: -1,
		},
		{
			book: &repository.Book{
				Title:  "new-title",
				Author: "new-author",
			},
			opt: dbkit.Equal(repository.BookTable.ID, 888),
			bookRepoFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE books SET title = $1, author = $2, updated_at = $3 WHERE id = $4`)).
					WithArgs("new-title", "new-author", sqlmock.AnyArg(), 888).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectedAffectedRow: 1,
		},
		{
			testName: "empty author",
			book: &repository.Book{
				Title: "new-title",
			},
			opt: dbkit.Equal(repository.BookTable.ID, 888),
			bookRepoFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE books SET title = $1, updated_at = $2 WHERE id = $3`)).
					WithArgs("new-title", sqlmock.AnyArg(), 888).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectedAffectedRow: 1,
		},
		{
			testName: "empty title",
			book: &repository.Book{
				Author: "new-author",
			},
			opt: dbkit.Equal(repository.BookTable.ID, 888),
			bookRepoFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE books SET author = $1, updated_at = $2 WHERE id = $3`)).
					WithArgs("new-author", sqlmock.AnyArg(), 888).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectedAffectedRow: 1,
		},
	}
	for _, tt := range testcases {
		t.Run(tt.testName, func(t *testing.T) {
			repo, db := createBookRepo(tt.bookRepoFn)
			defer db.Close()

			affectedRow, err := repo.Patch(context.Background(), tt.book, tt.opt)
			if tt.expectedErr != "" {
				require.EqualError(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedAffectedRow, affectedRow)
			}
		})
	}
}

func TestBookRepoImpl_Retrieve(t *testing.T) {
	now := time.Now()
	testcases := []struct {
		testName    string
		opts        []dbkit.SelectOption
		expected    []*repository.Book
		expectedErr string
		bookRepoFn  bookRepoFn
	}{
		{

			opts:        []dbkit.SelectOption{},
			expected:    []*repository.Book{},
			expectedErr: "some-error",
			bookRepoFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id, title, author, updated_at, created_at FROM books`).
					WillReturnError(errors.New("some-error"))
			},
		},
		{

			opts: []dbkit.SelectOption{},
			expected: []*repository.Book{
				&repository.Book{ID: 1234, Title: "some-title4", Author: "some-author4", UpdatedAt: now, CreatedAt: now},
				&repository.Book{ID: 1235, Title: "some-title5", Author: "some-author5", UpdatedAt: now, CreatedAt: now},
			},
			bookRepoFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id, title, author, updated_at, created_at FROM books`).
					WillReturnRows(sqlmock.
						NewRows([]string{"id", "title", "author", "updated_at", "created_at"}).
						AddRow("1234", "some-title4", "some-author4", now, now).
						AddRow("1235", "some-title5", "some-author5", now, now),
					)
			},
		},
	}
	for _, tt := range testcases {
		t.Run(tt.testName, func(t *testing.T) {
			repo, db := createBookRepo(tt.bookRepoFn)
			defer db.Close()

			books, err := repo.Retrieve(context.Background(), tt.opts...)
			if tt.expectedErr != "" {
				require.EqualError(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, books)
			}
		})
	}
}

func TestBookRepoImpl_Delete(t *testing.T) {
	testcases := []struct {
		testName            string
		opt                 dbkit.DeleteOption
		bookRepoFn          bookRepoFn
		expectedErr         string
		expectedAffectedRow int64
	}{
		{
			opt:         dbkit.Equal("id", 666),
			expectedErr: "some-error",
			bookRepoFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM books WHERE id = $1`)).
					WithArgs(666).
					WillReturnError(fmt.Errorf("some-error"))
			},
		},
		{
			opt: dbkit.Equal("id", 555),
			bookRepoFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM books WHERE id = $1`)).
					WithArgs(555).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectedAffectedRow: 1,
		},
	}
	for _, tt := range testcases {
		t.Run(tt.testName, func(t *testing.T) {
			repo, db := createBookRepo(tt.bookRepoFn)
			defer db.Close()

			affectedRow, err := repo.Delete(context.Background(), tt.opt)
			if tt.expectedErr != "" {
				require.EqualError(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedAffectedRow, affectedRow)
			}
		})
	}
}
