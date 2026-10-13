package lexer

// Lexer maintains the state required to tokenize Eloque source code.
// start marks the beginning of the current token, while current points
// to the next unread byte. Line and column positions are 1-based.
type Lexer struct {
	source        string
	start         int
	current       int
	line          int
	column        int
	startPosition Position
}

// New creates a lexer for the provided source text.
// Scanning begins at byte offset 0, line 1, column 1.
func New(source string) *Lexer {
	return &Lexer{
		source:        source,
		start:         0,
		current:       0,
		line:          1,
		column:        1,
		startPosition: Position{Line: 1, Column: 1},
	}
}

// advance consumes the next source byte and updates the cursor position.
// A newline moves the cursor to column 1 of the following line.
// The caller must ensure the lexer has not reached EOF.
func (l *Lexer) advance() byte {
	currentByte := l.source[l.current]
	if currentByte == '\n' {
		l.line += 1
		l.column = 1
	} else {
		l.column += 1
	}
	l.current += 1
	return currentByte
}

// peek returns the next unread byte without advancing the cursor.
// It returns 0 when the lexer has reached EOF.
func (l *Lexer) peek() byte {
	if l.isAtEnd() {
		return 0
	}
	return l.source[l.current]
}

// isAtEnd reports whether the cursor has reached or passed
// the end of the source text.
func (l *Lexer) isAtEnd() bool {
	return l.current >= len(l.source)
}

// isIdentifierStart reports whether a byte can begin an Eloque identifier.
// Identifiers currently begin with an ASCII letter or underscore.
func isIdentifierStart(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') ||
		(ch >= 'A' && ch <= 'Z') ||
		ch == '_'
}

// isIdentifierContinue reports whether a byte can appear after the
// first character of an identifier, including ASCII digits.
func isIdentifierContinue(ch byte) bool {
	return isIdentifierStart(ch) ||
		(ch >= '0' && ch <= '9')
}

// scanIdentifier consumes the remaining characters of an identifier.
// The first character must already have been consumed by advance.
// It returns the complete identifier text without classifying it.
func (l *Lexer) scanIdentifier() string {
	for isIdentifierContinue(l.peek()) {
		l.advance()
	}
	return l.source[l.start:l.current]
}

// scanIdentifierKind scans an identifier and determines its token kind.
// Reserved keywords receive their dedicated kinds; other words are identifiers.
func (l *Lexer) scanIdentifierKind() TokenKind {
	word := l.scanIdentifier()
	return lookupIdentifier(word)
}

// markTokenStart records the byte offset and source position at which
// the next token begins. The saved position remains unchanged while
// the token is being scanned.
func (l *Lexer) markTokenStart() {
	l.start = l.current
	l.startPosition = Position{
		Line:   l.line,
		Column: l.column,
	}
}

// makeToken constructs a token from the source range [start, current).
// The token retains its original lexeme and saved starting position.
func (l *Lexer) makeToken(kind TokenKind) Token {
	return Token{
		Kind:     kind,
		Lexeme:   l.source[l.start:l.current],
		Position: l.startPosition,
	}
}

// skipHorizontalWhitespace consumes consecutive ASCII space characters.
// It deliberately leaves newlines unread so they can become tokens.
// Tabs are not consumed by this helper.
func (l *Lexer) skipHorizontalWhitespace() {
	for l.peek() == ' ' {
		l.advance()
	}
}

// isDigit reports whether a byte is an ASCII decimal digit.
func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

// scanInteger consumes consecutive decimal digits and returns their
// original source text. The first digit must already be consumed.
// Literal validation is performed separately by nextToken.
func (l *Lexer) scanInteger() string {
	for isDigit(l.peek()) {
		l.advance()
	}
	return l.source[l.start:l.current]
}

// hasLeadingZero reports whether a decimal integer literal contains
// multiple digits and begins with zero. The literal "0" is permitted.
func hasLeadingZero(lexeme string) bool {
	return len(lexeme) > 1 && lexeme[0] == '0'
}

// nextToken scans and returns the next recognized token.
// It skips spaces, preserves newlines, and validates integer literals.
// On invalid input, it returns a zero-value token and a positioned
// lexical error, stopping tokenization at the first detected error.
func (l *Lexer) nextToken() (Token, error) {
	l.skipHorizontalWhitespace()
	l.markTokenStart()

	if l.isAtEnd() {
		return l.makeToken(TokenEOF), nil
	}

	char := l.advance()
	if char == '\n' {
		return l.makeToken(TokenNewline), nil
	}

	if isIdentifierStart(char) {
		kind := l.scanIdentifierKind()
		return l.makeToken(kind), nil
	}
	if isDigit(char) {
		lexeme := l.scanInteger()
		if isIdentifierStart(l.peek()) {
			return Token{}, &LexError{
				Position: l.startPosition,
				Message:  "invalid integer literal",
			}
		}
		if hasLeadingZero(lexeme) {
			return Token{}, &LexError{
				Position: l.startPosition,
				Message:  "integer literal cannot have leading zeros",
			}
		}
		return l.makeToken(TokenInt), nil
	}

	return Token{}, &LexError{
		Position: l.startPosition,
		Message:  "unexpected character",
	}
}
