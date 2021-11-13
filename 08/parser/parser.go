package parser

import (
	"log"
	"nand2tetris/projects/08/cmd"
	"nand2tetris/projects/08/tokenizer"
	"os"
	"strconv"
)

type Parser struct {
	tokenizer *tokenizer.Tokenizer

	arg1 string
	arg2 int16

	commandType cmd.COMMAND_TYPE

	currCmd cmd.Command
}

func New(file *os.File) *Parser {
	return &Parser{
		tokenizer: tokenizer.New(file),
	}
}

func (p *Parser) HasMoreCommands() bool {
	return p.tokenizer.HasMoreCommands()
}

func (p *Parser) Advance() {
	p.parseNextCommand()
}

func (p *Parser) CommandType() cmd.COMMAND_TYPE {
	return p.commandType
}

func (p *Parser) Arg1() string {
	if p.commandType == cmd.C_RETURN {
		log.Fatalf("invalid operation for command type: %s", p.commandType)
	}

	return p.arg1
}

func (p *Parser) Arg2() int16 {
	if p.commandType != cmd.C_PUSH &&
		p.commandType != cmd.C_POP &&
		p.commandType != cmd.C_FUNCTION &&
		p.commandType != cmd.C_CALL {
		log.Fatalf("invalid operation for command type: %s", p.commandType)
	}

	return p.arg2
}

func (p *Parser) Command() cmd.Command {
	return p.currCmd
}

func (p *Parser) parseNextCommand() {
	p.tokenizer.ReadNextCommand()

	op, isExist := p.tokenizer.NextToken() // first token must be type of command
	if !isExist {
		p.arg1 = ""
		p.arg2 = -1
		p.commandType = ""

		return
	}

	// parse command type & set command type
	cmdType := p.parseCommandType(op)
	p.commandType = cmdType

	switch cmdType {
	case cmd.C_ARITHMETIC:
		p.currCmd = p.parseArithmeticCommand(op)
	case cmd.C_PUSH:
		p.currCmd = p.parseMemoryAccessCommand(op)
	case cmd.C_POP:
		p.currCmd = p.parseMemoryAccessCommand(op)
	case cmd.C_LABEL:
		p.currCmd = p.parseFlowControlCommand(op)
	case cmd.C_GOTO:
		p.currCmd = p.parseFlowControlCommand(op)
	case cmd.C_FUNCTION:
		p.currCmd = p.parseFunctionCallCommand(op)
	case cmd.C_CALL:
		p.currCmd = p.parseFunctionCallCommand(op)
	case cmd.C_RETURN:
		p.currCmd = p.parseFunctionCallCommand(op)
	}
}

func (p *Parser) parseCommandType(command string) cmd.COMMAND_TYPE {
	switch command {
	case "add":
		return cmd.C_ARITHMETIC
	case "sub":
		return cmd.C_ARITHMETIC
	case "neg":
		return cmd.C_ARITHMETIC
	case "eq":
		return cmd.C_ARITHMETIC
	case "gt":
		return cmd.C_ARITHMETIC
	case "lt":
		return cmd.C_ARITHMETIC
	case "and":
		return cmd.C_ARITHMETIC
	case "or":
		return cmd.C_ARITHMETIC
	case "not":
		return cmd.C_ARITHMETIC
	case "push":
		return cmd.C_PUSH
	case "pop":
		return cmd.C_POP
	case "label":
		return cmd.C_LABEL
	case "goto":
		return cmd.C_GOTO
	case "if-goto":
		return cmd.C_IF
	case "function":
		return cmd.C_FUNCTION
	case "call":
		return cmd.C_CALL
	case "return":
		return cmd.C_RETURN
	default:
		return cmd.UNKNOWN_COMMAND_TYPE
	}
}

func (p *Parser) parseArithmeticCommand(op string) cmd.Command {
	p.arg1 = op
	p.arg2 = -1

	switch op {
	case "neg":
		return &cmd.NegCommand{Op: op}
	case "add":
		return &cmd.AddCommand{Op: op}
	case "sub":
		return &cmd.SubCommand{Op: op}
	case "eq":
		return &cmd.EqCommand{Op: op}
	case "gt":
		return &cmd.GtCommand{Op: op}
	case "lt":
		return &cmd.LtCommand{Op: op}
	case "and":
		return &cmd.AndCommand{Op: op}
	case "or":
		return &cmd.OrCommand{Op: op}
	case "not":
		return &cmd.NotCommand{Op: op}
	default:
		return &cmd.InvalidCommand{
			Op:   op,
			Arg1: "",
			Arg2: "",
		}
	}
}

func (p *Parser) parseMemoryAccessCommand(op string) cmd.Command {
	segment, isExist := p.tokenizer.NextToken()
	if !isExist {
		log.Fatalf("%s command needs arg1", op)
	}
	p.arg1 = segment

	arg2, isExist := p.tokenizer.NextToken()
	if !isExist {
		log.Fatalf("%s command needs arg2", op)
	}
	index, err := strconv.Atoi(arg2)
	if err != nil {
		log.Fatalf("arg2 of %s command must be int, not %T", op, index)
	}
	if index < 0 {
		log.Fatalf("arg2 can't be smaller than 0, when %s command", op)
	}
	p.arg2 = int16(index)

	switch op {
	case "push":
		return &cmd.PushCommand{
			Op:      op,
			Segment: segment,
			Index:   int16(index),
		}
	case "pop":
		return &cmd.PopCommand{
			Op:      op,
			Segment: segment,
			Index:   int16(index),
		}
	default:
		return &cmd.InvalidCommand{
			Op:   op,
			Arg1: segment,
			Arg2: arg2,
		}
	}
}

func (p *Parser) parseFlowControlCommand(op string) cmd.Command {
	label, isExist := p.tokenizer.NextToken()
	if !isExist {
		log.Fatalf("%s command needs arg1", op)
	}
	p.arg1 = label

	switch op {
	case "label":
		return &cmd.LabelCommand{
			Op:    op,
			Label: label,
		}
	case "goto":
		return &cmd.GotoCommand{
			Op:    op,
			Label: label,
		}
	case "if-goto":
		return &cmd.IfGotoCommand{
			Op:    op,
			Label: label,
		}
	default:
		return &cmd.InvalidCommand{
			Op:   op,
			Arg1: label,
			Arg2: "",
		}
	}
}

func (p *Parser) parseFunctionCallCommand(op string) cmd.Command {
	switch op {
	case "function":
		funcName, isExist := p.tokenizer.NextToken()
		if !isExist {
			log.Fatalf("%s command needs arg1", op)
		}
		p.arg1 = funcName

		arg2, isExist := p.tokenizer.NextToken()
		if !isExist {
			log.Fatalf("%s command needs arg2", op)
		}
		numOfArgs, err := strconv.Atoi(arg2)
		if err != nil {
			log.Fatalf("arg2 of %s command must be int, not %T", op, index)
		}
		if numOfArgs < 0 {
			log.Fatalf("arg2 can't be smaller than 0, when %s command", op)
		}
		p.arg2 = int16(numOfArgs)

		return &cmd.FunctionCommand{
			Op:        op,
			FuncName:  funcName,
			NumOfArgs: numOfArgs,
		}

	case "call":
		funcName, isExist := p.tokenizer.NextToken()
		if !isExist {
			log.Fatalf("%s command needs arg1", op)
		}
		p.arg1 = funcName

		arg2, isExist := p.tokenizer.NextToken()
		if !isExist {
			log.Fatalf("%s command needs arg2", op)
		}
		numOfArgs, err := strconv.Atoi(arg2)
		if err != nil {
			log.Fatalf("arg2 of %s command must be int, not %T", op, index)
		}
		if numOfArgs < 0 {
			log.Fatalf("arg2 can't be smaller than 0, when %s command", op)
		}
		p.arg2 = int16(numOfArgs)

		return &cmd.CallCommand{
			Op:        op,
			FuncName:  funcName,
			NumOfArgs: numOfArgs,
		}

	case "return":
		return &cmd.ReturnCommand{
			Op: op,
		}

	default:
		return &cmd.InvalidCommand{
			Op:   op,
			Arg1: "",
			Arg2: "",
		}
	}
}
