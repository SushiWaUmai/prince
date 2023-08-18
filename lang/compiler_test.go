package lang

import (
	"testing"

	"github.com/SushiWaUmai/prince/utils"
	"github.com/stretchr/testify/assert"
)

func TestCompiler(t *testing.T) {
	t.Run("Test Compiler with simple command", func(t *testing.T) {
		assert := assert.New(t)
		sample := "!ping hello | ping \"hello world\""
		tokens := Lex(sample)
		expressions, err := Parse(tokens)
		assert.Nil(err)
		commandInputs, err := Compile(expressions)
		assert.Nil(err)

		assert.Equal([]utils.CommandInput{
			{
				Name: "ping",
				Args: []string{"hello"},
			},
			{
				Name: "ping",
				Args: []string{"hello world"},
			},
		}, commandInputs)
	})
}
