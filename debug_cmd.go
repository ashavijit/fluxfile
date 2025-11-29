package main

import (
	"fmt"

	"github.com/ashavijit/fluxfile/internal/lexer"
	"github.com/ashavijit/fluxfile/internal/parser"
)

func main() {
	input := `task test:
    run:
        go fmt ./...
`

	l := lexer.New(input)
	p := parser.New(l)

	fluxFile, err := p.Parse()
	if err != nil {
		fmt.Printf("Parse error: %v\n", err)
		return
	}

	if len(fluxFile.Tasks) > 0 {
		fmt.Printf("Commands: %#v\n", fluxFile.Tasks[0].Run)
	}
}
