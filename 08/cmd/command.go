package cmd

import (
	"fmt"
	"strings"
)

type COMMAND_TYPE string

const (
	C_ARITHMETIC         = COMMAND_TYPE("C_ARITHMETIC")
	C_PUSH               = COMMAND_TYPE("C_PUSH")
	C_POP                = COMMAND_TYPE("C_POP")
	C_LABEL              = COMMAND_TYPE("C_LABEL")
	C_GOTO               = COMMAND_TYPE("C_GOTO")
	C_IF                 = COMMAND_TYPE("C_IF")
	C_FUNCTION           = COMMAND_TYPE("C_FUNCTION")
	C_RETURN             = COMMAND_TYPE("C_RETURN")
	C_CALL               = COMMAND_TYPE("C_CALL")
	UNKNOWN_COMMAND_TYPE = COMMAND_TYPE("Unknown command type")
)

type Command interface {
	Type() COMMAND_TYPE
}

type InvalidCommand struct {
	Op   string
	Arg1 string
	Arg2 string
}

func (i *InvalidCommand) Type() COMMAND_TYPE {
	return UNKNOWN_COMMAND_TYPE
}
func (i *InvalidCommand) ErrorMessage() string {
	return fmt.Sprintf("Invalid command: %s", strings.Join([]string{i.Op, i.Arg1, i.Arg2}, " "))
}
