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
