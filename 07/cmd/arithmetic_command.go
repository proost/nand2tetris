package cmd

type ArithmeticCommand interface {
	Command
	ArithmeticOpLiteral() string
}

type AddCommand struct {
	Op string
}

func (a *AddCommand) Type() COMMAND_TYPE {
	return C_ARITHMETIC
}
func (a *AddCommand) ArithmeticOpLiteral() string {
	return a.Op
}

type SubCommand struct {
	Op string
}

func (s *SubCommand) Type() COMMAND_TYPE {
	return C_ARITHMETIC
}
func (s *SubCommand) ArithmeticOpLiteral() string {
	return s.Op
}

type NegCommand struct {
	Op string
}

func (n *NegCommand) Type() COMMAND_TYPE {
	return C_ARITHMETIC
}
func (n *NegCommand) ArithmeticOpLiteral() string {
	return n.Op
}

type EqCommand struct {
	Op string
}

func (e *EqCommand) Type() COMMAND_TYPE {
	return C_ARITHMETIC
}
func (e *EqCommand) ArithmeticOpLiteral() string {
	return e.Op
}

type GtCommand struct {
	Op string
}

func (g *GtCommand) Type() COMMAND_TYPE {
	return C_ARITHMETIC
}
func (g *GtCommand) ArithmeticOpLiteral() string {
	return g.Op
}

type LtCommand struct {
	Op string
}

func (l *LtCommand) Type() COMMAND_TYPE {
	return C_ARITHMETIC
}
func (l *LtCommand) ArithmeticOpLiteral() string {
	return l.Op
}

type AndCommand struct {
	Op string
}

func (a *AndCommand) Type() COMMAND_TYPE {
	return C_ARITHMETIC
}
func (a *AndCommand) ArithmeticOpLiteral() string {
	return a.Op
}

type OrCommand struct {
	Op string
}

func (o *OrCommand) Type() COMMAND_TYPE {
	return C_ARITHMETIC
}
func (o *OrCommand) ArithmeticOpLiteral() string {
	return o.Op
}

type NotCommand struct {
	Op string
}

func (n *NotCommand) Type() COMMAND_TYPE {
	return C_ARITHMETIC
}
func (n *NotCommand) ArithmeticOpLiteral() string {
	return n.Op
}
