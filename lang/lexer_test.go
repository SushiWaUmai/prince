package lang_test

import (
	"testing"

	"github.com/SushiWaUmai/prince/env"
	"github.com/SushiWaUmai/prince/lang"
	"github.com/stretchr/testify/assert"
)

func TestLexer(t *testing.T) {
	t.Run("Test LexIdentifier", func(t *testing.T) {
		assert := assert.New(t)
		sample := []rune("hello")

		tokens := []lang.Token{}
		i := 0
		tokens, i, result := lang.LexIdentifier(tokens, sample, i)
		assert.True(result)
		assert.Equal(5, i)
		assert.Equal([]lang.Token{
			{
				Type:  lang.IDENTIFIER,
				Lexme: "hello",
			},
		}, tokens)
	})

	t.Run("Test LexIdentifier with simple command", func(t *testing.T) {
		assert := assert.New(t)
		sample := []rune("echo hello world")

		tokens := []lang.Token{}
		i := 0
		tokens, i, result := lang.LexIdentifier(tokens, sample, i)
		assert.True(result)
		assert.Equal(4, i)
		assert.Equal([]lang.Token{
			{
				Type:  lang.IDENTIFIER,
				Lexme: "echo",
			},
		}, tokens)
	})

	t.Run("Test LexPipe", func(t *testing.T) {
		assert := assert.New(t)
		sample := []rune("|")

		tokens := []lang.Token{}
		i := 0
		tokens, i, result := lang.LexPipe(tokens, sample, i)
		assert.True(result)
		assert.Equal(0, i)
		assert.Equal([]lang.Token{
			{
				Type:  lang.SEPARATOR,
				Lexme: "|",
			},
		}, tokens)
	})

	t.Run("Test LexDoubleQuoteString", func(t *testing.T) {
		assert := assert.New(t)
		sample := []rune("\"hello world\"")

		tokens := []lang.Token{}
		i := 0
		tokens, i, result := lang.LexDoubleQuoteString(tokens, sample, i)
		assert.True(result)
		assert.Equal(12, i)
		assert.Equal([]lang.Token{
			{
				Type:  lang.IDENTIFIER,
				Lexme: "hello world",
			},
		}, tokens)
	})

	t.Run("Test LexSingleQuoteString", func(t *testing.T) {
		assert := assert.New(t)
		sample := []rune("'hello world'")

		tokens := []lang.Token{}
		i := 0
		tokens, i, result := lang.LexSingleQuoteString(tokens, sample, i)
		assert.True(result)
		assert.Equal(12, i)
		assert.Equal([]lang.Token{
			{
				Type:  lang.IDENTIFIER,
				Lexme: "hello world",
			},
		}, tokens)
	})

	t.Run("Test LexTripleDoubleQuoteString", func(t *testing.T) {
		assert := assert.New(t)
		sample := []rune("\"\"\"hello world\"\"\"")

		tokens := []lang.Token{}
		i := 0
		tokens, i, result := lang.LexTripleDoubleQuoteString(tokens, sample, i)
		assert.True(result)
		assert.Equal(16, i)
		assert.Equal([]lang.Token{
			{
				Type:  lang.IDENTIFIER,
				Lexme: "hello world",
			},
		}, tokens)
	})

	t.Run("Test LexGraveString", func(t *testing.T) {
		assert := assert.New(t)
		sample := []rune("`hello world`")

		tokens := []lang.Token{}
		i := 0
		tokens, i, result := lang.LexGraveString(tokens, sample, i)
		assert.True(result)
		assert.Equal(12, i)
		assert.Equal([]lang.Token{
			{
				Type:  lang.IDENTIFIER,
				Lexme: "hello world",
			},
		}, tokens)
	})

	t.Run("Test Lex with simple command", func(t *testing.T) {
		assert := assert.New(t)
		sample := string(env.BOT_PREFIX) + "echo hello | echo \"hello world\""

		tokens := lang.Lex(sample)
		assert.Equal([]lang.Token{
			{
				Type:  lang.SEPARATOR,
				Lexme: "!",
			},
			{
				Type:  lang.IDENTIFIER,
				Lexme: "echo",
			},
			{
				Type:  lang.IDENTIFIER,
				Lexme: "hello",
			},
			{
				Type:  lang.SEPARATOR,
				Lexme: "|",
			},
			{
				Type:  lang.IDENTIFIER,
				Lexme: "echo",
			},
			{
				Type:  lang.IDENTIFIER,
				Lexme: "hello world",
			},
			{
				Type:  lang.EOF,
				Lexme: "",
			},
		}, tokens)
	})

	t.Run("Text Lex with TripleQuoteString", func(t *testing.T) {
		assert := assert.New(t)
		sample := string(env.BOT_PREFIX) + "echo \"\"\"hello world\"\"\" | ping"

		tokens := lang.Lex(sample)
		assert.Equal([]lang.Token{
			{
				Type:  lang.SEPARATOR,
				Lexme: string(env.BOT_PREFIX),
			},
			{
				Type:  lang.IDENTIFIER,
				Lexme: "echo",
			},
			{
				Type:  lang.IDENTIFIER,
				Lexme: "hello world",
			},
			{
				Type:  lang.SEPARATOR,
				Lexme: "|",
			},
			{
				Type:  lang.IDENTIFIER,
				Lexme: "ping",
			},
			{
				Type:  lang.EOF,
				Lexme: "",
			},
		}, tokens)
	})

	t.Run("Test Lex with DoubleQuote inside SingleQuote", func(t *testing.T) {
		assert := assert.New(t)
		sample := string(env.BOT_PREFIX) + "echo 'hello \"world\"'"

		tokens := lang.Lex(sample)
		assert.Equal([]lang.Token{
			{
				Type:  lang.SEPARATOR,
				Lexme: string(env.BOT_PREFIX),
			},
			{
				Type:  lang.IDENTIFIER,
				Lexme: "echo",
			},
			{
				Type:  lang.IDENTIFIER,
				Lexme: "hello \"world\"",
			},
			{
				Type:  lang.EOF,
				Lexme: "",
			},
		}, tokens)
	})

	t.Run("Test Lex with SingleQuote inside DoubleQuote", func(t *testing.T) {
		assert := assert.New(t)
		sample := string(env.BOT_PREFIX) + "echo \"hello 'world'\""

		tokens := lang.Lex(sample)
		assert.Equal([]lang.Token{
			{
				Type:  lang.SEPARATOR,
				Lexme: string(env.BOT_PREFIX),
			},
			{
				Type:  lang.IDENTIFIER,
				Lexme: "echo",
			},
			{
				Type:  lang.IDENTIFIER,
				Lexme: "hello 'world'",
			},
			{
				Type:  lang.EOF,
				Lexme: "",
			},
		}, tokens)
	})

	t.Run("Test Lex with DoubleQuote inside Grave", func(t *testing.T) {
		assert := assert.New(t)
		sample := string(env.BOT_PREFIX) + "echo `hello \"world\"`"

		tokens := lang.Lex(sample)
		assert.Equal([]lang.Token{
			{
				Type:  lang.SEPARATOR,
				Lexme: string(env.BOT_PREFIX),
			},
			{
				Type:  lang.IDENTIFIER,
				Lexme: "echo",
			},
			{
				Type:  lang.IDENTIFIER,
				Lexme: "hello \"world\"",
			},
			{
				Type:  lang.EOF,
				Lexme: "",
			},
		}, tokens)
	})

	t.Run("Test with non ASCII character", func(t *testing.T) {
		assert := assert.New(t)
		sample := string(env.BOT_PREFIX) + "echo こんにちは"

		tokens := lang.Lex(sample)
		assert.Equal([]lang.Token{
			{
				Type:  lang.SEPARATOR,
				Lexme: string(env.BOT_PREFIX),
			},
			{
				Type:  lang.IDENTIFIER,
				Lexme: "echo",
			},
			{
				Type:  lang.IDENTIFIER,
				Lexme: "こんにちは",
			},
			{
				Type:  lang.EOF,
				Lexme: "",
			},
		}, tokens)
	})

	t.Run("Test with emoji", func(t *testing.T) {
		assert := assert.New(t)
		sample := string(env.BOT_PREFIX) + "echo 🤔"

		tokens := lang.Lex(sample)
		assert.Equal([]lang.Token{
			{
				Type:  lang.SEPARATOR,
				Lexme: string(env.BOT_PREFIX),
			},
			{
				Type:  lang.IDENTIFIER,
				Lexme: "echo",
			},
			{
				Type:  lang.IDENTIFIER,
				Lexme: "🤔",
			},
			{
				Type:  lang.EOF,
				Lexme: "",
			},
		}, tokens)
	})

	t.Run("Test with short string", func(t *testing.T) {
		assert := assert.New(t)
		sampleDoubleQuote := string(env.BOT_PREFIX) + "echo \"a\""
		sampleSingleQuote := string(env.BOT_PREFIX) + "echo 'a'"
		sampleGraveQuote := string(env.BOT_PREFIX) + "echo `a`"
		sampleTripeDoubleQuote := string(env.BOT_PREFIX) + "echo \"\"\"a\"\"\""

		expected := []lang.Token{
			{
				Type:  lang.SEPARATOR,
				Lexme: string(env.BOT_PREFIX),
			},
			{
				Type:  lang.IDENTIFIER,
				Lexme: "echo",
			},
			{
				Type:  lang.IDENTIFIER,
				Lexme: "a",
			},
			{
				Type:  lang.EOF,
				Lexme: "",
			},
		}

		var tokens []lang.Token
		tokens = lang.Lex(sampleDoubleQuote)
		assert.Equal(expected, tokens)

		tokens = lang.Lex(sampleSingleQuote)
		assert.Equal(expected, tokens)

		tokens = lang.Lex(sampleGraveQuote)
		assert.Equal(expected, tokens)

		tokens = lang.Lex(sampleTripeDoubleQuote)
		assert.Equal(expected, tokens)
	})
}
