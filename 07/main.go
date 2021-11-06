package main

import (
	"flag"
	"log"
	"nand2tetris/07/codewriter"
	"nand2tetris/07/parser"
	"os"
	"strings"
)

var dirLoc = flag.String("dir", "", ".vm files location")
var fileLoc = flag.String("file", "", "a .vm file location")

func init() {
	flag.Parse()

	if *dirLoc == "" && *fileLoc == "" {
		log.Fatalf("dir option and file option can't be both empty")
	}
	if *dirLoc != "" && *fileLoc != "" {
		log.Fatalf("dir option and file option can't be both not empty")
	}

	if *fileLoc != "" {
		validateFileFormat(fileLoc, "vm")
	}
}

func validateFileFormat(name *string, format string) {
	temp := strings.Split(*name, ".")
	if temp[len(temp)-1] != format {
		log.Fatalf("Format of %s must be %s", *name, format)
	}
}

func main() {

	if *dirLoc != "" {
		list, err := os.ReadDir(*dirLoc)
		if err != nil {
			log.Fatalf("can't read %s", *dirLoc)
		}

		dir := *dirLoc
		if dir[len(dir)-1] == '\\' || (dir)[len(dir)-1] == '/' {
			dir = dir[:len(dir)-1]
		}

		for _, el := range list {
			if !el.IsDir() && checkFileFormat(el.Name(), "vm") {
				vmFileName := dir + string(os.PathSeparator) + el.Name()
				vmFile, err := os.Open(vmFileName)
				if err != nil {
					log.Fatalf("can't open file: %s", vmFileName)
				}

				asmFileName := createAssemblyFileName(dir, el.Name())

				p := parser.New(vmFile)
				w := codewriter.New(asmFileName)
				for p.HasMoreCommands() {
					p.Advance()

					command := p.Command()

					if command == nil {
						continue
					}

					w.WriteAssembly(command)
				}

				vmFile.Close()
				w.Close()
			}
		}
	} else if *fileLoc != "" {
		vmFile, err := os.Open(*fileLoc)
		if err != nil {
			log.Fatalf("can't open file: %s", *fileLoc)
		}
		defer vmFile.Close()

		dir := getDirFromFileLoc(*fileLoc)
		vmFileName := getFileNameFromFileLoc(*fileLoc)

		asmFileName := createAssemblyFileName(dir, vmFileName)

		p := parser.New(vmFile)
		w := codewriter.New(asmFileName)
		for p.HasMoreCommands() {
			p.Advance()

			command := p.Command()

			if command == nil {
				continue
			}

			w.WriteAssembly(command)
		}

		vmFile.Close()
		w.Close()
	}
}

func checkFileFormat(fileName string, format string) bool {
	temp := strings.Split(fileName, ".")
	return temp[len(temp)-1] == format
}

func createAssemblyFileName(dir string, vmFileName string) string {
	s := strings.Split(vmFileName, ".")

	return dir + string(os.PathSeparator) + strings.Join(s[:len(s)-1], "") + ".asm"
}

func getDirFromFileLoc(fileLoc string) string {
	s := strings.Split(fileLoc, "/")
	return strings.Join(s[:len(s)-1], "/")
}

func getFileNameFromFileLoc(fileLoc string) string {
	s := strings.Split(fileLoc, "/")
	return s[len(s)-1]
}
