package parser

import (
	"reflect"
	"testing"

	"github.com/manatunga/eloque/internal/ast"
	"github.com/manatunga/eloque/internal/lexer"
	"github.com/manatunga/eloque/internal/source"
)

func TestParseProgramEmpty(t *testing.T) {
	tests := []struct {
		name       string
		tokens     []lexer.Token
		wantCursor int
	}{
		{
			name: "immediate EOF",
			tokens: []lexer.Token{
				{
					Kind: lexer.TokenEOF,
					Span: source.Span{
						Start: source.Position{Line: 1, Column: 1},
						End:   source.Position{Line: 1, Column: 1},
					},
				},
			},
			wantCursor: 0,
		},
		{
			name: "blank lines only",
			tokens: []lexer.Token{
				{Kind: lexer.TokenNewline, Lexeme: "\n"},
				{Kind: lexer.TokenNewline, Lexeme: "\n"},
				{
					Kind: lexer.TokenEOF,
					Span: source.Span{
						Start: source.Position{Line: 3, Column: 1},
						End:   source.Position{Line: 3, Column: 1},
					},
				},
			},
			wantCursor: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New(tt.tokens)

			program, diagnostic := p.parseProgram()
			if diagnostic != nil {
				t.Fatalf("unexpected diagnostic: %+v", diagnostic)
			}
			if program == nil {
				t.Fatal("expected program, got nil")
			}
			if len(program.Statements) != 0 {
				t.Errorf("statement count = %d, want 0", len(program.Statements))
			}

			eofPosition := tt.tokens[len(tt.tokens)-1].Span.Start
			wantSpan := source.Span{
				Start: eofPosition,
				End:   eofPosition,
			}
			if program.Span() != wantSpan {
				t.Errorf("program span = %+v, want %+v", program.Span(), wantSpan)
			}
			if p.current != tt.wantCursor {
				t.Errorf("cursor = %d, want %d", p.current, tt.wantCursor)
			}
			if p.peek().Kind != lexer.TokenEOF {
				t.Errorf("next token = %v, want EOF", p.peek().Kind)
			}
		})
	}
}

func TestParseProgramStatements(t *testing.T) {
	tests := []struct {
		name       string
		tokens     []lexer.Token
		wantTypes  []ast.Statement
		wantSpan   source.Span
		wantCursor int
	}{
		{
			name: "single statement at EOF",
			tokens: []lexer.Token{
				{
					Kind: lexer.TokenSay, Lexeme: "say",
					Span: source.Span{
						Start: source.Position{Line: 1, Column: 1},
						End:   source.Position{Line: 1, Column: 4},
					},
				},
				{
					Kind: lexer.TokenInt, Lexeme: "42",
					Span: source.Span{
						Start: source.Position{Line: 1, Column: 5},
						End:   source.Position{Line: 1, Column: 7},
					},
				},
				{Kind: lexer.TokenEOF},
			},
			wantTypes: []ast.Statement{&ast.SayStatement{}},
			wantSpan: source.Span{
				Start: source.Position{Line: 1, Column: 1},
				End:   source.Position{Line: 1, Column: 7},
			},
			wantCursor: 2,
		},
		{
			name: "single statement with final newline",
			tokens: []lexer.Token{
				{
					Kind: lexer.TokenSay, Lexeme: "say",
					Span: source.Span{
						Start: source.Position{Line: 1, Column: 1},
						End:   source.Position{Line: 1, Column: 4},
					},
				},
				{
					Kind: lexer.TokenInt, Lexeme: "42",
					Span: source.Span{
						Start: source.Position{Line: 1, Column: 5},
						End:   source.Position{Line: 1, Column: 7},
					},
				},
				{Kind: lexer.TokenNewline, Lexeme: "\n"},
				{Kind: lexer.TokenEOF},
			},
			wantTypes: []ast.Statement{&ast.SayStatement{}},
			wantSpan: source.Span{
				Start: source.Position{Line: 1, Column: 1},
				End:   source.Position{Line: 1, Column: 7},
			},
			wantCursor: 3,
		},
		{
			name: "mixed statements with blank lines",
			tokens: []lexer.Token{
				{Kind: lexer.TokenNewline, Lexeme: "\n"},
				{
					Kind: lexer.TokenLet, Lexeme: "let",
					Span: source.Span{
						Start: source.Position{Line: 2, Column: 1},
						End:   source.Position{Line: 2, Column: 4},
					},
				},
				{Kind: lexer.TokenIdent, Lexeme: "age"},
				{Kind: lexer.TokenBe, Lexeme: "be"},
				{Kind: lexer.TokenInt, Lexeme: "20"},
				{Kind: lexer.TokenNewline, Lexeme: "\n"},
				{Kind: lexer.TokenNewline, Lexeme: "\n"},
				{Kind: lexer.TokenChange, Lexeme: "change"},
				{Kind: lexer.TokenIdent, Lexeme: "age"},
				{Kind: lexer.TokenTo, Lexeme: "to"},
				{Kind: lexer.TokenInt, Lexeme: "21"},
				{Kind: lexer.TokenNewline, Lexeme: "\n"},
				{
					Kind: lexer.TokenSay, Lexeme: "say",
					Span: source.Span{
						Start: source.Position{Line: 5, Column: 1},
						End:   source.Position{Line: 5, Column: 4},
					},
				},
				{
					Kind: lexer.TokenIdent, Lexeme: "age",
					Span: source.Span{
						Start: source.Position{Line: 5, Column: 5},
						End:   source.Position{Line: 5, Column: 8},
					},
				},
				{Kind: lexer.TokenNewline, Lexeme: "\n"},
				{Kind: lexer.TokenNewline, Lexeme: "\n"},
				{Kind: lexer.TokenEOF},
			},
			wantTypes: []ast.Statement{
				&ast.LetStatement{},
				&ast.ChangeStatement{},
				&ast.SayStatement{},
			},
			wantSpan: source.Span{
				Start: source.Position{Line: 2, Column: 1},
				End:   source.Position{Line: 5, Column: 8},
			},
			wantCursor: 16,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New(tt.tokens)

			program, diagnostic := p.parseProgram()
			if diagnostic != nil {
				t.Fatalf("unexpected diagnostic: %+v", diagnostic)
			}
			if program == nil {
				t.Fatal("expected program, got nil")
			}

			if len(program.Statements) != len(tt.wantTypes) {
				t.Fatalf("statement count = %d, want %d",
					len(program.Statements), len(tt.wantTypes))
			}

			for i, want := range tt.wantTypes {
				if reflect.TypeOf(program.Statements[i]) != reflect.TypeOf(want) {
					t.Errorf("statement %d type = %T, want %T",
						i, program.Statements[i], want)
				}
			}

			if program.Span() != tt.wantSpan {
				t.Errorf("program span = %+v, want %+v",
					program.Span(), tt.wantSpan)
			}
			if p.current != tt.wantCursor {
				t.Errorf("cursor = %d, want %d", p.current, tt.wantCursor)
			}
			if p.peek().Kind != lexer.TokenEOF {
				t.Errorf("next token = %v, want EOF", p.peek().Kind)
			}
		})
	}
}

func TestParseProgramGroupedExpressionSpan(t *testing.T) {
	p := New([]lexer.Token{
		{
			Kind: lexer.TokenSay, Lexeme: "say",
			Span: source.Span{
				Start: source.Position{Line: 1, Column: 1},
				End:   source.Position{Line: 1, Column: 4},
			},
		},
		{
			Kind: lexer.TokenLeftParen, Lexeme: "(",
			Span: source.Span{
				Start: source.Position{Line: 1, Column: 5},
				End:   source.Position{Line: 1, Column: 6},
			},
		},
		{
			Kind: lexer.TokenInt, Lexeme: "2",
			Span: source.Span{
				Start: source.Position{Line: 1, Column: 6},
				End:   source.Position{Line: 1, Column: 7},
			},
		},
		{Kind: lexer.TokenPlus, Lexeme: "+"},
		{
			Kind: lexer.TokenInt, Lexeme: "3",
			Span: source.Span{
				Start: source.Position{Line: 1, Column: 10},
				End:   source.Position{Line: 1, Column: 11},
			},
		},
		{
			Kind: lexer.TokenRightParen, Lexeme: ")",
			Span: source.Span{
				Start: source.Position{Line: 1, Column: 11},
				End:   source.Position{Line: 1, Column: 12},
			},
		},
		{Kind: lexer.TokenEOF},
	})

	program, diagnostic := p.parseProgram()
	if diagnostic != nil {
		t.Fatalf("unexpected diagnostic: %+v", diagnostic)
	}

	wantSpan := source.Span{
		Start: source.Position{Line: 1, Column: 1},
		End:   source.Position{Line: 1, Column: 12},
	}

	if program.Span() != wantSpan {
		t.Errorf("program span = %+v, want %+v", program.Span(), wantSpan)
	}
	if len(program.Statements) != 1 {
		t.Fatalf("statement count = %d, want 1", len(program.Statements))
	}
	if program.Statements[0].Span() != wantSpan {
		t.Errorf("statement span = %+v, want %+v",
			program.Statements[0].Span(), wantSpan)
	}
}

func TestParseProgramDiagnostics(t *testing.T) {
	tests := []struct {
		name         string
		tokens       []lexer.Token
		wantMessage  string
		wantExpected []lexer.TokenKind
		wantActual   lexer.TokenKind
		wantLocation source.Span
		wantCursor   int
	}{
		{
			name: "missing separator before say",
			tokens: []lexer.Token{
				{Kind: lexer.TokenLet, Lexeme: "let"},
				{Kind: lexer.TokenIdent, Lexeme: "age"},
				{Kind: lexer.TokenBe, Lexeme: "be"},
				{Kind: lexer.TokenInt, Lexeme: "20"},
				{
					Kind:   lexer.TokenSay,
					Lexeme: "say",
					Span: source.Span{
						Start: source.Position{Line: 1, Column: 15},
						End:   source.Position{Line: 1, Column: 18},
					},
				},
				{Kind: lexer.TokenIdent, Lexeme: "age"},
				{Kind: lexer.TokenEOF},
			},
			wantMessage: "expected a newline or end of file after statement",
			wantExpected: []lexer.TokenKind{
				lexer.TokenNewline,
				lexer.TokenEOF,
			},
			wantActual: lexer.TokenSay,
			wantLocation: source.Span{
				Start: source.Position{Line: 1, Column: 15},
				End:   source.Position{Line: 1, Column: 18},
			},
			wantCursor: 4,
		},
		{
			name: "unsupported statement starter",
			tokens: []lexer.Token{
				{
					Kind:   lexer.TokenIdent,
					Lexeme: "age",
					Span: source.Span{
						Start: source.Position{Line: 2, Column: 3},
						End:   source.Position{Line: 2, Column: 6},
					},
				},
				{Kind: lexer.TokenEOF},
			},
			wantMessage: "expected a statement ('let', 'change' or 'say')",
			wantExpected: []lexer.TokenKind{
				lexer.TokenLet,
				lexer.TokenChange,
				lexer.TokenSay,
			},
			wantActual: lexer.TokenIdent,
			wantLocation: source.Span{
				Start: source.Position{Line: 2, Column: 3},
				End:   source.Position{Line: 2, Column: 6},
			},
			wantCursor: 0,
		},
		{
			name: "missing let identifier",
			tokens: []lexer.Token{
				{Kind: lexer.TokenLet, Lexeme: "let"},
				{
					Kind:   lexer.TokenBe,
					Lexeme: "be",
					Span: source.Span{
						Start: source.Position{Line: 3, Column: 5},
						End:   source.Position{Line: 3, Column: 7},
					},
				},
				{Kind: lexer.TokenInt, Lexeme: "20"},
				{Kind: lexer.TokenEOF},
			},
			wantMessage:  "expected a variable name after 'let'",
			wantExpected: []lexer.TokenKind{lexer.TokenIdent},
			wantActual:   lexer.TokenBe,
			wantLocation: source.Span{
				Start: source.Position{Line: 3, Column: 5},
				End:   source.Position{Line: 3, Column: 7},
			},
			wantCursor: 1,
		},
		{
			name: "missing change to keyword",
			tokens: []lexer.Token{
				{Kind: lexer.TokenChange, Lexeme: "change"},
				{Kind: lexer.TokenIdent, Lexeme: "age"},
				{
					Kind:   lexer.TokenInt,
					Lexeme: "21",
					Span: source.Span{
						Start: source.Position{Line: 4, Column: 12},
						End:   source.Position{Line: 4, Column: 14},
					},
				},
				{Kind: lexer.TokenEOF},
			},
			wantMessage: "expected 'to' after variable name to assign a new variable value",
			wantExpected: []lexer.TokenKind{
				lexer.TokenTo,
			},
			wantActual: lexer.TokenInt,
			wantLocation: source.Span{
				Start: source.Position{Line: 4, Column: 12},
				End:   source.Position{Line: 4, Column: 14},
			},
			wantCursor: 2,
		},
		{
			name: "missing say expression",
			tokens: []lexer.Token{
				{Kind: lexer.TokenSay, Lexeme: "say"},
				{
					Kind:   lexer.TokenNewline,
					Lexeme: "\n",
					Span: source.Span{
						Start: source.Position{Line: 5, Column: 4},
						End:   source.Position{Line: 6, Column: 1},
					},
				},
				{Kind: lexer.TokenEOF},
			},
			wantMessage: "expected an expression",
			wantExpected: []lexer.TokenKind{
				lexer.TokenInt,
				lexer.TokenString,
				lexer.TokenIdent,
				lexer.TokenLeftParen,
				lexer.TokenPlus,
				lexer.TokenMinus,
			},
			wantActual: lexer.TokenNewline,
			wantLocation: source.Span{
				Start: source.Position{Line: 5, Column: 4},
				End:   source.Position{Line: 6, Column: 1},
			},
			wantCursor: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New(tt.tokens)

			program, diagnostic := p.parseProgram()
			if diagnostic == nil {
				t.Fatal("expected diagnostic, got nil")
			}
			if program != nil {
				t.Errorf("program = %+v, want nil", program)
			}

			if diagnostic.Message != tt.wantMessage {
				t.Errorf(
					"Message = %q, want %q",
					diagnostic.Message,
					tt.wantMessage,
				)
			}

			if !reflect.DeepEqual(diagnostic.Expected, tt.wantExpected) {
				t.Errorf(
					"Expected = %v, want %v",
					diagnostic.Expected,
					tt.wantExpected,
				)
			}

			if diagnostic.Actual == nil {
				t.Error("Actual = nil, want a token kind")
			} else if *diagnostic.Actual != tt.wantActual {
				t.Errorf(
					"Actual = %v, want %v",
					*diagnostic.Actual,
					tt.wantActual,
				)
			}

			if diagnostic.Location != tt.wantLocation {
				t.Errorf(
					"Location = %+v, want %+v",
					diagnostic.Location,
					tt.wantLocation,
				)
			}

			if p.current != tt.wantCursor {
				t.Errorf(
					"cursor = %d, want %d",
					p.current,
					tt.wantCursor,
				)
			}

			if p.peek().Kind != tt.wantActual {
				t.Errorf(
					"next token = %v, want %v",
					p.peek().Kind,
					tt.wantActual,
				)
			}
		})
	}
}

func TestParseProgramNestedExpressionDiagnostic(t *testing.T) {
	offendingSpan := source.Span{
		Start: source.Position{Line: 1, Column: 10},
		End:   source.Position{Line: 1, Column: 11},
	}

	tokens := []lexer.Token{
		{Kind: lexer.TokenSay, Lexeme: "say"},
		{Kind: lexer.TokenLeftParen, Lexeme: "("},
		{Kind: lexer.TokenInt, Lexeme: "2"},
		{Kind: lexer.TokenPlus, Lexeme: "+"},
		{
			Kind:   lexer.TokenRightParen,
			Lexeme: ")",
			Span:   offendingSpan,
		},
		{Kind: lexer.TokenEOF},
	}

	p := New(tokens)

	program, diagnostic := p.parseProgram()
	if diagnostic == nil {
		t.Fatal("expected diagnostic, got nil")
	}
	if program != nil {
		t.Errorf("program = %+v, want nil", program)
	}

	if diagnostic.Message != "expected an expression" {
		t.Errorf(
			"Message = %q, want %q",
			diagnostic.Message,
			"expected an expression",
		)
	}

	wantExpected := []lexer.TokenKind{
		lexer.TokenInt,
		lexer.TokenString,
		lexer.TokenIdent,
		lexer.TokenLeftParen,
		lexer.TokenPlus,
		lexer.TokenMinus,
	}

	if !reflect.DeepEqual(diagnostic.Expected, wantExpected) {
		t.Errorf(
			"Expected = %v, want %v",
			diagnostic.Expected,
			wantExpected,
		)
	}

	if diagnostic.Actual == nil ||
		*diagnostic.Actual != lexer.TokenRightParen {
		t.Errorf(
			"Actual = %v, want %v",
			diagnostic.Actual,
			lexer.TokenRightParen,
		)
	}

	if diagnostic.Location != offendingSpan {
		t.Errorf(
			"Location = %+v, want %+v",
			diagnostic.Location,
			offendingSpan,
		)
	}

	if p.current != 4 {
		t.Errorf("cursor = %d, want 4", p.current)
	}
	if p.peek().Kind != lexer.TokenRightParen {
		t.Errorf(
			"next token = %v, want RIGHT_PAREN",
			p.peek().Kind,
		)
	}
}

func TestParseProgramDiagnosticAfterValidStatements(t *testing.T) {
	eofSpan := source.Span{
		Start: source.Position{Line: 3, Column: 14},
		End:   source.Position{Line: 3, Column: 14},
	}

	tokens := []lexer.Token{
		{Kind: lexer.TokenLet, Lexeme: "let"},
		{Kind: lexer.TokenIdent, Lexeme: "age"},
		{Kind: lexer.TokenBe, Lexeme: "be"},
		{Kind: lexer.TokenInt, Lexeme: "20"},
		{Kind: lexer.TokenNewline, Lexeme: "\n"},

		{Kind: lexer.TokenSay, Lexeme: "say"},
		{Kind: lexer.TokenIdent, Lexeme: "age"},
		{Kind: lexer.TokenNewline, Lexeme: "\n"},

		{Kind: lexer.TokenChange, Lexeme: "change"},
		{Kind: lexer.TokenIdent, Lexeme: "age"},
		{Kind: lexer.TokenTo, Lexeme: "to"},
		{Kind: lexer.TokenEOF, Span: eofSpan},
	}

	p := New(tokens)

	program, diagnostic := p.parseProgram()
	if diagnostic == nil {
		t.Fatal("expected diagnostic, got nil")
	}
	if program != nil {
		t.Errorf("program = %+v, want nil", program)
	}

	if diagnostic.Message != "expected an expression" {
		t.Errorf(
			"Message = %q, want %q",
			diagnostic.Message,
			"expected an expression",
		)
	}

	wantExpected := []lexer.TokenKind{
		lexer.TokenInt,
		lexer.TokenString,
		lexer.TokenIdent,
		lexer.TokenLeftParen,
		lexer.TokenPlus,
		lexer.TokenMinus,
	}

	if !reflect.DeepEqual(diagnostic.Expected, wantExpected) {
		t.Errorf(
			"Expected = %v, want %v",
			diagnostic.Expected,
			wantExpected,
		)
	}

	if diagnostic.Actual == nil ||
		*diagnostic.Actual != lexer.TokenEOF {
		t.Errorf(
			"Actual = %v, want EOF",
			diagnostic.Actual,
		)
	}

	if diagnostic.Location != eofSpan {
		t.Errorf(
			"Location = %+v, want %+v",
			diagnostic.Location,
			eofSpan,
		)
	}

	if p.current != 11 {
		t.Errorf("cursor = %d, want 11", p.current)
	}

	if p.peek().Kind != lexer.TokenEOF {
		t.Errorf(
			"next token = %v, want EOF",
			p.peek().Kind,
		)
	}
}
