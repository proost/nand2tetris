package main

import (
	"assembler/parser"
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

var assembly = flag.String("asm", "", "Assembly file location")

func init() {
	flag.Parse()

	assertBlank(assembly, "asm")

	validateFileFormat(assembly, "asm")
}

func assertBlank(loc *string, flag string) {
	if loc == nil || *loc == "" {
		log.Fatalf("%s can't be empty", flag)
	}
}

func validateFileFormat(name *string, format string) {
	temp := strings.Split(*name, ".")
	if temp[len(temp)-1] != format {
		log.Fatalf("Format of %s must be %s", *name, format)
	}
}

func main() {

	file, err := os.Open(*assembly)
	if err != nil {
		log.Fatalf("Can't be open file: %s", *assembly)
	}
	defer file.Close()

	parser := parser.New(file)
	result := make([]string, 0)

	// phase 1
	for parser.HasMoreCommands() {
		parser.Advance()
	}

	parser.Rewind()

	// phase 2
	for parser.HasMoreCommands() {
		parser.Advance()

		binCode := parser.BinaryCode()
		if binCode != "" {
			result = append(result, binCode)
		}
	}

	binFile := makeBinFile(*assembly)
	defer binFile.Close()

	writer := bufio.NewWriter(binFile)
	for _, binaryCode := range result {
		_, err := writer.WriteString(binaryCode + "\n")
		if err != nil {
			fmt.Printf("%v", err)
		}
	}
	writer.Flush()
}

func makeBinFile(asmFileLoc string) *os.File {
	var binFileName string

	if strings.Contains(asmFileLoc, "/") {
		loc := strings.Split(asmFileLoc, "/")
		asmFile := loc[len(loc)-1]
		binFileName = "./" + strings.Split(asmFile, ".")[0] + ".hack"
	} else if strings.Contains(asmFileLoc, "\\") {
		loc := strings.Split(asmFileLoc, "\\")
		asmFile := loc[len(loc)-1]
		binFileName = ".\\" + strings.Split(asmFile, ".")[0] + ".hack"
	}

	binFile, err := os.Create(binFileName)
	if err != nil {
		log.Fatalf("Can't create binary file")
	}
	return binFile
}
