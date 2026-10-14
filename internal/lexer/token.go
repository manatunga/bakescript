package lexer

// TokenKind identifies the lexical category of a token.
// Each category represents a keyword, literal, identifier, or structural token.
type TokenKind uint8

// Token kinds currently supported by Eloque's lexer.
// Each kind has a compact integer representation.
// Additional kinds will be introduced as the language's syntax is implemented.
const (
	TokenInvalid TokenKind = iota
	TokenLet
	TokenBe
	TokenSay
	TokenIdent
	TokenInt
	TokenNewline
	TokenEOF
	TokenIndent
	TokenDedent
)

// String returns the human-readable name of a token kind.
// Unknown values are represented as UNKNOWN.
func (kind TokenKind) String() string {
	switch kind {
	case TokenInvalid:
		return "INVALID"
	case TokenLet:
		return "LET"
	case TokenBe:
		return "BE"
	case TokenSay:
		return "SAY"
	case TokenIdent:
		return "IDENT"
	case TokenInt:
		return "INT"
	case TokenNewline:
		return "NEWLINE"
	case TokenEOF:
		return "EOF"
	case TokenIndent:
		return "INDENT"
	case TokenDedent:
		return "DEDENT"
	default:
		return "UNKNOWN"
	}
}

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
