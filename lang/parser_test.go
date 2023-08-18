package lang

import (
	"testing"

	_ "github.com/SushiWaUmai/prince/commands"
	"github.com/SushiWaUmai/prince/db"
	"github.com/SushiWaUmai/prince/env"
	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	t.Run("Test Parser with simple command", func(t *testing.T) {
		assert := assert.New(t)
		sample := "!echo hello | echo \"hello world\""

		tokens := Lex(sample)

		expressions, err := Parse(tokens)
		assert.Nil(err)

		assert.Equal([]Expression{
			{
				Type:    COMMAND,
				Content: "echo",
			},
			{
				Type:    ARGUMENT,
				Content: "hello",
			},
			{
				Type:    COMMAND,
				Content: "echo",
			},
			{
				Type:    ARGUMENT,
				Content: "hello world",
			},
		}, expressions)
	})

	t.Run("Test Parser with invalid prefix", func(t *testing.T) {
		assert := assert.New(t)
		sample := "echo asdf | echo test"

		tokens := Lex(sample)

		expressions, err := Parse(tokens)
		assert.Nil(expressions)
		assert.NotNil(err)
	})

	t.Run("Test Parser with invalid pipe", func(t *testing.T) {
		assert := assert.New(t)
		sample := string(env.BOT_PREFIX) + "echo asdf | echo hello |"

		tokens := Lex(sample)

		expressions, err := Parse(tokens)
		assert.Nil(expressions)
		assert.NotNil(err)
	})

	t.Run("Test Parser with no argument", func(t *testing.T) {
		assert := assert.New(t)
		sample := string(env.BOT_PREFIX) + "download"

		tokens := Lex(sample)
		expressions, err := Parse(tokens)
		assert.Nil(err)

		assert.Equal(
			[]Expression{
				{
					Type:    COMMAND,
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

		tokens := Lex(sample)
		expressions, err := Parse(tokens)
		assert.Nil(err)

		assert.Equal(
			[]Expression{
				{
					Type:    COMMAND,
					Content: "echo",
				},
				{
					Type:    ARGUMENT,
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

		tokens := Lex(sample)
		expressions, err := Parse(tokens)
		assert.Nil(err)

		assert.Equal(
			[]Expression{
				{
					Type:    COMMAND,
					Content: "echo",
				},
				{
					Type:    ARGUMENT,
					Content: "hello",
				},
				{
					Type:    COMMAND,
					Content: "chat",
				},
				{
					Type:    ARGUMENT,
					Content: "how are you",
				},
				{
					Type:    COMMAND,
					Content: "echo",
				},
				{
					Type:    ARGUMENT,
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
		sample := "!echo hello | chat hello world"

		tokens := Lex(sample)
		expressions, err := Parse(tokens)
		assert.Nil(err)

		assert.Equal(
			[]Expression{
				{
					Type:    COMMAND,
					Content: "echo",
				},
				{
					Type:    ARGUMENT,
					Content: "hello",
				},
				{
					Type:    COMMAND,
					Content: "chat",
				},
				{
					Type:    ARGUMENT,
					Content: "hello",
				},
				{
					Type:    ARGUMENT,
					Content: "world",
				},
			},
			expressions,
		)
	})
}
