package parser

import (
	"assembler/code"
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type COMMAND_TYPE string

const (
	A_COMMAND = COMMAND_TYPE("A_COMMAND")
	C_COMMAND = COMMAND_TYPE("C_COMMAND")
	L_COMMAND = COMMAND_TYPE("L_COMMAND")
)

type Parser struct {
	scanner *bufio.Scanner
	commandType COMMAND_TYPE
	symbol string
	dest string
	comp string
	jump string
	binaryCode string
}

func New(file *os.File) *Parser{
	return &Parser{scanner: bufio.NewScanner(file)}
}

func (p *Parser) HasMoreCommands() bool {
	return p.scanner.Scan()
}

func (p *Parser) Advance() {
	p.nextCommand()
}

func (p *Parser) CommandType() COMMAND_TYPE {
	return p.commandType
}

func (p *Parser) Symbol() (string, error) {
	if p.CommandType() == A_COMMAND {
		return p.symbol, nil
	}

	if p.CommandType() == L_COMMAND {
		return p.symbol, nil
	}

	return "", fmt.Errorf("invalid command type for Symbol: %s", p.CommandType())
}

func (p *Parser) Dest() (string, error) {
	if p.CommandType() == C_COMMAND {
		return p.dest, nil
	}

	return "", fmt.Errorf("invalid command type for Dest: %s", p.CommandType())
}

func (p *Parser) Comp() (string, error) {
	if p.CommandType() == C_COMMAND {
		return p.comp, nil
	}

	return "", fmt.Errorf("invalid command type for Comp: %s", p.CommandType())
}

func (p *Parser) Jump() (string, error) {
	if p.CommandType() == C_COMMAND {
		return p.jump, nil
	}

	return "", fmt.Errorf("invalid command type for Jump: %s", p.CommandType())
}

func (p *Parser) BinaryCode() string {
	return p.binaryCode
}

func (p *Parser) nextCommand() {
	command := p.scanner.Text()

	if len(command) == 0 {
		// empty line
		p.symbol = ""
		p.binaryCode = ""

		return
	}

	if strings.HasPrefix(command, "//") {
		// comment-line

		p.symbol = ""
		p.binaryCode = ""

		return
	}

	command = p.trimComment(command)

	p.setCommandType(command)

	if p.commandType == A_COMMAND {

		p.setSymbol(command)

		p.dest = ""
		p.comp = ""
		p.jump = ""

		p.setBinaryCodeWhenAInstruction(command)
	} else if p.commandType == L_COMMAND {
		p.setSymbol(command)

		p.dest = ""
		p.comp = ""
		p.jump = ""

		p.setBinaryCodeWhenLInstruction(command)
	} else if p.commandType == C_COMMAND {
		p.symbol = ""

		p.setCInstruction(command)

		p.setBinaryCodeWhenCInstruction()
	}
}

func (p *Parser) trimComment(command string) string {
	return strings.Split(command, "//")[0]
}

func (p *Parser) setCommandType(command string) {
	if command[0] == '@' {
		p.commandType = A_COMMAND
	} else if command[0]== '(' && command[len(command)-1] == ')' {
		p.commandType = L_COMMAND
		return
	} else {
		p.commandType = C_COMMAND
	}
}

func (p *Parser) setSymbol(command string) {
	if command[0] == '@' {
		p.symbol = command[1:]
	}

	if command[0]== '(' && command[len(command)-1] == ')' {
		p.symbol = command[1:len(command)-1]
	}
}

func (p *Parser) setBinaryCodeWhenAInstruction(command string) {
	command = command[1:]

	p.binaryCode = binaryStringToByteArray(command)
}

func (p *Parser) setBinaryCodeWhenLInstruction(command string) {
	command = command[1:len(command)-1]

	p.binaryCode = binaryStringToByteArray(command)
}

func binaryStringToByteArray(command string) string {
	i, err := strconv.Atoi(command)
	if err != nil {
		log.Fatalf("Can't convert %s to integer", command)
	}

	return fmt.Sprintf("%016b", i)
}

func (p *Parser) setCInstruction(command string) {
	var (
		destMnemonic string
		compMnemonic string
		jumpMnemonic string
	)

	if strings.Contains(command, "=") {
		s := strings.Split(command, "=")
		destMnemonic = s[0]
		command = s[1]
	}

	if strings.Contains(command, ";") {
		s := strings.Split(command, ";")
		compMnemonic = s[0]
		jumpMnemonic = s[1]
	} else {
		compMnemonic = command
	}

	p.dest = destMnemonic
	p.comp = compMnemonic
	p.jump = jumpMnemonic
}

func (p *Parser) setBinaryCodeWhenCInstruction() {
	dest, ok := code.Dest(p.dest)
	if !ok {
		if p.dest != "" {
			log.Fatalf("Can't find dest mnemonic: %+v", dest)
		} else {
			dest, _ = code.Dest("null0")
		}
	}

	comp, ok := code.Comp(p.comp)
	if !ok && p.comp != "" {
		log.Fatalf("Can't find comp mnemonic: %+v", comp)
	}

	jump, ok := code.Jump(p.jump)
	if !ok {
		if p.jump != "" {
			log.Fatalf("Can't find jump mnemonic: %+v", jump)
		} else {
			jump, _ = code.Jump("null")
		}
	}

	binaryCode := "111"
	binaryCode += comp
	binaryCode += dest
	binaryCode += jump

	p.binaryCode = binaryCode
}

