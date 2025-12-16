package benchmark

import (
	"testing"

	"github.com/ashavijit/fluxfile/internal/ast"
	"github.com/ashavijit/fluxfile/internal/cache"
	"github.com/ashavijit/fluxfile/internal/executor"
	"github.com/ashavijit/fluxfile/internal/graph"
	"github.com/ashavijit/fluxfile/internal/lexer"
	"github.com/ashavijit/fluxfile/internal/parser"
)

var testFluxFile = `var PROJECT = test

task clean:
    desc: Clean build artifacts
    run:
        echo cleaning

task build:
    desc: Build the project
    deps: clean
    run:
        echo building

task test:
    desc: Run tests
    deps: build
    run:
        echo testing

task deploy:
    desc: Deploy application
    deps: test
    run:
        echo deploying
`

func BenchmarkLexer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		l := lexer.New(testFluxFile)
		l.Tokenize()
	}
}

func BenchmarkParser(b *testing.B) {
	for i := 0; i < b.N; i++ {
		l := lexer.New(testFluxFile)
		p := parser.New(l)
		_, _ = p.Parse()
	}
}

func BenchmarkGraphBuild(b *testing.B) {
	l := lexer.New(testFluxFile)
	p := parser.New(l)
	fluxFile, _ := p.Parse()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = graph.BuildGraph(fluxFile.Tasks)
	}
}

func BenchmarkTopologicalSort(b *testing.B) {
	l := lexer.New(testFluxFile)
	p := parser.New(l)
	fluxFile, _ := p.Parse()
	g, _ := graph.BuildGraph(fluxFile.Tasks)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = g.TopologicalSort()
	}
}

func BenchmarkExecutorCreate(b *testing.B) {
	l := lexer.New(testFluxFile)
	p := parser.New(l)
	fluxFile, _ := p.Parse()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = executor.New(fluxFile, b.TempDir(), true)
	}
}

func BenchmarkCacheHashString(b *testing.B) {
	data := "some test data to hash for benchmarking purposes"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = cache.HashString(data)
	}
}

func BenchmarkCacheNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = cache.New(b.TempDir())
	}
}

func BenchmarkCacheSetGet(b *testing.B) {
	c, _ := cache.New(b.TempDir())
	entry := &cache.CacheEntry{
		TaskName:  "test",
		InputHash: "abc123",
		Success:   true,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = c.Set(entry)
		_, _ = c.Get("test", "abc123")
	}
}

func BenchmarkASTNewTask(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = ast.NewTask("benchmark-task")
	}
}

func BenchmarkASTNewFluxFile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = ast.NewFluxFile()
	}
}

func BenchmarkLexerLargeFile(b *testing.B) {
	largeFile := ""
	for i := 0; i < 100; i++ {
		largeFile += testFluxFile
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l := lexer.New(largeFile)
		l.Tokenize()
	}
}

func BenchmarkParserLargeFile(b *testing.B) {
	largeFile := ""
	for i := 0; i < 50; i++ {
		largeFile += testFluxFile
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l := lexer.New(largeFile)
		p := parser.New(l)
		_, _ = p.Parse()
	}
}
