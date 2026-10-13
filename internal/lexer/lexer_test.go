package lexer

import "testing"

func TestAdvance(t *testing.T) {
	lexer := New("abc")
	first := lexer.advance()
	second := lexer.advance()

	if first != 'a' {
		t.Errorf("expected %q, got %q", 'a', first)
	}
	if second != 'b' {
		t.Errorf("expected %q, got %q", 'b', second)
	}
	if lexer.current != 2 {
		t.Errorf("expected current %d, got %d", 2, lexer.current)
	}
	if lexer.start != 0 {
		t.Errorf("expected start 0, got %d", lexer.start)
	}
}

func TestAdvancePosition(t *testing.T) {
	testCases := []struct {
		expectedByte    byte
		expectedCurrent int
		expectedLine    int
		expectedColumn  int
	}{
		{'a', 1, 1, 2},
		{'\n', 2, 2, 1},
		{'b', 3, 2, 2},
		{'\n', 4, 3, 1},
		{'\n', 5, 4, 1},
		{'c', 6, 4, 2},
	}
	lexer := New("a\nb\n\nc")

	if lexer.current != 0 || lexer.line != 1 || lexer.column != 1 {
		t.Errorf("unexpected initial position: current=%d, line=%d, column=%d",
			lexer.current, lexer.line, lexer.column)
	}

	for _, tC := range testCases {
		result := lexer.advance()

		if result != tC.expectedByte {
			t.Errorf("expected byte %q, got %q", tC.expectedByte, result)
		}
		if lexer.current != tC.expectedCurrent {
			t.Errorf("expected current %d, got %d",
				tC.expectedCurrent, lexer.current)
		}
		if lexer.line != tC.expectedLine {
			t.Errorf("expected line %d, got %d",
				tC.expectedLine, lexer.line)
		}
		if lexer.column != tC.expectedColumn {
			t.Errorf("expected column %d, got %d",
				tC.expectedColumn, lexer.column)
		}
	}
}

func TestIsAtEnd(t *testing.T) {
	lexer := New("abc")
	if lexer.isAtEnd() {
		t.Error("expected lexer should not be at end right after initialization.")
	}

	lexer.advance()
	lexer.advance()
	lexer.advance()
	if !lexer.isAtEnd() {
		t.Error("expected isAtEnd to return true when cursor equals source length")
	}

	lexer.current = 4
	if !lexer.isAtEnd() {
		t.Error("expected isAtEnd to return true when cursor exceeds source length")
	}
}

func TestPeek(t *testing.T) {
	lexer := New("abc")
	afterInit := lexer.peek()
	if afterInit != 'a' || lexer.current != 0 {
		t.Errorf("expected return %q and current %d, got %q and %d", 'a', 0, afterInit, lexer.current)
	}

	lexer.advance()
	afterFirst := lexer.peek()
	if afterFirst != 'b' || lexer.current != 1 {
		t.Errorf("expected return %q and current %d, got %q and %d", 'b', 1, afterFirst, lexer.current)
	}

	afterFirstAgain := lexer.peek()
	if afterFirstAgain != 'b' || lexer.current != 1 {
		t.Errorf("expected return %q and current %d, got %q and %d", 'b', 1, afterFirstAgain, lexer.current)
	}

	lexer.advance()
	lexer.advance()
	afterLast := lexer.peek()
	if afterLast != 0 || lexer.current != 3 {
		t.Errorf("expected return %d and current %d, got %d and %d", 0, 3, afterLast, lexer.current)
	}

	if lexer.start != 0 {
		t.Errorf("expected start to remain 0, got %d", lexer.start)
	}
}

func TestIsIdentifierStart(t *testing.T) {
	testCases := []struct {
		input    byte
		expected bool
	}{
		{'a', true},
		{'Z', true},
		{'_', true},
		{'5', false},
		{' ', false},
		{'-', false},
	}
	for _, tC := range testCases {
		t.Run(string(tC.input), func(t *testing.T) {
			result := isIdentifierStart(tC.input)
			if result != tC.expected {
				t.Errorf("input %q: expected %t, got %t", tC.input, tC.expected, result)
			}
		})
	}
}

func TestIsIdentifierContinue(t *testing.T) {
	testCases := []struct {
		input    byte
		expected bool
	}{
		{'a', true},
		{'Z', true},
		{'_', true},
		{'0', true},
		{'9', true},
		{' ', false},
		{'-', false},
		{'.', false},
	}
	for _, tC := range testCases {
		t.Run(string(tC.input), func(t *testing.T) {
			result := isIdentifierContinue(tC.input)
			if result != tC.expected {
				t.Errorf("input %q: expected %t, got %t", tC.input, tC.expected, result)
			}
		})
	}
}

func TestScanIdentifier(t *testing.T) {
	testCases := []struct {
		input           string
		expected        string
		expectedCurrent int
	}{
		{"age", "age", 3},
		{"player_1 be 20", "player_1", 8},
		{"_count42", "_count42", 8},
		{"x", "x", 1},
		{"name-value", "name", 4},
	}
	for _, tC := range testCases {
		t.Run(tC.input, func(t *testing.T) {
			lexer := New(tC.input)
			lexer.advance()
			result := lexer.scanIdentifier()
			if lexer.start != 0 {
				t.Errorf("expected start to remain 0, got %d", lexer.start)
			}
			if result != tC.expected || lexer.current != tC.expectedCurrent {
				t.Errorf("input %q: expected %q with current %d, got %q and %d",
					tC.input, tC.expected, tC.expectedCurrent, result, lexer.current)
			}
		})
	}
}

func TestScanIdentifierKind(t *testing.T) {
	testCases := []struct {
		input    string
		expected TokenKind
	}{
		{"let", TokenLet},
		{"be", TokenBe},
		{"say", TokenSay},
		{"player_1", TokenIdent},
		{"letter", TokenIdent},
		{"Let", TokenIdent},
		{"say hello", TokenSay},
	}
	for _, tC := range testCases {
		t.Run(tC.input, func(t *testing.T) {
			lexer := New(tC.input)
			lexer.advance()
			result := lexer.scanIdentifierKind()
			if result != tC.expected {
				t.Errorf("input %q: expected %q, got %q", tC.input, tC.expected, result)
			}
		})
	}
}

func TestMarkTokenStart(t *testing.T) {
	lexer := New("let age be 20")
	lexer.advance()
	lexer.advance()
	lexer.advance()
	lexer.advance()

	lexer.markTokenStart()
	if lexer.start != 4 {
		t.Errorf("expected start %d, got %d", 4, lexer.start)
	}
	if lexer.current != 4 {
		t.Errorf("expected current %d, got %d", 4, lexer.current)
	}
	if lexer.line != 1 {
		t.Errorf("expected line %d, got %d", 1, lexer.line)
	}
	if lexer.column != 5 {
		t.Errorf("expected column %d, got %d", 5, lexer.column)
	}
	if lexer.startPosition.Line != 1 {
		t.Errorf("expected next start line %d, got %d", 1, lexer.startPosition.Line)
	}
	if lexer.startPosition.Column != 5 {
		t.Errorf("expected next start column %d, got %d", 5, lexer.startPosition.Column)
	}

	lexer.advance()
	if lexer.column != 6 {
		t.Errorf("expected column %d, got %d", 6, lexer.column)
	}
	if lexer.startPosition.Column != 5 {
		t.Errorf("expected next start column %d, got %d", 5, lexer.startPosition.Column)
	}
	if lexer.start != 4 {
		t.Errorf("expected start %d, got %d", 4, lexer.start)
	}
}

func TestMakeToken(t *testing.T) {
	lexer := New("let age")
	lexer.markTokenStart()
	lexer.advance()
	lexer.advance()
	lexer.advance()

	token1 := lexer.makeToken(TokenLet)
	if token1.Kind != TokenLet {
		t.Errorf("expected kind %q, got %q", TokenLet, token1.Kind)
	}
	if token1.Lexeme != "let" {
		t.Errorf("expected lexeme %q, got %q", "let", token1.Lexeme)
	}
	if token1.Position.Line != 1 {
		t.Errorf("expected start line %d, got %d", 1, token1.Position.Line)
	}
	if token1.Position.Column != 1 {
		t.Errorf("expected start column %d, got %d", 1, token1.Position.Column)
	}

	if lexer.current != 3 {
		t.Errorf("expected current %d, got %d", 3, lexer.current)
	}
	if lexer.start != 0 {
		t.Errorf("expected source start %d, got %d", 0, lexer.start)
	}
	if lexer.column != 4 {
		t.Errorf("expected column %d, got %d", 4, lexer.column)
	}

	lexer.advance()
	lexer.markTokenStart()
	lexer.advance()
	lexer.advance()
	lexer.advance()

	token2 := lexer.makeToken(TokenIdent)
	if token2.Kind != TokenIdent {
		t.Errorf("expected kind %q, got %q", TokenIdent, token2.Kind)
	}
	if token2.Lexeme != "age" {
		t.Errorf("expected lexeme %q, got %q", "age", token2.Lexeme)
	}
	if token2.Position.Line != 1 {
		t.Errorf("expected start line %d, got %d", 1, token2.Position.Line)
	}
	if token2.Position.Column != 5 {
		t.Errorf("expected start column %d, got %d", 5, token2.Position.Column)
	}
}

func TestSkipHorizontalWhitespace(t *testing.T) {
	testCases := []struct {
		input           string
		expectedCurrent int
		expectedColumn  int
		expectedNext    byte
	}{
		{"   age", 3, 4, 'a'},
		{"age", 0, 1, 'a'},
		{"  \nage", 2, 3, '\n'},
		{"  \tage", 2, 3, '\t'},
		{"   ", 3, 4, 0},
	}
	for _, tC := range testCases {
		t.Run(tC.input, func(t *testing.T) {
			lexer := New(tC.input)
			lexer.skipHorizontalWhitespace()

			if lexer.current != tC.expectedCurrent {
				t.Errorf("input %q: expected current %d, got %d", tC.input, tC.expectedCurrent, lexer.current)
			}
			if lexer.column != tC.expectedColumn {
				t.Errorf("input %q: expected column %d, got %d", tC.input, tC.expectedColumn, lexer.column)
			}
			if lexer.peek() != tC.expectedNext {
				t.Errorf("input %q: expected next %q, got %q", tC.input, tC.expectedNext, lexer.peek())
			}
			if lexer.line != 1 {
				t.Errorf("input %q: expected line %d, got %d", tC.input, 1, lexer.line)
			}
		})
	}
}

func TestNewlineToken(t *testing.T) {
	lexer := New("\nsay")
	lexer.markTokenStart()
	lexer.advance()

	token := lexer.makeToken(TokenNewline)
	if token.Kind != TokenNewline {
		t.Errorf("expected kind %q, got %q", TokenNewline, token.Kind)
	}
	if token.Lexeme != "\n" {
		t.Errorf("expected lexeme %q, got %q", "\n", token.Lexeme)
	}
	if token.Position.Line != 1 {
		t.Errorf("expected start line %d, got %d", 1, token.Position.Line)
	}
	if token.Position.Column != 1 {
		t.Errorf("expected start column %d, got %d", 1, token.Position.Column)
	}

	if lexer.current != 1 {
		t.Errorf("expected current %d, got %d", 1, lexer.current)
	}
	if lexer.line != 2 {
		t.Errorf("expected line %d, got %d", 2, lexer.line)
	}
	if lexer.column != 1 {
		t.Errorf("expected column %d, got %d", 1, lexer.column)
	}
}

func TestNextToken(t *testing.T) {
	testCases := []struct {
		input          string
		expectedKind   TokenKind
		expectedLexeme string
		expectedLine   int
		expectedColumn int
	}{
		{"", TokenEOF, "", 1, 1},
		{"   ", TokenEOF, "", 1, 4},
		{"  \n", TokenNewline, "\n", 1, 3},
	}
	for _, tC := range testCases {
		t.Run(tC.input, func(t *testing.T) {
			lexer := New(tC.input)
			token, err := lexer.nextToken()

			if err != nil {
				t.Fatalf("unexpected lexical error: %v", err)
			}
			if token.Kind != tC.expectedKind {
				t.Errorf("input %q: expected kind %q, got %q", tC.input, tC.expectedKind, token.Kind)
			}
			if token.Lexeme != tC.expectedLexeme {
				t.Errorf("input %q: expected lexeme %q, got %q", tC.input, tC.expectedLexeme, token.Lexeme)
			}
			if token.Position.Line != tC.expectedLine {
				t.Errorf("input %q: expected start line %d, got %d", tC.input, tC.expectedLine, token.Position.Line)
			}
			if token.Position.Column != tC.expectedColumn {
				t.Errorf("input %q: expected start column %d, got %d", tC.input, tC.expectedColumn, token.Position.Column)
			}
		})
	}
}

func TestNextTokenSequence(t *testing.T) {
	testCases := []struct {
		expectedKind   TokenKind
		expectedLexeme string
		expectedLine   int
		expectedColumn int
	}{
		{TokenLet, "let", 1, 1},
		{TokenIdent, "player_1", 1, 5},
		{TokenBe, "be", 1, 14},
		{TokenNewline, "\n", 1, 16},
		{TokenSay, "say", 2, 1},
		{TokenIdent, "player_1", 2, 5},
		{TokenEOF, "", 2, 13},
	}
	lexer := New("let player_1 be\nsay player_1")
	for _, tC := range testCases {
		token, err := lexer.nextToken()
		if err != nil {
			t.Fatalf("unexpected lexical error: %v", err)
		}
		if token.Kind != tC.expectedKind {
			t.Errorf("expected kind %q, got %q", tC.expectedKind, token.Kind)
		}
		if token.Lexeme != tC.expectedLexeme {
			t.Errorf("expected lexeme %q, got %q", tC.expectedLexeme, token.Lexeme)
		}
		if token.Position.Line != tC.expectedLine {
			t.Errorf("expected start line %d, got %d", tC.expectedLine, token.Position.Line)
		}
		if token.Position.Column != tC.expectedColumn {
			t.Errorf("expected start column %d, got %d", tC.expectedColumn, token.Position.Column)
		}
	}
}

func TestIsDigit(t *testing.T) {
	testCases := []struct {
		input    byte
		expected bool
	}{
		{'0', true},
		{'5', true},
		{'9', true},
		{'a', false},
		{'_', false},
		{' ', false},
		{'-', false},
	}
	for _, tC := range testCases {
		t.Run(string(tC.input), func(t *testing.T) {
			result := isDigit(tC.input)
			if result != tC.expected {
				t.Errorf("input %q: expected %t, got %t", tC.input, tC.expected, result)
			}
		})
	}
}

func TestScanInteger(t *testing.T) {
	testCases := []struct {
		input           string
		expectedLexeme  string
		expectedCurrent int
	}{
		{"123", "123", 3},
		{"0", "0", 1},
		{"007", "007", 3},
		{"42 age", "42", 2},
		{"123abc", "123", 3},
		{"99\n", "99", 2},
	}
	for _, tC := range testCases {
		t.Run(tC.input, func(t *testing.T) {
			lexer := New(tC.input)
			lexer.advance()
			result := lexer.scanInteger()
			if lexer.start != 0 {
				t.Errorf("expected start to remain 0, got %d", lexer.start)
			}
			if result != tC.expectedLexeme {
				t.Errorf("input %q: expected lexeme %q, got %q", tC.input, tC.expectedLexeme, result)
			}
			if lexer.current != tC.expectedCurrent {
				t.Errorf("input %q: expected current %d, got %d", tC.input, tC.expectedCurrent, lexer.current)
			}

			if tC.input == "99\n" && lexer.peek() != '\n' {
				t.Errorf("input %q: expected next character %q, got %q", tC.input, '\n', lexer.peek())
			}
		})
	}
}

func TestNextTokenIntegerSequence(t *testing.T) {
	testCases := []struct {
		expectedKind   TokenKind
		expectedLexeme string
		expectedLine   int
		expectedColumn int
	}{
		{TokenLet, "let", 1, 1},
		{TokenIdent, "age", 1, 5},
		{TokenBe, "be", 1, 9},
		{TokenInt, "20", 1, 12},
		{TokenEOF, "", 1, 14},
	}
	lexer := New("let age be 20")
	for _, tC := range testCases {
		token, err := lexer.nextToken()
		if err != nil {
			t.Fatalf("unexpected lexical error: %v", err)
		}
		if token.Kind != tC.expectedKind {
			t.Errorf("expected kind %q, got %q", tC.expectedKind, token.Kind)
		}
		if token.Lexeme != tC.expectedLexeme {
			t.Errorf("expected lexeme %q, got %q", tC.expectedLexeme, token.Lexeme)
		}
		if token.Position.Line != tC.expectedLine {
			t.Errorf("expected start line %d, got %d", tC.expectedLine, token.Position.Line)
		}
		if token.Position.Column != tC.expectedColumn {
			t.Errorf("expected start column %d, got %d", tC.expectedColumn, token.Position.Column)
		}
	}
}

func TestHasLeadingZero(t *testing.T) {
	testCases := []struct {
		input    string
		expected bool
	}{
		{"0", false},
		{"7", false},
		{"42", false},
		{"1002", false},
		{"007", true},
		{"00", true},
		{"0123", true},
		{"", false},
	}
	for _, tC := range testCases {
		t.Run(tC.input, func(t *testing.T) {
			result := hasLeadingZero(tC.input)
			if result != tC.expected {
				t.Errorf("input %q: expected %t, got %t", tC.input, tC.expected, result)
			}
		})
	}
}

func TestNextTokenUnexpectedCharacter(t *testing.T) {
	lexer := New("let age be @")

	for i := 0; i < 3; i++ {
		_, err := lexer.nextToken()
		if err != nil {
			t.Fatalf("unexpected error before invalid character: %v", err)
		}
	}

	token, err := lexer.nextToken()
	if err == nil {
		t.Fatal("expected lexical error, got nil")
	}

	lexErr, ok := err.(*LexError)
	if !ok {
		t.Fatalf("expected *LexError, got %T", err)
	}

	expectedMessage := "unexpected character"
	if lexErr.Position.Line != 1 {
		t.Errorf("expected error line %d, got %d", 1, lexErr.Position.Line)
	}
	if lexErr.Position.Column != 12 {
		t.Errorf("expected error column %d, got %d", 12, lexErr.Position.Column)
	}
	if lexErr.Message != expectedMessage {
		t.Errorf("expected message %q, got %q", expectedMessage, lexErr.Message)
	}

	if token.Kind != TokenKind("") {
		t.Errorf("expected kind %q, got %q", TokenKind(""), token.Kind)
	}
}

func TestNextTokenLeadingZero(t *testing.T) {
	lexer := New("let age be 007")

	for i := 0; i < 3; i++ {
		_, err := lexer.nextToken()
		if err != nil {
			t.Fatalf("unexpected error before invalid integer literal: %v", err)
		}
	}

	token, err := lexer.nextToken()
	if err == nil {
		t.Fatal("expected lexical error, got nil")
	}

	lexErr, ok := err.(*LexError)
	if !ok {
		t.Fatalf("expected *LexError, got %T", err)
	}

	expectedMessage := "integer literal cannot have leading zeros"
	if lexErr.Position.Line != 1 {
		t.Errorf("expected error line %d, got %d", 1, lexErr.Position.Line)
	}
	if lexErr.Position.Column != 12 {
		t.Errorf("expected error column %d, got %d", 12, lexErr.Position.Column)
	}
	if lexErr.Message != expectedMessage {
		t.Errorf("expected message %q, got %q", expectedMessage, lexErr.Message)
	}

	if token.Kind != TokenKind("") {
		t.Errorf("expected kind %q, got %q", TokenKind(""), token.Kind)
	}
	if lexer.current != 14 {
		t.Errorf("expected current %d, got %d", 14, lexer.current)
	}
}

func TestNextTokenInvalidIntegerSuffix(t *testing.T) {
	testCases := []struct {
		input               string
		expectedLine        int
		expectedColumn      int
		expectedMessage     string
		expectedValidTokens int
	}{
		{"123abc", 1, 1, "invalid integer literal", 0},
		{"42_name", 1, 1, "invalid integer literal", 0},
		{"let age be 123abc", 1, 12, "invalid integer literal", 3},
		{"007abc", 1, 1, "invalid integer literal", 0},
	}
	for _, tC := range testCases {
		t.Run(tC.input, func(t *testing.T) {
			lexer := New(tC.input)

			for i := 0; i < tC.expectedValidTokens; i++ {
				_, err := lexer.nextToken()
				if err != nil {
					t.Fatalf("unexpected error before invalid integer literal: %v", err)
				}
			}

			token, err := lexer.nextToken()
			if err == nil {
				t.Fatal("expected lexical error, got nil")
			}

			lexErr, ok := err.(*LexError)
			if !ok {
				t.Fatalf("expected *LexError, got %T", err)
			}

			if lexErr.Position.Line != tC.expectedLine {
				t.Errorf("expected error line %d, got %d", tC.expectedLine, lexErr.Position.Line)
			}
			if lexErr.Position.Column != tC.expectedColumn {
				t.Errorf("expected error column %d, got %d", tC.expectedColumn, lexErr.Position.Column)
			}
			if lexErr.Message != tC.expectedMessage {
				t.Errorf("expected message %q, got %q", tC.expectedMessage, lexErr.Message)
			}

			if token.Kind != TokenKind("") {
				t.Errorf("expected kind %q, got %q", TokenKind(""), token.Kind)
			}
		})
	}
}

func TestNextTokenValidIntegers(t *testing.T) {
	testCases := []struct {
		input          string
		expectedLexeme string
	}{
		{"0", "0"},
		{"7", "7"},
		{"1002", "1002"},
		{"123 ", "123"},
		{"99\n", "99"},
	}

	for _, tC := range testCases {
		t.Run(tC.input, func(t *testing.T) {
			lexer := New(tC.input)

			token, err := lexer.nextToken()
			if err != nil {
				t.Fatalf("unexpected lexical error: %v", err)
			}

			if token.Kind != TokenInt {
				t.Errorf("expected kind %q, got %q", TokenInt, token.Kind)
			}
			if token.Lexeme != tC.expectedLexeme {
				t.Errorf("expected lexeme %q, got %q", tC.expectedLexeme, token.Lexeme)
			}
			if token.Position.Line != 1 {
				t.Errorf("expected start line %d, got %d", 1, token.Position.Line)
			}
			if token.Position.Column != 1 {
				t.Errorf("expected start column %d, got %d", 1, token.Position.Column)
			}
		})
	}
}
