package cmd

type MemoryAccessCommand interface {
	Command
	MemoryAccessOpLiteral() string
}

type PushCommand struct {
	Op      string
	Segment string
	Index   int16
}

func (p *PushCommand) Type() COMMAND_TYPE {
	return C_PUSH
}
func (p *PushCommand) MemoryAccessOpLiteral() string {
	return "push"
}

type PopCommand struct {
	Op      string
	Segment string
	Index   int16
}

func (p *PopCommand) Type() COMMAND_TYPE {
	return C_POP
}
func (p *PopCommand) MemoryAccessOpLiteral() string {
	return "pop"
}
