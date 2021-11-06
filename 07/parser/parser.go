package parser

import (
	"log"
	"nand2tetris/07/cmd"
	"nand2tetris/07/tokenizer"
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

	command, isExist := p.tokenizer.NextToken() // first token must be type of command
	if !isExist {
		p.arg1 = ""
		p.arg2 = -1
		p.commandType = ""

		return
	}

	// parse command type & set command type
	cmdType := p.parseCommandType(command)
	p.commandType = cmdType

	switch cmdType {
	case cmd.C_ARITHMETIC:
		p.currCmd = p.parseArithmeticCommand(command)
	case cmd.C_PUSH:
		p.currCmd = p.parseMemoryAccessCommand(command)
	case cmd.C_POP:
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
	default:
		return cmd.UNKNOWN_COMMAND_TYPE
	}
}

func (p *Parser) parseArithmeticCommand(command string) cmd.Command {
	p.arg1 = command
	p.arg2 = -1

	switch command {
	case "neg":
		return &cmd.NegCommand{Op: command}
	case "add":
		return &cmd.AddCommand{Op: command}
	case "sub":
		return &cmd.SubCommand{Op: command}
	case "eq":
		return &cmd.EqCommand{Op: command}
	case "gt":
		return &cmd.GtCommand{Op: command}
	case "lt":
		return &cmd.LtCommand{Op: command}
	case "and":
		return &cmd.AndCommand{Op: command}
	case "or":
		return &cmd.OrCommand{Op: command}
	case "not":
		return &cmd.NotCommand{Op: command}
	default:
		return &cmd.InvalidCommand{
			Op:   command,
			Arg1: "",
			Arg2: "",
		}
	}
}

func (p *Parser) parseMemoryAccessCommand(command string) cmd.Command {
	segment, isExist := p.tokenizer.NextToken()
	if !isExist {
		log.Fatalf("%s command needs arg1", command)
	}
	p.arg1 = segment

	arg2, isExist := p.tokenizer.NextToken()
	if !isExist {
		log.Fatalf("%s command needs arg2", command)
	}
	index, err := strconv.Atoi(arg2)
	if err != nil {
		log.Fatalf("arg2 of %s command must be int, not %T", command, index)
	}
	if index < 0 {
		log.Fatalf("arg2 can't be smaller than 0, when %s command", command)
	}
	p.arg2 = int16(index)

	switch command {
	case "push":
		return &cmd.PushCommand{
			Op:      command,
			Segment: segment,
			Index:   int16(index),
		}
	default:
		return &cmd.InvalidCommand{
			Op:   command,
			Arg1: segment,
			Arg2: arg2,
		}
	}
}
