package readme

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadmeRecipe(t *testing.T) {

	recipe := &Readme{
		Title:       "some-title",
		Description: "some-descrption",
	}
	recipe.SetSection("section-title1", func(md *Markdown) error {
		md.Writeln("some-content")
		return nil
	})
	// recipe.SetSection("section-title1", Section{Content: "some-content"})
	recipe.SetSection("section-title2", func(md *Markdown) error {
		md.Writeln("some-field1")
		return nil
	})
	recipe.SetSection("section-title1", func(md *Markdown) error {
		md.Writeln("revision-content")
		return nil
	})

	var builder strings.Builder
	recipe.Output(&builder)

	require.Equal(t, `<!-- Autogenerated by Typical-Go. DO NOT EDIT. -->

# some-title

some-descrption

## section-title1

revision-content

## section-title2

some-field1

`, builder.String())
}
