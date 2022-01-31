package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"

	"github.com/Ladicle/org2html/org"
)

const (
	inputFile  = "fixtures/input.org"
	outputFile = "fixtures/output.html"
)

func main() {
	if err := test(); err != nil {
		log.Fatalln(err)
	}
	log.Println("pass all tests")
}

func test() error {
	f, err := os.Open(inputFile)
	if err != nil {
	}

	tokens, err := org.DefaultTokenizer().Tokenize(f)
	if err != nil {
		return err
	}

	nodes, err := org.DefaultParser(tokens).Parse()
	if err != nil {
		return err
	}

	var out bytes.Buffer
	if err := org.Write(nodes, &out); err != nil {
		return err
	}
	return ioutil.WriteFile(outputFile, out.Bytes(), 0644)
}
