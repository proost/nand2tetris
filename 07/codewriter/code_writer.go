package codewriter

import (
	"log"
	"nand2tetris/07/cmd"
	"strconv"
)

type CodeWriter struct {
	writer      *fileWriter
	labelNumber int
}

func New(asmFileName string) *CodeWriter {
	return &CodeWriter{writer: newFileWriter(asmFileName)}
}

func (w *CodeWriter) SetFileName(fileName string) {
	w.writer.changeFile(fileName)
}

func (w *CodeWriter) WriteAssembly(command cmd.Command) {
	switch command.Type() {
	case cmd.C_ARITHMETIC:
		arithmeticCmd, ok := command.(cmd.ArithmeticCommand)
		if !ok {
			log.Fatalf("command is not ArithmeticCommand, %s", arithmeticCmd.Type())
		}
		w.WriteArithmetic(arithmeticCmd)
	case cmd.C_PUSH:
		fallthrough
	case cmd.C_POP:
		memoryCmd, ok := command.(cmd.MemoryAccessCommand)
		if !ok {
			log.Fatalf("command is not MemoryAccessCommand, %s", memoryCmd.Type())
		}
		w.WritePushPop(memoryCmd)
	}
}

func (w *CodeWriter) WriteArithmetic(command cmd.ArithmeticCommand) {
	switch command.ArithmeticOpLiteral() {
	case "add":
		w.pop()
		w.writer.writeString("@SP\n")
		w.writer.writeString("A=M-1\n")
		w.writer.writeString("M=M+D\n")
	case "sub":
		w.pop()
		w.writer.writeString("@SP\n")
		w.writer.writeString("A=M-1\n")
		w.writer.writeString("M=M-D\n")
	case "neg":
		w.writer.writeString("@SP\n")
		w.writer.writeString("A=M-1\n")
		w.writer.writeString("M=-M\n")
	case "eq":
		label := "FALSE" + strconv.Itoa(w.labelNumber)
		w.labelNumber++

		w.pop()
		w.writer.writeString("@SP\n")
		w.writer.writeString("A=M-1\n")
		w.writer.writeString("D=D-M\n")
		w.writer.writeString("@" + label + "\n")
		w.writer.writeString("D;JEQ\n")
		w.writer.writeString("D=1\n")
		w.writer.writeString("(" + label + ")\n")
		w.writer.writeString("@SP\n")
		w.writer.writeString("A=M-1\n")
		w.writer.writeString("M=!D\n")
	case "gt":
		falseLabel := "FALSE" + strconv.Itoa(w.labelNumber)
		trueLabel := "TRUE" + strconv.Itoa(w.labelNumber)
		w.labelNumber++

		w.pop()
		w.writer.writeString("@SP\n")
		w.writer.writeString("A=M-1\n")
		w.writer.writeString("D=M-D\n")
		w.writer.writeString("@" + trueLabel + "\n")
		w.writer.writeString("D;JGT\n")
		w.writer.writeString("D=0\n")
		w.writer.writeString("D=!D\n")
		w.writer.writeString("@" + falseLabel + "\n")
		w.writer.writeString("0;JMP\n")
		w.writer.writeString("(" + trueLabel + ")\n")
		w.writer.writeString("D=0\n")
		w.writer.writeString("(" + falseLabel + ")\n")
		w.writer.writeString("@SP\n")
		w.writer.writeString("A=M-1\n")
		w.writer.writeString("M=!D\n")
	case "lt":
		falseLabel := "FALSE" + strconv.Itoa(w.labelNumber)
		trueLabel := "TRUE" + strconv.Itoa(w.labelNumber)
		w.labelNumber++

		w.pop()
		w.writer.writeString("@SP\n")
		w.writer.writeString("A=M-1\n")
		w.writer.writeString("D=M-D\n")
		w.writer.writeString("@" + trueLabel + "\n")
		w.writer.writeString("D;JLT\n")
		w.writer.writeString("D=0\n")
		w.writer.writeString("D=!D\n")
		w.writer.writeString("@" + falseLabel + "\n")
		w.writer.writeString("0;JMP\n")
		w.writer.writeString("(" + trueLabel + ")\n")
		w.writer.writeString("D=0\n")
		w.writer.writeString("(" + falseLabel + ")\n")
		w.writer.writeString("@SP\n")
		w.writer.writeString("A=M-1\n")
		w.writer.writeString("M=!D\n")
	case "and":
		w.pop()
		w.writer.writeString("@SP\n")
		w.writer.writeString("A=M-1\n")
		w.writer.writeString("M=M&D\n")
	case "or":
		w.pop()
		w.writer.writeString("@SP\n")
		w.writer.writeString("A=M-1\n")
		w.writer.writeString("M=M|D\n")
	case "not":
		w.pop()
		w.writer.writeString("@SP\n")
		w.writer.writeString("A=M\n")
		w.writer.writeString("M=!D\n")
	}
}

func (w *CodeWriter) WritePushPop(command cmd.MemoryAccessCommand) {
	switch command.Type() {
	case cmd.C_PUSH:
		pushCmd := command.(*cmd.PushCommand)
		w.writePush(pushCmd)
	case cmd.C_POP:
		popCmd := command.(*cmd.PopCommand)
		w.writePop(popCmd)
	}
}

func (w *CodeWriter) Close() {
	w.writer.close()
}

// pop stack and put the value into register D
func (w *CodeWriter) pop() {
	w.writer.writeString("@SP\n")
	w.writer.writeString("M=M-1\n")
	w.writer.writeString("A=M\n")
	w.writer.writeString("D=M\n")
}

func (w *CodeWriter) push() {
	w.writer.writeString("@SP\n")
	w.writer.writeString("A=M\n")
	w.writer.writeString("M=D\n")
	w.writer.writeString("@SP\n")
	w.writer.writeString("M=M+1\n")
}

func (w *CodeWriter) writePush(command *cmd.PushCommand) {
	switch command.Segment {
	case "constant":
		w.writer.writeString("@" + strconv.Itoa(int(command.Index)) + "\n")
		w.writer.writeString("D=A\n")
		w.push()
	case "argument":
		w.getSegmentDataFromMemory("ARG")
		w.pushFromMemory(int(command.Index))
	case "local":
		w.getSegmentDataFromMemory("LCL")
		w.pushFromMemory(int(command.Index))
	case "this":
		w.getSegmentDataFromMemory("THIS")
		w.pushFromMemory(int(command.Index))
	case "that":
		w.getSegmentDataFromMemory("THAT")
		w.pushFromMemory(int(command.Index))
	case "temp":
		w.getDataByMemoryAddress(5)
		w.pushFromMemory(int(command.Index))
	case "pointer":
		w.getDataByMemoryAddress(3)
		w.pushFromMemory(int(command.Index))
	case "static":
		w.writer.writeString("@" + w.writer.fileName + "." + strconv.Itoa(int(command.Index)) + "\n")
		w.writer.writeString("M=D\n")
		w.push()
	}
}

func (w *CodeWriter) getSegmentDataFromMemory(segment string) {
	w.writer.writeString("@" + segment + "\n")
	w.writer.writeString("D=M\n")
}

func (w *CodeWriter) getDataByMemoryAddress(memoryAddr int) {
	w.writer.writeString("@" + strconv.Itoa(memoryAddr) + "\n")
	w.writer.writeString("D=A\n")
}

func (w *CodeWriter) pushFromMemory(index int) {
	w.writer.writeString("@" + strconv.Itoa(index) + "\n")
	w.writer.writeString("A=D+A\n")
	w.writer.writeString("D=M\n")
	w.push()
}

func (w *CodeWriter) writePop(command *cmd.PopCommand) {
	switch command.Segment {
	case "argument":
		w.getSegmentDataFromMemory("ARG")
		w.loadToMemory(int(command.Index))
	case "local":
		w.getSegmentDataFromMemory("LCL")
		w.loadToMemory(int(command.Index))
	case "this":
		w.getSegmentDataFromMemory("THIS")
		w.loadToMemory(int(command.Index))
	case "that":
		w.getSegmentDataFromMemory("THAT")
		w.loadToMemory(int(command.Index))
	case "temp":
		w.getDataByMemoryAddress(5)
		w.loadToMemory(int(command.Index))
	case "pointer":
		w.getDataByMemoryAddress(3)
		w.loadToMemory(int(command.Index))
	case "static":
		w.pop()
		w.writer.writeString("@" + w.writer.fileName + "." + strconv.Itoa(int(command.Index)) + "\n")
		w.writer.writeString("M=D\n")
	}
}

func (w *CodeWriter) loadToMemory(index int) {
	w.writer.writeString("@13\n")
	w.writer.writeString("M=D\n")
	w.writer.writeString("@" + strconv.Itoa(index) + "\n")
	w.writer.writeString("D=A\n")
	w.writer.writeString("@13\n")
	w.writer.writeString("M=M+D\n")
	w.writer.writeString("@SP\n")
	w.writer.writeString("M=M-1\n")
	w.writer.writeString("A=M\n")
	w.writer.writeString("D=M\n")
	w.writer.writeString("@13\n")
	w.writer.writeString("A=M\n")
	w.writer.writeString("M=D\n")
}
