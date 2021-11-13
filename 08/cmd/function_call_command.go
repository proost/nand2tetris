package cmd

type FunctionCallCommand interface {
	Command
	FunctionOpLiteral() string
}

type FunctionCommand struct {
	Op        string
	FuncName  string
	NumOfArgs int
}

func (f *FunctionCommand) Type() COMMAND_TYPE {
	return C_FUNCTION
}
func (f *FunctionCommand) FunctionOpLiteral() string {
	return "function"
}

type CallCommand struct {
	Op        string
	FuncName  string
	NumOfArgs int
}

func (c *CallCommand) Type() COMMAND_TYPE {
	return C_CALL
}
func (c *CallCommand) FunctionOpLiteral() string {
	return "call"
}

type ReturnCommand struct {
	Op string
}

func (r *ReturnCommand) Type() COMMAND_TYPE {
	return C_RETURN
}
func (r *ReturnCommand) FunctionOpLiteral() string {
	return "return"
}
