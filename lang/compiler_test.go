package lang_test

import (
	"testing"

	"github.com/SushiWaUmai/prince/env"
	"github.com/SushiWaUmai/prince/lang"
	"github.com/SushiWaUmai/prince/utils"
	"github.com/stretchr/testify/assert"
)

func TestCompiler(t *testing.T) {
	t.Run("Test Compiler with simple command", func(t *testing.T) {
		assert := assert.New(t)
		sample := string(env.BOT_PREFIX) + "echo hello | echo \"hello world\""
		tokens := lang.Lex(sample)
		expressions, err := lang.Parse(tokens)
		assert.Nil(err)
		commandInputs, err := lang.Compile(expressions)
		assert.Nil(err)

		assert.Equal([]utils.CommandInput{
			{
				Name: "echo",
				Args: []string{"hello"},
			},
			{
				Name: "echo",
				Args: []string{"hello world"},
			},
		}, commandInputs)
	})
}
