package sqkit_test

import (
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-rest-server/pkg/sqkit"
)

func TestSelectOption(t *testing.T) {
	expected := sq.Select("")
	selectOpt := sqkit.NewSelectOption(func(sq.SelectBuilder) sq.SelectBuilder {
		return expected
	})
	require.Equal(t, expected, selectOpt.CompileSelect(sq.Select("")))
}
