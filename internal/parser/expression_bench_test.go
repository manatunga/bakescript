package parser

import (
	"fmt"
	"testing"

	"github.com/manatunga/eloque/internal/lexer"
	"github.com/manatunga/eloque/internal/source"
)

// benchmarkToken creates a token with a valid synthetic source span.
func benchmarkToken(kind lexer.TokenKind, lexeme string, column int) lexer.Token {
	width := len(lexeme)
	if width == 0 {
		width = 1
	}

	return lexer.Token{
		Kind:   kind,
		Lexeme: lexeme,
		Span: source.Span{
			Start: source.Position{Line: 1, Column: column},
			End:   source.Position{Line: 1, Column: column + width},
		},
	}
}

// benchmarkExpression measures parsing a prepared token stream.
func benchmarkExpression(b *testing.B, tokens []lexer.Token) {
	b.Helper()
	b.ReportAllocs()

	for b.Loop() {
		p := New(tokens)

		_, diagnostic := p.parseExpression(0)
		if diagnostic != nil {
			b.Fatalf("unexpected diagnostic: %+v", diagnostic)
		}

		if p.peek().Kind != lexer.TokenEOF {
			b.Fatalf("unexpected remaining token: %v", p.peek().Kind)
		}
	}
}

// benchmarkChainTokens builds an expression containing repeated binary operators.
func benchmarkChainTokens(operators int, kind lexer.TokenKind, lexeme string) []lexer.Token {
	tokens := make([]lexer.Token, 0, operators*2+2)
	column := 1

	for i := 0; i < operators; i++ {
		tokens = append(tokens, benchmarkToken(lexer.TokenInt, "2", column))
		column++

		tokens = append(tokens, benchmarkToken(kind, lexeme, column))
		column += len(lexeme)
	}

	tokens = append(tokens, benchmarkToken(lexer.TokenInt, "2", column))
	tokens = append(tokens, benchmarkToken(lexer.TokenEOF, "", column+1))

	return tokens
}

// benchmarkGroupingTokens builds an expression with nested parentheses.
func benchmarkGroupingTokens(depth int) []lexer.Token {
	tokens := make([]lexer.Token, 0, depth*2+2)
	column := 1

	for i := 0; i < depth; i++ {
		tokens = append(tokens, benchmarkToken(lexer.TokenLeftParen, "(", column))
		column++
	}

	tokens = append(tokens, benchmarkToken(lexer.TokenInt, "42", column))
	column += 2

	for i := 0; i < depth; i++ {
		tokens = append(tokens, benchmarkToken(lexer.TokenRightParen, ")", column))
		column++
	}

	tokens = append(tokens, benchmarkToken(lexer.TokenEOF, "", column))

	return tokens
}

// BenchmarkParseSimpleArithmetic measures a small arithmetic expression.
func BenchmarkParseSimpleArithmetic(b *testing.B) {
	tokens := benchmarkChainTokens(2, lexer.TokenPlus, "+")
	benchmarkExpression(b, tokens)
}

// BenchmarkParseMixedPrecedence measures arithmetic with multiple precedence levels.
func BenchmarkParseMixedPrecedence(b *testing.B) {
	tokens := []lexer.Token{
		benchmarkToken(lexer.TokenInt, "2", 1),
		benchmarkToken(lexer.TokenPlus, "+", 2),
		benchmarkToken(lexer.TokenInt, "3", 3),
		benchmarkToken(lexer.TokenStar, "*", 4),
		benchmarkToken(lexer.TokenInt, "4", 5),
		benchmarkToken(lexer.TokenPower, "**", 6),
		benchmarkToken(lexer.TokenInt, "5", 8),
		benchmarkToken(lexer.TokenEOF, "", 9),
	}

	benchmarkExpression(b, tokens)
}

// BenchmarkParseAdditiveChain measures left-associative expression scaling.
func BenchmarkParseAdditiveChain(b *testing.B) {
	for _, operators := range []int{10, 100, 1000} {
		b.Run(fmt.Sprintf("Operators%d", operators), func(b *testing.B) {
			tokens := benchmarkChainTokens(operators, lexer.TokenPlus, "+")
			benchmarkExpression(b, tokens)
		})
	}
}

// BenchmarkParseExponentiationChain measures right-associative expression scaling.
func BenchmarkParseExponentiationChain(b *testing.B) {
	for _, operators := range []int{10, 100, 1000} {
		b.Run(fmt.Sprintf("Operators%d", operators), func(b *testing.B) {
			tokens := benchmarkChainTokens(operators, lexer.TokenPower, "**")
			benchmarkExpression(b, tokens)
		})
	}
}

// BenchmarkParseNestedParentheses measures recursive grouping overhead.
func BenchmarkParseNestedParentheses(b *testing.B) {
	for _, depth := range []int{10, 100, 1000} {
		b.Run(fmt.Sprintf("Depth%d", depth), func(b *testing.B) {
			tokens := benchmarkGroupingTokens(depth)
			benchmarkExpression(b, tokens)
		})
	}
}
