package lang

import (
	"errors"

	"github.com/SushiWaUmai/prince/db"
	"github.com/SushiWaUmai/prince/utils"
)

type ExpressionType int64

const MAX_ALIAS_DEPTH = 8

const (
	COMMAND ExpressionType = iota
	ARGUMENT
)

type Expression struct {
	Type    ExpressionType
	Content string
}

func CheckCommand(name string) bool {
	_, ok := utils.CommandMap[name]
	return ok
}

func FindAlias(name string) ([]Token, error) {
	ok := CheckCommand(name)

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
			return nil, errors.New("Invalid Syntax: Token \"" + tokens[i].Lexme + "\" is not a SEPARATOR")
		}

		i++
		if i >= length || tokens[i].Type != IDENTIFIER {
			return nil, errors.New("Invalid Syntax: Token \"" + tokens[i].Lexme + "\" is not a IDENTIFIER")
		}

		for j := 0; j < MAX_ALIAS_DEPTH; j++ {
			aliasTokens, err := FindAlias(tokens[i].Lexme)
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
		}

		ok := CheckCommand(tokens[i].Lexme)
		if !ok {
			return nil, errors.New("Command: \"" + tokens[i].Lexme + "\" not found")
		}

		expressions = append(expressions, Expression{
			Type:    COMMAND,
			Content: tokens[i].Lexme,
		})

		i++
		if i >= length {
			return nil, errors.New("Invalid syntax: Has to end with EOF")
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
