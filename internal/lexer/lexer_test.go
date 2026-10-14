package lexer

import (
	"reflect"
	"testing"
)

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
		t.Errorf("expected kind %s, got %s", TokenLet, token1.Kind)
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
		t.Errorf("expected kind %s, got %s", TokenIdent, token2.Kind)
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
		t.Errorf("expected kind %s, got %s", TokenNewline, token.Kind)
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
				t.Errorf("input %q: expected kind %s, got %s", tC.input, tC.expectedKind, token.Kind)
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
			t.Errorf("expected kind %s, got %s", tC.expectedKind, token.Kind)
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
			t.Errorf("expected kind %s, got %s", tC.expectedKind, token.Kind)
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

	if token.Kind != TokenInvalid {
		t.Errorf("expected kind %s, got %s", TokenInvalid, token.Kind)
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

	if token.Kind != TokenInvalid {
		t.Errorf("expected kind %s, got %s", TokenInvalid, token.Kind)
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

			if token.Kind != TokenInvalid {
				t.Errorf("expected kind %s, got %s", TokenInvalid, token.Kind)
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
				t.Errorf("expected kind %s, got %s", TokenInt, token.Kind)
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

func TestPopPendingToken(t *testing.T) {
	l := New("")

	first := Token{
		Kind:     TokenDedent,
		Position: Position{Line: 3, Column: 1},
	}
	second := Token{
		Kind:     TokenDedent,
		Position: Position{Line: 3, Column: 1},
	}

	l.pendingTokens = []Token{first, second}

	for i := 0; i < 2; i++ {
		token, ok := l.popPendingToken()
		if !ok {
			t.Fatalf("pop %d: expected a pending token", i+1)
		}
		if token.Kind != TokenDedent {
			t.Errorf("pop %d: expected DEDENT, got %q", i+1, token.Kind)
		}
		if token.Position != first.Position {
			t.Errorf("pop %d: expected position %+v, got %+v",
				i+1, first.Position, token.Position)
		}
	}

	token, ok := l.popPendingToken()
	if ok {
		t.Errorf("expected empty queue, got token %+v", token)
	}
	if token != (Token{}) {
		t.Errorf("expected zero-value token, got %+v", token)
	}
}

func TestPopPendingTokenReusesStorage(t *testing.T) {
	l := New("")

	l.pendingTokens = append(l.pendingTokens, Token{Kind: TokenIndent})
	originalAddress := &l.pendingTokens[0]

	first, ok := l.popPendingToken()
	if !ok || first.Kind != TokenIndent {
		t.Fatalf("first pop = (%v, %v), want INDENT, true", first, ok)
	}

	if len(l.pendingTokens) != 0 || l.pendingIndex != 0 {
		t.Fatalf(
			"drained queue: length = %d, index = %d; want 0, 0",
			len(l.pendingTokens), l.pendingIndex,
		)
	}

	// Append another token and verify that the backing array is reused.
	l.pendingTokens = append(l.pendingTokens, Token{Kind: TokenDedent})

	if &l.pendingTokens[0] != originalAddress {
		t.Fatal("pending-token queue did not reuse its backing array")
	}

	second, ok := l.popPendingToken()
	if !ok || second.Kind != TokenDedent {
		t.Fatalf("second pop = (%v, %v), want DEDENT, true", second, ok)
	}
}

func TestNextTokenPendingTokens(t *testing.T) {
	l := New("say")

	first := Token{
		Kind:     TokenDedent,
		Position: Position{Line: 2, Column: 1},
	}
	second := Token{
		Kind:     TokenIndent,
		Position: Position{Line: 3, Column: 5},
	}

	l.pendingTokens = []Token{first, second}

	expected := []Token{
		first,
		second,
		{
			Kind:     TokenSay,
			Lexeme:   "say",
			Position: Position{Line: 1, Column: 1},
		},
	}

	for i, want := range expected {
		got, err := l.nextToken()
		if err != nil {
			t.Fatalf("token %d: unexpected error: %v", i+1, err)
		}

		if got != want {
			t.Errorf("token %d: expected %+v, got %+v", i+1, want, got)
		}

		if i < 2 && l.current != 0 {
			t.Errorf(
				"token %d: expected source cursor to remain at 0, got %d",
				i+1,
				l.current,
			)
		}
	}

	if len(l.pendingTokens) != 0 {
		t.Errorf("expected empty pending-token queue, got %d tokens", len(l.pendingTokens))
	}
}

func TestMeasureIndentation(t *testing.T) {
	testCases := []struct {
		name          string
		source        string
		expectedWidth int
		expectedNext  byte
	}{
		{"no indentation", "say age", 0, 's'},
		{"one level", "    say age", 4, 's'},
		{"two levels", "        say age", 8, 's'},
		{"blank line", "    \n", 4, '\n'},
		{"comment line", "    # comment", 4, '#'},
		{"tab after spaces", "    \tsay", 4, '\t'},
		{"empty source", "", 0, 0},
		{"spaces only", "    ", 4, 0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			l := New(tc.source)

			width := l.measureIndentation()

			if width != tc.expectedWidth {
				t.Errorf(
					"expected indentation width %d, got %d",
					tc.expectedWidth,
					width,
				)
			}

			if l.current != tc.expectedWidth {
				t.Errorf(
					"expected cursor offset %d, got %d",
					tc.expectedWidth,
					l.current,
				)
			}

			if l.column != tc.expectedWidth+1 {
				t.Errorf(
					"expected column %d, got %d",
					tc.expectedWidth+1,
					l.column,
				)
			}

			if next := l.peek(); next != tc.expectedNext {
				t.Errorf(
					"expected next byte %q, got %q",
					tc.expectedNext,
					next,
				)
			}

			if len(l.indentStack) != 1 || l.indentStack[0] != 0 {
				t.Errorf(
					"expected unchanged indentation stack [0], got %v",
					l.indentStack,
				)
			}

			if !l.atLineStart {
				t.Error("expected atLineStart to remain true")
			}
		})
	}
}

func TestIsValidIndentationWidth(t *testing.T) {
	testCases := []struct {
		name     string
		width    int
		expected bool
	}{
		{"zero indentation", 0, true},
		{"one level", 4, true},
		{"two levels", 8, true},
		{"three levels", 12, true},
		{"two spaces", 2, false},
		{"six spaces", 6, false},
		{"ten spaces", 10, false},
		{"negative multiple", -4, false},
		{"negative non-multiple", -2, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := isValidIndentationWidth(tc.width)

			if got != tc.expected {
				t.Errorf(
					"isValidIndentationWidth(%d): expected %t, got %t",
					tc.width,
					tc.expected,
					got,
				)
			}
		})
	}
}
func TestProcessIndentation(t *testing.T) {
	tests := []struct {
		name           string
		initialStack   []int
		width          int
		expectedStack  []int
		expectedTokens []Token
		wantError      bool
	}{
		{
			name:          "increase from base level",
			initialStack:  []int{0},
			width:         4,
			expectedStack: []int{0, 4},
			expectedTokens: []Token{
				{Kind: TokenIndent, Position: Position{Line: 2, Column: 5}},
			},
		},
		{
			name:          "increase from existing level",
			initialStack:  []int{0, 4},
			width:         8,
			expectedStack: []int{0, 4, 8},
			expectedTokens: []Token{
				{Kind: TokenIndent, Position: Position{Line: 2, Column: 5}},
			},
		},
		{
			name:          "unchanged indentation",
			initialStack:  []int{0, 4},
			width:         4,
			expectedStack: []int{0, 4},
		},
		{
			name:          "invalid indentation width",
			initialStack:  []int{0, 4},
			width:         6,
			expectedStack: []int{0, 4},
			wantError:     true,
		},
		{
			name:          "negative indentation width",
			initialStack:  []int{0},
			width:         -4,
			expectedStack: []int{0},
			wantError:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New("")
			l.line = 2
			l.column = 5
			l.indentStack = append([]int(nil), tt.initialStack...)

			err := l.processIndentation(tt.width)

			if (err != nil) != tt.wantError {
				t.Fatalf("processIndentation(%d) error = %v, wantError = %t",
					tt.width, err, tt.wantError)
			}

			if !reflect.DeepEqual(l.indentStack, tt.expectedStack) {
				t.Errorf("indentStack = %v, want %v",
					l.indentStack, tt.expectedStack)
			}

			if len(l.pendingTokens) != len(tt.expectedTokens) {
				t.Fatalf("pendingTokens = %v, want %v",
					l.pendingTokens, tt.expectedTokens)
			}

			for i, token := range tt.expectedTokens {
				if l.pendingTokens[i] != token {
					t.Errorf("pendingTokens[%d] = %+v, want %+v",
						i, l.pendingTokens[i], token)
				}
			}
		})
	}
}

func TestProcessIndentationDedent(t *testing.T) {
	tests := []struct {
		name          string
		initialStack  []int
		width         int
		expectedStack []int
		dedentCount   int
		wantError     bool
	}{
		{
			name:          "single dedent",
			initialStack:  []int{0, 4, 8},
			width:         4,
			expectedStack: []int{0, 4},
			dedentCount:   1,
		},
		{
			name:          "multiple dedents",
			initialStack:  []int{0, 4, 8, 12},
			width:         4,
			expectedStack: []int{0, 4},
			dedentCount:   2,
		},
		{
			name:          "dedent to base",
			initialStack:  []int{0, 4, 8},
			width:         0,
			expectedStack: []int{0},
			dedentCount:   2,
		},
		{
			name:          "nonexistent indentation level",
			initialStack:  []int{0, 4, 12},
			width:         8,
			expectedStack: []int{0, 4, 12},
			wantError:     true,
		},
		{
			name:          "invalid indentation alignment",
			initialStack:  []int{0, 4, 8},
			width:         6,
			expectedStack: []int{0, 4, 8},
			wantError:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New("")
			l.line = 3
			l.column = tt.width + 1
			l.indentStack = append([]int(nil), tt.initialStack...)

			err := l.processIndentation(tt.width)

			if (err != nil) != tt.wantError {
				t.Fatalf("error = %v, wantError = %t", err, tt.wantError)
			}

			if !reflect.DeepEqual(l.indentStack, tt.expectedStack) {
				t.Errorf("indentStack = %v, want %v",
					l.indentStack, tt.expectedStack)
			}

			if len(l.pendingTokens) != tt.dedentCount {
				t.Fatalf("pending token count = %d, want %d",
					len(l.pendingTokens), tt.dedentCount)
			}

			for i, token := range l.pendingTokens {
				if token.Kind != TokenDedent {
					t.Errorf("token %d kind = %s, want DEDENT", i, token.Kind)
				}

				expectedPosition := Position{
					Line:   3,
					Column: tt.width + 1,
				}

				if token.Position != expectedPosition {
					t.Errorf("token %d position = %+v, want %+v",
						i, token.Position, expectedPosition)
				}

				if token.Lexeme != "" {
					t.Errorf("token %d lexeme = %q, want empty", i, token.Lexeme)
				}
			}
		})
	}
}

func TestNextTokenIndentation(t *testing.T) {
	source := "say 1\n    say 2\n        say 3\n    say 4\n"

	l := New(source)

	expected := []Token{
		{Kind: TokenSay, Lexeme: "say", Position: Position{1, 1}},
		{Kind: TokenInt, Lexeme: "1", Position: Position{1, 5}},
		{Kind: TokenNewline, Lexeme: "\n", Position: Position{1, 6}},

		{Kind: TokenIndent, Position: Position{2, 5}},
		{Kind: TokenSay, Lexeme: "say", Position: Position{2, 5}},
		{Kind: TokenInt, Lexeme: "2", Position: Position{2, 9}},
		{Kind: TokenNewline, Lexeme: "\n", Position: Position{2, 10}},

		{Kind: TokenIndent, Position: Position{3, 9}},
		{Kind: TokenSay, Lexeme: "say", Position: Position{3, 9}},
		{Kind: TokenInt, Lexeme: "3", Position: Position{3, 13}},
		{Kind: TokenNewline, Lexeme: "\n", Position: Position{3, 14}},

		{Kind: TokenDedent, Position: Position{4, 5}},
		{Kind: TokenSay, Lexeme: "say", Position: Position{4, 5}},
		{Kind: TokenInt, Lexeme: "4", Position: Position{4, 9}},
		{Kind: TokenNewline, Lexeme: "\n", Position: Position{4, 10}},
	}

	for i, want := range expected {
		got, err := l.nextToken()
		if err != nil {
			t.Fatalf("token %d: unexpected error: %v", i, err)
		}
		if got != want {
			t.Errorf("token %d: got %+v, want %+v", i, got, want)
		}
	}
}

func TestNextTokenBlankLinesPreserveIndentation(t *testing.T) {
	source := "say 1\n    say 2\n        \n\n    say 3\n"

	l := New(source)

	expectedKinds := []TokenKind{
		TokenSay, TokenInt, TokenNewline,
		TokenIndent, TokenSay, TokenInt, TokenNewline,
		TokenNewline,
		TokenNewline,
		TokenSay, TokenInt, TokenNewline,
	}

	for i, want := range expectedKinds {
		got, err := l.nextToken()
		if err != nil {
			t.Fatalf("token %d: unexpected error: %v", i, err)
		}

		if got.Kind != want {
			t.Fatalf("token %d: got %s, want %s", i, got.Kind, want)
		}
	}

	expectedStack := []int{0, 4}
	if !reflect.DeepEqual(l.indentStack, expectedStack) {
		t.Errorf("indentStack = %v, want %v", l.indentStack, expectedStack)
	}

	if len(l.pendingTokens) != 0 {
		t.Errorf("expected empty pending-token queue, got %v", l.pendingTokens)
	}
}

func TestNextTokenEOFDedentation(t *testing.T) {
	tests := []struct {
		name          string
		source        string
		expectedKinds []TokenKind
	}{
		{
			name:   "single indentation level without final newline",
			source: "say 1\n    say 2",
			expectedKinds: []TokenKind{
				TokenSay, TokenInt, TokenNewline,
				TokenIndent, TokenSay, TokenInt,
				TokenDedent, TokenEOF,
			},
		},
		{
			name:   "multiple indentation levels without final newline",
			source: "say 1\n    say 2\n        say 3",
			expectedKinds: []TokenKind{
				TokenSay, TokenInt, TokenNewline,
				TokenIndent, TokenSay, TokenInt, TokenNewline,
				TokenIndent, TokenSay, TokenInt,
				TokenDedent, TokenDedent, TokenEOF,
			},
		},
		{
			name:   "multiple indentation levels with final newline",
			source: "say 1\n    say 2\n        say 3\n",
			expectedKinds: []TokenKind{
				TokenSay, TokenInt, TokenNewline,
				TokenIndent, TokenSay, TokenInt, TokenNewline,
				TokenIndent, TokenSay, TokenInt, TokenNewline,
				TokenDedent, TokenDedent, TokenEOF,
			},
		},
		{
			name:   "trailing blank lines",
			source: "say 1\n    say 2\n    \n\n",
			expectedKinds: []TokenKind{
				TokenSay, TokenInt, TokenNewline,
				TokenIndent, TokenSay, TokenInt, TokenNewline,
				TokenNewline, TokenNewline,
				TokenDedent, TokenEOF,
			},
		},
		{
			name:   "no indentation",
			source: "say 1",
			expectedKinds: []TokenKind{
				TokenSay, TokenInt, TokenEOF,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(tt.source)

			for i, want := range tt.expectedKinds {
				got, err := l.nextToken()
				if err != nil {
					t.Fatalf("token %d: unexpected error: %v", i, err)
				}

				if got.Kind != want {
					t.Fatalf("token %d: got %s, want %s", i, got.Kind, want)
				}
			}

			if len(l.indentStack) != 1 || l.indentStack[0] != 0 {
				t.Errorf("indentStack = %v, want [0]", l.indentStack)
			}

			if len(l.pendingTokens) != 0 {
				t.Errorf("pendingTokens = %v, want empty queue", l.pendingTokens)
			}
		})
	}
}

func TestNextTokenRejectsTabs(t *testing.T) {
	tests := []struct {
		name          string
		source        string
		expectedError Position
	}{
		{
			name:          "tab at beginning of file",
			source:        "\tsay 1",
			expectedError: Position{Line: 1, Column: 1},
		},
		{
			name:          "tab after leading spaces",
			source:        "    \tsay 1",
			expectedError: Position{Line: 1, Column: 5},
		},
		{
			name:          "tab between tokens",
			source:        "say\t1",
			expectedError: Position{Line: 1, Column: 4},
		},
		{
			name:          "tab on subsequent line",
			source:        "say 1\n\t say 2",
			expectedError: Position{Line: 2, Column: 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(tt.source)

			for {
				token, err := l.nextToken()
				if err != nil {
					lexErr, ok := err.(*LexError)
					if !ok {
						t.Fatalf("expected *LexError, got %T: %v", err, err)
					}

					if lexErr.Position != tt.expectedError {
						t.Errorf("error position = %+v, want %+v",
							lexErr.Position, tt.expectedError)
					}
					return
				}

				if token.Kind == TokenEOF {
					t.Fatal("expected lexical error for tab, got EOF")
				}
			}
		})
	}
}

func TestSkipComment(t *testing.T) {
	tests := []struct {
		name           string
		source         string
		expectedCursor int
		expectedLine   int
		expectedColumn int
		expectedPeek   byte
	}{
		{
			name:           "comment followed by newline",
			source:         "# hello world\nsay 42",
			expectedCursor: 13,
			expectedLine:   1,
			expectedColumn: 14,
			expectedPeek:   '\n',
		},
		{
			name:           "comment ending at EOF",
			source:         "# hello world",
			expectedCursor: 13,
			expectedLine:   1,
			expectedColumn: 14,
			expectedPeek:   0,
		},
		{
			name:           "empty comment followed by newline",
			source:         "#\nsay 42",
			expectedCursor: 1,
			expectedLine:   1,
			expectedColumn: 2,
			expectedPeek:   '\n',
		},
		{
			name:           "empty comment at EOF",
			source:         "#",
			expectedCursor: 1,
			expectedLine:   1,
			expectedColumn: 2,
			expectedPeek:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(tt.source)

			l.skipComment()

			if l.current != tt.expectedCursor {
				t.Errorf("cursor = %d, want %d",
					l.current, tt.expectedCursor)
			}

			if l.line != tt.expectedLine {
				t.Errorf("line = %d, want %d",
					l.line, tt.expectedLine)
			}

			if l.column != tt.expectedColumn {
				t.Errorf("column = %d, want %d",
					l.column, tt.expectedColumn)
			}

			if got := l.peek(); got != tt.expectedPeek {
				t.Errorf("peek = %q, want %q",
					got, tt.expectedPeek)
			}

			if !l.atLineStart {
				t.Error("skipComment unexpectedly changed atLineStart")
			}
		})
	}
}

func TestNextTokenComments(t *testing.T) {
	tests := []struct {
		name          string
		source        string
		expectedKinds []TokenKind
	}{
		{
			name:   "inline comment",
			source: "say 42 # explanation\nsay 10",
			expectedKinds: []TokenKind{
				TokenSay, TokenInt, TokenNewline,
				TokenSay, TokenInt, TokenEOF,
			},
		},
		{
			name:   "comment-only line at deeper indentation",
			source: "say 1\n    say 2\n        # comment\n    say 3",
			expectedKinds: []TokenKind{
				TokenSay, TokenInt, TokenNewline,
				TokenIndent, TokenSay, TokenInt, TokenNewline,
				TokenNewline,
				TokenSay, TokenInt,
				TokenDedent, TokenEOF,
			},
		},
		{
			name:   "comment-only line at lower indentation",
			source: "say 1\n    say 2\n# comment\n    say 3",
			expectedKinds: []TokenKind{
				TokenSay, TokenInt, TokenNewline,
				TokenIndent, TokenSay, TokenInt, TokenNewline,
				TokenNewline,
				TokenSay, TokenInt,
				TokenDedent, TokenEOF,
			},
		},
		{
			name:   "comment ending at EOF",
			source: "say 42 # final comment",
			expectedKinds: []TokenKind{
				TokenSay, TokenInt, TokenEOF,
			},
		},
		{
			name:   "comment-only file",
			source: "# comment",
			expectedKinds: []TokenKind{
				TokenEOF,
			},
		},
		{
			name:   "multiple comment lines",
			source: "# first\n# second\nsay 1",
			expectedKinds: []TokenKind{
				TokenNewline, TokenNewline,
				TokenSay, TokenInt, TokenEOF,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(tt.source)

			for i, want := range tt.expectedKinds {
				got, err := l.nextToken()
				if err != nil {
					t.Fatalf("token %d: unexpected error: %v", i, err)
				}

				if got.Kind != want {
					t.Fatalf("token %d: got %s, want %s", i, got.Kind, want)
				}
			}

			if len(l.indentStack) != 1 || l.indentStack[0] != 0 {
				t.Errorf("indentStack = %v, want [0]", l.indentStack)
			}

			if len(l.pendingTokens) != 0 {
				t.Errorf("pendingTokens = %v, want empty queue", l.pendingTokens)
			}
		})
	}
}

func TestTokenKindString(t *testing.T) {
	tests := []struct {
		kind TokenKind
		want string
	}{
		{TokenInvalid, "INVALID"},
		{TokenLet, "LET"},
		{TokenBe, "BE"},
		{TokenSay, "SAY"},
		{TokenIdent, "IDENT"},
		{TokenInt, "INT"},
		{TokenNewline, "NEWLINE"},
		{TokenEOF, "EOF"},
		{TokenIndent, "INDENT"},
		{TokenDedent, "DEDENT"},
		{TokenKind(255), "UNKNOWN"},
	}

	for _, tt := range tests {
		if got := tt.kind.String(); got != tt.want {
			t.Errorf("TokenKind(%d).String() = %q, want %q",
				tt.kind, got, tt.want)
		}
	}
}
