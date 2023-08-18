package lang

import (
	"testing"

	_ "github.com/SushiWaUmai/prince/commands"
	"github.com/SushiWaUmai/prince/db"
	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	t.Run("Test Parser with simple command", func(t *testing.T) {
		assert := assert.New(t)
		sample := "!ping hello | ping \"hello world\""

		tokens := Lex(sample)

		expressions, err := Parse(tokens)
		assert.Nil(err)

		assert.Equal([]Expression{
			{
				Type:    COMMAND,
				Content: "ping",
			},
			{
				Type:    ARGUMENT,
				Content: "hello",
			},
			{
				Type:    COMMAND,
				Content: "ping",
			},
			{
				Type:    ARGUMENT,
				Content: "hello world",
			},
		}, expressions)
	})

	t.Run("Test Parser with invalid prefix", func(t *testing.T) {
		assert := assert.New(t)
		sample := "ping asdf | ping test"

		tokens := Lex(sample)

		expressions, err := Parse(tokens)
		assert.Nil(expressions)
		assert.NotNil(err)
	})

	t.Run("Test Parser with invalid pipe", func(t *testing.T) {
		assert := assert.New(t)
		sample := "!ping asdf | ping hello |"

		tokens := Lex(sample)

		expressions, err := Parse(tokens)
		assert.Nil(expressions)
		assert.NotNil(err)
	})

	t.Run("Test Parser with alias", func(t *testing.T) {
		assert := assert.New(t)

		db.CreateAlias("hello", "ping hello")
		defer db.DeleteAlias("hello")
		sample := "!hello"

		tokens := Lex(sample)
		expressions, err := Parse(tokens)
		assert.Nil(err)

		assert.Equal(
			[]Expression{
				{
					Type:    COMMAND,
					Content: "ping",
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

		db.CreateAlias("hello", "ping hello | chat \"how are you\"")
		defer db.DeleteAlias("hello")
		sample := "!hello | ping \"hello world\""

		tokens := Lex(sample)
		expressions, err := Parse(tokens)
		assert.Nil(err)

		assert.Equal(
			[]Expression{
				{
					Type:    COMMAND,
					Content: "ping",
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
					Content: "ping",
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

		db.CreateAlias("hello", "ping hello | chat \"how are you\"")
		defer db.DeleteAlias("hello")
		sample := "!ping hello | chat hello world"

		tokens := Lex(sample)
		expressions, err := Parse(tokens)
		assert.Nil(err)

		assert.Equal(
			[]Expression{
				{
					Type:    COMMAND,
					Content: "ping",
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
