package lang

import (
	"errors"

	"github.com/SushiWaUmai/prince/db"
	"github.com/SushiWaUmai/prince/utils"
)

type ExpressionType int64

const (
	COMMAND ExpressionType = iota
	ARGUMENT
)

type Expression struct {
	Type    ExpressionType
	Content string
}

func CheckCommand(name string) ([]Token, error) {
	_, ok := utils.CommandMap[name]

	if !ok {
		// get alias
		alias := db.GetAlias(name)
		if alias == nil {
			return nil, errors.New("Command: \"" + name + "\" not found")
		}

		aliasToken := Lex(alias.Content)
		return aliasToken, nil
	}

	return nil, nil
}

func Parse(tokens []Token) ([]Expression, error) {
	i := 0
	var expressions []Expression
	length := len(tokens)

	// Pattern:
	// SEPARATOR COMMAND ARG ARG ARG... SEPARATOR COMMAND ARG...
	for tokens[i].Type != EOF {
		if tokens[i].Type != SEPARATOR {
			return nil, errors.New("Invalid syntax")
		}

		i++
		if i >= length || tokens[i].Type != IDENTIFIER {
			return nil, errors.New("Invalid syntax")
		}

		aliasTokens, err := CheckCommand(tokens[i].Lexme)
		// replace the current token with alias tokens
		if err == nil && aliasTokens != nil {
			// Remove EOF
			aliasTokens = aliasTokens[:len(aliasTokens)-1]

			// Remove current token
			tokens = append(tokens[:i], tokens[i+1:]...)
			// Insert alias tokens
			tokens = append(tokens[:i], append(aliasTokens, tokens[i:]...)...)

			length = len(tokens)
		} else if err != nil {
			return nil, err
		}

		expressions = append(expressions, Expression{
			Type:    COMMAND,
			Content: tokens[i].Lexme,
		})

		i++
		if i >= length {
			return nil, errors.New("Invalid syntax")
		}

		for ; tokens[i].Type == IDENTIFIER; i++ {
			expressions = append(expressions, Expression{
				Type:    ARGUMENT,
				Content: tokens[i].Lexme,
			})
		}
	}
	return expressions, nil
}
