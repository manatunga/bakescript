package parser

import (
	"testing"

	"github.com/manatunga/eloque/internal/lexer"
)

func TestParserPeek(t *testing.T) {
	tokens := []lexer.Token{
		{Kind: lexer.TokenInt, Lexeme: "42"},
		{Kind: lexer.TokenEOF},
	}

	p := New(tokens)

	if got := p.peek(); got.Kind != lexer.TokenInt {
		t.Errorf("peek() = %v, want INT", got.Kind)
	}

	if p.current != 0 {
		t.Errorf("current = %d, want 0", p.current)
	}
}

func TestParserAdvance(t *testing.T) {
	tokens := []lexer.Token{
		{Kind: lexer.TokenInt, Lexeme: "42"},
		{Kind: lexer.TokenNewline},
		{Kind: lexer.TokenEOF},
	}

	p := New(tokens)

	if got := p.advance(); got.Kind != lexer.TokenInt {
		t.Errorf("first advance() = %v, want INT", got.Kind)
	}

	if got := p.advance(); got.Kind != lexer.TokenNewline {
		t.Errorf("second advance() = %v, want NEWLINE", got.Kind)
	}

	if got := p.advance(); got.Kind != lexer.TokenEOF {
		t.Errorf("third advance() = %v, want EOF", got.Kind)
	}

	if got := p.advance(); got.Kind != lexer.TokenEOF {
		t.Errorf("repeated advance() = %v, want EOF", got.Kind)
	}

	if p.current != 2 {
		t.Errorf("current = %d, want 2", p.current)
	}
}

func TestParserCheck(t *testing.T) {
	p := New([]lexer.Token{
		{Kind: lexer.TokenLet},
		{Kind: lexer.TokenEOF},
	})

	if !p.check(lexer.TokenLet) {
		t.Error("check(LET) = false, want true")
	}

	if p.check(lexer.TokenSay) {
		t.Error("check(SAY) = true, want false")
	}

	if p.current != 0 {
		t.Error("check() unexpectedly consumed a token")
	}
}

func TestParserMatch(t *testing.T) {
	p := New([]lexer.Token{
		{Kind: lexer.TokenPlus},
		{Kind: lexer.TokenInt, Lexeme: "5"},
		{Kind: lexer.TokenEOF},
	})

	if p.match(lexer.TokenMinus, lexer.TokenStar) {
		t.Error("match() succeeded for an absent token kind")
	}

	if p.current != 0 {
		t.Error("failed match() advanced the cursor")
	}

	if !p.match(lexer.TokenMinus, lexer.TokenPlus) {
		t.Error("match() failed to consume PLUS")
	}

	if p.current != 1 {
		t.Errorf("current = %d, want 1", p.current)
	}

	if p.match(lexer.TokenPlus) {
		t.Error("match() consumed an unexpected token")
	}

	if got := p.peek(); got.Kind != lexer.TokenInt {
		t.Errorf("peek() = %v, want INT", got.Kind)
	}
}

func TestParserMatchNoKinds(t *testing.T) {
	p := New([]lexer.Token{
		{Kind: lexer.TokenInt, Lexeme: "42"},
		{Kind: lexer.TokenEOF},
	})

	if p.match() {
		t.Error("match() succeeded without any token kinds")
	}

	if p.current != 0 {
		t.Errorf("current = %d, want 0", p.current)
	}
}

func TestParserMatchEOF(t *testing.T) {
	p := New([]lexer.Token{
		{Kind: lexer.TokenEOF},
	})

	if !p.match(lexer.TokenEOF) {
		t.Fatal("match(EOF) = false, want true")
	}

	if p.current != 0 {
		t.Errorf("current = %d, want 0", p.current)
	}

	if got := p.peek(); got.Kind != lexer.TokenEOF {
		t.Errorf("peek() = %v, want EOF", got.Kind)
	}
}
