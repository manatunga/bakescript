package lexer

// TokenKind identifies the lexical category of a token.
// Each category represents a keyword, literal, identifier, or structural token.
type TokenKind string

// Token kinds currently supported by Eloque's lexer.
// Additional kinds will be introduced as the language's syntax is implemented.
const (
	TokenLet     TokenKind = "LET"
	TokenBe      TokenKind = "BE"
	TokenSay     TokenKind = "SAY"
	TokenIdent   TokenKind = "IDENT"
	TokenInt     TokenKind = "INT"
	TokenNewline TokenKind = "NEWLINE"
	TokenEOF     TokenKind = "EOF"
)

// Position identifies a location in Eloque source code.
// Line and column numbers are 1-based.
type Position struct {
	Line   int
	Column int
}

// Token represents a recognized unit of Eloque source code.
// Lexeme preserves the original source text, while Position records
// the location of its first character.
type Token struct {
	Kind     TokenKind
	Lexeme   string
	Position Position
}

// lookupIdentifier determines whether a scanned word is a reserved keyword.
// Words that do not exactly match a recognized keyword are identifiers.
func lookupIdentifier(word string) TokenKind {
	switch word {
	case "let":
		return TokenLet
	case "be":
		return TokenBe
	case "say":
		return TokenSay
	default:
		return TokenIdent
	}
}
