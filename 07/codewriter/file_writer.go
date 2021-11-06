package codewriter

import (
	"bufio"
	"log"
	"os"
	"strings"
)

type fileWriter struct {
	file     *os.File
	fileName string
	writer   *bufio.Writer
}

func newFileWriter(fileName string) *fileWriter {
	asmFile, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Fatalf("can't create file: %s", fileName)
	}

	return &fileWriter{
		file:     asmFile,
		fileName: strings.Split(fileName, ".")[0],
		writer:   bufio.NewWriter(asmFile),
	}
}

func (w *fileWriter) changeFile(fileName string) {
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatalf("can't create file: %s", fileName)
	}

	w.file = file
	w.writer = bufio.NewWriter(file)
}

func (w *fileWriter) writeString(line string) {
	_, err := w.writer.WriteString(line)
	if err != nil {
		log.Fatalf("%v", err)
	}
}

func (w *fileWriter) fileNameWithoutExtension() string {
	return w.fileName
}

func (w *fileWriter) close() {
	err := w.writer.Flush()
	if err != nil {
		log.Fatalf("%v", err)
	}

	w.file.Close()
}
