package lang_test

import (
	"testing"

	_ "github.com/SushiWaUmai/prince/commands"
	"github.com/SushiWaUmai/prince/db"
	"github.com/SushiWaUmai/prince/env"
	"github.com/SushiWaUmai/prince/lang"
	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	t.Run("Test Parser with simple command", func(t *testing.T) {
		assert := assert.New(t)
		sample := string(env.BOT_PREFIX) + "echo hello | echo \"hello world\""

		tokens := lang.Lex(sample)

		expressions, err := lang.Parse(tokens)
		assert.Nil(err)

		assert.Equal([]lang.Expression{
			{
				Type:    lang.COMMAND,
				Content: "echo",
			},
			{
				Type:    lang.ARGUMENT,
				Content: "hello",
			},
			{
				Type:    lang.COMMAND,
				Content: "echo",
			},
			{
				Type:    lang.ARGUMENT,
				Content: "hello world",
			},
		}, expressions)
	})

	t.Run("Test Parser with invalid prefix", func(t *testing.T) {
		assert := assert.New(t)
		sample := "invalid-prefixecho asdf | echo test"

		tokens := lang.Lex(sample)

		expressions, err := lang.Parse(tokens)
		assert.Nil(expressions)
		assert.NotNil(err)
	})

	t.Run("Test Parser with invalid pipe", func(t *testing.T) {
		assert := assert.New(t)
		sample := string(env.BOT_PREFIX) + "echo asdf | echo hello |"

		tokens := lang.Lex(sample)

		expressions, err := lang.Parse(tokens)
		assert.Nil(expressions)
		assert.NotNil(err)
	})

	t.Run("Test Parser with no argument", func(t *testing.T) {
		assert := assert.New(t)
		sample := string(env.BOT_PREFIX) + "download"

		tokens := lang.Lex(sample)
		expressions, err := lang.Parse(tokens)
		assert.Nil(err)

		assert.Equal(
			[]lang.Expression{
				{
					Type:    lang.COMMAND,
					Content: "download",
				},
			},
			expressions,
		)
	})

	t.Run("Test Parser with alias", func(t *testing.T) {
		assert := assert.New(t)

		db.CreateAlias("hello", "echo hello")
		defer db.DeleteAlias("hello")
		sample := string(env.BOT_PREFIX) + "hello"

		tokens := lang.Lex(sample)
		expressions, err := lang.Parse(tokens)
		assert.Nil(err)

		assert.Equal(
			[]lang.Expression{
				{
					Type:    lang.COMMAND,
					Content: "echo",
				},
				{
					Type:    lang.ARGUMENT,
					Content: "hello",
				},
			},
			expressions,
		)
	})

	t.Run("Test Parser with alias with pipes", func(t *testing.T) {
		assert := assert.New(t)

		db.CreateAlias("hello", "echo hello | chat \"how are you\"")
		defer db.DeleteAlias("hello")
		sample := string(env.BOT_PREFIX) + "hello | echo \"hello world\""

		tokens := lang.Lex(sample)
		expressions, err := lang.Parse(tokens)
		assert.Nil(err)

		assert.Equal(
			[]lang.Expression{
				{
					Type:    lang.COMMAND,
					Content: "echo",
				},
				{
					Type:    lang.ARGUMENT,
					Content: "hello",
				},
				{
					Type:    lang.COMMAND,
					Content: "chat",
				},
				{
					Type:    lang.ARGUMENT,
					Content: "how are you",
				},
				{
					Type:    lang.COMMAND,
					Content: "echo",
				},
				{
					Type:    lang.ARGUMENT,
					Content: "hello world",
				},
			},
			expressions,
		)
	})

	t.Run("Test Parser with alias at argument position", func(t *testing.T) {
		assert := assert.New(t)

		db.CreateAlias("hello", "echo hello | chat \"how are you\"")
		defer db.DeleteAlias("hello")
		sample := string(env.BOT_PREFIX) + "echo hello | chat hello world"

		tokens := lang.Lex(sample)
		expressions, err := lang.Parse(tokens)
		assert.Nil(err)

		assert.Equal(
			[]lang.Expression{
				{
					Type:    lang.COMMAND,
					Content: "echo",
				},
				{
					Type:    lang.ARGUMENT,
					Content: "hello",
				},
				{
					Type:    lang.COMMAND,
					Content: "chat",
				},
				{
					Type:    lang.ARGUMENT,
					Content: "hello",
				},
				{
					Type:    lang.ARGUMENT,
					Content: "world",
				},
			},
			expressions,
		)
	})
}
