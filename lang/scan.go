package lang

import "github.com/SushiWaUmai/prince/utils"

func Scan(content string) ([]utils.CommandInput, error) {
	tokens := Lex(content)
	expressions, err := Parse(tokens)

	if err != nil {
		return nil, err
	}

	commandInput, err := Compile(expressions)
	if err != nil {
		return nil, err
	}

	return commandInput, nil
}
