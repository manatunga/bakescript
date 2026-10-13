package lexer

import "testing"

func TestTokenFields(t *testing.T) {
	position := Position{
		Line:   1,
		Column: 1,
	}
	token := Token{
		Kind:     TokenLet,
		Lexeme:   "let",
		Position: position,
	}

	if token.Kind != TokenLet {
		t.Errorf("expected %q, got %q", TokenLet, token.Kind)
	}
	if token.Lexeme != "let" {
		t.Errorf("expected %q, got %q", "let", token.Lexeme)
	}
	if token.Position.Line != 1 {
		t.Errorf("expected %d, got %d", 1, token.Position.Line)
	}
	if token.Position.Column != 1 {
		t.Errorf("expected %d, got %d", 1, token.Position.Column)
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
