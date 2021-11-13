package tokenizer

import (
	"bufio"
	"log"
	"os"
	"strings"
)

type Tokenizer struct {
	file    *os.File
	scanner *bufio.Scanner

	tokens []string
}

func New(file *os.File) *Tokenizer {
	return &Tokenizer{
		file:    file,
		scanner: bufio.NewScanner(file),
	}
}

func (t *Tokenizer) HasMoreCommands() bool {
	return t.scanner.Scan()
}

func (t *Tokenizer) ReadNextCommand() {
	cmd := t.scanner.Text()

	cmd = strings.TrimSpace(cmd)
	if t.isAbleToSkip(cmd) {
		t.tokens = []string{}

		return
	}

	cmd = t.trimComment(cmd)
	t.tokens = strings.Split(cmd, " ")
}

func (t *Tokenizer) isAbleToSkip(command string) bool {
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

func (t *Tokenizer) trimComment(command string) string {
	return strings.TrimSpace(strings.Split(command, "//")[0])
}

func (t *Tokenizer) HasMoreTokens() bool {
	return len(t.tokens) > 0
}

func (t *Tokenizer) NextToken() (string, bool) {
	if !t.HasMoreTokens() {
		return "", false
	}

	next := t.tokens[0]

	t.tokens = t.tokens[1:]

	return next, true
}

func (t *Tokenizer) Rewind() {
	_, err := t.file.Seek(0, 0)
	if err != nil {
		log.Fatalf("Can't start phase 2")
	}

	t.scanner = bufio.NewScanner(t.file)
}
