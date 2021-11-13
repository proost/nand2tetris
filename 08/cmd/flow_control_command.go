package cmd

type FlowControlCommand interface {
	Command
	FlowControlOpLiteral() string
}

type LabelCommand struct {
	Op    string
	Label string
}

func (l *LabelCommand) Type() COMMAND_TYPE {
	return C_LABEL
}
func (l *LabelCommand) FlowControlOpLiteral() string {
	return "label"
}

type GotoCommand struct {
	Op    string
	Label string
}

func (g *GotoCommand) Type() COMMAND_TYPE {
	return C_GOTO
}
func (g *GotoCommand) FlowControlOpLiteral() string {
	return "goto"
}

type IfGotoCommand struct {
	Op    string
	Label string
}

func (i *IfGotoCommand) Type() COMMAND_TYPE {
	return C_GOTO
}
func (i *IfGotoCommand) FlowControlOpLiteral() string {
	return "if-goto"
}
