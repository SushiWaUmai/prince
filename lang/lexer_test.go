package lang

import (
	"testing"

	"github.com/SushiWaUmai/prince/env"
	"github.com/stretchr/testify/assert"
)

func TestLexer(t *testing.T) {
	t.Run("Test LexIdentifier", func(t *testing.T) {
		assert := assert.New(t)
		sample := "hello"

		tokens := []Token{}
		i := 0
		tokens, i, result := LexIdentifier(tokens, sample, i)
		assert.True(result)
		assert.Equal(5, i)
		assert.Equal([]Token{
			{
				Type:  IDENTIFIER,
				Lexme: "hello",
			},
		}, tokens)
	})

	t.Run("Test LexIdentifier with simple command", func(t *testing.T) {
		assert := assert.New(t)
		sample := "echo hello world"

		tokens := []Token{}
		i := 0
		tokens, i, result := LexIdentifier(tokens, sample, i)
		assert.True(result)
		assert.Equal(4, i)
		assert.Equal([]Token{
			{
				Type:  IDENTIFIER,
				Lexme: "echo",
			},
		}, tokens)
	})

	t.Run("Test LexPipe", func(t *testing.T) {
		assert := assert.New(t)
		sample := "|"

		tokens := []Token{}
		i := 0
		tokens, i, result := LexPipe(tokens, sample, i)
		assert.True(result)
		assert.Equal(0, i)
		assert.Equal([]Token{
			{
				Type:  SEPARATOR,
				Lexme: "|",
			},
		}, tokens)
	})

	t.Run("Test LexDoubleQuoteString", func(t *testing.T) {
		assert := assert.New(t)
		sample := "\"hello world\""

		tokens := []Token{}
		i := 0
		tokens, i, result := LexDoubleQuoteString(tokens, sample, i)
		assert.True(result)
		assert.Equal(12, i)
		assert.Equal([]Token{
			{
				Type:  IDENTIFIER,
				Lexme: "hello world",
			},
		}, tokens)
	})

	t.Run("Test LexSingleQuoteString", func(t *testing.T) {
		assert := assert.New(t)
		sample := "'hello world'"

		tokens := []Token{}
		i := 0
		tokens, i, result := LexSingleQuoteString(tokens, sample, i)
		assert.True(result)
		assert.Equal(12, i)
		assert.Equal([]Token{
			{
				Type:  IDENTIFIER,
				Lexme: "hello world",
			},
		}, tokens)
	})

	t.Run("Test Lex with simple command", func(t *testing.T) {
		assert := assert.New(t)
		sample := string(env.BOT_PREFIX) + "echo hello | echo \"hello world\""

		tokens := Lex(sample)
		assert.Equal([]Token{
			{
				Type:  SEPARATOR,
				Lexme: "!",
			},
			{
				Type:  IDENTIFIER,
				Lexme: "echo",
			},
			{
				Type:  IDENTIFIER,
				Lexme: "hello",
			},
			{
				Type:  SEPARATOR,
				Lexme: "|",
			},
			{
				Type:  IDENTIFIER,
				Lexme: "echo",
			},
			{
				Type:  IDENTIFIER,
				Lexme: "hello world",
			},
			{
				Type:  EOF,
				Lexme: "",
			},
		}, tokens)
	})

	t.Run("Test Lex with DoubleQuote inside SingleQuote", func(t *testing.T) {
		assert := assert.New(t)
		sample := string(env.BOT_PREFIX) + "echo 'hello \"world\"'"

		tokens := Lex(sample)
		assert.Equal([]Token{
			{
				Type:  SEPARATOR,
				Lexme: "!",
			},
			{
				Type:  IDENTIFIER,
				Lexme: "echo",
			},
			{
				Type:  IDENTIFIER,
				Lexme: "hello \"world\"",
			},
			{
				Type:  EOF,
				Lexme: "",
			},
		}, tokens)
	})

	t.Run("Test Lex with SingleQuote inside DoubleQuote", func(t *testing.T) {
		assert := assert.New(t)
		sample := string(env.BOT_PREFIX) + "echo \"hello 'world'\""

		tokens := Lex(sample)
		assert.Equal([]Token{
			{
				Type:  SEPARATOR,
				Lexme: "!",
			},
			{
				Type:  IDENTIFIER,
				Lexme: "echo",
			},
			{
				Type:  IDENTIFIER,
				Lexme: "hello 'world'",
			},
			{
				Type:  EOF,
				Lexme: "",
			},
		}, tokens)
	})
}
