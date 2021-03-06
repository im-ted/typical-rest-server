package mylibrary_test

import (
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-rest-server/internal/app/domain/mylibrary"
	"github.com/typical-go/typical-rest-server/pkg/echokit"
)

func TestRoute(t *testing.T) {
	e := echo.New()
	echokit.SetRoute(e, &mylibrary.Router{})
	require.Equal(t, []string{
		"/mylibrary/books\tGET,POST",
		"/mylibrary/books/:id\tDELETE,GET,HEAD,PATCH,PUT",
	}, echokit.DumpEcho(e))
}
