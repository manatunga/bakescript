package lexer

import "testing"

func TestTokenFields(t *testing.T) {
	span := Span{
		Start: Position{Line: 1, Column: 1},
		End:   Position{Line: 1, Column: 4},
	}
	token := Token{
		Kind:   TokenLet,
		Lexeme: "let",
		Span:   span,
	}

	if token.Kind != TokenLet {
		t.Errorf("expected kind %q, got %q", TokenLet, token.Kind)
	}
	if token.Lexeme != "let" {
		t.Errorf("expected lexeme %q, got %q", "let", token.Lexeme)
	}
	if token.Span.Start.Line != 1 {
		t.Errorf("expected start line%d, got %d", 1, token.Span.Start.Line)
	}
	if token.Span.Start.Column != 1 {
		t.Errorf("expected start column %d, got %d", 1, token.Span.Start.Column)
	}
	if token.Span.End.Line != 1 {
		t.Errorf("expected end line %d, got %d", 1, token.Span.End.Line)
	}
	if token.Span.End.Column != 4 {
		t.Errorf("expected end column %d, got %d", 4, token.Span.End.Column)
	}
}

func TestLookupIdentifier(t *testing.T) {
	testCases := []struct {
		input    string
		expected TokenKind
	}{
		{"let", TokenLet},
		{"be", TokenBe},
		{"say", TokenSay},
		{"age", TokenIdent},
		{"letter", TokenIdent},
		{"Let", TokenIdent},
	}

	for _, tC := range testCases {
		t.Run(tC.input, func(t *testing.T) {
			result := lookupIdentifier(tC.input)
			if result != tC.expected {
				t.Errorf("expected %q, got %q", tC.expected, result)
			}
		})
	}
}
