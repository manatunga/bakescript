package parser

import (
	"fmt"
	"testing"

	"github.com/manatunga/eloque/internal/ast"
	"github.com/manatunga/eloque/internal/lexer"
)

// benchmarkProgramResult prevents the benchmark's parsed result
// from being discarded by compiler optimisations.
var benchmarkProgramResult *ast.Program

// makeBenchmarkProgramTokens constructs a token stream containing
// the requested number of statements, cycling through let,
// change, and say statements.
func makeBenchmarkProgramTokens(statementCount int) []lexer.Token {
	tokens := make([]lexer.Token, 0, statementCount*7+1)

	for i := 0; i < statementCount; i++ {
		switch i % 3 {
		case 0:
			tokens = append(tokens,
				lexer.Token{Kind: lexer.TokenLet, Lexeme: "let"},
				lexer.Token{Kind: lexer.TokenIdent, Lexeme: "value"},
				lexer.Token{Kind: lexer.TokenBe, Lexeme: "be"},
				lexer.Token{Kind: lexer.TokenInt, Lexeme: "10"},
				lexer.Token{Kind: lexer.TokenPlus, Lexeme: "+"},
				lexer.Token{Kind: lexer.TokenInt, Lexeme: "20"},
			)
		case 1:
			tokens = append(tokens,
				lexer.Token{Kind: lexer.TokenChange, Lexeme: "change"},
				lexer.Token{Kind: lexer.TokenIdent, Lexeme: "value"},
				lexer.Token{Kind: lexer.TokenTo, Lexeme: "to"},
				lexer.Token{Kind: lexer.TokenIdent, Lexeme: "value"},
				lexer.Token{Kind: lexer.TokenStar, Lexeme: "*"},
				lexer.Token{Kind: lexer.TokenInt, Lexeme: "2"},
			)
		case 2:
			tokens = append(tokens,
				lexer.Token{Kind: lexer.TokenSay, Lexeme: "say"},
				lexer.Token{Kind: lexer.TokenIdent, Lexeme: "value"},
			)
		}

		if i < statementCount-1 {
			tokens = append(tokens, lexer.Token{
				Kind:   lexer.TokenNewline,
				Lexeme: "\n",
			})
		}
	}

	tokens = append(tokens, lexer.Token{Kind: lexer.TokenEOF})
	return tokens
}

// benchmarkParseProgram measures program parsing independently
// of token construction and lexing.
func benchmarkParseProgram(b *testing.B, statementCount int) {
	tokens := makeBenchmarkProgramTokens(statementCount)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		p := New(tokens)

		program, diagnostic := p.parseProgram()
		if diagnostic != nil {
			b.Fatalf("unexpected diagnostic: %+v", diagnostic)
		}
		if program == nil {
			b.Fatal("expected a non-nil program")
		}
		if len(program.Statements) != statementCount {
			b.Fatalf(
				"statement count = %d, want %d",
				len(program.Statements),
				statementCount,
			)
		}

		benchmarkProgramResult = program
	}
}

// BenchmarkParseProgram measures parsing performance across
// small, medium, and large programs.
func BenchmarkParseProgram(b *testing.B) {
	for _, statementCount := range []int{3, 300, 3000} {
		b.Run(fmt.Sprintf("%dStatements", statementCount), func(b *testing.B) {
			benchmarkParseProgram(b, statementCount)
		})
	}
}
