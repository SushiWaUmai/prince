package lang

import (
	"unicode"

	"github.com/SushiWaUmai/prince/env"
)

type TokenType int64

const (
	IDENTIFIER TokenType = iota
	SEPARATOR
	EOF
)

type Token struct {
	Type  TokenType
	Lexme string
}

func isAlphaNumeric(c []rune, i int) bool {
	return !unicode.IsSpace(c[i])
}

func isWhiteSpace(c []rune, i int) bool {
	return unicode.IsSpace(c[i])
}

func isDoubleQuotation(c []rune, i int) bool {
	return c[i] == '"'
}

func isSingleQuotation(c []rune, i int) bool {
	return c[i] == '\''
}

func isTripleDoubleQuotation(c []rune, i int) bool {
	if i+2 >= len(c) {
		return false
	}
	return c[i] == '"' && c[i+1] == '"' && c[i+2] == '"'
}

func isGrave(c []rune, i int) bool {
	return c[i] == '`'
}

func StringTilChar(content []rune, i int, startCond func(c []rune, i int) bool, endCond func(c []rune, i int) bool) (string, int) {
	var c rune

	c = content[i]
	length := len(content)
	lexme := ""

	if startCond(content, i) {
		c = content[i]
		for {
			lexme += string(c)
			i++
			if i >= length {
				break
			}

			c = content[i]
			if endCond(content, i) {
				lexme += string(c)
				break
			}
		}
	}

	return lexme, i
}

func LexIdentifier(tokens []Token, content []rune, i int) ([]Token, int, bool) {
	lexme, i := StringTilChar(content, i, isAlphaNumeric, isWhiteSpace)
	result := false
	if lexme != "" {
		if i < len(content) {
			lexme = lexme[:len(lexme)-1]
		}

		tokens = append(tokens, Token{
			Type:  IDENTIFIER,
			Lexme: lexme,
		})
		result = true
	}

	return tokens, i, result
}

func LexDoubleQuoteString(tokens []Token, content []rune, i int) ([]Token, int, bool) {
	lexme, i := StringTilChar(content, i, isDoubleQuotation, isDoubleQuotation)
	result := false
	if lexme != "" {
		lexme = lexme[1 : len(lexme)-1]
		tokens = append(tokens, Token{
			Type:  IDENTIFIER,
			Lexme: lexme,
		})
		result = true
	}

	return tokens, i, result
}

func LexSingleQuoteString(tokens []Token, content []rune, i int) ([]Token, int, bool) {
	lexme, i := StringTilChar(content, i, isSingleQuotation, isSingleQuotation)
	result := false
	if lexme != "" {
		lexme = lexme[1 : len(lexme)-1]
		tokens = append(tokens, Token{
			Type:  IDENTIFIER,
			Lexme: lexme,
		})
		result = true
	}

	return tokens, i, result
}

func LexTripleDoubleQuoteString(tokens []Token, content []rune, i int) ([]Token, int, bool) {
	lexme, i := StringTilChar(content, i, isTripleDoubleQuotation, isTripleDoubleQuotation)
	result := false
	if lexme != "" {
		lexme = lexme[3 : len(lexme)-1]
		i += 2
		tokens = append(tokens, Token{
			Type:  IDENTIFIER,
			Lexme: lexme,
		})
		result = true
	}

	return tokens, i, result
}

func LexGraveString(tokens []Token, content []rune, i int) ([]Token, int, bool) {
	lexme, i := StringTilChar(content, i, isGrave, isGrave)
	result := false
	if lexme != "" {
		lexme = lexme[1 : len(lexme)-1]
		tokens = append(tokens, Token{
			Type:  IDENTIFIER,
			Lexme: lexme,
		})
		result = true
	}

	return tokens, i, result
}

func LexPipe(tokens []Token, content []rune, i int) ([]Token, int, bool) {
	c := content[i]
	result := false
	if c == '|' || c == env.BOT_PREFIX {
		tokens = append(tokens, Token{
			Type:  SEPARATOR,
			Lexme: string(c),
		})
		result = true
	}

	return tokens, i, result
}

func Lex(contentStr string) []Token {
	var tokens []Token
	var result bool

	content := []rune(contentStr)
	length := len(content)

	for i := 0; i < length; i++ {
		tokens, i, result = LexPipe(tokens, content, i)
		if result {
			continue
		}

		tokens, i, result = LexTripleDoubleQuoteString(tokens, content, i)
		if result {
			continue
		}

		tokens, i, result = LexDoubleQuoteString(tokens, content, i)
		if result {
			continue
		}

		tokens, i, result = LexSingleQuoteString(tokens, content, i)
		if result {
			continue
		}

		tokens, i, result = LexGraveString(tokens, content, i)
		if result {
			continue
		}

		tokens, i, result = LexIdentifier(tokens, content, i)
		if result {
			continue
		}
	}
	tokens = append(tokens, Token{EOF, ""})

	return tokens
}
