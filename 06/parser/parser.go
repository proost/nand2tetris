package parser

import (
	"assembler/code"
	"assembler/symbol"
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
	file    *os.File
	scanner *bufio.Scanner

	commandType COMMAND_TYPE

	symbol string
	dest   string
	comp   string
	jump   string

	addressCounter int // This is for new variable in A-Instruction
	lineCounter    int // This is for location of label
	symbolTable    *symbol.SymbolTable

	phase int

	binaryCode string
}

func New(file *os.File) *Parser {
	return &Parser{
		file:           file,
		scanner:        bufio.NewScanner(file),
		addressCounter: 16,
		lineCounter:    0,
		symbolTable:    symbol.New(),
		phase:          1,
	}
}

func (p *Parser) HasMoreCommands() bool {
	return p.scanner.Scan()
}

func (p *Parser) Advance() {
	if p.phase == 1 {
		p.parseNextOnPhase1()
	} else if p.phase == 2 {
		p.parseNextOnPhase2()
	}
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

func (p *Parser) Rewind() {
	// Rewind Parser for phase 2

	_, err := p.file.Seek(0, 0)
	if err != nil {
		log.Fatalf("Can't start phase 2")
	}

	p.scanner = bufio.NewScanner(p.file)

	p.phase = 2
}

func (p *Parser) parseNextOnPhase1() {
	command := p.scanner.Text()

	command = strings.TrimSpace(command)

	if isAbleToSkip(command) {
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

		p.lineCounter++
	} else if p.commandType == L_COMMAND {
		p.setSymbol(command)

		p.dest = ""
		p.comp = ""
		p.jump = ""

		// update label location
		sym, err := p.parseSymbol(command)
		if err != nil {
			log.Fatalf(err.Error())
		}

		p.symbolTable.AddEntry(sym, p.lineCounter)
	} else if p.commandType == C_COMMAND {
		p.symbol = ""

		p.setDestCompJumpWhenCInstruction(command)

		p.lineCounter++
	}
}

func (p *Parser) parseNextOnPhase2() {
	command := p.scanner.Text()

	command = strings.TrimSpace(command)

	if isAbleToSkip(command) {
		p.symbol = ""
		p.binaryCode = ""

		return
	}

	command = p.trimComment(command)

	p.setCommandType(command)

	if p.commandType == A_COMMAND {
		sym, err := p.parseSymbol(command)
		if err != nil {
			log.Fatalf(err.Error())
		}

		_, err = strconv.Atoi(sym)
		if err == nil {
			// A-Instruction does not have a symbol

			p.setBinaryCodeWhenAInstruction(command)

			return
		}

		address, isExist := p.symbolTable.GetAddress(sym)
		if !isExist {
			// If symbol is new variable, set variable to memory address
			address = p.addressCounter

			p.symbolTable.AddEntry(sym, address)

			p.addressCounter++
		}

		p.setSymbol(command)

		p.dest = ""
		p.comp = ""
		p.jump = ""

		p.binaryCode = binaryStringToByteArray(strconv.Itoa(address))
	} else if p.commandType == L_COMMAND {
		p.setSymbol(command)

		p.dest = ""
		p.comp = ""
		p.jump = ""

		p.binaryCode = ""
	} else if p.commandType == C_COMMAND {
		p.symbol = ""

		p.setDestCompJumpWhenCInstruction(command)

		p.setBinaryCodeWhenCInstruction()
	}
}

func (p *Parser) trimComment(command string) string {
	return strings.TrimSpace(strings.Split(command, "//")[0])
}

func (p *Parser) setCommandType(command string) {
	if command[0] == '@' {
		p.commandType = A_COMMAND
	} else if command[0] == '(' && command[len(command)-1] == ')' {
		p.commandType = L_COMMAND
		return
	} else {
		p.commandType = C_COMMAND
	}
}

func (p *Parser) parseSymbol(command string) (string, error) {
	if command[0] == '@' {
		return command[1:], nil
	}

	if command[0] == '(' && command[len(command)-1] == ')' {
		return command[1 : len(command)-1], nil
	}

	return "", fmt.Errorf("can't parse symbol: %s", command)
}

func (p *Parser) setSymbol(command string) {
	if command[0] == '@' {
		p.symbol = command[1:]
	}

	if command[0] == '(' && command[len(command)-1] == ')' {
		p.symbol = command[1 : len(command)-1]
	}
}

func (p *Parser) setBinaryCodeWhenAInstruction(command string) {
	command = command[1:]

	p.binaryCode = binaryStringToByteArray(command)
}

func binaryStringToByteArray(command string) string {
	i, err := strconv.Atoi(command)
	if err != nil {
		log.Fatalf("Can't convert %s to integer", command)
	}

	return fmt.Sprintf("%016b", i)
}

func (p *Parser) setDestCompJumpWhenCInstruction(command string) {
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
			log.Fatalf("Can't find dest mnemonic: %+v", p.dest)
		} else {
			dest, _ = code.Dest("null0")
		}
	}

	comp, ok := code.Comp(p.comp)
	if !ok && p.comp != "" {
		log.Fatalf("Can't find comp mnemonic: %+v, phase: %d", p.comp, p.phase)
	}

	jump, ok := code.Jump(p.jump)
	if !ok {
		if p.jump != "" {
			log.Fatalf("Can't find jump mnemonic: %+v", p.jump)
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

func isAbleToSkip(command string) bool {
	if len(command) == 0 {
		// empty line
		return true
	}

	if strings.HasPrefix(command, "//") {
		// comment-line

		return true
	}

	return false
}
