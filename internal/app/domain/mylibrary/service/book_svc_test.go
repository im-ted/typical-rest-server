package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/typical-go/typical-rest-server/internal/app/data_access/postgresdb"
	"github.com/typical-go/typical-rest-server/internal/generated/postgresdb_repo"
	"github.com/typical-go/typical-rest-server/internal/generated/postgresdb_repo_mock"

	"github.com/typical-go/typical-rest-server/internal/app/domain/mylibrary/service"
	"github.com/typical-go/typical-rest-server/pkg/dbkit"
)

type bookSvcFn func(mockRepo *postgresdb_repo_mock.MockBookRepo)

func createBookSvc(t *testing.T, fn bookSvcFn) (*service.BookSvcImpl, *gomock.Controller) {
	mock := gomock.NewController(t)
	mockRepo := postgresdb_repo_mock.NewMockBookRepo(mock)
	if fn != nil {
		fn(mockRepo)
	}

	return &service.BookSvcImpl{
		Repo: mockRepo,
	}, mock
}

func TestBookSvc_Create(t *testing.T) {
	testcases := []struct {
		testName    string
		bookSvcFn   bookSvcFn
		book        *postgresdb.Book
		expected    *postgresdb.Book
		expectedErr string
	}{
		{
			testName:    "validation error",
			book:        &postgresdb.Book{},
			expectedErr: "Key: 'Book.Title' Error:Field validation for 'Title' failed on the 'required' tag\nKey: 'Book.Author' Error:Field validation for 'Author' failed on the 'required' tag",
		},
		{
			testName:    "create error",
			book:        &postgresdb.Book{Author: "some-author", Title: "some-title"},
			expectedErr: "create-error",
			bookSvcFn: func(mockRepo *postgresdb_repo_mock.MockBookRepo) {
				mockRepo.EXPECT().
					Create(gomock.Any(), &postgresdb.Book{Author: "some-author", Title: "some-title"}).
					Return(int64(-1), errors.New("create-error"))
			},
		},
		{
			testName:    "Find error",
			book:        &postgresdb.Book{Author: "some-author", Title: "some-title"},
			expectedErr: "Find-error",
			bookSvcFn: func(mockRepo *postgresdb_repo_mock.MockBookRepo) {
				mockRepo.EXPECT().
					Create(gomock.Any(), &postgresdb.Book{Author: "some-author", Title: "some-title"}).
					Return(int64(1), nil)
				mockRepo.EXPECT().
					Find(gomock.Any(), dbkit.Eq{"id": int64(1)}).
					Return(nil, errors.New("Find-error"))
			},
		},
		{
			book: &postgresdb.Book{
				Author: "some-author",
				Title:  "some-title",
			},
			expected: &postgresdb.Book{Author: "some-author", Title: "some-title"},
			bookSvcFn: func(mockRepo *postgresdb_repo_mock.MockBookRepo) {
				mockRepo.EXPECT().
					Create(gomock.Any(), &postgresdb.Book{Author: "some-author", Title: "some-title"}).
					Return(int64(1), nil)
				mockRepo.EXPECT().
					Find(gomock.Any(), dbkit.Eq{"id": int64(1)}).
					Return([]*postgresdb.Book{{Author: "some-author", Title: "some-title"}}, nil)
			},
		},
	}

	for _, tt := range testcases {
		t.Run(tt.testName, func(t *testing.T) {
			svc, mock := createBookSvc(t, tt.bookSvcFn)
			defer mock.Finish()

			id, err := svc.Create(context.Background(), tt.book)
			if tt.expectedErr != "" {
				require.EqualError(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, id)
			}
		})
	}
}

func TestBookSvc_FindOne(t *testing.T) {
	testcases := []struct {
		testName    string
		bookSvcFn   bookSvcFn
		paramID     string
		expected    *postgresdb.Book
		expectedErr string
	}{
		{
			paramID:     "",
			expectedErr: "paramID is missing",
		},
		{
			paramID: "1",
			bookSvcFn: func(mockRepo *postgresdb_repo_mock.MockBookRepo) {
				mockRepo.EXPECT().
					Find(gomock.Any(), dbkit.Eq{"id": int64(1)}).
					Return(nil, errors.New("some-error"))
			},
			expectedErr: "some-error",
		},
		{
			paramID: "1",
			bookSvcFn: func(mockRepo *postgresdb_repo_mock.MockBookRepo) {
				mockRepo.EXPECT().
					Find(gomock.Any(), dbkit.Eq{"id": int64(1)}).
					Return([]*postgresdb.Book{
						{ID: 1, Title: "some-title"},
					}, nil)
			},
			expected: &postgresdb.Book{ID: 1, Title: "some-title"},
		},
	}

	for _, tt := range testcases {
		t.Run(tt.testName, func(t *testing.T) {
			svc, mock := createBookSvc(t, tt.bookSvcFn)
			defer mock.Finish()

			book, err := svc.FindOne(context.Background(), tt.paramID)
			if tt.expectedErr != "" {
				require.EqualError(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, book)
			}
		})
	}
}

func TestBookSvc_Find(t *testing.T) {
	testcases := []struct {
		testName    string
		bookSvcFn   bookSvcFn
		req         *service.FindReq
		expected    []*postgresdb.Book
		expectedErr string
	}{
		{
			bookSvcFn: func(mockRepo *postgresdb_repo_mock.MockBookRepo) {
				mockRepo.EXPECT().
					Find(gomock.Any(), &dbkit.OffsetPagination{}).
					Return([]*postgresdb.Book{
						{ID: 1, Title: "title1", Author: "author1"},
						{ID: 2, Title: "title2", Author: "author2"},
					}, nil)
			},
			req: &service.FindReq{},
			expected: []*postgresdb.Book{
				{ID: 1, Title: "title1", Author: "author1"},
				{ID: 2, Title: "title2", Author: "author2"},
			},
		},
		{
			bookSvcFn: func(mockRepo *postgresdb_repo_mock.MockBookRepo) {
				mockRepo.EXPECT().
					Find(gomock.Any(), &dbkit.OffsetPagination{Limit: 20, Offset: 10}, dbkit.Sorts{"title", "created_at"}).
					Return(nil, errors.New("some-error"))
			},
			req:         &service.FindReq{Limit: 20, Offset: 10, Sort: "title,created_at"},
			expectedErr: "some-error",
		},
	}
	for _, tt := range testcases {
		t.Run(tt.testName, func(t *testing.T) {
			svc, mock := createBookSvc(t, tt.bookSvcFn)
			defer mock.Finish()

			books, err := svc.Find(context.Background(), tt.req)
			if tt.expectedErr != "" {
				require.EqualError(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, books)
			}
		})
	}
}

func TestBookSvc_Delete(t *testing.T) {
	testcases := []struct {
		testName    string
		bookSvcFn   bookSvcFn
		paramID     string
		expectedErr string
	}{
		{
			paramID:     "",
			expectedErr: `paramID is missing`,
		},
		{
			paramID:     "1",
			expectedErr: `some-error`,
			bookSvcFn: func(mockRepo *postgresdb_repo_mock.MockBookRepo) {
				mockRepo.EXPECT().
					Delete(gomock.Any(), dbkit.Eq{postgresdb_repo.BookTable.ID: int64(1)}).
					Return(int64(0), errors.New("some-error"))
			},
		},
		{
			paramID: "1",
			bookSvcFn: func(mockRepo *postgresdb_repo_mock.MockBookRepo) {
				mockRepo.EXPECT().
					Delete(gomock.Any(), dbkit.Eq{postgresdb_repo.BookTable.ID: int64(1)}).
					Return(int64(1), nil)
			},
		},
		{
			testName: "success even if no affected row (idempotent)",
			paramID:  "1",
			bookSvcFn: func(mockRepo *postgresdb_repo_mock.MockBookRepo) {
				mockRepo.EXPECT().
					Delete(gomock.Any(), dbkit.Eq{postgresdb_repo.BookTable.ID: int64(1)}).
					Return(int64(0), nil)
			},
		},
	}
	for _, tt := range testcases {
		t.Run(tt.testName, func(t *testing.T) {
			svc, mock := createBookSvc(t, tt.bookSvcFn)
			defer mock.Finish()

			err := svc.Delete(context.Background(), tt.paramID)
			if tt.expectedErr != "" {
				require.EqualError(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestBookSvc_Update(t *testing.T) {
	testcases := []struct {
		testName    string
		bookSvcFn   bookSvcFn
		paramID     string
		book        *postgresdb.Book
		expected    *postgresdb.Book
		expectedErr string
	}{
		{
			testName:    "empty paramID",
			paramID:     "",
			expectedErr: `paramID is missing`,
		},
		{
			testName:    "zero paramID",
			paramID:     "0",
			expectedErr: `paramID is missing`,
		},
		{
			testName:    "bad request",
			paramID:     "1",
			book:        &postgresdb.Book{},
			expectedErr: "Key: 'Book.Title' Error:Field validation for 'Title' failed on the 'required' tag\nKey: 'Book.Author' Error:Field validation for 'Author' failed on the 'required' tag",
		},
		{
			testName:    "update error",
			paramID:     "1",
			book:        &postgresdb.Book{Author: "some-author", Title: "some-title"},
			expectedErr: "update error",
			bookSvcFn: func(mockRepo *postgresdb_repo_mock.MockBookRepo) {
				mockRepo.EXPECT().
					Update(gomock.Any(), &postgresdb.Book{Author: "some-author", Title: "some-title"}, dbkit.Eq{"id": int64(1)}).
					Return(int64(-1), errors.New("update error"))
			},
		},
		{
			testName:    "nothing to update",
			paramID:     "1",
			book:        &postgresdb.Book{Author: "some-author", Title: "some-title"},
			expectedErr: "sql: no rows in result set",
			bookSvcFn: func(mockRepo *postgresdb_repo_mock.MockBookRepo) {
				mockRepo.EXPECT().
					Update(gomock.Any(), &postgresdb.Book{Author: "some-author", Title: "some-title"}, dbkit.Eq{"id": int64(1)}).
					Return(int64(0), nil)
			},
		},
		{
			testName:    "Find error",
			paramID:     "1",
			book:        &postgresdb.Book{Author: "some-author", Title: "some-title"},
			expectedErr: "Find-error",
			bookSvcFn: func(mockRepo *postgresdb_repo_mock.MockBookRepo) {
				mockRepo.EXPECT().
					Update(gomock.Any(), &postgresdb.Book{Author: "some-author", Title: "some-title"}, dbkit.Eq{"id": int64(1)}).
					Return(int64(1), nil)
				mockRepo.EXPECT().
					Find(gomock.Any(), dbkit.Eq{"id": int64(1)}).
					Return(nil, errors.New("Find-error"))
			},
		},
		{
			testName:    "Find error",
			paramID:     "1",
			book:        &postgresdb.Book{Author: "some-author", Title: "some-title"},
			expectedErr: "Find-error",
			bookSvcFn: func(mockRepo *postgresdb_repo_mock.MockBookRepo) {
				mockRepo.EXPECT().
					Update(gomock.Any(), &postgresdb.Book{Author: "some-author", Title: "some-title"}, dbkit.Eq{"id": int64(1)}).
					Return(int64(1), nil)
				mockRepo.EXPECT().
					Find(gomock.Any(), dbkit.Eq{"id": int64(1)}).
					Return(nil, errors.New("Find-error"))
			},
		},
	}
	for _, tt := range testcases {
		t.Run(tt.testName, func(t *testing.T) {
			svc, mock := createBookSvc(t, tt.bookSvcFn)
			defer mock.Finish()

			book, err := svc.Update(context.Background(), tt.paramID, tt.book)
			if tt.expectedErr != "" {
				require.EqualError(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, book)
			}
		})
	}
}

func TestBookSvc_Patch(t *testing.T) {
	testcases := []struct {
		testName    string
		bookSvcFn   bookSvcFn
		paramID     string
		book        *postgresdb.Book
		expected    *postgresdb.Book
		expectedErr string
	}{
		{
			testName:    "empty paramID",
			paramID:     "",
			expectedErr: "paramID is missing",
		},
		{
			testName:    "zero paramID",
			paramID:     "0",
			expectedErr: "paramID is missing",
		},
		{
			testName:    "patch error",
			paramID:     "1",
			book:        &postgresdb.Book{Author: "some-author", Title: "some-title"},
			expectedErr: "patch-error",
			bookSvcFn: func(mockRepo *postgresdb_repo_mock.MockBookRepo) {
				mockRepo.EXPECT().
					Patch(gomock.Any(), &postgresdb.Book{Author: "some-author", Title: "some-title"}, dbkit.Eq{"id": int64(1)}).
					Return(int64(-1), errors.New("patch-error"))
			},
		},
		{
			testName:    "patch error",
			paramID:     "1",
			book:        &postgresdb.Book{Author: "some-author", Title: "some-title"},
			expectedErr: "sql: no rows in result set",
			bookSvcFn: func(mockRepo *postgresdb_repo_mock.MockBookRepo) {
				mockRepo.EXPECT().
					Patch(gomock.Any(), &postgresdb.Book{Author: "some-author", Title: "some-title"}, dbkit.Eq{"id": int64(1)}).
					Return(int64(0), nil)
			},
		},
		{
			testName:    "Find error",
			paramID:     "1",
			book:        &postgresdb.Book{Author: "some-author", Title: "some-title"},
			expectedErr: "Find-error",
			bookSvcFn: func(mockRepo *postgresdb_repo_mock.MockBookRepo) {
				mockRepo.EXPECT().
					Patch(gomock.Any(), &postgresdb.Book{Author: "some-author", Title: "some-title"}, dbkit.Eq{"id": int64(1)}).
					Return(int64(1), nil)
				mockRepo.EXPECT().
					Find(gomock.Any(), dbkit.Eq{"id": int64(1)}).
					Return(nil, errors.New("Find-error"))
			},
		},
		{
			paramID:  "1",
			book:     &postgresdb.Book{Author: "some-author", Title: "some-title"},
			expected: &postgresdb.Book{Author: "some-author", Title: "some-title"},
			bookSvcFn: func(mockRepo *postgresdb_repo_mock.MockBookRepo) {
				mockRepo.EXPECT().
					Patch(gomock.Any(), &postgresdb.Book{Author: "some-author", Title: "some-title"}, dbkit.Eq{"id": int64(1)}).
					Return(int64(1), nil)
				mockRepo.EXPECT().
					Find(gomock.Any(), dbkit.Eq{"id": int64(1)}).
					Return([]*postgresdb.Book{{Author: "some-author", Title: "some-title"}}, nil)
			},
		},
	}
	for _, tt := range testcases {
		t.Run(tt.testName, func(t *testing.T) {
			svc, mock := createBookSvc(t, tt.bookSvcFn)
			defer mock.Finish()

			book, err := svc.Patch(context.Background(), tt.paramID, tt.book)
			if tt.expectedErr != "" {
				require.EqualError(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, book)
			}
		})
	}
}
