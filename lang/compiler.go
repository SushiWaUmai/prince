package lang

import (
	"errors"

	"github.com/SushiWaUmai/prince/utils"
)

func Compile(expressions []Expression) ([]utils.CommandInput, error) {
	var result []utils.CommandInput
	length := len(expressions)

	for i := 0; i < length; i++ {
		if expressions[i].Type != COMMAND {
			return nil, errors.New("Could not compile expressions")
		}

		cmdName := expressions[i].Content
		i++

		var cmdArgs []string
		for ; i < length && expressions[i].Type == ARGUMENT; i++ {
			cmdArgs = append(cmdArgs, expressions[i].Content)
		}
		i--

		result = append(result, utils.CommandInput{
			Name: cmdName,
			Args: cmdArgs,
		})
	}

	return result, nil
}
